package api

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"time"
)

func SetupRouter(logger *zap.Logger) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"service": "product",
			"timestamp": time.Now().Unix(),
		})
	})

	// Metrics
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// API routes
	api := router.Group("/api/v1")
	{
		SetupRoutes(api, logger)
	}

	return router
}

func SetupRoutes(api *gin.RouterGroup, logger *zap.Logger) {
	// Service-specific routes based on OpenAPI spec
	api.GET("/products", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "List products endpoint",
			"service": "product",
		})
	})

	api.POST("/products", func(c *gin.Context) {
		c.JSON(201, gin.H{
			"message": "Create product endpoint",
			"service": "product",
		})
	})

	api.GET("/products/:id", func(c *gin.Context) {
		productID := c.Param("id")
		c.JSON(200, gin.H{
			"message": "Get product endpoint",
			"product_id": productID,
			"service": "product",
		})
	})

	api.PUT("/products/:id", func(c *gin.Context) {
		productID := c.Param("id")
		c.JSON(200, gin.H{
			"message": "Update product endpoint",
			"product_id": productID,
			"service": "product",
		})
	})

	api.DELETE("/products/:id", func(c *gin.Context) {
		productID := c.Param("id")
		c.JSON(200, gin.H{
			"message": "Delete product endpoint",
			"product_id": productID,
			"service": "product",
		})
	})

	api.GET("/products/:id/inventory", func(c *gin.Context) {
		productID := c.Param("id")
		c.JSON(200, gin.H{
			"message": "Get product inventory endpoint",
			"product_id": productID,
			"service": "product",
		})
	})
}