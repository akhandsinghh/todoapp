package controller

import (
	"strconv"
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
	res, err := c.service.List(ctx.Request.Context(), middleware.UserID(ctx), ctx.Query("status"), gid, limit, offset)
	if err != nil {
		util.Error(ctx, 500, err.Error())
		return
	}
	util.JSON(ctx, 200, res)
}

func (c *TaskController) Create(ctx *gin.Context) {
	var req model.TaskRequest
	if err := util.Decode(ctx, &req); err != nil {
		util.Error(ctx, 400, "invalid json")
		return
	}
	res, err := c.service.Create(ctx.Request.Context(), middleware.UserID(ctx), req)
	if err != nil {
		util.Error(ctx, 400, err.Error())
		return
	}
	util.JSON(ctx, 201, res)
}
func (c *TaskController) Update(ctx *gin.Context) {
	id, ok := pathID(ctx)
	if !ok {
		return
	}
	var req model.TaskRequest
	if err := util.Decode(ctx, &req); err != nil {
		util.Error(ctx, 400, "invalid json")
		return
	}
	res, err := c.service.Update(ctx.Request.Context(), middleware.UserID(ctx), id, req)
	if err != nil {
		util.Error(ctx, 400, err.Error())
		return
	}
	util.JSON(ctx, 200, res)
}
func (c *TaskController) Delete(ctx *gin.Context) {
	id, ok := pathID(ctx)
	if !ok {
		return
	}
	if err := c.service.Delete(ctx.Request.Context(), middleware.UserID(ctx), id); err != nil {
		util.Error(ctx, 400, err.Error())
		return
	}
	util.JSON(ctx, 200, model.MessageResponse{Message: "task deleted"})
}
