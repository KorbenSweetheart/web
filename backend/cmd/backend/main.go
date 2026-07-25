package main

import (
	"context"
	"errors"
	"log/slog"
	"match-me-api/internal/config"
	"match-me-api/internal/logger"
	"match-me-api/internal/pkg/tokenmgr"
	"match-me-api/internal/service"
	"match-me-api/internal/storage/postgres"
	"match-me-api/internal/transport/httpserver"
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

	log.Info("setup api server", slog.String("env", cfg.Env))

	// init storage/repo
	storage, err := postgres.NewPostgresDB(cfg.DB, log)
	if err != nil {
		log.Error("database connection failed", logger.Err(err))
		os.Exit(1)
	}

	if err := storage.AutoMigrate(); err != nil {
		log.Error("migration failed", logger.Err(err))
		os.Exit(1)
	}

	if err := storage.SeedData(); err != nil {
		log.Error("seeding failed", logger.Err(err))
		os.Exit(1)
	}

	tm := tokenmgr.NewTokenManager(
		cfg.TM.JWTSecretKey,
		cfg.TM.TokenIssuer,
	)

	// create services
	authService := service.NewAuthService(storage, tm, cfg.TM.AccessTokenTTL, cfg.TM.RefreshTokenTTL, log)
	userService := service.NewUserService(storage, log)

	// setup router/server (Echo)
	e := httpserver.SetupRouter(cfg, log, authService, userService)

	// server config
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

	log.Info("starting api server", slog.String("addr", cfg.HTTPServer.Address))

	if err := sc.Start(shutdownCtx, e); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("api server stopped with error", logger.Err(err))
		// db.Close()
		os.Exit(1)
	}

	log.Info("http server stopped, cleaning up resources...")

	// if err := db.Close(); err != nil {
	//     log.Error("Error during database shutdown", logger.Err(err))
	// }

	log.Info("api server shutdown completed successfully")
}
