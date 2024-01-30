package middleware

import (
	"github.com/8thgencore/passfort/internal/delivery/http/helper"
	"github.com/8thgencore/passfort/internal/delivery/http/response"
	"github.com/8thgencore/passfort/internal/domain"
	"github.com/gin-gonic/gin"
)

// AdminMiddleware is a middleware to check if the user is an admin
func AdminMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		payload := helper.GetAuthPayload(ctx, authorizationPayloadKey)

		isAdmin := payload.Role == domain.AdminRole
		if !isAdmin {
			err := domain.ErrForbidden
			response.HandleAbort(ctx, err)
			return
		}

		ctx.Next()
	}
}
