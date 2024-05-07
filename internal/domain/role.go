package domain

import "fmt"

// UserRole is an enum for user's role
type UserRoleEnum string

// UserRole enum values
const (
	AdminRole UserRoleEnum = "admin"
	UserRole  UserRoleEnum = "user"
)

// ParseUserRoleEnum parses a string into UserRoleEnum
func ParseUserRoleEnum(roleStr string) (UserRoleEnum, error) {
	switch roleStr {
	case string(AdminRole):
		return AdminRole, nil
	case string(UserRole):
		return UserRole, nil
	default:
		return "", fmt.Errorf("invalid role: %s", roleStr)
	}
}
