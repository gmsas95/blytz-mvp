package api

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"github.com/gmsas95/blytz-mvp/services/chat-service/internal/services"
	"github.com/gmsas95/blytz-mvp/services/chat-service/internal/config"
	"github.com/gmsas95/blytz-mvp/services/chat-service/internal/api/handlers"
	"github.com/gmsas95/blytz-mvp/shared/pkg/auth"
)

func SetupRouter(logger *zap.Logger) *gin.Engine {
	// Initialize config
	cfg := config.LoadConfig()

	// Initialize chat service
	chatService := services.NewChatService(logger, cfg)

	// Create router
	router := gin.Default()

	// Initialize auth client
	authClient := auth.NewAuthClient("http://auth-service:8084")

	// Create chat handler
	chatHandler := handlers.NewChatHandler(chatService, logger)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "chat"})
	})

	// Prometheus metrics endpoint
	router.GET("/chat-metrics", gin.WrapH(promhttp.Handler()))

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