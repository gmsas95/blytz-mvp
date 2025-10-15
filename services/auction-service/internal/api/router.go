package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/blytz/auction-service/internal/services"
	"github.com/blytz/auction-service/internal/config"
	"github.com/blytz/auction-service/internal/api/handlers"
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

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"service": "auction",
			"timestamp": time.Now().Unix(),
		})
	})

	// Metrics
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Initialize services
	auctionService := services.NewAuctionService(logger, cfg)
	auctionHandler := handlers.NewAuctionHandler(auctionService)

	// API routes
	api := router.Group("/api/v1")
	{
		// Auction routes
		auctions := api.Group("/auctions")
		{
			auctions.GET("", auctionHandler.ListAuctions)
			auctions.POST("", auctionHandler.CreateAuction)
			auctions.GET("/:auction_id", auctionHandler.GetAuction)
			auctions.PUT("/:auction_id", auctionHandler.UpdateAuction)
			auctions.DELETE("/:auction_id", auctionHandler.DeleteAuction)
			auctions.GET("/:auction_id/status", auctionHandler.GetAuctionStatus)

			// Bid routes
			auctions.POST("/:auction_id/bids", auctionHandler.PlaceBid)
			auctions.GET("/:auction_id/bids", auctionHandler.GetBids)
		}
	}

	return router
}