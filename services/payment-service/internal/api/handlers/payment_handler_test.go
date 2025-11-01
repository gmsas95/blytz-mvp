package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/gmsas95/blytz-mvp/services/payment-service/internal/models"
	"github.com/gmsas95/blytz-mvp/services/payment-service/internal/services"
	"github.com/gmsas95/blytz-mvp/services/payment-service/pkg/fiuu"
	"go.uber.org/zap/zaptest"
)

// MockPaymentService is a mock implementation of the PaymentService interface
type MockPaymentService struct {
	mock.Mock
}

func (m *MockPaymentService) CreatePayment(ctx context.Context, userID string, req *models.PaymentRequest) (*models.PaymentResponse, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PaymentResponse), args.Error(1)
}

func (m *MockPaymentService) GetPayment(ctx context.Context, userID, paymentID string) (*models.PaymentResponse, error) {
	args := m.Called(ctx, userID, paymentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PaymentResponse), args.Error(1)
}

func (m *MockPaymentService) GetPaymentHistory(ctx context.Context, userID string, page, limit int) ([]*models.PaymentResponse, int64, error) {
	args := m.Called(ctx, userID, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(2)
	}
	return args.Get(0).([]*models.PaymentResponse), args.Get(1).(int64), args.Error(2)
}

func (m *MockPaymentService) RefundPayment(ctx context.Context, userID, paymentID string, amount float64, reason string) (*models.RefundResponse, error) {
	args := m.Called(ctx, userID, paymentID, amount, reason)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RefundResponse), args.Error(1)
}

func (m *MockPaymentService) GetPaymentMethods(ctx context.Context, userID string) ([]*models.PaymentMethod, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.PaymentMethod), args.Error(1)
}

func (m *MockPaymentService) GetSeamlessConfig(ctx context.Context, userID, orderID string, amount float64, billName, billEmail, billMobile, billDesc string, channel fiuu.Channel) (map[string]interface{}, error) {
	args := m.Called(ctx, userID, orderID, amount, billName, billEmail, billMobile, billDesc, channel)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func (m *MockPaymentService) ProcessWebhook(ctx context.Context, webhookData *fiuu.WebhookData) error {
	args := m.Called(ctx, webhookData)
	return args.Error(0)
}

func setupTestRouter(mockService *MockPaymentService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	logger := zaptest.NewLogger(&testing.T{})
	handler := NewPaymentHandler(mockService, logger)

	// Setup routes
	v1 := router.Group("/api/v1")
	{
		payments := v1.Group("/payments")
		{
			payments.POST("", handler.CreatePayment)
			payments.GET("/:id", handler.GetPayment)
			payments.GET("", handler.GetPaymentHistory)
			payments.POST("/:id/refund", handler.RefundPayment)
			payments.GET("/methods", handler.GetPaymentMethods)
			payments.GET("/seamless/config", handler.GetSeamlessConfig)
		}

		v1.POST("/webhooks/fiuu", handler.ProcessFiuuWebhook)
	}

	return router
}

func TestPaymentHandler_CreatePayment_Success(t *testing.T) {
	mockService := new(MockPaymentService)
	router := setupTestRouter(mockService)

	// Mock successful payment creation
	expectedResp := &models.PaymentResponse{
		ID:            "payment123",
		OrderID:       "ORD123",
		UserID:        "user123",
		Amount:        100.50,
		Currency:      "MYR",
		Status:        "pending",
		PaymentMethod: "FPX",
		CreatedAt:     time.Now(),
	}

	mockService.On("CreatePayment", mock.Anything, "user123", mock.AnythingOfType("*models.PaymentRequest")).Return(expectedResp, nil)

	// Prepare request
	reqBody := models.PaymentRequest{
		OrderID:       "ORD123",
		Amount:        100.50,
		Currency:      "MYR",
		PaymentMethod: "FPX",
		BillName:      "John Doe",
		BillEmail:     "john@example.com",
		BillMobile:    "0123456789",
		BillDesc:      "Test Payment",
		ReturnURL:     "https://example.com/return",
		NotifyURL:     "https://example.com/notify",
	}

	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/payments", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Mock authentication context
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctx.Set("userID", "user123")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.PaymentResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, expectedResp.ID, response.ID)
	assert.Equal(t, expectedResp.OrderID, response.OrderID)
	assert.Equal(t, expectedResp.Amount, response.Amount)

	mockService.AssertExpectations(t)
}

