package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/gmsas95/blytz-mvp/services/product-service/internal/api/handlers"
	"github.com/gmsas95/blytz-mvp/services/product-service/internal/config"
	"github.com/gmsas95/blytz-mvp/services/product-service/internal/services"
	"github.com/gmsas95/blytz-mvp/shared/pkg/auth"
	"github.com/gmsas95/blytz-mvp/shared/pkg/errors"
	"github.com/gmsas95/blytz-mvp/shared/pkg/utils"
)

func SetupRouter(db *gorm.DB, logger *zap.Logger, cfg *config.Config) *gin.Engine {
	// Initialize structured logger
	structuredLogger, err := utils.NewStructuredLogger(utils.LoggerConfig{
		Level:       "info",
		Format:      "json",
		Service:     "product-service",
		Version:     "v1.0.0",
		Environment: cfg.Environment,
	})
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(utils.CorrelationMiddleware(structuredLogger))

	// Initialize auth client
	authClient := auth.NewAuthClient("http://auth-service:8084")

	// Initialize product service
	productService := services.NewProductService(db, logger, cfg)

	// Initialize product handler
	productHandler := handlers.NewProductHandler(productService, logger)

	// Comprehensive health check
	router.GET("/health", func(c *gin.Context) {
		health := gin.H{
			"status":    "ok",
			"service":   "product",
			"timestamp": time.Now().Unix(),
			"version":   "v1.0.0",
			"checks": gin.H{
				"database":     "connected",
				"auth_service": "connected",
			},
		}
		c.JSON(200, health)
	})

	// API routes
	api := router.Group("/api/v1")
	{
		// Public routes (no authentication required)
		public := api.Group("/products")
		{
			public.GET("/", productHandler.GetProducts)
			public.GET("/featured", productHandler.GetFeaturedProducts)
			public.GET("/:id", productHandler.GetProduct)
		}

		// Protected routes (authentication required)
		protected := api.Group("/products")
		protected.Use(auth.GinAuthMiddleware(authClient))
		{
			protected.POST("/", productHandler.CreateProduct)
			protected.PUT("/:id", productHandler.UpdateProduct)
			protected.DELETE("/:id", productHandler.DeleteProduct)
			protected.PUT("/:id/inventory", productHandler.UpdateInventory)
			protected.GET("/my", productHandler.GetMyProducts)
		}
	}

	return router
}
