package services

import (
	"context"
	"fmt"
	"time"

	"github.com/blytz/auction-service/internal/models"
	"github.com/blytz/auction-service/internal/config"
	"github.com/blytz/shared/errors"
	"go.uber.org/zap"
)

type AuctionService struct {
	logger *zap.Logger
	config *config.Config
}

func NewAuctionService(logger *zap.Logger, config *config.Config) *AuctionService {
	return &AuctionService{
		logger: logger,
		config: config,
	}
}

// CreateAuction creates a new auction
func (s *AuctionService) CreateAuction(ctx context.Context, sellerID string, req *models.CreateAuctionRequest) (*models.Auction, error) {
	s.logger.Info("CreateAuction called", zap.String("seller_id", sellerID))

	auction := &models.Auction{
		AuctionID:       generateAuctionID(),
		ProductID:       req.ProductID,
		SellerID:        sellerID,
		Title:           req.Title,
		Description:     req.Description,
		StartingPrice:   req.StartingPrice,
		CurrentPrice:    req.StartingPrice,
		ReservePrice:    req.ReservePrice,
		MinBidIncrement: req.MinBidIncrement,
		StartTime:       req.StartTime,
		EndTime:         req.EndTime,
		Status:          "scheduled",
		Type:            req.Type,
		IsActive:        true,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	return auction, nil
}

// GetAuction retrieves an auction by ID
func (s *AuctionService) GetAuction(ctx context.Context, auctionID string) (*models.Auction, error) {
	s.logger.Info("GetAuction called", zap.String("auction_id", auctionID))

	// Mock auction data
	return &models.Auction{
		AuctionID:       auctionID,
		ProductID:       "product123",
		SellerID:        "seller123",
		Title:           "Sample Auction",
		Description:     "This is a sample auction",
		StartingPrice:   100.00,
		CurrentPrice:    150.00,
		ReservePrice:    200.00,
		MinBidIncrement: 10.00,
		StartTime:       time.Now().Add(-1 * time.Hour),
		EndTime:         time.Now().Add(1 * time.Hour),
		Status:          "active",
		Type:            "live",
		IsActive:        true,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}, nil
}

// ListAuctions retrieves a list of auctions
func (s *AuctionService) ListAuctions(ctx context.Context, status string, page, limit int) (*models.AuctionsResponse, error) {
	s.logger.Info("ListAuctions called", zap.String("status", status), zap.Int("page", page), zap.Int("limit", limit))

	// Mock data
	auctions := []models.Auction{
		{
			AuctionID:       "auction123",
			ProductID:       "product123",
			SellerID:        "seller123",
			Title:           "Sample Auction 1",
			Description:     "This is a sample auction",
			StartingPrice:   100.00,
			CurrentPrice:    150.00,
			ReservePrice:    200.00,
			MinBidIncrement: 10.00,
			StartTime:       time.Now().Add(-1 * time.Hour),
			EndTime:         time.Now().Add(1 * time.Hour),
			Status:          "active",
			Type:            "live",
			IsActive:        true,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
	}

	return &models.AuctionsResponse{
		Auctions: auctions,
		Total:    1,
		Page:     page,
		Limit:    limit,
	}, nil
}

// UpdateAuction updates an auction
func (s *AuctionService) UpdateAuction(ctx context.Context, auctionID string, sellerID string, req *models.UpdateAuctionRequest) (*models.Auction, error) {
	s.logger.Info("UpdateAuction called", zap.String("auction_id", auctionID), zap.String("seller_id", sellerID))

	// Get existing auction
	auction, err := s.GetAuction(ctx, auctionID)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Title != "" {
		auction.Title = req.Title
	}
	if req.Description != "" {
		auction.Description = req.Description
	}
	if req.ReservePrice > 0 {
		auction.ReservePrice = req.ReservePrice
	}
	if req.MinBidIncrement > 0 {
		auction.MinBidIncrement = req.MinBidIncrement
	}
	if !req.StartTime.IsZero() {
		auction.StartTime = req.StartTime
	}
	if !req.EndTime.IsZero() {
		auction.EndTime = req.EndTime
	}

	auction.UpdatedAt = time.Now()
	return auction, nil
}

// DeleteAuction deletes an auction
func (s *AuctionService) DeleteAuction(ctx context.Context, auctionID string, sellerID string) error {
	s.logger.Info("DeleteAuction called", zap.String("auction_id", auctionID), zap.String("seller_id", sellerID))

	// In a real implementation, you would check if the auction exists and belongs to the seller
	return nil
}

// PlaceBid places a bid on an auction
func (s *AuctionService) PlaceBid(ctx context.Context, auctionID string, bidderID string, amount float64) (*models.Bid, error) {
	s.logger.Info("PlaceBid called", zap.String("auction_id", auctionID), zap.String("bidder_id", bidderID), zap.Float64("amount", amount))

	// Get auction
	auction, err := s.GetAuction(ctx, auctionID)
	if err != nil {
		return nil, err
	}

	// Validate bid
	if auction.Status != "active" {
		return nil, errors.ValidationError("AUCTION_NOT_ACTIVE", "Auction is not active")
	}

	if time.Now().After(auction.EndTime) {
		return nil, errors.ValidationError("AUCTION_ENDED", "Auction has ended")
	}

	if amount < auction.CurrentPrice+auction.MinBidIncrement {
		return nil, errors.ValidationError("BID_TOO_LOW", fmt.Sprintf("Bid must be at least $%.2f", auction.CurrentPrice+auction.MinBidIncrement))
	}

	bid := &models.Bid{
		BidID:     generateBidID(),
		AuctionID: auctionID,
		BidderID:  bidderID,
		Amount:    amount,
		IsWinning: true,
		BidTime:   time.Now(),
		CreatedAt: time.Now(),
	}

	// Update auction current price
	auction.CurrentPrice = amount
	auction.UpdatedAt = time.Now()

	return bid, nil
}

// GetBids retrieves bids for an auction
func (s *AuctionService) GetBids(ctx context.Context, auctionID string) (*models.BidsResponse, error) {
	s.logger.Info("GetBids called", zap.String("auction_id", auctionID))

	// Mock bids data
	bids := []models.Bid{
		{
			BidID:     "bid123",
			AuctionID: auctionID,
			BidderID:  "bidder123",
			Amount:    150.00,
			IsWinning: true,
			BidTime:   time.Now(),
			CreatedAt: time.Now(),
		},
	}

	return &models.BidsResponse{Bids: bids}, nil
}

// GetAuctionStatus retrieves the current status of an auction
func (s *AuctionService) GetAuctionStatus(ctx context.Context, auctionID string) (*models.AuctionStatus, error) {
	s.logger.Info("GetAuctionStatus called", zap.String("auction_id", auctionID))

	auction, err := s.GetAuction(ctx, auctionID)
	if err != nil {
		return nil, err
	}

	var timeRemaining string
	if time.Now().Before(auction.EndTime) {
		duration := auction.EndTime.Sub(time.Now())
		timeRemaining = duration.String()
	} else {
		timeRemaining = "Auction ended"
	}

	return &models.AuctionStatus{
		AuctionID:     auctionID,
		Status:        auction.Status,
		CurrentPrice:  auction.CurrentPrice,
		TotalBids:     1, // Mock data
		TimeRemaining: timeRemaining,
		UpdatedAt:     auction.UpdatedAt,
	}, nil
}

// Helper functions
func generateAuctionID() string {
	return "auction_" + time.Now().Format("20060102150405") + "_" + generateRandomString(8)
}

func generateBidID() string {
	return "bid_" + time.Now().Format("20060102150405") + "_" + generateRandomString(8)
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}