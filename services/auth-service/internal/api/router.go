package api

import (
	"fmt"
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

	// Add logging middleware
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

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "ok",
			"service":   "auth",
			"timestamp": time.Now().Unix(),
		})
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
