package token

import (
	"time"

	"github.com/8thgencore/passfort/internal/domain"
	"github.com/google/uuid"
)

// CreateToken creates a new paseto token
func (pt *TokenService) CreateToken(user *domain.User) (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", domain.ErrTokenCreation
	}

	payload := domain.TokenPayload{
		ID:        id,
		UserID:    user.ID,
		Role:      user.Role,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(pt.duration),
	}

	_ = payload
	// token, err := pt.paseto.Encrypt(pt.symmetricKey, payload, nil)
	// if err != nil {
	// 	return "", domain.ErrTokenCreation
	// }

	// return token, nil
	return "", nil
}

// VerifyToken verifies the paseto token
func (pt *TokenService) VerifyToken(token string) (*domain.TokenPayload, error) {
	var payload domain.TokenPayload

	// err := pt.paseto.Decrypt(token, pt.symmetricKey, &payload, nil)
	// if err != nil {
	// return nil, domain.ErrInvalidToken
	// }

	isExpired := time.Now().After(payload.ExpiredAt)
	if isExpired {
		return nil, domain.ErrExpiredToken
	}

	return &payload, nil
}
