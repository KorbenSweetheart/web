package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"match-me-api/internal/config"
	"match-me-api/internal/domain"
	"match-me-api/internal/logger"
	"match-me-api/internal/pkg/imgutil"
	"strings"
	"time"
)

const defaultPlaceholderObject = "default-profile.svg"

type UserRepository interface {
	AccountByEmail(ctx context.Context, email string) (*domain.Account, error)
	AccountByID(ctx context.Context, id int64) (*domain.Account, error)
	ProfileByID(ctx context.Context, id int64) (*domain.Profile, error)
	ProfilesByIDs(ctx context.Context, ids []int64) ([]*domain.Profile, error)
	UpdateProfileRecord(ctx context.Context, id int64, params *domain.ProfileUpdateParams) error
	FindConnectionRecord(ctx context.Context, fromUserID, toUserID int64) (*domain.Connection, error)
	IsDismissed(ctx context.Context, userID, targetUserID int64) (bool, error)
}

type MediaRepository interface {
	Upload(ctx context.Context, bucketName, objectName string, reader io.Reader, size int64, contentType string) (string, error)
	Delete(ctx context.Context, bucketName, objectName string) error
	ExtractObjectName(rawURL, bucketName string) string
	PublicURL(bucketName, objectName string) string
}

type UserService struct {
	repo  UserRepository
	media MediaRepository
	log   *slog.Logger
}

