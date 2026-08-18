package config

import (
	"os"
)

type Config struct {
	Port             string
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
		DatabaseHost:     getEnv("DATABASE_HOST", "localhost"),
		DatabasePort:     getEnv("DATABASE_PORT", "5432"),
		DatabaseUser:     getEnv("DATABASE_USER", "tasker"),
		DatabasePassword: getEnv("DATABASE_PASSWORD", "tasker_password"),
		DatabaseName:     getEnv("DATABASE_NAME", "tasker_db"),
		DatabaseSSLMode:  getEnv("DATABASE_SSLMODE", "disable"),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}
	return value
}
