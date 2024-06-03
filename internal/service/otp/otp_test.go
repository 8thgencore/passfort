package otp_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"log/slog"

	"github.com/8thgencore/passfort/internal/domain"
	"github.com/8thgencore/passfort/internal/service/adapters/cache/mocks"
	"github.com/8thgencore/passfort/internal/service/otp"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupOtpService() (*otp.OtpService, *mocks.CacheRepository) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	cacheMock := &mocks.CacheRepository{}
	otpService := otp.NewOtpService(logger, cacheMock)
	return otpService, cacheMock
}

func TestGenerateOTP(t *testing.T) {
	otpService, cache := setupOtpService()

	userID, _ := uuid.NewRandom()

	t.Run("success", func(t *testing.T) {
		cache.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

		generatedOTP, err := otpService.GenerateOTP(context.Background(), userID)
		assert.NoError(t, err)
		assert.Equal(t, generatedOTP, generatedOTP)
		cache.AssertExpectations(t)
		cache.ExpectedCalls = nil
	})

	t.Run("failure on cache set", func(t *testing.T) {
		cache.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("cache error"))

		_, err := otpService.GenerateOTP(context.Background(), userID)
		assert.Error(t, err)
		assert.Equal(t, domain.ErrInternal, err)
		cache.AssertExpectations(t)
		cache.ExpectedCalls = nil
	})
}

func TestVerifyOTP(t *testing.T) {
	otpService, cache := setupOtpService()

	userID, _ := uuid.NewRandom()
	otpCode := "123456"

	t.Run("success", func(t *testing.T) {
		cache.On("Get", mock.Anything, mock.Anything).Return([]byte(otpCode), nil)
		cache.On("Delete", mock.Anything, mock.Anything).Return(nil)

		err := otpService.VerifyOTP(context.Background(), userID, otpCode)
		assert.NoError(t, err)
		cache.AssertExpectations(t)
	})

	t.Run("invalid OTP code", func(t *testing.T) {
		cache.On("Get", mock.Anything, mock.Anything).Return([]byte(otpCode), nil)

		invalidOTPCode := "6543210"

		err := otpService.VerifyOTP(context.Background(), userID, invalidOTPCode)
		assert.Error(t, err)
		assert.Equal(t, domain.ErrInvalidOTP, err)
		cache.AssertNotCalled(t, "Delete")
		cache.ExpectedCalls = nil
	})

	t.Run("OTP not matching stored OTP", func(t *testing.T) {
		cache.On("Get", mock.Anything, mock.Anything).Return([]byte("654321"), nil)

		err := otpService.VerifyOTP(context.Background(), userID, otpCode)
		assert.Error(t, err)
		assert.Equal(t, domain.ErrInvalidOTP, err)
		cache.AssertNotCalled(t, "Delete")
		cache.ExpectedCalls = nil
	})

	t.Run("cache get error", func(t *testing.T) {
		cache.On("Get", mock.Anything, mock.Anything).Return(nil, errors.New("cache error"))

		err := otpService.VerifyOTP(context.Background(), userID, otpCode)
		assert.Error(t, err)
		assert.Equal(t, domain.ErrInvalidOTP, err)
		cache.AssertNotCalled(t, "Delete")
		cache.ExpectedCalls = nil
	})

	t.Run("cache delete error", func(t *testing.T) {
		cache.On("Get", mock.Anything, mock.Anything).Return([]byte(otpCode), nil)
		cache.On("Delete", mock.Anything, mock.Anything).Return(errors.New("cache error"))

		err := otpService.VerifyOTP(context.Background(), userID, otpCode)
		assert.Error(t, err)
		assert.Equal(t, domain.ErrInternal, err)
		cache.AssertExpectations(t)
		cache.ExpectedCalls = nil
	})
}

func TestCheckCacheForKey(t *testing.T) {
	otpService, cache := setupOtpService()

	userID, _ := uuid.NewRandom()

	t.Run("exists", func(t *testing.T) {
		cache.On("Exists", mock.Anything, mock.Anything).Return(true, nil)

		exists, err := otpService.CheckCacheForKey(context.Background(), userID)
		assert.NoError(t, err)
		assert.True(t, exists)
		cache.AssertExpectations(t)
		cache.ExpectedCalls = nil
	})

	t.Run("does not exist", func(t *testing.T) {
		cache.On("Exists", mock.Anything, mock.Anything).Return(false, nil)

		exists, err := otpService.CheckCacheForKey(context.Background(), userID)
		assert.NoError(t, err)
		assert.False(t, exists)
		cache.AssertExpectations(t)
		cache.ExpectedCalls = nil
	})

	t.Run("cache exists error", func(t *testing.T) {
		cache.On("Exists", mock.Anything, mock.Anything).Return(false, errors.New("cache error"))

		exists, err := otpService.CheckCacheForKey(context.Background(), userID)
		assert.Error(t, err)
		assert.Equal(t, domain.ErrInternal, err)
		assert.False(t, exists)
		cache.AssertExpectations(t)
		cache.ExpectedCalls = nil
	})
}
