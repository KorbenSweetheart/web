package httpserver

import (
	"log/slog"
	"match-me-api/internal/transport/httpserver/handlers"
	"match-me-api/internal/transport/httpserver/utils"
	"net/http"

	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func SetupRouter(
	JWTSecret string,
	as handlers.AuthService,
	log *slog.Logger,
) *echo.Echo {

	e := echo.New()

	// Middleware
	e.Use(middleware.RequestID())
	e.Use(middleware.RequestLoggerWithConfig(utils.LoggerConfig(log)))
	e.Use(middleware.Recover())
	e.Use(middleware.CORS("http://localhost:8080", "http://localhost:5173")) // TODO: move to config vars

	//
	authHandler := handlers.NewAuthHandler(as, log)
	// userHandler := handlers.NewUserHandler(us, log)

	// Public routes
	public := e.Group("")
	public.GET("/", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "Hello, World!"})
	})
	public.GET("/health", handlers.CheckHealth)
	public.POST("/auth/register", authHandler.Register)
	// public.POST("/auth/login", authHandler.Login)

	// Private routes
	private := e.Group("")
	private.Use(echojwt.JWT([]byte(JWTSecret))) // JWT Middleware

	// Users
	// private.GET("/users/:id", userHandler.GetBaseInfo)        // /users/{id}
	// private.GET("/users/:id/profile", userHandler.GetProfile) // /users/{id}/profile
	// private.GET("/users/:id/bio", userHandler.GetBio)         // /users/{id}/bio

	// Shortcuts
	// private.GET("/me", userHandler.GetMyBaseInfo)        // /me
	// private.GET("/me/profile", userHandler.GetMyProfile) // /me/profile
	// private.GET("/me/bio", userHandler.GetMyBio)         // /me/bio

	// Recommendations
	// private.GET("/recommendations", matchHandler.GetRecommendations)
	// private.GET("/connections", matchHandler.GetConnections)

	return e
}
