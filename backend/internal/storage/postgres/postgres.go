package postgres

import (
	"fmt"
	"log/slog"
	"match-me-api/internal/config"
	"match-me-api/internal/domain"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Storage struct {
	db  *gorm.DB
	log *slog.Logger
}

// NewPostgresDB creates a new Postgres storage object.
func NewPostgresDB(dbCfg config.Database, log *slog.Logger) (*Storage, error) {
	const op = "storage.postgres.New"

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Helsinki",
		dbCfg.Host, dbCfg.User, dbCfg.Pass, dbCfg.Name, dbCfg.Port)

	// Wrap our logger for GORM
	// https://gorm.io/docs/logger.html
	gormLogger := logger.New(
		slog.NewLogLogger(log.Handler(), slog.LevelInfo), // info, just to simplify things
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Info, // Log slow queries & errors
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	// Note: AutomaticPing: true
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: gormLogger})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %s, %w", op, err)
	}

	// For additional DB configuration
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

// AutoMigrate
func (s *Storage) AutoMigrate() error {
	const op = "storage.postgres.AutoMigrate"

	log := s.log.With(
		slog.String("op", op),
	)

	if err := s.db.AutoMigrate(
		&domain.InteractionMode{},
		&domain.Activity{},
		&domain.Account{},
		&domain.Profile{},
		&domain.RefreshToken{},
		&domain.ProfileActivity{},
		&domain.Connection{},
		&domain.Chat{},
		&domain.Message{},
	); err != nil {
		return fmt.Errorf("failed to auto-migrate: %s, %w", op, err)
	}

	log.Info("database auto-migration completed successfully")

	return nil
}

// SeedData adds dictionary elements and default values to the tables
func (s *Storage) SeedData() error {
	const op = "storage.postgres.SeedData"

	// Interaction Modes
	modes := []domain.InteractionMode{
		{ID: 1, Title: "Silent"},
		{ID: 2, Title: "Social"},
		{ID: 3, Title: "Dating"},
		{ID: 4, Title: "Open to anything"},
	}

	for _, mode := range modes {
		if err := s.db.FirstOrCreate(&mode, domain.InteractionMode{ID: mode.ID}).Error; err != nil {
			return fmt.Errorf("failed to seed interaction mode: %d:, op: %s, error: %w", mode.ID, op, err)
		}
	}

	// Activities
	activities := []domain.Activity{
		{ID: 1, Title: "Running"},
		{ID: 2, Title: "Padel"},
		{ID: 3, Title: "Gym"},
		{ID: 4, Title: "Cycling"},
		{ID: 4, Title: "Football"},
		{ID: 6, Title: "Tennis"},
		{ID: 7, Title: "Swimming"},
		{ID: 8, Title: "CrossFit"},
		{ID: 9, Title: "Yoga"},
		{ID: 10, Title: "Basketball"},
		{ID: 11, Title: "Climbing"},
		{ID: 12, Title: "Boxing"},
	}

	for _, act := range activities {
		if err := s.db.FirstOrCreate(&act, domain.Activity{ID: act.ID}).Error; err != nil {
			return fmt.Errorf("failed to seed activity, op: %s, id: %d, error: %w", op, act.ID, err)
		}
	}

	return nil
}
