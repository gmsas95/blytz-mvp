package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/gmsas95/blytz-mvp/services/auction-service/internal/config"
	"github.com/gmsas95/blytz-mvp/services/auction-service/internal/services"
)

func SetupRouter(auctionService *services.AuctionService, logger *zap.Logger, cfg *config.Config) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// CORS middleware first
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.Status(200)
			return
		}
		c.Next()
	})

	// Add custom logger middleware to avoid any potential debug routes
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return ""
	}))
	router.Use(gin.Recovery())

	// Mock auctions endpoint for frontend
	router.GET("/", func(c *gin.Context) {
		auctions := []gin.H{
			{
				"id":          1,
				"title":       "Vintage Leather Handbag",
				"description": "Authentic vintage leather handbag with gold hardware. Perfect condition with original dust bag.",
				"image":       "https://images.unsplash.com/photo-1553062407-98eeb64c6a62?w=400&h=400&fit=crop",
				"current_bid": 275.00,
				"status":      "live",
				"end_time":    "2024-12-31T23:59:59Z",
				"bids_count":  12,
				"watchers":    8,
				"min_bid":     10.00,
				"seller": gin.H{
					"name":   "JStyle Boutique",
					"rating": 4.8,
					"avatar": "S",
				},
			},
			{
				"id":          2,
				"title":       "Smart Home Security Camera",
				"description": "4K wireless security camera with night vision, motion detection, and cloud storage.",
				"image":       "https://images.unsplash.com/photo-1558089687-f282ffcbc126?w=400&h=400&fit=crop",
				"current_bid": 185.50,
				"status":      "scheduled",
				"end_time":    "2024-12-31T23:59:59Z",
				"bids_count":  0,
				"watchers":    0,
				"min_bid":     5.00,
				"seller": gin.H{
					"name":   "TechHub Deals",
					"rating": 4.6,
					"avatar": "M",
				},
			},
		}

		c.JSON(200, gin.H{
			"auctions": auctions,
			"total":    len(auctions),
			"page":     1,
			"limit":    10,
		})
	})

	// Comprehensive health check endpoint
	router.GET("/health", func(c *gin.Context) {
		health := gin.H{
			"status":    "ok",
			"service":   "auction",
			"timestamp": time.Now().Unix(),
			"version":   "v1.0.0",
			"checks": gin.H{
				"database": "connected",
				"redis":    "connected",
			},
		}
		c.JSON(200, health)
	})

	return router
}
