package fiuu

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"go.uber.org/zap"
)

// Client represents Fiuu payment client
type Client struct {
	merchantID string
	verifyKey  string
	sandbox    bool
	httpClient *http.Client
	logger     *zap.Logger
	baseURL    string
}

// NewClient creates a new Fiuu payment client
func NewClient(merchantID, verifyKey string, sandbox bool, logger *zap.Logger) *Client {
	baseURL := "https://pay.fiuu.com"
	if sandbox {
		baseURL = "https://sb-pay.fiuu.com"
	}

	return &Client{
		merchantID: merchantID,
		verifyKey:  verifyKey,
		sandbox:    sandbox,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger:  logger,
		baseURL: baseURL,
	}
}

// CreatePayment creates a new payment with Fiuu
func (c *Client) CreatePayment(ctx context.Context, req PaymentRequest) (*PaymentResponse, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid payment request: %w", err)
	}

	// Set merchant ID
	req.MerchantID = c.merchantID

	// Generate vcode
	req.VCode = GenerateVCode(req, c.verifyKey)

	// Prepare form data
	formData := req.ToFormData()

	// Create HTTP request
	url := fmt.Sprintf("%s/RMS/API/payment/PaymentRequest", c.baseURL)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBufferString(formData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Make request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		c.logger.Error("Failed to create Fiuu payment", zap.Error(err))
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse response
	var paymentResp PaymentResponse
	if err := json.Unmarshal(body, &paymentResp); err != nil {
		c.logger.Error("Failed to parse Fiuu response", zap.String("body", string(body)), zap.Error(err))
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for errors
	if resp.StatusCode != http.StatusOK || paymentResp.ErrorCode != "" {
		c.logger.Error("Fiuu payment error",
			zap.String("error_code", paymentResp.ErrorCode),
			zap.String("error_desc", paymentResp.ErrorDescription),
			zap.Int("status_code", resp.StatusCode))
		return &paymentResp, fmt.Errorf("payment failed: %s - %s", paymentResp.ErrorCode, paymentResp.ErrorDescription)
	}

	c.logger.Info("Fiuu payment created successfully",
		zap.String("transaction_id", paymentResp.TransactionID),
		zap.String("order_id", paymentResp.OrderID))

	return &paymentResp, nil
}

// CreateRefund creates a refund with Fiuu
func (c *Client) CreateRefund(ctx context.Context, req RefundRequest) (*RefundResponse, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid refund request: %w", err)
	}

	// Set merchant ID
	req.MerchantID = c.merchantID

	// Generate vcode
	req.VCode = GenerateRefundVCode(req, c.verifyKey)

	// Prepare form data
	formData := url.Values{}
	formData.Set("merchantid", req.MerchantID)
	formData.Set("orderid", req.OrderID)
	formData.Set("amount", fmt.Sprintf("%.2f", req.Amount))
	formData.Set("refundid", req.RefundID)
	formData.Set("refundreason", req.RefundReason)
	formData.Set("vcode", req.VCode)

	// Create HTTP request
	url := fmt.Sprintf("%s/RMS/API/payment/RefundRequest", c.baseURL)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create refund request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Make request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		c.logger.Error("Failed to create Fiuu refund", zap.Error(err))
		return nil, fmt.Errorf("failed to create refund: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read refund response: %w", err)
	}

	// Parse response
	var refundResp RefundResponse
	if err := json.Unmarshal(body, &refundResp); err != nil {
		c.logger.Error("Failed to parse Fiuu refund response", zap.String("body", string(body)), zap.Error(err))
		return nil, fmt.Errorf("failed to parse refund response: %w", err)
	}

	// Check for errors
	if resp.StatusCode != http.StatusOK || refundResp.ErrorCode != "" {
		c.logger.Error("Fiuu refund error",
			zap.String("error_code", refundResp.ErrorCode),
			zap.String("error_desc", refundResp.ErrorDescription),
			zap.Int("status_code", resp.StatusCode))
		return &refundResp, fmt.Errorf("refund failed: %s - %s", refundResp.ErrorCode, refundResp.ErrorDescription)
	}

	c.logger.Info("Fiuu refund created successfully",
		zap.String("refund_id", refundResp.RefundID),
		zap.String("order_id", refundResp.OrderID))

	return &refundResp, nil
}

// GetPaymentStatus retrieves payment status from Fiuu
func (c *Client) GetPaymentStatus(ctx context.Context, orderID string) (*PaymentResponse, error) {
	// Prepare request data
	formData := url.Values{}
	formData.Set("merchantid", c.merchantID)
	formData.Set("orderid", orderID)

	// Create HTTP request
	url := fmt.Sprintf("%s/RMS/API/payment/QueryPaymentStatus", c.baseURL)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create status request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Make request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		c.logger.Error("Failed to get Fiuu payment status", zap.Error(err))
		return nil, fmt.Errorf("failed to get payment status: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read status response: %w", err)
	}

	// Parse response
	var paymentResp PaymentResponse
	if err := json.Unmarshal(body, &paymentResp); err != nil {
		c.logger.Error("Failed to parse Fiuu status response", zap.String("body", string(body)), zap.Error(err))
		return nil, fmt.Errorf("failed to parse status response: %w", err)
	}

	// Check for errors
	if resp.StatusCode != http.StatusOK {
		c.logger.Error("Fiuu status check error",
			zap.String("error_code", paymentResp.ErrorCode),
			zap.String("error_desc", paymentResp.ErrorDescription),
			zap.Int("status_code", resp.StatusCode))
		return &paymentResp, fmt.Errorf("status check failed: %s - %s", paymentResp.ErrorCode, paymentResp.ErrorDescription)
	}

	return &paymentResp, nil
}

// GetSeamlessConfig returns configuration for frontend seamless integration
func (c *Client) GetSeamlessConfig(orderID string, amount float64, billName, billEmail, billMobile, billDesc string, channel Channel) map[string]interface{} {
	req := PaymentRequest{
		MerchantID: c.merchantID,
		Channel:    channel,
		Amount:     amount,
		OrderID:    orderID,
		BillName:   billName,
		BillEmail:  billEmail,
		BillMobile: billMobile,
		BillDesc:   billDesc,
		Currency:   CurrencyMYR, // Default to MYR
		LangCode:   "en",
	}

	// Generate vcode for frontend
	req.VCode = GenerateVCode(req, c.verifyKey)

	return map[string]interface{}{
		"mpsmerchantid":  req.MerchantID,
		"mpschannel":     string(req.Channel),
		"mpsamount":      fmt.Sprintf("%.2f", req.Amount),
		"mpsorderid":     req.OrderID,
		"mpsbill_name":   req.BillName,
		"mpsbill_email":  req.BillEmail,
		"mpsbill_mobile": req.BillMobile,
		"mpsbill_desc":   req.BillDesc,
		"mpscurrency":    string(req.Currency),
		"mpslangcode":    req.LangCode,
		"vcode":          req.VCode,
		"sandbox":        c.sandbox,
	}
}

// IsSandbox returns whether client is in sandbox mode
func (c *Client) IsSandbox() bool {
	return c.sandbox
}

// GetMerchantID returns the merchant ID
func (c *Client) GetMerchantID() string {
	return c.merchantID
}

// GetBaseURL returns the base URL
func (c *Client) GetBaseURL() string {
	return c.baseURL
}
