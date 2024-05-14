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
func (ts *TokenService) GenerateToken(userID uuid.UUID, role domain.UserRoleEnum) (string, string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	tokenId, err := uuid.NewRandom()
	if err != nil {
		return "", "", domain.ErrTokenCreation
	}

	var accessTokenTTL, refreshTokenTTL time.Duration
	accessTokenTTL = ts.accessTokenTTL
	refreshTokenTTL = ts.refreshTokenTTL

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(accessTokenTTL).Unix()
	claims["iat"] = time.Now().Unix()
	claims["id"] = tokenId
	claims["user_id"] = userID
	claims["role"] = role

	accessToken, err := token.SignedString([]byte(ts.signingKey))
	if err != nil {
		return "", "", err
	}

	// Reset the expiration time for the refresh token
	claims = token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(refreshTokenTTL).Unix()

	refreshToken, err := token.SignedString([]byte(ts.signingKey))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// ParseUserClaims parses the access token and returns the user claims.
func (ts *TokenService) ParseUserClaims(accessToken string) (*domain.UserClaims, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(ts.signingKey), nil
	})
	if err != nil {
		return &domain.UserClaims{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return &domain.UserClaims{}, fmt.Errorf("error getting user claims from access token")
	}

	tokenID, err := uuid.Parse(fmt.Sprintf("%v", claims["id"]))
	if err != nil {
		return &domain.UserClaims{}, fmt.Errorf("error parsing token ID")
	}

	userID, err := uuid.Parse(fmt.Sprintf("%v", claims["user_id"]))
	if err != nil {
		return &domain.UserClaims{}, fmt.Errorf("error parsing user ID")
	}

	roleStr, ok := claims["role"].(string)
	if !ok {
		return &domain.UserClaims{}, fmt.Errorf("error getting role from claims")
	}
	role, err := domain.ParseUserRoleEnum(roleStr)
	if err != nil {
		return &domain.UserClaims{}, fmt.Errorf("error parsing user role")
	}

	return &domain.UserClaims{
		ID:     tokenID,
		UserID: userID,
		Role:   role,
	}, nil
}

// RevokeToken revokes the specified JWT token.
func (ts *TokenService) RevokeToken(ctx context.Context, token uuid.UUID) error {
	// Caching a revoked token
	cacheKey := util.GenerateCacheKey("token", token.ID)
	userSerialized, err := util.Serialize(token)
	if err != nil {
		return domain.ErrInternal
	}

	err = ts.cache.Set(ctx, cacheKey, userSerialized, ts.refreshTokenTTL)
	if err != nil {
		return domain.ErrInternal
	}

	return nil
}

// CheckJWTTokenRevoked checks if the JWT token is revoked.
func (ts *TokenService) CheckJWTTokenRevoked(ctx context.Context, tokenId uuid.UUID) (bool, error) {
	cacheKey := util.GenerateCacheKey("token", tokenId)
	_, err := util.Serialize(tokenId)
	if err != nil {
		return false, domain.ErrInternal
	}

	// Check if the value exists in the cache for the given key
	exists, err := ts.cache.Exists(ctx, cacheKey)
	if err != nil {
		return false, domain.ErrInternal
	}

	return exists, nil
}
