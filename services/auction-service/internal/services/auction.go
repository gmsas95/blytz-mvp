package services

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/gmsas95/blytz-mvp/services/auction-service/internal/models"
	"github.com/gmsas95/blytz-mvp/services/auction-service/internal/config"
	"github.com/gmsas95/blytz-mvp/services/auction-service/internal/repository"
	shared_errors "github.com/gmsas95/blytz-mvp/shared/pkg/errors"
)


type AuctionService struct {
	repo   repository.AuctionRepo
	logger *zap.Logger
	config *config.Config
}

func NewAuctionService(repo repository.AuctionRepo, logger *zap.Logger, config *config.Config) *AuctionService {
	return &AuctionService{repo: repo, logger: logger, config: config}
}

func (s *AuctionService) CreateAuction(ctx context.Context, auction *models.Auction) error {
	if auction.StartTime.IsZero() {
		auction.StartTime = time.Now()
	}
	if auction.EndTime.IsZero() {
		auction.EndTime = auction.StartTime.Add(s.config.DefaultAuctionDuration)
	}

	if err := s.repo.Create(ctx, auction); err != nil {
		s.logger.Error("Failed to create auction", zap.Error(err))
		return shared_errors.ErrInternalServer
	}
	return nil
}

func (s *AuctionService) GetAuction(ctx context.Context, id string) (*models.Auction, error) {
	auction, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get auction", zap.String("auction_id", id), zap.Error(err))
		return nil, shared_errors.ErrNotFound
	}
	return auction, nil
}

func (s *AuctionService) PlaceBid(ctx context.Context, bid *models.Bid) error {
	auction, err := s.repo.GetByID(ctx, bid.AuctionID)
	if err != nil {
		s.logger.Error("Failed to get auction for bid", zap.String("auction_id", bid.AuctionID), zap.Error(err))
		return shared_errors.ErrNotFound
	}

	if time.Now().After(auction.EndTime) {
		return shared_errors.ErrAuctionEnded
	}

	if bid.Amount <= auction.CurrentPrice {
		return shared_errors.ErrBidTooLow
	}

	if err := s.repo.CreateBid(ctx, bid); err != nil {
		s.logger.Error("Failed to place bid", zap.Any("bid", bid), zap.Error(err))
		return shared_errors.ErrInternalServer
	}

	if err := s.repo.UpdateAuctionPrice(ctx, bid.AuctionID, bid.Amount); err != nil {
		s.logger.Error("Failed to update auction price", zap.String("auction_id", bid.AuctionID), zap.Float64("new_price", bid.Amount), zap.Error(err))
		// This is an internal error, the bid was placed but the price update failed.
		// The system should handle this inconsistency, for now we log it.
	}

	return nil
}