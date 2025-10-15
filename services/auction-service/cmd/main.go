package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/blytz/auction-service/internal/api"
	"github.com/blytz/auction-service/internal/config"
	"github.com/blytz/auction-service/internal/redis"
)

func main() {
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rdb, err := redis.NewClient(ctx, cfg.RedisURL)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer rdb.Close()

	if err := rdb.Ping(ctx); err != nil {
		log.Fatalf("Redis ping failed: %v", err)
	}

	if err := rdb.LoadScripts(ctx); err != nil {
		log.Fatalf("Failed to load Lua scripts: %v", err)
	}

	router := api.NewRouter(rdb, cfg)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting auction service on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}