package service

import (
	"context"
	"log/slog"
	"match-me-api/internal/domain"
	"match-me-api/internal/logger"
)

type UserProvider interface {
	AccountByEmail(ctx context.Context, email string) (*domain.Account, error)
	IsEmailTaken(ctx context.Context, email string) (bool, error)
	AccountByID(ctx context.Context, id int64) (*domain.Account, error)
	ProfileByID(ctx context.Context, id int64) (*domain.Profile, error)
	// UpdateProfile(ctx context.Context, profile *domain.Profile) error
}

type UserService struct {
	storage UserProvider
	log     *slog.Logger
}

func NewUserService(up UserProvider, logger *slog.Logger) *UserService {
	return &UserService{storage: up, log: logger}
}

// Account returns a user account data struct from db.
func (us *UserService) Account(ctx context.Context, id int64) (*domain.Account, error) {
	const op = "service.userService.Account"
	log := us.log.With(slog.String("op", op))

	user, err := us.storage.AccountByID(ctx, id)
	if err != nil {
		log.Debug("failed to get account by id", "id", id, "error", logger.Err(err))
		return nil, err
	}

	return user, nil
}

// Profile returns a user profile data struct from db.
func (us *UserService) Profile(ctx context.Context, id int64) (*domain.Profile, error) {
	const op = "service.userService.Profile"
	log := us.log.With(slog.String("op", op))

	profile, err := us.storage.ProfileByID(ctx, id)
	if err != nil {
		log.Debug("failed to get profile by id", "id", id, "error", logger.Err(err))
		return nil, err
	}

	return profile, nil
}

// // UpdateProfile updates profile with provided data and returns it back with changes.
// func (us *UserService) UpdateProfile(ctx context.Context, profile domain.Profile) (*domain.Profile, error) {
// 	const op = "service.userService.Profile"
// 	log := us.log.With(slog.String("op", op))

// 	profile, err := us.storage.ProfileByID(ctx, id)
// 	if err != nil {
// 		log.Debug("failed to get profile by id", "id", id, "error", logger.Err(err))
// 		return nil, err
// 	}

// 	return profile, nil
// }

// func (us *UserService) UpdateLocation(ctx context.Context, id int64, lat, lon float64) error {
// 	const op = "service.userService.Profile"
// 	log := us.log.With(slog.String("op", op))

// 	profile, err := us.storage.ProfileByID(ctx, id)
// 	if err != nil {
// 		log.Debug("failed to get profile by id", "id", id, "error", logger.Err(err))
// 		return err
// 	}

// 	profile.Lat = lat
// 	profile.Lon = lon

// 	return nil
// }
