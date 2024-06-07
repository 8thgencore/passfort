package helper

import (
	"strconv"

	"github.com/8thgencore/passfort/internal/delivery/http/response"
	"github.com/8thgencore/passfort/internal/domain"
	"github.com/gin-gonic/gin"
)

// GetAuthPayload is a helper function to get the auth payload from the context
func GetAuthPayload(ctx *gin.Context, key string) *domain.UserClaims {
	return ctx.MustGet(key).(*domain.UserClaims)
}

// GetEncryptionKey is a helper function to get the encryption key from the context
func GetEncryptionKey(ctx *gin.Context, key string) string {
	return ctx.MustGet(key).(string)
}

// StringToUint64 is a helper function to convert a string to uint64
func StringToUint64(str string) (uint64, error) {
	num, err := strconv.ParseUint(str, 10, 64)

	return num, err
}

// ToMap is a helper function to add meta and data to a map
func ToMap(m response.Meta, data any, key string) map[string]any {
	return map[string]any{
		"meta": m,
		key:    data,
	}
}
