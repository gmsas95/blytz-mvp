package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gmsas95/blytz-mvp/services/order-service/internal/api/handlers"
	"github.com/gmsas95/blytz-mvp/services/order-service/internal/config"
	"github.com/gmsas95/blytz-mvp/services/order-service/internal/models"
	"github.com/gmsas95/blytz-mvp/services/order-service/internal/services"
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

	// Auto-migrate database schema
	if err := db.AutoMigrate(
		&models.Order{},
		&models.OrderItem{},
		&models.Cart{},
		&models.CartItem{},
	); err != nil {
		logger.Fatal("Failed to migrate database", zap.Error(err))
	}

	// Initialize order service
	orderService := services.NewOrderService(db, logger, cfg)

	// Create router
	router := gin.Default()

	// Add correlation ID middleware for structured logging
	router.Use(utils.CorrelationMiddleware())

	// Initialize auth client
	authClient := auth.NewAuthClient("http://auth-service:8084")

	// Create order and cart handlers
	orderHandler := handlers.NewOrderHandler(orderService, logger)
	cartHandler := handlers.NewCartHandler(orderService, logger)

	// Enhanced health check endpoint
	router.GET("/health", func(c *gin.Context) {
		correlationID := c.GetHeader("X-Correlation-ID")
		if correlationID == "" {
			correlationID = c.GetString("correlation_id")
		}

		health := gin.H{
			"status":         "ok",
			"service":        "order",
			"version":        "1.0.0",
			"timestamp":      time.Now().Unix(),
			"correlation_id": correlationID,
			"environment":    cfg.Environment,
		}

		// Check database connectivity
		if orderService != nil {
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
		}

		c.JSON(http.StatusOK, health)
	})

	// Prometheus metrics endpoint

	// Order endpoints
	orderRoutes := router.Group("/api/v1/orders")
	orderRoutes.Use(auth.GinAuthMiddleware(authClient))
	{
		orderRoutes.POST("/", orderHandler.CreateOrder)
		orderRoutes.GET("/:id", orderHandler.GetOrder)
		orderRoutes.GET("/user/:userId", orderHandler.GetUserOrders)
		orderRoutes.PUT("/:id/status", orderHandler.UpdateOrderStatus)
		orderRoutes.DELETE("/:id", orderHandler.CancelOrder)
	}

	// Cart endpoints
	cartRoutes := router.Group("/api/v1/cart")
	cartRoutes.Use(auth.GinAuthMiddleware(authClient))
	{
		cartRoutes.GET("/", cartHandler.GetCart)
		cartRoutes.POST("/add", cartHandler.AddToCart)
		cartRoutes.DELETE("/remove/:itemId", cartHandler.RemoveFromCart)
		cartRoutes.PUT("/update/:itemId", cartHandler.UpdateCartItemQuantity)
		cartRoutes.DELETE("/clear", cartHandler.ClearCart)
	}

	return router
}
