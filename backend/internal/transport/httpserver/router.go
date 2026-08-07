package httpserver

import (
	"log/slog"
	"match-me-api/internal/config"
	"match-me-api/internal/transport/httpserver/handlers"
	mdlwrconfig "match-me-api/internal/transport/httpserver/middleware"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func SetupRouter(
	cfg *config.Config,
	log *slog.Logger,
	h handlers.Handlers,
) *echo.Echo {

	e := echo.New()

	// Middleware
	e.Use(middleware.RequestID())
	e.Use(middleware.RequestLoggerWithConfig(mdlwrconfig.LoggerConfig(log)))
	e.Use(middleware.Recover())
	e.Use(middleware.CORS("http://localhost:8080", "http://localhost:5173")) // TODO: move to config vars

	// Healthchecks routes
	e.GET("/healthz", h.Health.Healthz)
	e.GET("/readyz", h.Health.Readyz)

	// Public routes
	public := e.Group("")
	public.POST("/auth/register", h.Auth.Register)
	public.POST("/auth/login", h.Auth.Login)
	// public.POST("/auth/refresh", h.Auth.Refresh)

	// Private routes
	private := e.Group("")
	private.Use(mdlwrconfig.JWTMiddlewareWithConfig(cfg.TM.JWTSecretKey, handlers.AccessTokenCookieName))
	private.POST("/auth/logout", h.Auth.Logout)

	// Users
	private.GET("/users/:id", h.User.UserSummary)         // /users/{id}
	private.GET("/users/:id/profile", h.User.UserProfile) // /users/{id}/profile
	private.GET("/users/:id/bio", h.User.UserBio)         // /users/{id}/bio
	// private.GET("/connections", h.User.Connections) // /connections
	private.GET("/activities", h.Dictionary.Activities) // /actvities

	// Shortcuts
	private.GET("/me", h.User.MySummary)         // /me
	private.GET("/me/profile", h.User.MyProfile) // /me/profile
	private.GET("/me/bio", h.User.MyBio)         // /me/bio
	private.PATCH("/me/profile", h.User.UpdateProfile)

	// Recommendations
	private.GET("/recommendations", h.Match.Recommendations)

	return e
}
