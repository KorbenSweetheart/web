package handlers

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"match-me-api/internal/domain"
	"match-me-api/internal/logger"
	"match-me-api/internal/pkg/imgutil"
	"match-me-api/internal/transport/httpserver/dto"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

const maxProfilePictureSize = 1 << 20 // 1 MB

type UserManager interface {
	Profile(ctx context.Context, myID int64, targetUserID int64) (*domain.Profile, error)
	UpdateProfile(ctx context.Context, id int64, params *domain.ProfileUpdateParams) error
	UpdateProfilePicture(ctx context.Context, userID int64, rawFile io.Reader) (string, error)
	DeleteProfilePicture(ctx context.Context, userID int64) (string, error)
}

type UserHandler struct {
	userService UserManager
	validator   *validator.Validate
	log         *slog.Logger
}

func NewUserHandler(um UserManager, v *validator.Validate, logger *slog.Logger) *UserHandler {
	return &UserHandler{userService: um, validator: v, log: logger}
}

// UserSummary returns the user's id, name, and link to the profile picture.
// /users/{id}
func (h *UserHandler) UserSummary(c *echo.Context) error {
	ctx := c.Request().Context()

	myID, ok := c.Get("user_id").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]any{"error": "Unauthorized"})
	}

	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "Invalid user id: NAN"})
	}

	profile, err := h.userService.Profile(ctx, myID, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			return c.JSON(http.StatusNotFound, map[string]any{
				"id":    userID,
				"error": "User not found",
			})
		case errors.Is(err, domain.ErrNoPermissionViewProfile):
			return c.JSON(http.StatusForbidden, map[string]any{
				"error": "No permission to view profile",
			})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]any{"error": "Failed to get user"})
		}
	}

	return c.JSON(http.StatusOK, dto.UserSummaryResponse{
		ID:         profile.UserID,
		Name:       profile.Name,
		PictureURL: profile.PictureURL,
	})
}

// // UserProfile returns the user's id and "about me" type information.
// // /users/{id}/profile
func (h *UserHandler) UserProfile(c *echo.Context) error {
	ctx := c.Request().Context()

	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "Invalid user id: NAN"})
	}

	myID, ok := c.Get("user_id").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]any{"error": "Unauthorized"})
	}

	profile, err := h.userService.Profile(ctx, myID, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			return c.JSON(http.StatusNotFound, map[string]any{
				"id":    userID,
				"error": "User not found",
			})
		case errors.Is(err, domain.ErrNoPermissionViewProfile):
			return c.JSON(http.StatusForbidden, map[string]any{
				"error": "No permission to view profile",
			})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]any{"error": "Failed to get user"})
		}
	}

	return c.JSON(http.StatusOK, dto.ProfileResponse{
		ID:  profile.UserID,
		Age: profile.Age,
		Bio: profile.Bio,
		// IsOnline: profile.IsOnline,
	})
}

// UserBio returns the user's id and biographical data (the data used to power recommendations).
// /users/{id}/bio
func (h *UserHandler) UserBio(c *echo.Context) error {
	ctx := c.Request().Context()

	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "Invalid user id: NAN"})
	}

	myID, ok := c.Get("user_id").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]any{"error": "Unauthorized"})
	}

	profile, err := h.userService.Profile(ctx, myID, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			return c.JSON(http.StatusNotFound, map[string]any{
				"id":    userID,
				"error": "User not found",
			})
		case errors.Is(err, domain.ErrNoPermissionViewProfile):
			return c.JSON(http.StatusForbidden, map[string]any{
				"error": "No permission to view profile",
			})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]any{"error": "Failed to get user"})
		}
	}

	activitiesResponse := make([]dto.UserActivityResponse, 0, len(profile.Activities))

	if len(profile.Activities) > 0 {
		for _, a := range profile.Activities {
			activity := dto.UserActivityResponse{
				ID:            a.Activity.ID,
				Title:         a.Activity.Title,
				Experience:    int(a.Experience),
				InterestLevel: int(a.InterestLevel),
			}
			activitiesResponse = append(activitiesResponse, activity)
		}
	}

	return c.JSON(http.StatusOK, dto.UserBioResponse{
		ID:              profile.UserID,
		MaxRadius:       profile.MaxRadius,
		InteractionMode: int(profile.InteractionMode),
		Activities:      activitiesResponse,
		Lat:             profile.Lat,
		Lon:             profile.Lon,
	})
}

// MySummary returns the user's id, name, and link to the profile picture for the authorized user.
// /me
func (h *UserHandler) MySummary(c *echo.Context) error {
	ctx := c.Request().Context()

	myID, ok := c.Get("user_id").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]any{"error": "Unauthorized"})
	}

	profile, err := h.userService.Profile(ctx, myID, myID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, map[string]any{
				"id":    myID,
				"error": "User not found",
			})
		} else {
			return c.JSON(http.StatusInternalServerError, map[string]any{"error": "Failed to get user"})
		}
	}

	return c.JSON(http.StatusOK, dto.MySummaryResponse{
		ID:         profile.UserID,
		Email:      profile.Email,
		Name:       profile.Name,
		PictureURL: profile.PictureURL,
	})
}

