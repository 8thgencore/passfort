package collection

import (
	"log/slog"

	"github.com/8thgencore/passfort/internal/service/adapters/storage"
)

/**
 * CollectionService implements service.CollectionService interface
 * and provides an access to the collection repository
 */
type CollectionService struct {
	log     *slog.Logger
	storage storage.CollectionRepository
}

// NewCollectionService creates a new collection service instance
func NewCollectionService(log *slog.Logger, storage storage.CollectionRepository) *CollectionService {
	return &CollectionService{
		log,
		storage,
	}
}
