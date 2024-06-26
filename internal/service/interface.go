package service

import (
	"context"

	"github.com/8thgencore/passfort/internal/domain"
	"github.com/google/uuid"
)

// TokenService represents a service for handling tokens.
type TokenService interface {
	// GenerateToken generates a new JWT token pair based on the provided user claims.
	GenerateToken(userID uuid.UUID, role domain.UserRoleEnum) (string, string, error)
	// ParseUserClaims parses the access token and returns the user claims.
	ParseUserClaims(accessToken string) (*domain.UserClaims, error)
	// RevokeToken revokes the specified JWT token.
	RevokeToken(ctx context.Context, token uuid.UUID) error
	// CheckJWTTokenRevoked checks if the JWT token is revoked.
	CheckJWTTokenRevoked(ctx context.Context, token uuid.UUID) (bool, error)
}

// OtpService
type OtpService interface {
	// GenerateOTP generates a new OTP for the given user ID
	GenerateOTP(ctx context.Context, userID uuid.UUID) (string, error)
	// VerifyOTP verifies if the provided OTP is valid for the given user ID
	VerifyOTP(ctx context.Context, userID uuid.UUID, otp2 string) error
	// CheckCacheForKey checks if a value exists in the cache for the given user ID
	CheckCacheForKey(ctx context.Context, userID uuid.UUID) (bool, error)
}

// AuthService is an interface for interacting with user authentication-related business logic
type AuthService interface {
	// Login authenticates a user by email and password and returns a token
	Login(ctx context.Context, email, password string) (string, string, error)

	// Register registers a new user
	Register(ctx context.Context, user *domain.User) (*domain.User, error)
	// ConfirmRegistration confirms user registration with OTP code
	ConfirmRegistration(ctx context.Context, email, otp string) error
	// RequestNewRegistrationCode requests a new registration confirmation code for a user
	RequestNewRegistrationCode(ctx context.Context, email string) error

	// Logout invalidates the access token, logging the user out
	Logout(ctx context.Context, token *domain.UserClaims) error

	// RefreshToken refreshes the access token for the user
	RefreshToken(ctx context.Context, refreshToken string) (string, string, error)

	// ChangePassword changes the password for the authenticated user
	ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error

	// ForgotPassword initiates the process of resetting a forgotten password
	ForgotPassword(ctx context.Context, email string) error
	// ResetPassword confirms password reset with OTP code
	ResetPassword(ctx context.Context, email, newPassword, otp string) error
}

// MasterPasswordService is an interface for interacting with master password-related business logic
type MasterPasswordService interface {
	// MasterPasswordExists checks if a master password already exists for the given user
	MasterPasswordExists(ctx context.Context, userID uuid.UUID) (bool, error)
	// SaveMasterPassword saves or updates the master password for the given user
	SaveMasterPassword(ctx context.Context, userID uuid.UUID, password string) error
	// ChangeMasterPassword changes the master password for the given user.
	ChangeMasterPassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error
	// ActivateMasterPassword validates the master password for the given user
	ActivateMasterPassword(ctx context.Context, userID uuid.UUID, password string) error
	// GetEncryptionKey is required to encrypt or decrypt the password
	GetEncryptionKey(ctx context.Context, userID uuid.UUID) ([]byte, error)
}

// UserService is an interface for interacting with user-related business logic
type UserService interface {
	// GetUser returns a user by id
	GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	// ListUsers returns a list of users with pagination
	ListUsers(ctx context.Context, skip, limit uint64) ([]domain.User, error)
	// UpdateUser updates a user
	UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	// DeleteUser deletes a user
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

// CollectionService is an interface for interacting with collection-related business logic
type CollectionService interface {
	// CreateCollection inserts a new collection into the database
	CreateCollection(ctx context.Context, userID uuid.UUID, collection *domain.Collection) (*domain.Collection, error)
	// ListCollectionsByUserID returns a list of collections by user id with pagination
	ListCollectionsByUserID(ctx context.Context, userID uuid.UUID, skip, limit uint64) ([]domain.Collection, error)
	// GetCollection returns a collection by id
	GetCollection(ctx context.Context, userID, collectionID uuid.UUID) (*domain.Collection, error)
	// UpdateCollection updates a collection
	UpdateCollection(ctx context.Context, userID uuid.UUID, collection *domain.Collection) (*domain.Collection, error)
	// DeleteCollection deletes a collection
	DeleteCollection(ctx context.Context, userID, collectionID uuid.UUID) error
}

// SecretService is an interface for interacting with secret-related business logic
type SecretService interface {
	// CreateSecret inserts a new secret into the database
	CreateSecret(ctx context.Context, userID uuid.UUID, secret *domain.Secret, encryptionKey []byte) (*domain.Secret, error)
	// ListSecretsByCollectionID returns a list of secrets by collection ID with pagination
	ListSecretsByCollectionID(ctx context.Context, userID uuid.UUID, collectionID uuid.UUID, skip, limit uint64) ([]domain.Secret, error)
	// GetSecret returns a secret by id
	GetSecret(ctx context.Context, userID, collectionID, secretID uuid.UUID, encryptionKey []byte) (*domain.Secret, error)
	// UpdateSecret updates a secret
	UpdateSecret(ctx context.Context, userID, collectionID uuid.UUID, secret *domain.Secret, encryptionKey []byte) (*domain.Secret, error)
	// DeleteSecret deletes a secret
	DeleteSecret(ctx context.Context, userID, collectionID, secretID uuid.UUID) error
	// ReencryptAllSecrets reencrypt all secrets
	ReencryptAllSecrets(ctx context.Context, userID uuid.UUID, oldEncryptionKey, newEncryptionKey []byte) error
}
