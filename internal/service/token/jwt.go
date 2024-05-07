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

// GenerateAccessToken generates a new JWT access token based on the provided user claims.
func (ts *TokenService) GenerateAccessToken(userID uuid.UUID, role domain.UserRoleEnum) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	tokenId, err := uuid.NewRandom()
	if err != nil {
		return "", domain.ErrTokenCreation
	}

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(ts.accessTokenTTL).Unix()
	claims["iat"] = time.Now().Unix()
	claims["id"] = tokenId
	claims["user_id"] = userID
	claims["role"] = role

	tokenString, err := token.SignedString([]byte(ts.signingKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GenerateRefreshToken generates a new refresh token.
func (ts *TokenService) GenerateRefreshToken(userID uuid.UUID) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	tokenId, err := uuid.NewRandom()
	if err != nil {
		return "", domain.ErrTokenCreation
	}

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(ts.refreshTokenTTL).Unix()
	claims["iat"] = time.Now().Unix()
	claims["id"] = tokenId
	claims["user_id"] = userID

	tokenString, err := token.SignedString([]byte(ts.signingKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
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