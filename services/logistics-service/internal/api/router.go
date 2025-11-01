package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/gmsas95/blytz-mvp/services/logistics-service/internal/api/handlers"
	"github.com/gmsas95/blytz-mvp/services/logistics-service/internal/config"
	"github.com/gmsas95/blytz-mvp/services/logistics-service/internal/services"
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

	// Initialize logistics service
	logisticsService := services.NewLogisticsService(db, logger, cfg)

	// Initialize Ninja Van service
	ninjaVanService := services.NewNinjaVanService(db, logger, cfg)

	// Create router
	router := gin.Default()

	// Initialize auth client
	authClient := auth.NewAuthClient("http://auth-service:8084")

	// Create handlers
	logisticsHandler := handlers.NewLogisticsHandler(logisticsService, logger)
	ninjaVanHandler := handlers.NewNinjaVanHandler(ninjaVanService, logger)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "logistics"})
	})

	// Prometheus metrics endpoint

	// Logistics endpoints
	logisticsRoutes := router.Group("/api/v1/logistics")
	logisticsRoutes.Use(auth.GinAuthMiddleware(authClient))
	{
		logisticsRoutes.POST("/shipments", logisticsHandler.CreateShipment)
		logisticsRoutes.GET("/shipments/:id", logisticsHandler.GetShipment)
		logisticsRoutes.PUT("/shipments/:id/status", logisticsHandler.UpdateShipmentStatus)
		logisticsRoutes.GET("/shipments/order/:orderId", logisticsHandler.GetShipmentByOrder)
		logisticsRoutes.GET("/tracking/:trackingNumber", logisticsHandler.TrackShipment)

		// Ninja Van integration endpoints
		ninjaVanRoutes := logisticsRoutes.Group("/ninjavan")
		{
			ninjaVanRoutes.POST("/shipments", ninjaVanHandler.CreateNinjaVanShipment)
			ninjaVanRoutes.POST("/shipments/:id/cancel", ninjaVanHandler.CancelNinjaVanShipment)
			ninjaVanRoutes.POST("/tariff", ninjaVanHandler.GetShippingCost)
			ninjaVanRoutes.GET("/pudo-points", ninjaVanHandler.GetPUDOPoints)
		}
	}

	// Public webhook endpoint (no auth required)
	router.POST("/api/v1/logistics/ninjavan/webhook", ninjaVanHandler.ProcessWebhook)

	return router
}
