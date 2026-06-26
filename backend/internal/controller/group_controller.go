package controller

import (
	"strconv"
	apperr "todo-app/backend/internal/errors"
	"todo-app/backend/internal/middleware"
	"todo-app/backend/internal/model"
	"todo-app/backend/internal/service"
	"todo-app/backend/internal/util"

	"github.com/gin-gonic/gin"
)

type GroupController struct{ service *service.GroupService }

func NewGroupController(s *service.GroupService) *GroupController {
	return &GroupController{service: s}
}
func (c *GroupController) List(ctx *gin.Context) {
	res, err := c.service.List(ctx.Request.Context(), middleware.UserID(ctx))
	if err != nil {
		util.HandleError(ctx, err)
		return
	}
	util.Success(ctx, 200, "groups fetched successfully", res)
}
func (c *GroupController) Create(ctx *gin.Context) {
	var req model.GroupRequest
	if err := util.Decode(ctx, &req); err != nil {
		util.HandleError(ctx, apperr.BadRequest("invalid json"))
		return
	}
	res, err := c.service.Create(ctx.Request.Context(), middleware.UserID(ctx), req)
	if err != nil {
		util.HandleError(ctx, err)
		return
	}
	util.Success(ctx, 201, "group created successfully", res)
}
func (c *GroupController) Update(ctx *gin.Context) {
	id, ok := pathID(ctx)
	if !ok {
		return
	}
	var req model.GroupRequest
	if err := util.Decode(ctx, &req); err != nil {
		util.HandleError(ctx, apperr.BadRequest("invalid json"))
		return
	}
	if err := c.service.Update(ctx.Request.Context(), middleware.UserID(ctx), id, req); err != nil {
		util.HandleError(ctx, err)
		return
	}
	util.Success(ctx, 200, "group updated", model.MessageResponse{Message: "group updated"})
}
func (c *GroupController) Share(ctx *gin.Context) {
	id, ok := pathID(ctx)
	if !ok {
		return
	}
	var req model.ShareGroupRequest
	if err := util.Decode(ctx, &req); err != nil {
		util.HandleError(ctx, apperr.BadRequest("invalid json"))
		return
	}
	if err := c.service.Share(ctx.Request.Context(), middleware.UserID(ctx), id, req); err != nil {
		util.HandleError(ctx, err)
		return
	}
	util.Success(ctx, 200, "group shared", model.MessageResponse{Message: "group shared"})
}
func (c *GroupController) Delete(ctx *gin.Context) {
	id, ok := pathID(ctx)
	if !ok {
		return
	}
	if err := c.service.Delete(ctx.Request.Context(), middleware.UserID(ctx), id); err != nil {
		util.HandleError(ctx, err)
		return
	}
	util.Success(ctx, 200, "group deleted", model.MessageResponse{Message: "group deleted"})
}
func pathID(ctx *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil || id == 0 {
		util.HandleError(ctx, apperr.BadRequest("invalid id"))
		return 0, false
	}
	return id, true
}
