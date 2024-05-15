package middleware

import (
	"github.com/8thgencore/passfort/internal/delivery/http/response"
	"github.com/8thgencore/passfort/internal/domain"
	"github.com/8thgencore/passfort/internal/service"
	"github.com/gin-gonic/gin"
)

// MasterPasswordMiddleware is a middleware to check if the master password is activated recently
func MasterPasswordMiddleware(masterPasswordService service.MasterPasswordService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authPayload, exists := ctx.Get(AuthorizationPayloadKey)
		if !exists {
			response.HandleAbort(ctx, domain.ErrUnauthorized)
			return
		}

		payload, ok := authPayload.(*domain.UserClaims)
		if !ok {
			response.HandleAbort(ctx, domain.ErrUnauthorized)
			return
		}

		userID := payload.UserID

		activated, err := masterPasswordService.IsMasterPasswordActivated(ctx, userID)
		if err != nil {
			response.HandleAbort(ctx, err)
			return
		}
		if !activated {
			err := domain.ErrMasterPasswordActivationExpired
			response.HandleAbort(ctx, err)
			return
		}

		ctx.Next()
	}
}
