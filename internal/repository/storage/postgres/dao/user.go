package dao

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type UserDAO struct {
	ID             uuid.UUID      `db:"id"`
	Name           string         `db:"name"`
	Email          string         `db:"email"`
	Password       string         `db:"password"`
	MasterPassword sql.NullString `db:"master_password"`
	IsVerified     bool           `db:"is_verified"`
	Role           string         `db:"role"`
	CreatedAt      time.Time      `db:"created_at"`
	UpdatedAt      time.Time      `db:"updated_at"`
}
