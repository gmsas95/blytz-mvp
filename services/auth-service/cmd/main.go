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
	"github.com/gmsas95/blytz-mvp/services/auth-service/internal/api"
	"github.com/gmsas95/blytz-mvp/services/auth-service/internal/config"
	"github.com/gmsas95/blytz-mvp/services/auth-service/internal/models"
	"github.com/gmsas95/blytz-mvp/services/auth-service/internal/services"
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

	nodeEnv := os.Getenv("NODE_ENV")
	logger.Info("NODE_ENV", zap.String("value", nodeEnv))

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
	}

	logger.Info("Configuration loaded",
		zap.String("port", cfg.ServicePort),
		zap.String("environment", cfg.Environment))

	// Set Gin mode based on environment
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin router
	router := gin.New()
	router.Use(gin.Recovery())

	// Initialize database connection
	db, err := config.InitDB(cfg)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	// Auto-migrate the User model
	db.AutoMigrate(&models.User{})

	// Initialize auth service
	authService := services.NewAuthService(db, cfg)

	// Setup routes using the API router
	api.NewRouter(router, authService, logger, cfg)

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + cfg.ServicePort,
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

