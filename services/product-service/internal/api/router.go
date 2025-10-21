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
	// Specific routes must come before parameterized routes

	// Featured products endpoint (must be before /products/:id)
	api.GET("/products/featured", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"success": true,
			"data": []gin.H{
				{
					"id": "1",
					"name": "Vintage Leather Jacket",
					"description": "Authentic vintage leather jacket in excellent condition",
					"price": 299.99,
					"image": "https://images.unsplash.com/photo-1551028719-00167b16eac5?w=400",
					"category": "Fashion",
					"seller_id": "seller1",
					"status": "active",
					"created_at": "2024-01-15T10:00:00Z",
					"updated_at": "2024-01-15T10:00:00Z",
				},
				{
					"id": "2",
					"name": "Wireless Headphones Pro",
					"description": "Premium noise-canceling wireless headphones",
					"price": 199.99,
					"image": "https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=400",
					"category": "Electronics",
					"seller_id": "seller2",
					"status": "active",
					"created_at": "2024-01-14T15:30:00Z",
					"updated_at": "2024-01-14T15:30:00Z",
				},
				{
					"id": "3",
					"name": "Smart Watch Ultra",
					"description": "Latest smartwatch with health tracking and GPS",
					"price": 399.99,
					"image": "https://images.unsplash.com/photo-1523275335684-37898b6baf30?w=400",
					"category": "Electronics",
					"seller_id": "seller3",
					"status": "active",
					"created_at": "2024-01-13T09:15:00Z",
					"updated_at": "2024-01-13T09:15:00Z",
				},
			},
		})
	})

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