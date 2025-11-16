package main

import (
	"context"
	"os"

	"pull-request-review/config"
	"pull-request-review/internal/app"
	"pull-request-review/internal/infrastructure/adapters/logger"
	"pull-request-review/internal/infrastructure/database"
)

func main() {
	log := logger.NewZerologLogger()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error(err, "Failed to load configuration")
		os.Exit(1)
	}

	log.Info(
		"Starting PR Reviewer Service",
		logger.F("port", cfg.Server.Port),
	)

	ctx := context.Background()
	db := database.NewDatabase(cfg.Database, log)
	if err := db.Connect(ctx); err != nil {
		log.Error(err, "Failed to connect to database")
		os.Exit(1)
	}
	defer db.Close()

	app.Run(cfg, db, log)
}