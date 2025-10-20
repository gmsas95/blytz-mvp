package services

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/gmsas95/blytz-mvp/services/payment-service/internal/config"
	"github.com/gmsas95/blytz-mvp/services/payment-service/internal/models"
	"github.com/gmsas95/blytz-mvp/shared/pkg/errors"
)

type PaymentService struct {
	db     *gorm.DB
	logger *zap.Logger
	config *config.Config
}

func NewPaymentService(db *gorm.DB, logger *zap.Logger, config *config.Config) *PaymentService {
	return &PaymentService{
		db:     db,
		logger: logger,
		config: config,
	}
}

func (s *PaymentService) ProcessPayment(ctx context.Context, userID string, req *models.ProcessPaymentRequest) (*models.Payment, error) {
	s.logger.Info("Processing payment", zap.String("user_id", userID), zap.String("order_id", req.OrderID))

	// Validate request
	if err := s.validatePaymentRequest(req); err != nil {
		return nil, err
	}

	// Create payment record
	payment := &models.Payment{
		UserID:        userID,
		OrderID:       req.OrderID,
		Amount:        req.Amount,
		Currency:      req.Currency,
		Status:        string(models.PaymentStatusProcessing),
		PaymentMethod: req.PaymentMethod,
		Provider:      req.Provider,
	}

	if err := s.db.Create(payment).Error; err != nil {
		s.logger.Error("Failed to create payment record", zap.Error(err))
		return nil, errors.ErrInternalServer
	}

	// Process payment with provider (mock implementation)
	providerID, err := s.processWithProvider(req)
	if err != nil {
		payment.Status = string(models.PaymentStatusFailed)
		payment.FailureReason = err.Error()
		s.db.Save(payment)
		s.logger.Error("Payment processing failed", zap.Error(err))
		return nil, err
	}

	// Update payment as completed
	payment.Status = string(models.PaymentStatusCompleted)
	payment.ProviderID = providerID
	if err := s.db.Save(payment).Error; err != nil {
		s.logger.Error("Failed to update payment status", zap.Error(err))
		return nil, errors.ErrInternalServer
	}

	s.logger.Info("Payment processed successfully", zap.String("payment_id", payment.ID))
	return payment, nil
}

func (s *PaymentService) GetPayment(ctx context.Context, paymentID string, userID string) (*models.Payment, error) {
	s.logger.Info("Getting payment", zap.String("payment_id", paymentID), zap.String("user_id", userID))

	var payment models.Payment
	if err := s.db.Where("id = ? AND user_id = ?", paymentID, userID).First(&payment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		s.logger.Error("Failed to get payment", zap.Error(err))
		return nil, errors.ErrInternalServer
	}

	return &payment, nil
}

func (s *PaymentService) GetPaymentHistory(ctx context.Context, userID string, limit int) ([]*models.Payment, error) {
	s.logger.Info("Getting payment history", zap.String("user_id", userID))

	var payments []*models.Payment
	if err := s.db.Where("user_id = ?", userID).Order("created_at DESC").Limit(limit).Find(&payments).Error; err != nil {
		s.logger.Error("Failed to get payment history", zap.Error(err))
		return nil, errors.ErrInternalServer
	}

	return payments, nil
}

func (s *PaymentService) ProcessRefund(ctx context.Context, paymentID string, userID string, amount int64, reason string) (*models.Payment, error) {
	s.logger.Info("Processing refund", zap.String("payment_id", paymentID), zap.String("user_id", userID), zap.Int64("amount", amount))

	var payment models.Payment
	if err := s.db.Where("id = ? AND user_id = ?", paymentID, userID).First(&payment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		s.logger.Error("Failed to get payment for refund", zap.Error(err))
		return nil, errors.ErrInternalServer
	}

	// Check if payment can be refunded
	if payment.Status != string(models.PaymentStatusCompleted) {
		return nil, fmt.Errorf("payment cannot be refunded in current status: %s", payment.Status)
	}

	if amount > payment.Amount {
		return nil, fmt.Errorf("refund amount cannot exceed payment amount")
	}

	// Process refund with provider (mock implementation)
	if err := s.processRefundWithProvider(payment.ProviderID, amount); err != nil {
		s.logger.Error("Refund processing failed", zap.Error(err))
		return nil, err
	}

	// Update payment record
	payment.RefundedAmount += amount
	if payment.RefundedAmount == payment.Amount {
		payment.Status = string(models.PaymentStatusRefunded)
	}
	now := time.Now()
	payment.RefundedAt = &now

	if err := s.db.Save(&payment).Error; err != nil {
		s.logger.Error("Failed to update payment after refund", zap.Error(err))
		return nil, errors.ErrInternalServer
	}

	s.logger.Info("Refund processed successfully", zap.String("payment_id", payment.ID))
	return &payment, nil
}

func (s *PaymentService) GetPaymentMethods(ctx context.Context) ([]*models.PaymentMethodInfo, error) {
	// Return available payment methods
	methods := []*models.PaymentMethodInfo{
		{
			ID:          "card",
			Name:        "Credit/Debit Card",
			Type:        "card",
			Description: "Pay with credit or debit card",
			Enabled:     true,
		},
		{
			ID:          "paypal",
			Name:        "PayPal",
			Type:        "wallet",
			Description: "Pay with PayPal",
			Enabled:     true,
		},
		{
			ID:          "apple_pay",
			Name:        "Apple Pay",
			Type:        "wallet",
			Description: "Pay with Apple Pay",
			Enabled:     false, // Disabled for now
		},
		{
			ID:          "google_pay",
			Name:        "Google Pay",
			Type:        "wallet",
			Description: "Pay with Google Pay",
			Enabled:     false, // Disabled for now
		},
	}

	return methods, nil
}

func (s *PaymentService) validatePaymentRequest(req *models.ProcessPaymentRequest) error {
	if req.Amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}
	if len(req.Currency) != 3 {
		return fmt.Errorf("currency must be 3 characters")
	}
	if req.Token == "" {
		return fmt.Errorf("payment token is required")
	}
	return nil
}

func (s *PaymentService) processWithProvider(req *models.ProcessPaymentRequest) (string, error) {
	// Mock payment processing - in real implementation, this would integrate with Stripe, PayPal, etc.
	s.logger.Info("Processing payment with provider", zap.String("provider", req.Provider))

	// Simulate processing time
	time.Sleep(100 * time.Millisecond)

	// Generate mock provider ID
	providerID := fmt.Sprintf("%s_%d", req.Provider, time.Now().UnixNano())

	return providerID, nil
}

func (s *PaymentService) processRefundWithProvider(providerID string, amount int64) error {
	// Mock refund processing - in real implementation, this would integrate with payment provider
	s.logger.Info("Processing refund with provider", zap.String("provider_id", providerID), zap.Int64("amount", amount))

	// Simulate processing time
	time.Sleep(100 * time.Millisecond)

	return nil
}