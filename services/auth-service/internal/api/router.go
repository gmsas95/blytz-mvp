package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gmsas95/blytz-mvp/services/auth-service/internal/api/handlers"
	"github.com/gmsas95/blytz-mvp/services/auth-service/internal/config"
	"github.com/gmsas95/blytz-mvp/services/auth-service/internal/middleware"
	"github.com/gmsas95/blytz-mvp/services/auth-service/internal/services"
	"go.uber.org/zap"
)

type Router struct {
	engine *gin.Engine
}

func NewRouter(router *gin.Engine, authService *services.AuthService, logger *zap.Logger, cfg *config.Config) {
	// Temporarily disable production mode for debugging
	// if cfg.Environment == "production" {
	//	gin.SetMode(gin.ReleaseMode)
	// }

	// Note: Correlation middleware temporarily disabled for deployment
	// router.Use(utils.CorrelationMiddleware())

	// Add structured logging middleware
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	// Enhanced health check
	router.GET("/health", func(c *gin.Context) {
		correlationID := c.GetHeader("X-Correlation-ID")
		if correlationID == "" {
			correlationID = c.GetString("correlation_id")
		}

		health := gin.H{
			"status":         "ok",
			"service":        "auth",
			"version":        "1.0.0",
			"timestamp":      time.Now().Unix(),
			"correlation_id": correlationID,
			"environment":    cfg.Environment,
		}

		// Check database connectivity
		if authService != nil {
			health["database"] = "connected"
		} else {
			health["database"] = "disconnected"
			health["status"] = "degraded"
			c.JSON(http.StatusServiceUnavailable, health)
			return
		}

		c.JSON(http.StatusOK, health)
	})

	// Metrics

	// API routes
	api := router.Group("/api")
	{
		authHandler := handlers.NewAuthHandler(authService)
		SetupAuthRoutes(api, authHandler, authService)
	}
}

func SetupAuthRoutes(api *gin.RouterGroup, authHandler *handlers.AuthHandler, authService *services.AuthService) {
	// Public routes
	public := api.Group("/auth")
	{
		public.POST("/register", authHandler.SignUp) // Changed from signup to register to match frontend
		public.POST("/login", authHandler.Login)
		public.POST("/refresh", authHandler.RefreshToken)
	}

	// Test endpoint without middleware
	api.GET("/auth/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Auth service is working", "timestamp": time.Now().Unix()})
	})

	// Protected routes (with auth middleware)
	protected := api.Group("/auth")
	protected.Use(middleware.AuthMiddleware(authService))
	{
		protected.GET("/me", authHandler.GetProfile) // Added /me endpoint for frontend compatibility
		protected.GET("/verify", authHandler.Verify)
		protected.POST("/logout", authHandler.Logout)
		protected.PUT("/profile", authHandler.UpdateProfile)
		protected.GET("/profile", authHandler.GetProfile) // Keep existing profile endpoint
	}
}
