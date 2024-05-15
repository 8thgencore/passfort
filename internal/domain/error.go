package domain

import (
	"errors"
	"strings"
)

var (
	// Internal Errors
	// ErrInternal is an error for when an internal service fails to process the request
	ErrInternal = errors.New("internal error")
	// ErrFailedToSendEmail is an error for when sending an email fails
	ErrFailedToSendEmail = errors.New("failed to send email")

	// Data Errors
	// ErrDataNotFound is an error for when requested data is not found
	ErrDataNotFound = errors.New("data not found")
	// ErrNoUpdatedData is an error for when no data is provided to update
	ErrNoUpdatedData = errors.New("no data to update")
	// ErrConflictingData is an error for when data conflicts with existing data
	ErrConflictingData = errors.New("data conflicts with existing data")

	// Token Errors
	// ErrTokenDuration is an error for when the token duration format is invalid
	ErrTokenDuration = errors.New("invalid token duration format")
	// ErrTokenCreation is an error for when the token creation fails
	ErrTokenCreation = errors.New("error creating token")
	// ErrExpiredToken is an error for when the access token is expired
	ErrExpiredToken = errors.New("access token has expired")
	// ErrInvalidToken is an error for when the access token is invalid
	ErrInvalidToken = errors.New("access token is invalid")
	// ErrInvalidRefreshToken is an error for when the refresh token is invalid
	ErrInvalidRefreshToken = errors.New("refresh token is invalid")

	// Authentication Errors
	// ErrInvalidCredentials is an error for when the credentials are invalid
	ErrInvalidCredentials = errors.New("invalid email or password")
	// ErrPasswordsDoNotMatch is an error for when the provided passwords do not match
	ErrPasswordsDoNotMatch = errors.New("provided passwords do not match")
	// ErrInvalidOTP is an error for when the OTP (One-Time Password) is invalid
	ErrInvalidOTP = errors.New("invalid OTP")
	// ErrOTPAlreadySent is an error for when an OTP (One-Time Password) has already been sent to the user
	ErrOTPAlreadySent = errors.New("an OTP has already been sent to the user")

	// Authorization Errors
	// ErrEmptyAuthorizationHeader is an error for when the authorization header is empty
	ErrEmptyAuthorizationHeader = errors.New("authorization header is not provided")
	// ErrInvalidAuthorizationHeader is an error for when the authorization header is invalid
	ErrInvalidAuthorizationHeader = errors.New("authorization header format is invalid")
	// ErrInvalidAuthorizationType is an error for when the authorization type is invalid
	ErrInvalidAuthorizationType = errors.New("authorization type is not supported")
	// ErrUnauthorized is an error for when the user is unauthorized
	ErrUnauthorized = errors.New("user is unauthorized to access the resource")
	// ErrForbidden is an error for when the user is forbidden to access the resource
	ErrForbidden = errors.New("user is forbidden to access the resource")

	// User Errors
	// ErrUserNotVerified is an error for when a user is not verified
	ErrUserNotVerified = errors.New("user not verified")
	// ErrDeleteOwnAccount is an error for when a user tries to delete their own account
	ErrDeleteOwnAccount = errors.New("you cannot delete your own account")

	// Master Password Errors
	// ErrMasterPasswordActivationExpired is an error for when master password validation has expired
	ErrMasterPasswordActivationExpired = errors.New("master password validation has expired")
	// ErrMasterPasswordNotSet is an error for when a master password has not been set by the user
	ErrMasterPasswordNotSet = errors.New("master password has not been set")
	// ErrInvalidMasterPassword is an error for when the master password provided is invalid
	ErrInvalidMasterPassword = errors.New("invalid master password")
	// ErrMasterPasswordAlreadyExists is an error for when a master password already exists for the user
	ErrMasterPasswordAlreadyExists = errors.New("master password already exists")
)

// IsUniqueConstraintViolationError checks if the error is a unique constraint violation error
func IsUniqueConstraintViolationError(err error) bool {
	return strings.Contains(err.Error(), "23505")
}