func TestPaymentHandler_CreatePayment_ValidationError(t *testing.T) {
	mockService := new(MockPaymentService)
	router := setupTestRouter(mockService)

	// Prepare invalid request (missing required fields)
	reqBody := models.PaymentRequest{
		Amount: 100.50,
		// Missing other required fields
	}

	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/payments", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var errorResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	require.NoError(t, err)

	assert.Contains(t, errorResponse["error"], "validation failed")
}

func TestPaymentHandler_CreatePayment_ServiceError(t *testing.T) {
	mockService := new(MockPaymentService)
	router := setupTestRouter(mockService)

	// Mock service error
	mockService.On("CreatePayment", mock.Anything, "user123", mock.AnythingOfType("*models.PaymentRequest")).Return(nil, assert.AnError)

	// Prepare request
	reqBody := models.PaymentRequest{
		OrderID:       "ORD123",
		Amount:        100.50,
		Currency:      "MYR",
		PaymentMethod: "FPX",
		BillName:      "John Doe",
		BillEmail:     "john@example.com",
		BillMobile:    "0123456789",
		BillDesc:      "Test Payment",
	}

	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/payments", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var errorResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	require.NoError(t, err)

	assert.Contains(t, errorResponse["error"], "failed to create payment")
}

func TestPaymentHandler_GetPayment_Success(t *testing.T) {
	mockService := new(MockPaymentService)
	router := setupTestRouter(mockService)

	// Mock successful payment retrieval
	expectedResp := &models.PaymentResponse{
		ID:            "payment123",
		OrderID:       "ORD123",
		UserID:        "user123",
		Amount:        100.50,
		Currency:      "MYR",
		Status:        "completed",
		PaymentMethod: "FPX",
		CreatedAt:     time.Now(),
	}

	mockService.On("GetPayment", mock.Anything, "user123", "payment123").Return(expectedResp, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/payments/payment123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.PaymentResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, expectedResp.ID, response.ID)
	assert.Equal(t, expectedResp.OrderID, response.OrderID)

	mockService.AssertExpectations(t)
}

func TestPaymentHandler_GetPayment_NotFound(t *testing.T) {
	mockService := new(MockPaymentService)
	router := setupTestRouter(mockService)

	// Mock payment not found
	mockService.On("GetPayment", mock.Anything, "user123", "payment123").Return(nil, services.ErrPaymentNotFound)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/payments/payment123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var errorResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	require.NoError(t, err)

	assert.Contains(t, errorResponse["error"], "payment not found")
}

func TestPaymentHandler_GetPaymentHistory_Success(t *testing.T) {
	mockService := new(MockPaymentService)
	router := setupTestRouter(mockService)

	// Mock successful payment history retrieval
	expectedPayments := []*models.PaymentResponse{
		{
			ID:            "payment123",
			OrderID:       "ORD123",
			UserID:        "user123",
			Amount:        100.50,
			Currency:      "MYR",
			Status:        "completed",
			PaymentMethod: "FPX",
			CreatedAt:     time.Now(),
		},
		{
			ID:            "payment124",
			OrderID:       "ORD124",
			UserID:        "user123",
			Amount:        50.00,
			Currency:      "MYR",
			Status:        "pending",
			PaymentMethod: "GrabPay",
			CreatedAt:     time.Now(),
		},
	}

	mockService.On("GetPaymentHistory", mock.Anything, "user123", 1, 10).Return(expectedPayments, int64(2), nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/payments?page=1&limit=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, float64(2), response["total"])
	assert.Equal(t, float64(1), response["page"])
	assert.Equal(t, float64(10), response["limit"])

	payments := response["payments"].([]interface{})
	assert.Len(t, payments, 2)

	mockService.AssertExpectations(t)
}

