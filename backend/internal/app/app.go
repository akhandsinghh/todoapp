package app

import (
	"context"
	"database/sql"
	"log"
	"path/filepath"

	"todo-app/backend/internal/config"
	"todo-app/backend/internal/controller"
	"todo-app/backend/internal/db/sqlc"
	"todo-app/backend/internal/middleware"
	"todo-app/backend/internal/repository"
	"todo-app/backend/internal/scheduler"
	"todo-app/backend/internal/service"

	"github.com/gin-gonic/gin"
)

// App is the core application structure.
type App struct {
	Config *config.AppConfig
	DB     *sql.DB
	Router *gin.Engine
}

// NewApp creates a new App instance with the provided config dependencies.
func NewApp(cfg *config.AppConfig) *App {
	return &App{
		Config: cfg,
	}
}

// Init initializes the application and returns the configured app.
func (a *App) Init() *App {
	// Build DB config from centralized application config
	    dbConfig := config.DatabaseConfig{
		Username: a.Config.Database.Username,
		Password: a.Config.Database.Password,
		Host:     a.Config.Database.Host,
		Port:     a.Config.Database.Port,
		Database: a.Config.Database.Database,
	}

	// Connect to the database using the provided config
	database, err := config.Connect(dbConfig)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	a.DB = database

	// Run migrations
	root := config.BackendRoot()
	if err := config.RunMigrations(database, filepath.Join(root, "internal", "db", "migrations")); err != nil {
		log.Fatalf("migrations failed: %v", err)
	}

	// Initialize database queries
	q := sqlc.New(database)
	userRepo := repository.NewUserRepository(q)
	groupRepo := repository.NewGroupRepository(q)
	taskRepo := repository.NewTaskRepository(q)
	reminderRepo := repository.NewReminderRepository(q)

	// Initialize services using config
	authSvc := service.NewAuthService(userRepo, a.Config.Auth.JWTSecret)
	groupSvc := service.NewGroupService(groupRepo)
	taskSvc := service.NewTaskService(taskRepo)
	reminderSvc := service.NewReminderService(reminderRepo, taskRepo)

	// Start scheduler in background
	scheduler.NewReminderScheduler(reminderSvc).Start(context.Background())

	// Build router inside app.go
	a.Router = gin.New()
	a.Router.Use(middleware.Recovery(), middleware.Logging(), middleware.CORS(a.Config.Server.CORSOrigin))
	a.Router.GET("/health", func(ctx *gin.Context) { ctx.String(200, "ok") })

	authCtrl := controller.NewAuthController(authSvc)
	groupCtrl := controller.NewGroupController(groupSvc)
	taskCtrl := controller.NewTaskController(taskSvc)
	reminderCtrl := controller.NewReminderController(reminderSvc)

	api := a.Router.Group("/api")
	api.POST("/auth/register", authCtrl.Register)
	api.POST("/auth/login", authCtrl.Login)
	api.POST("/auth/forgot-password", authCtrl.ForgotPassword)

	authenticated := api.Group("")
	authenticated.Use(middleware.Auth(a.Config.Auth.JWTSecret))
	authenticated.GET("/auth/me", authCtrl.Me)
	authenticated.POST("/auth/change-password", authCtrl.ChangePassword)

	authenticated.GET("/tasks", taskCtrl.List)
	authenticated.POST("/tasks", taskCtrl.Create)
	authenticated.PUT("/tasks/:id", taskCtrl.Update)
	authenticated.PATCH("/tasks/:id", taskCtrl.Update)
	authenticated.DELETE("/tasks/:id", taskCtrl.Delete)

	authenticated.GET("/groups", groupCtrl.List)
	authenticated.POST("/groups", groupCtrl.Create)
	authenticated.POST("/groups/:id/share", groupCtrl.Share)
	authenticated.PUT("/groups/:id", groupCtrl.Update)
	authenticated.PATCH("/groups/:id", groupCtrl.Update)
	authenticated.DELETE("/groups/:id", groupCtrl.Delete)

	authenticated.GET("/reminders", reminderCtrl.List)
	authenticated.POST("/reminders", reminderCtrl.Create)
	authenticated.DELETE("/reminders/:id", reminderCtrl.Delete)

	return a
}

// Start starts the HTTP server.
func (a *App) Start() error {
	if a.Router == nil {
		a.Init()
	}
	return a.Router.Run(":" + a.Config.Server.Port)
}
