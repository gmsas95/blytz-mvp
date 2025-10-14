package api

import (
	"github.com/blytz/payment-service/internal/config"
	"github.com/blytz/payment-service/internal/database"
	"github.com/blytz/payment-service/internal/redis"
	"github.com/blytz/payment-service/internal/stripe"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Router struct {
	engine *gin.Engine
	db     *database.Connection
	rdb    *redis.Client
	stripe *stripe.Client
	cfg    *config.Config
}

func NewRouter(db *database.Connection, rdb *redis.Client, stripe *stripe.Client, cfg *config.Config) *Router {
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := &Router{
		engine: gin.New(),
		db:     db,
		rdb:    rdb,
		stripe: stripe,
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
	api := r.engine.Group("/payment")
	{
		// Payment endpoints
		api.POST("/create", r.createPayment)
		api.GET("/status/:payment_id", r.getPaymentStatus)
		api.POST("/confirm", r.confirmPayment)
		api.POST("/cancel", r.cancelPayment)

		// Refund endpoints
		api.POST("/refund", r.createRefund)
		api.GET("/refund/:refund_id", r.getRefundStatus)

		// User payment history
		api.GET("/history/:user_id", r.getUserPaymentHistory)

		// Stripe webhook
		api.POST("/webhook", r.handleStripeWebhook)

		// Payment methods
		api.POST("/methods", r.createPaymentMethod)
		api.GET("/methods/:user_id", r.getUserPaymentMethods)
		api.DELETE("/methods/:method_id", r.deletePaymentMethod)
	}
}

func (r *Router) Run(addr string) error {
	return r.engine.Run(addr)
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Authorization, Accept, X-Requested-With, Stripe-Signature")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}