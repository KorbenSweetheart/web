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

// Err creates a standardized slog.Attr for errors, including unwrapped cause chains.
func Err(err error) slog.Attr {
	if err == nil {
		return slog.Attr{}
	}

	attrs := []slog.Attr{
		slog.String("msg", err.Error()),
	}

	// Unroll error chain if cause exists
	if cause := errors.Unwrap(err); cause != nil {
		var chain []string
		for curr := cause; curr != nil; curr = errors.Unwrap(curr) {
			chain = append(chain, curr.Error())
		}
		attrs = append(attrs, slog.Any("cause_chain", chain))
	}

	// Capture LogValuer details if supported anywhere in the chain
	var valuer slog.LogValuer
	if errors.As(err, &valuer) {
		attrs = append(attrs, slog.Any("details", valuer.LogValue()))
	}

	return slog.Attr{
		Key:   "error",
		Value: slog.GroupValue(attrs...),
	}
}
