package collection

import (
	"context"
	"fmt"

	"github.com/8thgencore/passfort/internal/domain"
	"github.com/8thgencore/passfort/internal/repository/storage/postgres/converter"
	"github.com/google/uuid"
)

// CreateCollection creates a new collection
func (cs *CollectionService) CreateCollection(ctx context.Context, userID uuid.UUID, collection *domain.Collection) (*domain.Collection, error) {
	collectionDAO, err := cs.storage.CreateCollection(ctx, userID, converter.ToCollectionDAO(collection))
	if err != nil {
		cs.log.Error("Error creating collection:", "error", err.Error())
		return nil, domain.ErrDataNotAdded
	}
	return converter.ToCollection(collectionDAO), err
}

// ListCollectionsByUserID retrieves a list of collections with pagination for a specific user
func (cs *CollectionService) ListCollectionsByUserID(ctx context.Context, userID uuid.UUID, skip, limit uint64) ([]domain.Collection, error) {
	collectionsDAO, err := cs.storage.ListCollectionsByUserID(ctx, userID, skip, limit)
	if err != nil {
		cs.log.Error(fmt.Sprintf("Error listing collections for user %d:", userID), "error", err.Error())
		return nil, domain.ErrInternal

	}
	var collections []domain.Collection
	for _, collectionDAO := range collectionsDAO {
		collections = append(collections, *converter.ToCollection(&collectionDAO))
	}

	return collections, err
}

// GetCollection retrieves a collection by ID
func (cs *CollectionService) GetCollection(ctx context.Context, userID, collectionID uuid.UUID) (*domain.Collection, error) {
	collectionDAO, err := cs.storage.GetCollectionByID(ctx, collectionID)
	if err != nil {
		cs.log.Error(fmt.Sprintf("Error getting collection %d:", collectionID), "error", err.Error())
		return nil, domain.ErrDataNotAdded
	}

	// Check if the user is part of the collection
	if !cs.isUserPartOfCollection(ctx, userID, collectionDAO.ID) {
		return nil, domain.ErrUnauthorized
	}

	return converter.ToCollection(collectionDAO), nil
}

// UpdateCollection updates a collection by ID, checking if the user is part of the collection
func (cs *CollectionService) UpdateCollection(ctx context.Context, userID uuid.UUID, collection *domain.Collection) (*domain.Collection, error) {
	// Check if the user is part of the collection
	if !cs.isUserPartOfCollection(ctx, userID, collection.ID) {
		return nil, domain.ErrUnauthorized
	}

	updatedCollectionDAO, err := cs.storage.UpdateCollection(ctx, converter.ToCollectionDAO(collection))
	if err != nil {
		cs.log.Error(fmt.Sprintf("Error updating collection %d:", collection.ID), "error", err.Error())
		return nil, domain.ErrNoUpdatedData
	}

	return converter.ToCollection(updatedCollectionDAO), nil
}

// DeleteCollection deletes a collection by ID, checking if the user is part of the collection
func (cs *CollectionService) DeleteCollection(ctx context.Context, userID, collectionID uuid.UUID) error {
	if !cs.isUserPartOfCollection(ctx, userID, collectionID) {
		return domain.ErrUnauthorized
	}

	err := cs.storage.DeleteCollection(ctx, collectionID)
	if err != nil {
		cs.log.Error(fmt.Sprintf("Error deleting collection %d:", collectionID), "error", err.Error())
		return domain.ErrDataNotDeleted
	}

	return nil
}

// isUserPartOfCollection checks if the user is part of the given collection
func (cs *CollectionService) isUserPartOfCollection(ctx context.Context, userID, collectionID uuid.UUID) bool {
	isPartOfCollection, err := cs.storage.IsUserPartOfCollection(ctx, userID, collectionID)
	if err != nil {
		cs.log.Error(fmt.Sprintf("Error checking if user %d is part of collection %d:", userID, collectionID), "error", err.Error())
		return false
	}

	return isPartOfCollection
}
