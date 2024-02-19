package storage

import (
	"context"

	"github.com/8thgencore/passfort/internal/repository/storage/postgres/dao"
	"github.com/google/uuid"
)

// UserRepository is an interface for interacting with user-related data
type UserRepository interface {
	// CreateUser inserts a new user into the database
	CreateUser(ctx context.Context, user *dao.UserDAO) (*dao.UserDAO, error)
	// GetUserByID selects a user by id
	GetUserByID(ctx context.Context, id uuid.UUID) (*dao.UserDAO, error)
	// GetUserByEmail selects a user by email
	GetUserByEmail(ctx context.Context, email string) (*dao.UserDAO, error)
	// ListUsers selects a list of users with pagination
	ListUsers(ctx context.Context, skip, limit uint64) ([]dao.UserDAO, error)
	// UpdateUser updates a user
	UpdateUser(ctx context.Context, user *dao.UserDAO) (*dao.UserDAO, error)
	// DeleteUser deletes a user
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

// CollectionRepository is an interface for interacting with collection-related data
type CollectionRepository interface {
	// CreateCollection inserts a new collection into the database
	CreateCollection(ctx context.Context, userID uuid.UUID, collection *dao.CollectionDAO) (*dao.CollectionDAO, error)
	// GetCollectionByID selects a collection by id
	GetCollectionByID(ctx context.Context, id uuid.UUID) (*dao.CollectionDAO, error)
	// ListCollectionsByUserID selects a list of collections for a specific user ID
	ListCollectionsByUserID(ctx context.Context, userID uuid.UUID, skip, limit uint64) ([]dao.CollectionDAO, error)
	// UpdateCollection updates a collection
	UpdateCollection(ctx context.Context, collection *dao.CollectionDAO) (*dao.CollectionDAO, error)
	// DeleteCollection deletes a collection
	DeleteCollection(ctx context.Context, id uuid.UUID) error
	// IsUserPartOfCollection checks if the user is part of the specified collection
	IsUserPartOfCollection(ctx context.Context, userID, collectionID uuid.UUID) (bool, error)
}
