package otp

import (
	"context"
	"time"

	"github.com/8thgencore/passfort/internal/domain"
	"github.com/8thgencore/passfort/pkg/util"
	"github.com/google/uuid"
)

// GenerateOTP generates a new OTP for the given user ID
func (svc *OtpService) GenerateOTP(ctx context.Context, userID uuid.UUID) (string, error) {
	otpCode := util.GenerateOTP()

	// Store the generated OTP in your repository
	cacheKey := util.GenerateCacheKey("user_otp", userID)
	serializedOtp := []byte(otpCode)

	if err := svc.cache.Set(ctx, cacheKey, serializedOtp, 10*time.Minute); err != nil {
		svc.log.Error("Error storing OTP:", "error", err.Error())
		return "", domain.ErrInternal
	}

	// Return the generated OTP
	return otpCode, nil
}

// VerifyOTP verifies if the provided OTP is valid for the given user ID
func (svc *OtpService) VerifyOTP(ctx context.Context, userID uuid.UUID, otpCode string) error {
	var oldOtpCode string
	// Retrieve the stored OTP for the user
	cacheKey := util.GenerateCacheKey("user_otp", userID)

	storedOTP, err := svc.cache.Get(ctx, cacheKey)
	if err == nil {
		oldOtpCode = string(storedOTP)
	}

	// Verify the provided OTP against the stored OTP
	valid := util.ValidateOTP(otpCode)
	if !valid {
		svc.log.Error("Invalid OTP provided")
		return domain.ErrInvalidOTP
	}

	if otpCode != oldOtpCode {
		svc.log.Error("Provided OTP does not match stored OTP")
		return domain.ErrInvalidOTP
	}

	// OTP is valid, remove it from storage to ensure one-time use
	if err := svc.cache.Delete(ctx, cacheKey); err != nil {
		svc.log.Error("Error removing OTP:", "error", err.Error())
		return domain.ErrInternal
	}

	return nil
}

// CheckCacheForKey checks if a value exists in the cache for the given key
func (svc *OtpService) CheckCacheForKey(ctx context.Context, userID uuid.UUID) (bool, error) {
	// Retrieve the stored OTP for the user
	key := util.GenerateCacheKey("user_otp", userID)

	// Check if the value exists in the cache for the given key
	exists, err := svc.cache.Exists(ctx, key)
	if err != nil {
		svc.log.Error("Error checking cache for key:", "error", err.Error())
		return false, domain.ErrInternal
	}

	return exists, nil
}
