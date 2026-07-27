package handlers

import (
	"context"
	"errors"
	"log/slog"
	"match-me-api/internal/domain"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

type UserService interface {
	Account(ctx context.Context, id int64) (*domain.Account, error)
	Profile(ctx context.Context, id int64) (*domain.Profile, error)
	// UpdateProfile(ctx context.Context, profile *domain.Profile) error
	// AccountByEmail(ctx context.Context, email string) (*domain.Account, error)
}

type UserHandler struct {
	userService UserService
	validator   *validator.Validate
	log         *slog.Logger
}

func NewUserHandler(us UserService, v *validator.Validate, logger *slog.Logger) *UserHandler {
	return &UserHandler{userService: us, validator: v, log: logger}
}

// User returns the user's id, name, and link to the profile picture.
// /users/{id}
func (h *UserHandler) UserSummary(c *echo.Context) error {
	ctx := c.Request().Context()

	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "Invalid user id: NAN"})
	}

	profile, err := h.userService.Profile(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, map[string]any{
				"id":    userID,
				"error": "User not found",
			})
		} else {
			return c.JSON(http.StatusInternalServerError, map[string]any{"error": "Failed to get user"})
		}
	}

	return c.JSON(http.StatusOK, UserSummaryResponse{
		ID:         profile.UserID,
		Name:       profile.Name,
		PictureURL: profile.PictureURL,
	})
}

// // Profile returns the user's id and "about me" type information.
// // /users/{id}/profile
func (h *UserHandler) UserProfile(c *echo.Context) error {
	ctx := c.Request().Context()

	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "Invalid user id: NAN"})
	}

	profile, err := h.userService.Profile(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, map[string]any{
				"id":    userID,
				"error": "User not found",
			})
		} else {
			return c.JSON(http.StatusInternalServerError, map[string]any{"error": "Failed to get user"})
		}
	}

	return c.JSON(http.StatusOK, ProfileResponse{
		ID:         profile.UserID,
		Name:       profile.Name,
		Age:        profile.Age,
		PictureURL: profile.PictureURL,
		Bio:        profile.Bio,
		IsOnline:   profile.IsOnline,
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

	profile, err := h.userService.Profile(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, map[string]any{
				"id":    userID,
				"error": "User not found",
			})
		} else {
			return c.JSON(http.StatusInternalServerError, map[string]any{"error": "Failed to get user"})
		}
	}

	activitiesResponce := []ActivityResp{}

	if len(profile.Activities) > 0 {
		for _, a := range profile.Activities {
			activity := ActivityResp{
				ID:            a.Activity.ID,
				Title:         a.Activity.Title,
				Experience:    a.Experience,
				InterestLevel: a.InterestLevel,
			}
			activitiesResponce = append(activitiesResponce, activity)
		}
	}

	return c.JSON(http.StatusOK, ProfileResponse{
		ID:                   profile.UserID,
		MaxRadius:            profile.MaxRadius,
		InteractionModeID:    profile.InteractionModeID,
		InteractionModeTitle: profile.InteractionMode.Title,
		Activities:           activitiesResponce,
	})
}

// // UpdateProfile
// func (h *UserHandler) UpdateProfile(c *echo.Context) error {
// 	ctx := c.Request().Context()

// 	var req UpdateProfileRequest

// 	if err := c.Bind(&req); err != nil {
// 		return c.JSON(http.StatusBadRequest, map[string]any{
// 			"error": "invalid json",
// 		})
// 	}

// 	if err := h.validator.Struct(req); err != nil {
// 		return c.JSON(http.StatusBadRequest, map[string]any{
// 			"error": "validation failed: " + err.Error(),
// 		})
// 	}

// 	// Build domain model out of DTO
// 	profile := &domain.Profile{
// 		Name:      req.Name,
//		PictureURL req.Picture // add/remove/change their profile picture
// 		Age:       req.Age,
// 		Bio:       req.Bio,
// 		MaxRadius: req.MaxRadius,
// 	}

// 	if err := h.userService.Update(ctx, profile); err != nil {
// 		return err
// 	}

// 	return c.JSON(http.StatusOK, profile)
// }
