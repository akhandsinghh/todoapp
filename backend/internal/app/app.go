package app

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"path/filepath"

	"todo-app/backend/internal/config"
	"todo-app/backend/internal/controller"
	database "todo-app/backend/internal/db"
	"todo-app/backend/internal/db/sqlc"
	"todo-app/backend/internal/middleware"
	"todo-app/backend/internal/repository"
	"todo-app/backend/internal/scheduler"
	"todo-app/backend/internal/service"

	"github.com/gin-gonic/gin"
)

// App is the core application structure.
type App struct {
	config    *config.AppConfig
	db        *sql.DB
	router    *gin.Engine
	server    *http.Server
	connector database.DatabaseConnector
}

// NewApp creates a new App instance with the provided config dependencies.
func NewApp(cfg *config.AppConfig, db *sql.DB) *App {
	return &App{config: cfg, db: db}
}

// NewAppWithConnector creates a new App instance and injects a database connector.
func NewAppWithConnector(cfg *config.AppConfig, connector database.DatabaseConnector, db *sql.DB) *App {
	if connector == nil && cfg != nil && cfg.Database() != nil {
		connector = database.NewDatabaseConnector(*cfg.Database())
	}
	return &App{config: cfg, connector: connector, db: db}
}

// Init initializes the application and returns the configured app.
func (a *App) Init() *App {
	if a.db == nil {
		if a.connector == nil {
			a.connector = database.NewDatabaseConnector(*a.config.Database())
		}

		databaseConn, err := a.connector.Connect()
		if err != nil {
			log.Fatalf("database connection failed: %v", err)
		}
		a.db = databaseConn
	}

	if a.connector == nil && a.config != nil && a.config.Database() != nil {
		a.connector = database.NewDatabaseConnector(*a.config.Database())
	}
	if a.connector != nil {
		root := config.BackendRoot()
		if err := a.connector.Migrate(a.db, filepath.Join(root, "internal", "db", "migrations")); err != nil {
			log.Fatalf("migrations failed: %v", err)
		}
	}

	// Initialize database queries
	q := sqlc.New(a.db)
	userRepo := repository.NewUserRepository(q)
	groupRepo := repository.NewGroupRepository(q)
	taskRepo := repository.NewTaskRepository(q)
	reminderRepo := repository.NewReminderRepository(q)

	// Initialize services using config
	authSvc := service.NewAuthService(userRepo, a.config.Auth().JWTSecret())
	groupSvc := service.NewGroupService(groupRepo)
	taskSvc := service.NewTaskService(taskRepo)
	reminderSvc := service.NewReminderService(reminderRepo, taskRepo)

	// Start scheduler in background
	scheduler.NewReminderScheduler(reminderSvc).Start(context.Background())

	// Build router inside app.go
	a.router = gin.New()
	a.router.Use(middleware.Recovery(), middleware.Logging(), middleware.CORS(a.config.Server().CORSOrigin()))
	a.router.GET("/health", func(ctx *gin.Context) { ctx.String(200, "ok") })

	authCtrl := controller.NewAuthController(authSvc)
	groupCtrl := controller.NewGroupController(groupSvc)
	taskCtrl := controller.NewTaskController(taskSvc)
	reminderCtrl := controller.NewReminderController(reminderSvc)

	api := a.router.Group("/api")
	api.POST("/auth/register", authCtrl.Register)
	api.POST("/auth/login", authCtrl.Login)
	api.POST("/auth/forgot-password", authCtrl.ForgotPassword)

	authenticated := api.Group("")
	authenticated.Use(middleware.Auth(a.config.Auth().JWTSecret()))
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
	if a.router == nil {
		a.Init()
	}
	if a.server == nil {
		a.server = &http.Server{
			Addr:    ":" + a.config.Server().Port(),
			Handler: a.router,
		}
	}
	return a.server.ListenAndServe()
}

// Stop gracefully shuts down the HTTP server.
func (a *App) Stop(ctx context.Context) error {
	if a.server == nil {
		return nil
	}
	log.Println("shutting down HTTP server")
	return a.server.Shutdown(ctx)
}
