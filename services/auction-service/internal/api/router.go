package api

import (
	"context"
	"net/http"

	"github.com/blytz/auction-service/internal/config"
	"github.com/blytz/auction-service/internal/redis"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Router struct {
	engine *gin.Engine
	rdb    *redis.Client
	cfg    *config.Config
}

func NewRouter(rdb *redis.Client, cfg *config.Config) *Router {
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := &Router{
		engine: gin.New(),
		rdb:    rdb,
		cfg:    cfg,
	}

	r.setupMiddleware()
	r.setupRoutes()

	return r
}

func (r *Router) setupMiddleware() {
	r.engine.Use(gin.Logger())
	r.engine.Use(gin.Recovery())
	r.engine.Use(corsMiddleware())
}

func (r *Router) setupRoutes() {
	// Health check
	r.engine.GET("/health", r.healthCheck)

	// Prometheus metrics
	r.engine.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// API routes
	api := r.engine.Group("/auction")
	{
		api.GET("/auctions/:id", r.getAuction)
		api.POST("/auctions/:id/bid", r.placeBid)
		api.POST("/auctions", r.createAuction)
		api.POST("/auctions/:id/end", r.endAuction)
		api.GET("/auctions/:id/stream", r.auctionStream)
	}
}

func (r *Router) Run(addr string) error {
	return r.engine.Run(addr)
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Authorization, Accept, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}