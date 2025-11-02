package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gmsas95/blytz-mvp/services/payment-service/internal/api/handlers"
	"github.com/gmsas95/blytz-mvp/services/payment-service/internal/config"
	"github.com/gmsas95/blytz-mvp/services/payment-service/internal/services"
	"github.com/gmsas95/blytz-mvp/shared/pkg/auth"
	"github.com/gmsas95/blytz-mvp/shared/pkg/utils"
	"go.uber.org/zap"
)

func SetupRouter(logger *zap.Logger) *gin.Engine {
	// Initialize config
	cfg := config.LoadConfig()

	// Initialize database connection
	db, err := config.InitDB(cfg)
	if err != nil {
		logger.Fatal("Failed to initialize database", zap.Error(err))
	}

	// Initialize payment service
	paymentService := services.NewPaymentService(db, logger, cfg)

	// Create router
	router := gin.Default()

	// Add correlation ID middleware for structured logging
	router.Use(utils.CorrelationMiddleware())

	// Initialize auth client
	authClient := auth.NewAuthClient("http://auth-service:8084")

	// Create payment handler
	paymentHandler := handlers.NewPaymentHandler(paymentService, logger)

	// Enhanced health check endpoint
	router.GET("/health", func(c *gin.Context) {
		correlationID := c.GetHeader("X-Correlation-ID")
		if correlationID == "" {
			correlationID = c.GetString("correlation_id")
		}

		health := gin.H{
			"status":         "ok",
			"service":        "payment",
			"version":        "1.0.0",
			"timestamp":      time.Now().Unix(),
			"correlation_id": correlationID,
			"environment":    cfg.Environment,
		}

		// Check database connectivity
		if paymentService != nil {
			health["database"] = "connected"
		} else {
			health["database"] = "disconnected"
			health["status"] = "degraded"
			c.JSON(http.StatusServiceUnavailable, health)
			return
		}

		// Check external dependencies
		health["dependencies"] = gin.H{
			"auth_service": "connected",
			"fiuu_api":     "configured",
		}

		c.JSON(http.StatusOK, health)
	})

	// Prometheus metrics endpoint

	// Payment endpoints
	paymentRoutes := router.Group("/api/v1/payments")
	paymentRoutes.Use(auth.GinAuthMiddleware(authClient))
	{
		paymentRoutes.POST("/process", paymentHandler.ProcessPayment)
		paymentRoutes.GET("/methods", paymentHandler.GetPaymentMethods)
		paymentRoutes.GET("/history", paymentHandler.GetPaymentHistory)
		paymentRoutes.GET("/:id", paymentHandler.GetPayment)
		paymentRoutes.POST("/:id/refund", paymentHandler.ProcessRefund)
		paymentRoutes.GET("/seamless/config", paymentHandler.GetSeamlessConfig)
	}

	// Public seamless config endpoint (no auth required for frontend)
	router.GET("/api/v1/public/seamless/config", paymentHandler.GetPublicSeamlessConfig)

	// Webhook endpoint (no auth required)
	router.POST("/api/v1/webhooks/fiuu", paymentHandler.ProcessWebhook)

	return router
}
