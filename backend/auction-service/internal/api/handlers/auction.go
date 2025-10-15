package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/blytz/auction-service/internal/models"
	"github.com/blytz/auction-service/internal/services"
	"github.com/blytz/shared/utils"
	"github.com/blytz/shared/errors"
)

type AuctionHandler struct {
	auctionService *services.AuctionService
}

func NewAuctionHandler(auctionService *services.AuctionService) *AuctionHandler {
	return &AuctionHandler{auctionService: auctionService}
}

// CreateAuction creates a new auction
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

	auction, err := h.auctionService.CreateAuction(c.Request.Context(), userID.(string), &req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

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

// PlaceBid places a bid on an auction
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

	bid, err := h.auctionService.PlaceBid(c.Request.Context(), auctionID, userID.(string), req.Amount)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

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