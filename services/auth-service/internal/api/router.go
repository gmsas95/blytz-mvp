package api

import (
	"time"
	"github.com/gin-gonic/gin"
	"github.com/gmsas95/blytz-mvp/services/auth-service/internal/services"
	"github.com/gmsas95/blytz-mvp/services/auth-service/internal/config"
	"github.com/gmsas95/blytz-mvp/services/auth-service/internal/api/handlers"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type Router struct {
	engine *gin.Engine
}

func NewRouter(router *gin.Engine, authService *services.AuthService, logger *zap.Logger, cfg *config.Config) {
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"service": "auth",
			"timestamp": time.Now().Unix(),
		})
	})

	// Metrics
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// API routes
	api := router.Group("/api/v1")
	{
		authHandler := handlers.NewAuthHandler(authService)
		SetupAuthRoutes(api, authHandler)
	}
}

func SetupAuthRoutes(api *gin.RouterGroup, authHandler *handlers.AuthHandler) {
	// Public routes
	public := api.Group("/auth")
	{
		public.POST("/signup", authHandler.SignUp)
		public.POST("/login", authHandler.Login)
		public.POST("/refresh", authHandler.RefreshToken)
	}

	// Protected routes (would need auth middleware in real implementation)
	protected := api.Group("/auth")
	{
		protected.GET("/verify", authHandler.Verify)
		protected.POST("/logout", authHandler.Logout)
		protected.PUT("/profile", authHandler.UpdateProfile)
		protected.GET("/profile", authHandler.GetProfile)
	}
}