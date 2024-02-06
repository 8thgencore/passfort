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

// CollectionHandler represents the HTTP handler for collection-related requests
type CollectionHandler struct {
	svc service.CollectionService
}

// NewCollectionHandler creates a new CollectionHandler instance
func NewCollectionHandler(svc service.CollectionService) *CollectionHandler {
	return &CollectionHandler{
		svc,
	}
}

// createCollectionRequest represents the request body for creating a collection
type createCollectionRequest struct {
	Name        string `json:"name" binding:"required" example:"My Collection"`
	Description string `json:"description" example:"A collection of items"`
}

// CreateCollection godoc
//
//	@Summary		Create a new collection
//	@Description	Create a new collection
//	@Tags			Collections
//	@Accept			json
//	@Produce		json
//	@Param			request	body		createCollectionRequest		true	"Create Collection Request"
//	@Success		201		{object}	response.CollectionResponse	"Collection created"
//	@Failure		400		{object}	response.ErrorResponse		"Validation error"
//	@Failure		500		{object}	response.ErrorResponse		"Internal server error"
//	@Router			/collections [post]
//	@Security		BearerAuth
func (ch *CollectionHandler) CreateCollection(ctx *gin.Context) {
	var req createCollectionRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidationError(ctx, err)
		return
	}

	// Assuming your CollectionService has a method CreateCollection
	newCollection := domain.Collection{
		Name:        req.Name,
		Description: req.Description,
	}

	authPayload := helper.GetAuthPayload(ctx, middleware.AuthorizationPayloadKey)

	createdCollection, err := ch.svc.CreateCollection(ctx, authPayload.UserID, &newCollection)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	rsp := response.NewCollectionResponse(createdCollection)

	response.HandleSuccess(ctx, rsp)
}

// listMeCollectionsRequest represents the request body for listing collections by user ID
type listMeCollectionsRequest struct {
	Skip  uint64 `form:"skip" binding:"required,min=0" example:"0"`
	Limit uint64 `form:"limit" binding:"required,min=5" example:"5"`
}

// ListMeCollections godoc
//
//	@Summary		List me collections
//	@Description	List me collections associated with pagination
//	@Tags			Collections
//	@Accept			json
//	@Produce		json
//	@Param			skip	query		uint64					true	"Skip"
//	@Param			limit	query		uint64					true	"Limit"
//	@Success		200		{object}	response.Meta			"Collections displayed"
//	@Failure		400		{object}	response.ErrorResponse	"Validation error"
//	@Failure		500		{object}	response.ErrorResponse	"Internal server error"
//	@Router			/collections/me [get]
//	@Security		BearerAuth
func (ch *CollectionHandler) ListMeCollections(ctx *gin.Context) {
	var req listMeCollectionsRequest
	var collectionsList []response.CollectionResponse

	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.ValidationError(ctx, err)
		return
	}

	authPayload := helper.GetAuthPayload(ctx, middleware.AuthorizationPayloadKey)

	// Assuming your CollectionService has a method ListCollectionsByUserID
	collections, err := ch.svc.ListCollectionsByUserID(ctx, authPayload.UserID, req.Skip, req.Limit)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	for _, collection := range collections {
		collectionsList = append(collectionsList, response.NewCollectionResponse(&collection))
	}

	total := uint64(len(collectionsList))
	meta := response.NewMeta(total, req.Limit, req.Skip)
	rsp := helper.ToMap(meta, collectionsList, "collections")

	response.HandleSuccess(ctx, rsp)
}

// getCollectionRequest represents the request body for getting a collection
type getCollectionRequest struct {
	ID string `uri:"id" binding:"required" example:"5950a459-5126-40b7-bd8e-82f7b91c2cf1"`
}

