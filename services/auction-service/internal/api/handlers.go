package api

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Health check endpoint
func (r *Router) healthCheck(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := r.rdb.Ping(ctx); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "unhealthy",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"service": "auction",
		"timestamp": time.Now().Unix(),
	})
}

// Get auction details
func (r *Router) getAuction(c *gin.Context) {
	auctionID := c.Param("id")
	if auctionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "auction ID required"})
		return
	}

	ctx := context.Background()
	auction, err := r.rdb.GetAuction(ctx, auctionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get auction"})
		return
	}

	if len(auction) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "auction not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"auction_id": auctionID,
		"data": auction,
	})
}

// Place a bid on an auction
func (r *Router) placeBid(c *gin.Context) {
	auctionID := c.Param("id")
	if auctionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "auction ID required"})
		return
	}

	var bidRequest struct {
		UserID string `json:"user_id" binding:"required"`
		Amount int64  `json:"amount" binding:"required,gt=0"`
	}

	if err := c.ShouldBindJSON(&bidRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	success, err := r.rdb.ProcessBid(ctx, auctionID, bidRequest.UserID, bidRequest.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "bid processing failed"})
		return
	}

	if !success {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bid rejected - auction ended or bid too low"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"amount": bidRequest.Amount,
		"auction_id": auctionID,
	})
}

// Create a new auction
func (r *Router) createAuction(c *gin.Context) {
	var createRequest struct {
		AuctionID    string `json:"auction_id" binding:"required"`
		StartTime    int64  `json:"start_time" binding:"required"`
		EndTime      int64  `json:"end_time" binding:"required,gtfield=StartTime"`
		ReservePrice int64  `json:"reserve_price" binding:"gte=0"`
	}

	if err := c.ShouldBindJSON(&createRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	err := r.rdb.CreateAuction(ctx, createRequest.AuctionID, createRequest.StartTime, createRequest.EndTime, createRequest.ReservePrice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create auction"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"auction_id": createRequest.AuctionID,
		"status": "created",
	})
}

// End an auction
func (r *Router) endAuction(c *gin.Context) {
	auctionID := c.Param("id")
	if auctionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "auction ID required"})
		return
	}

	ctx := context.Background()
	err := r.rdb.EndAuction(ctx, auctionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to end auction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"auction_id": auctionID,
		"status": "ended",
	})
}

// Stream auction updates (WebSocket placeholder)
func (r *Router) auctionStream(c *gin.Context) {
	auctionID := c.Param("id")
	if auctionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "auction ID required"})
		return
	}

	// For now, return a placeholder - WebSocket implementation will be added
	c.JSON(http.StatusOK, gin.H{
		"auction_id": auctionID,
		"message": "WebSocket streaming will be implemented here",
		"endpoint": "/auction/" + auctionID + "/stream",
	})
}