package models

import (
	"time"
	"gorm.io/gorm"
)

type Payment struct {
	ID              string         `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID          string         `json:"user_id" gorm:"not null;index"`
	OrderID         string         `json:"order_id" gorm:"not null;index"`
	Amount          int64          `json:"amount" gorm:"not null"` // Amount in cents
	Currency        string         `json:"currency" gorm:"not null;default:'USD'"`
	Status          string         `json:"status" gorm:"not null;default:'pending'"`
	PaymentMethod   string         `json:"payment_method" gorm:"not null"`
	Provider        string         `json:"provider" gorm:"not null"`
	ProviderID      string         `json:"provider_id" gorm:"index"`
	FailureReason   string         `json:"failure_reason,omitempty"`
	RefundedAmount  int64          `json:"refunded_amount" gorm:"default:0"`
	RefundedAt      *time.Time     `json:"refunded_at,omitempty"`
	Metadata        string         `json:"metadata,omitempty" gorm:"type:text"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusProcessing PaymentStatus = "processing"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
	PaymentStatusCancelled PaymentStatus = "cancelled"
)

type PaymentMethod string

const (
	PaymentMethodCard        PaymentMethod = "card"
	PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
	PaymentMethodPayPal      PaymentMethod = "paypal"
	PaymentMethodApplePay    PaymentMethod = "apple_pay"
	PaymentMethodGooglePay   PaymentMethod = "google_pay"
)

type ProcessPaymentRequest struct {
	OrderID       string `json:"order_id" binding:"required"`
	Amount        int64  `json:"amount" binding:"required,min=1"`
	Currency      string `json:"currency" binding:"required,len=3"`
	PaymentMethod string `json:"payment_method" binding:"required"`
	Provider      string `json:"provider" binding:"required"`
	Token         string `json:"token" binding:"required"`
}

type RefundRequest struct {
	Amount int64  `json:"amount" binding:"required,min=1"`
	Reason string `json:"reason" binding:"required"`
}

type PaymentResponse struct {
	ID             string  `json:"id"`
	UserID         string  `json:"user_id"`
	OrderID        string  `json:"order_id"`
	Amount         int64   `json:"amount"`
	Currency       string  `json:"currency"`
	Status         string  `json:"status"`
	PaymentMethod  string  `json:"payment_method"`
	Provider       string  `json:"provider"`
	ProviderID     string  `json:"provider_id,omitempty"`
	FailureReason  string  `json:"failure_reason,omitempty"`
	RefundedAmount int64   `json:"refunded_amount"`
	RefundedAt     *string `json:"refunded_at,omitempty"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

type PaymentHistoryResponse struct {
	Payments []PaymentResponse `json:"payments"`
	Total    int64             `json:"total"`
}

type PaymentMethodsResponse struct {
	Methods []PaymentMethodInfo `json:"methods"`
}

type PaymentMethodInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}