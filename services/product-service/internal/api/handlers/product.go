package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/gmsas95/blytz-mvp/services/product-service/internal/models"
	"github.com/gmsas95/blytz-mvp/services/product-service/internal/services"
	"github.com/gmsas95/blytz-mvp/shared/pkg/errors"
	"github.com/gmsas95/blytz-mvp/shared/pkg/utils"
)

type ProductHandler struct {
	productService *services.ProductService
	logger         *zap.Logger
}

func NewProductHandler(productService *services.ProductService, logger *zap.Logger) *ProductHandler {
	return &ProductHandler{
		productService: productService,
		logger:         logger,
	}
}

// CreateProduct handles product creation
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		utils.ErrorResponse(c, errors.ErrUnauthorized)
		return
	}

	var req models.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, errors.ErrInvalidRequestBody)
		return
	}

	product, err := h.productService.CreateProduct(c.Request.Context(), userID, &req)
	if err != nil {
		h.logger.Error("Failed to create product", zap.Error(err))
		utils.ErrorResponse(c, err)
		return
	}

	response := h.mapProductToResponse(product)
	utils.SuccessResponse(c, response)
}

// GetProduct handles getting a single product
func (h *ProductHandler) GetProduct(c *gin.Context) {
	productID := c.Param("id")
	if productID == "" {
		utils.ErrorResponse(c, errors.ErrInvalidRequest)
		return
	}

	product, err := h.productService.GetProduct(c.Request.Context(), productID)
	if err != nil {
		if err == errors.ErrNotFound {
			utils.ErrorResponse(c, errors.ErrNotFound)
			return
		}
		h.logger.Error("Failed to get product", zap.String("product_id", productID), zap.Error(err))
		utils.ErrorResponse(c, err)
		return
	}

	response := h.mapProductToResponse(product)
	utils.SuccessResponse(c, response)
}

// GetProducts handles getting a list of products with filtering
func (h *ProductHandler) GetProducts(c *gin.Context) {
	// Parse query parameters
	page := 1
	pageSize := 20

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if ps := c.Query("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}

	// Build filter
	filter := &models.ProductFilter{
		Category:    c.Query("category"),
		Subcategory: c.Query("subcategory"),
		SellerID:    c.Query("seller_id"),
		Status:      c.Query("status"),
		Search:      c.Query("search"),
	}

	// Parse price range
	if minPrice := c.Query("min_price"); minPrice != "" {
		if parsed, err := strconv.ParseInt(minPrice, 10, 64); err == nil && parsed > 0 {
			filter.MinPrice = parsed
		}
	}
	if maxPrice := c.Query("max_price"); maxPrice != "" {
		if parsed, err := strconv.ParseInt(maxPrice, 10, 64); err == nil && parsed > 0 {
			filter.MaxPrice = parsed
		}
	}

	// Parse featured flag
	if featured := c.Query("featured"); featured != "" {
		featuredBool := featured == "true"
		filter.IsFeatured = &featuredBool
	}

	// Parse tags
	if tags := c.Query("tags"); tags != "" {
		filter.Tags = []string{tags} // Simple implementation, could be enhanced for comma-separated
	}

	response, err := h.productService.GetProducts(c.Request.Context(), filter, page, pageSize)
	if err != nil {
		h.logger.Error("Failed to get products", zap.Error(err))
		utils.ErrorResponse(c, err)
		return
	}

	// Map products to responses
	productResponses := make([]models.ProductResponse, len(response.Products))
	for i, product := range response.Products {
		productResponses[i] = *h.mapProductToResponse(&response.Products[i])
	}

	utils.SuccessResponse(c, gin.H{
		"products":  productResponses,
		"total":     response.Total,
		"page":      response.Page,
		"page_size": response.PageSize,
		"has_next":  response.HasNext,
	})
}

// GetFeaturedProducts handles getting featured products
func (h *ProductHandler) GetFeaturedProducts(c *gin.Context) {
	limit := 10
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 50 {
			limit = parsed
		}
	}

	products, err := h.productService.GetFeaturedProducts(c.Request.Context(), limit)
	if err != nil {
		h.logger.Error("Failed to get featured products", zap.Error(err))
		utils.ErrorResponse(c, err)
		return
	}

	// Map products to responses
	productResponses := make([]models.ProductResponse, len(products))
	for i, product := range products {
		productResponses[i] = *h.mapProductToResponse(&products[i])
	}

	utils.SuccessResponse(c, gin.H{
		"products": productResponses,
		"total":    len(products),
	})
}

