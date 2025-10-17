package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/blytz/auction-service/internal/api"
	"github.com/blytz/auction-service/internal/config"
	"github.com/blytz/auction-service/internal/repository"
	"github.com/blytz/auction-service/internal/services"
	"github.com/blytz/shared/pkg/utils"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	logger, err := utils.InitLogger("development")
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// Initialize database connection
	db, err := initDatabase(logger)
	if err != nil {
		logger.Fatal("Failed to initialize database", zap.Error(err))
	}
	defer db.Close()

	// Initialize repository
	auctionRepo := repository.NewPostgresAuctionRepository(db, logger)

	// Test database connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := auctionRepo.Ping(ctx); err != nil {
		logger.Error("Database connection test failed", zap.Error(err))
		// Continue with mock data for demo purposes
		auctionRepo = nil
	} else {
		logger.Info("Database connection successful")

		// Initialize database schema
		if err := initDatabaseSchema(ctx, db, logger); err != nil {
			logger.Error("Failed to initialize database schema", zap.Error(err))
			// Continue with existing schema
		}
	}

	// Create services with repository
	auctionService := services.NewAuctionServiceWithRepo(logger, &config.Config{}, auctionRepo)

	// Create API router
	router := api.SetupRouterWithServices(logger, auctionService)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
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

	logger.Info("Auction service started", zap.String("port", port))
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal("Server failed to start", zap.Error(err))
	}

	logger.Info("Server stopped")
}

func initDatabase(logger *zap.Logger) (*sql.DB, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://postgres:password@localhost:5432/auction_db?sslmode=disable"
	}

	logger.Info("Connecting to database", zap.String("url", databaseURL))

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}

func initDatabaseSchema(ctx context.Context, db *sql.DB, logger *zap.Logger) error {
	logger.Info("Initializing database schema")

	// Create tables if they don't exist
	schema := `
		CREATE TABLE IF NOT EXISTS auctions (
			auction_id VARCHAR(255) PRIMARY KEY,
			product_id VARCHAR(255) NOT NULL,
			seller_id VARCHAR(255) NOT NULL,
			title VARCHAR(500) NOT NULL,
			description TEXT,
			starting_price DECIMAL(10,2) NOT NULL,
			current_price DECIMAL(10,2) NOT NULL,
			reserve_price DECIMAL(10,2),
			min_bid_increment DECIMAL(10,2) NOT NULL DEFAULT 1.00,
			start_time TIMESTAMP NOT NULL,
			end_time TIMESTAMP NOT NULL,
			status VARCHAR(50) NOT NULL DEFAULT 'scheduled',
			type VARCHAR(50) NOT NULL DEFAULT 'live',
			is_active BOOLEAN NOT NULL DEFAULT true,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS bids (
			bid_id VARCHAR(255) PRIMARY KEY,
			auction_id VARCHAR(255) NOT NULL,
			bidder_id VARCHAR(255) NOT NULL,
			amount DECIMAL(10,2) NOT NULL,
			is_winning BOOLEAN NOT NULL DEFAULT false,
			bid_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (auction_id) REFERENCES auctions(auction_id) ON DELETE CASCADE
		);

		CREATE INDEX IF NOT EXISTS idx_auctions_status ON auctions(status);
		CREATE INDEX IF NOT EXISTS idx_auctions_seller_id ON auctions(seller_id);
		CREATE INDEX IF NOT EXISTS idx_bids_auction_id ON bids(auction_id);
	`

	_, err := db.ExecContext(ctx, schema)
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	logger.Info("Database schema initialized successfully")
	return nil
}