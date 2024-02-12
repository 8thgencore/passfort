package handler

import (
	"github.com/8thgencore/passfort/internal/delivery/http/helper"
	"github.com/8thgencore/passfort/internal/delivery/http/middleware"
	"github.com/8thgencore/passfort/internal/delivery/http/response"
	"github.com/8thgencore/passfort/internal/domain"
	"github.com/8thgencore/passfort/internal/service"
	"github.com/gin-gonic/gin"
)

// AuthHandler represents the HTTP handler for authentication-related requests
type AuthHandler struct {
	svc service.AuthService
}

// NewAuthHandler creates a new AuthHandler instance
func NewAuthHandler(svc service.AuthService) *AuthHandler {
	return &AuthHandler{
		svc,
	}
}


// loginRequest represents the request body for logging in a user
type loginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"test@example.com"`
	Password string `json:"password" binding:"required,min=8" example:"12345678" minLength:"8"`
}

// Login godoc
//
//	@Summary		Login and get an access token
//	@Description	Logs in a registered user and returns an access token if the credentials are valid.
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			request	body		loginRequest			true	"Login request body"
//	@Success		200		{object}	response.AuthResponse	"Succesfully logged in"
//	@Failure		400		{object}	response.ErrorResponse	"Validation error"
//	@Failure		401		{object}	response.ErrorResponse	"Unauthorized error"
//	@Failure		500		{object}	response.ErrorResponse	"Internal server error"
//	@Router			/auth/login [post]
func (ah *AuthHandler) Login(ctx *gin.Context) {
	var req loginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidationError(ctx, err)
		return
	}

	token, err := ah.svc.Login(ctx, req.Email, req.Password)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	rsp := response.NewAuthResponse(token)

	response.HandleSuccess(ctx, rsp)
}

type changePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=8" example:"oldpassword"`
	NewPassword string `json:"new_password" binding:"required,min=8" example:"newpassword"`
}

// ChangeOwnPassword godoc
//
//	@Summary		Change own password
//	@Description	Change the authenticated user's password by providing the old and new passwords
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			changePasswordRequest	body		changePasswordRequest	true	"Change password request"
//	@Success		200						{object}	response.Response		"Password changed successfully"
//	@Failure		400						{object}	response.ErrorResponse	"Validation error"
//	@Failure		401						{object}	response.ErrorResponse	"Unauthorized error"
//	@Failure		403						{object}	response.ErrorResponse	"Forbidden error"
//	@Failure		404						{object}	response.ErrorResponse	"Data not found error"
//	@Failure		422						{object}	response.ErrorResponse	"Passwords do not match"
//	@Failure		500						{object}	response.ErrorResponse	"Internal server error"
//	@Router			/auth/change-password [put]
//	@Security		BearerAuth
func (ah *AuthHandler) ChangePassword(ctx *gin.Context) {
	var req changePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidationError(ctx, err)
		return
	}

	authPayload := helper.GetAuthPayload(ctx, middleware.AuthorizationPayloadKey)

	err := ah.svc.ChangePassword(ctx, authPayload.UserID, req.OldPassword, req.NewPassword)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.HandleSuccess(ctx, nil)
}
