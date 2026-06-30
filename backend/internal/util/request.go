package util

import (
    "strconv"
    apperr "todo-app/backend/internal/errors"
    "github.com/gin-gonic/gin"
)

// Extracting the "id" parameter from the URL type casting to int64
func ParsePathID(ctx *gin.Context) (int64, bool) {
    id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
    if err != nil || id == 0 {
        HandleError(ctx, apperr.BadRequest("invalid id"))
        return 0, false
    }
    return id, true
}