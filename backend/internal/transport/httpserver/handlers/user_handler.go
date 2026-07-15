package handlers

import (
	"match-me-api/internal/service"
	"net/http"

	"github.com/labstack/echo/v5"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(us *service.UserService) *UserHandler {
	return &UserHandler{userService: us}
}

// /users/{id}
func (h *UserHandler) User(c *echo.Context) error {
	id := c.Param("id")
	ctx := c.Request().Context()

	user, err := h.userService.GetProfile(ctx, id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	// /users/{id} (id, name, avatar)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":         user.ID,
		"name":       user.UserName,
		"avatar_url": user.AvatarURL,
	})
}

// /users/{id}/profile
func (h *UserHandler) UserProfile(c *echo.Context) error {
	id := c.Param("id")
	ctx := c.Request().Context()

	user, err := h.userService.GetProfile(ctx, id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Profile not found"})
	}

	// Возвращаем "about me" информацию
	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":       user.ID,
		"about_me": user.AboutMe,
	})
}

// GET /users/{id}/bio
func (h *UserHandler) AboutUser(c *echo.Context) error {
	id := c.Param("id")
	ctx := c.Request().Context()

	user, err := h.userService.GetProfile(ctx, id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "About Me information is not found"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"id":        user.ID,
		"interests": user.Interests,
		// other logic to match users
	})
}
