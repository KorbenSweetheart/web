package main

import (
	"context"
	"log/slog"
	"match-me-api/internal/config"
	"match-me-api/internal/logger"
	"match-me-api/internal/storage/postgres"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	shutdownCtx, shutdown := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer shutdown()

	// init config
	cfg := config.MustLoad()

	// setup logger
	log := logger.SetupLogger(cfg.Env)

	log.Info("starting database migration service", slog.String("env", cfg.Env))

	// connect to repository
	dbCtx, dbCancel := context.WithTimeout(shutdownCtx, 20*time.Second)
	defer dbCancel()

	storage, err := postgres.NewPostgresDB(dbCtx, cfg.DB, log)
	if err != nil {
		log.Error("migration failed: database connection error", logger.Err(err))
		os.Exit(1)
	}

	// repository migrations
	log.Info("running auto-migrations...")

	migrateCtx, migrateCancel := context.WithTimeout(shutdownCtx, 20*time.Second)
	defer migrateCancel()

	if err := storage.AutoMigrate(migrateCtx); err != nil {
		log.Error("auto-migration failed", logger.Err(err))
		os.Exit(1)
	}

	log.Info("auto-migrations completed successfully")

	// Seeding dictionaries and test users into repository
	log.Info("seeding system dictionaries...")

	seedCtx, seedCancel := context.WithTimeout(shutdownCtx, 15*time.Second)
	defer seedCancel()

	if err := storage.SeedDictionaries(seedCtx); err != nil {
		log.Error("failed seeding dictionaries", logger.Err(err))
		os.Exit(1)
	}

	log.Info("system dictionaries seeded successfully")

	// Seeding dummy data and users for local and dev env.
	if cfg.Env == "local" || cfg.Env == "dev" || cfg.Env == "development" {
		log.Info("seeding test dummy users (development mode)...")
		dummyCtx, dummyCancel := context.WithTimeout(shutdownCtx, 15*time.Second)
		defer dummyCancel()

		if err := storage.SeedDummyUsers(dummyCtx); err != nil {
			log.Error("failed seeding dummy users", logger.Err(err))
			os.Exit(1)
		}
		log.Info("dummy users seeded successfully")
	}

	log.Info("migrator finished work successfully")
}
