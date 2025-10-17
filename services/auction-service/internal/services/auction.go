package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/blytz/auction-service/internal/models"
	"github.com/blytz/auction-service/internal/config"
	"github.com/blytz/auction-service/internal/repository"
	"github.com/blytz/shared/pkg/errors"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type AuctionService struct {
	logger *zap.Logger
	config *config.Config
	repo   repository.AuctionRepository
}

func NewAuctionService(logger *zap.Logger, config *config.Config) *AuctionService {
	return &AuctionService{
		logger: logger,
		config: config,
		repo:   nil, // Will be set after database connection
	}
}

func NewAuctionServiceWithRepo(logger *zap.Logger, config *config.Config, repo repository.AuctionRepository) *AuctionService {
	return &AuctionService{
		logger: logger,
		config: config,
		repo:   repo,
	}
}

// CreateAuction creates a new auction
func (s *AuctionService) CreateAuction(ctx context.Context, sellerID string, req *models.CreateAuctionRequest) (*models.Auction, error) {
	s.logger.Info("CreateAuction called", zap.String("seller_id", sellerID))

	// CRITICAL: No fallback to mock data - database persistence is required
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

	if err := s.repo.CreateAuction(ctx, auction); err != nil {
		s.logger.Error("Failed to create auction", zap.Error(err))
		return nil, err
	}

	return auction, nil
}

// GetAuction retrieves an auction by ID
func (s *AuctionService) GetAuction(ctx context.Context, auctionID string) (*models.Auction, error) {
	s.logger.Info("GetAuction called", zap.String("auction_id", auctionID))

	// CRITICAL: No fallback to mock data - database persistence is required
	auction, err := s.repo.GetAuction(ctx, auctionID)
	if err != nil {
		s.logger.Error("Failed to get auction", zap.Error(err))
		return nil, err
	}

	return auction, nil
}

// ListAuctions retrieves a list of auctions
func (s *AuctionService) ListAuctions(ctx context.Context, status string, page, limit int) (*models.AuctionsResponse, error) {
	s.logger.Info("ListAuctions called", zap.String("status", status), zap.Int("page", page), zap.Int("limit", limit))

	// CRITICAL: No fallback to mock data - database persistence is required
	response, err := s.repo.GetAuctions(ctx, status, page, limit)
	if err != nil {
		s.logger.Error("Failed to list auctions", zap.Error(err))
		return nil, err
	}

	return response, nil
}

// UpdateAuction updates an auction
func (s *AuctionService) UpdateAuction(ctx context.Context, auctionID string, sellerID string, req *models.UpdateAuctionRequest) (*models.Auction, error) {
	s.logger.Info("UpdateAuction called", zap.String("auction_id", auctionID), zap.String("seller_id", sellerID))

	// Get existing auction to validate ownership
	auction, err := s.repo.GetAuction(ctx, auctionID)
	if err != nil {
		s.logger.Error("Failed to get auction for update", zap.Error(err))
		return nil, err
	}

	// Verify seller ownership
	if auction.SellerID != sellerID {
		return nil, errors.ValidationError("UNAUTHORIZED", "Only the auction seller can update the auction")
	}

	// Only allow updates if auction is not ended
	if auction.Status == "ended" {
		return nil, errors.ValidationError("AUCTION_ENDED", "Cannot update ended auction")
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

	// Save to database
	if err := s.repo.UpdateAuction(ctx, auction); err != nil {
		s.logger.Error("Failed to update auction", zap.Error(err))
		return nil, err
	}

	return auction, nil
}

// DeleteAuction deletes an auction
func (s *AuctionService) DeleteAuction(ctx context.Context, auctionID string, sellerID string) error {
	s.logger.Info("DeleteAuction called", zap.String("auction_id", auctionID), zap.String("seller_id", sellerID))

	// Get existing auction to validate ownership
	auction, err := s.repo.GetAuction(ctx, auctionID)
	if err != nil {
		s.logger.Error("Failed to get auction for deletion", zap.Error(err))
		return err
	}

	// Verify seller ownership
	if auction.SellerID != sellerID {
		return errors.ValidationError("UNAUTHORIZED", "Only the auction seller can delete the auction")
	}

	// Only allow deletion if auction has no bids (business rule)
	bidsResponse, err := s.repo.GetBids(ctx, auctionID)
	if err != nil {
		s.logger.Error("Failed to check auction bids", zap.Error(err))
		return err
	}

	if len(bidsResponse.Bids) > 0 {
		return errors.ValidationError("AUCTION_HAS_BIDS", "Cannot delete auction that has bids")
	}

	// Delete from database
	if err := s.repo.DeleteAuction(ctx, auctionID); err != nil {
		s.logger.Error("Failed to delete auction", zap.Error(err))
		return err
	}

	return nil
}

// PlaceBid places a bid on an auction
func (s *AuctionService) PlaceBid(ctx context.Context, auctionID string, bidderID string, amount float64) (*models.Bid, error) {
	s.logger.Info("PlaceBid called", zap.String("auction_id", auctionID), zap.String("bidder_id", bidderID), zap.Float64("amount", amount))

	// CRITICAL: Database persistence is required - no fallback to mock data
	// Get auction to validate
	auction, err := s.repo.GetAuction(ctx, auctionID)
	if err != nil {
		s.logger.Error("Failed to get auction for bid validation", zap.Error(err))
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

	// Create bid
	bid := &models.Bid{
		BidID:     generateBidID(),
		AuctionID: auctionID,
		BidderID:  bidderID,
		Amount:    amount,
		IsWinning: true,
		BidTime:   time.Now(),
		CreatedAt: time.Now(),
	}

	if err := s.repo.CreateBid(ctx, bid); err != nil {
		s.logger.Error("Failed to create bid", zap.Error(err))
		return nil, err
	}

	// Update auction current price
	auction.CurrentPrice = amount
	auction.UpdatedAt = time.Now()
	if err := s.repo.UpdateAuction(ctx, auction); err != nil {
		s.logger.Error("Failed to update auction after bid", zap.Error(err))
		return nil, err
	}

	return bid, nil
}

// GetBids retrieves bids for an auction
func (s *AuctionService) GetBids(ctx context.Context, auctionID string) (*models.BidsResponse, error) {
	s.logger.Info("GetBids called", zap.String("auction_id", auctionID))

	// CRITICAL: Database persistence is required - no fallback to mock data
	response, err := s.repo.GetBids(ctx, auctionID)
	if err != nil {
		s.logger.Error("Failed to get bids", zap.Error(err))
		return nil, err
	}

	return response, nil
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