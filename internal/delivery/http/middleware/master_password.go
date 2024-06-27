package middleware

import (
	"github.com/8thgencore/passfort/internal/delivery/http/response"
	"github.com/8thgencore/passfort/internal/domain"
	"github.com/8thgencore/passfort/internal/service"
	"github.com/8thgencore/passfort/pkg/base64_util"
	"github.com/gin-gonic/gin"
)

const (
	// EncryptionKey is the key for encryption and decryption secrets
	EncryptionKey = "encryption_key"
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

		encryptionKey, err := masterPasswordService.GetEncryptionKey(ctx, payload.UserID)
		if err != nil {
			response.HandleAbort(ctx, err)
			return
		}
		if encryptionKey == nil {
			err := domain.ErrMasterPasswordActivationExpired
			response.HandleAbort(ctx, err)
			return
		}

		ctx.Set(EncryptionKey, base64_util.BytesToBase64(encryptionKey))
		ctx.Next()
	}
}
