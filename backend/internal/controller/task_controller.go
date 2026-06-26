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

type TaskController struct{ service *service.TaskService }

func NewTaskController(s *service.TaskService) *TaskController {
	return &TaskController{service: s}
}

func (c *TaskController) List(ctx *gin.Context) {
	limit, offset := util.Pagination(ctx)
	gid, _ := strconv.ParseInt(ctx.Query("group_id"), 10, 64)
	ungrouped := ctx.Query("ungrouped") == "true" || ctx.Query("ungrouped") == "1"
	res, err := c.service.List(
		ctx.Request.Context(),
		middleware.UserID(ctx),
		ctx.Query("status"),
		gid,
		ungrouped,
		limit,
		offset,
		ctx.Query("sort_by"),
		ctx.Query("sort_order"),
	)
	if err != nil {
		util.HandleError(ctx, err)
		return
	}
	util.Success(ctx, 200, "tasks fetched successfully", res)
}

func (c *TaskController) Create(ctx *gin.Context) {
	var req model.TaskRequest
	if err := util.Decode(ctx, &req); err != nil {
		util.HandleError(ctx, apperr.BadRequest("invalid json"))
		return
	}
	res, err := c.service.Create(ctx.Request.Context(), middleware.UserID(ctx), req)
	if err != nil {
		util.HandleError(ctx, err)
		return
	}
	util.Success(ctx, 201, "task created successfully", res)
}
func (c *TaskController) Update(ctx *gin.Context) {
	id, ok := pathID(ctx)
	if !ok {
		return
	}
	var req model.TaskRequest
	if err := util.Decode(ctx, &req); err != nil {
		util.HandleError(ctx, apperr.BadRequest("invalid json"))
		return
	}
	res, err := c.service.Update(ctx.Request.Context(), middleware.UserID(ctx), id, req)
	if err != nil {
		util.HandleError(ctx, err)
		return
	}
	util.Success(ctx, 200, "task updated successfully", res)
}
func (c *TaskController) Delete(ctx *gin.Context) {
	id, ok := pathID(ctx)
	if !ok {
		return
	}
	if err := c.service.Delete(ctx.Request.Context(), middleware.UserID(ctx), id); err != nil {
		util.HandleError(ctx, err)
		return
	}
	util.Success(ctx, 200, "task deleted", model.MessageResponse{Message: "task deleted"})
}
