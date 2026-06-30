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

type AuthServiceInterface interface {
	Register(ctx context.Context, req dto.RegisterRequest) (dto.AuthResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (dto.AuthResponse, error)
	Me(ctx context.Context, userID int64) (model.UserResponse, error)
	ChangePassword(ctx context.Context, userID int64, req dto.ChangePasswordRequest) error
	ForgotPassword(ctx context.Context, req dto.ForgotPasswordRequest) error
}

type AuthController struct {
	service AuthServiceInterface
}

func NewAuthController(s AuthServiceInterface) *AuthController {
	return &AuthController{service: s}
}

func (c *AuthController) Register(ctx *gin.Context) {
	var req dto.RegisterRequest
	if err := util.Decode(ctx, &req); err != nil {
		util.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	res, err := c.service.Register(ctx.Request.Context(), req)
	if err != nil {
		util.HandleError(ctx, err)
		return
	}
	util.Success(ctx, 201, "registered successfully", res)
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req dto.LoginRequest
	if err := util.Decode(ctx, &req); err != nil {
		util.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	res, err := c.service.Login(ctx.Request.Context(), req)
	if err != nil {
		util.HandleError(ctx, err)
		return
	}
	util.Success(ctx, 200, "logged in successfully", res)
}

func (c *AuthController) Me(ctx *gin.Context) {
	res, err := c.service.Me(ctx.Request.Context(), middleware.UserID(ctx))
	if err != nil {
		util.HandleError(ctx, err)
		return
	}
	util.Success(ctx, 200, "user fetched successfully", dto.ConvertUserDomainToResponse(res))
}

func (c *AuthController) ChangePassword(ctx *gin.Context) {
	var req dto.ChangePasswordRequest
	if err := util.Decode(ctx, &req); err != nil {
		util.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if err := c.service.ChangePassword(ctx.Request.Context(), middleware.UserID(ctx), req); err != nil {
		util.HandleError(ctx, err)
		return
	}
	util.Success(ctx, 200, "password updated", dto.MessageResponse{Message: "password updated"})
}

func (c *AuthController) ForgotPassword(ctx *gin.Context) {
	var req dto.ForgotPasswordRequest
	if err := util.Decode(ctx, &req); err != nil {
		util.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if err := c.service.ForgotPassword(ctx.Request.Context(), req); err != nil {
		util.HandleError(ctx, err)
		return
	}
	util.Success(ctx, 200, "password reset", dto.MessageResponse{Message: "password reset"})
}
