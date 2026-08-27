package main

import (
	"context"
	"log/slog"
	"match-me-api/internal/config"
	"match-me-api/internal/logger"
	"match-me-api/internal/storage/minio"
	"match-me-api/internal/storage/postgres"
	"os"
	"os/signal"
	"strings"
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

	// Seeding dummy data and users
	if (cfg.Env == "dev" || cfg.Env == "development") && cfg.SeedUsers {
		log.Info("seeding test dummy users...")
		dummyCtx, dummyCancel := context.WithTimeout(shutdownCtx, 15*time.Second)
		defer dummyCancel()

		if err := storage.SeedDummyUsers(dummyCtx); err != nil {
			log.Error("failed seeding dummy users", logger.Err(err))
			os.Exit(1)
		}
		log.Info("dummy users seeded successfully")
	} else {
		log.Info("skipping test dummy users seeding (SEED_USERS is disabled)")
	}

	// Initializing S3-buckets in MinIO
	bucketCtx, bucketCancel := context.WithTimeout(shutdownCtx, 10*time.Second)
	defer bucketCancel()

	minioStore, err := minio.NewMinioStorage(bucketCtx, cfg.MinIO, log)
	if err != nil {
		log.Error("failed to connect to minio in migrator", logger.Err(err))
		os.Exit(1)
	}

	// List of buckets
	requiredBuckets := []string{
		config.MinIOProfilePicturesBucket,
		// config.MinIOChatMediaBucket, // "chat-attachments",
	}

	for _, bucket := range requiredBuckets {
		policy, err := minio.MakeReadOnlyPolicy(bucket)
		if err != nil {
			log.Error("failed to generate bucket policy", slog.String("bucket", bucket), logger.Err(err))
			os.Exit(1)
		}
		if err := minioStore.InitBucket(bucketCtx, bucket, policy); err != nil {
			log.Error("failed to init bucket", slog.String("bucket", bucket), logger.Err(err))
			os.Exit(1)
		}
	}

	// Seed default profile picture placeholder if it doesn't exist
	placeholderName := "default-profile.svg"
	exists, err := minioStore.ObjectExists(bucketCtx, config.MinIOProfilePicturesBucket, placeholderName)
	if err != nil {
		log.Warn("could not check default profile placeholder existence", logger.Err(err))
	} else if !exists {
		defaultSVG := `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 200 200" width="100%" height="100%"><rect width="200" height="200" fill="#E2E8F0"/><circle cx="100" cy="75" r="35" fill="#94A3B8"/><path d="M45,170 C45,130 70,120 100,120 C130,120 155,130 155,170 Z" fill="#94A3B8"/></svg>`
		svgReader := strings.NewReader(defaultSVG)
		if _, err := minioStore.Upload(bucketCtx, config.MinIOProfilePicturesBucket, placeholderName, svgReader, int64(len(defaultSVG)), "image/svg+xml"); err != nil {
			log.Warn("failed to seed default profile placeholder", logger.Err(err))
		} else {
			log.Info("default profile placeholder seeded successfully", slog.String("object", placeholderName))
		}
	}

	log.Info("all migrations and bucket initializations completed successfully")
}
