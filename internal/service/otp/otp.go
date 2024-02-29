package otp

import (
	"log/slog"

	"github.com/8thgencore/passfort/internal/service/adapters/cache"
)

/**
 * OtpService implements service.OtpService interface
 * and provides OTP (One-Time Password) related functionality
 */
type OtpService struct {
	log   *slog.Logger
	cache cache.CacheRepository
}

// NewOtpService creates a new OTP service instance
func NewOtpService(log *slog.Logger, cache cache.CacheRepository) *OtpService {
	return &OtpService{
		log,
		cache,
	}
}
