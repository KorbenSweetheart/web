package handlers

import (
	"context"
	"errors"
	"log/slog"
	"match-me-api/internal/domain"
	"match-me-api/internal/transport/httpserver/dto"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

const (
	AccessTokenCookieName  = "access_token"
	RefreshTokenCookieName = "refresh_token"
)

type Authenticator interface {
	Register(ctx context.Context, name, email, password string) (*domain.Account, error)
	Login(ctx context.Context, email, password string) (string, string, error)
	Refresh(ctx context.Context, rawRefreshToken string) (string, string, error)
	Logout(ctx context.Context, userID int64) error
}

type AuthHandler struct {
	auth             Authenticator
	validator        *validator.Validate
	accessCookieTTL  time.Duration
	refreshCookieTTL time.Duration
	log              *slog.Logger
}

func NewAuthHandler(auth Authenticator, v *validator.Validate, atTTL, rtTTL time.Duration, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{auth: auth, validator: v, accessCookieTTL: atTTL, refreshCookieTTL: rtTTL, log: logger}
}

// @Register godoc
// @Summary      Register a new user
// @Description  Create a new user account with name, email, and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.RegisterRequest true "User Registration Info"
// @Success      201 {object} dto.RegistrationResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *echo.Context) error {
	ctx := c.Request().Context()

	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid request body",
		})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid request body",
		})
	}

	account, err := h.auth.Register(ctx, req.Name, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrEmailIsTaken) {
			return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Message: "Email is already taken",
			})
		} else {
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Message: "Internal server error",
			})
		}
	}

	return c.JSON(http.StatusCreated, dto.RegistrationResponse{
		ID:      account.ID,
		Email:   account.Email,
		Name:    account.Profile.Name,
		Message: "User registered successfully",
	})
}

// @Login godoc
// @Summary      Login a user
// @Description  Login with email and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.LoginRequest true "User Login Info"
// @Success      200 {object} dto.LoginResponse
// @Failure      400 {object} dto.ErrorResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *echo.Context) error {
	ctx := c.Request().Context()

	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid request body",
		})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid request body",
		})
	}

	accessToken, refreshToken, err := h.auth.Login(ctx, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) || errors.Is(err, domain.ErrInvalidCreds) {
			return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Message: "Invalid credentials",
			})
		} else {
			return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Message: "Internal server error",
			})
		}
	}

	h.setAuthCookies(c, accessToken, refreshToken)

	return c.JSON(http.StatusOK, dto.LoginResponse{
		Message:     "success",
		AccessToken: accessToken,
	})
}

// @Refresh godoc
// @Summary      Refresh access token
// @Description  Refresh access token using refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200 {object} dto.LoginResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /auth/refresh [post]
func (h *AuthHandler) Refresh(c *echo.Context) error {
	ctx := c.Request().Context()

	cookie, err := c.Cookie(RefreshTokenCookieName)
	if err != nil || cookie.Value == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "missing refresh token",
		})
	}

	accessToken, refreshToken, err := h.auth.Refresh(ctx, cookie.Value)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidOrExpiredToken) {
			h.clearAuthCookies(c)
			return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Message: "Invalid or expired token",
			})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Internal server error",
		})
	}

	h.setAuthCookies(c, accessToken, refreshToken)

	return c.JSON(http.StatusOK, dto.LoginResponse{
		Message:     "success",
		AccessToken: accessToken,
	})
}

// @Logout godoc
// @Summary      Logout a user
// @Description  Logout a user
// @Tags         auth
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} dto.OKResponse
// @Failure      401 {object} dto.ErrorResponse
// @Failure      500 {object} dto.ErrorResponse
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *echo.Context) error {
	ctx := c.Request().Context()

	userID, ok := c.Get("user_id").(int64)
	if !ok {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Unauthorized",
		})
	}

	if err := h.auth.Logout(ctx, userID); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Internal server error",
		})
	}

	h.clearAuthCookies(c)

	return c.JSON(http.StatusOK, dto.OKResponse{
		Message: "Logged out",
	})
}

func (h *AuthHandler) setAuthCookies(c *echo.Context, accessToken, refreshToken string) {
	accessCookie := &http.Cookie{
		Name:     AccessTokenCookieName,
		Value:    accessToken,
		Expires:  time.Now().Add(h.accessCookieTTL),
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // true in production (HTTPS)
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(accessCookie)

	refreshCookie := &http.Cookie{
		Name:     RefreshTokenCookieName,
		Value:    refreshToken,
		Expires:  time.Now().Add(h.refreshCookieTTL),
		Path:     "/auth/refresh",
		HttpOnly: true,
		Secure:   false, // true in production (HTTPS)
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(refreshCookie)
}

func (h *AuthHandler) clearAuthCookies(c *echo.Context) {
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
}
