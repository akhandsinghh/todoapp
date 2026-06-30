package controller

import (
	"context"
	"net/http"

	"todo-app/backend/internal/controller/dto"
	"todo-app/backend/internal/middleware"
	"todo-app/backend/internal/model"
	"todo-app/backend/internal/util"

	"github.com/gin-gonic/gin"
)

type GroupServiceInterface interface {
	List(ctx context.Context, userID int64) ([]model.GroupDTO, error)
	Create(ctx context.Context, userID int64, req dto.GroupRequest) (model.GroupDTO, error)
	Update(ctx context.Context, userID, id int64, req dto.GroupUpdateRequest) error
	Share(ctx context.Context, userID, id int64, req dto.ShareGroupRequest) error
	Delete(ctx context.Context, userID, id int64) error
}

type GroupController struct {
	service GroupServiceInterface
}

func NewGroupController(s GroupServiceInterface) *GroupController {
	return &GroupController{service: s}
}

func (c *GroupController) List(ctx *gin.Context) {
	res, err := c.service.List(ctx.Request.Context(), middleware.UserID(ctx))
	if err != nil {
		util.HandleError(ctx, err)
		return
	}
	util.Success(ctx, 200, "groups fetched successfully", dto.ConvertGroupListDomainToResponse(res))
}

func (c *GroupController) Create(ctx *gin.Context) {
	var req dto.GroupRequest
	if err := util.Decode(ctx, &req); err != nil {
		util.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	res, err := c.service.Create(ctx.Request.Context(), middleware.UserID(ctx), req)
	if err != nil {
		util.HandleError(ctx, err)
		return
	}
	util.Success(ctx, 201, "group created successfully", dto.ConvertGroupDomainToResponse(res))
}

func (c *GroupController) Update(ctx *gin.Context) {
	id, ok := util.ParsePathID(ctx)
	if !ok {
		return
	}
	var req dto.GroupUpdateRequest
	if err := util.Decode(ctx, &req); err != nil {
		util.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if err := c.service.Update(ctx.Request.Context(), middleware.UserID(ctx), id, req); err != nil {
		util.HandleError(ctx, err)
		return
	}
	util.Success(ctx, 200, "group updated", dto.MessageResponse{Message: "group updated"})
}

func (c *GroupController) Share(ctx *gin.Context) {
	id, ok := util.ParsePathID(ctx)
	if !ok {
		return
	}
	var req dto.ShareGroupRequest
	if err := util.Decode(ctx, &req); err != nil {
		util.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if err := c.service.Share(ctx.Request.Context(), middleware.UserID(ctx), id, req); err != nil {
		util.HandleError(ctx, err)
		return
	}
	util.Success(ctx, 200, "group shared", dto.MessageResponse{Message: "group shared"})
}

func (c *GroupController) Delete(ctx *gin.Context) {
	id, ok := util.ParsePathID(ctx)
	if !ok {
		return
	}
	if err := c.service.Delete(ctx.Request.Context(), middleware.UserID(ctx), id); err != nil {
		util.HandleError(ctx, err)
		return
	}
	util.Success(ctx, 200, "group deleted", dto.MessageResponse{Message: "group deleted"})
}
