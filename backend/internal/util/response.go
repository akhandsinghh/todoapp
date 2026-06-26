package util

import (
	"strconv"
	apperr "todo-app/backend/internal/errors"

	"github.com/gin-gonic/gin"
)

// SuccessResponse defines the standardized success payload.
type SuccessResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Model   any    `json:"model"`
}

// ErrorResponse defines the standardized error payload.
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Model   any    `json:"model"`
}

func JSON(ctx *gin.Context, status int, payload any) {
	ctx.JSON(status, payload)
}

func Success(ctx *gin.Context, status int, message string, model any) {
	JSON(ctx, status, SuccessResponse{Code: strconv.Itoa(status), Message: message, Model: model})
}

func Error(ctx *gin.Context, status int, message string) {
	JSON(ctx, status, ErrorResponse{Code: strconv.Itoa(status), Message: message, Model: nil})
}

// HandleError extracts the status code from AppError if available, otherwise defaults to 500
func HandleError(ctx *gin.Context, err error) {
	statusCode := 500
	message := err.Error()

	if appErr, ok := err.(*apperr.AppError); ok {
		statusCode = appErr.StatusCode
	}

	Error(ctx, statusCode, message)
}

func Decode(ctx *gin.Context, target any) error {
	return ctx.ShouldBindJSON(target)
}
