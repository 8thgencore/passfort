package paseto

import (
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/8thgencore/passfort/internal/config"
	"github.com/8thgencore/passfort/internal/domain"
	"github.com/8thgencore/passfort/internal/service"
	"github.com/8thgencore/passfort/internal/service/adapters/cache"
)

/**
 * Token implements service.TokenService interface
 * and provides an access to the paseto library
 */
type Token struct {
	token    *paseto.Token
	key      *paseto.V4SymmetricKey
	parser   *paseto.Parser
	duration time.Duration
	cache    cache.CacheRepository
}

// New creates a new paseto instance
func New(config *config.Token, cache cache.CacheRepository) (service.TokenService, error) {
	durationStr := config.Duration
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return nil, domain.ErrTokenDuration
	}

	token := paseto.NewToken()
	key := paseto.NewV4SymmetricKey()
	parser := paseto.NewParser()

	return &Token{
		&token,
		&key,
		&parser,
		duration,
		cache,
	}, nil
}
