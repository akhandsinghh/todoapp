package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"todo-app/backend/internal/controller"
	"todo-app/backend/internal/db"
	"todo-app/backend/internal/db/sqlc"
	"todo-app/backend/internal/repository"
	"todo-app/backend/internal/routes"
	"todo-app/backend/internal/scheduler"
	"todo-app/backend/internal/service"
)

func main() {
	root := backendRoot()
	db.LoadEnv(filepath.Join(root, ".env"))
	database, err := db.Connect(db.FromEnv())
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer database.Close()
	if err := db.RunMigrations(database, filepath.Join(root, "internal", "db", "migrations")); err != nil {
		log.Fatalf("migrations failed: %v", err)
	}

	q := sqlc.New(database)
	userRepo := repository.NewUserRepository(q)
	groupRepo := repository.NewGroupRepository(q)
	taskRepo := repository.NewTaskRepository(q)
	reminderRepo := repository.NewReminderRepository(q)

	secret := env("JWT_SECRET", "change-this-secret")
	authSvc := service.NewAuthService(userRepo, secret)
	groupSvc := service.NewGroupService(groupRepo)
	taskSvc := service.NewTaskService(taskRepo)
	reminderSvc := service.NewReminderService(reminderRepo, taskRepo)

	scheduler.NewReminderScheduler(reminderSvc).Start(context.Background())

	router := routes.New(routes.Controllers{
		Auth:      controller.NewAuthController(authSvc),
		Tasks:     controller.NewTaskController(taskSvc),
		Groups:    controller.NewGroupController(groupSvc),
		Reminders: controller.NewReminderController(reminderSvc),
	}, routes.Config{
		JWTSecret:         secret,
		CORSAllowedOrigin: env("CORS_ALLOWED_ORIGIN", "http://localhost:3000"),
	})

	addr := ":" + env("APP_PORT", "8080")
	log.Printf("backend listening on %s", addr)
	log.Fatal(router.Run(addr))
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func backendRoot() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "."
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}
