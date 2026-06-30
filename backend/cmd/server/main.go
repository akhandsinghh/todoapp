package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"todo-app/backend/internal/app"
	"todo-app/backend/internal/config"
	database "todo-app/backend/internal/db"
)

func main() {
	root := config.BackendRoot()
	config.LoadEnv(filepath.Join(root, ".env"))

	appConfig := config.LoadConfig()
	connector := database.NewDatabaseConnector(*appConfig.Database())
	dbConn, err := connector.Connect()
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer func() {
		if dbConn != nil {
			_ = dbConn.Close()
		}
	}()

	application := app.NewApp(appConfig, dbConn)
	application.Init()

	log.Printf("backend listening on :%s (environment: %s)", appConfig.Server().Port(), appConfig.Server().AppEnv())

	go func() {
		if err := application.Start(); err != nil && err != http.ErrServerClosed {
			log.Printf("server stopped: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := application.Stop(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
}
