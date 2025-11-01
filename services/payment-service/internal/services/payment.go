package services

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/gmsas95/blytz-mvp/services/payment-service/internal/config"
	"github.com/gmsas95/blytz-mvp/services/payment-service/internal/models"
	"github.com/gmsas95/blytz-mvp/services/payment-service/pkg/fiuu"
	"github.com/gmsas95/blytz-mvp/shared/pkg/errors"
)

type PaymentService struct {
	db     *gorm.DB
	logger *zap.Logger
	config *config.Config
	fiuu   *fiuu.Client
}

func NewPaymentService(db *gorm.DB, logger *zap.Logger, config *config.Config) *PaymentService {
	// Initialize Fiuu client if credentials are provided
	var fiuuClient *fiuu.Client
	if config.FiuuMerchantID != "" && config.FiuuVerifyKey != "" {
		fiuuClient = fiuu.NewClient(config.FiuuMerchantID, config.FiuuVerifyKey, config.FiuuSandbox, logger)
		logger.Info("Fiuu payment client initialized",
			zap.String("merchant_id", config.FiuuMerchantID),
			zap.Bool("sandbox", config.FiuuSandbox))
	} else {
		logger.Warn("Fiuu credentials not provided, using mock payment processing")
	}

	return &PaymentService{
		db:     db,
		logger: logger,
		config: config,
		fiuu:   fiuuClient,
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
		Channel:       req.Channel,
	}

	if err := s.db.Create(payment).Error; err != nil {
		s.logger.Error("Failed to create payment record", zap.Error(err))
		return nil, errors.ErrInternalServer
	}

	// Process payment with provider
	if s.fiuu != nil && req.Provider == "fiuu" {
		// Use Fiuu payment processing
		providerID, transactionID, paymentRefID, err := s.processWithFiuu(ctx, req)
		if err != nil {
			payment.Status = string(models.PaymentStatusFailed)
			payment.FailureReason = err.Error()
			s.db.Save(payment)
			s.logger.Error("Fiuu payment processing failed", zap.Error(err))
			return nil, err
		}

		// Update payment with Fiuu details
		payment.Status = string(models.PaymentStatusCompleted)
		payment.ProviderID = providerID
		payment.TransactionID = transactionID
		payment.PaymentRefID = paymentRefID
	} else {
		// Use mock payment processing
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
	}

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

	// Process refund with provider
	if s.fiuu != nil && payment.Provider == "fiuu" {
		// Use Fiuu refund processing
		if err := s.processRefundWithFiuu(ctx, &payment, amount, reason); err != nil {
			s.logger.Error("Fiuu refund processing failed", zap.Error(err))
			return nil, err
		}
	} else {
		// Use mock refund processing
		if err := s.processRefundWithProvider(payment.ProviderID, amount); err != nil {
			s.logger.Error("Refund processing failed", zap.Error(err))
			return nil, err
		}
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
	var methods []*models.PaymentMethodInfo

	// Add Fiuu payment methods if client is available
	if s.fiuu != nil {
		fiuuChannels := fiuu.GetAvailableChannels()
		for _, channel := range fiuuChannels {
			if channel.Enabled {
				methods = append(methods, &models.PaymentMethodInfo{
					ID:          string(channel.Code),
					Name:        channel.Name,
					Type:        channel.Type,
					Description: channel.Description,
					Enabled:     channel.Enabled,
					Channel:     string(channel.Code),
					Currency:    string(channel.Currency),
				})
			}
		}
	} else {
		// Fallback to basic payment methods
		methods = append(methods, &models.PaymentMethodInfo{
			ID:          "card",
			Name:        "Credit/Debit Card",
			Type:        "card",
			Description: "Pay with credit or debit card",
			Enabled:     true,
		})
		methods = append(methods, &models.PaymentMethodInfo{
			ID:          "paypal",
			Name:        "PayPal",
			Type:        "wallet",
			Description: "Pay with PayPal",
			Enabled:     true,
		})
		methods = append(methods, &models.PaymentMethodInfo{
			ID:          "apple_pay",
			Name:        "Apple Pay",
			Type:        "wallet",
			Description: "Pay with Apple Pay",
			Enabled:     false, // Disabled for now
		})
		methods = append(methods, &models.PaymentMethodInfo{
			ID:          "google_pay",
			Name:        "Google Pay",
			Type:        "wallet",
			Description: "Pay with Google Pay",
			Enabled:     false, // Disabled for now
		})
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

// processWithFiuu processes payment using Fiuu payment gateway
func (s *PaymentService) processWithFiuu(ctx context.Context, req *models.ProcessPaymentRequest) (string, string, string, error) {
	// Convert amount from cents to ringgit
	amount := float64(req.Amount) / 100.0

	// Create Fiuu payment request
	fiuuReq := fiuu.PaymentRequest{
		Channel:     fiuu.Channel(req.Channel),
		Amount:      amount,
		OrderID:     req.OrderID,
		BillName:    req.BillName,
		BillEmail:   req.BillEmail,
		BillMobile:  req.BillMobile,
		BillDesc:    req.BillDesc,
		Currency:    fiuu.Currency(req.Currency),
		LangCode:    "en",
		ReturnURL:   s.config.FiuuReturnURL,
		NotifyURL:   s.config.FiuuNotifyURL,
		CallbackURL: s.config.FiuuCallbackURL,
		CancelURL:   s.config.FiuuCancelURL,
	}

	// Process payment with Fiuu
	resp, err := s.fiuu.CreatePayment(ctx, fiuuReq)
	if err != nil {
		return "", "", "", fmt.Errorf("fiuu payment failed: %w", err)
	}

	return resp.TransactionID, resp.TransactionID, resp.PaymentRefID, nil
}

// processRefundWithFiuu processes refund using Fiuu payment gateway
func (s *PaymentService) processRefundWithFiuu(ctx context.Context, payment *models.Payment, amount int64, reason string) error {
	// Convert amount from cents to ringgit
	amountFloat := float64(amount) / 100.0

	// Generate unique refund ID
	refundID := fmt.Sprintf("refund_%s_%d", payment.ID, time.Now().Unix())

	// Create Fiuu refund request
	fiuuReq := fiuu.RefundRequest{
		OrderID:      payment.OrderID,
		Amount:       amountFloat,
		RefundID:     refundID,
		RefundReason: reason,
	}

	// Process refund with Fiuu
	_, err := s.fiuu.CreateRefund(ctx, fiuuReq)
	if err != nil {
		return fmt.Errorf("fiuu refund failed: %w", err)
	}

	return nil
}

// GetSeamlessConfig returns configuration for frontend seamless integration
func (s *PaymentService) GetSeamlessConfig(orderID string, amount int64, billName, billEmail, billMobile, billDesc, channel string) (*models.FiuuSeamlessConfig, error) {
	if s.fiuu == nil {
		return nil, fmt.Errorf("fiuu client not initialized")
	}

	// Convert amount from cents to ringgit
	amountFloat := float64(amount) / 100.0

	// Get seamless config from Fiuu client
	config := s.fiuu.GetSeamlessConfig(orderID, amountFloat, billName, billEmail, billMobile, billDesc, fiuu.Channel(channel))

	// Convert to our model
	seamlessConfig := &models.FiuuSeamlessConfig{
		MerchantID: config["mpsmerchantid"].(string),
		Channel:    config["mpschannel"].(string),
		Amount:     config["mpsamount"].(string),
		OrderID:    config["mpsorderid"].(string),
		BillName:   config["mpsbill_name"].(string),
		BillEmail:  config["mpsbill_email"].(string),
		BillMobile: config["mpsbill_mobile"].(string),
		BillDesc:   config["mpsbill_desc"].(string),
		Currency:   config["mpscurrency"].(string),
		LangCode:   config["mpslangcode"].(string),
		VCode:      config["vcode"].(string),
		Sandbox:    config["sandbox"].(bool),
	}

	return seamlessConfig, nil
}

// ProcessWebhook handles Fiuu webhook notifications
func (s *PaymentService) ProcessWebhook(ctx context.Context, webhook *models.FiuuWebhookRequest) error {
	s.logger.Info("Processing Fiuu webhook",
		zap.String("transaction_id", webhook.TransactionID),
		zap.String("order_id", webhook.OrderID),
		zap.String("payment_status", webhook.PaymentStatus))

	// Find payment by order ID
	var payment models.Payment
	if err := s.db.Where("order_id = ?", webhook.OrderID).First(&payment).Error; err != nil {
		s.logger.Error("Payment not found for webhook", zap.String("order_id", webhook.OrderID))
		return fmt.Errorf("payment not found: %w", err)
	}

	// Update payment details
	payment.TransactionID = webhook.TransactionID
	payment.PaymentRefID = webhook.PaymentRefID
	payment.Channel = webhook.ChannelCode

	// Update payment status based on webhook
	switch webhook.PaymentStatus {
	case "00": // Successful
		payment.Status = string(models.PaymentStatusCompleted)
	case "11": // Pending
		payment.Status = string(models.PaymentStatusProcessing)
	default: // Failed
		payment.Status = string(models.PaymentStatusFailed)
		if webhook.ErrorDescription != "" {
			payment.FailureReason = webhook.ErrorDescription
		}
	}

	// Save updated payment
	if err := s.db.Save(&payment).Error; err != nil {
		s.logger.Error("Failed to update payment from webhook", zap.Error(err))
		return fmt.Errorf("failed to update payment: %w", err)
	}

	s.logger.Info("Webhook processed successfully",
		zap.String("payment_id", payment.ID),
		zap.String("status", payment.Status))

	return nil
}
