package services

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/gmsas95/blytz-mvp/services/product-service/internal/config"
	"github.com/gmsas95/blytz-mvp/services/product-service/internal/models"
	shared_errors "github.com/gmsas95/blytz-mvp/shared/pkg/errors"
)

type ProductService struct {
	db     *gorm.DB
	logger *zap.Logger
	config *config.Config
}

func NewProductService(db *gorm.DB, logger *zap.Logger, config *config.Config) *ProductService {
	return &ProductService{
		db:     db,
		logger: logger,
		config: config,
	}
}

// CreateProduct creates a new product
func (s *ProductService) CreateProduct(ctx context.Context, userID string, req *models.CreateProductRequest) (*models.Product, error) {
	product := &models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Currency:    req.Currency,
		ImageURL:    req.ImageURL,
		SellerID:    userID,
		Stock:       req.Stock,
		Category:    req.Category,
		Subcategory: req.Subcategory,
		Status:      models.ProductStatusActive,
		IsActive:    true,
	}

	// Set images array
	product.SetImagesArray(req.Images)

	// Set tags array
	product.SetTagsArray(req.Tags)

	if err := s.db.WithContext(ctx).Create(product).Error; err != nil {
		s.logger.Error("Failed to create product", zap.Error(err))
		return nil, shared_errors.ErrInternalServer
	}

	return product, nil
}

// GetProduct retrieves a product by ID
func (s *ProductService) GetProduct(ctx context.Context, productID string) (*models.Product, error) {
	var product models.Product
	if err := s.db.WithContext(ctx).Where("product_id = ?", productID).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared_errors.ErrNotFound
		}
		s.logger.Error("Failed to get product", zap.String("product_id", productID), zap.Error(err))
		return nil, shared_errors.ErrInternalServer
	}

	return &product, nil
}

// GetProducts retrieves a list of products with filtering and pagination
func (s *ProductService) GetProducts(ctx context.Context, filter *models.ProductFilter, page, pageSize int) (*models.ProductListResponse, error) {
	var products []models.Product
	var total int64

	query := s.db.WithContext(ctx).Model(&models.Product{})

	// Apply filters
	if filter != nil {
		if filter.Category != "" {
			query = query.Where("category = ?", filter.Category)
		}
		if filter.Subcategory != "" {
			query = query.Where("subcategory = ?", filter.Subcategory)
		}
		if filter.MinPrice > 0 {
			query = query.Where("price >= ?", filter.MinPrice)
		}
		if filter.MaxPrice > 0 {
			query = query.Where("price <= ?", filter.MaxPrice)
		}
		if filter.SellerID != "" {
			query = query.Where("seller_id = ?", filter.SellerID)
		}
		if filter.Status != "" {
			query = query.Where("status = ?", filter.Status)
		}
		if filter.IsFeatured != nil {
			query = query.Where("is_featured = ?", *filter.IsFeatured)
		}
		if filter.Search != "" {
			searchPattern := "%" + filter.Search + "%"
			query = query.Where("name ILIKE ? OR description ILIKE ?", searchPattern, searchPattern)
		}
		if len(filter.Tags) > 0 {
			for _, tag := range filter.Tags {
				query = query.Where("tags ILIKE ?", "%"+tag+"%")
			}
		}
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		s.logger.Error("Failed to count products", zap.Error(err))
		return nil, shared_errors.ErrInternalServer
	}

	// Apply pagination
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&products).Error; err != nil {
		s.logger.Error("Failed to get products", zap.Error(err))
		return nil, shared_errors.ErrInternalServer
	}

	// Calculate available stock for each product
	for i := range products {
		products[i].Available = products[i].GetAvailable()
	}

	totalPages := int(total)/pageSize + 1
	if int(total)%pageSize == 0 {
		totalPages = int(total) / pageSize
	}

	return &models.ProductListResponse{
		Products: products,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		HasNext:  page < totalPages,
	}, nil
}

// GetFeaturedProducts retrieves featured products
func (s *ProductService) GetFeaturedProducts(ctx context.Context, limit int) ([]models.Product, error) {
	var products []models.Product

	if limit <= 0 {
		limit = 10
	}

	if err := s.db.WithContext(ctx).
		Where("is_featured = ? AND is_active = ? AND status = ?", true, true, models.ProductStatusActive).
		Limit(limit).
		Order("created_at DESC").
		Find(&products).Error; err != nil {
		s.logger.Error("Failed to get featured products", zap.Error(err))
		return nil, shared_errors.ErrInternalServer
	}

	// Calculate available stock for each product
	for i := range products {
		products[i].Available = products[i].GetAvailable()
	}

	return products, nil
}

