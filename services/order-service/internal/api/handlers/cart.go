package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/gmsas95/blytz-mvp/services/order-service/internal/models"
	"github.com/gmsas95/blytz-mvp/services/order-service/internal/services"
)

type CartHandler struct {
	orderService *services.OrderService
	logger       *zap.Logger
}

func NewCartHandler(orderService *services.OrderService, logger *zap.Logger) *CartHandler {
	return &CartHandler{
		orderService: orderService,
		logger:       logger,
	}
}

// GetCart retrieves the user's cart
func (h *CartHandler) GetCart(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	db := h.orderService.GetDB()

	var cart models.Cart
	err := db.Preload("Items").Where("user_id = ?", userID).First(&cart).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create empty cart for user
			cart = models.Cart{
				UserID:    userID,
				Items:     []models.CartItem{},
				Total:     0,
				ItemCount: 0,
			}
			if createErr := db.Create(&cart).Error; createErr != nil {
				h.logger.Error("Failed to create cart", zap.Error(createErr))
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cart"})
				return
			}
		} else {
			h.logger.Error("Failed to get cart", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cart"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": cart})
}

// AddToCart adds an item to the user's cart
func (h *CartHandler) AddToCart(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req struct {
		ProductID string `json:"productId" binding:"required"`
		Quantity  int    `json:"quantity" binding:"required,min=1"`
		AuctionID string `json:"auctionId,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := h.orderService.GetDB()

	// Get or create user's cart
	var cart models.Cart
	err := db.Preload("Items").Where("user_id = ?", userID).First(&cart).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			cart = models.Cart{
				UserID:    userID,
				Items:     []models.CartItem{},
				Total:     0,
				ItemCount: 0,
			}
			if createErr := db.Create(&cart).Error; createErr != nil {
				h.logger.Error("Failed to create cart", zap.Error(createErr))
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cart"})
				return
			}
		} else {
			h.logger.Error("Failed to get cart", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cart"})
			return
		}
	}

	// Check if item already exists in cart
	var existingItem models.CartItem
	query := db.Where("cart_id = ? AND product_id = ?", cart.ID, req.ProductID)
	if req.AuctionID != "" {
		query = query.Where("auction_id = ?", req.AuctionID)
	} else {
		query = query.Where("auction_id IS NULL")
	}

	err = query.First(&existingItem).Error
	if err == nil {
		// Update existing item quantity
		existingItem.Quantity += req.Quantity
		existingItem.Total = int64(existingItem.Quantity) * existingItem.Price

		if updateErr := db.Save(&existingItem).Error; updateErr != nil {
			h.logger.Error("Failed to update cart item", zap.Error(updateErr))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart item"})
			return
		}
	} else if err == gorm.ErrRecordNotFound {
		// For demo purposes, we'll use a mock price
		// In a real implementation, you'd fetch this from the product service
		price := int64(10000) // $100.00 in cents

		// Create new cart item
		newItem := models.CartItem{
			ID:        uuid.New().String(),
			CartID:    cart.ID,
			ProductID: req.ProductID,
			AuctionID: &req.AuctionID,
			Quantity:  req.Quantity,
			Price:     price,
			Total:     int64(req.Quantity) * price,
		}

		if req.AuctionID == "" {
			newItem.AuctionID = nil
		}

		if createErr := db.Create(&newItem).Error; createErr != nil {
			h.logger.Error("Failed to create cart item", zap.Error(createErr))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add item to cart"})
			return
		}
	} else {
		h.logger.Error("Failed to check existing cart item", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check cart item"})
		return
	}

	// Recalculate cart totals
	if recalcErr := h.recalculateCartTotals(db, cart.ID); recalcErr != nil {
		h.logger.Error("Failed to recalculate cart totals", zap.Error(recalcErr))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart totals"})
		return
	}

	// Get updated cart
	var updatedCart models.Cart
	if getErr := db.Preload("Items").First(&updatedCart, cart.ID).Error; getErr != nil {
		h.logger.Error("Failed to get updated cart", zap.Error(getErr))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": updatedCart})
}

// RemoveFromCart removes an item from the user's cart
func (h *CartHandler) RemoveFromCart(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	itemID := c.Param("itemId")
	if itemID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Item ID is required"})
		return
	}

	db := h.orderService.GetDB()

	// Verify the cart item belongs to the user
	var cartItem models.CartItem
	err := db.Joins("JOIN carts ON cart_items.cart_id = carts.id").
		Where("cart_items.id = ? AND carts.user_id = ?", itemID, userID).
		First(&cartItem).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cart item not found"})
			return
		}
		h.logger.Error("Failed to get cart item", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cart item"})
		return
	}

	// Remove the item
	if deleteErr := db.Delete(&cartItem).Error; deleteErr != nil {
		h.logger.Error("Failed to remove cart item", zap.Error(deleteErr))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove item from cart"})
		return
	}

	// Recalculate cart totals
	if recalcErr := h.recalculateCartTotals(db, cartItem.CartID); recalcErr != nil {
		h.logger.Error("Failed to recalculate cart totals", zap.Error(recalcErr))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart totals"})
		return
	}

	// Get updated cart
	var updatedCart models.Cart
	if getErr := db.Preload("Items").First(&updatedCart, cartItem.CartID).Error; getErr != nil {
		h.logger.Error("Failed to get updated cart", zap.Error(getErr))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": updatedCart})
}

// UpdateCartItemQuantity updates the quantity of a cart item
func (h *CartHandler) UpdateCartItemQuantity(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	itemID := c.Param("itemId")
	if itemID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Item ID is required"})
		return
	}

	var req struct {
		Quantity int `json:"quantity" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := h.orderService.GetDB()

	// Verify the cart item belongs to the user
	var cartItem models.CartItem
	err := db.Joins("JOIN carts ON cart_items.cart_id = carts.id").
		Where("cart_items.id = ? AND carts.user_id = ?", itemID, userID).
		First(&cartItem).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cart item not found"})
			return
		}
		h.logger.Error("Failed to get cart item", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cart item"})
		return
	}

	// Update quantity and total
	cartItem.Quantity = req.Quantity
	cartItem.Total = int64(req.Quantity) * cartItem.Price

	if updateErr := db.Save(&cartItem).Error; updateErr != nil {
		h.logger.Error("Failed to update cart item", zap.Error(updateErr))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart item"})
		return
	}

	// Recalculate cart totals
	if recalcErr := h.recalculateCartTotals(db, cartItem.CartID); recalcErr != nil {
		h.logger.Error("Failed to recalculate cart totals", zap.Error(recalcErr))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart totals"})
		return
	}

	// Get updated cart
	var updatedCart models.Cart
	if getErr := db.Preload("Items").First(&updatedCart, cartItem.CartID).Error; getErr != nil {
		h.logger.Error("Failed to get updated cart", zap.Error(getErr))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": updatedCart})
}

// ClearCart removes all items from the user's cart
func (h *CartHandler) ClearCart(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	db := h.orderService.GetDB()

	// Get user's cart
	var cart models.Cart
	err := db.Where("user_id = ?", userID).First(&cart).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
			return
		}
		h.logger.Error("Failed to get cart", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve cart"})
		return
	}

	// Delete all cart items
	if deleteErr := db.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error; deleteErr != nil {
		h.logger.Error("Failed to clear cart items", zap.Error(deleteErr))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear cart"})
		return
	}

	// Update cart totals
	cart.Total = 0
	cart.ItemCount = 0
	if updateErr := db.Save(&cart).Error; updateErr != nil {
		h.logger.Error("Failed to update cart", zap.Error(updateErr))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": cart})
}

// Helper method to recalculate cart totals
func (h *CartHandler) recalculateCartTotals(db *gorm.DB, cartID string) error {
	var cartItems []models.CartItem
	if err := db.Where("cart_id = ?", cartID).Find(&cartItems).Error; err != nil {
		return err
	}

	var total int64 = 0
	var itemCount int = 0

	for _, item := range cartItems {
		total += item.Total
		itemCount += item.Quantity
	}

	// Update cart
	return db.Model(&models.Cart{}).Where("id = ?", cartID).Updates(map[string]interface{}{
		"total":      total,
		"item_count": itemCount,
	}).Error
}
