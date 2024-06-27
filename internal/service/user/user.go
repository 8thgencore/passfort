package user

import (
	"context"

	"github.com/8thgencore/passfort/internal/domain"
	"github.com/8thgencore/passfort/internal/repository/storage/postgres/converter"
	"github.com/8thgencore/passfort/pkg/util"
	"github.com/google/uuid"
)

// GetUserByID gets a user by ID
func (svc *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user *domain.User

	cacheKey := util.GenerateCacheKey("user", id)
	cachedUser, err := svc.cache.Get(ctx, cacheKey)
	if err == nil {
		if err := util.Deserialize(cachedUser, &user); err != nil {
			return nil, domain.ErrInternal
		}

		return user, nil
	}

	userDAO, err := svc.storage.GetUserByID(ctx, id)
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

	if err = svc.cache.Set(ctx, cacheKey, serializedUser, 0); err != nil {
		return nil, domain.ErrInternal
	}

	return user, nil
}

// ListUsers lists all users
func (svc *UserService) ListUsers(ctx context.Context, skip, limit uint64) ([]domain.User, error) {
	var users []domain.User

	params := util.GenerateCacheKeyParams(skip, limit)
	cacheKey := util.GenerateCacheKey("users", params)

	cachedUsers, err := svc.cache.Get(ctx, cacheKey)
	if err == nil {
		if err := util.Deserialize(cachedUsers, &users); err != nil {
			return nil, domain.ErrInternal
		}

		return users, nil
	}

	usersDAO, err := svc.storage.ListUsers(ctx, skip, limit)
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

	if err = svc.cache.Set(ctx, cacheKey, usersSerialized, 0); err != nil {
		return nil, domain.ErrInternal
	}

	return users, nil
}

// UpdateUser updates a user's name, email, and password
func (svc *UserService) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	existingUserDAO, err := svc.storage.GetUserByID(ctx, user.ID)
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

	updatedUserDAO, err := svc.storage.UpdateUser(ctx, converter.ToUserDAO(user))
	if err != nil {
		if err == domain.ErrConflictingData {
			return nil, err
		}
		return nil, domain.ErrInternal
	}
	updatedUser := converter.ToUser(updatedUserDAO)

	cacheKey := util.GenerateCacheKey("user", user.ID)
	if err = svc.cache.Delete(ctx, cacheKey); err != nil {
		return nil, domain.ErrInternal
	}

	serializedUser, err := util.Serialize(user)
	if err != nil {
		return nil, domain.ErrInternal
	}

	if err = svc.cache.Set(ctx, cacheKey, serializedUser, 0); err != nil {
		return nil, domain.ErrInternal
	}

	if err = svc.cache.DeleteByPrefix(ctx, "users:*"); err != nil {
		return nil, domain.ErrInternal
	}

	return updatedUser, nil
}

// DeleteUser deletes a user by ID
func (svc *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := svc.storage.GetUserByID(ctx, id)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return err
		}
		return domain.ErrInternal
	}

	cacheKey := util.GenerateCacheKey("user", id)

	if err = svc.cache.Delete(ctx, cacheKey); err != nil {
		return domain.ErrInternal
	}

	if err = svc.cache.DeleteByPrefix(ctx, "users:*"); err != nil {
		return domain.ErrInternal
	}

	return svc.storage.DeleteUser(ctx, id)
}
