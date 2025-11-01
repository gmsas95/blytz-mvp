package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/gmsas95/blytz-mvp/services/auction-service/internal/api/handlers"
	"github.com/gmsas95/blytz-mvp/services/auction-service/internal/config"
	"github.com/gmsas95/blytz-mvp/services/auction-service/internal/services"
	"github.com/gmsas95/blytz-mvp/services/auction-service/pkg/firebase"
	"github.com/gmsas95/blytz-mvp/shared/pkg/auth"
)

func SetupRouter(auctionService *services.AuctionService, logger *zap.Logger, cfg *config.Config) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Initialize auth client
	authClient := auth.NewAuthClient("http://auth-service:8084")

	// Initialize Firebase client
	firebaseClient := firebase.NewClient(logger)

	auctionHandler := handlers.NewAuctionHandler(auctionService, logger, firebaseClient)
	livekitHandler := handlers.NewLiveKitHandler(logger)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Auction endpoints
	auctionRoutes := router.Group("/api/v1/auctions")
	{
		auctionRoutes.GET("/", auctionHandler.ListAuctions)
		auctionRoutes.GET("/active", auctionHandler.GetActiveAuctions)
		auctionRoutes.GET("/:id", auctionHandler.GetAuction)

		// Protected routes
		protected := auctionRoutes.Group("")
		protected.Use(auth.GinAuthMiddleware(authClient))
		{
			protected.POST("/", auctionHandler.CreateAuction)
			protected.POST("/:id/bids", auctionHandler.PlaceBid)
		}
	}

	// LiveKit token endpoint
	livekitRoutes := router.Group("/api/v1/livekit")
	livekitRoutes.Use(auth.GinAuthMiddleware(authClient))
	{
		livekitRoutes.GET("/token", livekitHandler.GenerateToken)
	}

	return router
}
