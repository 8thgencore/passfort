package secret

import (
	"log/slog"

	"github.com/8thgencore/passfort/internal/service/adapters/cache"
	"github.com/8thgencore/passfort/internal/service/adapters/storage"
)

/**
 * SecretService implements the service.SecretService interface
 * and provides access to the secret repository
 */
type SecretService struct {
	log               *slog.Logger
	secretStorage     storage.SecretRepository
	collectionStorage storage.CollectionRepository
	cache             cache.CacheRepository
}

// NewSecretService creates a new secret service instance
func NewSecretService(
	log *slog.Logger,
	secretStorage storage.SecretRepository,
	collectionStorage storage.CollectionRepository,
	cache cache.CacheRepository,
) *SecretService {
	return &SecretService{
		log,
		secretStorage,
		collectionStorage,
		cache,
	}
}
