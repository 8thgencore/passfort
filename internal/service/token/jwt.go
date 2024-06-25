package token

import (
	"context"
	"fmt"
	"time"

	"errors"

	"github.com/8thgencore/passfort/internal/domain"
	"github.com/8thgencore/passfort/pkg/util"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// GenerateToken generates a new JWT token pair based on the provided user claims.
func (svc *TokenService) GenerateToken(userID uuid.UUID, role domain.UserRoleEnum) (string, string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	tokenID, err := uuid.NewRandom()
	if err != nil {
		return "", "", domain.ErrTokenCreation
	}

	now := time.Now()
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = now.Add(svc.accessTokenTTL).Unix()
	claims["iat"] = now.Unix()
	claims["id"] = tokenID
	claims["user_id"] = userID
	claims["role"] = role

	accessToken, err := token.SignedString([]byte(svc.signingKey))
	if err != nil {
		return "", "", errors.New("failed to sign access token")
	}

	claims["exp"] = now.Add(svc.refreshTokenTTL).Unix()

	refreshToken, err := token.SignedString([]byte(svc.signingKey))
	if err != nil {
		return "", "", errors.New("failed to sign refresh token")
	}

	return accessToken, refreshToken, nil
}

// ParseUserClaims parses the access token and returns the user claims.
func (svc *TokenService) ParseUserClaims(accessToken string) (*domain.UserClaims, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(svc.signingKey), nil
	})
	if err != nil {
		svc.log.Debug("Error parsing access token: %v", err)
		return nil, domain.ErrExpiredToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		svc.log.Debug("Error getting user claims from access token")
		return nil, domain.ErrInvalidToken
	}

	tokenID, err := uuid.Parse(fmt.Sprintf("%v", claims["id"]))
	if err != nil {
		svc.log.Debug("Error parsing token ID: %v", err)
		return nil, domain.ErrInvalidToken
	}

	userID, err := uuid.Parse(fmt.Sprintf("%v", claims["user_id"]))
	if err != nil {
		svc.log.Debug("Error parsing user ID: %v", err)
		return nil, domain.ErrInvalidToken
	}

	roleStr, ok := claims["role"].(string)
	if !ok {
		svc.log.Debug("Error getting role from claims")
		return nil, domain.ErrInvalidToken
	}
	role, err := domain.ParseUserRoleEnum(roleStr)
	if err != nil {
		svc.log.Debug("Error parsing user role: %v", err)
		return nil, domain.ErrInvalidToken
	}

	return &domain.UserClaims{
		ID:     tokenID,
		UserID: userID,
		Role:   role,
	}, nil
}

// RevokeToken revokes the specified JWT token.
func (svc *TokenService) RevokeToken(ctx context.Context, tokenID uuid.UUID) error {
	cacheKey := util.GenerateCacheKey("token", tokenID)
	userSerialized, err := util.Serialize(tokenID)
	if err != nil {
		return errors.New("failed to serialize token ID")
	}

	err = svc.cache.Set(ctx, cacheKey, userSerialized, svc.refreshTokenTTL)
	if err != nil {
		return errors.New("failed to cache revoked token")
	}

	return nil
}

// CheckJWTTokenRevoked checks if the JWT token is revoked.
func (svc *TokenService) CheckJWTTokenRevoked(ctx context.Context, tokenID uuid.UUID) (bool, error) {
	cacheKey := util.GenerateCacheKey("token", tokenID)

	exists, err := svc.cache.Exists(ctx, cacheKey)
	if err != nil {
		return false, errors.New("failed to check if token is revoked")
	}

	return exists, nil
}
