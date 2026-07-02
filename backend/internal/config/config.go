package config

import (
	"log"
	"os"
	"time"
)

type Config struct {
	Env        string
	JWTSecret  string
	DB         Database
	HTTPServer HTTPServer
}

type Database struct {
	Name string
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

func MustLoad() *Config {
	/*
		ENVIRONMENT=dev
		JWT_SECRET=super-secure-random-string

		# DB variables
		DB_NAME=postgres
		DB_PORT=5432
		DB_USER=dbuser
		DB_PASS=my_db_password

		# Server variables
		SERVER_ADDRESS=localhost:8080
		# Timeout in seconds
		SERVER_TIMEOUT=4
		SERVER_IDLE_TIMEOUT=60
		SERVER_SHUTDOWN_TIMEOUT=5

	*/

	var cfg Config

	cfg.Env = os.Getenv("ENVIRONMENT")
	if cfg.Env == "" {
		cfg.Env = "dev"
	}

	cfg.JWTSecret = os.Getenv("JWT_SECRET")
	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET enviroment variable is required")
	}

	// DB Env load
	cfg.DB.Name = os.Getenv("DB_NAME")
	if cfg.DB.Name == "" {
		cfg.DB.Name = "postgres"
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
		log.Fatal("DB_PASS enviroment variable is required")
	}

	// Server Env load
	// SERVER_ADDRESS=localhost:8080
	// SERVER_TIMEOUT=4
	// 	SERVER_IDLE_TIMEOUT=60
	// 	SERVER_SHUTDOWN_TIMEOUT=5

	cfg.HTTPServer.Address = os.Getenv("SERVER_ADDRESS")
	if cfg.HTTPServer.Address == "" {
		cfg.HTTPServer.Address = "localhost:8080"
	}

	srvTimeoutStr := os.Getenv("SERVER_TIMEOUT")
	srvTimeout, err := time.ParseDuration(srvTimeoutStr)
	if err != nil {
		cfg.HTTPServer.Timeout = 4 * time.Second
	} else {
		cfg.HTTPServer.Timeout = srvTimeout
	}

	return &cfg
}
