package handlers

import (
	"context"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/blytz/auction-service/internal/models"
	"github.com/blytz/auction-service/internal/services"
	"github.com/blytz/auction-service/pkg/firebase"
	"github.com/blytz/shared/pkg/utils"
	"github.com/blytz/shared/pkg/errors"
	"go.uber.org/zap"
)

type AuctionHandler struct {
	auctionService *services.AuctionService
	firebaseClient *firebase.Client
	logger         *zap.Logger
}

func NewAuctionHandler(auctionService *services.AuctionService, firebaseClient *firebase.Client, logger *zap.Logger) *AuctionHandler {
	return &AuctionHandler{
		auctionService: auctionService,
		firebaseClient: firebaseClient,
		logger:         logger,
	}
}

// CreateAuction creates a new auction with Firebase persistence
func (h *AuctionHandler) CreateAuction(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, errors.AuthenticationError("NO_USER_ID", "User ID not found"))
		return
	}

	var req models.CreateAuctionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, "Invalid request data", gin.H{"error": err.Error()})
		return
	}

	// Step 1: Create auction in Redis (primary storage - fast)
	auction, err := h.auctionService.CreateAuction(c.Request.Context(), userID.(string), &req)
	if err != nil {
		h.logger.Error("Failed to create auction in Redis", zap.Error(err))
		utils.ErrorResponse(c, err)
		return
	}

	// Step 2: Persist to Firebase (secondary storage - durable)
	go func() {
		ctx := context.Background()
		firebaseData := firebase.AuctionData{
			Title:         req.Title,
			Description:   req.Description,
			StartingPrice: req.StartingPrice,
			Duration:      int(req.EndTime.Sub(req.StartTime).Hours()),
			Category:      "general", // Map from your product categories
			Images:        req.Images,
		}

		firebaseResp, err := h.firebaseClient.CreateAuction(ctx, firebaseData)
		if err != nil {
			h.logger.Error("Failed to persist auction to Firebase", zap.Error(err),
				zap.String("auction_id", auction.AuctionID),
				zap.String("user_id", userID.(string)))
		} else {
			h.logger.Info("Auction persisted to Firebase successfully",
				zap.String("auction_id", auction.AuctionID),
				zap.String("firebase_id", firebaseResp.AuctionID))
		}
	}()

	// Return immediately - don't wait for Firebase (async for performance)
	utils.SuccessResponse(c, models.AuctionResponse{Auction: *auction})
}

// GetAuction retrieves an auction by ID
func (h *AuctionHandler) GetAuction(c *gin.Context) {
	auctionID := c.Param("auction_id")
	if auctionID == "" {
		utils.ErrorResponse(c, errors.ValidationError("MISSING_AUCTION_ID", "Auction ID is required"))
		return
	}

	auction, err := h.auctionService.GetAuction(c.Request.Context(), auctionID)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, models.AuctionResponse{Auction: *auction})
}

// ListAuctions retrieves a list of auctions
func (h *AuctionHandler) ListAuctions(c *gin.Context) {
	status := c.Query("status")
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	response, err := h.auctionService.ListAuctions(c.Request.Context(), status, page, limit)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, response)
}

// GetActiveAuctions retrieves only active auctions
func (h *AuctionHandler) GetActiveAuctions(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	// Call ListAuctions with "active" status
	response, err := h.auctionService.ListAuctions(c.Request.Context(), "active", page, limit)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, response)
}

// UpdateAuction updates an auction
func (h *AuctionHandler) UpdateAuction(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, errors.AuthenticationError("NO_USER_ID", "User ID not found"))
		return
	}

	auctionID := c.Param("auction_id")
	if auctionID == "" {
		utils.ErrorResponse(c, errors.ValidationError("MISSING_AUCTION_ID", "Auction ID is required"))
		return
	}

	var req models.UpdateAuctionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, "Invalid request data", gin.H{"error": err.Error()})
		return
	}

	auction, err := h.auctionService.UpdateAuction(c.Request.Context(), auctionID, userID.(string), &req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, models.AuctionResponse{Auction: *auction})
}

// DeleteAuction deletes an auction
func (h *AuctionHandler) DeleteAuction(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, errors.AuthenticationError("NO_USER_ID", "User ID not found"))
		return
	}

	auctionID := c.Param("auction_id")
	if auctionID == "" {
		utils.ErrorResponse(c, errors.ValidationError("MISSING_AUCTION_ID", "Auction ID is required"))
		return
	}

	if err := h.auctionService.DeleteAuction(c.Request.Context(), auctionID, userID.(string)); err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "Auction deleted successfully"})
}

// PlaceBid places a bid on an auction with Firebase persistence
func (h *AuctionHandler) PlaceBid(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, errors.AuthenticationError("NO_USER_ID", "User ID not found"))
		return
	}

	auctionID := c.Param("auction_id")
	if auctionID == "" {
		utils.ErrorResponse(c, errors.ValidationError("MISSING_AUCTION_ID", "Auction ID is required"))
		return
	}

	var req models.PlaceBidRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, "Invalid request data", gin.H{"error": err.Error()})
		return
	}

	// Step 1: Place bid in Redis (primary storage - fast)
	bid, err := h.auctionService.PlaceBid(c.Request.Context(), auctionID, userID.(string), req.Amount)
	if err != nil {
		h.logger.Error("Failed to place bid in Redis", zap.Error(err))
		utils.ErrorResponse(c, err)
		return
	}

	// Step 2: Persist bid to Firebase (secondary storage - durable)
	go func() {
		ctx := context.Background()
		bidData := firebase.BidData{
			AuctionID: auctionID,
			Amount:    req.Amount,
		}

		bidResp, err := h.firebaseClient.PlaceBid(ctx, bidData)
		if err != nil {
			h.logger.Error("Failed to persist bid to Firebase", zap.Error(err),
				zap.String("auction_id", auctionID),
				zap.String("user_id", userID.(string)),
				zap.Float64("amount", req.Amount))
		} else {
			h.logger.Info("Bid persisted to Firebase successfully",
				zap.String("auction_id", auctionID),
				zap.String("bid_id", bidResp.BidID),
				zap.Float64("amount", req.Amount))
		}
	}()

	// Return immediately - don't wait for Firebase (async for performance)
	utils.SuccessResponse(c, models.BidResponse{Bid: *bid})
}

// GetBids retrieves bids for an auction
func (h *AuctionHandler) GetBids(c *gin.Context) {
	auctionID := c.Param("auction_id")
	if auctionID == "" {
		utils.ErrorResponse(c, errors.ValidationError("MISSING_AUCTION_ID", "Auction ID is required"))
		return
	}

	response, err := h.auctionService.GetBids(c.Request.Context(), auctionID)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, response)
}

// GetAuctionStatus retrieves the current status of an auction
func (h *AuctionHandler) GetAuctionStatus(c *gin.Context) {
	auctionID := c.Param("auction_id")
	if auctionID == "" {
		utils.ErrorResponse(c, errors.ValidationError("MISSING_AUCTION_ID", "Auction ID is required"))
		return
	}

	status, err := h.auctionService.GetAuctionStatus(c.Request.Context(), auctionID)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, status)
}