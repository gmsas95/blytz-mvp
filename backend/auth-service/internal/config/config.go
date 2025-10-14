package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port        string
	Environment string
	LogLevel    string
	DatabaseURL string
	RedisURL    string
	JWTSecret   string
	JWTExpiry   time.Duration
	Issuer      string
	Audience    string
}

func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8081"),
		Environment: getEnv("GO_ENV", "development"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:password@localhost:5432/auth?sslmode=disable"),
		RedisURL:    getEnv("REDIS_URL", "redis://localhost:6379"),
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-key"),
		JWTExpiry:   getDurationEnv("JWT_EXPIRY", 24*time.Hour),
		Issuer:      getEnv("JWT_ISSUER", "blytz.auth"),
		Audience:    getEnv("JWT_AUDIENCE", "blytz.app"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}