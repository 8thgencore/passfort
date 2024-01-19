package user

import (
	"context"

	"github.com/8thgencore/passfort/internal/domain"
	"github.com/8thgencore/passfort/internal/service/adapters/cache"
	"github.com/8thgencore/passfort/internal/service/adapters/storage"
	"github.com/8thgencore/passfort/pkg/util"
)

/**
 * UserService implements port.UserService interface
 * and provides an access to the user repository
 * and cache service
 */
type UserService struct {
	storage storage.UserRepository
	cache   cache.CacheRepository
}

// NewUserService creates a new user service instance
func NewUserService(storage storage.UserRepository, cache cache.CacheRepository) *UserService {
	return &UserService{
		storage,
		cache,
	}
}

// Register creates a new user
func (us *UserService) Register(ctx context.Context, user *domain.User) (*domain.User, error) {
	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = hashedPassword

	_, err = us.storage.CreateUser(ctx, user)
	if err != nil {
		if domain.IsUniqueConstraintViolationError(err) {
			return nil, domain.ErrConflictingData
		}

		return nil, err
	}

	cacheKey := util.GenerateCacheKey("user", user.ID)
	userSerialized, err := util.Serialize(user)
	if err != nil {
		return nil, err
	}

	err = us.cache.Set(ctx, cacheKey, userSerialized, 0)
	if err != nil {
		return nil, err
	}

	err = us.cache.DeleteByPrefix(ctx, "users:*")
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUser gets a user by ID
func (us *UserService) GetUser(ctx context.Context, id uint64) (*domain.User, error) {
	var user *domain.User

	cacheKey := util.GenerateCacheKey("user", id)
	cachedUser, err := us.cache.Get(ctx, cacheKey)
	if err == nil {
		err := util.Deserialize(cachedUser, &user)
		if err != nil {
			return nil, err
		}

		return user, nil
	}

	user, err = us.storage.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	userSerialized, err := util.Serialize(user)
	if err != nil {
		return nil, err
	}

	err = us.cache.Set(ctx, cacheKey, userSerialized, 0)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// ListUsers lists all users
func (us *UserService) ListUsers(ctx context.Context, skip, limit uint64) ([]domain.User, error) {
	var users []domain.User

	params := util.GenerateCacheKeyParams(skip, limit)
	cacheKey := util.GenerateCacheKey("users", params)

	cachedUsers, err := us.cache.Get(ctx, cacheKey)
	if err == nil {
		err := util.Deserialize(cachedUsers, &users)
		if err != nil {
			return nil, err
		}

		return users, nil
	}

	users, err = us.storage.ListUsers(ctx, skip, limit)
	if err != nil {
		return nil, err
	}

	usersSerialized, err := util.Serialize(users)
	if err != nil {
		return nil, err
	}

	err = us.cache.Set(ctx, cacheKey, usersSerialized, 0)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// UpdateUser updates a user's name, email, and password
func (us *UserService) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	existingUser, err := us.storage.GetUserByID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	emptyData := user.Name == "" &&
		user.Email == "" &&
		user.Password == "" &&
		user.Role == ""
	sameData := existingUser.Name == user.Name &&
		existingUser.Email == user.Email &&
		existingUser.Role == user.Role
	if emptyData || sameData {
		return nil, domain.ErrNoUpdatedData
	}

	var hashedPassword string

	if user.Password != "" {
		hashedPassword, err = util.HashPassword(user.Password)
		if err != nil {
			return nil, err
		}
	}

	user.Password = hashedPassword

	_, err = us.storage.UpdateUser(ctx, user)
	if err != nil {
		if domain.IsUniqueConstraintViolationError(err) {
			return nil, domain.ErrConflictingData
		}

		return nil, err
	}

	cacheKey := util.GenerateCacheKey("user", user.ID)
	_ = us.cache.Delete(ctx, cacheKey)

	userSerialized, err := util.Serialize(user)
	if err != nil {
		return nil, err
	}

	err = us.cache.Set(ctx, cacheKey, userSerialized, 0)
	if err != nil {
		return nil, err
	}

	err = us.cache.DeleteByPrefix(ctx, "users:*")
	if err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser deletes a user by ID
func (us *UserService) DeleteUser(ctx context.Context, id uint64) error {
	_, err := us.storage.GetUserByID(ctx, id)
	if err != nil {
		return err
	}

	cacheKey := util.GenerateCacheKey("user", id)
	_ = us.cache.Delete(ctx, cacheKey)

	err = us.cache.DeleteByPrefix(ctx, "users:*")
	if err != nil {
		return err
	}

	return us.storage.DeleteUser(ctx, id)
}
