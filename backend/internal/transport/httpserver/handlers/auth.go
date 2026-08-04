package handlers

import (
	"context"
	"errors"
	"log/slog"
	"match-me-api/internal/domain"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

const (
	AccessTokenCookieName  = "access_token"
	RefreshTokenCookieName = "refresh_token"
)

type AuthService interface {
	Register(ctx context.Context, name, email, password string) (*domain.Account, error)
	Login(ctx context.Context, email, password string) (string, string, error)
	Logout(ctx context.Context, userID int64) error
}

type AuthHandler struct {
	authService      AuthService
	validator        *validator.Validate
	accessCookieTTL  time.Duration
	refreshCookieTTL time.Duration
	log              *slog.Logger
}

func NewAuthHandler(as AuthService, v *validator.Validate, atTTL, rtTTL time.Duration, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{authService: as, validator: v, accessCookieTTL: atTTL, refreshCookieTTL: rtTTL, log: logger}
}

// Register handler
func (h *AuthHandler) Register(c *echo.Context) error {
	ctx := c.Request().Context()

	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error": "invalid json",
		})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error": "validation failed: " + err.Error(),
		})
	}

	account, err := h.authService.Register(ctx, req.Name, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrEmailIsTaken) {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"email":   req.Email,
				"name":    req.Name,
				"error":   domain.ErrEmailIsTaken.Error(),
				"message": "An account with this email already exists.",
			})
		} else {
			return c.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
		}
	}

	// Login in or redirect to login upon successful registration

	return c.JSON(http.StatusCreated, map[string]any{
		"id":      account.ID,
		"email":   account.Email,
		"name":    account.Profile.Name,
		"message": "user registered successfully",
	})
}

// Login handler
func (h *AuthHandler) Login(c *echo.Context) error {
	ctx := c.Request().Context()

	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error": "faild to read body",
		})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{
			"error": "validation failed: " + err.Error(),
		})
	}

	accessToken, refreshToken, err := h.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) || errors.Is(err, domain.ErrInvalidCreds) {
			return c.JSON(http.StatusUnauthorized, map[string]any{
				"email":   req.Email,
				"error":   domain.ErrInvalidCreds.Error(),
				"message": "invalid input body",
			})
		} else {
			return c.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
		}
	}

	// Issue access token cookie
	accessCookie := &http.Cookie{
		Name:     AccessTokenCookieName,
		Value:    accessToken,
		Expires:  time.Now().Add(h.accessCookieTTL),
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // true в production (HTTPS)
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(accessCookie)

	// Issue refresh token cookie
	refreshCookie := &http.Cookie{
		Name:     RefreshTokenCookieName,
		Value:    refreshToken,
		Expires:  time.Now().Add(h.refreshCookieTTL),
		Path:     "/auth/refresh",
		HttpOnly: true,
		Secure:   false, // true в production (HTTPS)
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(refreshCookie)

	return c.JSON(http.StatusOK, map[string]any{
		"message":      "success",
		"access_token": accessToken,
	})
}

func (h *AuthHandler) Logout(c *echo.Context) error {
	ctx := c.Request().Context()

	userID, ok := c.Get("user_id").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]any{"error": "Unauthorized"})
	}

	if err := h.authService.Logout(ctx, userID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
	}

	c.SetCookie(&http.Cookie{
		Name:     AccessTokenCookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
	})

	c.SetCookie(&http.Cookie{
		Name:     RefreshTokenCookieName,
		Value:    "",
		Path:     "/auth/refresh",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
	})

	return c.JSON(http.StatusOK, map[string]any{
		"message": "Logged out",
	})
}
