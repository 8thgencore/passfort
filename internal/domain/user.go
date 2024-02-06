package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	Name      string
	Email     string
	Password  string
	Role      UserRoleEnum
	CreatedAt time.Time
	UpdatedAt time.Time
}
