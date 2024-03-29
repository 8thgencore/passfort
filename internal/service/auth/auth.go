package auth

import (
	"context"

	"github.com/8thgencore/passfort/internal/domain"
	"github.com/8thgencore/passfort/internal/repository/storage/postgres/converter"
	"github.com/8thgencore/passfort/pkg/util"
	"github.com/google/uuid"
)

// Login gives a registered user an access token if the credentials are valid
func (as *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	userDAO, err := as.storage.GetUserByEmail(ctx, email)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return "", domain.ErrInvalidCredentials
		}
		as.log.Error("failed to get the user by email", "error", err.Error())

		return "", domain.ErrInternal
	}
	user := converter.ToUser(userDAO)

	if !user.IsVerified {
		return "", domain.ErrUserNotVerified
	}

	err = util.ComparePassword(password, user.Password)
	if err != nil {
		return "", domain.ErrInvalidCredentials
	}

	accessToken, err := as.tokenService.CreateToken(user)
	if err != nil {
		return "", domain.ErrTokenCreation
	}

	return accessToken, nil
}

// Register creates a new user
func (as *AuthService) Register(ctx context.Context, user *domain.User) (*domain.User, error) {
	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		return nil, domain.ErrInternal
	}

	user.Password = hashedPassword

	userDAO, err := as.storage.CreateUser(ctx, converter.ToUserDAO(user))
	if err != nil {
		if err == domain.ErrConflictingData {
			return nil, err
		}
		as.log.Error("failed to create a user", "error", err.Error())

		return nil, domain.ErrInternal
	}
	user = converter.ToUser(userDAO)

	// Send confirm otp code
	otp, err := as.otp.GenerateOTP(ctx, user.ID)
	if err != nil {
		return nil, domain.ErrInternal
	}
	_, err = as.mailClient.SendConfirmationEmail(ctx, user.Email, otp)
	if err != nil {
		as.log.Error("failed send confirmation email", "error", err.Error())
		return nil, domain.ErrInternal
	}

	// Update cache
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

// ConfirmRegistration confirms user registration with OTP code
func (as *AuthService) ConfirmRegistration(ctx context.Context, email, otp string) error {
	// Implement the logic to confirm user registration with OTP code
	// Retrieve user by email
	userDAO, err := as.storage.GetUserByEmail(ctx, email)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return domain.ErrInvalidOTP
		}
		as.log.Error("failed to get the user by email", "error", err.Error())
		return domain.ErrInternal
	}

	// Validate OTP
	if err := as.otp.VerifyOTP(ctx, userDAO.ID, otp); err != nil {
		return domain.ErrInvalidOTP
	}

	userDAO.IsVerified = true
	_, err = as.storage.UpdateUser(ctx, userDAO)
	if err != nil {
		as.log.Error("failed to update user", "error", err.Error())
		return domain.ErrInternal
	}

	return nil
}

// ChangePassword implements the ChangePassword method of the AuthService interface
func (as *AuthService) ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
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

// RequestResetPassword initiates the process of resetting a forgotten password
func (as *AuthService) RequestResetPassword(ctx context.Context, email string) error {
	// Implement the logic to initiate the process of resetting a forgotten password
	// Retrieve user by email
	userDAO, err := as.storage.GetUserByEmail(ctx, email)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return domain.ErrInvalidCredentials
		}
		as.log.Error("failed to get the user by email", "error", err.Error())
		return domain.ErrInternal
	}
	user := converter.ToUser(userDAO)

	// Generate and send reset OTP
	resetOTP, err := as.otp.GenerateOTP(ctx, user.ID)
	if err != nil {
		as.log.Error("failed to update user", "error", err.Error())
		return domain.ErrInternal
	}

	_, err = as.mailClient.SendPasswordReset(ctx, user.Email, resetOTP)
	if err != nil {
		as.log.Error("failed send reset password", "error", err.Error())
		return domain.ErrInternal
	}

	return nil
}

// ConfirmResetPassword confirms password reset with OTP code
func (as *AuthService) ConfirmResetPassword(ctx context.Context, email, otp string) error {
	// Implement the logic to confirm password reset with OTP code
	// Retrieve user by email
	userDAO, err := as.storage.GetUserByEmail(ctx, email)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return domain.ErrInvalidOTP
		}
		as.log.Error("failed to get the user by email", "error", err.Error())
		return domain.ErrInternal
	}
	user := converter.ToUser(userDAO)

	// Validate OTP
	if err := as.otp.VerifyOTP(ctx, user.ID, otp); err != nil {
		return domain.ErrInvalidOTP
	}

	return nil
}

// SetNewPassword resets user password after confirmation with OTP code
func (as *AuthService) SetNewPassword(ctx context.Context, email, newPassword, otp string) error {
	// Implement the logic to set a new password after confirmation with OTP code
	// Retrieve user by email
	userDAO, err := as.storage.GetUserByEmail(ctx, email)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return domain.ErrInvalidOTP
		}
		as.log.Error("failed to get the user by email", "error", err.Error())
		return domain.ErrInternal
	}
	user := converter.ToUser(userDAO)

	// Validate OTP
	if err := as.otp.VerifyOTP(ctx, user.ID, otp); err != nil {
		return domain.ErrInvalidOTP
	}

	// Update the password
	hashedPassword, err := util.HashPassword(newPassword)
	if err != nil {
		return domain.ErrInternal
	}

	user.Password = hashedPassword
	_, err = as.storage.UpdateUser(ctx, converter.ToUserDAO(user))
	if err != nil {
		as.log.Error("failed to update user", "error", err.Error())
		return domain.ErrInternal
	}

	return nil
}