// UpdateProduct handles product updates
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		utils.ErrorResponse(c, errors.ErrUnauthorized)
		return
	}

	productID := c.Param("id")
	if productID == "" {
		utils.ErrorResponse(c, errors.ErrInvalidRequest)
		return
	}

	var req models.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, errors.ErrInvalidRequestBody)
		return
	}

	product, err := h.productService.UpdateProduct(c.Request.Context(), userID, productID, &req)
	if err != nil {
		if err == errors.ErrNotFound {
			utils.ErrorResponse(c, errors.ErrNotFound)
			return
		}
		h.logger.Error("Failed to update product", zap.String("product_id", productID), zap.Error(err))
		utils.ErrorResponse(c, err)
		return
	}

	response := h.mapProductToResponse(product)
	utils.SuccessResponse(c, response)
}

// DeleteProduct handles product deletion
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		utils.ErrorResponse(c, errors.ErrUnauthorized)
		return
	}

	productID := c.Param("id")
	if productID == "" {
		utils.ErrorResponse(c, errors.ErrInvalidRequest)
		return
	}

	err := h.productService.DeleteProduct(c.Request.Context(), userID, productID)
	if err != nil {
		if err == errors.ErrNotFound {
			utils.ErrorResponse(c, errors.ErrNotFound)
			return
		}
		h.logger.Error("Failed to delete product", zap.String("product_id", productID), zap.Error(err))
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "Product deleted successfully"})
}

// UpdateInventory handles inventory updates
func (h *ProductHandler) UpdateInventory(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		utils.ErrorResponse(c, errors.ErrUnauthorized)
		return
	}

	productID := c.Param("id")
	if productID == "" {
		utils.ErrorResponse(c, errors.ErrInvalidRequest)
		return
	}

	var req models.InventoryUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, errors.ErrInvalidRequestBody)
		return
	}

	err := h.productService.UpdateInventory(c.Request.Context(), productID, &req)
	if err != nil {
		if err == errors.ErrNotFound {
			utils.ErrorResponse(c, errors.ErrNotFound)
			return
		}
		h.logger.Error("Failed to update inventory", zap.String("product_id", productID), zap.Error(err))
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "Inventory updated successfully"})
}

// GetMyProducts handles getting products for the authenticated user
func (h *ProductHandler) GetMyProducts(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		utils.ErrorResponse(c, errors.ErrUnauthorized)
		return
	}

	// Parse pagination parameters
	page := 1
	pageSize := 20

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	if ps := c.Query("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}

	response, err := h.productService.GetProductsBySeller(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		h.logger.Error("Failed to get user products", zap.Error(err))
		utils.ErrorResponse(c, err)
		return
	}

	// Map products to responses
	productResponses := make([]models.ProductResponse, len(response.Products))
	for i, product := range response.Products {
		productResponses[i] = *h.mapProductToResponse(&response.Products[i])
	}

	utils.SuccessResponse(c, gin.H{
		"products":  productResponses,
		"total":     response.Total,
		"page":      response.Page,
		"page_size": response.PageSize,
		"has_next":  response.HasNext,
	})
}

// mapProductToResponse maps a Product model to ProductResponse
func (h *ProductHandler) mapProductToResponse(product *models.Product) *models.ProductResponse {
	return &models.ProductResponse{
		Product:   *product,
		Available: product.GetAvailable(),
	}
}
