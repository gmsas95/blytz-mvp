package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/gmsas95/blytz-mvp/services/auction-service/internal/config"
	"github.com/gmsas95/blytz-mvp/services/auction-service/internal/services"
)

func SetupRouter(auctionService *services.AuctionService, logger *zap.Logger, cfg *config.Config) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Add custom logger middleware to avoid any potential debug routes
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return ""
	}))
	router.Use(gin.Recovery())

	// Health check endpoint only
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return router
}
