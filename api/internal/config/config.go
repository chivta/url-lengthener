package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	DatabaseURL    string
	RedisURL       string
	Port           string
	AllowedOrigins string
	Environment    string
	LogLevel       string
}

func Load() (*Config, error) {
	cfg := &Config{
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		RedisURL:       os.Getenv("REDIS_URL"),
		Port:           getEnvOrDefault("PORT", "8080"),
		AllowedOrigins: os.Getenv("ALLOWED_ORIGINS"),
		Environment:    getEnvOrDefault("ENVIRONMENT", "development"),
		LogLevel:       getEnvOrDefault("LOG_LEVEL", "info"),
	}

	var missing []string
	if cfg.DatabaseURL == "" {
		missing = append(missing, "DATABASE_URL")
	}
	if cfg.RedisURL == "" {
		missing = append(missing, "REDIS_URL")
	}
	if cfg.AllowedOrigins == "" {
		missing = append(missing, "ALLOWED_ORIGINS")
	}
	if len(missing) > 0 {
		return nil, fmt.Errorf("missing required env vars: %s", strings.Join(missing, ", "))
	}

	return cfg, nil
}

func getEnvOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
