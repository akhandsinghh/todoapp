package controller

import (
	"todo-app/backend/internal/middleware"
	"todo-app/backend/internal/model"
	"todo-app/backend/internal/service"
	"todo-app/backend/internal/util"

	"github.com/gin-gonic/gin"
)

type ReminderController struct{ service *service.ReminderService }

func NewReminderController(s *service.ReminderService) *ReminderController {
	return &ReminderController{service: s}
}
func (c *ReminderController) List(ctx *gin.Context) {
	res, err := c.service.List(ctx.Request.Context(), middleware.UserID(ctx))
	if err != nil {
		util.Error(ctx, 500, err.Error())
		return
	}
	util.JSON(ctx, 200, res)
}
func (c *ReminderController) Create(ctx *gin.Context) {
	var req model.ReminderRequest
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
func (c *ReminderController) Delete(ctx *gin.Context) {
	id, ok := pathID(ctx)
	if !ok {
		return
	}
	if err := c.service.Delete(ctx.Request.Context(), middleware.UserID(ctx), id); err != nil {
		util.Error(ctx, 400, err.Error())
		return
	}
	util.JSON(ctx, 200, model.MessageResponse{Message: "reminder deleted"})
}
