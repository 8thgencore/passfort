package domain

import (
	"github.com/google/uuid"
)

// UserClaims is an entity that represents the payload of the token
type UserClaims struct {
	ID     uuid.UUID
	UserID uuid.UUID
	Role   UserRoleEnum
}
