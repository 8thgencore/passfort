package paseto

import (
	"context"
	"time"

	"github.com/8thgencore/passfort/internal/domain"
	"github.com/8thgencore/passfort/pkg/util"
	"github.com/google/uuid"
)

// CreateToken creates a new paseto token
func (pt *Token) CreateToken(user *domain.User) (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", domain.ErrTokenCreation
	}

	payload := &domain.TokenPayload{
		ID:     id,
		UserID: user.ID,
		Role:   user.Role,
	}

	err = pt.token.Set("payload", payload)
	if err != nil {
		return "", domain.ErrTokenCreation
	}

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(pt.duration)

	pt.token.SetIssuedAt(issuedAt)
	pt.token.SetNotBefore(issuedAt)
	pt.token.SetExpiration(expiredAt)

	token := pt.token.V4Encrypt(*pt.key, nil)

	return token, nil
}

// VerifyToken verifies the paseto token
func (pt *Token) VerifyToken(token string) (*domain.TokenPayload, error) {
	var payload *domain.TokenPayload

	parsedToken, err := pt.parser.ParseV4Local(*pt.key, token, nil)
	if err != nil {
		if err.Error() == "this token has expired" {
			return nil, domain.ErrExpiredToken
		}
		return nil, domain.ErrInvalidToken
	}

	err = parsedToken.Get("payload", &payload)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	return payload, nil
}

// CheckTokenRevoked checks if the token is invalidated or outdated
func (pt *Token) CheckTokenRevoked(ctx context.Context, token *domain.TokenPayload) (bool, error) {
	cacheKey := util.GenerateCacheKey("token", token.ID)
	_, err := util.Serialize(token)
	if err != nil {
		return false, domain.ErrInternal
	}

	// Check if the value exists in the cache for the given key
	exists, err := pt.cache.Exists(ctx, cacheKey)
	if err != nil {
		return false, domain.ErrInternal
	}

	return exists, nil
}
