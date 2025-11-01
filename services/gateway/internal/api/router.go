package api

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

func SetupRouter(logger *zap.Logger) *gin.Engine {
	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.Status(http.StatusOK)
			return
		}
		c.Next()
	})

	// Root endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": "Blytz API Gateway",
			"version": "v1",
			"status":  "running",
			"endpoints": map[string]string{
				"health":        "/health",
				"public_api":    "/api/public/",
				"livekit_token": "/api/public/livekit/token",
			},
		})
	})

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "gateway", "timestamp": time.Now().UTC()})
	})

	// Simple ping endpoint (no dependencies)
	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	// API routes
	api := router.Group("/api")
	{
		// Public routes
		public := api.Group("/public")
		{
			public.GET("/health", func(c *gin.Context) {
				c.JSON(200, gin.H{"status": "ok"})
			})

			public.GET("/test", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "test works"})
			})

			// LiveKit token generation
			public.GET("/livekit/token", createLiveKitTokenHandler(logger))
		}
	}

	return router
}

func createLiveKitTokenHandler(logger *zap.Logger) gin.HandlerFunc {
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
		token, err := createLiveKitToken(room, role, identity, name)
		if err != nil {
			logger.Error("Failed to create LiveKit token", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to generate token",
			})
			return
		}

		// Get LiveKit URL from environment
		livekitURL := getEnv("LIVEKIT_URL", "wss://blytz-live-u5u72ozx.livekit.cloud")

		response := gin.H{
			"token":    token,
			"url":      livekitURL,
			"room":     room,
			"identity": identity,
		}

		c.JSON(http.StatusOK, response)
	}
}

func createLiveKitToken(room, role, identity, name string) (string, error) {
	// Get LiveKit credentials from environment
	apiKey := getEnv("LIVEKIT_API_KEY", "")
	apiSecret := getEnv("LIVEKIT_API_SECRET", "")

	if apiKey == "" || apiSecret == "" {
		return "", fmt.Errorf("LiveKit API key and secret are required")
	}

	// Set role-specific permissions
	videoClaims := map[string]interface{}{
		"room":     room,
		"roomJoin": true,
	}

	switch role {
	case "host", "broadcaster":
		videoClaims["roomAdmin"] = true
		videoClaims["canPublish"] = true
		videoClaims["canPublishData"] = true
		videoClaims["canSubscribe"] = true
	default: // viewer
		videoClaims["canPublish"] = false
		videoClaims["canPublishData"] = false
		videoClaims["canSubscribe"] = true
	}

	// Create JWT claims
	claims := jwt.MapClaims{
		"iss":      apiKey,
		"sub":      identity,
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(6 * time.Hour).Unix(), // 6 hours
		"video":    videoClaims,
		"metadata": fmt.Sprintf(`{"role":"%s","name":"%s","room":"%s"}`, role, name, room),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with LiveKit secret
	tokenString, err := token.SignedString([]byte(apiSecret))
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
