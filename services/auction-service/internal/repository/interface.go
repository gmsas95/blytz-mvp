package repository

import (
	"context"
	"github.com/blytz/auction-service/internal/models"
)

// AuctionRepository defines the interface for auction data operations
type AuctionRepository interface {
	// Auction operations
	CreateAuction(ctx context.Context, auction *models.Auction) error
	GetAuction(ctx context.Context, auctionID string) (*models.Auction, error)
	GetAuctions(ctx context.Context, status string, page, limit int) (*models.AuctionsResponse, error)
	UpdateAuction(ctx context.Context, auction *models.Auction) error
	DeleteAuction(ctx context.Context, auctionID string) error

	// Bid operations
	CreateBid(ctx context.Context, bid *models.Bid) error
	GetBids(ctx context.Context, auctionID string) (*models.BidsResponse, error)
	GetWinningBid(ctx context.Context, auctionID string) (*models.Bid, error)

	// Auction status operations
	UpdateAuctionStatus(ctx context.Context, auctionID string, status string) error
	GetActiveAuctions(ctx context.Context) ([]models.Auction, error)

	// Database operations
	Ping(ctx context.Context) error
	Close() error
}