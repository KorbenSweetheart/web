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
	UpdateProfileRecord(ctx context.Context, id int64, params *domain.ProfileUpdateParams) error
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

// UpdateProfile updates profile with provided data and returns it back with changes.
func (us *UserService) UpdateProfile(ctx context.Context, id int64, params *domain.ProfileUpdateParams) error {
	const op = "service.userService.Profile"
	log := us.log.With(slog.String("op", op))

	// Бизнес-rules
	// if params.PictureURL != nil && *params.PictureURL == "" {
	// 	defaultAvatar := "https://your-cdn.com/avatars/default_placeholder.png"
	// 	params.PictureURL = &defaultAvatar
	// }

	// При необходимости здесь можно добавить дополнительные проверки доступа и валидацию бизнес-правил

	if err := us.storage.UpdateProfileRecord(ctx, id, params); err != nil {
		log.Debug("failed to get profile by id", "id", id, "error", logger.Err(err))
		return err
	}

	return nil
}
