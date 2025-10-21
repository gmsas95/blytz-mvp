package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/gmsas95/blytz-mvp/services/gateway/internal/api"
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
		port = "8080"
	}

	// Start the server
	logger.Info("Starting server on port " + port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}