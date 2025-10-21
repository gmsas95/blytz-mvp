package config

import (
	"fmt"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the auth service
type Config struct {
	DatabaseURL      string
	PostgresUser     string `env:"POSTGRES_USER"`
	PostgresPassword string `env:"POSTGRES_PASSWORD"`
	PostgresHost     string `env:"POSTGRES_HOST"`
	PostgresPort     string `env:"POSTGRES_PORT"`
	PostgresDB       string `env:"POSTGRES_DB"`
	BetterAuthSecret string `env:"BETTER_AUTH_SECRET"`
	JWTSecret        string `env:"JWT_SECRET"`
	ServicePort      string `env:"PORT"`
	Environment      string `env:"ENVIRONMENT"`
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists (for development)
	_ = godotenv.Load()

	cfg := &Config{
		PostgresUser:     getEnvOrDefault("POSTGRES_USER", "blytz"),
		PostgresPassword: getEnvOrDefault("POSTGRES_PASSWORD", ""),
		PostgresHost:     getEnvOrDefault("POSTGRES_HOST", "postgres"),
		PostgresPort:     getEnvOrDefault("POSTGRES_PORT", "5432"),
		PostgresDB:       getEnvOrDefault("POSTGRES_DB", "blytz_prod"),
		BetterAuthSecret: getEnvOrDefault("BETTER_AUTH_SECRET", "better-auth-secret-key-change-in-production"),
		JWTSecret:        getEnvOrDefault("JWT_SECRET", "jwt-secret-key-change-in-production"),
		ServicePort:      getEnvOrDefault("PORT", "8084"),
		Environment:      getEnvOrDefault("NODE_ENV", "development"),
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

	return cfg, nil
}

// getEnvOrDefault gets environment variable or returns default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}
