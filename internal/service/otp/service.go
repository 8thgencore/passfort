package otp

import (
	"context"
	"time"

	"github.com/8thgencore/passfort/internal/domain"
	"github.com/8thgencore/passfort/pkg/util"
	"github.com/google/uuid"
)

// GenerateOTP generates a new OTP for the given user ID
func (os *OtpService) GenerateOTP(ctx context.Context, userID uuid.UUID) (string, error) {
	otpCode := util.GenerateOTP()

	// Store the generated OTP in your repository
	cacheKey := util.GenerateCacheKey("user_otp", userID)
	serializedOtp := []byte(otpCode)

	err := os.cache.Set(ctx, cacheKey, serializedOtp, 10*time.Minute)
	if err != nil {
		os.log.Error("Error storing OTP:", "error", err.Error())
		return "", domain.ErrInternal
	}

	// Return the generated OTP
	return otpCode, nil
}

// VerifyOTP verifies if the provided OTP is valid for the given user ID
func (os *OtpService) VerifyOTP(ctx context.Context, userID uuid.UUID, otpCode string) error {
	var oldOtpCode string
	// Retrieve the stored OTP for the user
	cacheKey := util.GenerateCacheKey("user_otp", userID)

	storedOTP, err := os.cache.Get(ctx, cacheKey)
	if err == nil {
		oldOtpCode = string(storedOTP)
	}

	// Verify the provided OTP against the stored OTP
	valid := util.ValidateOTP(otpCode)
	if !valid {
		os.log.Error("Invalid OTP provided")
		return domain.ErrInvalidOTP
	}

	os.log.Info("sdfsdf", otpCode, oldOtpCode)

	if otpCode != oldOtpCode {
		os.log.Error("Provided OTP does not match stored OTP")
		return domain.ErrInvalidOTP
	}

	// OTP is valid, remove it from storage to ensure one-time use
	err = os.cache.Delete(ctx, cacheKey)
	if err != nil {
		os.log.Error("Error removing OTP:", "error", err.Error())
		return domain.ErrInternal
	}

	return nil
}
