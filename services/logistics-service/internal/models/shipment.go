package models

import (
	"time"
	"gorm.io/gorm"
)

type Shipment struct {
	ID              string         `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	OrderID         string         `json:"order_id" gorm:"not null;index"`
	UserID          string         `json:"user_id" gorm:"not null;index"`
	TrackingNumber  string         `json:"tracking_number" gorm:"uniqueIndex"`
	Carrier         string         `json:"carrier" gorm:"not null"`
	Service         string         `json:"service" gorm:"not null"`
	Status          string         `json:"status" gorm:"not null;default:'pending'"`
	OriginAddress   Address        `json:"origin_address" gorm:"embedded;embeddedPrefix:origin_"`
	DestinationAddress Address     `json:"destination_address" gorm:"embedded;embeddedPrefix:destination_"`
	Weight          float64        `json:"weight" gorm:"not null"`
	Dimensions      Dimensions     `json:"dimensions" gorm:"embedded;embeddedPrefix:dimensions_"`
	EstimatedDelivery *time.Time   `json:"estimated_delivery,omitempty"`
	ActualDelivery  *time.Time     `json:"actual_delivery,omitempty"`
	Cost            float64        `json:"cost" gorm:"not null"`
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

type Dimensions struct {
	Length float64 `json:"length" gorm:"not null"`
	Width  float64 `json:"width" gorm:"not null"`
	Height float64 `json:"height" gorm:"not null"`
}

type ShipmentStatus string

const (
	ShipmentStatusPending   ShipmentStatus = "pending"
	ShipmentStatusLabelCreated ShipmentStatus = "label_created"
	ShipmentStatusPickedUp  ShipmentStatus = "picked_up"
	ShipmentStatusInTransit ShipmentStatus = "in_transit"
	ShipmentStatusOutForDelivery ShipmentStatus = "out_for_delivery"
	ShipmentStatusDelivered ShipmentStatus = "delivered"
	ShipmentStatusFailed    ShipmentStatus = "failed"
	ShipmentStatusReturned  ShipmentStatus = "returned"
)

type TrackingEvent struct {
	ID          string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ShipmentID  string    `json:"shipment_id" gorm:"not null;index"`
	Status      string    `json:"status" gorm:"not null"`
	Description string    `json:"description" gorm:"not null"`
	Location    string    `json:"location,omitempty"`
	Timestamp   time.Time `json:"timestamp" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at"`
}