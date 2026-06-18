package middleware

import (
	"log"
	"todo-app/backend/internal/util"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("panic: %v", rec)
				util.Error(ctx, 500, "internal server error")
				ctx.Abort()
			}
		}()
		ctx.Next()
	}
}
