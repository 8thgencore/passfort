package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID
	Name           string
	Email          string
	Password       string
	MasterPassword string
	IsVerified     bool
	Role           UserRoleEnum
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
