package masterpassword

import (
	"log/slog"
	"time"

	"github.com/8thgencore/passfort/internal/service/adapters/cache"
	"github.com/8thgencore/passfort/internal/service/adapters/storage"
	"github.com/8thgencore/passfort/internal/service/secret"
)

/**
 * MasterPasswordService implements service.MasterPasswordService interface
 * and provides an access to the user repository
 * and cache for storing master password activation states
 */
type MasterPasswordService struct {
	log               *slog.Logger
	userStorage       storage.UserRepository
	cache             cache.CacheRepository
	secretSvc         secret.SecretService
	masterPasswordTTL time.Duration
}

// NewMasterPasswordService creates a new master password service instance
func NewMasterPasswordService(
	log *slog.Logger,
	userStorage storage.UserRepository,
	cache cache.CacheRepository,
	secretSvc secret.SecretService,
	masterPasswordTTL time.Duration,
) *MasterPasswordService {
	return &MasterPasswordService{
		log:               log,
		userStorage:       userStorage,
		cache:             cache,
		secretSvc:         secretSvc,
		masterPasswordTTL: masterPasswordTTL,
	}
}
