package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

const (
	tokenIssuer               = "match-me-api"
	accessTokenTTL            = 60 * time.Minute    // 15 min
	refreshTokenTTL           = 30 * 24 * time.Hour // 30 days
	defaultSrvTimeout         = 4 * time.Second
	defaultSrvIdleTimeout     = 60 * time.Second
	defaultSrvShutdownTimeout = 10 * time.Second
	defaultMinIOUseSSL        = false
)

type Config struct {
	Env        string
	TM         TokenMgr
	DB         Database
	MinIO      MinIOConfig
	HTTPServer HTTPServer
}

type Database struct {
	Name string
	Host string
	Port string
	User string
	Pass string
}

type HTTPServer struct {
	Address         string
	Timeout         time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

type TokenMgr struct {
	JWTSecretKey    string
	TokenIssuer     string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

type MinIOConfig struct {
	Endpoint              string
	PublicEndpoint        string
	AccessKey             string
	SecretKey             string
	UseSSL                bool
	ProfilePicturesBucket string
}

func MustLoad() *Config {
	var cfg Config

	cfg.Env = os.Getenv("ENVIRONMENT")
	if cfg.Env == "" {
		cfg.Env = "dev"
	}

	// Token Manager
	cfg.TM.JWTSecretKey = os.Getenv("JWT_SECRET")
	if cfg.TM.JWTSecretKey == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	cfg.TM.TokenIssuer = tokenIssuer
	cfg.TM.AccessTokenTTL = accessTokenTTL
	cfg.TM.RefreshTokenTTL = refreshTokenTTL

	// DB Env load
	cfg.DB.Name = os.Getenv("DB_NAME")
	if cfg.DB.Name == "" {
		cfg.DB.Name = "postgres"
	}

	cfg.DB.Host = os.Getenv("DB_HOST")
	if cfg.DB.Host == "" {
		cfg.DB.Host = "db"
	}

	cfg.DB.Port = os.Getenv("DB_PORT")
	if cfg.DB.Port == "" {
		cfg.DB.Port = "5432"
	}

	cfg.DB.User = os.Getenv("DB_USER")
	if cfg.DB.User == "" {
		cfg.DB.User = "dbuser"
	}

	cfg.DB.Pass = os.Getenv("DB_PASS")
	if cfg.DB.Pass == "" {
		log.Fatal("DB_PASS environment variable is required")
	}

	// Server Env load
	cfg.HTTPServer.Address = os.Getenv("SERVER_ADDRESS")
	if cfg.HTTPServer.Address == "" {
		cfg.HTTPServer.Address = "localhost:8080"
	}

	cfg.HTTPServer.Timeout = parseDurationEnv("SERVER_TIMEOUT", defaultSrvTimeout)
	cfg.HTTPServer.IdleTimeout = parseDurationEnv("SERVER_IDLE_TIMEOUT", defaultSrvIdleTimeout)
	cfg.HTTPServer.ShutdownTimeout = parseDurationEnv("SERVER_SHUTDOWN_TIMEOUT", defaultSrvShutdownTimeout)

	// MinIO Env load
	cfg.MinIO.Endpoint = os.Getenv("MINIO_ENDPOINT")
	if cfg.MinIO.Endpoint == "" {
		cfg.MinIO.Endpoint = "minio:9000"
	}

	cfg.MinIO.PublicEndpoint = os.Getenv("MINIO_PUBLIC_ENDPOINT")
	if cfg.MinIO.PublicEndpoint == "" {
		cfg.MinIO.PublicEndpoint = "http://localhost:9000"
	}

	cfg.MinIO.AccessKey = os.Getenv("MINIO_ROOT_USER")
	if cfg.MinIO.AccessKey == "" {
		cfg.MinIO.AccessKey = "minio_admin"
	}

	cfg.MinIO.SecretKey = os.Getenv("MINIO_ROOT_PASSWORD")
	if cfg.MinIO.SecretKey == "" {
		cfg.MinIO.SecretKey = "minio_password"
	}

	cfg.MinIO.ProfilePicturesBucket = os.Getenv("MINIO_PROFILE_PICTURES_BUCKET")
	if cfg.MinIO.ProfilePicturesBucket == "" {
		cfg.MinIO.ProfilePicturesBucket = "profile_pictures"
	}

	cfg.MinIO.UseSSL = parseBoolEnv("MINIO_USE_SSL", defaultMinIOUseSSL)

	return &cfg
}

// parseDurationEnv is a helper function to parse time.Duration values from env variables.
func parseDurationEnv(envKey string, defaultVal time.Duration) time.Duration {
	valStr := os.Getenv(envKey)
	if valStr == "" {
		return defaultVal
	}

	val, err := time.ParseDuration(valStr)
	if err != nil {
		return defaultVal
	}

	return val
}

// parseBoolEnv is a helper function to parse bool env variable (e.g., "true", "false", "1", "0").
func parseBoolEnv(envKey string, defaultVal bool) bool {
	valStr := os.Getenv(envKey)
	if valStr == "" {
		return defaultVal
	}

	val, err := strconv.ParseBool(valStr)
	if err != nil {
		return defaultVal
	}

	return val
}
