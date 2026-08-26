package main

import (
	"context"
	"errors"
	"log/slog"
	"match-me-api/internal/config"
	"match-me-api/internal/logger"
	"match-me-api/internal/pkg/tokenmgr"
	"match-me-api/internal/service"
	"match-me-api/internal/storage/memory"
	"match-me-api/internal/storage/minio"
	"match-me-api/internal/storage/postgres"
	"match-me-api/internal/transport/httpserver"
	"match-me-api/internal/transport/httpserver/handlers"
	"match-me-api/internal/transport/websocket"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
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

	// connect to repository
	dbCtx, dbCancel := context.WithTimeout(shutdownCtx, 10*time.Second)
	defer dbCancel()

	storage, err := postgres.NewPostgresDB(dbCtx, cfg.DB, log)
	if err != nil {
		log.Error("database connection failed", logger.Err(err))
		os.Exit(1)
	}

	log.Info("database connection established successfully")

	// create token manager
	tm := tokenmgr.NewTokenManager(
		cfg.TM.JWTSecretKey,
		cfg.TM.TokenIssuer,
	)

	minioStorage, err := minio.NewMinioStorage(shutdownCtx, cfg.MinIO, log)
	if err != nil {
		log.Error("minio storage connection failed", logger.Err(err))
		os.Exit(1)
	}

	// create storage
	memoryStorage := memory.NewStorage()

	// create services
	authService := service.NewAuthService(storage, storage, tm, cfg.TM.AccessTokenTTL, cfg.TM.RefreshTokenTTL, log)
	dictionaryService := service.NewDictionaryService(storage, log)
	userService := service.NewUserService(storage, minioStorage, storage, storage, log)
	matchService := service.NewMatchService(storage, log)
	connectionService := service.NewConnectionService(storage, storage, log)
	chatService := service.NewChatService(storage, storage, log)
	presenceService := service.NewPresenceService(memoryStorage, log)

	// init validator
	validator := validator.New(validator.WithRequiredStructEnabled())

	// create websocket hub & handler
	wsHub := websocket.NewHub(chatService, presenceService, validator, log)
	go wsHub.Run(shutdownCtx)
	wsHandler := websocket.NewHandler(wsHub, log)

	// Handlers
	healthHandler := handlers.NewHealthHandler(storage)
	dictionaryHandler := handlers.NewDictionaryHandler(dictionaryService, validator, log)
	authHandler := handlers.NewAuthHandler(authService, validator, cfg.TM.AccessTokenTTL, cfg.TM.RefreshTokenTTL, log)
	userHandler := handlers.NewUserHandler(userService, validator, log)
	matchHandler := handlers.NewMatchHandler(matchService, validator, log)
	connectionHandler := handlers.NewConnectionHandler(connectionService, validator, log)
	chatHandler := handlers.NewChatHandler(chatService, validator, log)

	// setup router/server (Echo)
	e := httpserver.SetupRouter(cfg, log, handlers.Handlers{
		Health:     healthHandler,
		Dictionary: dictionaryHandler,
		Auth:       authHandler,
		User:       userHandler,
		Match:      matchHandler,
		Conn:       connectionHandler,
		Chat:       chatHandler,
		WS:         wsHandler,
	})

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
		// TODO: storage.Close()
		os.Exit(1)
	}

	log.Info("http server stopped, cleaning up resources...")

	// if err := db.Close(); err != nil {
	//     log.Error("Error during database shutdown", logger.Err(err))
	// }

	log.Info("api server shutdown completed successfully")
}
