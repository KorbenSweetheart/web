package logger

import (
	"errors"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func SetupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}

// Err safely wraps an error and deeply unwraps it to find structured logs.
func Err(err error) slog.Attr {
	if err == nil {
		return slog.Attr{}
	}

	// Traverse the error chain looking for a slog.LogValuer
	var valuer slog.LogValuer
	if errors.As(err, &valuer) {
		return slog.Attr{
			Key: "error",
			// We still log the full error string to preserve the "wrap:" context,
			// but we append the structured data from the inner error as well.
			Value: slog.GroupValue(
				slog.String("msg", err.Error()),
				slog.Any("details", valuer.LogValue()),
			),
		}
	}

	// Maximum Performance Path for standard errors
	return slog.String("error", err.Error())
}
