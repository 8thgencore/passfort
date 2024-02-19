package dao

import (
	"time"

	"github.com/google/uuid"
)

// SecretDAO is a model of a secret in a data store.
type SecretDAO struct {
	ID           uuid.UUID `db:"id"`
	CollectionID uuid.UUID `db:"collection_id"`
	SecretType   string    `db:"secret_type"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
	CreatedBy    uuid.UUID `db:"created_by"`
	UpdatedBy    uuid.UUID `db:"updated_by"`
}
