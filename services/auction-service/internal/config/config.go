package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port                   string
	Environment            string
	LogLevel               string
	DatabaseURL            string
	RedisURL               string
	JWTSecret              string
	MetricsPort            string
	ServiceName            string
	DefaultAuctionDuration time.Duration
}

func Load() (*Config, error) {
	return &Config{
		Port:            getEnv("PORT", "8083"),
		Environment:     getEnv("ENVIRONMENT", "development"),
		LogLevel:        getEnv("LOG_LEVEL", "info"),
		DatabaseURL:     getEnv("DATABASE_URL", "postgres://postgres:password@localhost:5432/auction_db?sslmode=disable"),
		RedisURL:        getEnv("REDIS_URL", "localhost:6379"),
		JWTSecret:       getEnv("JWT_SECRET", "your-secret-key"),
		MetricsPort:     getEnv("METRICS_PORT", "9083"),
		ServiceName:            getEnv("SERVICE_NAME", "auction-service"),
		DefaultAuctionDuration: getEnvAsDuration("DEFAULT_AUCTION_DURATION", 24*time.Hour),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}