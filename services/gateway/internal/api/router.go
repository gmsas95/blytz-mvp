package api

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	shared_auth "github.com/gmsas95/blytz-mvp/shared/pkg/auth"
)

func SetupRouter(logger *zap.Logger) *gin.Engine {
	router := gin.Default()

	// Initialize auth client
	authClient := shared_auth.NewAuthClient("http://auth-service:8084")

	// Root endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": "Blytz API Gateway",
			"version": "v1",
			"status":  "running",
			"endpoints": map[string]string{
				"health":        "/health",
				"metrics":       "/metrics",
				"public_api":    "/api/public/",
				"auth_api":      "/api/auth/",
				"protected_api": "/api/v1/",
				"livekit_token": "/api/public/livekit/token",
			},
		})
	})

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "gateway"})
	})

	// Prometheus metrics endpoint
	router.GET("/gateway-metrics", gin.WrapH(promhttp.Handler()))

	// API routes
	api := router.Group("/api")
	{
		// Public routes
		public := api.Group("/public")
		{
			public.GET("/health", func(c *gin.Context) {
				c.JSON(200, gin.H{"status": "ok"})
			})

			// LiveKit token generation (public for demo purposes)
			public.Any("/livekit/token", proxyToService("http://livekit-service:8089", logger))
		}

		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.Any("/*proxyPath", proxyToServiceWithPath("http://auth-service:8084", "/api/v1/auth", logger))
		}

		// Webhook routes (public - no auth required) - moved before protected routes
		webhooks := api.Group("/v1/webhooks")
		{
			webhooks.Any("/fiuu", func(c *gin.Context) {
				// Direct proxy to payment service webhook endpoint
				target, err := url.Parse("http://payment-service:8086")
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid target URL"})
					return
				}

				proxy := httputil.NewSingleHostReverseProxy(target)
				proxy.Director = func(req *http.Request) {
					req.URL.Scheme = target.Scheme
					req.URL.Host = target.Host
					req.URL.Path = "/api/v1/webhooks/fiuu"
					req.Host = target.Host
				}

				proxy.ServeHTTP(c.Writer, c.Request)
			})

		}

		// Protected routes
		protected := api.Group("/v1")
		protected.Use(shared_auth.GinAuthMiddleware(authClient))
		{
			// Proxy to other services
			protected.Any("/auctions/*proxyPath", proxyToService("http://auction-service:8083", logger))
			protected.Any("/products/*proxyPath", proxyToService("http://product-service:8082", logger))
			protected.Any("/orders/*proxyPath", proxyToService("http://order-service:8085", logger))
			protected.Any("/payments/*proxyPath", proxyToService("http://payment-service:8086", logger))
			protected.Any("/logistics/*proxyPath", proxyToService("http://logistics-service:8087", logger))
			protected.Any("/chat/*proxyPath", proxyToService("http://chat-service:8088", logger))
		}
	}

	return router
}

func proxyToServiceWithPath(targetURL string, targetPath string, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse the target URL
		target, err := url.Parse(targetURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid target URL"})
			return
		}

		// Create a single-host reverse proxy
		proxy := httputil.NewSingleHostReverseProxy(target)

		// Customize the director to properly handle the path
		originalDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			originalDirector(req)

			// Extract the proxy path and prepend the target path
			proxyPath := c.Param("proxyPath")
			if proxyPath != "" && !strings.HasPrefix(proxyPath, "/") {
				proxyPath = "/" + proxyPath
			}

			// Combine target path with proxy path
			fullPath := targetPath + proxyPath

			// Set the proper path and raw path
			req.URL.Path = fullPath
			req.URL.RawPath = fullPath

			// Copy headers from original request
			for key, values := range c.Request.Header {
				for _, value := range values {
					req.Header.Add(key, value)
				}
			}

			// Set the proper host
			req.Host = target.Host
		}

		// Handle proxy errors
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			logger.Error("Proxy error",
				zap.String("target", targetURL),
				zap.String("path", r.URL.Path),
				zap.Error(err))
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte(`{"error": "Service unavailable", "message": "The requested service is not available"}`))
		}

		// Set a timeout for the proxy request
		proxy.Transport = &http.Transport{
			ResponseHeaderTimeout: 30 * time.Second,
		}

		// Serve the request through the proxy
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func proxyToService(targetURL string, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse the target URL
		target, err := url.Parse(targetURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid target URL"})
			return
		}

		// Create a single-host reverse proxy
		proxy := httputil.NewSingleHostReverseProxy(target)

		// Customize the director to properly handle the path
		originalDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			originalDirector(req)

			// Extract the proxy path and clean it
			proxyPath := c.Param("proxyPath")
			if proxyPath != "" && !strings.HasPrefix(proxyPath, "/") {
				proxyPath = "/" + proxyPath
			}

			// Set the proper path and raw path
			req.URL.Path = proxyPath
			req.URL.RawPath = proxyPath

			// Copy headers from original request
			for key, values := range c.Request.Header {
				for _, value := range values {
					req.Header.Add(key, value)
				}
			}

			// Set the proper host
			req.Host = target.Host
		}

		// Handle proxy errors
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			logger.Error("Proxy error",
				zap.String("target", targetURL),
				zap.String("path", r.URL.Path),
				zap.Error(err))
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte(`{"error": "Service unavailable", "message": "The requested service is not available"}`))
		}

		// Set a timeout for the proxy request
		proxy.Transport = &http.Transport{
			ResponseHeaderTimeout: 30 * time.Second,
		}

		// Serve the request through the proxy
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