func NewUserService(repo UserRepository, media MediaRepository, logger *slog.Logger) *UserService {
	return &UserService{repo: repo, media: media, log: logger}
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

// Profile returns a user profile data struct from db if the user has permission to view it.
func (us *UserService) Profile(ctx context.Context, userID, targetUserID int64) (*domain.Profile, error) {
	const op = "service.userService.Profile"
	log := us.log.With(slog.String("op", op))

	// Self-profile view
	if userID == targetUserID {
		profile, err := us.repo.ProfileByID(ctx, targetUserID)
		if err != nil {
			log.Debug("failed to get profile by id", "targetUserID", targetUserID, "error", logger.Err(err))
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		if profile.PictureURL == "" {
			profile.PictureURL = us.defaultPictureURL()
		}
		return profile, nil
	}

	// 1. Connection-first check: if connected (or pending), view is permitted immediately
	conn, err := us.repo.FindConnectionRecord(ctx, userID, targetUserID)
	if err == nil && conn.Status != domain.Declined {
		profile, err := us.repo.ProfileByID(ctx, targetUserID)
		if err != nil {
			log.Debug("failed to get profile by id", "targetUserID", targetUserID, "error", logger.Err(err))
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		if profile.PictureURL == "" {
			profile.PictureURL = us.defaultPictureURL()
		}
		return profile, nil
	}

	if err != nil && !errors.Is(err, domain.ErrConnectionNotFound) {
		log.Debug("failed to find connection", "error", logger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// 2. Non-connected flow: fetch both profiles in single batch query
	profiles, err := us.repo.ProfilesByIDs(ctx, []int64{userID, targetUserID})
	if err != nil {
		log.Debug("failed to get profiles by ids", "error", logger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if len(profiles) < 2 {
		return nil, fmt.Errorf("%s: %w", op, domain.ErrUserNotFound)
	}

	var myProfile, targetProfile *domain.Profile
	if profiles[0].UserID == userID {
		myProfile, targetProfile = profiles[0], profiles[1]
	} else {
		myProfile, targetProfile = profiles[1], profiles[0]
	}

	// 3. In-memory candidate check (radii, activity, and score threshold)
	if !domain.IsCandidate(myProfile, targetProfile) {
		log.Debug("not a match candidate", "targetUserID", targetUserID)
		return nil, fmt.Errorf("%s: %w", op, domain.ErrNoPermissionViewProfile)
	}

	// 4. Check if recommendation was dismissed by either user
	isDismissed, err := us.repo.IsDismissed(ctx, userID, targetUserID)
	if err != nil {
		log.Debug("failed to check dismissed status", "error", logger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if isDismissed {
		log.Debug("recommendation dismissed", "targetUserID", targetUserID)
		return nil, fmt.Errorf("%s: %w", op, domain.ErrNoPermissionViewProfile)
	}

	if targetProfile.PictureURL == "" {
		targetProfile.PictureURL = us.defaultPictureURL()
	}

	return targetProfile, nil
}

// UpdateProfile updates profile with provided data.
func (us *UserService) UpdateProfile(ctx context.Context, id int64, params *domain.ProfileUpdateParams) error {
	const op = "service.userService.UpdateProfile"
	log := us.log.With(slog.String("op", op))

	// business rules
	if params.PictureURL != nil && *params.PictureURL == "" {
		defaultAvatar := us.defaultPictureURL()
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

// UpdateProfilePicture sanitizes, crops, and resizes image, uploads to MinIO, updates DB, and cleans old image.
func (us *UserService) UpdateProfilePicture(ctx context.Context, userID int64, rawFile io.Reader) (string, error) {
	const op = "service.userService.UpdateProfilePicture"
	log := us.log.With(slog.String("op", op))

	// Process & sanitize image (square crop, 400x400, clean JPEG 85%)
	processedReader, processedSize, contentType, err := imgutil.ProcessProfilePicture(
		rawFile,
		imgutil.DefaultTargetDim,
		imgutil.DefaultQuality,
	)
	if err != nil {
		return "", fmt.Errorf("%s: failed to process image: %w", op, err)
	}

	profile, err := us.repo.ProfileByID(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	oldPictureURL := profile.PictureURL

	// Upload sanitized JPEG to MinIO
	objectName := fmt.Sprintf("users/%d/picture-%d.jpg", userID, time.Now().Unix())
	newPictureURL, err := us.media.Upload(
		ctx,
		config.MinIOProfilePicturesBucket,
		objectName,
		processedReader,
		processedSize,
		contentType,
	)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	params := &domain.ProfileUpdateParams{
		PictureURL: &newPictureURL,
	}

	if err := us.repo.UpdateProfileRecord(ctx, userID, params); err != nil {
		log.Debug("failed to update profile picture url", "id", userID, "error", logger.Err(err))
		_ = us.media.Delete(ctx, config.MinIOProfilePicturesBucket, objectName)
		return "", fmt.Errorf("%s: %w", op, err)
	}

	// Clean up old image if it was a custom upload
	if us.isCustomPicture(oldPictureURL) {
		oldObjectName := us.media.ExtractObjectName(oldPictureURL, config.MinIOProfilePicturesBucket)
		if oldObjectName != "" {
			go func(obj string) {
				_ = us.media.Delete(context.Background(), config.MinIOProfilePicturesBucket, obj)
			}(oldObjectName)
		}
	}

	return newPictureURL, nil
}

// DeleteProfilePicture resets user's profile picture to default placeholder and removes custom image from MinIO.
func (us *UserService) DeleteProfilePicture(ctx context.Context, userID int64) (string, error) {
	const op = "service.userService.DeleteProfilePicture"
	log := us.log.With(slog.String("op", op))

	profile, err := us.repo.ProfileByID(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	defaultURL := us.defaultPictureURL()
	oldPictureURL := profile.PictureURL

	if oldPictureURL == defaultURL {
		return defaultURL, nil
	}

	params := &domain.ProfileUpdateParams{
		PictureURL: &defaultURL,
	}

	if err := us.repo.UpdateProfileRecord(ctx, userID, params); err != nil {
		log.Debug("failed to reset profile picture", "id", userID, "error", logger.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if us.isCustomPicture(oldPictureURL) {
		oldObjectName := us.media.ExtractObjectName(oldPictureURL, config.MinIOProfilePicturesBucket)
		if oldObjectName != "" {
			go func(obj string) {
				_ = us.media.Delete(context.Background(), config.MinIOProfilePicturesBucket, obj)
			}(oldObjectName)
		}
	}

	return defaultURL, nil
}

func (us *UserService) defaultPictureURL() string {
	if us.media != nil {
		return us.media.PublicURL(config.MinIOProfilePicturesBucket, defaultPlaceholderObject)
	}
	return ""
}

func (us *UserService) isCustomPicture(pictureURL string) bool {
	if pictureURL == "" || strings.HasSuffix(pictureURL, defaultPlaceholderObject) || strings.Contains(pictureURL, "placehold.net") {
		return false
	}
	return true
}
