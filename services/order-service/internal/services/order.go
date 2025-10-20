package services

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/gmsas95/blytz-mvp/services/order-service/internal/config"
	"github.com/gmsas95/blytz-mvp/services/order-service/internal/models"
	"github.com/gmsas95/blytz-mvp/shared/pkg/errors"
)

type OrderService struct {
	db     *gorm.DB
	logger *zap.Logger
	config *config.Config
}

func NewOrderService(db *gorm.DB, logger *zap.Logger, config *config.Config) *OrderService {
	return &OrderService{
		db:     db,
		logger: logger,
		config: config,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, userID string, req *CreateOrderRequest) (*models.Order, error) {
	s.logger.Info("Creating order", zap.String("user_id", userID), zap.String("product_id", req.ProductID))

	// Validate request
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Calculate total amount
	totalAmount := req.Price * int64(req.Quantity)

	// Create order
	order := &models.Order{
		UserID:       userID,
		AuctionID:    req.AuctionID,
		ProductID:    req.ProductID,
		ProductName:  req.ProductName,
		ProductImage: req.ProductImage,
		Quantity:     req.Quantity,
		Price:        req.Price,
		TotalAmount:  totalAmount,
		Currency:     req.Currency,
		Status:       string(models.OrderStatusPending),
		PaymentStatus: string(models.PaymentStatusPending),
		ShippingAddress: models.Address{
			Name:        req.ShippingAddress.Name,
			Street:      req.ShippingAddress.Street,
			City:        req.ShippingAddress.City,
			State:       req.ShippingAddress.State,
			PostalCode:  req.ShippingAddress.PostalCode,
			Country:     req.ShippingAddress.Country,
			PhoneNumber: req.ShippingAddress.PhoneNumber,
		},
		BillingAddress: models.Address{
			Name:        req.BillingAddress.Name,
			Street:      req.BillingAddress.Street,
			City:        req.BillingAddress.City,
			State:       req.BillingAddress.State,
			PostalCode:  req.BillingAddress.PostalCode,
			Country:     req.BillingAddress.Country,
			PhoneNumber: req.BillingAddress.PhoneNumber,
		},
		Notes: req.Notes,
	}

	if err := s.db.Create(order).Error; err != nil {
		s.logger.Error("Failed to create order", zap.Error(err))
		return nil, errors.ErrInternalServer
	}

	s.logger.Info("Order created successfully", zap.String("order_id", order.ID))
	return order, nil
}

func (s *OrderService) GetOrder(ctx context.Context, orderID string, userID string) (*models.Order, error) {
	s.logger.Info("Getting order", zap.String("order_id", orderID), zap.String("user_id", userID))

	var order models.Order
	if err := s.db.Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		s.logger.Error("Failed to get order", zap.Error(err))
		return nil, errors.ErrInternalServer
	}

	return &order, nil
}

func (s *OrderService) GetUserOrders(ctx context.Context, userID string, limit, offset int) ([]*models.Order, int64, error) {
	s.logger.Info("Getting user orders", zap.String("user_id", userID))

	var orders []*models.Order
	var total int64

	query := s.db.Where("user_id = ?", userID)

	// Get total count
	if err := query.Model(&models.Order{}).Count(&total).Error; err != nil {
		s.logger.Error("Failed to count user orders", zap.Error(err))
		return nil, 0, errors.ErrInternalServer
	}

	// Get orders with pagination
	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&orders).Error; err != nil {
		s.logger.Error("Failed to get user orders", zap.Error(err))
		return nil, 0, errors.ErrInternalServer
	}

	return orders, total, nil
}

func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderID string, userID string, status models.OrderStatus) (*models.Order, error) {
	s.logger.Info("Updating order status", zap.String("order_id", orderID), zap.String("user_id", userID), zap.String("status", string(status)))

	var order models.Order
	if err := s.db.Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		s.logger.Error("Failed to get order for status update", zap.Error(err))
		return nil, errors.ErrInternalServer
	}

	// Validate status transition
	if !s.isValidStatusTransition(order.Status, string(status)) {
		return nil, errors.ErrBadRequest
	}

	order.Status = string(status)
	order.UpdatedAt = time.Now()

	if err := s.db.Save(&order).Error; err != nil {
		s.logger.Error("Failed to update order status", zap.Error(err))
		return nil, errors.ErrInternalServer
	}

	s.logger.Info("Order status updated successfully", zap.String("order_id", order.ID))
	return &order, nil
}

func (s *OrderService) CancelOrder(ctx context.Context, orderID string, userID string) error {
	s.logger.Info("Cancelling order", zap.String("order_id", orderID), zap.String("user_id", userID))

	var order models.Order
	if err := s.db.Where("id = ? AND user_id = ?", orderID, userID).First(&order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.ErrNotFound
		}
		s.logger.Error("Failed to get order for cancellation", zap.Error(err))
		return errors.ErrInternalServer
	}

	// Check if order can be cancelled
	if !s.canCancelOrder(order.Status) {
		return fmt.Errorf("order cannot be cancelled in current status: %s", order.Status)
	}

	order.Status = string(models.OrderStatusCancelled)
	order.PaymentStatus = string(models.PaymentStatusCancelled)
	order.UpdatedAt = time.Now()

	if err := s.db.Save(&order).Error; err != nil {
		s.logger.Error("Failed to cancel order", zap.Error(err))
		return errors.ErrInternalServer
	}

	s.logger.Info("Order cancelled successfully", zap.String("order_id", order.ID))
	return nil
}

func (s *OrderService) isValidStatusTransition(currentStatus, newStatus string) bool {
	// Define valid status transitions
	validTransitions := map[string][]string{
		string(models.OrderStatusPending):    {string(models.OrderStatusProcessing), string(models.OrderStatusCancelled)},
		string(models.OrderStatusProcessing): {string(models.OrderStatusConfirmed), string(models.OrderStatusCancelled)},
		string(models.OrderStatusConfirmed):  {string(models.OrderStatusShipped), string(models.OrderStatusCancelled)},
		string(models.OrderStatusShipped):    {string(models.OrderStatusDelivered)},
		string(models.OrderStatusDelivered):  {},
		string(models.OrderStatusCancelled):  {},
		string(models.OrderStatusRefunded):   {},
	}

	allowedStatuses, exists := validTransitions[currentStatus]
	if !exists {
		return false
	}

	for _, allowed := range allowedStatuses {
		if allowed == newStatus {
			return true
		}
	}

	return false
}

func (s *OrderService) canCancelOrder(status string) bool {
	// Orders can only be cancelled if they're not already shipped, delivered, or cancelled
	return status == string(models.OrderStatusPending) ||
		status == string(models.OrderStatusProcessing) ||
		status == string(models.OrderStatusConfirmed)
}