package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/blytz/payment-service/internal/api"
	"github.com/blytz/payment-service/internal/config"
	"github.com/blytz/payment-service/internal/database"
	"github.com/blytz/payment-service/internal/redis"
	"github.com/blytz/payment-service/internal/stripe"
)

func main() {
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Initialize Redis client
	rdb, err := redis.NewClient(ctx, cfg.RedisURL)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer rdb.Close()

	// Initialize database
	db, err := database.NewConnection(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize Stripe
	stripeClient := stripe.NewClient(cfg.StripeSecretKey)

	// Auto-migrate database tables
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	router := api.NewRouter(db, rdb, stripeClient, cfg)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("Starting payment service on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}