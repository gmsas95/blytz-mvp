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

	// Check if DATABASE_URL is provided (Dokploy style)
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		cfg.DatabaseURL = databaseURL
		// Parse DATABASE_URL to extract components for fallback usage
		if parsedURL, err := url.Parse(databaseURL); err == nil {
			// Extract user info
			if parsedURL.User != nil {
				cfg.PostgresUser = parsedURL.User.Username()
				if password, ok := parsedURL.User.Password(); ok {
					cfg.PostgresPassword = password
				}
			}
			// Extract host and port
			cfg.PostgresHost = parsedURL.Hostname()
			if port := parsedURL.Port(); port != "" {
				cfg.PostgresPort = port
			}
			// Extract database name
			if len(parsedURL.Path) > 1 {
				cfg.PostgresDB = parsedURL.Path[1:] // Remove leading slash
			}
		}
	} else {
		// Construct the database URL from individual components
		encodedPassword := url.QueryEscape(cfg.PostgresPassword)
		cfg.DatabaseURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			cfg.PostgresUser, encodedPassword, cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresDB)
	}

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