// GetCollection godoc
//
//	@Summary		Get a collection
//	@Description	Get a collection by id
//	@Tags			Collections
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string						true	"Collection ID"
//	@Success		200	{object}	response.CollectionResponse	"Collection displayed"
//	@Failure		400	{object}	response.ErrorResponse		"Validation error"
//	@Failure		404	{object}	response.ErrorResponse		"Data not found error"
//	@Failure		500	{object}	response.ErrorResponse		"Internal server error"
//	@Router			/collections/{id} [get]
//	@Security		BearerAuth
func (ch *CollectionHandler) GetCollection(ctx *gin.Context) {
	var req getCollectionRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		response.ValidationError(ctx, err)
		return
	}

	uuidString := req.ID
	uuid, err := uuid.Parse(uuidString)
	if err != nil {
		response.ValidationError(ctx, err)
		return
	}

	authPayload := helper.GetAuthPayload(ctx, middleware.AuthorizationPayloadKey)

	collection, err := ch.svc.GetCollection(ctx, authPayload.UserID, uuid)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	rsp := response.NewCollectionResponse(collection)

	response.HandleSuccess(ctx, rsp)
}

// updateCollectionRequest represents the request body for updating a collection
type updateCollectionRequest struct {
	Name        string `json:"name" binding:"omitempty,required" example:"My Collection"`
	Description string `json:"description,omitempty" example:"Collection description"`
}

// UpdateCollection godoc
//
//	@Summary		Update a collection
//	@Description	Update a collection's name or description by id
//	@Tags			Collections
//	@Accept			json
//	@Produce		json
//	@Param			id						path		string						true	"Collection ID"
//	@Param			updateCollectionRequest	body		updateCollectionRequest		true	"Update collection request"
//	@Success		200						{object}	response.CollectionResponse	"Collection updated"
//	@Failure		400						{object}	response.ErrorResponse		"Validation error"
//	@Failure		401						{object}	response.ErrorResponse		"Unauthorized error"
//	@Failure		403						{object}	response.ErrorResponse		"Forbidden error"
//	@Failure		404						{object}	response.ErrorResponse		"Data not found error"
//	@Failure		500						{object}	response.ErrorResponse		"Internal server error"
//	@Router			/collections/{id} [put]
//	@Security		BearerAuth
func (ch *CollectionHandler) UpdateCollection(ctx *gin.Context) {
	var req updateCollectionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ValidationError(ctx, err)
		return
	}

	uuidString := ctx.Param("id")
	uuid, err := uuid.Parse(uuidString)
	if err != nil {
		response.ValidationError(ctx, err)
		return
	}

	collection := domain.Collection{
		ID:          uuid,
		Name:        req.Name,
		Description: req.Description,
	}

	authPayload := helper.GetAuthPayload(ctx, middleware.AuthorizationPayloadKey)

	updatedCollection, err := ch.svc.UpdateCollection(ctx, authPayload.UserID, &collection)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	rsp := response.NewCollectionResponse(updatedCollection)

	response.HandleSuccess(ctx, rsp)
}

// deleteCollectionRequest represents the request body for deleting a collection
type deleteCollectionRequest struct {
	ID string `uri:"id" binding:"required" example:"5950a459-5126-40b7-bd8e-82f7b91c2cf1"`
}

// DeleteCollection godoc
//
//	@Summary		Delete a collection
//	@Description	Delete a collection by id
//	@Tags			Collections
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string					true	"Collection ID"
//	@Success		200	{object}	response.Response		"Collection deleted"
//	@Failure		400	{object}	response.ErrorResponse	"Validation error"
//	@Failure		401	{object}	response.ErrorResponse	"Unauthorized error"
//	@Failure		403	{object}	response.ErrorResponse	"Forbidden error"
//	@Failure		404	{object}	response.ErrorResponse	"Data not found error"
//	@Failure		500	{object}	response.ErrorResponse	"Internal server error"
//	@Router			/collections/{id} [delete]
//	@Security		BearerAuth
func (ch *CollectionHandler) DeleteCollection(ctx *gin.Context) {
	var req deleteCollectionRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		response.ValidationError(ctx, err)
		return
	}

	uuidString := ctx.Param("id")
	uuid, err := uuid.Parse(uuidString)
	if err != nil {
		response.ValidationError(ctx, err)
		return
	}

	authPayload := helper.GetAuthPayload(ctx, middleware.AuthorizationPayloadKey)

	err = ch.svc.DeleteCollection(ctx, authPayload.UserID, uuid)
	if err != nil {
		response.HandleError(ctx, err)
		return
	}

	response.HandleSuccess(ctx, nil)
}
