package service

import (
	"context"
	"fmt"
	"log/slog"
	"match-me-api/internal/domain"
	"match-me-api/internal/logger"
)

type UserRepository interface {
	AccountByEmail(ctx context.Context, email string) (*domain.Account, error)
	AccountByID(ctx context.Context, id int64) (*domain.Account, error)
	ProfileByID(ctx context.Context, id int64) (*domain.Profile, error)
	UpdateProfileRecord(ctx context.Context, id int64, params *domain.ProfileUpdateParams) error
	// IsEmailTaken(ctx context.Context, email string) (bool, error)
}

type UserService struct {
	repo UserRepository
	log  *slog.Logger
}

func NewUserService(r UserRepository, logger *slog.Logger) *UserService {
	return &UserService{repo: r, log: logger}
}

// Account returns a user account data struct from db.
func (us *UserService) Account(ctx context.Context, id int64) (*domain.Account, error) {
	const op = "service.userService.Account"
	log := us.log.With(slog.String("op", op))

	user, err := us.repo.AccountByID(ctx, id)
	if err != nil {
		log.Debug("failed to get account by id", "id", id, "error", logger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

// Profile returns a user profile data struct from db.
func (us *UserService) Profile(ctx context.Context, id int64) (*domain.Profile, error) {
	const op = "service.userService.Profile"
	log := us.log.With(slog.String("op", op))

	profile, err := us.repo.ProfileByID(ctx, id)
	if err != nil {
		log.Debug("failed to get profile by id", "id", id, "error", logger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return profile, nil
}

// UpdateProfile updates profile with provided data.
func (us *UserService) UpdateProfile(ctx context.Context, id int64, params *domain.ProfileUpdateParams) error {
	const op = "service.userService.Profile"
	log := us.log.With(slog.String("op", op))

	// business rules
	if params.PictureURL != nil && *params.PictureURL == "" {
		defaultAvatar := DefaultAvatar
		params.PictureURL = &defaultAvatar
	}

	// TODO: Maybe add MaxRadius rules?
	if params.MaxRadius != nil && *params.MaxRadius <= 0 {
		defaultRadius := DefaultRadius
		params.MaxRadius = &defaultRadius
	}

	if err := us.repo.UpdateProfileRecord(ctx, id, params); err != nil {
		log.Debug("failed to update profile", "id", id, "error", logger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// isCompleted verifies Profile bio parameters and returns "true" if all touch points are set.
// After that the profile can get recommendations and be used for recommendations.
// func isCompleted(p *domain.Profile) bool {
// 	const op = "service.userService.isCompleted"
// 	// log := us.log.With(slog.String("op", op))

// 	if p.Lat == 0 && p.Lon == 0 {
// 		return false
// 	}
// 	if len(p.Activities) == 0 {
// 		return false
// 	}
// 	return true
// }
