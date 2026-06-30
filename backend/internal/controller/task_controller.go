package controller

import (
	"context"
	"net/http"
	"strconv"

	"todo-app/backend/internal/controller/dto"
	"todo-app/backend/internal/middleware"
	"todo-app/backend/internal/model"
	"todo-app/backend/internal/util"

	"github.com/gin-gonic/gin"
)

type TaskServiceInterface interface {
	List(ctx context.Context, userID int64, status string, groupID int64, ungrouped bool, limit, offset int32, sortBy, sortOrder string) (dto.TaskListResponse, error)
	Create(ctx context.Context, userID int64, req dto.TaskRequest) (model.TaskDTO, error)
	Update(ctx context.Context, userID, id int64, req dto.TaskUpdateRequest) (model.TaskDTO, error)
	Delete(ctx context.Context, userID, id int64) error
}

type TaskController struct {
	service TaskServiceInterface
}

func NewTaskController(s TaskServiceInterface) *TaskController {
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
	var req dto.TaskRequest
	if err := util.Decode(ctx, &req); err != nil {
		util.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	res, err := c.service.Create(ctx.Request.Context(), middleware.UserID(ctx), req)
	if err != nil {
		util.HandleError(ctx, err)
		return
	}
	util.Success(ctx, 201, "task created successfully", dto.ConvertTaskDomainToResponse(res))
}

func (c *TaskController) Update(ctx *gin.Context) {
	id, ok := util.ParsePathID(ctx)
	if !ok {
		return
	}
	var req dto.TaskUpdateRequest
	if err := util.Decode(ctx, &req); err != nil {
		util.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	res, err := c.service.Update(ctx.Request.Context(), middleware.UserID(ctx), id, req)
	if err != nil {
		util.HandleError(ctx, err)
		return
	}
	util.Success(ctx, 200, "task updated successfully", dto.ConvertTaskDomainToResponse(res))
}

func (c *TaskController) Delete(ctx *gin.Context) {
	id, ok := util.ParsePathID(ctx)
	if !ok {
		return
	}
	if err := c.service.Delete(ctx.Request.Context(), middleware.UserID(ctx), id); err != nil {
		util.HandleError(ctx, err)
		return
	}
	util.Success(ctx, 200, "task deleted", dto.MessageResponse{Message: "task deleted"})
}
