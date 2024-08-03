package collection

import (
	"context"

	"github.com/8thgencore/passfort/internal/domain"
	"github.com/8thgencore/passfort/internal/repository/storage/postgres/converter"
	"github.com/8thgencore/passfort/pkg/logger/sl"
	"github.com/google/uuid"
)

// CreateCollection creates a new collection
func (svc *CollectionService) CreateCollection(ctx context.Context, userID uuid.UUID, collection *domain.Collection) (*domain.Collection, error) {
	collectionDAO, err := svc.storage.CreateCollection(ctx, userID, converter.ToCollectionDAO(collection))
	if err != nil {
		svc.log.Error("Error creating collection:", sl.Err(err))
		return nil, domain.ErrDataNotAdded
	}
	return converter.ToCollection(collectionDAO), err
}

// ListCollectionsByUserID retrieves a list of collections with pagination for a specific user
func (svc *CollectionService) ListCollectionsByUserID(ctx context.Context, userID uuid.UUID, skip, limit uint64) ([]domain.Collection, error) {
	collectionsDAO, err := svc.storage.ListCollectionsByUserID(ctx, userID, skip, limit)
	if err != nil {
		svc.log.Error("Error listing collections", "user", userID, sl.Err(err))
		return nil, domain.ErrInternal

	}
	var collections []domain.Collection
	for _, collectionDAO := range collectionsDAO {
		collections = append(collections, *converter.ToCollection(&collectionDAO))
	}

	return collections, err
}

// GetCollection retrieves a collection by ID
func (svc *CollectionService) GetCollection(ctx context.Context, userID, collectionID uuid.UUID) (*domain.Collection, error) {
	collectionDAO, err := svc.storage.GetCollectionByID(ctx, collectionID)
	if err != nil {
		svc.log.Error("Error getting collection", "collection", collectionID, sl.Err(err))
		return nil, domain.ErrDataNotAdded
	}

	// Check if the user is part of the collection
	if !svc.isUserPartOfCollection(ctx, userID, collectionDAO.ID) {
		return nil, domain.ErrUnauthorized
	}

	return converter.ToCollection(collectionDAO), nil
}

// UpdateCollection updates a collection by ID, checking if the user is part of the collection
func (svc *CollectionService) UpdateCollection(ctx context.Context, userID uuid.UUID, collection *domain.Collection) (*domain.Collection, error) {
	// Check if the user is part of the collection
	if !svc.isUserPartOfCollection(ctx, userID, collection.ID) {
		return nil, domain.ErrUnauthorized
	}

	updatedCollectionDAO, err := svc.storage.UpdateCollection(ctx, converter.ToCollectionDAO(collection))
	if err != nil {
		svc.log.Error("Error updating collection", "collection", collection.ID, sl.Err(err))
		return nil, domain.ErrNoUpdatedData
	}

	return converter.ToCollection(updatedCollectionDAO), nil
}

// DeleteCollection deletes a collection by ID, checking if the user is part of the collection
func (svc *CollectionService) DeleteCollection(ctx context.Context, userID, collectionID uuid.UUID) error {
	if !svc.isUserPartOfCollection(ctx, userID, collectionID) {
		return domain.ErrUnauthorized
	}

	err := svc.storage.DeleteCollection(ctx, collectionID)
	if err != nil {
		svc.log.Error("Error deleting collection", "collection", collectionID, sl.Err(err))
		return domain.ErrDataNotDeleted
	}

	return nil
}

// isUserPartOfCollection checks if the user is part of the given collection
func (svc *CollectionService) isUserPartOfCollection(ctx context.Context, userID, collectionID uuid.UUID) bool {
	isPartOfCollection, err := svc.storage.IsUserPartOfCollection(ctx, userID, collectionID)
	if err != nil {
		svc.log.Error("Error checking if user is part of collection", "user", userID, "collection", collectionID, sl.Err(err))
		return false
	}

	return isPartOfCollection
}
