package user

import (
	"github.com/8thgencore/passfort/internal/service/adapters/cache"
	"github.com/8thgencore/passfort/internal/service/adapters/storage"
)

/**
 * UserService implements service.UserService interface
 * and provides an access to the user repository
 * and cache service
 */
type UserService struct {
	storage storage.UserRepository
	cache   cache.CacheRepository
}

// NewUserService creates a new user service instance
func NewUserService(storage storage.UserRepository, cache cache.CacheRepository) *UserService {
	return &UserService{
		storage,
		cache,
	}
}
