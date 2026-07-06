package restapi

import (
	"context"
	"log/slog"

	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func SetupRouter(log *slog.Logger, JWTSecret string) *echo.Echo {

	e := echo.New()

	// 		add middleware
	e.Use(middleware.RequestID())
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogMethod:    true,
		LogURI:       true,
		LogStatus:    true,
		LogRequestID: true,
		LogLatency:   true,
		HandleError:  true,
		LogValuesFunc: func(c *echo.Context, v middleware.RequestLoggerValues) error {
			attrs := []slog.Attr{
				slog.String("method", v.Method),
				slog.String("uri", v.URI),
				slog.Int("status", v.Status),
				slog.String("requestid", v.RequestID),
				slog.Duration("duration", v.Latency),
			}

			level := slog.LevelInfo
			msg := "REQUEST"

			if v.Error != nil {
				level = slog.LevelError
				msg = "REQUEST_ERROR"
				attrs = append(attrs, slog.String("err", v.Error.Error()))
			}

			log.LogAttrs(context.Background(), level, msg, attrs...)

			return nil
		},
	}))
	e.Use(middleware.Recover())

	// public routes
	// public := e.Group("")
	// public.POST("/auth/register", authHandler.SignUp)
	// public.POST("/auth/login", authHandler.SignIn)

	// private routes
	private := e.Group("")
	private.Use(echojwt.JWT([]byte(JWTSecret)))

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
