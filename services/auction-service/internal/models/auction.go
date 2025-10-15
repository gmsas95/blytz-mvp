package models

import (
	"time"
)

type Auction struct {
	AuctionID       string          `json:"auction_id" gorm:"primaryKey"`
	ProductID       string          `json:"product_id"`
	SellerID        string          `json:"seller_id"`
	Title           string          `json:"title"`
	Description     string          `json:"description"`
	StartingPrice   float64         `json:"starting_price"`
	CurrentPrice    float64         `json:"current_price"`
	ReservePrice    float64         `json:"reserve_price"`
	MinBidIncrement float64         `json:"min_bid_increment"`
	StartTime       time.Time       `json:"start_time"`
	EndTime         time.Time       `json:"end_time"`
	Status          string          `json:"status"` // scheduled, active, ended, cancelled
	Type            string          `json:"type"`   // live, scheduled
	IsActive        bool            `json:"is_active"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	Bids            []Bid           `json:"bids,omitempty" gorm:"foreignKey:AuctionID"`
	Images          []AuctionImage  `json:"images,omitempty" gorm:"foreignKey:AuctionID"`
}

type Bid struct {
	BidID      string    `json:"bid_id" gorm:"primaryKey"`
	AuctionID  string    `json:"auction_id"`
	BidderID   string    `json:"bidder_id"`
	Amount     float64   `json:"amount"`
	IsWinning  bool      `json:"is_winning"`
	BidTime    time.Time `json:"bid_time"`
	CreatedAt  time.Time `json:"created_at"`
}

type AuctionImage struct {
	ImageID   string `json:"image_id" gorm:"primaryKey"`
	AuctionID string `json:"auction_id"`
	ImageURL  string `json:"image_url"`
	AltText   string `json:"alt_text"`
	Order     int    `json:"order"`
}

type CreateAuctionRequest struct {
	ProductID       string    `json:"product_id" binding:"required"`
	Title           string    `json:"title" binding:"required"`
	Description     string    `json:"description"`
	StartingPrice   float64   `json:"starting_price" binding:"required,min=0"`
	ReservePrice    float64   `json:"reserve_price"`
	MinBidIncrement float64   `json:"min_bid_increment" binding:"required,min=0"`
	StartTime       time.Time `json:"start_time" binding:"required"`
	EndTime         time.Time `json:"end_time" binding:"required"`
	Type            string    `json:"type" binding:"required,oneof=live scheduled"`
	Images          []string  `json:"images"`
}

type UpdateAuctionRequest struct {
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	ReservePrice    float64   `json:"reserve_price"`
	MinBidIncrement float64   `json:"min_bid_increment"`
	StartTime       time.Time `json:"start_time"`
	EndTime         time.Time `json:"end_time"`
	Images          []string  `json:"images"`
}

type PlaceBidRequest struct {
	Amount float64 `json:"amount" binding:"required,min=0"`
}

type AuctionResponse struct {
	Auction Auction `json:"auction"`
}

type BidResponse struct {
	Bid   Bid   `json:"bid"`
}

type BidsResponse struct {
	Bids []Bid `json:"bids"`
}

type AuctionsResponse struct {
	Auctions []Auction `json:"auctions"`
	Total    int64     `json:"total"`
	Page     int       `json:"page"`
	Limit    int       `json:"limit"`
}

type AuctionStatus struct {
	AuctionID     string    `json:"auction_id"`
	Status        string    `json:"status"`
	CurrentPrice  float64   `json:"current_price"`
	WinningBidID  string    `json:"winning_bid_id,omitempty"`
	TotalBids     int       `json:"total_bids"`
	TimeRemaining string    `json:"time_remaining"`
	UpdatedAt     time.Time `json:"updated_at"`
}