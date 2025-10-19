package api

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"github.com/gmsas95/blytz-mvp/services/auction-service/internal/services"
	"github.com/gmsas95/blytz-mvp/services/auction-service/internal/config"
	"github.com/gmsas95/blytz-mvp/services/auction-service/internal/api/handlers"
	"github.com/gmsas95/blytz-mvp/services/auction-service/pkg/firebase"
	"github.com/gmsas95/blytz-mvp/shared/pkg/auth"
)

func SetupRouter(auctionService *services.AuctionService, logger *zap.Logger, cfg *config.Config) *gin.Engine {
	router := gin.Default()

	authMiddleware := auth.NewAuthMiddleware(cfg.JWTSecret)

	auctionHandler := handlers.NewAuctionHandler(auctionService, logger, nil)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Prometheus metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Auction endpoints
	auctionRoutes := router.Group("/auctions")
	{
		auctionRoutes.POST("/", authMiddleware.Middleware(), auctionHandler.CreateAuction)
		auctionRoutes.GET("/:id", auctionHandler.GetAuction)
		auctionRoutes.POST("/:id/bids", authMiddleware.Middleware(), auctionHandler.PlaceBid)
	}

	return router
}