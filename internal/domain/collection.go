package domain

import "time"

type Collection struct {
	ID          uint64
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
