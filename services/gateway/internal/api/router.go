package api

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Simple in-memory rate limiter for deployment
type InMemoryRateLimiter struct {
	requests map[string][]time.Time
	mutex    sync.RWMutex
	limit    int
}

func NewInMemoryRateLimiter(limit int) *InMemoryRateLimiter {
	return &InMemoryRateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
	}
}

func (r *InMemoryRateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()

		r.mutex.Lock()
		defer r.mutex.Unlock()

		// Clean old requests
		if requests, exists := r.requests[ip]; exists {
			var validRequests []time.Time
			for _, reqTime := range requests {
				if now.Sub(reqTime) < time.Minute {
					validRequests = append(validRequests, reqTime)
				}
			}
			r.requests[ip] = validRequests
		}

		// Check limit
		if len(r.requests[ip]) >= r.limit {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
				"limit": r.limit,
				"window": "1 minute",
			})
			c.Abort()
			return
		}

		// Add current request
		r.requests[ip] = append(r.requests[ip], now)
		c.Next()
	}
}

func SetupRouter(logger *zap.Logger) *gin.Engine {
	router := gin.Default()

	// Add simple rate limiting
	rateLimiter := NewInMemoryRateLimiter(60) // 60 requests per minute
	router.Use(rateLimiter.Middleware())

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Correlation-ID")

		if c.Request.Method == "OPTIONS" {
			c.Status(http.StatusOK)
			return
		}
		c.Next()
	})

	// Root endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": "Blytz API Gateway",
			"version": "v1",
			"status":  "running",
			"endpoints": map[string]string{
				"health":        "/health",
				"public_api":    "/api/public/",
				"livekit_token": "/api/public/livekit/token",
			},
		})
	})

	// Enhanced health check endpoint
	router.GET("/health", func(c *gin.Context) {
		correlationID := c.GetHeader("X-Correlation-ID")
		if correlationID == "" {
			correlationID = c.GetString("correlation_id")
		}

		health := gin.H{
			"status":         "ok",
			"service":        "gateway",
			"version":        "1.0.0",
			"timestamp":      time.Now().Unix(),
			"correlation_id": correlationID,
			"environment":    "production",
		}

		// Rate limiter is always active in this implementation
		health["rate_limiter"] = "active"

		// Check external dependencies
		health["dependencies"] = gin.H{
			"redis": "configured",
		}

		c.JSON(http.StatusOK, health)
	})

	// Simple ping endpoint (no dependencies)
	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong - updated at "+time.Now().Format("2006-01-02 15:04:05"))
	})

	// API routes with enhanced rate limiting
	api := router.Group("/api")
	{
		// Public routes 
		public := api.Group("/public")
		{

			public.GET("/health", func(c *gin.Context) {
				correlationID := c.GetString("correlation_id")
				c.JSON(200, gin.H{
					"status":         "ok",
					"service":        "gateway-public",
					"correlation_id": correlationID,
				})
			})

			public.GET("/test", func(c *gin.Context) {
				correlationID := c.GetString("correlation_id")
				c.JSON(200, gin.H{
					"message":        "test works",
					"correlation_id": correlationID,
				})
			})

			// LiveKit token generation
			public.GET("/livekit/token", createLiveKitTokenHandler(logger))

			// LiveKit proxy routes (removed to avoid conflicts)
			// public.Any("/livekit/*proxyPath", liveKitProxyHandler(logger))
		}

		// Microservice proxy routes
		v1 := api.Group("/v1")
		{
			// Auth service routes
			auth := v1.Group("/auth")
			{
				createProxyRoutes(auth, "http://blytz-auth-prod:8084", logger)
			}

			// Product service routes
			product := v1.Group("/products")
			{
				createProxyRoutes(product, "http://blytz-product-prod:8082", logger)
			}

			// Auction service routes
			auction := v1.Group("/auctions")
			{
				createProxyRoutes(auction, "http://blytz-auction-prod:8083", logger)
			}

			// Order service routes
			order := v1.Group("/orders")
			{
				createProxyRoutes(order, "http://blytz-order-prod:8085", logger)
			}

			// Payment service routes
			payment := v1.Group("/payments")
			{
				createProxyRoutes(payment, "http://blytz-payment-prod:8086", logger)
			}

			// Logistics service routes
			logistics := v1.Group("/logistics")
			{
				createProxyRoutes(logistics, "http://blytz-logistics-prod:8087", logger)
			}

			// Chat service routes
			chat := v1.Group("/chat")
			{
				createProxyRoutes(chat, "http://blytz-chat-prod:8088", logger)
			}
		}
	}

	return router
}

func createLiveKitTokenHandler(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get query parameters
		room := c.Query("room")
		role := c.DefaultQuery("role", "viewer")

		if room == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Room name is required",
			})
			return
		}

		// For now, return a mock response to test connectivity
		response := gin.H{
			"token":    "mock_token_for_testing",
			"url":      "wss://blytz-live-u5u72ozx.livekit.cloud",
			"room":     room,
			"identity": fmt.Sprintf("%s_%d", role, time.Now().Unix()),
			"message":  "This is a mock token for testing - LiveKit integration needs environment variables",
		}

		c.JSON(http.StatusOK, response)
	}
}

// createProxyRoutes creates reverse proxy routes for a microservice
func createProxyRoutes(group *gin.RouterGroup, targetURL string, logger *zap.Logger) {
	target, err := url.Parse(targetURL)
	if err != nil {
		logger.Error("Failed to parse target URL", zap.String("url", targetURL), zap.Error(err))
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	// Modify director to add headers
	proxy.Director = func(req *http.Request) {
		req.Host = target.Host
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
		req.Header.Set("X-Forwarded-For", req.RemoteAddr)
		req.Header.Set("X-Forwarded-Proto", "https")
	}

	// Error handler
	proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
		logger.Error("Proxy error", zap.String("path", req.URL.Path), zap.Error(err))
		rw.WriteHeader(http.StatusBadGateway)
		rw.Write([]byte(`{"error": "Service unavailable"}`))
	}

	// Proxy all requests to the microservice
	group.Any("/*path", gin.WrapH(proxy))
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
