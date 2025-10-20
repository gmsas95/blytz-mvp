package config

import (
	"os"
	"strconv"
)

type Config struct {
	Environment      string
	ServicePort      string
	DatabaseURL      string
	RedisURL         string
	RedisPassword    string
	AuthServiceURL   string
	JWTSecret        string
	LogLevel         string
	StripeSecretKey  string
	StripePublishableKey string
}

func LoadConfig() *Config {
	return &Config{
		Environment:      getEnv("ENVIRONMENT", "development"),
		ServicePort:      getEnv("PORT", "8086"),
		DatabaseURL:      getEnv("DATABASE_URL", "postgres://postgres:password@postgres:5432/payments?sslmode=disable"),
		RedisURL:         getEnv("REDIS_URL", "redis:6379"),
		RedisPassword:    getEnv("REDIS_PASSWORD", ""),
		AuthServiceURL:   getEnv("AUTH_SERVICE_URL", "http://auth-service:8084"),
		JWTSecret:        getEnv("JWT_SECRET", "your-secret-key"),
		LogLevel:         getEnv("LOG_LEVEL", "info"),
		StripeSecretKey:  getEnv("STRIPE_SECRET_KEY", "sk_test_..."),
		StripePublishableKey: getEnv("STRIPE_PUBLISHABLE_KEY", "pk_test_..."),
	}
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