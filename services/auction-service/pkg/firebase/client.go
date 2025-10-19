package firebase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Client handles communication with Firebase Functions
type Client struct {
	baseURL    string
	httpClient *http.Client
	logger     *zap.Logger
}

// NewClient creates a new Firebase client
func NewClient(logger *zap.Logger) *Client {
	return &Client{
		baseURL: "http://localhost:5001/demo-blytz-mvp/us-central1",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

// callFirebase makes a POST request to Firebase functions
func (c *Client) callFirebase(ctx context.Context, function string, data interface{}, result interface{}) error {
	payload := map[string]interface{}{
		"data": data,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request data: %w", err)
	}

	url := fmt.Sprintf("%s/%s", c.baseURL, function)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	c.logger.Info("Calling Firebase function",
		zap.String("function", function),
		zap.String("url", url))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call Firebase function: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResp struct {
			Error struct {
				Message string `json:"message"`
				Status  string `json:"status"`
			} `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err == nil {
			return fmt.Errorf("Firebase function error: %s (status: %s)",
				errorResp.Error.Message, errorResp.Error.Status)
		}
		return fmt.Errorf("Firebase function returned status %d", resp.StatusCode)
	}

	var response struct {
		Result interface{} `json:"result"`
	}
	response.Result = result

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	c.logger.Info("Firebase function call successful",
		zap.String("function", function))

	return nil
}

// UpdateAuction updates an auction via Firebase
func (c *Client) UpdateAuction(ctx context.Context, auctionID string, updateData interface{}) error {
	// For now, we'll use the existing EndAuction method as a placeholder
	// In a real implementation, you'd have a separate updateAuction function
	data := map[string]interface{}{
		"auctionId": auctionID,
		"update":    updateData,
	}
	return c.callFirebase(ctx, "updateAuction", data, nil)
}

// GetAuction gets auction details via Firebase
func (c *Client) GetAuction(ctx context.Context, auctionID string) (map[string]interface{}, error) {
	result, err := c.GetAuctionDetails(ctx, auctionID)
	if err != nil {
		return nil, err
	}
	return result.Auction, nil
}

// ProcessPayment processes a payment via Firebase
func (c *Client) ProcessPayment(ctx context.Context, paymentData interface{}) error {
	return c.callFirebase(ctx, "processPayment", paymentData, nil)
}