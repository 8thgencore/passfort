package handler

import (
	"github.com/8thgencore/passfort/internal/delivery/http/helper"
	"github.com/8thgencore/passfort/internal/delivery/http/response"
	"github.com/8thgencore/passfort/internal/domain"
	"github.com/8thgencore/passfort/internal/service"
	"github.com/gin-gonic/gin"
)

// UserHandler represents the HTTP handler for user-related requests
type UserHandler struct {
	svc service.UserService
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{
		svc,
	}
}

// listUsersRequest represents the request body for listing users
type listUsersRequest struct {
	Skip  uint64 `form:"skip" binding:"required,min=0" example:"0"`
	Limit uint64 `form:"limit" binding:"required,min=5" example:"5"`
}

// ListUsers godoc
//
//	@Summary		List users
//	@Description	List users with pagination
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			skip	query		uint64					true	"Skip"
//	@Param			limit	query		uint64					true	"Limit"
//	@Success		200		{object}	response.Meta			"Users displayed"
//	@Failure		400		{object}	response.ErrorResponse	"Validation error"
//	@Failure		500		{object}	response.ErrorResponse	"Internal server error"
//	@Router			/users [get]
//	@Security		BearerAuth
func (uh *UserHandler) ListUsers(ctx *gin.Context) {
	var req listUsersRequest
	var usersList []response.UserResponse

	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.ValidationError(ctx, err)
		return
	}

	users, err := uh.svc.ListUsers(ctx, req.Skip, req.Limit)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	for _, user := range users {
		usersList = append(usersList, response.NewUserResponse(&user))
	}

	total := uint64(len(usersList))
	meta := response.NewMeta(total, req.Limit, req.Skip)
	rsp := helper.ToMap(meta, usersList, "users")

	response.HandleSuccess(ctx, rsp)
}

// getUserRequest represents the request body for getting a user
type getUserRequest struct {
	ID uint64 `uri:"id" binding:"required,min=1" example:"1"`
}

// GetUser godoc
//
//	@Summary		Get a user
//	@Description	Get a user by id
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		uint64					true	"User ID"
//	@Success		200	{object}	response.UserResponse	"User displayed"
//	@Failure		400	{object}	response.ErrorResponse	"Validation error"
//	@Failure		404	{object}	response.ErrorResponse	"Data not found error"
//	@Failure		500	{object}	response.ErrorResponse	"Internal server error"
//	@Router			/users/{id} [get]
//	@Security		BearerAuth
func (uh *UserHandler) GetUser(ctx *gin.Context) {
	var req getUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		response.ValidationError(ctx, err)
		return
	}

	user, err := uh.svc.GetUser(ctx, req.ID)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	rsp := response.NewUserResponse(user)

	response.HandleSuccess(ctx, rsp)
}

// updateUserRequest represents the request body for updating a user
type updateUserRequest struct {
	Name  string              `json:"name" binding:"omitempty,required" example:"John Doe"`
	Email string              `json:"email" binding:"omitempty,required,email" example:"test@example.com"`
	Role  domain.UserRoleEnum `json:"role" binding:"omitempty,required,user_role" example:"admin"`
}

// UpdateUser godoc
//
//	@Summary		Update a user
//	@Description	Update a user's name, email, password, or role by id
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id					path		uint64					true	"User ID"
//	@Param			updateUserRequest	body		updateUserRequest		true	"Update user request"
//	@Success		200					{object}	response.UserResponse	"User updated"
//	@Failure		400					{object}	response.ErrorResponse	"Validation error"
//	@Failure		401					{object}	response.ErrorResponse	"Unauthorized error"
//	@Failure		403					{object}	response.ErrorResponse	"Forbidden error"
//	@Failure		404					{object}	response.ErrorResponse	"Data not found error"
//	@Failure		500					{object}	response.ErrorResponse	"Internal server error"
//	@Router			/users/{id} [put]
//	@Security		BearerAuth
func (uh *UserHandler) UpdateUser(ctx *gin.Context) {
	var req updateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidationError(ctx, err)
		return
	}

	idStr := ctx.Param("id")
	id, err := helper.StringToUint64(idStr)
	if err != nil {
		response.ValidationError(ctx, err)
		return
	}

	user := domain.User{
		ID:    id,
		Name:  req.Name,
		Email: req.Email,
		Role:  req.Role,
	}

	updatedUser, err := uh.svc.UpdateUser(ctx, &user)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	rsp := response.NewUserResponse(updatedUser)

	response.HandleSuccess(ctx, rsp)
}

// deleteUserRequest represents the request body for deleting a user
type deleteUserRequest struct {
	ID uint64 `uri:"id" binding:"required,min=1" example:"1"`
}

// DeleteUser godoc
//
//	@Summary		Delete a user
//	@Description	Delete a user by id
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		uint64					true	"User ID"
//	@Success		200	{object}	response.Response		"User deleted"
//	@Failure		400	{object}	response.ErrorResponse	"Validation error"
//	@Failure		401	{object}	response.ErrorResponse	"Unauthorized error"
//	@Failure		403	{object}	response.ErrorResponse	"Forbidden error"
//	@Failure		404	{object}	response.ErrorResponse	"Data not found error"
//	@Failure		500	{object}	response.ErrorResponse	"Internal server error"
//	@Router			/users/{id} [delete]
//	@Security		BearerAuth
func (uh *UserHandler) DeleteUser(ctx *gin.Context) {
	var req deleteUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		response.ValidationError(ctx, err)
		return
	}

	err := uh.svc.DeleteUser(ctx, req.ID)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.HandleSuccess(ctx, nil)
}
