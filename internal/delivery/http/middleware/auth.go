package middleware

import (
	"strings"

	"github.com/8thgencore/passfort/internal/delivery/http/response"
	"github.com/8thgencore/passfort/internal/domain"
	"github.com/8thgencore/passfort/internal/service/token"
	"github.com/gin-gonic/gin"
)

const (
	// AuthorizationHeaderKey is the key for authorization header in the request
	AuthorizationHeaderKey = "authorization"
	// AuthorizationType is the accepted authorization type
	AuthorizationType = "bearer"
	// AuthorizationPayloadKey is the key for authorization payload in the context
	AuthorizationPayloadKey = "authorization_payload"
)

// AuthMiddleware is a middleware to check if the user is authenticated
func AuthMiddleware(tokenService token.TokenService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(AuthorizationHeaderKey)

		isEmpty := len(authorizationHeader) == 0
		if isEmpty {
			err := domain.ErrEmptyAuthorizationHeader
			response.HandleAbort(ctx, err)
			return
		}

		fields := strings.Fields(authorizationHeader)
		isValid := len(fields) == 2
		if !isValid {
			err := domain.ErrInvalidAuthorizationHeader
			response.HandleAbort(ctx, err)
			return
		}

		currentAuthorizationType := strings.ToLower(fields[0])
		if currentAuthorizationType != AuthorizationType {
			err := domain.ErrInvalidAuthorizationType
			response.HandleAbort(ctx, err)
			return
		}	

		payload, err := tokenService.ParseUserClaims(fields[1])
		if err != nil {
			response.HandleAbort(ctx, err)
			return
		}

		exists, err := tokenService.CheckJWTTokenRevoked(ctx, payload.ID)
		if err != nil {
			response.HandleAbort(ctx, domain.ErrInvalidToken)
			return
		}
		if exists {
			err := domain.ErrUnauthorized
			response.HandleAbort(ctx, err)
			return
		}

		ctx.Set(AuthorizationPayloadKey, payload)
		ctx.Next()
	}
}
