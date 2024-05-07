package token

import (
	"time"

	"github.com/8thgencore/passfort/internal/service/adapters/cache"
)

// TokenService handles operations related to tokens.
type TokenService struct {
	signingKey      string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	cache           cache.CacheRepository
}

// New creates a new instance of TokenService.
func New(
	signingKey string,
	accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration,
	cache cache.CacheRepository,
) *TokenService {
	return &TokenService{
		signingKey:      signingKey,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
		cache:           cache,
	}
}
