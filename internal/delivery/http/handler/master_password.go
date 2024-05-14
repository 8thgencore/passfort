package handler

import (
	"github.com/8thgencore/passfort/internal/delivery/http/helper"
	"github.com/8thgencore/passfort/internal/delivery/http/middleware"
	"github.com/8thgencore/passfort/internal/delivery/http/response"
	"github.com/8thgencore/passfort/internal/domain"
	"github.com/8thgencore/passfort/internal/service"

	"github.com/gin-gonic/gin"
)

// MasterPasswordHandler represents the HTTP handler for master password-related requests
type MasterPasswordHandler struct {
	svc service.MasterPasswordService
}

// NewMasterPasswordHandler creates a new MasterPasswordHandler instance
func NewMasterPasswordHandler(svc service.MasterPasswordService) *MasterPasswordHandler {
	return &MasterPasswordHandler{
		svc,
	}
}

// createMasterPasswordRequest represents the request body for creating a master password
type createMasterPasswordRequest struct {
	Password string `json:"password" binding:"required,min=8" example:"masterpassword"`
}

// CreateMasterPassword godoc
//
//	@Summary		Create master password
//	@Description	Create a master password for the authenticated user
//	@Tags			MasterPassword
//	@Accept			json
//	@Produce		json
//	@Param			request	body		createMasterPasswordRequest	true	"Create master password request"
//	@Success		200		{object}	response.Response			"Master password created successfully"
//	@Failure		400		{object}	response.ErrorResponse		"Validation error"
//	@Failure		409		{object}	response.ErrorResponse		"Master password already exists"
//	@Failure		500		{object}	response.ErrorResponse		"Internal server error"
//	@Router			/master-password [post]
//	@Security		BearerAuth
func (h *MasterPasswordHandler) CreateMasterPassword(ctx *gin.Context) {
	var req createMasterPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidationError(ctx, err)
		return
	}

	authPayload := helper.GetAuthPayload(ctx, middleware.AuthorizationPayloadKey)
	userID := authPayload.UserID

	// Check if master password already exists
	exists, err := h.svc.MasterPasswordExists(ctx, userID)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	if exists {
		response.HandleError(ctx, domain.ErrMasterPasswordAlreadyExists)
		return
	}

	// Save the new master password (hashed)
	if err := h.svc.SaveMasterPassword(ctx, userID, req.Password); err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.HandleSuccess(ctx, nil)
}

// changeMasterPasswordRequest represents the request body for changing a master password
type changeMasterPasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required,min=8" example:"currentmasterpassword"`
	NewPassword     string `json:"new_password" binding:"required,min=8" example:"newmasterpassword"`
}

// ChangeMasterPassword godoc
//
//	@Summary		Change master password
//	@Description	Change the master password for the authenticated user
//	@Tags			MasterPassword
//	@Accept			json
//	@Produce		json
//	@Param			request	body		changeMasterPasswordRequest	true	"Change master password request"
//	@Success		200		{object}	response.Response			"Master password changed successfully"
//	@Failure		400		{object}	response.ErrorResponse		"Validation error"
//	@Failure		401		{object}	response.ErrorResponse		"Unauthorized error"
//	@Failure		500		{object}	response.ErrorResponse		"Internal server error"
//	@Router			/master-password [put]
//	@Security		BearerAuth
func (h *MasterPasswordHandler) ChangeMasterPassword(ctx *gin.Context) {
	var req changeMasterPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidationError(ctx, err)
		return
	}

	authPayload := helper.GetAuthPayload(ctx, middleware.AuthorizationPayloadKey)
	userID := authPayload.UserID

	// Save the new master password (hashed)
	if err := h.svc.SaveMasterPassword(ctx, userID, req.NewPassword); err != nil {
		response.HandleError(ctx, err)
		return
	}

	// Validate current master password
	if err := h.svc.ValidateMasterPassword(ctx, userID, req.CurrentPassword); err != nil {
		response.HandleError(ctx, domain.ErrInvalidMasterPassword)
		return
	}

	response.HandleSuccess(ctx, nil)
}

// validateMasterPasswordRequest represents the request body for validating a master password
type validateMasterPasswordRequest struct {
	Password string `json:"password" binding:"required,min=8" example:"masterpassword"`
}

// ValidateMasterPassword godoc
//
//	@Summary		Validate master password
//	@Description	Validate the master password for the authenticated user
//	@Tags			MasterPassword
//	@Accept			json
//	@Produce		json
//	@Param			request	body		validateMasterPasswordRequest	true	"Validate master password request"
//	@Success		200		{object}	response.Response				"Master password is valid"
//	@Failure		400		{object}	response.ErrorResponse			"Validation error"
//	@Failure		401		{object}	response.ErrorResponse			"Invalid master password"
//	@Failure		500		{object}	response.ErrorResponse			"Internal server error"
//	@Router			/master-password/validate [post]
//	@Security		BearerAuth
func (h *MasterPasswordHandler) ValidateMasterPassword(ctx *gin.Context) {
	var req validateMasterPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidationError(ctx, err)
		return
	}

	authPayload := helper.GetAuthPayload(ctx, middleware.AuthorizationPayloadKey)
	userID := authPayload.UserID

	// Validate current master password
	if err := h.svc.ValidateMasterPassword(ctx, userID, req.Password); err != nil {
		response.HandleError(ctx, domain.ErrInvalidMasterPassword)
		return
	}

	response.HandleSuccess(ctx, nil)
}
