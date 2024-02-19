package user

import (
	"context"

	"github.com/8thgencore/passfort/internal/domain"
	"github.com/8thgencore/passfort/internal/repository/storage/postgres/converter"
	"github.com/8thgencore/passfort/pkg/util"
	"github.com/google/uuid"
)

// Register creates a new user
func (us *UserService) Register(ctx context.Context, user *domain.User) (*domain.User, error) {
	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		return nil, domain.ErrInternal
	}

	user.Password = hashedPassword

	userDAO, err := us.storage.CreateUser(ctx, converter.ToUserDAO(user))
	if err != nil {
		if err == domain.ErrConflictingData {
			return nil, err
		}
		us.log.Error("failed to create a user", "error", err.Error())

		return nil, domain.ErrInternal
	}
	user = converter.ToUser(userDAO)

	cacheKey := util.GenerateCacheKey("user", user.ID)
	userSerialized, err := util.Serialize(user)
	if err != nil {
		return nil, domain.ErrInternal
	}

	err = us.cache.Set(ctx, cacheKey, userSerialized, 0)
	if err != nil {
		return nil, domain.ErrInternal
	}

	err = us.cache.DeleteByPrefix(ctx, "users:*")
	if err != nil {
		return nil, domain.ErrInternal
	}

	return user, nil
}

// GetUser gets a user by ID
func (us *UserService) GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user *domain.User

	cacheKey := util.GenerateCacheKey("user", id)
	cachedUser, err := us.cache.Get(ctx, cacheKey)
	if err == nil {
		err := util.Deserialize(cachedUser, &user)
		if err != nil {
			return nil, domain.ErrInternal
		}

		return user, nil
	}

	us.log.Debug(id.String())

	userDAO, err := us.storage.GetUserByID(ctx, id)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}
	user = converter.ToUser(userDAO)

	serializedUser, err := util.Serialize(user)
	if err != nil {
		return nil, domain.ErrInternal
	}

	err = us.cache.Set(ctx, cacheKey, serializedUser, 0)
	if err != nil {
		return nil, domain.ErrInternal
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
			return nil, domain.ErrInternal
		}

		return users, nil
	}

	usersDAO, err := us.storage.ListUsers(ctx, skip, limit)
	if err != nil {
		return nil, domain.ErrInternal
	}
	for _, userDAO := range usersDAO {
		users = append(users, *converter.ToUser(&userDAO))
	}

	usersSerialized, err := util.Serialize(users)
	if err != nil {
		return nil, domain.ErrInternal
	}

	err = us.cache.Set(ctx, cacheKey, usersSerialized, 0)
	if err != nil {
		return nil, domain.ErrInternal
	}

	return users, nil
}

// UpdateUser updates a user's name, email, and password
func (us *UserService) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	existingUserDAO, err := us.storage.GetUserByID(ctx, user.ID)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}
	existingUser := converter.ToUser(existingUserDAO)

	emptyData := user.Name == "" &&
		user.Email == "" &&
		user.Role == ""
	sameData := existingUser.Name == user.Name &&
		existingUser.Email == user.Email &&
		existingUser.Role == user.Role
	if emptyData || sameData {
		return nil, domain.ErrNoUpdatedData
	}

	updatedUserDAO, err := us.storage.UpdateUser(ctx, converter.ToUserDAO(user))
	if err != nil {
		if err == domain.ErrConflictingData {
			return nil, err
		}
		return nil, domain.ErrInternal
	}
	updatedUser := converter.ToUser(updatedUserDAO)

	cacheKey := util.GenerateCacheKey("user", user.ID)
	err = us.cache.Delete(ctx, cacheKey)
	if err != nil {
		return nil, domain.ErrInternal
	}

	serializedUser, err := util.Serialize(user)
	if err != nil {
		return nil, domain.ErrInternal
	}

	err = us.cache.Set(ctx, cacheKey, serializedUser, 0)
	if err != nil {
		return nil, domain.ErrInternal
	}

	err = us.cache.DeleteByPrefix(ctx, "users:*")
	if err != nil {
		return nil, domain.ErrInternal
	}

	return updatedUser, nil
}

// DeleteUser deletes a user by ID
func (us *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := us.storage.GetUserByID(ctx, id)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return err
		}
		return domain.ErrInternal
	}

	cacheKey := util.GenerateCacheKey("user", id)

	err = us.cache.Delete(ctx, cacheKey)
	if err != nil {
		return domain.ErrInternal
	}

	err = us.cache.DeleteByPrefix(ctx, "users:*")
	if err != nil {
		return domain.ErrInternal
	}

	return us.storage.DeleteUser(ctx, id)
}
