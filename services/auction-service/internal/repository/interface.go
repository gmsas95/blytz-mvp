package repository

import (
	"context"
	"database/sql"
	"github.com/gmsas95/blytz-mvp/services/auction-service/internal/models"
)

// DBTX is an interface that can represent either a *sql.DB or a *sql.Tx
type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

type AuctionRepo interface {
	Create(ctx context.Context, auction *models.Auction) error
	GetByID(ctx context.Context, id string) (*models.Auction, error)
	UpdateAuctionPrice(ctx context.Context, id string, price float64) error
	CreateBid(ctx context.Context, bid *models.Bid) error
	BeginTx(ctx context.Context) (*sql.Tx, error)
	WithTx(tx *sql.Tx) AuctionRepo
}