package controller

import (
	apperr "todo-app/backend/internal/errors"
	"todo-app/backend/internal/middleware"
	"todo-app/backend/internal/model"
	"todo-app/backend/internal/service"
	"todo-app/backend/internal/util"

	"github.com/gin-gonic/gin"
)

type AuthController struct{ service *service.AuthService }

func NewAuthController(s *service.AuthService) *AuthController { return &AuthController{service: s} }
func (c *AuthController) Register(ctx *gin.Context) {
	var req model.RegisterRequest
	if err := util.Decode(ctx, &req); err != nil {
		util.HandleError(ctx, apperr.BadRequest("invalid json"))
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
	var req model.LoginRequest
	if err := util.Decode(ctx, &req); err != nil {
		util.HandleError(ctx, apperr.BadRequest("invalid json"))
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
	util.Success(ctx, 200, "user fetched successfully", res)
}

func (c *AuthController) ChangePassword(ctx *gin.Context) {
	var req model.ChangePasswordRequest
	if err := util.Decode(ctx, &req); err != nil {
		util.HandleError(ctx, apperr.BadRequest("invalid json"))
		return
	}
	if err := c.service.ChangePassword(ctx.Request.Context(), middleware.UserID(ctx), req); err != nil {
		util.HandleError(ctx, err)
		return
	}
	util.Success(ctx, 200, "password updated", model.MessageResponse{Message: "password updated"})
}

func (c *AuthController) ForgotPassword(ctx *gin.Context) {
	var req model.ForgotPasswordRequest
	if err := util.Decode(ctx, &req); err != nil {
		util.HandleError(ctx, apperr.BadRequest("invalid json"))
		return
	}
	if err := c.service.ForgotPassword(ctx.Request.Context(), req); err != nil {
		util.HandleError(ctx, err)
		return
	}
	util.Success(ctx, 200, "password reset", model.MessageResponse{Message: "password reset"})
}
