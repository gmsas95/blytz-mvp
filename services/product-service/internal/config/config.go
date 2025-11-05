package config

import (
	"fmt"
	"net/url"
	"os"
)

type Config struct {
	Port             string
	Environment      string
	LogLevel         string
	DatabaseURL      string
	PostgresUser     string `env:"POSTGRES_USER"`
	PostgresPassword string `env:"POSTGRES_PASSWORD"`
	PostgresHost     string `env:"POSTGRES_HOST"`
	PostgresPort     string `env:"POSTGRES_PORT"`
	PostgresDB       string `env:"POSTGRES_DB"`
	RedisURL         string
}

func Load() *Config {
	cfg := &Config{
		Port:             getEnv("PORT", "8082"),
		Environment:      getEnv("NODE_ENV", "development"),
		LogLevel:         getEnv("LOG_LEVEL", "info"),
		PostgresUser:     getEnv("POSTGRES_USER", "blytz"),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", ""),
		PostgresHost:     getEnv("POSTGRES_HOST", "postgres"),
		PostgresPort:     getEnv("POSTGRES_PORT", "5432"),
		PostgresDB:       getEnv("POSTGRES_DB", "blytz_prod"),
		RedisURL:         getEnv("REDIS_URL", "redis://localhost:6379"),
	}

	// Check if DATABASE_URL is provided (Dokploy style)
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		cfg.DatabaseURL = databaseURL
		// Parse the DATABASE_URL to extract components for fallback usage
		if parsedURL, err := url.Parse(databaseURL); err == nil {
			if parsedURL.User != nil {
				cfg.PostgresUser = parsedURL.User.Username()
				if password, ok := parsedURL.User.Password(); ok {
					cfg.PostgresPassword = password
				}
			}
			if parsedURL.Hostname() != "" {
				cfg.PostgresHost = parsedURL.Hostname()
			}
			if parsedURL.Port() != "" {
				cfg.PostgresPort = parsedURL.Port()
			}
			// Extract database name from path
			if len(parsedURL.Path) > 1 {
				cfg.PostgresDB = parsedURL.Path[1:] // Remove leading slash
			}
		}
	} else {
		// Construct the database URL from individual components (original behavior)
		encodedPassword := url.QueryEscape(cfg.PostgresPassword)
		cfg.DatabaseURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			cfg.PostgresUser, encodedPassword, cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresDB)
	}

	return cfg
}

	// Check if DATABASE_URL is provided (Dokploy style)
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		cfg.DatabaseURL = databaseURL
		// Parse the DATABASE_URL to extract components for fallback usage
		if parsedURL, err := url.Parse(databaseURL); err == nil {
			if parsedURL.User != nil {
				cfg.PostgresUser = parsedURL.User.Username()
				if password, ok := parsedURL.User.Password(); ok {
					cfg.PostgresPassword = password
				}
			}
			if parsedURL.Hostname() != "" {
				cfg.PostgresHost = parsedURL.Hostname()
			}
			if parsedURL.Port() != "" {
				cfg.PostgresPort = parsedURL.Port()
			}
			// Extract database name from path
			if len(parsedURL.Path) > 1 {
				cfg.PostgresDB = parsedURL.Path[1:] // Remove leading slash
			}
		}
	} else {
		// Construct the database URL from individual components (original behavior)
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