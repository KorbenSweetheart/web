package middleware

import (
	"context"
	"log/slog"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func LoggerConfig(log *slog.Logger) middleware.RequestLoggerConfig {

	var skipper = func(c *echo.Context) bool {
		// Skip the health check endpoint.
		return c.Request().URL.Path == "/health"
	}

	loggerConfig := middleware.RequestLoggerConfig{
		LogMethod:    true,
		LogURI:       true,
		Skipper:      skipper,
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
	}
	return loggerConfig
}
