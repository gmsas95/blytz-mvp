package services

import (
	"context"
	"database/sql"
	"time"

	"go.uber.org/zap"

	"github.com/gmsas95/blytz-mvp/services/auction-service/internal/models"
	"github.com/gmsas95/blytz-mvp/services/auction-service/internal/config"
	"github.com/gmsas95/blytz-mvp/services/auction-service/internal/repository"
	shared_errors "github.com/gmsas95/blytz-mvp/shared/pkg/errors"
)


type AuctionService struct {
	db     *sql.DB
	logger *zap.Logger
	config *config.Config
}

func NewAuctionService(db *sql.DB, logger *zap.Logger, config *config.Config) *AuctionService {
	return &AuctionService{db: db, logger: logger, config: config}
}

func (s *AuctionService) CreateAuction(ctx context.Context, auction *models.Auction) error {
	if auction.StartTime.IsZero() {
		auction.StartTime = time.Now()
	}
	if auction.EndTime.IsZero() {
		auction.EndTime = auction.StartTime.Add(s.config.DefaultAuctionDuration)
	}

	repo := repository.NewPostgresRepo(s.db, s.logger)
	if err := repo.Create(ctx, auction); err != nil {
		s.logger.Error("Failed to create auction", zap.Error(err))
		return shared_errors.ErrInternalServer
	}
	return nil
}

func (s *AuctionService) GetAuction(ctx context.Context, id string) (*models.Auction, error) {
	repo := repository.NewPostgresRepo(s.db, s.logger)
	auction, err := repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get auction", zap.String("auction_id", id), zap.Error(err))
		return nil, shared_errors.ErrNotFound
	}
	return auction, nil
}

func (s *AuctionService) PlaceBid(ctx context.Context, bid *models.Bid) error {
	// Create a new repository instance for the transaction
	repo := repository.NewPostgresRepo(s.db, s.logger)

	// Begin a transaction
	tx, err := repo.BeginTx(ctx)
	if err != nil {
		s.logger.Error("Failed to begin transaction", zap.Error(err))
		return shared_errors.ErrInternalServer
	}
	defer tx.Rollback() // Rollback is a no-op if the transaction is already committed

	// Create a new repository with the transaction
	txRepo := repo.WithTx(tx)

	auction, err := txRepo.GetByID(ctx, bid.AuctionID)
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

	if err := txRepo.CreateBid(ctx, bid); err != nil {
		s.logger.Error("Failed to place bid", zap.Any("bid", bid), zap.Error(err))
		return shared_errors.ErrInternalServer
	}

	if err := txRepo.UpdateAuctionPrice(ctx, bid.AuctionID, bid.Amount); err != nil {
		s.logger.Error("Failed to update auction price", zap.String("auction_id", bid.AuctionID), zap.Float64("new_price", bid.Amount), zap.Error(err))
		return shared_errors.ErrInternalServer
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		s.logger.Error("Failed to commit transaction", zap.Error(err))
		return shared_errors.ErrInternalServer
	}

	return nil
}