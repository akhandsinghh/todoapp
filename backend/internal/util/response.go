package util

import "github.com/gin-gonic/gin"

type ErrorResponse struct {
	Error string `json:"error"`
}

func JSON(ctx *gin.Context, status int, payload any) {
	ctx.JSON(status, payload)
}

func Error(ctx *gin.Context, status int, message string) {
	JSON(ctx, status, ErrorResponse{Error: message})
}

func Decode(ctx *gin.Context, target any) error {
	return ctx.ShouldBindJSON(target)
}
