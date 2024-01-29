package token

import (
	"time"

	"github.com/8thgencore/passfort/internal/config"
	"github.com/8thgencore/passfort/internal/domain"
	"github.com/8thgencore/passfort/internal/service"
	"golang.org/x/crypto/chacha20poly1305"
)

/**
 * TokenService implements service.TokenService interface
 * and provides an access to the paseto library
 */
type TokenService struct {
	symmetricKey []byte
	duration     time.Duration
}

// New creates a new paseto instance
func New(config *config.Token) (service.TokenService, error) {
	symmetricKey := config.SymmetricKey
	durationStr := config.Duration

	validSymmetricKey := len(symmetricKey) == chacha20poly1305.KeySize
	if !validSymmetricKey {
		return nil, domain.ErrInvalidTokenSymmetricKey
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return nil, err
	}

	return &TokenService{
		[]byte(symmetricKey),
		duration,
	}, nil
}
