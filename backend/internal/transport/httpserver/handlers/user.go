package handlers

// import (
// 	"fmt"
// 	"log/slog"
// 	"match-me-api/internal/service"
// 	"net/http"
// 	"strconv"

// 	"github.com/labstack/echo/v5"
// )

// type UserService interface {
// 	// UserByEmail(ctx context.Context, email string) (*domain.User, error)
// 	// UserByID(ctx context.Context, id int64) (*domain.User, error)
// 	// ProfileByID(ctx context.Context, id int64) (*domain.Profile, error)
// 	// UpdateProfile(ctx context.Context, profile *domain.Profile) error
// }

// type UserHandler struct {
// userService UserService
// 	log         *slog.Logger
// }

// func NewUserHandler(us UserService, logger *slog.Logger) *UserHandler {
// 	return &UserHandler{userService: us, log: logger}
// }

// // /users/{id}
// func (h *UserHandler) User(c *echo.Context) error {
// 	ctx := c.Request().Context()

// 	id, err := idValidation(c.Param("id"))
// 	if err != nil {
// 		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID: NAN"})
// 	}

// 	user, err := h.userService.GetProfile(ctx, id)
// 	if err != nil {
// 		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
// 	}

// 	// /users/{id} (id, name, avatar)
// 	return c.JSON(http.StatusOK, map[string]interface{}{
// 		"id":         user.UserID,
// 		"name":       user.Name,
// 		"avatar_url": user.PictureURL,
// 	})
// }

// // /users/{id}/profile
// func (h *UserHandler) Profile(c *echo.Context) error {
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

// // GET /users/{id}/bio
// func (h *UserHandler) AboutUser(c *echo.Context) error {
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

// func idValidation(idStr string) (int64, error) {
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		return 0, fmt.Errorf("failed to parse userID, error: %w", err)
// 	}

// 	return int64(id), nil
// }
