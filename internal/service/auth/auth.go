package auth

import (
	"context"

	"github.com/8thgencore/passfort/internal/domain"
	"github.com/8thgencore/passfort/internal/repository/storage/postgres/converter"
	"github.com/8thgencore/passfort/pkg/logger/sl"
	"github.com/8thgencore/passfort/pkg/util"
	"github.com/google/uuid"
)

// Login gives a registered user an access token if the credentials are valid
func (svc *AuthService) Login(ctx context.Context, email, password string) (string, string, error) {
	userDAO, err := svc.storage.GetUserByEmail(ctx, email)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return "", "", domain.ErrInvalidCredentials
		}
		svc.log.Error("failed to get the user by email", sl.Err(err))

		return "", "", domain.ErrInternal
	}
	user := converter.ToUser(userDAO)

	if !user.IsVerified {
		return "", "", domain.ErrUserNotVerified
	}

	err = util.CompareHash(password, user.Password)
	if err != nil {
		return "", "", domain.ErrInvalidCredentials
	}

	accessToken, refreshToken, err := svc.tokenService.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", "", domain.ErrTokenCreation
	}

	return accessToken, refreshToken, nil
}

// Register creates a new user
func (svc *AuthService) Register(ctx context.Context, user *domain.User) (*domain.User, error) {
	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		return nil, domain.ErrInternal
	}

	user.Password = hashedPassword

	userDAO, err := svc.storage.CreateUser(ctx, converter.ToUserDAO(user))
	if err != nil {
		if err == domain.ErrConflictingData {
			return nil, err
		}
		svc.log.Error("failed to create a user", sl.Err(err))

		return nil, domain.ErrInternal
	}
	user = converter.ToUser(userDAO)

	// Send confirm otp code
	otp, err := svc.otp.GenerateOTP(ctx, user.ID)
	if err != nil {
		return nil, domain.ErrInternal
	}
	_, err = svc.mailClient.SendConfirmationEmail(ctx, user.Email, otp)
	if err != nil {
		svc.log.Error("failed send confirmation email", sl.Err(err))
		return nil, domain.ErrInternal
	}

	// Update cache
	cacheKey := util.GenerateCacheKey("user", user.ID)
	userSerialized, err := util.Serialize(user)
	if err != nil {
		return nil, domain.ErrInternal
	}

	err = svc.cache.Set(ctx, cacheKey, userSerialized, 0)
	if err != nil {
		return nil, domain.ErrInternal
	}

	err = svc.cache.DeleteByPrefix(ctx, "users:*")
	if err != nil {
		return nil, domain.ErrInternal
	}

	return user, nil
}

// ConfirmRegistration confirms user registration with OTP code
func (svc *AuthService) ConfirmRegistration(ctx context.Context, email, otp string) error {
	// Implement the logic to confirm user registration with OTP code
	// Retrieve user by email
	userDAO, err := svc.storage.GetUserByEmail(ctx, email)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return domain.ErrInvalidOTP
		}
		svc.log.Error("failed to get the user by email", sl.Err(err))
		return domain.ErrInternal
	}

	// Validate OTP
	if err := svc.otp.VerifyOTP(ctx, userDAO.ID, otp); err != nil {
		return domain.ErrInvalidOTP
	}

	userDAO.IsVerified = true
	_, err = svc.storage.UpdateUser(ctx, userDAO)
	if err != nil {
		svc.log.Error("failed to update user", sl.Err(err))
		return domain.ErrInternal
	}

	return nil
}

// RequestNewRegistrationCode requests a new registration confirmation code for a user
func (svc *AuthService) RequestNewRegistrationCode(ctx context.Context, email string) error {
	// Retrieve user by email
	userDAO, err := svc.storage.GetUserByEmail(ctx, email)
	if err != nil {
		svc.log.Error("failed to get the user by email", sl.Err(err))
		return nil
	}
	user := converter.ToUser(userDAO)

	// Check if OTP already exists for the user
	exists, err := svc.otp.CheckCacheForKey(ctx, user.ID)
	if err != nil {
		svc.log.Error("failed to check cache for OTP", sl.Err(err))
		return domain.ErrInternal
	}
	if exists {
		return domain.ErrOTPAlreadySent
	}

	// Generate and send new registration confirmation OTP
	otp, err := svc.otp.GenerateOTP(ctx, user.ID)
	if err != nil {
		svc.log.Error("failed to generate new registration confirmation OTP", sl.Err(err))
		return domain.ErrInternal
	}

	_, err = svc.mailClient.SendConfirmationEmail(ctx, user.Email, otp)
	if err != nil {
		svc.log.Error("failed to send new registration confirmation email", sl.Err(err))
		return domain.ErrInternal
	}

	return nil
}

