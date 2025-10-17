package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/blytz/auth-service/internal/api"
	"github.com/blytz/auth-service/internal/config"
	"github.com/blytz/auth-service/internal/services"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer logger.Sync()

	logger.Info("Starting auth service...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
	}

	logger.Info("Configuration loaded", 
		zap.String("port", cfg.GetServicePort()),
		zap.String("environment", cfg.GetEnvironment()))

	// Set Gin mode based on environment
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin router
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Initialize auth service
	authService, err := services.NewAuthService(cfg, logger)
	if err != nil {
		logger.Fatal("Failed to initialize auth service", zap.Error(err))
	}

	// Create auth handler
	authHandler := api.NewAuthHandler(authService, logger)

	// Setup routes
	setupRoutes(router, authHandler, logger)

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + cfg.GetServicePort(),
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		logger.Info("Starting HTTP server", zap.String("address", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}

// setupRoutes configures all API routes
func setupRoutes(router *gin.Engine, authHandler *api.AuthHandler, logger *zap.Logger) {
	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"service":   "auth",
			"timestamp": time.Now().Unix(),
		})
	})

	// Metrics endpoint
	router.GET("/metrics", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "metrics endpoint",
		})
	})

	// API v1 routes
	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.RegisterUser)
			auth.POST("/login", authHandler.LoginUser)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/validate", authHandler.ValidateToken) // Internal endpoint
			
			// Protected routes
			auth.GET("/me", authHandler.GetCurrentUser)
			auth.PUT("/profile", authHandler.UpdateProfile)
			auth.POST("/change-password", authHandler.ChangePassword)
		}
	}
}
