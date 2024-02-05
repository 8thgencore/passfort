package storage

import (
	"context"

	"github.com/8thgencore/passfort/internal/domain"
)

// UserRepository is an interface for interacting with user-related data
type UserRepository interface {
	// CreateUser inserts a new user into the database
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	// GetUserByID selects a user by id
	GetUserByID(ctx context.Context, id uint64) (*domain.User, error)
	// GetUserByEmail selects a user by email
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	// ListUsers selects a list of users with pagination
	ListUsers(ctx context.Context, skip, limit uint64) ([]domain.User, error)
	// UpdateUser updates a user
	UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	// DeleteUser deletes a user
	DeleteUser(ctx context.Context, id uint64) error
}

// CollectionRepository is an interface for interacting with collection-related data
type CollectionRepository interface {
	// CreateCollection inserts a new collection into the database
	CreateCollection(ctx context.Context, userID uint64, collection *domain.Collection) (*domain.Collection, error)
	// GetCollectionByID selects a collection by id
	GetCollectionByID(ctx context.Context, id uint64) (*domain.Collection, error)
	// ListCollectionsByUserID selects a list of collections for a specific user ID
	ListCollectionsByUserID(ctx context.Context, userID uint64, skip, limit uint64) ([]domain.Collection, error)
	// UpdateCollection updates a collection
	UpdateCollection(ctx context.Context, collection *domain.Collection) (*domain.Collection, error)
	// DeleteCollection deletes a collection
	DeleteCollection(ctx context.Context, id uint64) error
	// IsUserPartOfCollection checks if the user is part of the specified collection
	IsUserPartOfCollection(ctx context.Context, userID, collectionID uint64) (bool, error)
}
