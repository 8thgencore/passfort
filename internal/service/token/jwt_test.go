package token_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"log/slog"

	"github.com/8thgencore/passfort/internal/domain"
	"github.com/8thgencore/passfort/internal/service/adapters/cache/mocks"
	"github.com/8thgencore/passfort/internal/service/token"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	cache           = new(mocks.CacheRepository)
	log             = slog.Default()
	signingKey      = "test-signing-key"
	accessTokenTTL  = 15 * time.Minute
	refreshTokenTTL = 7 * 24 * time.Hour
)

func newTokenService() *token.TokenService {
	return token.NewTokenService(log, signingKey, accessTokenTTL, refreshTokenTTL, cache)
}

func TestGenerateToken(t *testing.T) {
	ts := newTokenService()

	t.Run("Successfully generates token pair", func(t *testing.T) {
		userID := uuid.New()
		role := domain.UserRole

		accessToken, refreshToken, err := ts.GenerateToken(userID, role)

		assert.NoError(t, err)
		assert.NotEmpty(t, accessToken)
		assert.NotEmpty(t, refreshToken)

		// Verify the tokens
		accessTokenClaims := jwt.MapClaims{}
		_, err = jwt.ParseWithClaims(accessToken, accessTokenClaims, func(token *jwt.Token) (interface{}, error) {
			return []byte(signingKey), nil
		})
		assert.NoError(t, err)

		refreshTokenClaims := jwt.MapClaims{}
		_, err = jwt.ParseWithClaims(refreshToken, refreshTokenClaims, func(token *jwt.Token) (interface{}, error) {
			return []byte(signingKey), nil
		})
		assert.NoError(t, err)
	})
}

func TestParseUserClaims(t *testing.T) {
	ts := newTokenService()

	t.Run("Successfully parses user claims", func(t *testing.T) {
		userID := uuid.New()
		role := domain.UserRole

		accessToken, _, err := ts.GenerateToken(userID, role)
		assert.NoError(t, err)

		claims, err := ts.ParseUserClaims(accessToken)
		assert.NoError(t, err)
		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, role, claims.Role)
	})

	t.Run("Returns error when token parsing fails", func(t *testing.T) {
		accessToken := "invalid-token"

		_, err := ts.ParseUserClaims(accessToken)
		assert.Error(t, err)
	})

	t.Run("Returns error when token claims are not a map", func(t *testing.T) {
		token := jwt.New(jwt.SigningMethodHS256)
		token.Claims = jwt.MapClaims{}

		accessToken, err := token.SignedString([]byte(signingKey))
		assert.NoError(t, err)

		_, err = ts.ParseUserClaims(accessToken)
		assert.Error(t, err)
	})

	t.Run("Returns error when token ID is invalid", func(t *testing.T) {
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["id"] = "invalid-id"

		accessToken, err := token.SignedString([]byte(signingKey))
		assert.NoError(t, err)

		_, err = ts.ParseUserClaims(accessToken)
		assert.Error(t, err)
	})

	t.Run("Returns error when user ID is invalid", func(t *testing.T) {
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["user_id"] = "invalid-user-id"

		accessToken, err := token.SignedString([]byte(signingKey))
		assert.NoError(t, err)

		_, err = ts.ParseUserClaims(accessToken)
		assert.Error(t, err)
	})

	t.Run("Returns error when role is invalid", func(t *testing.T) {
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["role"] = "invalid-role"

		accessToken, err := token.SignedString([]byte(signingKey))
		assert.NoError(t, err)

		_, err = ts.ParseUserClaims(accessToken)
		assert.Error(t, err)
	})
}

func TestRevokeToken(t *testing.T) {
	ts := newTokenService()

	t.Run("Successfully revokes token", func(t *testing.T) {
		tokenID := uuid.New()

		cache.On("Set", mock.Anything, mock.Anything, mock.Anything, refreshTokenTTL).Return(nil)

		err := ts.RevokeToken(context.Background(), tokenID)
		assert.NoError(t, err)

		cache.AssertExpectations(t)
		cache.ExpectedCalls = nil
	})

	t.Run("Returns error when cache set fails", func(t *testing.T) {
		tokenID := uuid.New()

		cache.On("Set", mock.Anything, mock.Anything, mock.Anything, refreshTokenTTL).Return(errors.New("failed to set cache"))

		_ = ts.RevokeToken(context.Background(), tokenID)
		assert.Error(t, domain.ErrInternal)

		cache.AssertExpectations(t)
		cache.ExpectedCalls = nil
	})
}

func TestCheckJWTTokenRevoked(t *testing.T) {
	ts := newTokenService()

	t.Run("Token is revoked", func(t *testing.T) {
		tokenID := uuid.New()

		cache.On("Exists", mock.Anything, mock.Anything).Return(true, nil)

		revoked, err := ts.CheckJWTTokenRevoked(context.Background(), tokenID)
		assert.NoError(t, err)
		assert.True(t, revoked)

		cache.AssertExpectations(t)
		cache.ExpectedCalls = nil
	})

	t.Run("Token is not revoked", func(t *testing.T) {
		tokenID := uuid.New()

		cache.On("Exists", mock.Anything, mock.Anything).Return(false, nil)

		revoked, err := ts.CheckJWTTokenRevoked(context.Background(), tokenID)

		assert.NoError(t, err)
		assert.False(t, revoked)

		cache.AssertExpectations(t)
		cache.ExpectedCalls = nil
	})

	t.Run("Returns error when cache exists check fails", func(t *testing.T) {
		tokenID := uuid.New()

		cache.On("Exists", mock.Anything, mock.Anything).Return(false, errors.New("failed to check cache"))

		_, _ = ts.CheckJWTTokenRevoked(context.Background(), tokenID)
		assert.Error(t, domain.ErrInternal)

		cache.AssertExpectations(t)
		cache.ExpectedCalls = nil
	})
}