func TestPaymentHandler_RefundPayment_Success(t *testing.T) {
	mockService := new(MockPaymentService)
	router := setupTestRouter(mockService)

	// Mock successful refund
	expectedResp := &models.RefundResponse{
		ID:            "refund123",
		PaymentID:     "payment123",
		Amount:        50.00,
		Currency:      "MYR",
		Status:        "processing",
		RefundReason:  "Customer requested refund",
		CreatedAt:     time.Now(),
	}

	mockService.On("RefundPayment", mock.Anything, "user123", "payment123", 50.00, "Customer requested refund").Return(expectedResp, nil)

	// Prepare request
	reqBody := map[string]interface{}{
		"amount":       50.00,
		"refund_reason": "Customer requested refund",
	}

	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/payments/payment123/refund", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.RefundResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, expectedResp.ID, response.ID)
	assert.Equal(t, expectedResp.PaymentID, response.PaymentID)
	assert.Equal(t, expectedResp.Amount, response.Amount)

	mockService.AssertExpectations(t)
}

func TestPaymentHandler_GetPaymentMethods_Success(t *testing.T) {
	mockService := new(MockPaymentService)
	router := setupTestRouter(mockService)

	// Mock successful payment methods retrieval
	expectedMethods := []*models.PaymentMethod{
		{
			ID:          "fpx",
			Name:        "FPX Online Banking",
			Type:        "bank_transfer",
			Enabled:     true,
			Icon:        "https://example.com/fpx.png",
			Currencies:  []string{"MYR"},
			Description: "Pay directly from your bank account",
		},
		{
			ID:          "grabpay",
			Name:        "GrabPay",
			Type:        "ewallet",
			Enabled:     true,
			Icon:        "https://example.com/grabpay.png",
			Currencies:  []string{"MYR"},
			Description: "Pay using your GrabPay e-wallet",
		},
	}

	mockService.On("GetPaymentMethods", mock.Anything, "user123").Return(expectedMethods, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/payments/methods", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.PaymentMethod
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Len(t, response, 2)
	assert.Equal(t, "fpx", response[0].ID)
	assert.Equal(t, "FPX Online Banking", response[0].Name)

	mockService.AssertExpectations(t)
}

func TestPaymentHandler_GetSeamlessConfig_Success(t *testing.T) {
	mockService := new(MockPaymentService)
	router := setupTestRouter(mockService)

	// Mock successful seamless config retrieval
	expectedConfig := map[string]interface{}{
		"mpsmerchantid":  "test_merchant",
		"mpschannel":     "FPX",
		"mpsamount":      "100.50",
		"mpsorderid":     "ORD123",
		"mpsbill_name":   "John Doe",
		"mpsbill_email":  "john@example.com",
		"mpsbill_mobile": "0123456789",
		"mpsbill_desc":   "Test Payment",
		"mpscurrency":    "MYR",
		"mpslangcode":    "en",
		"vcode":          "test_vcode",
		"sandbox":        true,
		"scriptUrl":      "https://sandbox.merchant.razer.com/RMS2/IPGSeamless/IPGSeamless.js",
	}

	mockService.On("GetSeamlessConfig", mock.Anything, "user123", "ORD123", 100.50, "John Doe", "john@example.com", "0123456789", "Test Payment", fiuu.ChannelFPX).Return(expectedConfig, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/payments/seamless/config?order_id=ORD123&amount=100.50&bill_name=John Doe&bill_email=john@example.com&bill_mobile=0123456789&bill_desc=Test Payment&channel=FPX", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "test_merchant", response["mpsmerchantid"])
	assert.Equal(t, "FPX", response["mpschannel"])
	assert.Equal(t, "100.50", response["mpsamount"])

	mockService.AssertExpectations(t)
}

func TestPaymentHandler_ProcessFiuuWebhook_Success(t *testing.T) {
	mockService := new(MockPaymentService)
	router := setupTestRouter(mockService)

	// Mock successful webhook processing
	mockService.On("ProcessWebhook", mock.Anything, mock.AnythingOfType("*fiuu.WebhookData")).Return(nil)

	// Prepare webhook request
	webhookData := fiuu.WebhookData{
		TransactionID:  "TX123456",
		OrderID:        "ORD123",
		Amount:         100.50,
		Currency:       "MYR",
		PaymentStatus:  "1", // Success
		ErrorCode:      "",
		ErrorDescription: "",
		Signature:      "test_signature",
	}

	jsonBody, _ := json.Marshal(webhookData)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/webhooks/fiuu", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "success", response["status"])

	mockService.AssertExpectations(t)
}

func TestPaymentHandler_ProcessFiuuWebhook_InvalidSignature(t *testing.T) {
	mockService := new(MockPaymentService)
	router := setupTestRouter(mockService)

	// Mock webhook validation error
	mockService.On("ProcessWebhook", mock.Anything, mock.AnythingOfType("*fiuu.WebhookData")).Return(services.ErrInvalidWebhookSignature)

	// Prepare webhook request with invalid signature
	webhookData := fiuu.WebhookData{
		TransactionID:  "TX123456",
		OrderID:        "ORD123",
		Amount:         100.50,
		Currency:       "MYR",
		PaymentStatus:  "1",
		Signature:      "invalid_signature",
	}

	jsonBody, _ := json.Marshal(webhookData)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/webhooks/fiuu", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var errorResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	require.NoError(t, err)

	assert.Contains(t, errorResponse["error"], "invalid webhook signature")
}

// Performance tests
func TestPaymentHandler_CreatePayment_Performance(t *testing.T) {
	mockService := new(MockPaymentService)
	router := setupTestRouter(mockService)

	// Mock successful payment creation
	expectedResp := &models.PaymentResponse{
		ID:            "payment123",
		OrderID:       "ORD123",
		UserID:        "user123",
		Amount:        100.50,
		Currency:      "MYR",
		Status:        "pending",
		PaymentMethod: "FPX",
		CreatedAt:     time.Now(),
	}

	mockService.On("CreatePayment", mock.Anything, "user123", mock.AnythingOfType("*models.PaymentRequest")).Return(expectedResp, nil)

	// Prepare request
	reqBody := models.PaymentRequest{
		OrderID:       "ORD123",
		Amount:        100.50,
		Currency:      "MYR",
		PaymentMethod: "FPX",
		BillName:      "John Doe",
		BillEmail:     "john@example.com",
		BillMobile:    "0123456789",
		BillDesc:      "Test Payment",
	}

	// Measure performance of multiple concurrent requests
	const numRequests := 50
	start := time.Now()

	for i := 0; i < numRequests; i++ {
		jsonBody, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/payments", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	}

	duration := time.Since(start)
	avgDuration := duration / numRequests

	t.Logf("Processed %d requests in %v (avg: %v per request)", numRequests, duration, avgDuration)

	// Assert reasonable performance (each request should be under 100ms in test environment)
	assert.Less(t, avgDuration, 100*time.Millisecond, "Average request duration should be under 100ms")

	mockService.AssertExpectations(t)
}

// Edge case tests
func TestPaymentHandler_CreatePayment_MalformedJSON(t *testing.T) {
	mockService := new(MockPaymentService)
	router := setupTestRouter(mockService)

	// Prepare malformed JSON request
	req := httptest.NewRequest(http.MethodPost, "/api/v1/payments", bytes.NewBuffer([]byte("{invalid json}")))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var errorResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	require.NoError(t, err)

	assert.Contains(t, errorResponse["error"], "invalid request body")
}

func TestPaymentHandler_GetPayment_InvalidPaymentID(t *testing.T) {
	mockService := new(MockPaymentService)
	router := setupTestRouter(mockService)

	// Test with empty payment ID
	req := httptest.NewRequest(http.MethodGet, "/api/v1/payments/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code) // Route not found
}