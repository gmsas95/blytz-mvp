package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/gmsas95/blytz-mvp/services/payment-service/internal/models"
	"github.com/gmsas95/blytz-mvp/services/payment-service/internal/services"
	"github.com/gmsas95/blytz-mvp/shared/pkg/errors"
	"github.com/gmsas95/blytz-mvp/shared/pkg/utils"
)

type PaymentHandler struct {
	paymentService *services.PaymentService
	logger         *zap.Logger
}

func NewPaymentHandler(paymentService *services.PaymentService, logger *zap.Logger) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
		logger:         logger,
	}
}

func (h *PaymentHandler) ProcessPayment(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		utils.RespondWithError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req models.ProcessPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	payment, err := h.paymentService.ProcessPayment(c.Request.Context(), userID, &req)
	if err != nil {
		h.logger.Error("Failed to process payment", zap.Error(err))
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to process payment")
		return
	}

	response := h.mapPaymentToResponse(payment)
	utils.RespondWithJSON(c, http.StatusCreated, response)
}

func (h *PaymentHandler) GetPayment(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		utils.RespondWithError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	paymentID := c.Param("id")
	if paymentID == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Payment ID is required")
		return
	}

	payment, err := h.paymentService.GetPayment(c.Request.Context(), paymentID, userID)
	if err != nil {
		if err == errors.ErrNotFound {
			utils.RespondWithError(c, http.StatusNotFound, "Payment not found")
			return
		}
		h.logger.Error("Failed to get payment", zap.Error(err))
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get payment")
		return
	}

	response := h.mapPaymentToResponse(payment)
	utils.RespondWithJSON(c, http.StatusOK, response)
}

func (h *PaymentHandler) GetPaymentHistory(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		utils.RespondWithError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	limit := 20
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	payments, err := h.paymentService.GetPaymentHistory(c.Request.Context(), userID, limit)
	if err != nil {
		h.logger.Error("Failed to get payment history", zap.Error(err))
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get payment history")
		return
	}

	response := make([]models.PaymentResponse, len(payments))
	for i, payment := range payments {
		response[i] = *h.mapPaymentToResponse(payment)
	}

	utils.RespondWithJSON(c, http.StatusOK, models.PaymentHistoryResponse{
		Payments: response,
		Total:    int64(len(payments)),
	})
}

func (h *PaymentHandler) GetPaymentMethods(c *gin.Context) {
	methods, err := h.paymentService.GetPaymentMethods(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to get payment methods", zap.Error(err))
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get payment methods")
		return
	}

	utils.RespondWithJSON(c, http.StatusOK, models.PaymentMethodsResponse{Methods: methods})
}

func (h *PaymentHandler) ProcessRefund(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		utils.RespondWithError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	paymentID := c.Param("id")
	if paymentID == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Payment ID is required")
		return
	}

	var req models.RefundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	payment, err := h.paymentService.ProcessRefund(c.Request.Context(), paymentID, userID, req.Amount, req.Reason)
	if err != nil {
		if err == errors.ErrNotFound {
			utils.RespondWithError(c, http.StatusNotFound, "Payment not found")
			return
		}
		h.logger.Error("Failed to process refund", zap.Error(err))
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to process refund")
		return
	}

	response := h.mapPaymentToResponse(payment)
	utils.RespondWithJSON(c, http.StatusOK, response)
}

func (h *PaymentHandler) mapPaymentToResponse(payment *models.Payment) *models.PaymentResponse {
	var refundedAt *string
	if payment.RefundedAt != nil {
		str := payment.RefundedAt.Format("2006-01-02T15:04:05Z")
		refundedAt = &str
	}

	return &models.PaymentResponse{
		ID:             payment.ID,
		UserID:         payment.UserID,
		OrderID:        payment.OrderID,
		Amount:         payment.Amount,
		Currency:       payment.Currency,
		Status:         payment.Status,
		PaymentMethod:  payment.PaymentMethod,
		Provider:       payment.Provider,
		ProviderID:     payment.ProviderID,
		FailureReason:  payment.FailureReason,
		RefundedAmount: payment.RefundedAmount,
		RefundedAt:     refundedAt,
		CreatedAt:      payment.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:      payment.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func (h *PaymentHandler) mapPaymentMethodToResponse(method *models.PaymentMethodInfo) models.PaymentMethodInfo {
	return models.PaymentMethodInfo{
		ID:          method.ID,
		Name:        method.Name,
		Type:        method.Type,
		Description: method.Description,
		Enabled:     method.Enabled,
	}
}