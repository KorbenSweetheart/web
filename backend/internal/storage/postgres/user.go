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

	if err := s.db.WithContext(ctx).Create(account).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return domain.ErrEmailIsTaken
		}
		return fmt.Errorf("%s: failed to create account: %w", op, err)
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
			return nil, fmt.Errorf("%s: failed to get account, id %d: %w", op, id, err)
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
			return nil, fmt.Errorf("%s: failed to get account by email: %w", op, err)
		}
	}

	return &account, nil
}

// ProfileByID returns a single user profile record if it exists in the database.
func (s *Storage) ProfileByID(ctx context.Context, id int64) (*domain.Profile, error) {
	const op = "storage.postgres.ProfileByID"

	var profile domain.Profile

	err := s.db.WithContext(ctx).
		Preload("Activities.Activity").
		Where("user_id = ?", id).
		First(&profile).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		} else {
			return nil, fmt.Errorf("%s: failed to get profile by id: %d: %w", op, id, err)
		}
	}

	return &profile, nil
}

// ProfilesByIDs returns profile records matching given IDs with preloaded activities.
func (s *Storage) ProfilesByIDs(ctx context.Context, ids []int64) ([]*domain.Profile, error) {
	const op = "storage.postgres.ProfilesByIDs"

	var profiles []*domain.Profile
	err := s.db.WithContext(ctx).
		Preload("Activities.Activity").
		Where("user_id IN (?)", ids).
		Find(&profiles).Error

	if err != nil {
		return nil, fmt.Errorf("%s: failed to get profiles by ids: %w", op, err)
	}

	return profiles, nil
}

// UpdateProfileRecord updates profile db record based on the provided parameters.
func (s *Storage) UpdateProfileRecord(ctx context.Context, id int64, params *domain.ProfileUpdateParams) error {
	const op = "storage.postgres.UpdateProfileRecord"

	updates := make(map[string]any)

	if params.Name != nil {
		updates["name"] = *params.Name
	}
	if params.Age != nil {
		updates["age"] = *params.Age
	}
	if params.PictureURL != nil {
		updates["picture_url"] = *params.PictureURL
	}
	if params.Bio != nil {
		updates["bio"] = *params.Bio
	}
	if params.MaxRadius != nil {
		updates["max_radius"] = *params.MaxRadius
	}
	if params.InteractionMode != nil {
		updates["interaction_mode"] = int64(*params.InteractionMode)
	}
	if params.Lat != nil {
		updates["lat"] = *params.Lat
	}
	if params.Lon != nil {
		updates["lon"] = *params.Lon
	}

	// Scenario 1: Activities haven't been changed
	if params.Activities == nil {
		if len(updates) == 0 {
			return nil
		}
		err := s.db.WithContext(ctx).Model(&domain.Profile{}).Where("user_id = ?", id).Updates(updates).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return domain.ErrUserNotFound
			} else {
				return fmt.Errorf("%s: failed to update profile, id: %d: %w", op, id, err)
			}
		}
		return nil
	}

	// Scenario 2: Activities have been changed
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if len(updates) > 0 {
			if err := tx.Model(&domain.Profile{}).Where("user_id = ?", id).Updates(updates).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return domain.ErrUserNotFound
				} else {
					return fmt.Errorf("%s: failed to update profile, id: %d: %w", op, id, err)
				}
			}
		}

		if err := tx.Where("profile_user_id = ?", id).Delete(&domain.ProfileActivity{}).Error; err != nil {
			return fmt.Errorf("%s: failed to delete profile activities, id: %d: %w", op, id, err)
		}

		if len(*params.Activities) > 0 {

			newActivities := make([]domain.ProfileActivity, len(*params.Activities))
			for i, act := range *params.Activities {
				newActivities[i] = domain.ProfileActivity{
					ProfileUserID: id,
					ActivityID:    act.ActivityID,
					Experience:    act.Experience,
					InterestLevel: act.InterestLevel,
				}
			}
			if err := tx.Create(&newActivities).Error; err != nil {
				return fmt.Errorf("%s: failed to create new profile activities, id: %d: %w", op, id, err)
			}
		}

		return nil
	})
}
