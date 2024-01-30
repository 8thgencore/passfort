package helper

import (
	"github.com/8thgencore/passfort/internal/domain"
	"github.com/go-playground/validator/v10"
)

// UserRoleValidator is a custom validator for validating user roles
var UserRoleValidator validator.Func = func(fl validator.FieldLevel) bool {
	userRole := fl.Field().Interface().(domain.UserRoleEnum)

	switch userRole {
	case domain.AdminRole, domain.UserRole:
		return true
	default:
		return false
	}
}
