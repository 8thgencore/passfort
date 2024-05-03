package middleware

import (
	"strings"

	"github.com/8thgencore/passfort/internal/delivery/http/response"
	"github.com/8thgencore/passfort/internal/domain"
	"github.com/8thgencore/passfort/internal/service"
	"github.com/gin-gonic/gin"
)

const (
	// authorizationHeaderKey is the key for authorization header in the request
	authorizationHeaderKey = "authorization"
	// authorizationType is the accepted authorization type
	authorizationType = "bearer"
	// AuthorizationPayloadKey is the key for authorization payload in the context
	AuthorizationPayloadKey = "authorization_payload"
)

// AuthMiddleware is a middleware to check if the user is authenticated
func AuthMiddleware(token service.TokenService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)

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
		if currentAuthorizationType != authorizationType {
			err := domain.ErrInvalidAuthorizationType
			response.HandleAbort(ctx, err)
			return
		}

		accessToken := fields[1]
		payload, err := token.VerifyToken(accessToken)
		if err != nil {
			response.HandleAbort(ctx, err)
			return
		}

		exists, err := token.CheckTokenRevoked(ctx, payload)
		if err != nil {
			response.HandleAbort(ctx, err)
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
