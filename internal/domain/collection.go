package domain

import (
	"time"

	"github.com/google/uuid"
)

type Collection struct {
	ID          uuid.UUID
	Name        string
	Description string
	CreatedBy   uuid.UUID
	UpdatedBy   uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
