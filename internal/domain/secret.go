package domain

import (
	"time"

	"github.com/google/uuid"
)

// SecretTypeEnum is an enum for secrets's type
type SecretTypeEnum string

// SecretTypeEnum enum values
const (
	Password SecretTypeEnum = "password"
	Text     SecretTypeEnum = "text"
	File     SecretTypeEnum = "file"
)

type Secret struct {
	ID           uuid.UUID
	CollectionID uuid.UUID
	SecretType   SecretTypeEnum
	CreatedAt    time.Time
	UpdatedAt    time.Time
	CreatedBy    uuid.UUID
	UpdatedBy    uuid.UUID
}
