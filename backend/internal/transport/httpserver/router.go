package httpserver

import (
	"log/slog"
	"match-me-api/internal/config"
	"match-me-api/internal/transport/httpserver/handlers"
	"match-me-api/internal/transport/httpserver/utils"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func SetupRouter(
	cfg *config.Config,
	log *slog.Logger,
	as handlers.AuthService,
	us handlers.UserService,
) *echo.Echo {

	e := echo.New()

	// Middleware
	e.Use(middleware.RequestID())
	e.Use(middleware.RequestLoggerWithConfig(utils.LoggerConfig(log)))
	e.Use(middleware.Recover())
	e.Use(middleware.CORS("http://localhost:8080", "http://localhost:5173")) // TODO: move to config vars

	validate := validator.New(validator.WithRequiredStructEnabled())

	// Handlers
	authHandler := handlers.NewAuthHandler(as, validate, cfg.TM.AccessTokenTTL, cfg.TM.RefreshTokenTTL, log)
	userHandler := handlers.NewUserHandler(us, validate, log)

	// Public routes
	public := e.Group("")
	public.GET("/health", handlers.CheckHealth)
	public.POST("/auth/register", authHandler.Register)
	public.POST("/auth/login", authHandler.Login)
	// public.POST("/auth/refresh", authHandler.Refresh)
	// public.POST("/auth/logout", authHandler.Logout)

	// Private routes
	private := e.Group("")
	private.Use(utils.JWTMiddlewareWithConfig(cfg.TM.JWTSecretKey, handlers.AccessTokenCookieName))

	// Users
	private.GET("/users/:id", userHandler.UserSummary)         // /users/{id}
	private.GET("/users/:id/profile", userHandler.UserProfile) // /users/{id}/profile
	private.GET("/users/:id/bio", userHandler.UserBio)         // /users/{id}/bio

	// Shortcuts
	// private.GET("/me", userHandler.GetMyBaseInfo)        // /me
	// private.GET("/me/profile", userHandler.GetMyProfile) // /me/profile
	// private.GET("/me/bio", userHandler.GetMyBio)         // /me/bio

	// Recommendations
	// private.GET("/recommendations", matchHandler.GetRecommendations)
	// private.GET("/connections", matchHandler.GetConnections)

	return e
}
