package config

import (
	"log/slog"
	"os"
	"strings"
)

type Config struct {
	Port             string
	LogLevel        slog.Level
	DatabaseHost     string
	DatabasePort     string
	DatabaseUser     string
	DatabasePassword string
	DatabaseName     string
	DatabaseSSLMode  string
}

func Load() Config {
	return Config{
		Port:             getEnv("PORT", "8080"),
		LogLevel: parseLogLevel(getEnv("LOG_LEVEL", "info")),
		DatabaseHost:     getEnv("DATABASE_HOST", "localhost"),
		DatabasePort:     getEnv("DATABASE_PORT", "5432"),
		DatabaseUser:     getEnv("DATABASE_USER", "tasker"),
		DatabasePassword: getEnv("DATABASE_PASSWORD", "tasker_password"),
		DatabaseName:     getEnv("DATABASE_NAME", "tasker_db"),
		DatabaseSSLMode:  getEnv("DATABASE_SSLMODE", "disable"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func parseLogLevel(raw string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "debug":
		return slog.LevelDebug
		case "warn", "warning":
			return slog.LevelWarn
	case "error":
			return slog.LevelError
	default:
			return slog.LevelInfo
	}
}
