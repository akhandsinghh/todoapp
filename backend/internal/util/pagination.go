package util

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func Pagination(ctx *gin.Context) (int32, int32) {
	limit := int32(5)
	page := int32(1)
	if v, err := strconv.Atoi(ctx.Query("limit")); err == nil && v > 0 && v <= 100 {
		limit = int32(v)
	}
	if v, err := strconv.Atoi(ctx.Query("page")); err == nil && v > 0 {
		page = int32(v)
	}
	return limit, (page - 1) * limit
}
