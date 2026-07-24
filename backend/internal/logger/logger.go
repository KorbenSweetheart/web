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

// ErrValue converts an error into a slog.Value for use in key-value logging pairs.
// TODO: Rethink or just get rid of it and user slog.Attr
func ErrValue(err error) slog.Value {
	if err == nil {
		return slog.StringValue("")
	}

	var valuer slog.LogValuer
	if errors.As(err, &valuer) {
		return slog.GroupValue(
			slog.String("msg", err.Error()),
			slog.Any("details", valuer.LogValue()),
		)
	}

	return slog.StringValue(err.Error())
}
