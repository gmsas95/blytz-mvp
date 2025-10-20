package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
)

type Config struct {
	Environment      string
	ServicePort      string
	DatabaseURL      string
	PostgresUser     string `env:"POSTGRES_USER"`
	PostgresPassword string `env:"POSTGRES_PASSWORD"`
	PostgresHost     string `env:"POSTGRES_HOST"`
	PostgresPort     string `env:"POSTGRES_PORT"`
	PostgresDB       string `env:"POSTGRES_DB"`
	RedisURL         string
	RedisPassword    string
	AuthServiceURL   string
	JWTSecret        string
	LogLevel         string
}

func LoadConfig() *Config {
	cfg := &Config{
		Environment:      getEnv("NODE_ENV", "development"),
		ServicePort:      getEnv("PORT", "8085"),
		PostgresUser:     getEnv("POSTGRES_USER", "blytz"),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", ""),
		PostgresHost:     getEnv("POSTGRES_HOST", "postgres"),
		PostgresPort:     getEnv("POSTGRES_PORT", "5432"),
		PostgresDB:       getEnv("POSTGRES_DB", "blytz_prod"),
		RedisURL:         getEnv("REDIS_URL", "redis:6379"),
		RedisPassword:    getEnv("REDIS_PASSWORD", ""),
		AuthServiceURL:   getEnv("AUTH_SERVICE_URL", "http://auth-service:8084"),
		JWTSecret:        getEnv("JWT_SECRET", "your-secret-key"),
		LogLevel:         getEnv("LOG_LEVEL", "info"),
	}

	// Construct the database URL
	encodedPassword := url.QueryEscape(cfg.PostgresPassword)
	cfg.DatabaseURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.PostgresUser, encodedPassword, cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresDB)

	return cfg
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