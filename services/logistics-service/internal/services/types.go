package services

import (
	"fmt"
)

type CreateShipmentRequest struct {
	OrderID           string          `json:"order_id" binding:"required"`
	Carrier           string          `json:"carrier" binding:"required"`
	Service           string          `json:"service" binding:"required"`
	OriginAddress     AddressRequest  `json:"origin_address" binding:"required"`
	DestinationAddress AddressRequest `json:"destination_address" binding:"required"`
	Weight            float64         `json:"weight" binding:"required,min=0"`
	Dimensions        DimensionsRequest `json:"dimensions" binding:"required"`
	Cost              float64         `json:"cost" binding:"required,min=0"`
	Notes             string          `json:"notes,omitempty"`
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

type DimensionsRequest struct {
	Length float64 `json:"length" binding:"required,min=0"`
	Width  float64 `json:"width" binding:"required,min=0"`
	Height float64 `json:"height" binding:"required,min=0"`
}

type UpdateShipmentStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending label_created picked_up in_transit out_for_delivery delivered failed returned"`
}

type ShipmentResponse struct {
	ID                  string            `json:"id"`
	OrderID             string            `json:"order_id"`
	UserID              string            `json:"user_id"`
	TrackingNumber      string            `json:"tracking_number"`
	Carrier             string            `json:"carrier"`
	Service             string            `json:"service"`
	Status              string            `json:"status"`
	OriginAddress       AddressResponse   `json:"origin_address"`
	DestinationAddress  AddressResponse   `json:"destination_address"`
	Weight              float64           `json:"weight"`
	Dimensions          DimensionsResponse `json:"dimensions"`
	EstimatedDelivery   *string           `json:"estimated_delivery,omitempty"`
	ActualDelivery      *string           `json:"actual_delivery,omitempty"`
	Cost                float64           `json:"cost"`
	Notes               string            `json:"notes,omitempty"`
	CreatedAt           string            `json:"created_at"`
	UpdatedAt           string            `json:"updated_at"`
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

type DimensionsResponse struct {
	Length float64 `json:"length"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

type TrackingResponse struct {
	Shipment ShipmentResponse   `json:"shipment"`
	Events   []TrackingEventResponse `json:"events"`
}

type TrackingEventResponse struct {
	ID          string `json:"id"`
	Status      string `json:"status"`
	Description string `json:"description"`
	Location    string `json:"location,omitempty"`
	Timestamp   string `json:"timestamp"`
}

func (r *CreateShipmentRequest) Validate() error {
	if r.Weight <= 0 {
		return fmt.Errorf("weight must be greater than 0")
	}
	if r.Cost < 0 {
		return fmt.Errorf("cost must be non-negative")
	}
	if r.Dimensions.Length <= 0 || r.Dimensions.Width <= 0 || r.Dimensions.Height <= 0 {
		return fmt.Errorf("dimensions must be greater than 0")
	}
	return nil
}