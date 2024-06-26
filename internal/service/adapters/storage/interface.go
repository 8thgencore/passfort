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

// SecretRepository is an interface for interacting with secret-related data
type SecretRepository interface {
	// CreateSecret inserts a new secret into the database
	CreateSecret(ctx context.Context, collectionID uuid.UUID, secret *dao.SecretDAO) (*dao.SecretDAO, error)
	// GetSecretByID selects a secret by id
	GetSecretByID(ctx context.Context, id uuid.UUID) (*dao.SecretDAO, error)
	// ListSecretsByCollectionID selects a list of secrets for a specific collection ID
	ListSecretsByCollectionID(ctx context.Context, collectionID uuid.UUID, skip, limit uint64) ([]dao.SecretDAO, error)
	// UpdateSecret updates a secret
	UpdateSecret(ctx context.Context, secret *dao.SecretDAO) (*dao.SecretDAO, error)
	// DeleteSecret deletes a secret
	DeleteSecret(ctx context.Context, id uuid.UUID) error
	// CreatePasswordSecret creates a new password secret in the data warehouse
	CreatePasswordSecret(ctx context.Context, secret *dao.PasswordSecretDAO) (*dao.PasswordSecretDAO, error)
	// GetPasswordSecretByID selects a password secret by id
	GetPasswordSecretByID(ctx context.Context, id uuid.UUID) (*dao.PasswordSecretDAO, error)
	// UpdatePasswordSecret updates a password secret
	UpdatePasswordSecret(ctx context.Context, secret *dao.PasswordSecretDAO) (*dao.PasswordSecretDAO, error)
	// CreateTextSecret creates a new text secret in the data warehouse
	CreateTextSecret(ctx context.Context, secret *dao.TextSecretDAO) (*dao.TextSecretDAO, error)
	// GetTextSecretByID selects a text secret by id
	GetTextSecretByID(ctx context.Context, id uuid.UUID) (*dao.TextSecretDAO, error)
	// UpdateTextSecret updates a text secret
	UpdateTextSecret(ctx context.Context, secret *dao.TextSecretDAO) (*dao.TextSecretDAO, error)
}
