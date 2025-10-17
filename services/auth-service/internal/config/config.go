package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the auth service
type Config struct {
	DatabaseURL      string `env:"DATABASE_URL"`
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
		DatabaseURL:      getEnvOrDefault("DATABASE_URL", "postgres://user:pass@localhost:5432/authdb"),
		BetterAuthSecret: getEnvOrDefault("BETTER_AUTH_SECRET", "better-auth-secret-key-change-in-production"),
		JWTSecret:        getEnvOrDefault("JWT_SECRET", "jwt-secret-key-change-in-production"),
		ServicePort:      getEnvOrDefault("PORT", "8084"),
		Environment:      getEnvOrDefault("ENVIRONMENT", "development"),
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
