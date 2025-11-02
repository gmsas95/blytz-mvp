package api

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gmsas95/blytz-mvp/shared/pkg/middleware"
	"github.com/gmsas95/blytz-mvp/shared/pkg/utils"
	"go.uber.org/zap"
)

func SetupRouter(logger *zap.Logger) *gin.Engine {
	router := gin.Default()

	// Add correlation ID middleware for structured logging
	router.Use(utils.CorrelationMiddleware())

	// Initialize rate limiter
	rateLimiter, err := middleware.NewRateLimiter(middleware.RateLimiterConfig{
		RequestsPerMinute: 60, // 60 requests per minute per IP
		BurstSize:         10,
		RedisURL:          "redis:6379",
		Logger:            logger,
	})
	if err != nil {
		logger.Error("Failed to initialize rate limiter", zap.Error(err))
	} else {
		// Apply rate limiting to all API routes
		router.Use("/api", rateLimiter.RateLimit())
	}

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Correlation-ID")

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

	// Enhanced health check endpoint
	router.GET("/health", func(c *gin.Context) {
		correlationID := c.GetHeader("X-Correlation-ID")
		if correlationID == "" {
			correlationID = c.GetString("correlation_id")
		}

		health := gin.H{
			"status":         "ok",
			"service":        "gateway",
			"version":        "1.0.0",
			"timestamp":      time.Now().Unix(),
			"correlation_id": correlationID,
			"environment":    "production",
		}

		// Check rate limiter status
		if rateLimiter != nil {
			health["rate_limiter"] = "active"
		} else {
			health["rate_limiter"] = "inactive"
			health["status"] = "degraded"
		}

		// Check external dependencies
		health["dependencies"] = gin.H{
			"redis": "configured",
		}

		c.JSON(http.StatusOK, health)
	})

	// Simple ping endpoint (no dependencies)
	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong - updated at "+time.Now().Format("2006-01-02 15:04:05"))
	})

	// API routes with enhanced rate limiting
	api := router.Group("/api")
	{
		// Public routes with stricter rate limiting
		public := api.Group("/public")
		{
			if rateLimiter != nil {
				// Apply stricter rate limiting to sensitive endpoints
				public.Use("/livekit/token", rateLimiter.RateLimitByPath(map[string]int{
					"livekit/token": 10, // 10 requests per minute for token generation
				}))
			}

			public.GET("/health", func(c *gin.Context) {
				correlationID := c.GetString("correlation_id")
				c.JSON(200, gin.H{
					"status":         "ok",
					"service":        "gateway-public",
					"correlation_id": correlationID,
				})
			})

			public.GET("/test", func(c *gin.Context) {
				correlationID := c.GetString("correlation_id")
				c.JSON(200, gin.H{
					"message":        "test works",
					"correlation_id": correlationID,
				})
			})

			// LiveKit token generation
			public.GET("/livekit/token", createLiveKitTokenHandler(logger))

			// LiveKit proxy routes (removed to avoid conflicts)
			// public.Any("/livekit/*proxyPath", liveKitProxyHandler(logger))
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
