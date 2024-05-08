package domain

import (
	"fmt"

	"github.com/google/uuid"
)

// UserClaims is an entity that represents the payload of the token
type UserClaims struct {
	ID     uuid.UUID
	UserID uuid.UUID
	Role   UserRoleEnum
}

// TokenEnum is an enum for user's role
type TokenEnum string

// Token enum values
const (
	AccessToken  TokenEnum = "access"
	RefreshToken TokenEnum = "refresh"
)

// ParseTokenEnum parses a string into TokenEnum
func ParseTokenEnum(token string) (TokenEnum, error) {
	switch token {
	case string(AccessToken):
		return AccessToken, nil
	case string(RefreshToken):
		return RefreshToken, nil
	default:
		return "", fmt.Errorf("invalid role: %s", token)
	}
}
