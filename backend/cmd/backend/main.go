package main

import (
	"context"
	"errors"
	"log/slog"
	"match-me-api/internal/config"
	"match-me-api/internal/controller/restapi"
	"match-me-api/internal/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
)

func main() {
	_ = godotenv.Load()

	shutdownCtx, shutdown := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer shutdown()

	// init config
	cfg := config.MustLoad()

	// setup logger
	log := logger.SetupLogger(cfg.Env)

	log.Info("Setup API server", slog.String("env", cfg.Env))

	// init repository

	// create usecases

	// setup router (Echo)
	e := restapi.SetupRouter(log, cfg.JWTSecret)

	// start server
	sc := echo.StartConfig{
		Address:         cfg.HTTPServer.Address,
		GracefulTimeout: cfg.HTTPServer.ShutdownTimeout,
		BeforeServeFunc: func(s *http.Server) error {
			s.ReadTimeout = cfg.HTTPServer.Timeout
			s.WriteTimeout = cfg.HTTPServer.Timeout
			s.IdleTimeout = cfg.HTTPServer.IdleTimeout
			return nil
		},
	}

	log.Info("Starting API server", slog.String("addr", cfg.HTTPServer.Address))

	if err := sc.Start(shutdownCtx, e); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("API server stopped with error", logger.Err(err))
		// db.Close()
		os.Exit(1)
	}

	log.Info("HTTP server stopped, cleaning up resources...")

	// if err := db.Close(); err != nil {
	//     log.Error("Error during database shutdown", logger.Err(err))
	// }

	log.Info("API server shutdown completed successfully")
}