// RefreshToken refreshes the access token for the user
func (svc *AuthService) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	token, err := svc.tokenService.ParseUserClaims(refreshToken)
	if err != nil {
		return "", "", domain.ErrInvalidRefreshToken
	}

	// Caching a revoked token
	err = svc.tokenService.RevokeToken(ctx, token.ID)
	if err != nil {
		return "", "", err
	}

	userDAO, err := svc.storage.GetUserByID(ctx, token.UserID)
	if err != nil {
		return "", "", domain.ErrDataNotFound
	}

	user := converter.ToUser(userDAO)

	accessToken, refreshToken, err := svc.tokenService.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", "", domain.ErrTokenCreation
	}

	return accessToken, refreshToken, nil
}

// Logout invalidates the access token, logging the user out
func (svc *AuthService) Logout(ctx context.Context, token *domain.UserClaims) error {
	_, err := svc.storage.GetUserByID(ctx, token.UserID)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return err
		}
		return domain.ErrInternal
	}

	// Caching a revoked token
	err = svc.tokenService.RevokeToken(ctx, token.ID)
	if err != nil {
		return err
	}

	return nil
}

// ChangePassword implements the ChangePassword method of the AuthService interface
func (svc *AuthService) ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	// Retrieve the user based on the userID
	user, err := svc.storage.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	// Verify the old password
	err = util.CompareHash(oldPassword, user.Password)
	if err != nil {
		return domain.ErrPasswordsDoNotMatch
	}

	// Update the password
	hashedPassword, err := util.HashPassword(newPassword)
	if err != nil {
		return domain.ErrInternal
	}

	user.Password = hashedPassword
	_, err = svc.storage.UpdateUser(ctx, user)
	return err
}

// ForgotPassword initiates the process of resetting a forgotten password
func (svc *AuthService) ForgotPassword(ctx context.Context, email string) error {
	// Retrieve user by email
	userDAO, err := svc.storage.GetUserByEmail(ctx, email)
	if err != nil {
		svc.log.Error("failed to get the user by email", sl.Err(err))
		return nil
	}
	user := converter.ToUser(userDAO)

	// Generate and send reset OTP
	resetOTP, err := svc.otp.GenerateOTP(ctx, user.ID)
	if err != nil {
		svc.log.Error("failed to update user", sl.Err(err))
		return domain.ErrInternal
	}

	_, err = svc.mailClient.SendPasswordReset(ctx, user.Email, resetOTP)
	if err != nil {
		svc.log.Error("failed send reset password", sl.Err(err))
		return domain.ErrInternal
	}

	return nil
}

// ResetPassword confirms password reset with OTP code
func (svc *AuthService) ResetPassword(ctx context.Context, email, newPassword, otp string) error {
	// Retrieve user by email
	userDAO, err := svc.storage.GetUserByEmail(ctx, email)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return domain.ErrInvalidOTP
		}
		svc.log.Error("failed to get the user by email", sl.Err(err))
		return domain.ErrInternal
	}
	user := converter.ToUser(userDAO)

	// Validate OTP
	if err := svc.otp.VerifyOTP(ctx, user.ID, otp); err != nil {
		return domain.ErrInvalidOTP
	}

	// Update the password
	hashedPassword, err := util.HashPassword(newPassword)
	if err != nil {
		return domain.ErrInternal
	}

	user.Password = hashedPassword
	_, err = svc.storage.UpdateUser(ctx, converter.ToUserDAO(user))
	if err != nil {
		svc.log.Error("failed to update user", sl.Err(err))
		return domain.ErrInternal
	}

	return nil
}
