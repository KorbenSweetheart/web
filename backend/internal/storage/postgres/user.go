package postgres

import (
	"context"
	"errors"
	"fmt"
	"match-me-api/internal/domain"

	"gorm.io/gorm"
)

// CreateAccount adds a single user account record to database.
func (s *Storage) CreateAccount(ctx context.Context, account *domain.Account) error {
	const op = "storage.postgres.CreateAccount"
	// log := s.log.With(slog.String("op", op))

	if err := s.db.WithContext(ctx).Create(account).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return domain.ErrEmailIsTaken
		}
		return fmt.Errorf("failed to create account, op: %s, error: %w", op, err)
	}

	return nil
}

func (s *Storage) SaveRefreshToken(ctx context.Context, rtRecord *domain.RefreshToken) error {
	const op = "storage.postgres.SaveRefreshToken"
	// log := s.log.With(slog.String("op", op))

	if err := s.db.WithContext(ctx).Create(rtRecord).Error; err != nil {
		return fmt.Errorf("failed to create refresh token, op: %s, error: %w", op, err)
	}

	return nil
}

// AccountByID returns a single account record if it exists in the database.
func (s *Storage) AccountByID(ctx context.Context, id int64) (*domain.Account, error) {
	const op = "storage.postgres.AccountByID"

	var account domain.Account

	err := s.db.WithContext(ctx).
		Where("id = ?", id).
		First(&account).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		} else {
			return nil, fmt.Errorf("failed to get account by id, op: %s, error: %w", op, err)
		}
	}

	return &account, nil
}

// AccountByEmail returns a single account record if it exists in the database.
func (s *Storage) AccountByEmail(ctx context.Context, email string) (*domain.Account, error) {
	const op = "storage.postgres.AccountByEmail"

	var account domain.Account

	err := s.db.WithContext(ctx).
		Where("email = ?", email).
		First(&account).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		} else {
			return nil, fmt.Errorf("failed to get account by email, op: %s, error: %w", op, err)
		}
	}

	return &account, nil
}

// ProfileByID returns a single user profile record if it exists in the database.
func (s *Storage) ProfileByID(ctx context.Context, id int64) (*domain.Profile, error) {
	const op = "storage.postgres.ProfileByID"

	var profile domain.Profile

	err := s.db.WithContext(ctx).
		Preload("InteractionMode").
		Preload("Activities.Activity").
		Where("user_id = ?", id).
		First(&profile).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		} else {
			return nil, fmt.Errorf("failed to get profile by id, op: %s, error: %w", op, err)
		}
	}

	return &profile, nil
}

// func (s *Storage) UpdateProfile(ctx context.Context, profile *domain.Profile) error

// IsEmailTaken checks whether the email is already taken.
// Currently not used anywhere
func (s *Storage) IsEmailTaken(ctx context.Context, email string) (bool, error) {
	const op = "storage.postgres.IsEmailTaken"

	var count int64

	err := s.db.Model(&domain.Account{}).
		Where("email = ?", email).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("failed to check email existence, op: %s, error: %w", op, err)
	}

	return count > 0, nil
}
