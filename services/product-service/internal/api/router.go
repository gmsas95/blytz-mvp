package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/gmsas95/blytz-mvp/services/product-service/internal/api/handlers"
	"github.com/gmsas95/blytz-mvp/services/product-service/internal/config"
	"github.com/gmsas95/blytz-mvp/services/product-service/internal/services"
	"github.com/gmsas95/blytz-mvp/shared/pkg/auth"
)

func SetupRouter(db *gorm.DB, logger *zap.Logger, cfg *config.Config) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Initialize auth client
	authClient := auth.NewAuthClient("http://auth-service:8084")

	// Initialize product service
	productService := services.NewProductService(db, logger, cfg)

	// Initialize product handler
	productHandler := handlers.NewProductHandler(productService, logger)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "ok",
			"service":   "product",
			"timestamp": time.Now().Unix(),
		})
	})

	// Metrics
	router.GET("/product-metrics", gin.WrapH(promhttp.Handler()))

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
