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

		if room == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Room name is required",
			})
			return
		}

		// For now, return a mock response to test connectivity
		response := gin.H{
			"token":    "mock_token_for_testing",
			"url":      "wss://blytz-live-u5u72ozx.livekit.cloud",
			"room":     room,
			"identity": fmt.Sprintf("%s_%d", role, time.Now().Unix()),
			"message":  "This is a mock token for testing - LiveKit integration needs environment variables",
		}

		c.JSON(http.StatusOK, response)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
