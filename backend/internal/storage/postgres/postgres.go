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

// AutoMigrate
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

func (s *Storage) SeedDummyUsers(ctx context.Context) error {
	const op = "storage.postgres.SeedDummyUsers"

	// Seed users
	var count int64
	if err := s.db.WithContext(ctx).Model(&domain.Account{}).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to count existing accounts: op: %s, error: %w", op, err)
	}

	if count > 0 {
		return nil
	}

	users := generateDummyUsers()
	const batchLimit = 50

	if err := s.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "email"}},
		DoNothing: true,
	}).CreateInBatches(&users, batchLimit).Error; err != nil {
		return fmt.Errorf("failed to seed users, op: %s, error: %w", op, err)
	}

	return nil
}

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

	pingCtx, pingCtxCancel := context.WithTimeout(ctx, 2*time.Second)
	defer pingCtxCancel()
	if err := sqlDB.PingContext(pingCtx); err != nil {
		return nil, sqlDB, fmt.Errorf("failed to ping sql.DB: op: %s, error: %w", op, err)
	}

	// Adding PostGIS extension
	extCtx, extCancel := context.WithTimeout(ctx, 3*time.Second)
	defer extCancel()
	if err := db.WithContext(extCtx).Exec("CREATE EXTENSION IF NOT EXISTS postgis;").Error; err != nil {
		return nil, sqlDB, fmt.Errorf("failed to create postgis extension: op: %s, error: %w", op, err)
	}

	return db, sqlDB, nil
}

// generateSeedUsers creates fake users for testing purposes
func generateDummyUsers() []domain.Account {
	users := make([]domain.Account, 0, 100)

	user1 := domain.Account{
		Email:        "obiwan@matchme.com",
		PasswordHash: "$2a$12$15dw2.nyH6xOf10DjQezcOIDY.PL.Jkr6ZjJOjpmqcL3xHtVeTWIq", // 12345678
		Profile: domain.Profile{
			Name: "Obi-Wan Kenobi", // TODO: maybe add it during registration
			Age:  35,
			Bio: `A disciplined mind, a patient approach, and a good cup of tea are my essentials.
			I value loyalty, strategy, and staying calm in chaos. Always down for a witty debate or a long walk.`,
			MaxRadius:       10,
			InteractionMode: 2,
			Activities: []domain.ProfileActivity{
				{
					ActivityID:    1, // Running
					Experience:    3, // 1-5 levels: "Beginner", "Active Novice", "Intermediate", "Advanced", "Professional"
					InterestLevel: 5,
				},
				{
					ActivityID:    14,
					Experience:    5, // 1-5 levels: "Beginner", "Active Novice", "Intermediate", "Advanced", "Professional"
					InterestLevel: 4, // 1-5: "Not interested", "Open to it" , "Interested" , "Highly interested", "Actively looking"
				},
			},
			Lat: 60.1699,
			Lon: 24.9384,
		},
	}

	user2 := domain.Account{
		Email:        "anakin@matchme.com",
		PasswordHash: "$2a$12$15dw2.nyH6xOf10DjQezcOIDY.PL.Jkr6ZjJOjpmqcL3xHtVeTWIq", // 12345678
		Profile: domain.Profile{
			Name: "Anakin Skywalker", // TODO: maybe add it during registration
			Age:  19,
			Bio: `I live for speed, high stakes, and pushing limits.
			I trust my gut, speak my mind, and never back down from a challenge.
			If it's fast, intense, or "impossible", count me in.`,
			MaxRadius:       20,
			InteractionMode: 2,
			Activities: []domain.ProfileActivity{
				{
					ActivityID:    13, // Running
					Experience:    4,  // 1-5 levels: "Beginner", "Active Novice", "Intermediate", "Advanced", "Professional"
					InterestLevel: 4,
				},
				{
					ActivityID:    1,
					Experience:    3, // 1-5 levels: "Beginner", "Active Novice", "Intermediate", "Advanced", "Professional"
					InterestLevel: 5, // 1-5: "Not interested", "Open to it" , "Interested" , "Highly interested", "Actively looking"
				},
			},
			Lat: 60.1699,
			Lon: 24.9384,
		},
	}

	// TODO: add 100 randomly generated users

	users = append(users, user1, user2)

	return users
}
