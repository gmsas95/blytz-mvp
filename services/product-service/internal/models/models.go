package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Product represents a product in the system
type Product struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Basic info
	ProductID   string `gorm:"uniqueIndex;not null" json:"product_id"`
	Name        string `gorm:"not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	Price       int64  `gorm:"not null" json:"price"` // Price in cents
	Currency    string `gorm:"not null;default:'USD'" json:"currency"`

	// Images and media
	ImageURL string `json:"image_url"`
	Images   string `gorm:"type:text" json:"images"` // JSON array of image URLs

	// Seller info
	SellerID   string `gorm:"not null;index" json:"seller_id"`
	SellerName string `json:"seller_name"`

	// Inventory
	Stock     int `gorm:"default:0" json:"stock"`
	Reserved  int `gorm:"default:0" json:"reserved"`
	Available int `gorm:"-" json:"available"` // Calculated field

	// Status
	Status     string `gorm:"default:'active'" json:"status"`
	IsFeatured bool   `gorm:"default:false" json:"is_featured"`
	IsActive   bool   `gorm:"default:true" json:"is_active"`

	// Categories
	Category    string `json:"category"`
	Subcategory string `json:"subcategory"`
	Tags        string `gorm:"type:text" json:"tags"` // JSON array

	// Metadata
	Metadata string `gorm:"type:text" json:"metadata,omitempty"`
}

// BeforeCreate hook to generate ProductID
func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.ProductID == "" {
		p.ProductID = uuid.New().String()
	}
	return nil
}

// GetAvailable calculates available stock
func (p *Product) GetAvailable() int {
	return p.Stock - p.Reserved
}

// GetImagesArray returns images as string array
func (p *Product) GetImagesArray() []string {
	if p.Images == "" {
		return []string{}
	}
	var images []string
	json.Unmarshal([]byte(p.Images), &images)
	return images
}

// SetImagesArray sets images from string array
func (p *Product) SetImagesArray(images []string) {
	if len(images) == 0 {
		p.Images = ""
		return
	}
	data, _ := json.Marshal(images)
	p.Images = string(data)
}

// GetTagsArray returns tags as string array
func (p *Product) GetTagsArray() []string {
	if p.Tags == "" {
		return []string{}
	}
	var tags []string
	json.Unmarshal([]byte(p.Tags), &tags)
	return tags
}

// SetTagsArray sets tags from string array
func (p *Product) SetTagsArray(tags []string) {
	if len(tags) == 0 {
		p.Tags = ""
		return
	}
	data, _ := json.Marshal(tags)
	p.Tags = string(data)
}

// ProductStatus constants
const (
	ProductStatusActive   = "active"
	ProductStatusDraft    = "draft"
	ProductStatusArchived = "archived"
	ProductStatusSoldOut  = "sold_out"
)

// CreateProductRequest represents a product creation request
type CreateProductRequest struct {
	Name        string   `json:"name" binding:"required,min=3,max=200"`
	Description string   `json:"description" binding:"required,min=10,max=2000"`
	Price       int64    `json:"price" binding:"required,gt=0"`
	Currency    string   `json:"currency" binding:"required,len=3"`
	ImageURL    string   `json:"image_url" binding:"required,url"`
	Images      []string `json:"images"`
	Stock       int      `json:"stock" binding:"required,min=0"`
	Category    string   `json:"category" binding:"required"`
	Subcategory string   `json:"subcategory"`
	Tags        []string `json:"tags"`
}

// UpdateProductRequest represents a product update request
type UpdateProductRequest struct {
	Name        string   `json:"name" binding:"omitempty,min=3,max=200"`
	Description string   `json:"description" binding:"omitempty,min=10,max=2000"`
	Price       int64    `json:"price" binding:"omitempty,gt=0"`
	Currency    string   `json:"currency" binding:"omitempty,len=3"`
	ImageURL    string   `json:"image_url" binding:"omitempty,url"`
	Images      []string `json:"images"`
	Stock       int      `json:"stock" binding:"omitempty,min=0"`
	Category    string   `json:"category" binding:"omitempty"`
	Subcategory string   `json:"subcategory"`
	Tags        []string `json:"tags"`
	Status      string   `json:"status" binding:"omitempty,oneof=active draft archived sold_out"`
	IsFeatured  bool     `json:"is_featured"`
}

// InventoryUpdateRequest represents an inventory update request
type InventoryUpdateRequest struct {
	Stock    int `json:"stock" binding:"required,min=0"`
	Reserved int `json:"reserved" binding:"min=0"`
}

// ProductResponse represents a product response
type ProductResponse struct {
	Product   Product `json:"product"`
	Available int     `json:"available"`
}

// ProductListResponse represents a paginated product list response
type ProductListResponse struct {
	Products []Product `json:"products"`
	Total    int64     `json:"total"`
	Page     int       `json:"page"`
	PageSize int       `json:"page_size"`
	HasNext  bool      `json:"has_next"`
}

// ProductFilter represents product filtering options
type ProductFilter struct {
	Category    string
	Subcategory string
	MinPrice    int64
	MaxPrice    int64
	SellerID    string
	Status      string
	IsFeatured  *bool
	Search      string
	Tags        []string
}
