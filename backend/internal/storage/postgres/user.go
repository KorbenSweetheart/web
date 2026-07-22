package postgres

import (
	"context"
	"errors"
	"fmt"
	"match-me-api/internal/domain"
	"match-me-api/internal/storage"

	"gorm.io/gorm"
)

// CreateUser adds a single user record to database.
func (s *Storage) CreateUser(ctx context.Context, user *domain.User) error {

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Create(user).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return storage.ErrEmailIsTaken
			}
			return fmt.Errorf("failed to create user record: %w", err)
		}

		profile := &domain.Profile{
			UserID: user.ID,
			// Name:   "",
			// MaxRadius:         10,
			// InteractionModeID: 4,
		}

		if err := tx.WithContext(ctx).Create(profile).Error; err != nil {
			return fmt.Errorf("failed to create profile record: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to execute the user creation transaction: %w", err)
	}

	return nil
}

// UserByID returns a single user record if it exists in the database.
func (s *Storage) UserByID(ctx context.Context, id int64) (*domain.User, error) {

	var user domain.User

	err := s.db.WithContext(ctx).
		Where("id = ?", id).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, storage.ErrUserNotFound
		} else {
			return nil, fmt.Errorf("failed to get user by id: %w", err)
		}
	}

	return &user, nil
}

// UserByEmail returns a single user record if it exists in the database.
func (s *Storage) UserByEmail(ctx context.Context, email string) (*domain.User, error) {

	var user domain.User

	err := s.db.WithContext(ctx).
		Where("email = ?", email).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, storage.ErrUserNotFound
		} else {
			return nil, fmt.Errorf("failed to get user by email: %w", err)
		}
	}

	return &user, nil
}

// ProfileByID returns a single user profile record if it exists in the database.
func (s *Storage) ProfileByID(ctx context.Context, id int64) (*domain.Profile, error) {
	var profile domain.Profile

	err := s.db.WithContext(ctx).
		Where("user_id = ?", id).
		First(&profile).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, storage.ErrUserNotFound
		} else {
			return nil, fmt.Errorf("failed to get profile by id: %w", err)
		}
	}

	return &profile, nil
}

// func (s *Storage) UpdateProfile(ctx context.Context, profile *domain.Profile) error

// IsEmailTaken checks whether the email is already taken.
func (s *Storage) IsEmailTaken(ctx context.Context, email string) (bool, error) {

	var count int64

	err := s.db.Model(&domain.User{}).
		Where("email = ?", email).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return count > 0, nil
}
