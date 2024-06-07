package domain

import (
	"time"

	"github.com/google/uuid"
)

// SecretTypeEnum is an enum for secrets's type
type SecretTypeEnum string

const (
	PasswordSecretType SecretTypeEnum = "password"
	TextSecretType     SecretTypeEnum = "text"
	FileSecretType     SecretTypeEnum = "file"
)

type Secret struct {
	ID             uuid.UUID
	CollectionID   uuid.UUID
	SecretType     SecretTypeEnum
	Name           string
	Description    string
	CreatedBy      uuid.UUID
	UpdatedBy      uuid.UUID
	CreatedAt      time.Time
	UpdatedAt      time.Time
	LinkedSecretId uuid.UUID
	PasswordSecret *PasswordSecret
	TextSecret     *TextSecret
}

// PasswordSecret represents a password secret.
type PasswordSecret struct {
	ID       uuid.UUID
	URL      string
	Login    string
	Password string
}

// TextSecret represents a text secret.
type TextSecret struct {
	ID   uuid.UUID
	Text string
}
