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
	"go.uber.org/zap"

	"github.com/gmsas95/blytz-mvp/services/logistics-service/internal/api"
	"github.com/gmsas95/blytz-mvp/shared/pkg/utils"
)

func main() {
	// Initialize logger
	logger, err := utils.InitLogger("production")
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync() // flushes buffer, if any

	// Set Gin to release mode
	gin.SetMode(gin.ReleaseMode)

	// Create a new Gin router
	router := api.SetupRouter(logger)


	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8087"
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + port,
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