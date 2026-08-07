package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"match-me-api/internal/config"
	"match-me-api/internal/domain"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

const (
	maxRetries        = 5
	backoffMultiplier = 2
)

type Storage struct {
	db  *gorm.DB
	log *slog.Logger
}

// NewPostgresDB creates a new Postgres storage object.
func NewPostgresDB(ctx context.Context, dbCfg config.Database, log *slog.Logger) (*Storage, error) {
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

	var (
		db      *gorm.DB
		sqlDB   *sql.DB
		err     error
		backoff = 500 * time.Millisecond
	)

	for attempt := 1; attempt <= maxRetries; attempt++ {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("failed to connect, context canceled before attempt: op: %s, error: %w", op, ctx.Err())
		default:
		}

		db, sqlDB, err = connectAndSetup(ctx, dsn, gormLogger)
		if err == nil {
			// Configuring DB
			// Connection Pool
			sqlDB.SetMaxOpenConns(25)
			sqlDB.SetMaxIdleConns(5)
			sqlDB.SetConnMaxLifetime(10 * time.Minute)
			sqlDB.SetConnMaxIdleTime(5 * time.Minute)

			return &Storage{db: db, log: log}, nil
		}

		if sqlDB != nil {
			_ = sqlDB.Close()
		}

		if attempt == maxRetries {
			break
		}

		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("failed to connect, context canceled during backoff: op: %s, error: %w", op, ctx.Err())
		case <-time.After(backoff):
		}

		backoff *= backoffMultiplier
	}

	return nil, fmt.Errorf("failed to connect to db after %d attempts: op: %s, error: %w", maxRetries, op, err)
}

// AutoMigrate creates required tables in database
func (s *Storage) AutoMigrate(ctx context.Context) error {
	const op = "storage.postgres.AutoMigrate"

	log := s.log.With(
		slog.String("op", op),
	)

	if err := s.db.WithContext(ctx).AutoMigrate(
		&domain.Activity{},
		&domain.Account{},
		&domain.Profile{},
		&domain.RefreshToken{},
		&domain.ProfileActivity{},
		&domain.Connection{},
		&domain.Chat{},
		&domain.Message{},
	); err != nil {
		return fmt.Errorf("failed to auto-migrate: op: %s, error: %w", op, err)
	}

	locationColQuery := `
        ALTER TABLE profiles 
        ADD COLUMN IF NOT EXISTS location geography(Point, 4326) 
        GENERATED ALWAYS AS (ST_SetSRID(ST_MakePoint(lon, lat), 4326)::geography) STORED;
    `
	if err := s.db.WithContext(ctx).Exec(locationColQuery).Error; err != nil {
		return fmt.Errorf("failed to add generated location column: op: %s, error: %w", op, err)
	}

	// GIST index for radius search
	gistIndexQuery := `CREATE INDEX IF NOT EXISTS idx_profiles_location ON profiles USING GIST (location);`
	if err := s.db.WithContext(ctx).Exec(gistIndexQuery).Error; err != nil {
		return fmt.Errorf("failed to create GIST index: op: %s, error: %w", op, err)
	}

	log.Info("database auto-migration completed successfully")

	return nil
}

// SeedData adds dictionary elements and default values to the tables
func (s *Storage) SeedDictionaries(ctx context.Context) error {
	const op = "storage.postgres.SeedDictionaries"

	// Activities
	activities := []domain.Activity{
		{ID: 1, Title: "Running"},
		{ID: 2, Title: "Padel"},
		{ID: 3, Title: "Gym"},
		{ID: 4, Title: "Cycling"},
		{ID: 5, Title: "Football"},
		{ID: 6, Title: "Tennis"},
		{ID: 7, Title: "Swimming"},
		{ID: 8, Title: "CrossFit"},
		{ID: 9, Title: "Yoga"},
		{ID: 10, Title: "Basketball"},
		{ID: 11, Title: "Climbing"},
		{ID: 12, Title: "Boxing"},
		{ID: 13, Title: "MMA"},
		{ID: 14, Title: "Aikido"},
		{ID: 15, Title: "Jiu-Jitsu"},
	}

	if err := s.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoNothing: true,
	}).Create(&activities).Error; err != nil {
		return fmt.Errorf("failed to seed activities: op: %s, error: %w", op, err)
	}

	// Reset Postgres serial sequences so dynamic INSERTs don't crash on primary key conflicts
	table := "activities"
	query := fmt.Sprintf(
		"SELECT setval(pg_get_serial_sequence('%s', 'id'), COALESCE(MAX(id), 1)) FROM %s;",
		table, table,
	)
	if err := s.db.WithContext(ctx).Exec(query).Error; err != nil {
		return fmt.Errorf("failed to execute reset sequences for table %s: op: %s, error: %w", table, op, err)
	}

	return nil
}

func (s *Storage) Ping(ctx context.Context) error {
	const op = "storage.postgres.Ping"

	sqlDB, err := s.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: op: %s, error: %w", op, err)
	}

	pingCtx, pingCtxCancel := context.WithTimeout(ctx, 2*time.Second)
	defer pingCtxCancel()

	if err := sqlDB.PingContext(pingCtx); err != nil {
		return fmt.Errorf("failed to ping sql.DB: op: %s, error: %w", op, err)
	}

	return nil
}

// connectAndSetup is a helper "facade" function that wraps the db setup, including context and timeouts, to improve readability.
func connectAndSetup(ctx context.Context, dsn string, gormLogger logger.Interface) (*gorm.DB, *sql.DB, error) {
	const op = "storage.postgres.connectAndSetup"

	// Note: AutomaticPing: true
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: gormLogger})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to postgres: op: %s, error: %w", op, err)
	}

	// For additional DB configuration
	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get sql.DB: op: %s, error: %w", op, err)
	}

	pingCtx, pingCtxCancel := context.WithTimeout(ctx, 5*time.Second)
	defer pingCtxCancel()
	if err := sqlDB.PingContext(pingCtx); err != nil {
		return nil, sqlDB, fmt.Errorf("failed to ping sql.DB: op: %s, error: %w", op, err)
	}

	// Adding PostGIS extension
	extCtx, extCancel := context.WithTimeout(ctx, 5*time.Second)
	defer extCancel()
	if err := db.WithContext(extCtx).Exec("CREATE EXTENSION IF NOT EXISTS postgis;").Error; err != nil {
		return nil, sqlDB, fmt.Errorf("failed to create postgis extension: op: %s, error: %w", op, err)
	}

	return db, sqlDB, nil
}
