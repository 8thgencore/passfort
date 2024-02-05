package service

import (
	"context"

	"github.com/8thgencore/passfort/internal/domain"
)

// TokenService is an interface for interacting with token-related business logic
type TokenService interface {
	// CreateToken creates a new token for a given user
	CreateToken(user *domain.User) (string, error)
	// VerifyToken verifies the token and returns the payload
	VerifyToken(token string) (*domain.TokenPayload, error)
}

// UserService is an interface for interacting with user authentication-related business logic
type AuthService interface {
	// Register registers a new user
	Register(ctx context.Context, user *domain.User) (*domain.User, error)
	// Login authenticates a user by email and password and returns a token
	Login(ctx context.Context, email, password string) (string, error)
	// ChangePassword changes the password for the authenticated user
	ChangePassword(ctx context.Context, userID uint64, oldPassword, newPassword string) error
}

// UserService is an interface for interacting with user-related business logic
type UserService interface {
	// GetUser returns a user by id
	GetUser(ctx context.Context, id uint64) (*domain.User, error)
	// ListUsers returns a list of users with pagination
	ListUsers(ctx context.Context, skip, limit uint64) ([]domain.User, error)
	// UpdateUser updates a user
	UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	// DeleteUser deletes a user
	DeleteUser(ctx context.Context, id uint64) error
}

// CollectionService is an interface for interacting with collection-related business logic
type CollectionService interface {
	// CreateCollection inserts a new collection into the database
	CreateCollection(ctx context.Context, userID uint64, collection *domain.Collection) (*domain.Collection, error)
	// ListCollectionsByUserID returns a list of collections by user id with pagination
	ListCollectionsByUserID(ctx context.Context, userID, skip, limit uint64) ([]domain.Collection, error)
	// GetCollection returns a collection by id
	GetCollection(ctx context.Context, userID, collectionID uint64) (*domain.Collection, error)
	// UpdateCollection updates a collection
	UpdateCollection(ctx context.Context, userID uint64, collection *domain.Collection) (*domain.Collection, error)
	// DeleteCollection deletes a collection
	DeleteCollection(ctx context.Context, userID, collectionID uint64) error
}
