package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/gmsas95/blytz-mvp/services/auction-service/internal/api"
	"github.com/gmsas95/blytz-mvp/services/auction-service/internal/config"
	"github.com/gmsas95/blytz-mvp/services/auction-service/internal/services"
	"github.com/gmsas95/blytz-mvp/shared/pkg/utils"
)

func main() {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	logger, err := utils.InitLogger("production")
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync() // flushes buffer, if any

	// Initialize database
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// Initialize services
	auctionService := services.NewAuctionService(db, logger, cfg)

	// Create a new Gin router
	fmt.Println("DEBUG: About to setup router")
	router := api.SetupRouter(auctionService, logger, cfg)
	fmt.Println("DEBUG: Router setup completed")

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
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
