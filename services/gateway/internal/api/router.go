package api

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"github.com/gmsas95/blytz-mvp/shared/pkg/auth"
)

func SetupRouter(logger *zap.Logger) *gin.Engine {
	router := gin.Default()

	// Initialize auth client
	authClient := auth.NewAuthClient("http://auth-service:8084")

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "gateway"})
	})

	// Prometheus metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// API routes
	api := router.Group("/api")
	{
		// Public routes
		public := api.Group("/public")
		{
			public.GET("/health", func(c *gin.Context) {
				c.JSON(200, gin.H{"status": "ok"})
			})
		}

		// Protected routes
		protected := api.Group("/v1")
		protected.Use(auth.GinAuthMiddleware(authClient))
		{
			// Proxy to other services
			protected.Any("/auctions/*proxyPath", proxyToService("http://auction-service:8083"))
			protected.Any("/products/*proxyPath", proxyToService("http://product-service:8082"))
			protected.Any("/orders/*proxyPath", proxyToService("http://order-service:8085"))
			protected.Any("/payments/*proxyPath", proxyToService("http://payment-service:8086"))
			protected.Any("/logistics/*proxyPath", proxyToService("http://logistics-service:8087"))
			protected.Any("/chat/*proxyPath", proxyToService("http://chat-service:8088"))
		}
	}

	return router
}

func proxyToService(targetURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Simple proxy implementation
		proxyPath := c.Param("proxyPath")
		fullURL := targetURL + proxyPath

		// For now, just return a placeholder response
		c.JSON(200, gin.H{
			"message": "Gateway proxy",
			"target":  fullURL,
			"path":    proxyPath,
		})
	}
}