package routes

import (
	"todo-app/backend/internal/controller"
	"todo-app/backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

type Controllers struct {
	Auth      *controller.AuthController
	Tasks     *controller.TaskController
	Groups    *controller.GroupController
	Reminders *controller.ReminderController
}

type Config struct {
	JWTSecret         string
	CORSAllowedOrigin string
}

func New(c Controllers, cfg Config) *gin.Engine {
	router := gin.New()
	router.Use(middleware.Recovery(), middleware.Logging(), middleware.CORS(cfg.CORSAllowedOrigin))

	router.GET("/health", func(ctx *gin.Context) { ctx.String(200, "ok") })

	api := router.Group("/api")
	api.POST("/auth/register", c.Auth.Register)
	api.POST("/auth/login", c.Auth.Login)

	authenticated := api.Group("")
	authenticated.Use(middleware.Auth(cfg.JWTSecret))
	authenticated.GET("/auth/me", c.Auth.Me)

	authenticated.GET("/tasks", c.Tasks.List)
	authenticated.POST("/tasks", c.Tasks.Create)
	authenticated.PUT("/tasks/:id", c.Tasks.Update)
	authenticated.PATCH("/tasks/:id", c.Tasks.Update)
	authenticated.DELETE("/tasks/:id", c.Tasks.Delete)

	authenticated.GET("/groups", c.Groups.List)
	authenticated.POST("/groups", c.Groups.Create)
	authenticated.PUT("/groups/:id", c.Groups.Update)
	authenticated.PATCH("/groups/:id", c.Groups.Update)
	authenticated.DELETE("/groups/:id", c.Groups.Delete)

	authenticated.GET("/reminders", c.Reminders.List)
	authenticated.POST("/reminders", c.Reminders.Create)
	authenticated.DELETE("/reminders/:id", c.Reminders.Delete)

	return router
}
