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

type ReminderServiceInterface interface {
	List(ctx context.Context, userID int64) ([]model.ReminderDTO, error)
	Create(ctx context.Context, userID int64, req dto.ReminderRequest) (model.ReminderDTO, error)
	Delete(ctx context.Context, userID, id int64) error
}

type ReminderController struct {
	service ReminderServiceInterface
}

func NewReminderController(s ReminderServiceInterface) *ReminderController {
	return &ReminderController{service: s}
}

func (c *ReminderController) List(ctx *gin.Context) {
	res, err := c.service.List(ctx.Request.Context(), middleware.UserID(ctx))
	if err != nil {
		util.HandleError(ctx, err)
		return
	}
	util.Success(ctx, 200, "reminders fetched successfully", dto.ConvertReminderListDomainToResponse(res))
}

func (c *ReminderController) Create(ctx *gin.Context) {
	var req dto.ReminderRequest
	if err := util.Decode(ctx, &req); err != nil {
		util.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	res, err := c.service.Create(ctx.Request.Context(), middleware.UserID(ctx), req)
	if err != nil {
		util.HandleError(ctx, err)
		return
	}
	util.Success(ctx, 201, "reminder created successfully", dto.ConvertReminderDomainToResponse(res))
}

func (c *ReminderController) Delete(ctx *gin.Context) {
	id, ok := util.ParsePathID(ctx)
	if !ok {
		return
	}
	if err := c.service.Delete(ctx.Request.Context(), middleware.UserID(ctx), id); err != nil {
		util.HandleError(ctx, err)
		return
	}
	util.Success(ctx, 200, "reminder deleted", dto.MessageResponse{Message: "reminder deleted"})
}
