package models

import (
	"time"

	"gorm.io/gorm"
)

// Payment represents a payment transaction
type Payment struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// Payment details
	PaymentID     string         `gorm:"uniqueIndex;not null" json:"payment_id"`
	OrderID       string         `gorm:"not null;index" json:"order_id"`
	UserID        string         `gorm:"not null;index" json:"user_id"`
	AuctionID     string         `gorm:"not null;index" json:"auction_id"`

	// Amount details
	Amount        int64          `gorm:"not null" json:"amount"`        // in cents
	Currency      string         `gorm:"not null" json:"currency"`
	Fee           int64          `json:"fee"`                           // processing fee in cents
	NetAmount     int64          `json:"net_amount"`                    // net amount after fees

	// Payment status
	Status        string         `gorm:"not null;index" json:"status"`
	PaymentMethod string         `json:"payment_method"`
	Provider      string         `gorm:"not null" json:"provider"`      // stripe, etc.

	// Provider details
	ProviderID    string         `json:"provider_id"`
	ProviderData  string         `gorm:"type:text" json:"provider_data"`

	// Refund information
	RefundStatus  string         `json:"refund_status"`
	RefundAmount  int64          `json:"refund_amount"`
	RefundedAt    *time.Time     `json:"refunded_at"`

	// Metadata
	Description   string         `json:"description"`
	Metadata      string         `gorm:"type:text" json:"metadata"`
}

// PaymentStatus constants
const (
	PaymentStatusPending    = "pending"
	PaymentStatusProcessing = "processing"
	PaymentStatusSucceeded  = "succeeded"
	PaymentStatusFailed     = "failed"
	PaymentStatusCanceled   = "canceled"
	PaymentStatusRefunded   = "refunded"
	PaymentStatusPartiallyRefunded = "partially_refunded"
)

// RefundStatus constants
const (
	RefundStatusNotRefunded = "not_refunded"
	RefundStatusPending     = "refund_pending"
	RefundStatusSucceeded   = "refund_succeeded"
	RefundStatusFailed      = "refund_failed"
)

// CreatePaymentRequest represents a payment creation request
type CreatePaymentRequest struct {
	OrderID       string `json:"order_id" binding:"required"`
	UserID        string `json:"user_id" binding:"required"`
	AuctionID     string `json:"auction_id" binding:"required"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,len=3"`
	PaymentMethod string `json:"payment_method" binding:"required"`
	Description   string `json:"description"`
	ReturnURL     string `json:"return_url"`
	CancelURL     string `json:"cancel_url"`
}

// PaymentResponse represents a payment response
type PaymentResponse struct {
	PaymentID     string `json:"payment_id"`
	Status        string `json:"status"`
	Amount        int64  `json:"amount"`
	Currency      string `json:"currency"`
	PaymentMethod string `json:"payment_method"`
	ClientSecret  string `json:"client_secret,omitempty"`
	NextAction    string `json:"next_action,omitempty"`
	RedirectURL   string `json:"redirect_url,omitempty"`
}

// RefundRequest represents a refund request
type RefundRequest struct {
	PaymentID     string `json:"payment_id" binding:"required"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Reason        string `json:"reason"`
}

// WebhookEvent represents a payment provider webhook event
type WebhookEvent struct {
	ID            string                 `json:"id"`
	Type          string                 `json:"type"`
	PaymentID     string                 `json:"payment_id"`
	Status        string                 `json:"status"`
	Amount        int64                  `json:"amount"`
	Currency      string                 `json:"currency"`
	ProviderData  map[string]interface{} `json:"provider_data"`
}