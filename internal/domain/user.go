package domain

import "time"

type User struct {
	ID        uint64
	Name      string
	Email     string
	Password  string
	Role      UserRoleEnum
	CreatedAt time.Time
	UpdatedAt time.Time
}
