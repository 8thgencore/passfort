package user

import (
	"github.com/8thgencore/passfort/internal/service/adapters/cache"
	"github.com/8thgencore/passfort/internal/service/adapters/storage"
	"log/slog"
)

/**
 * UserService implements service.UserService interface
 * and provides an access to the user repository
 * and cache service
 */
type UserService struct {
	log     *slog.Logger
	storage storage.UserRepository
	cache   cache.CacheRepository
}

// NewUserService creates a new user service instance
func NewUserService(log *slog.Logger, storage storage.UserRepository, cache cache.CacheRepository) *UserService {
	return &UserService{
		log,
		storage,
		cache,
	}
}
