package collection

import (
	"context"
	"fmt"

	"github.com/8thgencore/passfort/internal/domain"
)

// CreateCollection creates a new collection
func (cs *CollectionService) CreateCollection(ctx context.Context, userID uint64, collection *domain.Collection) (*domain.Collection, error) {
	coll, err := cs.storage.CreateCollection(ctx, userID, collection)
	if err != nil {
		cs.log.Error("Error creating collection:", "error", err.Error())
	}
	return coll, err
}

// ListCollectionsByUserID retrieves a list of collections with pagination for a specific user
func (cs *CollectionService) ListCollectionsByUserID(ctx context.Context, userID, skip, limit uint64) ([]domain.Collection, error) {
	collections, err := cs.storage.ListCollectionsByUserID(ctx, userID, skip, limit)
	if err != nil {
		cs.log.Error(fmt.Sprintf("Error listing collections for user %d:", userID), "error", err.Error())
	}
	return collections, err
}

// GetCollection retrieves a collection by ID
func (cs *CollectionService) GetCollection(ctx context.Context, userID, collectionID uint64) (*domain.Collection, error) {
	collection, err := cs.storage.GetCollectionByID(ctx, collectionID)
	if err != nil {
		cs.log.Error(fmt.Sprintf("Error getting collection %d:", collectionID), "error", err.Error())
		return nil, err
	}

	// Check if the user is part of the collection
	if !cs.isUserPartOfCollection(ctx, userID, collection.ID) {
		return nil, domain.ErrUnauthorized
	}

	return collection, nil
}

// UpdateCollection updates a collection by ID, checking if the user is part of the collection
func (cs *CollectionService) UpdateCollection(ctx context.Context, userID uint64, collection *domain.Collection) (*domain.Collection, error) {
	// Check if the user is part of the collection
	if !cs.isUserPartOfCollection(ctx, userID, collection.ID) {
		return nil, domain.ErrUnauthorized
	}

	updatedColl, err := cs.storage.UpdateCollection(ctx, collection)
	if err != nil {
		cs.log.Error(fmt.Sprintf("Error updating collection %d:", collection.ID), "error", err.Error())
		return nil, err
	}

	return updatedColl, nil
}

// DeleteCollection deletes a collection by ID, checking if the user is part of the collection
func (cs *CollectionService) DeleteCollection(ctx context.Context, userID, collectionID uint64) error {
	// Check if the user is part of the collection
	collection, err := cs.storage.GetCollectionByID(ctx, collectionID)
	if err != nil {
		cs.log.Error(fmt.Sprintf("Error getting collection %d:", collectionID), "error", err.Error())
		return err
	}

	if !cs.isUserPartOfCollection(ctx, userID, collection.ID) {
		return domain.ErrUnauthorized
	}

	err = cs.storage.DeleteCollection(ctx, collectionID)
	if err != nil {
		cs.log.Error(fmt.Sprintf("Error deleting collection %d:", collectionID), "error", err.Error())
		return err
	}

	return nil
}

// isUserPartOfCollection checks if the user is part of the given collection
func (cs *CollectionService) isUserPartOfCollection(ctx context.Context, userID, collectionID uint64) bool {
	isPartOfCollection, err := cs.storage.IsUserPartOfCollection(ctx, userID, collectionID)
	if err != nil {
		cs.log.Error(fmt.Sprintf("Error checking if user %d is part of collection %d:", userID, collectionID), "error", err.Error())
		return false
	}

	return isPartOfCollection
}
