package handler

import (
	"github.com/8thgencore/passfort/internal/delivery/http/helper"
	"github.com/8thgencore/passfort/internal/delivery/http/middleware"
	"github.com/8thgencore/passfort/internal/delivery/http/response"
	"github.com/8thgencore/passfort/internal/domain"
	"github.com/8thgencore/passfort/internal/service"
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
	SecretType domain.SecretTypeEnum `json:"secret_type" binding:"required,secret_type" example:"password"`
}

// CreateSecret godoc
//
//	@Summary		Create a new secret
//	@Description	Create a new secret
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

	newSecret := domain.Secret{
		CollectionID: collectionID,
		SecretType:   req.SecretType,
	}

	authPayload := helper.GetAuthPayload(ctx, middleware.AuthorizationPayloadKey)

	createdSecret, err := sh.svc.CreateSecret(ctx, authPayload.UserID, &newSecret)
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

	secret, err := sh.svc.GetSecret(ctx, authPayload.UserID, collectionID, secretID)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	rsp := response.NewSecretResponse(secret)

	response.HandleSuccess(ctx, rsp)
}

// TODO:
// updateSecretRequest represents the request body for updating a secret
// type updateSecretRequest struct {
// 	// Define fields to be updated
// }
// UpdateSecret godoc
//
//	@Summary		Update a secret
//	@Description	Update a secret's fields by id
//	@Tags			Secrets
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string					true	"Secret ID"
//	@Param			request	body		updateSecretRequest		true	"Update secret request"
//	@Success		200		{object}	response.SecretResponse	"Secret updated"
//	@Failure		400		{object}	response.ErrorResponse	"Validation error"
//	@Failure		401		{object}	response.ErrorResponse	"Unauthorized error"
//	@Failure		403		{object}	response.ErrorResponse	"Forbidden error"
//	@Failure		404		{object}	response.ErrorResponse	"Data not found error"
//	@Failure		500		{object}	response.ErrorResponse	"Internal server error"
//	@Router			/collections/{collection_id}/secrets/{id} [put]
//	@Security		BearerAuth
// func (sh *SecretHandler) UpdateSecret(ctx *gin.Context) {
// 	var req updateSecretRequest
// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		response.ValidationError(ctx, err)
// 		return
// 	}

// 	secretID, err := uuid.Parse(ctx.Param("id"))
// 	if err != nil {
// 		response.ValidationError(ctx, err)
// 		return
// 	}

// 	authPayload := helper.GetAuthPayload(ctx, middleware.AuthorizationPayloadKey)

// 	sh.svc.UpdateSecret(ctx, authPayload.UserID)
// 	// Update the secret using the SecretService's UpdateSecret method
// 	// ...

// 	response.HandleSuccess(ctx, nil)
// }

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
