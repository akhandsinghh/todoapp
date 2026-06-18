package controller

import (
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
		util.Error(ctx, 400, "invalid json")
		return
	}
	res, err := c.service.Register(ctx.Request.Context(), req)
	if err != nil {
		util.Error(ctx, 400, err.Error())
		return
	}
	util.JSON(ctx, 201, res)
}
func (c *AuthController) Login(ctx *gin.Context) {
	var req model.LoginRequest
	if err := util.Decode(ctx, &req); err != nil {
		util.Error(ctx, 400, "invalid json")
		return
	}
	res, err := c.service.Login(ctx.Request.Context(), req)
	if err != nil {
		util.Error(ctx, 401, err.Error())
		return
	}
	util.JSON(ctx, 200, res)
}
func (c *AuthController) Me(ctx *gin.Context) {
	res, err := c.service.Me(ctx.Request.Context(), middleware.UserID(ctx))
	if err != nil {
		util.Error(ctx, 404, "user not found")
		return
	}
	util.JSON(ctx, 200, res)
}
