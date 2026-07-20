package postgres

import (
	"fmt"
	"log/slog"
	"match-me-api/internal/config"
	"match-me-api/internal/domain"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Storage struct {
	db  *gorm.DB
	log *slog.Logger
}

// New creates a new Postgres storage object.
func NewPostgresDB(dbCfg config.Database, log *slog.Logger) (*Storage, error) {
	const op = "storage.postgres.New"

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Helsinki",
		dbCfg.Host, dbCfg.User, dbCfg.Pass, dbCfg.Name, dbCfg.Port)

	// Note: AutomaticPing: true
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %s, %w", op, err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB from gorm: %s, %w", op, err)
	}

	// Configuring DB
	// Connection Pool
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(10 * time.Minute)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping postgres db: %s, %w", op, err)
	}

	return &Storage{
		db:  db,
		log: log,
	}, nil
}

func (s *Storage) AutoMigrate() error {
	const op = "storage.postgres.AutoMigrate"

	if err := s.db.AutoMigrate(&domain.User{}, &domain.Profile{}); err != nil {
		return fmt.Errorf("failed to migrate structs using gorm: %s, %w", op, err)
	}

	return nil
}
