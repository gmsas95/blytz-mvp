package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/blytz/auction-service/internal/services"
	"github.com/blytz/auction-service/internal/config"
	"github.com/blytz/auction-service/internal/api/handlers"
	"github.com/blytz/auction-service/pkg/firebase"
	"github.com/blytz/shared/pkg/auth"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func SetupRouter(logger *zap.Logger) *gin.Engine {
	// Initialize config
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}

	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Health check (public)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"service": "auction",
			"timestamp": time.Now().Unix(),
		})
	})

	// Metrics (public)
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Initialize auth client for service-to-service communication
	authClient := auth.NewAuthClient("http://auth-service:8084")

	// Initialize services
	auctionService := services.NewAuctionService(logger, cfg)

	// Initialize Firebase client
	firebaseClient := firebase.NewClient(logger)
	auctionHandler := handlers.NewAuctionHandler(auctionService, firebaseClient, logger)

	// API routes
	api := router.Group("/api/v1")
	{
		// Public auction routes (view only)
		auctions := api.Group("/auctions")
		{
			auctions.GET("", auctionHandler.ListAuctions)
			auctions.GET("/:auction_id", auctionHandler.GetAuction)
			auctions.GET("/:auction_id/status", auctionHandler.GetAuctionStatus)
			auctions.GET("/:auction_id/bids", auctionHandler.GetBids)
		}

		// Protected auction routes (require authentication)
		protectedAuctions := api.Group("/auctions")
		protectedAuctions.Use(auth.GinAuthMiddleware(authClient))
		{
			protectedAuctions.POST("", auctionHandler.CreateAuction)
			protectedAuctions.PUT("/:auction_id", auctionHandler.UpdateAuction)
			protectedAuctions.DELETE("/:auction_id", auctionHandler.DeleteAuction)
			protectedAuctions.POST("/:auction_id/bids", auctionHandler.PlaceBid)
		}
	}

	return router
}