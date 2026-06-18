package middleware

import (
	"strings"
	"todo-app/backend/internal/util"

	"github.com/gin-gonic/gin"
)

const UserIDKey = "userID"

func Auth(secret string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			util.Error(ctx, 401, "missing bearer token")
			ctx.Abort()
			return
		}
		claims, err := util.VerifyToken(strings.TrimPrefix(header, "Bearer "), secret)
		if err != nil {
			util.Error(ctx, 401, "invalid token")
			ctx.Abort()
			return
		}
		ctx.Set(UserIDKey, claims.UserID)
		ctx.Next()
	}
}

func UserID(ctx *gin.Context) int64 {
	id, _ := ctx.Get(UserIDKey)
	value, _ := id.(int64)
	return value
}
