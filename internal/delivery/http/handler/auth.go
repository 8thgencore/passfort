package handler

import (
	"github.com/8thgencore/passfort/internal/delivery/http/helper"
	"github.com/8thgencore/passfort/internal/delivery/http/middleware"
	"github.com/8thgencore/passfort/internal/delivery/http/response"
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

// confirmRegistrationRequest represents the request body for confirming registration with OTP code
type confirmRegistrationRequest struct {
	Email string `json:"email" binding:"required,email" example:"test@example.com"`
	OTP   string `json:"otp" binding:"required" example:"123456"`
}

// ConfirmRegistration godoc
//
//	@Summary		Confirm user registration with OTP code
//	@Description	Confirm user registration by providing the email and OTP code
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			request	body		confirmRegistrationRequest	true	"Confirm registration request"
//	@Success		200		{object}	response.Response			"Successfully confirmed registration"
//	@Failure		400		{object}	response.ErrorResponse		"Validation error"
//	@Failure		404		{object}	response.ErrorResponse		"Data not found error"
//	@Failure		500		{object}	response.ErrorResponse		"Internal server error"
//	@Router			/auth/confirm-registration [post]
func (ah *AuthHandler) ConfirmRegistration(ctx *gin.Context) {
	var req confirmRegistrationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidationError(ctx, err)
		return
	}

	err := ah.svc.ConfirmRegistration(ctx, req.Email, req.OTP)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.HandleSuccess(ctx, nil)
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

// resetPasswordRequest represents the request body for requesting a reset of a forgotten password
type resetPasswordRequest struct {
	Email string `json:"email" binding:"required,email" example:"user@example.com"`
}

// RequestResetPassword godoc
//
//	@Summary		Request to reset forgotten password
//	@Description	Initiate the process of resetting a forgotten password by providing the user's email
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			request	body		resetPasswordRequest	true	"Request reset forgot password request body"
//	@Success		200		{object}	response.Response		"Password reset request initiated successfully"
//	@Failure		400		{object}	response.ErrorResponse	"Validation error"
//	@Failure		404		{object}	response.ErrorResponse	"User not found error"
//	@Failure		500		{object}	response.ErrorResponse	"Internal server error"
//	@Router			/auth/request-reset-password [post]
func (ah *AuthHandler) RequestResetPassword(ctx *gin.Context) {
	var req resetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidationError(ctx, err)
		return
	}

	err := ah.svc.RequestResetPassword(ctx, req.Email)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.HandleSuccess(ctx, nil)
}

// confirmResetPasswordRequest represents the request body for confirming password reset with OTP code
type confirmResetPasswordRequest struct {
	Email string `json:"email" binding:"required,email" example:"test@example.com"`
	OTP   string `json:"otp" binding:"required" example:"123456"`
}

// ConfirmResetPassword godoc
//
//	@Summary		Confirm password reset with OTP code
//	@Description	Confirm password reset by providing the email and OTP code
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			request	body		confirmResetPasswordRequest	true	"Confirm password reset request"
//	@Success		200		{object}	response.Response			"Successfully confirmed password reset"
//	@Failure		400		{object}	response.ErrorResponse		"Validation error"
//	@Failure		404		{object}	response.ErrorResponse		"Data not found error"
//	@Failure		500		{object}	response.ErrorResponse		"Internal server error"
//	@Router			/auth/confirm-reset-password [post]
func (ah *AuthHandler) ConfirmResetPassword(ctx *gin.Context) {
	var req confirmResetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidationError(ctx, err)
		return
	}

	err := ah.svc.ConfirmResetPassword(ctx, req.Email, req.OTP)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.HandleSuccess(ctx, nil)
}

// setNewPasswordRequest represents the request body for resetting password
type setNewPasswordRequest struct {
	Email       string `json:"email" binding:"required,email" example:"test@example.com"`
	NewPassword string `json:"new_password" binding:"required,min=8" example:"newpassword" minLength:"8"`
	OTP         string `json:"otp" binding:"required" example:"123456"`
}

// SetNewPassword godoc
//
//	@Summary		Reset user password after confirmation with OTP code
//	@Description	Reset user password by providing the email, new password, and OTP code
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			request	body		setNewPasswordRequest	true	"Reset password request"
//	@Success		200		{object}	response.Response		"Successfully reset password"
//	@Failure		400		{object}	response.ErrorResponse	"Validation error"
//	@Failure		404		{object}	response.ErrorResponse	"Data not found error"
//	@Failure		422		{object}	response.ErrorResponse	"Passwords do not match"
//	@Failure		500		{object}	response.ErrorResponse	"Internal server error"
//	@Router			/auth/set-new-password [put]
func (ah *AuthHandler) SetNewPassword(ctx *gin.Context) {
	var req setNewPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidationError(ctx, err)
		return
	}

	err := ah.svc.SetNewPassword(ctx, req.Email, req.NewPassword, req.OTP)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.HandleSuccess(ctx, nil)
}
