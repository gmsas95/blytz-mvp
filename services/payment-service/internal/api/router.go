package api

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"github.com/gmsas95/blytz-mvp/services/payment-service/internal/services"
	"github.com/gmsas95/blytz-mvp/services/payment-service/internal/config"
	"github.com/gmsas95/blytz-mvp/services/payment-service/internal/api/handlers"
	"github.com/gmsas95/blytz-mvp/shared/pkg/auth"
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

	// Initialize auth client
	authClient := auth.NewAuthClient("http://auth-service:8084")

	// Create payment handler
	paymentHandler := handlers.NewPaymentHandler(paymentService, logger)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "payment"})
	})

	// Prometheus metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Payment endpoints
	paymentRoutes := router.Group("/api/v1/payments")
	paymentRoutes.Use(auth.GinAuthMiddleware(authClient))
	{
		paymentRoutes.POST("/process", paymentHandler.ProcessPayment)
		paymentRoutes.GET("/methods", paymentHandler.GetPaymentMethods)
		paymentRoutes.GET("/history", paymentHandler.GetPaymentHistory)
		paymentRoutes.GET("/:id", paymentHandler.GetPayment)
		paymentRoutes.POST("/:id/refund", paymentHandler.ProcessRefund)
	}

	return router
}