package auth

import (
	"log/slog"

	"github.com/8thgencore/passfort/internal/service"
	"github.com/8thgencore/passfort/internal/service/adapters/cache"
	"github.com/8thgencore/passfort/internal/service/adapters/storage"
)

/**
 * AuthService implements service.AuthService interface
 * and provides an access to the user repository
 * and token service
 */
type AuthService struct {
	log     *slog.Logger
	storage storage.UserRepository
	cache   cache.CacheRepository
	ts      service.TokenService
}

// NewAuthService creates a new auth service instance
func NewAuthService(log *slog.Logger, storage storage.UserRepository, cache cache.CacheRepository, ts service.TokenService) *AuthService {
	return &AuthService{
		log,
		storage,
		cache,
		ts,
	}
}
