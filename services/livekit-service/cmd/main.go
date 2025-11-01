package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gmsas95/blytz-mvp/shared/pkg/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/livekit/livekit-go"
	"go.uber.org/zap"
)

type Config struct {
	Port          string
	LiveKitAPIKey string
	LiveKitSecret string
	LiveKitURL    string
	AuthClient    *auth.AuthClient
	Logger        *zap.Logger
}

type TokenRequest struct {
	Room string `json:"room" binding:"required"`
	Role string `json:"role"`
	Name string `json:"name"`
}

type TokenResponse struct {
	Token    string `json:"token"`
	URL      string `json:"url"`
	Room     string `json:"room"`
	Identity string `json:"identity"`
}

func main() {
	// Initialize logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Load configuration
	config := &Config{
		Port:          getEnv("PORT", "8089"),
		LiveKitAPIKey: getEnv("LIVEKIT_API_KEY", ""),
		LiveKitSecret: getEnv("LIVEKIT_API_SECRET", ""),
		LiveKitURL:    getEnv("LIVEKIT_URL", "https://livekit.blytz.app"),
		Logger:        logger,
	}

	// Initialize auth client
	authClient := auth.NewAuthClient(getEnv("AUTH_SERVICE_URL", "http://auth-service:8084"))
	config.AuthClient = authClient

	// Validate configuration
	if config.LiveKitAPIKey == "" || config.LiveKitSecret == "" {
		logger.Fatal("LiveKit API key and secret are required")
	}

	// Initialize Gin router
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Enable CORS
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// Routes
	api := router.Group("/api/v1")
	{
		// Public endpoints
		api.GET("/health", healthHandler(config))
		api.GET("/livekit/token", createTokenHandler(config))
	}

	// Legacy endpoint for compatibility
	router.GET("/api/livekit/token", createTokenHandler(config))

	logger.Info("LiveKit service starting", zap.String("port", config.Port))

	// Start server
	if err := router.Run(":" + config.Port); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}

func healthHandler(config *Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"service":   "livekit-service",
			"timestamp": time.Now().UTC(),
			"version":   "1.0.0",
		})
	}
}

func createTokenHandler(config *Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get query parameters
		room := c.Query("room")
		role := c.DefaultQuery("role", "viewer")
		name := c.Query("name")

		if room == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Room name is required",
			})
			return
		}

		// Validate role
		if role != "viewer" && role != "host" && role != "broadcaster" {
			role = "viewer"
		}

		// Generate unique identity
		identity := fmt.Sprintf("%s_%d_%d", role, time.Now().Unix(), time.Now().Nanosecond())

		// Set default name if not provided
		if name == "" {
			name = identity
		}

		// Create LiveKit token
		token, err := createLiveKitToken(config, room, role, identity, name)
		if err != nil {
			config.Logger.Error("Failed to create LiveKit token", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to generate token",
			})
			return
		}

		response := TokenResponse{
			Token:    token,
			URL:      config.LiveKitURL,
			Room:     room,
			Identity: identity,
		}

		c.JSON(http.StatusOK, response)
	}
}

func createLiveKitToken(config *Config, room, role, identity, name string) (string, error) {
	// Create JWT claims
	claims := &jwt.MapClaims{
		"iss": config.LiveKitAPIKey,
		"nbf": time.Now().Unix(),
		"exp": time.Now().Add(6 * time.Hour).Unix(), // 6 hours
		"sub": identity,
		"jti": fmt.Sprintf("%s_%d", identity, time.Now().Unix()),
		"video": map[string]interface{}{
			"room":     room,
			"roomJoin": true,
		},
		"metadata": fmt.Sprintf(`{"role":"%s","name":"%s","room":"%s"}`, role, name, room),
	}

	// Set role-specific permissions
	switch role {
	case "host", "broadcaster":
		claims["video"].(map[string]interface{})["roomAdmin"] = true
		claims["video"].(map[string]interface{})["canPublish"] = true
		claims["video"].(map[string]interface{})["canPublishData"] = true
		claims["video"].(map[string]interface{})["canSubscribe"] = true
	default: // viewer
		claims["video"].(map[string]interface{})["canPublish"] = false
		claims["video"].(map[string]interface{})["canPublishData"] = false
		claims["video"].(map[string]interface{})["canSubscribe"] = true
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with LiveKit secret
	tokenString, err := token.SignedString([]byte(config.LiveKitSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
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

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
