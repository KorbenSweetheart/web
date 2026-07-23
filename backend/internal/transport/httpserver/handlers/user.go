package handlers

import (
	"context"
	"errors"
	"log/slog"
	"match-me-api/internal/domain"
	"match-me-api/internal/storage"
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

// User returns the user's name and link to the profile picture.
// /users/{id}
func (h *UserHandler) User(c *echo.Context) error {
	ctx := c.Request().Context()

	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "Invalid user id: NAN"})
	}

	profile, err := h.userService.Profile(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, map[string]any{"error": "User not found"})
		} else {
			return c.JSON(http.StatusInternalServerError, map[string]any{"error": "Failed to get user"})
		}
	}

	return c.JSON(http.StatusOK, UserSummaryResponse{
		Name:       profile.Name,
		PictureURL: profile.PictureURL,
	})
}

// // Profile returns the users "about me" type information.
// // /users/{id}/profile
// func (h *UserHandler) UserProfile(c *echo.Context) error {
// 	ctx := c.Request().Context()

// 	id, err := idValidation(c.Param("id"))
// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID: NAN"})
// 	}

// 	user, err := h.userService.GetProfile(ctx, id)
// 	if err != nil {
// 		return c.JSON(http.StatusNotFound, map[string]string{"error": "Profile not found"})
// 	}

// 	return c.JSON(http.StatusOK, map[string]interface{}{
// 		"id":       user.UserID,
// 		"about_me": user.Bio,
// 	})
// }

// // UserBio returns the users biographical data (the data used to power recommendations).
// // /users/{id}/bio
// func (h *UserHandler) UserBio(c *echo.Context) error {
// 	ctx := c.Request().Context()

// 	id, err := idValidation(c.Param("id"))
// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID: NAN"})
// 	}

// 	user, err := h.userService.GetProfile(ctx, id)
// 	if err != nil {
// 		return c.JSON(http.StatusNotFound, map[string]string{"error": "About Me information is not found"})
// 	}

// 	return c.JSON(http.StatusOK, map[string]interface{}{
// 		"id":        user.UserID,
// 		"interests": user.Activities,
// 		// other logic to match users
// 	})
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
// 		Age:       req.Age,
// 		Bio:       req.Bio,
// 		MaxRadius: req.MaxRadius,
// 	}

// 	if err := h.userService.Update(ctx, profile); err != nil {
// 		return err
// 	}

// 	return c.JSON(http.StatusOK, profile)
// }