// MyProfile returns the user's id and "about me" type information for the authorized user.
// /me/profile
func (h *UserHandler) MyProfile(c *echo.Context) error {
	ctx := c.Request().Context()

	myID, ok := c.Get("user_id").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]any{"error": "Unauthorized"})
	}

	profile, err := h.userService.Profile(ctx, myID, myID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, map[string]any{
				"id":    myID,
				"error": "User not found",
			})
		} else {
			return c.JSON(http.StatusInternalServerError, map[string]any{"error": "Failed to get user"})
		}
	}

	return c.JSON(http.StatusOK, dto.ProfileResponse{
		ID:  profile.UserID,
		Age: profile.Age,
		Bio: profile.Bio,
	})
}

// MyBio returns the user's id and biographical data (the data used to power recommendations) for the authorized user.
// /me/bio
func (h *UserHandler) MyBio(c *echo.Context) error {
	ctx := c.Request().Context()

	myID, ok := c.Get("user_id").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]any{"error": "Unauthorized"})
	}

	profile, err := h.userService.Profile(ctx, myID, myID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, map[string]any{
				"id":    myID,
				"error": "User not found",
			})
		} else {
			return c.JSON(http.StatusInternalServerError, map[string]any{"error": "Failed to get user"})
		}
	}

	activitiesResponce := []dto.UserActivityResponse{}

	if len(profile.Activities) > 0 {
		for _, a := range profile.Activities {
			activity := dto.UserActivityResponse{
				ID:            a.Activity.ID,
				Title:         a.Activity.Title,
				Experience:    int(a.Experience),
				InterestLevel: int(a.InterestLevel),
			}
			activitiesResponce = append(activitiesResponce, activity)
		}
	}

	return c.JSON(http.StatusOK, dto.UserBioResponse{
		ID:              profile.UserID,
		MaxRadius:       profile.MaxRadius,
		InteractionMode: int(profile.InteractionMode),
		Activities:      activitiesResponce,
		// TODO: remove lat and lon from responce, added for testing
		Lat: profile.Lat,
		Lon: profile.Lon,
	})
}

// UpdateProfile updates existin user profile fully or partially based on the provided data.
func (h *UserHandler) UpdateProfile(c *echo.Context) error {
	const op = "httpserver.handlers.UpdateProfile"
	// log := h.log.With(slog.String("op", op))

	ctx := c.Request().Context()

	userID, ok := c.Get("user_id").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]any{"error": "Unauthorized"})
	}

	var req dto.ProfileUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "Invalid request body"})
	}
	if err := h.validator.Struct(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
	}

	// 3. Map DTO into domain struct
	params := &domain.ProfileUpdateParams{
		Name:       req.Name,
		Age:        req.Age,
		PictureURL: req.PictureURL,
		Bio:        req.Bio,
		MaxRadius:  req.MaxRadius,
		Lat:        req.Lat,
		Lon:        req.Lon,
	}

	if req.InteractionMode != nil {
		params.InteractionMode = new(domain.InteractionMode(*req.InteractionMode))
	}

	if req.Activities != nil {
		activities := make([]domain.ActivityInput, len(*req.Activities))
		for i, act := range *req.Activities {
			activities[i] = domain.ActivityInput{
				ActivityID:    act.ID,
				Experience:    domain.ExperienceLevel(act.Experience),
				InterestLevel: domain.InterestLevel(act.InterestLevel),
			}
		}
		params.Activities = &activities
	}

	if err := h.userService.UpdateProfile(ctx, userID, params); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{"error": "Failed to update profile"})
	}

	return c.JSON(http.StatusOK, map[string]any{"status": "Profile updated successfully"})
}

// UploadProfilePicture uploads, processes, and updates the profile picture for the authorized user.
// POST /me/picture
func (h *UserHandler) UploadProfilePicture(c *echo.Context) error {
	const op = "httpserver.handlers.UploadProfilePicture"
	ctx := c.Request().Context()

	userID, ok := c.Get("user_id").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]any{"error": "Unauthorized"})
	}

	fileHeader, err := c.FormFile("picture")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "Form field 'picture' is required"})
	}

	if fileHeader.Size > maxProfilePictureSize {
		return c.JSON(http.StatusRequestEntityTooLarge, map[string]any{"error": "File size exceeds 1MB limit"})
	}

	file, err := fileHeader.Open()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "Failed to read uploaded file"})
	}
	defer file.Close()

	pictureURL, err := h.userService.UpdateProfilePicture(ctx, userID, file)
	if err != nil {
		if errors.Is(err, imgutil.ErrUnsupportedFormat) || errors.Is(err, imgutil.ErrInvalidDimensions) {
			return c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
		}
		h.log.Error("failed to update profile picture", slog.String("op", op), "user_id", userID, logger.Err(err))
		return c.JSON(http.StatusInternalServerError, map[string]any{"error": "Failed to update profile picture"})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"status":      "Profile picture updated successfully",
		"picture_url": pictureURL,
	})
}

// DeleteProfilePicture resets the profile picture for the authorized user to the default placeholder.
// DELETE /me/picture
func (h *UserHandler) DeleteProfilePicture(c *echo.Context) error {
	const op = "httpserver.handlers.DeleteProfilePicture"
	ctx := c.Request().Context()

	userID, ok := c.Get("user_id").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]any{"error": "Unauthorized"})
	}

	defaultURL, err := h.userService.DeleteProfilePicture(ctx, userID)
	if err != nil {
		h.log.Error("failed to delete profile picture", slog.String("op", op), "user_id", userID, logger.Err(err))
		return c.JSON(http.StatusInternalServerError, map[string]any{"error": "Failed to reset profile picture"})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"status":      "Profile picture removed successfully",
		"picture_url": defaultURL,
	})
}
