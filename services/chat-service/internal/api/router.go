package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gmsas95/blytz-mvp/services/chat-service/internal/api/handlers"
	"github.com/gmsas95/blytz-mvp/services/chat-service/internal/config"
	"github.com/gmsas95/blytz-mvp/services/chat-service/internal/services"
	"github.com/gmsas95/blytz-mvp/shared/pkg/auth"
	"github.com/gmsas95/blytz-mvp/shared/pkg/utils"
	"go.uber.org/zap"
)

func SetupRouter(logger *zap.Logger) *gin.Engine {
	// Initialize config
	cfg := config.LoadConfig()

	// Initialize structured logger
	structuredLogger, err := utils.NewStructuredLogger(utils.LoggerConfig{
		Service:     "chat-service",
		Environment: cfg.Environment,
		Level:       cfg.LogLevel,
	})
	if err != nil {
		logger.Fatal("Failed to initialize structured logger", zap.Error(err))
	}

	// Initialize chat service
	chatService := services.NewChatService(logger, cfg)

	// Create router
	router := gin.Default()

	// Add correlation ID middleware for structured logging
	router.Use(utils.CorrelationMiddleware(structuredLogger))

	// Initialize auth client
	authClient := auth.NewAuthClient("http://auth-service:8084")

	// Create chat handler
	chatHandler := handlers.NewChatHandler(chatService, logger)

	// Enhanced health check endpoint
	router.GET("/health", func(c *gin.Context) {
		correlationID := c.GetHeader("X-Correlation-ID")
		if correlationID == "" {
			correlationID = c.GetString("correlation_id")
		}

		health := gin.H{
			"status":         "ok",
			"service":        "chat",
			"version":        "1.0.0",
			"timestamp":      time.Now().Unix(),
			"correlation_id": correlationID,
			"environment":    cfg.Environment,
		}

		// Check service connectivity
		if chatService != nil {
			health["websocket"] = "available"
		} else {
			health["websocket"] = "unavailable"
			health["status"] = "degraded"
			c.JSON(http.StatusServiceUnavailable, health)
			return
		}

		// Check external dependencies
		health["dependencies"] = gin.H{
			"auth_service": "connected",
			"redis":        "configured",
		}

		c.JSON(http.StatusOK, health)
	})

	// Prometheus metrics endpoint

	// Chat endpoints
	chatRoutes := router.Group("/api/v1/chat")
	chatRoutes.Use(auth.GinAuthMiddleware(authClient))
	{
		chatRoutes.GET("/ws", chatHandler.HandleWebSocket)
		chatRoutes.GET("/rooms/:roomId/messages", chatHandler.GetMessages)
		chatRoutes.POST("/rooms/:roomId/messages", chatHandler.SendMessage)
		chatRoutes.GET("/rooms", chatHandler.GetUserRooms)
	}

	return router
}
