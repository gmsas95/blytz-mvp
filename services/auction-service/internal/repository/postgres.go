package repository

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/gmsas95/blytz-mvp/services/auction-service/internal/models"
)


type PostgresRepo struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewPostgresRepo(db *sql.DB, logger *zap.Logger) *PostgresRepo {
	return &PostgresRepo{db: db, logger: logger}
}

func (r *PostgresRepo) Create(ctx context.Context, auction *models.Auction) error {
	query := `INSERT INTO auctions (id, item, start_time, end_time, start_price, created_by) VALUES (, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query, auction.ID, auction.Item, auction.StartTime, auction.EndTime, auction.StartPrice, auction.CreatedBy)
	return err
}

func (r *PostgresRepo) GetByID(ctx context.Context, id string) (*models.Auction, error) {
	query := `SELECT id, item, start_time, end_time, start_price, current_price, created_by, winner_id, status FROM auctions WHERE id = `
	row := r.db.QueryRowContext(ctx, query, id)

	auction := &models.Auction{}
	if err := row.Scan(&auction.ID, &auction.Item, &auction.StartTime, &auction.EndTime, &auction.StartPrice, &auction.CurrentPrice, &auction.CreatedBy, &auction.WinnerID, &auction.Status); err != nil {
		return nil, err
	}
	return auction, nil
}

func (r *PostgresRepo) UpdateAuctionPrice(ctx context.Context, id string, price float64) error {
	query := `UPDATE auctions SET current_price =  WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, price, id)
	return err
}

func (r *PostgresRepo) CreateBid(ctx context.Context, bid *models.Bid) error {
	query := `INSERT INTO bids (id, auction_id, user_id, amount, timestamp) VALUES (, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query, bid.ID, bid.AuctionID, bid.UserID, bid.Amount, bid.Timestamp)
	return err
}


func (r *PostgresAuctionRepository) GetBids(ctx context.Context, auctionID string) (*models.BidsResponse, error) {
	r.logger.Info("Getting bids from database", zap.String("auction_id", auctionID))

	query := `
		SELECT bid_id, auction_id, bidder_id, amount, is_winning, bid_time, created_at
		FROM bids
		WHERE auction_id = $1
		ORDER BY bid_time DESC
	`

	rows, err := r.db.QueryContext(ctx, query, auctionID)
	if err != nil {
		r.logger.Error("Failed to query bids", zap.Error(err))
		return nil, fmt.Errorf("failed to query bids: %w", err)
	}
	defer rows.Close()

	var bids []models.Bid
	for rows.Next() {
		var bid models.Bid
		err := rows.Scan(
			&bid.BidID, &bid.AuctionID, &bid.BidderID, &bid.Amount, &bid.IsWinning, &bid.BidTime, &bid.CreatedAt,
		)
		if err != nil {
			r.logger.Error("Failed to scan bid", zap.Error(err))
			return nil, fmt.Errorf("failed to scan bid: %w", err)
		}
		bids = append(bids, bid)
	}

	return &models.BidsResponse{Bids: bids}, nil
}

func (r *PostgresAuctionRepository) GetWinningBid(ctx context.Context, auctionID string) (*models.Bid, error) {
	r.logger.Info("Getting winning bid from database", zap.String("auction_id", auctionID))

	query := `
		SELECT bid_id, auction_id, bidder_id, amount, is_winning, bid_time, created_at
		FROM bids
		WHERE auction_id = $1 AND is_winning = true
		ORDER BY bid_time DESC
		LIMIT 1
	`

	var bid models.Bid
	err := r.db.QueryRowContext(ctx, query, auctionID).Scan(
		&bid.BidID, &bid.AuctionID, &bid.BidderID, &bid.Amount, &bid.IsWinning, &bid.BidTime, &bid.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no winning bid found for auction: %s", auctionID)
		}
		r.logger.Error("Failed to get winning bid", zap.Error(err))
		return nil, fmt.Errorf("failed to get winning bid: %w", err)
	}

	return &bid, nil
}

func (r *PostgresAuctionRepository) UpdateAuctionStatus(ctx context.Context, auctionID string, status string) error {
	r.logger.Info("Updating auction status in database", zap.String("auction_id", auctionID), zap.String("status", status))

	query := "UPDATE auctions SET status = $2, updated_at = $3 WHERE auction_id = $1"
	_, err := r.db.ExecContext(ctx, query, auctionID, status, time.Now())

	if err != nil {
		r.logger.Error("Failed to update auction status", zap.Error(err))
		return fmt.Errorf("failed to update auction status: %w", err)
	}

	return nil
}

func (r *PostgresAuctionRepository) GetActiveAuctions(ctx context.Context) ([]models.Auction, error) {
	r.logger.Info("Getting active auctions from database")

	query := `
		SELECT auction_id, product_id, seller_id, title, description,
			starting_price, current_price, reserve_price, min_bid_increment,
			start_time, end_time, status, type, is_active, created_at, updated_at
		FROM auctions
		WHERE status = 'active' AND end_time > $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, time.Now())
	if err != nil {
		r.logger.Error("Failed to query active auctions", zap.Error(err))
		return nil, fmt.Errorf("failed to query active auctions: %w", err)
	}
	defer rows.Close()

	var auctions []models.Auction
	for rows.Next() {
		var auction models.Auction
		err := rows.Scan(
			&auction.AuctionID, &auction.ProductID, &auction.SellerID, &auction.Title, &auction.Description,
			&auction.StartingPrice, &auction.CurrentPrice, &auction.ReservePrice, &auction.MinBidIncrement,
			&auction.StartTime, &auction.EndTime, &auction.Status, &auction.Type, &auction.IsActive,
			&auction.CreatedAt, &auction.UpdatedAt,
		)
		if err != nil {
			r.logger.Error("Failed to scan auction", zap.Error(err))
			return nil, fmt.Errorf("failed to scan auction: %w", err)
		}
		auctions = append(auctions, auction)
	}

	return auctions, nil
}

func (r *PostgresAuctionRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

func (r *PostgresAuctionRepository) Close() error {
	return r.db.Close()
}