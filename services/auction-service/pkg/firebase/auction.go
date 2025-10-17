package firebase

import (
	"context"
	"fmt"
)

// AuctionData represents auction creation request
type AuctionData struct {
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	StartingPrice float64  `json:"startingPrice"`
	Duration      int      `json:"duration"`
	Category      string   `json:"category"`
	Images        []string `json:"images,omitempty"`
}

// AuctionResponse represents auction creation response
type AuctionResponse struct {
	Success   bool   `json:"success"`
	AuctionID string `json:"auctionId"`
	Message   string `json:"message"`
}

// BidData represents bid placement request
type BidData struct {
	AuctionID string  `json:"auctionId"`
	Amount    float64 `json:"amount"`
}

// BidResponse represents bid placement response
type BidResponse struct {
	Success bool   `json:"success"`
	BidID   string `json:"bidId"`
	Message string `json:"message"`
}

// AuctionDetails represents detailed auction information
type AuctionDetails struct {
	Success   bool                   `json:"success"`
	Auction   map[string]interface{} `json:"auction"`
	Bids      []map[string]interface{} `json:"bids"`
	TotalBids int                    `json:"totalBids"`
}

// CreateAuction creates a new auction via Firebase
func (c *Client) CreateAuction(ctx context.Context, auctionData AuctionData) (*AuctionResponse, error) {
	var result AuctionResponse
	err := c.callFirebase(ctx, "createAuction", auctionData, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to create auction: %w", err)
	}
	return &result, nil
}

// PlaceBid places a bid on an auction via Firebase
func (c *Client) PlaceBid(ctx context.Context, bidData BidData) (*BidResponse, error) {
	var result BidResponse
	err := c.callFirebase(ctx, "placeBid", bidData, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to place bid: %w", err)
	}
	return &result, nil
}

// EndAuction ends an auction via Firebase
func (c *Client) EndAuction(ctx context.Context, auctionID string) (*EndAuctionResponse, error) {
	data := map[string]string{"auctionId": auctionID}
	var result EndAuctionResponse
	err := c.callFirebase(ctx, "endAuction", data, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to end auction: %w", err)
	}
	return &result, nil
}

// GetAuctionDetails gets detailed auction information
func (c *Client) GetAuctionDetails(ctx context.Context, auctionID string) (*AuctionDetails, error) {
	data := map[string]string{"auctionId": auctionID}
	var result AuctionDetails
	err := c.callFirebase(ctx, "getAuctionDetails", data, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to get auction details: %w", err)
	}
	return &result, nil
}

// EndAuctionResponse represents auction completion response
type EndAuctionResponse struct {
	Success    bool    `json:"success"`
	WinnerID   *string `json:"winnerId"`
	WinningBid float64 `json:"winningBid"`
	Message    string  `json:"message"`
}