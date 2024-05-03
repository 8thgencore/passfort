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

// RequestNewRegistrationCode requests a new registration confirmation code for a user
func (as *AuthService) RequestNewRegistrationCode(ctx context.Context, email string) error {
	// Retrieve user by email
	userDAO, err := as.storage.GetUserByEmail(ctx, email)
	if err != nil {
		as.log.Error("failed to get the user by email", "error", err.Error())
		return nil
	}
	user := converter.ToUser(userDAO)

	// Check if OTP already exists for the user
	exists, err := as.otp.CheckCacheForKey(ctx, user.ID)
	if err != nil {
		as.log.Error("failed to check cache for OTP", "error", err.Error())
		return domain.ErrInternal
	}
	if exists {
		return domain.ErrOTPAlreadySent
	}

	// Generate and send new registration confirmation OTP
	otp, err := as.otp.GenerateOTP(ctx, user.ID)
	if err != nil {
		as.log.Error("failed to generate new registration confirmation OTP", "error", err.Error())
		return domain.ErrInternal
	}

	_, err = as.mailClient.SendConfirmationEmail(ctx, user.Email, otp)
	if err != nil {
		as.log.Error("failed to send new registration confirmation email", "error", err.Error())
		return domain.ErrInternal
	}

	return nil
}

// Logout invalidates the access token, logging the user out
func (as *AuthService) Logout(ctx context.Context, token *domain.TokenPayload) error {
	_, err := as.storage.GetUserByID(ctx, token.UserID)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return err
		}
		return domain.ErrInternal
	}

	cacheKey := util.GenerateCacheKey("token", token.ID)
	userSerialized, err := util.Serialize(token)
	if err != nil {
		return domain.ErrInternal
	}

	err = as.cache.Set(ctx, cacheKey, userSerialized, 0)
	if err != nil {
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

// ForgotPassword initiates the process of resetting a forgotten password
func (as *AuthService) ForgotPassword(ctx context.Context, email string) error {
	// Retrieve user by email
	userDAO, err := as.storage.GetUserByEmail(ctx, email)
	if err != nil {
		as.log.Error("failed to get the user by email", "error", err.Error())
		return nil
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

// ResetPassword confirms password reset with OTP code
func (as *AuthService) ResetPassword(ctx context.Context, email, newPassword, otp string) error {
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
