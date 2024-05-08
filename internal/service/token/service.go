package token

import (
	"log/slog"
	"time"

	"github.com/8thgencore/passfort/internal/service/adapters/cache"
)

// TokenService handles operations related to tokens.
type TokenService struct {
	log             *slog.Logger
	signingKey      string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	cache           cache.CacheRepository
}

// New creates a new instance of TokenService.
func NewTokenService(
	log *slog.Logger,
	signingKey string,
	accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration,
	cache cache.CacheRepository,
) *TokenService {
	return &TokenService{
		log:             log,
		signingKey:      signingKey,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
		cache:           cache,
	}
}
