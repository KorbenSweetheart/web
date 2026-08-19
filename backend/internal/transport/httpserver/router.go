package httpserver

import (
	"log/slog"
	"match-me-api/internal/config"
	"match-me-api/internal/transport/httpserver/handlers"
	mdlwrconfig "match-me-api/internal/transport/httpserver/middleware"
	"net/http"

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
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:8080", "http://localhost:5173"}, // TODO: move to config vars
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, echo.HeaderXCSRFToken},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete, http.MethodOptions},
	}))

	// Healthchecks routes
	e.GET("/healthz", h.Health.Healthz)
	e.GET("/readyz", h.Health.Readyz)

	// Public routes
	public := e.Group("")
	public.POST("/auth/register", h.Auth.Register)
	public.POST("/auth/login", h.Auth.Login)
	public.POST("/auth/refresh", h.Auth.Refresh)

	// Private routes
	private := e.Group("")
	private.Use(mdlwrconfig.JWTMiddlewareWithConfig(cfg.TM.JWTSecretKey, handlers.AccessTokenCookieName))
	private.POST("/auth/logout", h.Auth.Logout)

	// Users
	private.GET("/users/:id", h.User.UserSummary)         // /users/{id}
	private.GET("/users/:id/profile", h.User.UserProfile) // /users/{id}/profile
	private.GET("/users/:id/bio", h.User.UserBio)         // /users/{id}/bio

	// "Me" Shortcuts
	private.GET("/me", h.User.MySummary)         // /me
	private.GET("/me/profile", h.User.MyProfile) // /me/profile
	private.GET("/me/bio", h.User.MyBio)         // /me/bio
	private.PATCH("/me/profile", h.User.UpdateProfile)
	private.POST("/me/picture", h.User.UploadProfilePicture)
	private.DELETE("/me/picture", h.User.DeleteProfilePicture)

	// Connections
	private.GET("/connections", h.Conn.Connections) // /connections
	private.GET("/connections/requests", h.Conn.ConnectionRequests)
	private.POST("/connections", h.Conn.CreateConnection)            // request body: {"to_user_id": 123}, Response: 200 OK {"status": "pending"}
	private.PATCH("/connections/:id", h.Conn.UpdateConnectionStatus) // request body: {"status": "accepted" | "declined"}, Response: 200 OK
	private.DELETE("/connections/:id", h.Conn.DeleteConnection)      // Response: 204 No Content

	// Recommendations
	private.GET("/recommendations", h.Match.Recommendations)

	// Chat
	private.GET("/ws", h.WS.Upgrade)

	// Dictionary
	private.GET("/activities", h.Dictionary.Activities) // /actvities

	return e
}
