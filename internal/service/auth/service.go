package auth

import (
	"github.com/8thgencore/passfort/internal/service"
	"github.com/8thgencore/passfort/internal/service/adapters/storage"
)

/**
 * AuthService implements service.AuthService interface
 * and provides an access to the user repository
 * and token service
 */
type AuthService struct {
	repo storage.UserRepository
	ts   service.TokenService
}

// NewAuthService creates a new auth service instance
func NewAuthService(repo storage.UserRepository, ts service.TokenService) *AuthService {
	return &AuthService{
		repo,
		ts,
	}
}
