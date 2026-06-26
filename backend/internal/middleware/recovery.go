package middleware

import (
	apperr "todo-app/backend/internal/errors"
	"todo-app/backend/internal/util"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				logger.Error("panic recovered: %v", rec)
				util.HandleError(ctx, apperr.Internal("internal server error"))
				ctx.Abort()
			}
		}()
		ctx.Next()
	}
}
