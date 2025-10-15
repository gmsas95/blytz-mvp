package config

import (
	"os"
)

type Config struct {
	RedisURL            string
	LiveKitAPIKey       string
	LiveKitAPISecret    string
	FirestoreProjectID  string
	Environment         string
	Port                string
}

func Load() *Config {
	return &Config{
		RedisURL:            getEnv("REDIS_URL", "redis:6379"),
		LiveKitAPIKey:       getEnv("LIVEKIT_API_KEY", ""),
		LiveKitAPISecret:    getEnv("LIVEKIT_API_SECRET", ""),
		FirestoreProjectID:  getEnv("FIRESTORE_PROJECT_ID", ""),
		Environment:         getEnv("GO_ENV", "development"),
		Port:                getEnv("PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}