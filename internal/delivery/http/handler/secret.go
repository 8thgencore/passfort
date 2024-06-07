package handler

import (
	"errors"

	"github.com/8thgencore/passfort/internal/delivery/http/helper"
	"github.com/8thgencore/passfort/internal/delivery/http/middleware"
	"github.com/8thgencore/passfort/internal/delivery/http/response"
	"github.com/8thgencore/passfort/internal/domain"
	"github.com/8thgencore/passfort/internal/service"
	"github.com/8thgencore/passfort/pkg/base64_util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SecretHandler represents the HTTP handler for secret-related requests
type SecretHandler struct {
	svc service.SecretService
}

// NewSecretHandler creates a new SecretHandler instance
func NewSecretHandler(svc service.SecretService) *SecretHandler {
	return &SecretHandler{
		svc,
	}
}

// createSecretRequest represents the request body for creating a secret
type createSecretRequest struct {
	Name        string `json:"name" binding:"required" example:"My Secret"`
	Description string `json:"description" binding:"required" example:"This is a secret"`
	URL         string `json:"url,omitempty" example:"https://example.com"`       // Optional for PasswordSecret
	Login       string `json:"login,omitempty" example:"user@example.com"`        // Optional for PasswordSecret
	Password    string `json:"password,omitempty" example:"password123"`          // Optional for PasswordSecret
	Text        string `json:"text,omitempty" example:"This is some secret text"` // Optional for TextSecret
	SecretType  string `json:"secret_type" binding:"required" example:"password"` // "password" or "text"
}

// CreateSecret godoc
//
//	@Summary		Create a new secret
//	@Description	Create a new secret (password or text)
//	@Tags			Secrets
//	@Accept			json
//	@Produce		json
//	@Param			collection_id	path		string					true	"Collection ID"
//	@Param			request			body		createSecretRequest		true	"Create Secret Request"
//	@Success		201				{object}	response.SecretResponse	"Secret created"
//	@Failure		400				{object}	response.ErrorResponse	"Validation error"
//	@Failure		500				{object}	response.ErrorResponse	"Internal server error"
//	@Router			/collections/{collection_id}/secrets [post]
//	@Security		BearerAuth
func (sh *SecretHandler) CreateSecret(ctx *gin.Context) {
	var req createSecretRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidationError(ctx, err)
		return
	}

	collectionID, err := uuid.Parse(ctx.Param("collection_id"))
	if err != nil {
		response.ValidationError(ctx, err)
		return
	}

	authPayload := helper.GetAuthPayload(ctx, middleware.AuthorizationPayloadKey)

	var newSecret *domain.Secret
	switch req.SecretType {
	case string(domain.PasswordSecretType):
		if req.URL == "" || req.Login == "" || req.Password == "" {
			response.ValidationError(ctx, errors.New("missing fields for password secret"))
			return
		}
		newSecret = &domain.Secret{
			CollectionID: collectionID,
			SecretType:   domain.PasswordSecretType,
			Name:         req.Name,
			Description:  req.Description,
			CreatedBy:    authPayload.UserID,
			UpdatedBy:    authPayload.UserID,
			PasswordSecret: &domain.PasswordSecret{
				URL:      req.URL,
				Login:    req.Login,
				Password: req.Password,
			},
		}
	case string(domain.TextSecretType):
		if req.Text == "" {
			response.ValidationError(ctx, errors.New("missing fields for text secret"))
			return
		}
		newSecret = &domain.Secret{
			CollectionID: collectionID,
			SecretType:   domain.TextSecretType,
			Name:         req.Name,
			Description:  req.Description,
			CreatedBy:    authPayload.UserID,
			UpdatedBy:    authPayload.UserID,
			TextSecret: &domain.TextSecret{
				Text: req.Text,
			},
		}
	default:
		response.ValidationError(ctx, domain.ErrInvalidSecretType)
		return
	}

	encryptionKey, err := base64_util.Base64ToBytes(helper.GetEncryptionKey(ctx, middleware.EncryptionKey))
	if err != nil {
		response.HandleError(ctx, domain.ErrInternal)
		return
	}

	createdSecret, err := sh.svc.CreateSecret(ctx, authPayload.UserID, newSecret, encryptionKey)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	rsp := response.NewSecretResponse(createdSecret)

	response.HandleSuccess(ctx, rsp)
}

// listMeSecretsRequest represents the request body for listing secrets by user ID
type listMeSecretsRequest struct {
	Skip  uint64 `form:"skip" binding:"required,min=0" example:"0"`
	Limit uint64 `form:"limit" binding:"required,min=5" example:"5"`
}

