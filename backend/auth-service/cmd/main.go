package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/blytz/auth-service/internal/api"
	"github.com/blytz/auth-service/internal/config"
	"github.com/blytz/auth-service/internal/services"
	"github.com/blytz/shared/utils"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	logger, err := utils.InitLogger("development")
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// Load configuration
	cfg := config.Load()

	// Create context with cancellation
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Initialize auth service
	authService := services.NewAuthService(logger, cfg)

	// Create API router
	router := api.NewRouter(authService, logger, cfg)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.Port
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan

		logger.Info("Shutting down server...")
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Error("Server forced to shutdown", zap.Error(err))
		}
	}()

	logger.Info("Auth service started", zap.String("port", port))
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal("Server failed to start", zap.Error(err))
	}

	logger.Info("Server stopped")
}