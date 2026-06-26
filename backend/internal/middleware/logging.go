package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type Logger struct {
	prefix string
}

// NewLogger creates a new logger instance.
func NewLogger(prefix string) *Logger {
	return &Logger{prefix: prefix}
}

// Info logs an informational message with timestamp.
func (l *Logger) Info(format string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, args...)
	fmt.Fprintf(os.Stdout, "[%s] [INFO] %s: %s\n", timestamp, l.prefix, message)
}

// Error logs an error message with timestamp.
func (l *Logger) Error(format string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, args...)
	fmt.Fprintf(os.Stderr, "[%s] [ERROR] %s: %s\n", timestamp, l.prefix, message)
}

var logger = NewLogger("app")

func Logging() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		ctx.Next()
		duration := time.Since(start)
		logger.Info("%s %s %d %s", ctx.Request.Method, ctx.Request.URL.Path, ctx.Writer.Status(), duration)
	}
}