// ListMeSecrets godoc
//
//	@Summary		List me secrets
//	@Description	List me secrets associated with pagination
//	@Tags			Secrets
//	@Accept			json
//	@Produce		json
//	@Param			collection_id	path		string					true	"Collection ID"
//	@Param			skip			query		uint64					true	"Skip"
//	@Param			limit			query		uint64					true	"Limit"
//	@Success		200				{object}	response.Meta			"Secrets displayed"
//	@Failure		400				{object}	response.ErrorResponse	"Validation error"
//	@Failure		500				{object}	response.ErrorResponse	"Internal server error"
//	@Router			/collections/{collection_id}/secrets [get]
//	@Security		BearerAuth
func (sh *SecretHandler) ListMeSecrets(ctx *gin.Context) {
	var req listMeSecretsRequest
	var secretsList []response.SecretResponse

	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.ValidationError(ctx, err)
		return
	}

	collectionID, err := uuid.Parse(ctx.Param("collection_id"))
	if err != nil {
		response.ValidationError(ctx, err)
		return
	}

	authPayload := helper.GetAuthPayload(ctx, middleware.AuthorizationPayloadKey)

	secrets, err := sh.svc.ListSecretsByCollectionID(ctx, authPayload.UserID, collectionID, req.Skip, req.Limit)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	for _, secret := range secrets {
		secretsList = append(secretsList, response.NewSecretResponse(&secret))
	}

	total := uint64(len(secretsList))
	meta := response.NewMeta(total, req.Limit, req.Skip)
	rsp := helper.ToMap(meta, secretsList, "secrets")

	response.HandleSuccess(ctx, rsp)
}

// getSecretRequest represents the request body for getting a secret
type getSecretRequest struct {
	CollectionID string `uri:"collection_id" binding:"required"`
	SecretID     string `uri:"secret_id" binding:"required"`
}

// GetSecret godoc
//
//	@Summary		Get a secret
//	@Description	Get a secret by id
//	@Tags			Secrets
//	@Accept			json
//	@Produce		json
//	@Param			collection_id	path		string					true	"Collection ID"
//	@Param			secret_id		path		string					true	"Secret ID"
//	@Success		200				{object}	response.SecretResponse	"Secret displayed"
//	@Failure		400				{object}	response.ErrorResponse	"Validation error"
//	@Failure		404				{object}	response.ErrorResponse	"Data not found error"
//	@Failure		500				{object}	response.ErrorResponse	"Internal server error"
//	@Router			/collections/{collection_id}/secrets/{secret_id} [get]
//	@Security		BearerAuth
func (sh *SecretHandler) GetSecret(ctx *gin.Context) {
	var req getSecretRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		response.ValidationError(ctx, err)
		return
	}

	collectionID, err := uuid.Parse(req.CollectionID)
	if err != nil {
		response.ValidationError(ctx, err)
		return
	}

	secretID, err := uuid.Parse(req.SecretID)
	if err != nil {
		response.ValidationError(ctx, err)
		return
	}

	authPayload := helper.GetAuthPayload(ctx, middleware.AuthorizationPayloadKey)

	encryptionKey, err := base64_util.Base64ToBytes(helper.GetEncryptionKey(ctx, middleware.EncryptionKey))
	if err != nil {
		response.HandleError(ctx, domain.ErrInternal)
		return
	}

	secret, err := sh.svc.GetSecret(ctx, authPayload.UserID, collectionID, secretID, encryptionKey)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	rsp := response.NewSecretResponse(secret)

	response.HandleSuccess(ctx, rsp)
}

type updateSecretRequest struct {
	Name        string `json:"name" binding:"required" example:"My Secret"`
	Description string `json:"description" binding:"required" example:"This is a secret"`
	URL         string `json:"url,omitempty" example:"https://example.com"`       // Optional for PasswordSecret
	Login       string `json:"login,omitempty" example:"user@example.com"`        // Optional for PasswordSecret
	Password    string `json:"password,omitempty" example:"password123"`          // Optional for PasswordSecret
	Text        string `json:"text,omitempty" example:"This is some secret text"` // Optional for TextSecret
	SecretType  string `json:"secret_type" binding:"required" example:"password"` // "password" or "text"
}

