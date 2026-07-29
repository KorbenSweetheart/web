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
	// UpdateLocation(ctx context.Context, id int64, lat, lon float64) error
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

// UserSummary returns the user's id, name, and link to the profile picture.
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

// // UserProfile returns the user's id and "about me" type information.
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

	activitiesResponce := []ActivityResponse{}

	if len(profile.Activities) > 0 {
		for _, a := range profile.Activities {
			activity := ActivityResponse{
				ID:            a.Activity.ID,
				Title:         a.Activity.Title,
				Experience:    a.Experience,
				InterestLevel: a.InterestLevel,
			}
			activitiesResponce = append(activitiesResponce, activity)
		}
	}

	return c.JSON(http.StatusOK, UserBioResponse{
		ID:              profile.UserID,
		MaxRadius:       profile.MaxRadius,
		InteractionMode: profile.InteractionMode,
		Activities:      activitiesResponce,
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

	profile, err := h.userService.Profile(ctx, myID)
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

	return c.JSON(http.StatusOK, UserSummaryResponse{
		ID:         profile.UserID,
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

	profile, err := h.userService.Profile(ctx, myID)
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

	return c.JSON(http.StatusOK, ProfileResponse{
		ID:         profile.UserID,
		Name:       profile.Name,
		Age:        profile.Age,
		PictureURL: profile.PictureURL,
		Bio:        profile.Bio,
		IsOnline:   profile.IsOnline,
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

	profile, err := h.userService.Profile(ctx, myID)
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

	activitiesResponce := []ActivityResponse{}

	if len(profile.Activities) > 0 {
		for _, a := range profile.Activities {
			activity := ActivityResponse{
				ID:            a.Activity.ID,
				Title:         a.Activity.Title,
				Experience:    a.Experience,
				InterestLevel: a.InterestLevel,
			}
			activitiesResponce = append(activitiesResponce, activity)
		}
	}

	return c.JSON(http.StatusOK, UserBioResponse{
		ID:              profile.UserID,
		MaxRadius:       profile.MaxRadius,
		InteractionMode: profile.InteractionMode,
		Activities:      activitiesResponce,
	})
}

// func (h *UserHandler) UpdateLocation(c *echo.Context) error {
// 	ctx := c.Request().Context()

// 	var req UpdateLocationRequest
// 	if err := c.Bind(&req); err != nil {
// 		return c.JSON(http.StatusBadRequest, map[string]any{"error": "invalid payload"})
// 	}

// 	userID := c.Get("user_id").(int64)

// 	if err := h.userService.UpdateLocation(ctx, userID, req.Lat, req.Lon); err != nil {
// 		return c.JSON(http.StatusInternalServerError, map[string]any{"error": "failed to update location"})
// 	}

// 	return c.JSON(http.StatusOK, map[string]any{"status": "location updated"})
// }

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
