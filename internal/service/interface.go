package service
import (
	"context"

	"github.com/8thgencore/passfort/internal/domain"
)

// UserService is an interface for interacting with user-related business logic
type UserService interface {
	// Register registers a new user
	Register(ctx context.Context, user *domain.User) (*domain.User, error)
	// GetUser returns a user by id
	GetUser(ctx context.Context, id uint64) (*domain.User, error)
	// ListUsers returns a list of users with pagination
	ListUsers(ctx context.Context, skip, limit uint64) ([]domain.User, error)
	// UpdateUser updates a user
	UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	// DeleteUser deletes a user
	DeleteUser(ctx context.Context, id uint64) error
}
