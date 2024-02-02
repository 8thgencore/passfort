package auth

import (
	"context"

	"github.com/8thgencore/passfort/internal/domain"
	"github.com/8thgencore/passfort/pkg/util"
)

// Register creates a new user
func (as *AuthService) Register(ctx context.Context, user *domain.User) (*domain.User, error) {
	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		return nil, domain.ErrInternal
	}

	user.Password = hashedPassword

	user, err = as.storage.CreateUser(ctx, user)
	if err != nil {
		if err == domain.ErrConflictingData {
			return nil, err
		}

		return nil, domain.ErrInternal
	}

	cacheKey := util.GenerateCacheKey("user", user.ID)
	userSerialized, err := util.Serialize(user)
	if err != nil {
		return nil, domain.ErrInternal
	}

	err = as.cache.Set(ctx, cacheKey, userSerialized, 0)
	if err != nil {
		return nil, domain.ErrInternal
	}

	err = as.cache.DeleteByPrefix(ctx, "users:*")
	if err != nil {
		return nil, domain.ErrInternal
	}

	return user, nil
}

// Login gives a registered user an access token if the credentials are valid
func (as *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := as.storage.GetUserByEmail(ctx, email)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return "", domain.ErrInvalidCredentials
		}
		return "", domain.ErrInternal
	}

	err = util.ComparePassword(password, user.Password)
	if err != nil {
		return "", domain.ErrInvalidCredentials
	}

	accessToken, err := as.ts.CreateToken(user)
	if err != nil {
		return "", domain.ErrTokenCreation
	}

	return accessToken, nil
}

// ChangePassword implements the ChangePassword method of the AuthService interface
func (as *AuthService) ChangePassword(ctx context.Context, userID uint64, oldPassword, newPassword string) error {
	// Retrieve the user based on the userID
	user, err := as.storage.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	// Verify the old password
	err = util.ComparePassword(oldPassword, user.Password)
	if err != nil {
		return domain.ErrPasswordsDoNotMatch
	}

	// Update the password
	hashedPassword, err := util.HashPassword(newPassword)
	if err != nil {
		return domain.ErrInternal
	}

	user.Password = hashedPassword
	_, err = as.storage.UpdateUser(ctx, user)
	return err
}
