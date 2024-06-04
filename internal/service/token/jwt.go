package token

import (
	"context"
	"fmt"
	"time"

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

	var accessTokenTTL, refreshTokenTTL time.Duration
	accessTokenTTL = svc.accessTokenTTL
	refreshTokenTTL = svc.refreshTokenTTL

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(accessTokenTTL).Unix()
	claims["iat"] = time.Now().Unix()
	claims["id"] = tokenID
	claims["user_id"] = userID
	claims["role"] = role

	accessToken, err := token.SignedString([]byte(svc.signingKey))
	if err != nil {
		return "", "", err
	}

	// Reset the expiration time for the refresh token
	claims = token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(refreshTokenTTL).Unix()

	refreshToken, err := token.SignedString([]byte(svc.signingKey))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// ParseUserClaims parses the access token and returns the user claims.
func (svc *TokenService) ParseUserClaims(accessToken string) (*domain.UserClaims, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(svc.signingKey), nil
	})
	if err != nil {
		svc.log.Debug("Error parsing access token: %v", err)
		return &domain.UserClaims{}, domain.ErrExpiredToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		svc.log.Debug("Error getting user claims from access token")
		return &domain.UserClaims{}, domain.ErrInvalidToken
	}

	tokenID, err := uuid.Parse(fmt.Sprintf("%v", claims["id"]))
	if err != nil {
		svc.log.Debug("Error parsing token ID: %v", err)
		return &domain.UserClaims{}, domain.ErrInvalidToken
	}

	userID, err := uuid.Parse(fmt.Sprintf("%v", claims["user_id"]))
	if err != nil {
		svc.log.Debug("Error parsing user ID: %v", err)
		return &domain.UserClaims{}, domain.ErrInvalidToken
	}

	roleStr, ok := claims["role"].(string)
	if !ok {
		svc.log.Debug("Error getting role from claims")
		return &domain.UserClaims{}, domain.ErrInvalidToken
	}
	role, err := domain.ParseUserRoleEnum(roleStr)
	if err != nil {
		svc.log.Debug("Error parsing user role: %v", err)
		return &domain.UserClaims{}, domain.ErrInvalidToken
	}

	return &domain.UserClaims{
		ID:     tokenID,
		UserID: userID,
		Role:   role,
	}, nil
}

// RevokeToken revokes the specified JWT token.
func (svc *TokenService) RevokeToken(ctx context.Context, tokenID uuid.UUID) error {
	// Caching a revoked token
	cacheKey := util.GenerateCacheKey("token", tokenID)
	userSerialized, err := util.Serialize(tokenID)
	if err != nil {
		return domain.ErrInternal
	}

	err = svc.cache.Set(ctx, cacheKey, userSerialized, svc.refreshTokenTTL)
	if err != nil {
		return domain.ErrInternal
	}

	return nil
}

// CheckJWTTokenRevoked checks if the JWT token is revoked.
func (svc *TokenService) CheckJWTTokenRevoked(ctx context.Context, tokenID uuid.UUID) (bool, error) {
	cacheKey := util.GenerateCacheKey("token", tokenID)
	fmt.Println(cacheKey)
	// Check if the value exists in the cache for the given key
	exists, err := svc.cache.Exists(ctx, cacheKey)
	if err != nil {
		return false, domain.ErrInternal
	}

	return exists, nil
}
