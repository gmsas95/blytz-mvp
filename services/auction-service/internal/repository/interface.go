package repository

import (
	"context"
	"github.com/gmsas95/blytz-mvp/services/auction-service/internal/models"
)


type AuctionRepo interface {
	Create(ctx context.Context, auction *models.Auction) error
	GetByID(ctx context.Context, id string) (*models.Auction, error)
	UpdateAuctionPrice(ctx context.Context, id string, price float64) error
	CreateBid(ctx context.Context, bid *models.Bid) error
}