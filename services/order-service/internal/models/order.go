package models

import (
	"gorm.io/gorm"
	"time"
)

type Order struct {
	ID              string         `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID          string         `json:"user_id" gorm:"not null;index"`
	AuctionID       *string        `json:"auction_id,omitempty" gorm:"index"`
	ProductID       string         `json:"product_id" gorm:"not null;index"`
	ProductName     string         `json:"product_name" gorm:"not null"`
	ProductImage    string         `json:"product_image,omitempty"`
	Quantity        int            `json:"quantity" gorm:"not null;default:1"`
	Price           int64          `json:"price" gorm:"not null"`        // Price in cents
	TotalAmount     int64          `json:"total_amount" gorm:"not null"` // Total in cents
	Currency        string         `json:"currency" gorm:"not null;default:'USD'"`
	Status          string         `json:"status" gorm:"not null;default:'pending'"`
	PaymentStatus   string         `json:"payment_status" gorm:"not null;default:'pending'"`
	PaymentMethod   string         `json:"payment_method,omitempty"`
	ShippingAddress Address        `json:"shipping_address" gorm:"embedded;embeddedPrefix:shipping_"`
	BillingAddress  Address        `json:"billing_address" gorm:"embedded;embeddedPrefix:billing_"`
	Notes           string         `json:"notes,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

type Address struct {
	Name        string `json:"name" gorm:"not null"`
	Street      string `json:"street" gorm:"not null"`
	City        string `json:"city" gorm:"not null"`
	State       string `json:"state" gorm:"not null"`
	PostalCode  string `json:"postal_code" gorm:"not null"`
	Country     string `json:"country" gorm:"not null;default:'US'"`
	PhoneNumber string `json:"phone_number,omitempty"`
}

type OrderItem struct {
	ID          string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	OrderID     string    `json:"order_id" gorm:"not null;index"`
	ProductID   string    `json:"product_id" gorm:"not null"`
	ProductName string    `json:"product_name" gorm:"not null"`
	Quantity    int       `json:"quantity" gorm:"not null"`
	Price       int64     `json:"price" gorm:"not null"`       // Price per unit in cents
	TotalPrice  int64     `json:"total_price" gorm:"not null"` // Total for this item in cents
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusConfirmed  OrderStatus = "confirmed"
	OrderStatusShipped    OrderStatus = "shipped"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
	OrderStatusRefunded   OrderStatus = "refunded"
)

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusPaid      PaymentStatus = "paid"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
	PaymentStatusCancelled PaymentStatus = "cancelled"
)

// Cart models
type Cart struct {
	ID        string     `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID    string     `json:"user_id" gorm:"not null;uniqueIndex"`
	Items     []CartItem `json:"items" gorm:"foreignKey:CartID"`
	Total     int64      `json:"total" gorm:"not null;default:0"` // Total in cents
	ItemCount int        `json:"item_count" gorm:"not null;default:0"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type CartItem struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CartID    string    `json:"cart_id" gorm:"not null;index"`
	ProductID string    `json:"product_id" gorm:"not null"`
	AuctionID *string   `json:"auction_id,omitempty" gorm:"index"`
	Quantity  int       `json:"quantity" gorm:"not null;default:1"`
	Price     int64     `json:"price" gorm:"not null"` // Price per unit in cents
	Total     int64     `json:"total" gorm:"not null"` // Total for this item in cents
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
