package masterpassword

import (
	"log/slog"
	"time"

	"github.com/8thgencore/passfort/internal/service/adapters/cache"
	"github.com/8thgencore/passfort/internal/service/adapters/storage"
)

/**
 * MasterPasswordService implements service.MasterPasswordService interface
 * and provides an access to the user repository
 * and cache for storing master password validation states
 */
type MasterPasswordService struct {
	log               *slog.Logger
	storage           storage.UserRepository
	cache             cache.CacheRepository
	masterPasswordTTL time.Duration
}

// NewMasterPasswordService creates a new master password service instance
func NewMasterPasswordService(
	log *slog.Logger,
	storage storage.UserRepository,
	cache cache.CacheRepository,
	masterPasswordTTL time.Duration,
) *MasterPasswordService {
	return &MasterPasswordService{
		log:               log,
		storage:           storage,
		cache:             cache,
		masterPasswordTTL: masterPasswordTTL,
	}
}
