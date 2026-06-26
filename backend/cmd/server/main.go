package main

import (
	"log"
	"path/filepath"

	"todo-app/backend/internal/app"
	"todo-app/backend/internal/config"
)

func main() {
	root := config.BackendRoot()
	config.LoadEnv(filepath.Join(root, ".env"))

	appConfig := config.LoadConfig()
	application := app.NewApp(appConfig)
	application.Init()

	log.Printf("backend listening on :%s (environment: %s)", appConfig.Server.Port, appConfig.Server.AppEnv)
	log.Fatal(application.Start())
}
