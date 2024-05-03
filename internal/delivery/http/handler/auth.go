package handler

import (
	"net/http"

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

// registerRequet represents the request body for creating a user
type registerRequest struct {
	Name     string `json:"name" binding:"required" example:"John Doe"`
	Email    string `json:"email" binding:"required,email" example:"test@example.com"`
	Password string `json:"password" binding:"required,min=8" example:"12345678"`
}

// Register godoc
//
//	@Summary		Register a new user
//	@Description	create a new user account with default role "user"
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			registerRequest	body		registerRequest			true	"Register request"
//	@Success		200				{object}	response.UserResponse	"User created"
//	@Failure		400				{object}	response.ErrorResponse	"Validation error"
//	@Failure		401				{object}	response.ErrorResponse	"Unauthorized error"
//	@Failure		404				{object}	response.ErrorResponse	"Data not found error"
//	@Failure		409				{object}	response.ErrorResponse	"Data conflict error"
//	@Failure		500				{object}	response.ErrorResponse	"Internal server error"
//	@Router			/auth/register [post]
func (ah *AuthHandler) Register(ctx *gin.Context) {
	var req registerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidationError(ctx, err)
		return
	}

	user := domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	_, err := ah.svc.Register(ctx, &user)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	resp := response.NewResponse(true, "Registration successful. OTP code sent to your email.", nil)

	ctx.JSON(http.StatusOK, resp)
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
//	@Router			/auth/register/confirm [post]
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

// newRegistrationCodeRequest represents the request body for requesting a new OTP code
type newRegistrationCodeRequest struct {
	Email string `json:"email" binding:"required,email" example:"user@example.com"`
}

// RequestNewRegistrationCode godoc
//
//	@Summary		Request a new OTP code for registration confirmation
//	@Description	Requests a new OTP code for confirming user registration. If the previous OTP code
//					has expired or the user hasn't requested one before, a new OTP code will be generated
//					and sent to the user's email.
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			request	body		newRegistrationCodeRequest	true	"Request new OTP request body"
//	@Success		200		{object}	response.Response			"OTP code requested successfully"
//	@Failure		400		{object}	response.ErrorResponse		"Validation error"
//	@Failure		429		{object}	response.ErrorResponse		"Too many requests, try again later"
//	@Failure		500		{object}	response.ErrorResponse		"Internal server error"
//	@Router			/auth/register/request-new-code [post]
func (ah *AuthHandler) RequestNewRegistrationCode(ctx *gin.Context) {
	var req newRegistrationCodeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidationError(ctx, err)
		return
	}

	// Call RequestNewRegistrationCode method
	err := ah.svc.RequestNewRegistrationCode(ctx.Request.Context(), req.Email)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.HandleSuccess(ctx, nil)
}

// Logout godoc
//
//	@Summary		Logout a user
//	@Description	Logs out a user by invalidating the access token
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Security		Bearer
//	@Success		200		string	"Successfully logged out"
//	@Failure		401		{object}	response.ErrorResponse	"Unauthorized error"
//	@Router			/auth/logout [post]
func (ah *AuthHandler) Logout(ctx *gin.Context) {
	// Retrieve the user ID from the context (assuming it's stored during authentication)
	authPayload := helper.GetAuthPayload(ctx, middleware.AuthorizationPayloadKey)

	err := ah.svc.Logout(ctx, authPayload)
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

// ResetPassword godoc
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

// forgotPasswordRequest represents the request body for requesting a reset of a forgotten password
type forgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email" example:"user@example.com"`
}

// RequestResetPassword godoc
//
//	@Summary		Request to reset forgotten password
//	@Description	Initiate the process of resetting a forgotten password by providing the user's email
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			request	body		forgotPasswordRequest	true	"Request reset forgot password request body"
//	@Success		200		{object}	response.Response		"Password reset request initiated successfully"
//	@Failure		400		{object}	response.ErrorResponse	"Validation error"
//	@Failure		404		{object}	response.ErrorResponse	"User not found error"
//	@Failure		500		{object}	response.ErrorResponse	"Internal server error"
//	@Router			/auth/forgot-password [post]
func (ah *AuthHandler) ForgotPassword(ctx *gin.Context) {
	var req forgotPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidationError(ctx, err)
		return
	}

	err := ah.svc.ForgotPassword(ctx, req.Email)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.HandleSuccess(ctx, nil)
}

// resetPasswordRequest represents the request body for resetting password
type resetPasswordRequest struct {
	Email       string `json:"email" binding:"required,email" example:"test@example.com"`
	NewPassword string `json:"new_password" binding:"required,min=8" example:"new_password" minLength:"8"`
	OTP         string `json:"otp" binding:"required" example:"123456"`
}

// ResetPassword godoc
//
//	@Summary		Reset user's password
//	@Description	Resets user's password after confirmation with OTP code.
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//
//	@Param			request	body		resetPasswordRequest	true	"Request reset password request body"
//
//	@Success		200		{string}	string					"Password reset successfully"
//	@Failure		400		{string}	string					"Invalid email or password format"
//	@Failure		401		{string}	string					"Invalid OTP code"
//	@Failure		404		{string}	string					"User not found"
//	@Failure		500		{string}	string					"Internal server error"
//	@Router			/auth/reset-password [post]
func (ah *AuthHandler) ResetPassword(ctx *gin.Context) {
	var req resetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidationError(ctx, err)
		return
	}

	err := ah.svc.ResetPassword(ctx, req.Email, req.NewPassword, req.OTP)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.HandleSuccess(ctx, nil)
}
