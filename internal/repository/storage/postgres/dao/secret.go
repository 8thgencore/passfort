package dao

import (
	"time"

	"github.com/google/uuid"
)

// SecretDAO is a model of a secret in a data store.
type SecretDAO struct {
	ID             uuid.UUID   `db:"id"`
	CollectionID   uuid.UUID   `db:"collection_id"`
	SecretType     string      `db:"secret_type"`
	Name           string      `db:"name"`
	Description    string      `db:"description"`
	CreatedBy      uuid.UUID   `db:"created_by"`
	UpdatedBy      uuid.UUID   `db:"updated_by"`
	CreatedAt      time.Time   `db:"created_at"`
	UpdatedAt      time.Time   `db:"updated_at"`
	LinkedSecretId uuid.UUID   `db:"linked_secret_id"`
	LinkedSecret   interface{} // Holds either PasswordSecret or TextSecret
}

// PasswordSecretDAO represents a password secret.
type PasswordSecretDAO struct {
	ID       uuid.UUID `db:"id"`
	URL      string    `db:"url"`
	Login    string    `db:"login"`
	Password string    `db:"password"`
}

// TextSecretDAO represents a text secret.
type TextSecretDAO struct {
	ID   uuid.UUID `db:"id"`
	Text string    `db:"text"`
}
