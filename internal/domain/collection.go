package domain

import (
	"time"

	"github.com/google/uuid"
)

type Collection struct {
	ID          uuid.UUID
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
