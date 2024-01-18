package domain

// UserRole is an enum for user's role
type UserRoleEnum string

// UserRole enum values
const (
	AdminRole UserRoleEnum = "admin"
	UserRole  UserRoleEnum = "user"
)
