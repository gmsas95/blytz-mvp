package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/blytz/auction-service/internal/models"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type PostgresAuctionRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewPostgresAuctionRepository(db *sql.DB, logger *zap.Logger) AuctionRepository {
	return &PostgresAuctionRepository{
		db:     db,
		logger: logger,
	}
}

func (r *PostgresAuctionRepository) CreateAuction(ctx context.Context, auction *models.Auction) error {
	r.logger.Info("Creating auction in database", zap.String("auction_id", auction.AuctionID))

	query := `
		INSERT INTO auctions (
			auction_id, product_id, seller_id, title, description,
			starting_price, current_price, reserve_price, min_bid_increment,
			start_time, end_time, status, type, is_active, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	`

	_, err := r.db.ExecContext(ctx, query,
		auction.AuctionID, auction.ProductID, auction.SellerID, auction.Title, auction.Description,
		auction.StartingPrice, auction.CurrentPrice, auction.ReservePrice, auction.MinBidIncrement,
		auction.StartTime, auction.EndTime, auction.Status, auction.Type, auction.IsActive,
		auction.CreatedAt, auction.UpdatedAt,
	)

	if err != nil {
		r.logger.Error("Failed to create auction", zap.Error(err))
		return fmt.Errorf("failed to create auction: %w", err)
	}

	return nil
}

func (r *PostgresAuctionRepository) GetAuction(ctx context.Context, auctionID string) (*models.Auction, error) {
	r.logger.Info("Getting auction from database", zap.String("auction_id", auctionID))

	query := `
		SELECT auction_id, product_id, seller_id, title, description,
			starting_price, current_price, reserve_price, min_bid_increment,
			start_time, end_time, status, type, is_active, created_at, updated_at
		FROM auctions
		WHERE auction_id = $1
	`

	var auction models.Auction
	err := r.db.QueryRowContext(ctx, query, auctionID).Scan(
		&auction.AuctionID, &auction.ProductID, &auction.SellerID, &auction.Title, &auction.Description,
		&auction.StartingPrice, &auction.CurrentPrice, &auction.ReservePrice, &auction.MinBidIncrement,
		&auction.StartTime, &auction.EndTime, &auction.Status, &auction.Type, &auction.IsActive,
		&auction.CreatedAt, &auction.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("auction not found: %s", auctionID)
		}
		r.logger.Error("Failed to get auction", zap.Error(err))
		return nil, fmt.Errorf("failed to get auction: %w", err)
	}

	return &auction, nil
}

func (r *PostgresAuctionRepository) GetAuctions(ctx context.Context, status string, page, limit int) (*models.AuctionsResponse, error) {
	r.logger.Info("Getting auctions from database", zap.String("status", status), zap.Int("page", page), zap.Int("limit", limit))

	offset := (page - 1) * limit

	query := `
		SELECT auction_id, product_id, seller_id, title, description,
			starting_price, current_price, reserve_price, min_bid_increment,
			start_time, end_time, status, type, is_active, created_at, updated_at
		FROM auctions
	`

	var args []interface{}
	var whereClause string

	if status != "" {
		whereClause = " WHERE status = $1"
		args = append(args, status)
	}

	query += whereClause + " ORDER BY created_at DESC LIMIT $" + fmt.Sprintf("%d", len(args)+1) + " OFFSET $" + fmt.Sprintf("%d", len(args)+2)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		r.logger.Error("Failed to query auctions", zap.Error(err))
		return nil, fmt.Errorf("failed to query auctions: %w", err)
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

	// Get total count
	var total int
	countQuery := "SELECT COUNT(*) FROM auctions"
	if whereClause != "" {
		countQuery += whereClause
	}
	err = r.db.QueryRowContext(ctx, countQuery, args[:len(args)-2]...).Scan(&total)
	if err != nil {
		r.logger.Error("Failed to count auctions", zap.Error(err))
		return nil, fmt.Errorf("failed to count auctions: %w", err)
	}

	return &models.AuctionsResponse{
		Auctions: auctions,
		Total:    total,
		Page:     page,
		Limit:    limit,
	}, nil
}

func (r *PostgresAuctionRepository) UpdateAuction(ctx context.Context, auction *models.Auction) error {
	r.logger.Info("Updating auction in database", zap.String("auction_id", auction.AuctionID))

	query := `
		UPDATE auctions SET
			title = $2, description = $3, current_price = $4, reserve_price = $5,
			min_bid_increment = $6, start_time = $7, end_time = $8, status = $9,
			updated_at = $10
		WHERE auction_id = $1
	`

	_, err := r.db.ExecContext(ctx, query,
		auction.AuctionID, auction.Title, auction.Description,
		auction.CurrentPrice, auction.ReservePrice, auction.MinBidIncrement,
		auction.StartTime, auction.EndTime, auction.Status, auction.UpdatedAt,
	)

	if err != nil {
		r.logger.Error("Failed to update auction", zap.Error(err))
		return fmt.Errorf("failed to update auction: %w", err)
	}

	return nil
}

func (r *PostgresAuctionRepository) DeleteAuction(ctx context.Context, auctionID string) error {
	r.logger.Info("Deleting auction from database", zap.String("auction_id", auctionID))

	query := "DELETE FROM auctions WHERE auction_id = $1"
	_, err := r.db.ExecContext(ctx, query, auctionID)

	if err != nil {
		r.logger.Error("Failed to delete auction", zap.Error(err))
		return fmt.Errorf("failed to delete auction: %w", err)
	}

	return nil
}

func (r *PostgresAuctionRepository) CreateBid(ctx context.Context, bid *models.Bid) error {
	r.logger.Info("Creating bid in database", zap.String("bid_id", bid.BidID), zap.String("auction_id", bid.AuctionID))

	query := `
		INSERT INTO bids (bid_id, auction_id, bidder_id, amount, is_winning, bid_time, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(ctx, query,
		bid.BidID, bid.AuctionID, bid.BidderID, bid.Amount, bid.IsWinning, bid.BidTime, bid.CreatedAt,
	)

	if err != nil {
		r.logger.Error("Failed to create bid", zap.Error(err))
		return fmt.Errorf("failed to create bid: %w", err)
	}

	return nil
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