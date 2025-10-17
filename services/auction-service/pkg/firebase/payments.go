package firebase

import (
	"context"
	"fmt"
)

// PaymentIntentData represents payment intent creation request
type PaymentIntentData struct {
	Amount    float64 `json:"amount"`
	Currency  string  `json:"currency,omitempty"`
	AuctionID string  `json:"auctionId"`
	BidID     string  `json:"bidId"`
}

// PaymentIntentResponse represents payment intent creation response
type PaymentIntentResponse struct {
	Success         bool   `json:"success"`
	ClientSecret    string `json:"clientSecret"`
	PaymentIntentID string `json:"paymentIntentId"`
}

// ConfirmPaymentData represents payment confirmation request
type ConfirmPaymentData struct {
	PaymentIntentID string `json:"paymentIntentId"`
	AuctionID       string `json:"auctionId"`
	BidID           string `json:"bidId"`
}

// ConfirmPaymentResponse represents payment confirmation response
type ConfirmPaymentResponse struct {
	Success bool   `json:"success"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

// CreatePaymentIntent creates a payment intent for bidding
func (c *Client) CreatePaymentIntent(ctx context.Context, amount float64, auctionID, bidID string) (*PaymentIntentResponse, error) {
	data := PaymentIntentData{
		Amount:    amount,
		Currency:  "usd",
		AuctionID: auctionID,
		BidID:     bidID,
	}

	var result PaymentIntentResponse
	err := c.callFirebase(ctx, "createPaymentIntent", data, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment intent: %w", err)
	}
	return &result, nil
}

// ConfirmPayment confirms a payment and updates bid status
func (c *Client) ConfirmPayment(ctx context.Context, paymentIntentID, auctionID, bidID string) (*ConfirmPaymentResponse, error) {
	data := ConfirmPaymentData{
		PaymentIntentID: paymentIntentID,
		AuctionID:       auctionID,
		BidID:           bidID,
	}

	var result ConfirmPaymentResponse
	err := c.callFirebase(ctx, "confirmPayment", data, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm payment: %w", err)
	}
	return &result, nil
}