package dao

import (
	"time"

	"github.com/google/uuid"
)

type CollectionDAO struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	CreatedBy    uuid.UUID `db:"created_by"`
	UpdatedBy    uuid.UUID `db:"updated_by"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
