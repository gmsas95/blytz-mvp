package config

import (
	"os"
)

type Config struct {
	DatabaseURL      string
	RedisURL         string
	StripeSecretKey  string
	StripePublishKey string
	StripeWebhookKey string
	Environment      string
	Port             string
	JWTSecret        string
}

func Load() *Config {
	return &Config{
		DatabaseURL:      getEnv("DATABASE_URL", "postgres://postgres:password@postgres:5432/payments?sslmode=disable"),
		RedisURL:         getEnv("REDIS_URL", "redis:6379"),
		StripeSecretKey:  getEnv("STRIPE_SECRET_KEY", ""),
		StripePublishKey: getEnv("STRIPE_PUBLISH_KEY", ""),
		StripeWebhookKey: getEnv("STRIPE_WEBHOOK_SECRET", ""),
		Environment:      getEnv("GO_ENV", "development"),
		Port:             getEnv("PORT", "8082"),
		JWTSecret:        getEnv("JWT_SECRET", "your-jwt-secret-key"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}