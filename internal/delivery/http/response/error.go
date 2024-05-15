package response

import (
	"errors"
	"net/http"

	"github.com/8thgencore/passfort/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// errorStatusMap is a map of defined error messages and their corresponding http status codes
var errorStatusMap = map[error]int{
	// Internal Errors
	domain.ErrInternal:          http.StatusInternalServerError,
	domain.ErrFailedToSendEmail: http.StatusInternalServerError,

	// Data Errors
	domain.ErrDataNotFound:    http.StatusNotFound,
	domain.ErrNoUpdatedData:   http.StatusBadRequest,
	domain.ErrConflictingData: http.StatusConflict,

	// Token Errors
	domain.ErrTokenDuration:       http.StatusBadRequest,
	domain.ErrTokenCreation:       http.StatusInternalServerError,
	domain.ErrExpiredToken:        http.StatusUnauthorized,
	domain.ErrInvalidToken:        http.StatusUnauthorized,
	domain.ErrInvalidRefreshToken: http.StatusUnauthorized,

	// Authentication Errors
	domain.ErrInvalidCredentials:  http.StatusUnauthorized,
	domain.ErrPasswordsDoNotMatch: http.StatusBadRequest,
	domain.ErrInvalidOTP:          http.StatusUnauthorized,
	domain.ErrOTPAlreadySent:      http.StatusTooManyRequests,

	// Authorization Errors
	domain.ErrEmptyAuthorizationHeader:   http.StatusUnauthorized,
	domain.ErrInvalidAuthorizationHeader: http.StatusUnauthorized,
	domain.ErrInvalidAuthorizationType:   http.StatusUnauthorized,
	domain.ErrUnauthorized:               http.StatusUnauthorized,
	domain.ErrForbidden:                  http.StatusForbidden,

	// User Errors
	domain.ErrUserNotVerified:  http.StatusUnauthorized,
	domain.ErrDeleteOwnAccount: http.StatusForbidden,

	// Master Password Errors
	domain.ErrMasterPasswordActivationExpired: http.StatusUnauthorized,
	domain.ErrMasterPasswordNotSet:            http.StatusUnauthorized,
	domain.ErrInvalidMasterPassword:           http.StatusUnauthorized,
	domain.ErrMasterPasswordAlreadyExists:     http.StatusConflict,
}

// ValidationError sends an error response for some specific request validation error
func ValidationError(ctx *gin.Context, err error) {
	errMsgs := ParseError(err)
	errRsp := NewErrorResponse(errMsgs)
	ctx.JSON(http.StatusBadRequest, errRsp)
}

// HandleError determines the status code of an error and returns a JSON response with the error message and status code
func HandleError(ctx *gin.Context, err error) {
	statusCode, ok := errorStatusMap[err]
	if !ok {
		statusCode = http.StatusInternalServerError
	}

	errMsg := ParseError(err)
	errRsp := NewErrorResponse(errMsg)
	ctx.JSON(statusCode, errRsp)
}

// HandleAbort sends an error response and aborts the request with the specified status code and error message
func HandleAbort(ctx *gin.Context, err error) {
	statusCode, ok := errorStatusMap[err]
	if !ok {
		statusCode = http.StatusInternalServerError
	}

	errMsg := ParseError(err)
	errRsp := NewErrorResponse(errMsg)
	ctx.AbortWithStatusJSON(statusCode, errRsp)
}

// ParseError parses error messages from the error object and returns a slice of error messages
func ParseError(err error) []string {
	var errMsgs []string

	if errors.As(err, &validator.ValidationErrors{}) {
		for _, err := range err.(validator.ValidationErrors) {
			errMsgs = append(errMsgs, err.Error())
		}
	} else {
		errMsgs = append(errMsgs, err.Error())
	}

	return errMsgs
}

// ErrorResponse represents an error response body format
type ErrorResponse struct {
	Success  bool     `json:"success" example:"false"`
	Messages []string `json:"messages" example:"Error message 1, Error message 2"`
}

// NewErrorResponse is a helper function to create an error response body
func NewErrorResponse(errMsgs []string) ErrorResponse {
	return ErrorResponse{
		Success:  false,
		Messages: errMsgs,
	}
}

// HandleSuccess sends a success response with the specified status code and optional data
func HandleSuccess(ctx *gin.Context, data any) {
	rsp := NewResponse(true, "Success", data)
	ctx.JSON(http.StatusOK, rsp)
}
