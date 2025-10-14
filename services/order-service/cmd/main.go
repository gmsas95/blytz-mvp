package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/blytz/order-service/internal/api"
	"github.com/blytz/order-service/internal/config"
	"github.com/blytz/order-service/internal/database"
	"github.com/blytz/order-service/internal/redis"
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

	// Auto-migrate database tables
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	router := api.NewRouter(db, rdb, cfg)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	log.Printf("Starting order service on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}