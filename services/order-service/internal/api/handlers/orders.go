package handlers

import (

	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/gmsas95/blytz-mvp/services/order-service/internal/models"
	"github.com/gmsas95/blytz-mvp/services/order-service/internal/services"
	"github.com/gmsas95/blytz-mvp/shared/pkg/errors"
	"github.com/gmsas95/blytz-mvp/shared/pkg/utils"
)

type OrderHandler struct {
	orderService *services.OrderService
	logger       *zap.Logger
}

func NewOrderHandler(orderService *services.OrderService, logger *zap.Logger) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
		logger:       logger,
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		utils.ErrorResponse(c, errors.ErrUnauthorized)
		return
	}

	var req services.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, errors.ErrInvalidRequestBody)
		return
	}

	order, err := h.orderService.CreateOrder(c.Request.Context(), userID, &req)
	if err != nil {
		h.logger.Error("Failed to create order", zap.Error(err))
		utils.ErrorResponse(c, err)
		return
	}

	response := h.mapOrderToResponse(order)
	utils.SuccessResponse(c, response)
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		utils.ErrorResponse(c, errors.ErrUnauthorized)
		return
	}

	orderID := c.Param("id")
	if orderID == "" {
		utils.ErrorResponse(c, errors.ErrInvalidRequest)
		return
	}

	order, err := h.orderService.GetOrder(c.Request.Context(), orderID, userID)
	if err != nil {
		if err == errors.ErrNotFound {
			utils.ErrorResponse(c, errors.ErrNotFound)
			return
		}
		h.logger.Error("Failed to get order", zap.Error(err))
		utils.ErrorResponse(c, err)
		return
	}

	response := h.mapOrderToResponse(order)
	utils.SuccessResponse(c, response)
}

func (h *OrderHandler) GetUserOrders(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		utils.ErrorResponse(c, errors.ErrUnauthorized)
		return
	}

	// Get pagination parameters
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

	offset := (page - 1) * pageSize

	orders, total, err := h.orderService.GetUserOrders(c.Request.Context(), userID, pageSize, offset)
	if err != nil {
		h.logger.Error("Failed to get user orders", zap.Error(err))
		utils.ErrorResponse(c, err)
		return
	}

	// Map orders to responses
	orderResponses := make([]services.OrderResponse, len(orders))
	for i, order := range orders {
		orderResponses[i] = *h.mapOrderToResponse(order)
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	response := services.OrdersListResponse{
		Orders:     orderResponses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}

	utils.SuccessResponse(c, response)
}

func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		utils.ErrorResponse(c, errors.ErrUnauthorized)
		return
	}

	orderID := c.Param("id")
	if orderID == "" {
		utils.ErrorResponse(c, errors.ErrInvalidRequest)
		return
	}

	var req services.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, errors.ErrInvalidRequestBody)
		return
	}

	order, err := h.orderService.UpdateOrderStatus(c.Request.Context(), orderID, userID, models.OrderStatus(req.Status))
	if err != nil {
		if err == errors.ErrNotFound {
			utils.ErrorResponse(c, errors.ErrNotFound)
			return
		}
		if err == errors.ErrInvalidRequest {
			utils.ErrorResponse(c, errors.ErrInvalidRequest)
			return
		}
		h.logger.Error("Failed to update order status", zap.Error(err))
		utils.ErrorResponse(c, err)
		return
	}

	response := h.mapOrderToResponse(order)
	utils.SuccessResponse(c, response)
}

func (h *OrderHandler) CancelOrder(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		utils.ErrorResponse(c, errors.ErrUnauthorized)
		return
	}

	orderID := c.Param("id")
	if orderID == "" {
		utils.ErrorResponse(c, errors.ErrInvalidRequest)
		return
	}

	err := h.orderService.CancelOrder(c.Request.Context(), orderID, userID)
	if err != nil {
		if err == errors.ErrNotFound {
			utils.ErrorResponse(c, errors.ErrNotFound)
			return
		}
		h.logger.Error("Failed to cancel order", zap.Error(err))
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "Order cancelled successfully"})
}

func (h *OrderHandler) mapOrderToResponse(order *models.Order) *services.OrderResponse {
	return &services.OrderResponse{
		ID:           order.ID,
		UserID:       order.UserID,
		AuctionID:    order.AuctionID,
		ProductID:    order.ProductID,
		ProductName:  order.ProductName,
		ProductImage: order.ProductImage,
		Quantity:     order.Quantity,
		Price:        order.Price,
		TotalAmount:  order.TotalAmount,
		Currency:     order.Currency,
		Status:       order.Status,
		PaymentStatus: order.PaymentStatus,
		PaymentMethod: order.PaymentMethod,
		ShippingAddress: services.AddressResponse{
			Name:        order.ShippingAddress.Name,
			Street:      order.ShippingAddress.Street,
			City:        order.ShippingAddress.City,
			State:       order.ShippingAddress.State,
			PostalCode:  order.ShippingAddress.PostalCode,
			Country:     order.ShippingAddress.Country,
			PhoneNumber: order.ShippingAddress.PhoneNumber,
		},
		BillingAddress: services.AddressResponse{
			Name:        order.BillingAddress.Name,
			Street:      order.BillingAddress.Street,
			City:        order.BillingAddress.City,
			State:       order.BillingAddress.State,
			PostalCode:  order.BillingAddress.PostalCode,
			Country:     order.BillingAddress.Country,
			PhoneNumber: order.BillingAddress.PhoneNumber,
		},
		Notes:     order.Notes,
		CreatedAt: order.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: order.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