// UpdateProduct updates an existing product
func (s *ProductService) UpdateProduct(ctx context.Context, userID string, productID string, req *models.UpdateProductRequest) (*models.Product, error) {
	var product models.Product

	// First, get the product and verify ownership
	if err := s.db.WithContext(ctx).Where("product_id = ? AND seller_id = ?", productID, userID).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared_errors.ErrNotFound
		}
		s.logger.Error("Failed to get product for update", zap.String("product_id", productID), zap.Error(err))
		return nil, shared_errors.ErrInternalServer
	}

	// Update fields if provided
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price > 0 {
		product.Price = req.Price
	}
	if req.Currency != "" {
		product.Currency = req.Currency
	}
	if req.ImageURL != "" {
		product.ImageURL = req.ImageURL
	}
	if req.Images != nil {
		product.SetImagesArray(req.Images)
	}
	if req.Stock >= 0 {
		product.Stock = req.Stock
	}
	if req.Category != "" {
		product.Category = req.Category
	}
	if req.Subcategory != "" {
		product.Subcategory = req.Subcategory
	}
	if req.Tags != nil {
		product.SetTagsArray(req.Tags)
	}
	if req.Status != "" {
		product.Status = req.Status
	}
	product.IsFeatured = req.IsFeatured

	product.UpdatedAt = time.Now()

	if err := s.db.WithContext(ctx).Save(&product).Error; err != nil {
		s.logger.Error("Failed to update product", zap.String("product_id", productID), zap.Error(err))
		return nil, shared_errors.ErrInternalServer
	}

	return &product, nil
}

// DeleteProduct soft deletes a product
func (s *ProductService) DeleteProduct(ctx context.Context, userID string, productID string) error {
	result := s.db.WithContext(ctx).Where("product_id = ? AND seller_id = ?", productID, userID).Delete(&models.Product{})
	if result.Error != nil {
		s.logger.Error("Failed to delete product", zap.String("product_id", productID), zap.Error(result.Error))
		return shared_errors.ErrInternalServer
	}
	if result.RowsAffected == 0 {
		return shared_errors.ErrNotFound
	}

	return nil
}

// UpdateInventory updates product inventory
func (s *ProductService) UpdateInventory(ctx context.Context, productID string, req *models.InventoryUpdateRequest) error {
	result := s.db.WithContext(ctx).Model(&models.Product{}).
		Where("product_id = ?", productID).
		Updates(map[string]interface{}{
			"stock":      req.Stock,
			"reserved":   req.Reserved,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		s.logger.Error("Failed to update inventory", zap.String("product_id", productID), zap.Error(result.Error))
		return shared_errors.ErrInternalServer
	}
	if result.RowsAffected == 0 {
		return shared_errors.ErrNotFound
	}

	return nil
}

// GetProductsBySeller retrieves products for a specific seller
func (s *ProductService) GetProductsBySeller(ctx context.Context, sellerID string, page, pageSize int) (*models.ProductListResponse, error) {
	filter := &models.ProductFilter{
		SellerID: sellerID,
		Status:   models.ProductStatusActive,
	}

	return s.GetProducts(ctx, filter, page, pageSize)
}

// ReserveStock reserves stock for a product (used by order service)
func (s *ProductService) ReserveStock(ctx context.Context, productID string, quantity int) error {
	result := s.db.WithContext(ctx).Model(&models.Product{}).
		Where("product_id = ? AND stock - reserved >= ?", productID, quantity).
		Update("reserved", gorm.Expr("reserved + ?", quantity))

	if result.Error != nil {
		s.logger.Error("Failed to reserve stock", zap.String("product_id", productID), zap.Error(result.Error))
		return shared_errors.ErrInternalServer
	}
	if result.RowsAffected == 0 {
		return shared_errors.ErrInsufficientStock
	}

	return nil
}

// ReleaseStock releases reserved stock for a product
func (s *ProductService) ReleaseStock(ctx context.Context, productID string, quantity int) error {
	result := s.db.WithContext(ctx).Model(&models.Product{}).
		Where("product_id = ?", productID).
		Update("reserved", gorm.Expr("GREATEST(reserved - ?, 0)", quantity))

	if result.Error != nil {
		s.logger.Error("Failed to release stock", zap.String("product_id", productID), zap.Error(result.Error))
		return shared_errors.ErrInternalServer
	}
	if result.RowsAffected == 0 {
		return shared_errors.ErrNotFound
	}

	return nil
}

// ConfirmStockDeduction confirms stock deduction after order completion
func (s *ProductService) ConfirmStockDeduction(ctx context.Context, productID string, quantity int) error {
	result := s.db.WithContext(ctx).Model(&models.Product{}).
		Where("product_id = ?", productID).
		Updates(map[string]interface{}{
			"stock":      gorm.Expr("stock - ?", quantity),
			"reserved":   gorm.Expr("GREATEST(reserved - ?, 0)", quantity),
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		s.logger.Error("Failed to confirm stock deduction", zap.String("product_id", productID), zap.Error(result.Error))
		return shared_errors.ErrInternalServer
	}
	if result.RowsAffected == 0 {
		return shared_errors.ErrNotFound
	}

	return nil
}