// UpdateSecret godoc
//
//	@Summary		Update a secret
//	@Description	Update a secret (password or text) by id
//	@Tags			Secrets
//	@Accept			json
//	@Produce		json
//	@Param			collection_id	path		string					true	"Collection ID"
//	@Param			secret_id		path		string					true	"Secret ID"
//	@Param			request			body		updateSecretRequest		true	"Update Secret Request"
//	@Success		200				{object}	response.SecretResponse	"Secret updated"
//	@Failure		400				{object}	response.ErrorResponse	"Validation error"
//	@Failure		404				{object}	response.ErrorResponse	"Data not found error"
//	@Failure		500				{object}	response.ErrorResponse	"Internal server error"
//	@Router			/collections/{collection_id}/secrets/{secret_id} [put]
//	@Security		BearerAuth
func (sh *SecretHandler) UpdateSecret(ctx *gin.Context) {
	var req updateSecretRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidationError(ctx, err)
		return
	}

	collectionID, err := uuid.Parse(ctx.Param("collection_id"))
	if err != nil {
		response.ValidationError(ctx, err)
		return
	}

	secretID, err := uuid.Parse(ctx.Param("secret_id"))
	if err != nil {
		response.ValidationError(ctx, err)
		return
	}

	authPayload := helper.GetAuthPayload(ctx, middleware.AuthorizationPayloadKey)

	var secret *domain.Secret
	switch req.SecretType {
	case string(domain.PasswordSecretType):
		if req.URL == "" || req.Login == "" || req.Password == "" {
			response.ValidationError(ctx, errors.New("missing fields for password secret"))
			return
		}
		secret = &domain.Secret{
			ID:           secretID,
			CollectionID: collectionID,
			Name:         req.Name,
			Description:  req.Description,
			UpdatedBy:    authPayload.UserID,
			SecretType:   domain.PasswordSecretType,
			PasswordSecret: &domain.PasswordSecret{
				URL:      req.URL,
				Login:    req.Login,
				Password: req.Password,
			},
		}
	case string(domain.TextSecretType):
		if req.Text == "" {
			response.ValidationError(ctx, errors.New("missing fields for text secret"))
			return
		}
		secret = &domain.Secret{
			ID:           secretID,
			CollectionID: collectionID,
			Name:         req.Name,
			Description:  req.Description,
			UpdatedBy:    authPayload.UserID,
			SecretType:   domain.TextSecretType,
			TextSecret: &domain.TextSecret{
				Text: req.Text,
			},
		}
	default:
		response.ValidationError(ctx, domain.ErrInvalidSecretType)
		return
	}

	encryptionKey, err := base64_util.Base64ToBytes(helper.GetEncryptionKey(ctx, middleware.EncryptionKey))
	if err != nil {
		response.HandleError(ctx, domain.ErrInternal)
		return
	}

	updatedSecret, err := sh.svc.UpdateSecret(ctx, authPayload.UserID, collectionID, secret, encryptionKey)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	rsp := response.NewSecretResponse(updatedSecret)
	response.HandleSuccess(ctx, rsp)
}

// deleteSecretRequest represents the request body for deleting a secret
type deleteSecretRequest struct {
	CollectionID string `uri:"collection_id" binding:"required"`
	SecretID     string `uri:"secret_id" binding:"required"`
}

// DeleteSecret godoc
//
//	@Summary		Delete a secret
//	@Description	Delete a secret by id
//	@Tags			Secrets
//	@Accept			json
//	@Produce		json
//	@Param			collection_id	path		string					true	"Collection ID"
//	@Param			secret_id		path		string					true	"Secret ID"
//	@Success		200				{object}	response.Response		"Secret deleted"
//	@Failure		400				{object}	response.ErrorResponse	"Validation error"
//	@Failure		401				{object}	response.ErrorResponse	"Unauthorized error"
//	@Failure		403				{object}	response.ErrorResponse	"Forbidden error"
//	@Failure		404				{object}	response.ErrorResponse	"Data not found error"
//	@Failure		500				{object}	response.ErrorResponse	"Internal server error"
//	@Router			/collections/{collection_id}/secrets/{secret_id} [delete]
//	@Security		BearerAuth
func (sh *SecretHandler) DeleteSecret(ctx *gin.Context) {
	var req deleteSecretRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		response.ValidationError(ctx, err)
		return
	}

	collectionID, err := uuid.Parse(req.CollectionID)
	if err != nil {
		response.ValidationError(ctx, err)
		return
	}

	secretID, err := uuid.Parse(req.SecretID)
	if err != nil {
		response.ValidationError(ctx, err)
		return
	}

	authPayload := helper.GetAuthPayload(ctx, middleware.AuthorizationPayloadKey)

	err = sh.svc.DeleteSecret(ctx, authPayload.UserID, collectionID, secretID)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.HandleSuccess(ctx, nil)
}
