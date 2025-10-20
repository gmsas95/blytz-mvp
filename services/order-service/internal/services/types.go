package services

import (
	"fmt"
	"strings"
)

type CreateOrderRequest struct {
	AuctionID       *string        `json:"auction_id,omitempty"`
	ProductID       string         `json:"product_id" binding:"required"`
	ProductName     string         `json:"product_name" binding:"required"`
	ProductImage    string         `json:"product_image,omitempty"`
	Quantity        int            `json:"quantity" binding:"required,min=1"`
	Price           int64          `json:"price" binding:"required,min=1"` // Price in cents
	Currency        string         `json:"currency" binding:"required,len=3"`
	ShippingAddress AddressRequest `json:"shipping_address" binding:"required"`
	BillingAddress  AddressRequest `json:"billing_address" binding:"required"`
	Notes           string         `json:"notes,omitempty"`
}

type AddressRequest struct {
	Name        string `json:"name" binding:"required"`
	Street      string `json:"street" binding:"required"`
	City        string `json:"city" binding:"required"`
	State       string `json:"state" binding:"required"`
	PostalCode  string `json:"postal_code" binding:"required"`
	Country     string `json:"country" binding:"required,len=2"`
	PhoneNumber string `json:"phone_number,omitempty"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending processing confirmed shipped delivered cancelled refunded"`
}

type OrderResponse struct {
	ID              string         `json:"id"`
	UserID          string         `json:"user_id"`
	AuctionID       *string        `json:"auction_id,omitempty"`
	ProductID       string         `json:"product_id"`
	ProductName     string         `json:"product_name"`
	ProductImage    string         `json:"product_image,omitempty"`
	Quantity        int            `json:"quantity"`
	Price           int64          `json:"price"` // Price in cents
	TotalAmount     int64          `json:"total_amount"` // Total in cents
	Currency        string         `json:"currency"`
	Status          string         `json:"status"`
	PaymentStatus   string         `json:"payment_status"`
	PaymentMethod   string         `json:"payment_method,omitempty"`
	ShippingAddress AddressResponse `json:"shipping_address"`
	BillingAddress  AddressResponse `json:"billing_address"`
	Notes           string         `json:"notes,omitempty"`
	CreatedAt       string         `json:"created_at"`
	UpdatedAt       string         `json:"updated_at"`
}

type AddressResponse struct {
	Name        string `json:"name"`
	Street      string `json:"street"`
	City        string `json:"city"`
	State       string `json:"state"`
	PostalCode  string `json:"postal_code"`
	Country     string `json:"country"`
	PhoneNumber string `json:"phone_number,omitempty"`
}

type OrdersListResponse struct {
	Orders     []OrderResponse `json:"orders"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
	TotalPages int             `json:"total_pages"`
}

func (r *CreateOrderRequest) Validate() error {
	if r.Quantity <= 0 {
		return fmt.Errorf("quantity must be greater than 0")
	}
	if r.Price <= 0 {
		return fmt.Errorf("price must be greater than 0")
	}
	if len(r.Currency) != 3 {
		return fmt.Errorf("currency must be 3 characters")
	}
	if strings.TrimSpace(r.ProductID) == "" {
		return fmt.Errorf("product_id is required")
	}
	if strings.TrimSpace(r.ProductName) == "" {
		return fmt.Errorf("product_name is required")
	}
	return nil
}