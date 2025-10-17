package firebase

import (
	"context"
	"fmt"
)

// NotificationData represents notification request
type NotificationData struct {
	UserID string            `json:"userId"`
	Title  string            `json:"title"`
	Body   string            `json:"body"`
	Data   map[string]string `json:"data,omitempty"`
}

// NotificationResponse represents notification response
type NotificationResponse struct {
	Success   bool   `json:"success"`
	MessageID string `json:"messageId"`
	Message   string `json:"message"`
}

// AuctionUpdateData represents auction update notification
type AuctionUpdateData struct {
	AuctionID string `json:"auctionId"`
	Type      string `json:"type"`
	Message   string `json:"message"`
}

// AuctionUpdateResponse represents auction update response
type AuctionUpdateResponse struct {
	Success            int                      `json:"success"`
	NotificationsSent  int                      `json:"notificationsSent"`
	Notifications      []map[string]interface{} `json:"notifications"`
}

// SendNotification sends a push notification to a user
func (c *Client) SendNotification(ctx context.Context, userID, title, body string, data map[string]string) (*NotificationResponse, error) {
	reqData := NotificationData{
		UserID: userID,
		Title:  title,
		Body:   body,
		Data:   data,
	}

	var result NotificationResponse
	err := c.callFirebase(ctx, "sendNotification", reqData, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to send notification: %w", err)
	}
	return &result, nil
}

// SendAuctionUpdate sends auction update notifications to participants
func (c *Client) SendAuctionUpdate(ctx context.Context, auctionID, updateType, message string) (*AuctionUpdateResponse, error) {
	data := AuctionUpdateData{
		AuctionID: auctionID,
		Type:      updateType,
		Message:   message,
	}

	var result AuctionUpdateResponse
	err := c.callFirebase(ctx, "sendAuctionUpdate", data, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to send auction update: %w", err)
	}
	return &result, nil
}