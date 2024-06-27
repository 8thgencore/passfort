package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                uuid.UUID    `json:"id"`
	Name              string       `json:"name"`
	Email             string       `json:"email"`
	Password          string       `json:"-"` // Hide the password field
	MasterPassword    string       `json:"-"` // Hide the master password field
	MasterPasswordSet bool         `json:"master_password_set"`
	Salt              []byte       `json:"-"` // Hide the salt field
	IsVerified        bool         `json:"is_verified"`
	Role              UserRoleEnum `json:"role"`
	CreatedAt         time.Time    `json:"created_at"`
	UpdatedAt         time.Time    `json:"updated_at"`
}
