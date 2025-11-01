package models

import (
	"gorm.io/gorm"
	"time"
)

type Payment struct {
	ID             string         `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID         string         `json:"user_id" gorm:"not null;index"`
	OrderID        string         `json:"order_id" gorm:"not null;index"`
	Amount         int64          `json:"amount" gorm:"not null"` // Amount in cents
	Currency       string         `json:"currency" gorm:"not null;default:'MYR'"`
	Status         string         `json:"status" gorm:"not null;default:'pending'"`
	PaymentMethod  string         `json:"payment_method" gorm:"not null"`
	Provider       string         `json:"provider" gorm:"not null"`
	ProviderID     string         `json:"provider_id" gorm:"index"`
	Channel        string         `json:"channel,omitempty"`        // Fiuu payment channel
	TransactionID  string         `json:"transaction_id,omitempty"` // Fiuu transaction ID
	PaymentRefID   string         `json:"payment_ref_id,omitempty"` // Fiuu payment reference ID
	FailureReason  string         `json:"failure_reason,omitempty"`
	RefundedAmount int64          `json:"refunded_amount" gorm:"default:0"`
	RefundedAt     *time.Time     `json:"refunded_at,omitempty"`
	Metadata       string         `json:"metadata,omitempty" gorm:"type:text"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

type PaymentStatus string

const (
	PaymentStatusPending    PaymentStatus = "pending"
	PaymentStatusProcessing PaymentStatus = "processing"
	PaymentStatusCompleted  PaymentStatus = "completed"
	PaymentStatusFailed     PaymentStatus = "failed"
	PaymentStatusRefunded   PaymentStatus = "refunded"
	PaymentStatusCancelled  PaymentStatus = "cancelled"
)

type PaymentMethod string

const (
	PaymentMethodCard         PaymentMethod = "card"
	PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
	PaymentMethodPayPal       PaymentMethod = "paypal"
	PaymentMethodApplePay     PaymentMethod = "apple_pay"
	PaymentMethodGooglePay    PaymentMethod = "google_pay"
	PaymentMethodFPX          PaymentMethod = "fpx"
	PaymentMethodEwallet      PaymentMethod = "ewallet"
	PaymentMethodBNPL         PaymentMethod = "bnpl"
	PaymentMethodQR           PaymentMethod = "qr"
)

type ProcessPaymentRequest struct {
	OrderID       string `json:"order_id" binding:"required"`
	Amount        int64  `json:"amount" binding:"required,min=1"`
	Currency      string `json:"currency" binding:"required,len=3"`
	PaymentMethod string `json:"payment_method" binding:"required"`
	Provider      string `json:"provider" binding:"required"`
	Channel       string `json:"channel,omitempty"` // Fiuu payment channel
	Token         string `json:"token" binding:"required"`
	BillName      string `json:"bill_name,omitempty"`
	BillEmail     string `json:"bill_email,omitempty"`
	BillMobile    string `json:"bill_mobile,omitempty"`
	BillDesc      string `json:"bill_desc,omitempty"`
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
	Channel     string `json:"channel,omitempty"`  // Fiuu channel code
	Currency    string `json:"currency,omitempty"` // Supported currency
}

// FiuuWebhookRequest represents Fiuu webhook payload
type FiuuWebhookRequest struct {
	TransactionID    string `json:"tranID"`
	OrderID          string `json:"order_id"`
	Amount           string `json:"amount"`
	Currency         string `json:"currency"`
	PaymentStatus    string `json:"payment_status"`
	PaymentChannel   string `json:"payment_channel"`
	ChannelCode      string `json:"channel_code"`
	PaymentRefID     string `json:"payment_ref_id"`
	PayDate          string `json:"pay_date"`
	PayTime          string `json:"pay_time"`
	ErrorCode        string `json:"error_code"`
	ErrorDescription string `json:"error_desc"`
	Signature        string `json:"signature"`
}

// FiuuSeamlessConfig represents configuration for frontend seamless integration
type FiuuSeamlessConfig struct {
	MerchantID string                 `json:"mpsmerchantid"`
	Channel    string                 `json:"mpschannel"`
	Amount     string                 `json:"mpsamount"`
	OrderID    string                 `json:"mpsorderid"`
	BillName   string                 `json:"mpsbill_name"`
	BillEmail  string                 `json:"mpsbill_email"`
	BillMobile string                 `json:"mpsbill_mobile"`
	BillDesc   string                 `json:"mpsbill_desc"`
	Currency   string                 `json:"mpscurrency"`
	LangCode   string                 `json:"mpslangcode"`
	VCode      string                 `json:"vcode"`
	Sandbox    bool                   `json:"sandbox"`
	ScriptURL  string                 `json:"scriptUrl"`
	Additional map[string]interface{} `json:"additional,omitempty"`
}
