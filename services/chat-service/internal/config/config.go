package config

import (
	"os"
)

type Config struct {
	Environment    string
	ServicePort    string
	RedisURL       string
	RedisPassword  string
	AuthServiceURL string
	JWTSecret      string
	LogLevel       string
}

func LoadConfig() *Config {
	return &Config{
		Environment:    getEnv("ENVIRONMENT", "development"),
		ServicePort:    getEnv("PORT", "8088"),
		RedisURL:       getEnv("REDIS_URL", "redis:6379"),
		RedisPassword:  getEnv("REDIS_PASSWORD", ""),
		AuthServiceURL: getEnv("AUTH_SERVICE_URL", "http://auth-service:8084"),
		JWTSecret:      getEnv("JWT_SECRET", "your-secret-key"),
		LogLevel:       getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}