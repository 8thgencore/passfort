package response

import (
	"time"

	"github.com/8thgencore/passfort/internal/domain"
	"github.com/google/uuid"
)

// Response represents a Response body format
type Response struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Success"`
	Data    any    `json:"data,omitempty"`
}

// NewResponse is a helper function to create a response body
func NewResponse(success bool, message string, data any) Response {
	return Response{
		Success: success,
		Message: message,
		Data:    data,
	}
}

// Meta represents metadata for a paginated response
type Meta struct {
	Total uint64 `json:"total" example:"100"`
	Limit uint64 `json:"limit" example:"10"`
	Skip  uint64 `json:"skip" example:"0"`
}

// NewMeta is a helper function to create metadata for a paginated response
func NewMeta(total, limit, skip uint64) Meta {
	return Meta{
		Total: total,
		Limit: limit,
		Skip:  skip,
	}
}

// AuthResponse represents an authentication response body
type AuthResponse struct {
	AccessToken string `json:"token" example:"v2.local.Gdh5kiOTyyaQ3_bNykYDeYHO21Jg2..."`
}

// NewAuthResponse is a helper function to create a response body for handling authentication data
func NewAuthResponse(token string) AuthResponse {
	return AuthResponse{
		AccessToken: token,
	}
}

// UserResponse represents a user response body
type UserResponse struct {
	ID        uuid.UUID `json:"id" example:"bb073c91-f09b-4858-b2d1-d14116e73b8d"`
	Name      string    `json:"name" example:"John Doe"`
	Email     string    `json:"email" example:"test@example.com"`
	CreatedAt time.Time `json:"created_at" example:"1970-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"1970-01-01T00:00:00Z"`
}

// NewUserResponse is a helper function to create a response body for handling user data
func NewUserResponse(user *domain.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// RegistrationResponse represents a successful registration response body
type RegistrationResponse struct {
	Message    string `json:"message" example:"Registration successful. OTP code sent to your email."`
	OTPToken   string `json:"otp_token" example:"123456"`
	UserDetail UserResponse
}

// NewRegistrationResponse is a helper function to create a response body for successful registration
func NewRegistrationResponse(otpToken string, user *domain.User) RegistrationResponse {
	return RegistrationResponse{
		Message:    "Registration successful. OTP code sent to your email.",
		OTPToken:   otpToken,
		UserDetail: NewUserResponse(user),
	}
}

// CollectionResponse represents a collection response body
type CollectionResponse struct {
	ID          uuid.UUID `json:"id" example:"bb073c91-f09b-4858-b2d1-d14116e73b8d"`
	Name        string    `json:"name" example:"My Collection"`
	Description string    `json:"description,omitempty" example:"Collection description"`
	CreatedAt   time.Time `json:"created_at" example:"1970-01-01T00:00:00Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"1970-01-01T00:00:00Z"`
}

// NewCollectionResponse is a helper function to create a response body for handling collection data
func NewCollectionResponse(collection *domain.Collection) CollectionResponse {
	return CollectionResponse{
		ID:          collection.ID,
		Name:        collection.Name,
		Description: collection.Description,
		CreatedAt:   collection.CreatedAt,
		UpdatedAt:   collection.UpdatedAt,
	}
}

// SecretResponse represents a secret response body
type SecretResponse struct {
	ID           uuid.UUID             `json:"id" example:"bb073c91-f09b-4858-b2d1-d14116e73b8d"`
	CollectionID uuid.UUID             `json:"collection_id" example:"fab8dfe9-7cd0-4cd7-a387-7d6835a910d3"`
	SecretType   domain.SecretTypeEnum `json:"secret_type" example:"password"`
	CreatedAt    time.Time             `json:"created_at" example:"1970-01-01T00:00:00Z"`
	UpdatedAt    time.Time             `json:"updated_at" example:"1970-01-01T00:00:00Z"`
	CreatedBy    uuid.UUID             `json:"created_by" example:"f10ff052-b316-47f0-9788-ae8ebfa91b86"`
	UpdatedBy    uuid.UUID             `json:"updated_by" example:"f10ff052-b316-47f0-9788-ae8ebfa91b86"`
}

// NewSecretResponse is a helper function to create a response body for handling secret data
func NewSecretResponse(secret *domain.Secret) SecretResponse {
	return SecretResponse{
		ID:           secret.ID,
		CollectionID: secret.CollectionID,
		SecretType:   secret.SecretType,
		CreatedAt:    secret.CreatedAt,
		UpdatedAt:    secret.UpdatedAt,
		CreatedBy:    secret.CreatedBy,
		UpdatedBy:    secret.UpdatedBy,
	}
}
