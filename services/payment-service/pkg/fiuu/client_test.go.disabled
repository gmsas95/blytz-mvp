package fiuu

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestNewClient(t *testing.T) {
	logger := zaptest.NewLogger(t)

	tests := []struct {
		name       string
		merchantID string
		verifyKey  string
		sandbox    bool
		wantURL    string
	}{
		{
			name:       "production client",
			merchantID: "test_merchant",
			verifyKey:  "test_key",
			sandbox:    false,
			wantURL:    "https://pay.fiuu.com",
		},
		{
			name:       "sandbox client",
			merchantID: "test_merchant",
			verifyKey:  "test_key",
			sandbox:    true,
			wantURL:    "https://sb-pay.fiuu.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.merchantID, tt.verifyKey, tt.sandbox, logger)

			assert.Equal(t, tt.merchantID, client.GetMerchantID())
			assert.Equal(t, tt.verifyKey, client.verifyKey)
			assert.Equal(t, tt.sandbox, client.IsSandbox())
			assert.Equal(t, tt.wantURL, client.GetBaseURL())
			assert.NotNil(t, client.httpClient)
			assert.Equal(t, 30*time.Second, client.httpClient.Timeout)
		})
	}
}

func TestClient_CreatePayment_Success(t *testing.T) {
	logger := zaptest.NewLogger(t)

	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/RMS/API/payment/PaymentRequest", r.URL.Path)
		assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))

		// Mock successful response
		response := PaymentResponse{
			TransactionID:  "TX123456",
			OrderID:        "ORD123",
			Amount:         100.50,
			Currency:       "MYR",
			PaymentStatus:  "1",
			ErrorCode:      "",
			ErrorDescription: "",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient("test_merchant", "test_key", true, logger)
	client.baseURL = server.URL // Override with test server

	req := PaymentRequest{
		Channel:     ChannelFPX,
		Amount:      100.50,
		OrderID:     "ORD123",
		BillName:    "John Doe",
		BillEmail:   "john@example.com",
		BillMobile:  "0123456789",
		BillDesc:    "Test Payment",
		Currency:    CurrencyMYR,
		LangCode:    "en",
		ReturnURL:   "https://example.com/return",
		NotifyURL:   "https://example.com/notify",
	}

	ctx := context.Background()
	resp, err := client.CreatePayment(ctx, req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "TX123456", resp.TransactionID)
	assert.Equal(t, "ORD123", resp.OrderID)
	assert.Equal(t, 100.50, resp.Amount)
	assert.Equal(t, "MYR", resp.Currency)
	assert.Equal(t, "1", resp.PaymentStatus)
	assert.Equal(t, "", resp.ErrorCode)
}

func TestClient_CreatePayment_ValidationError(t *testing.T) {
	logger := zaptest.NewLogger(t)
	client := NewClient("test_merchant", "test_key", true, logger)

	req := PaymentRequest{
		// Missing required fields
		Amount: 100.50,
	}

	ctx := context.Background()
	resp, err := client.CreatePayment(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "invalid payment request")
}

func TestClient_CreatePayment_APIError(t *testing.T) {
	logger := zaptest.NewLogger(t)

	// Mock server with error response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := PaymentResponse{
			ErrorCode:        "INVALID_PARAMETER",
			ErrorDescription: "Invalid merchant ID",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient("invalid_merchant", "test_key", true, logger)
	client.baseURL = server.URL

	req := PaymentRequest{
		Channel:     ChannelFPX,
		Amount:      100.50,
		OrderID:     "ORD123",
		BillName:    "John Doe",
		BillEmail:   "john@example.com",
		BillMobile:  "0123456789",
		BillDesc:    "Test Payment",
		Currency:    CurrencyMYR,
		LangCode:    "en",
	}

	ctx := context.Background()
	resp, err := client.CreatePayment(ctx, req)

	assert.Error(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "INVALID_PARAMETER", resp.ErrorCode)
	assert.Contains(t, err.Error(), "INVALID_PARAMETER - Invalid merchant ID")
}

func TestClient_CreatePayment_NetworkError(t *testing.T) {
	logger := zaptest.NewLogger(t)
	client := NewClient("test_merchant", "test_key", true, logger)

	// Use invalid URL to simulate network error
	client.baseURL = "http://invalid-url-that-does-not-exist.com"

	req := PaymentRequest{
		Channel:     ChannelFPX,
		Amount:      100.50,
		OrderID:     "ORD123",
		BillName:    "John Doe",
		BillEmail:   "john@example.com",
		BillMobile:  "0123456789",
		BillDesc:    "Test Payment",
		Currency:    CurrencyMYR,
		LangCode:    "en",
	}

	ctx := context.Background()
	resp, err := client.CreatePayment(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "failed to create payment")
}

func TestClient_CreatePayment_ContextTimeout(t *testing.T) {
	logger := zaptest.NewLogger(t)

	// Mock server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond) // Delay longer than context timeout
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(PaymentResponse{})
	}))
	defer server.Close()

	client := NewClient("test_merchant", "test_key", true, logger)
	client.baseURL = server.URL

	req := PaymentRequest{
		Channel:     ChannelFPX,
		Amount:      100.50,
		OrderID:     "ORD123",
		BillName:    "John Doe",
		BillEmail:   "john@example.com",
		BillMobile:  "0123456789",
		BillDesc:    "Test Payment",
		Currency:    CurrencyMYR,
		LangCode:    "en",
	}

	// Context with very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	resp, err := client.CreatePayment(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "context deadline exceeded")
}

func TestClient_CreateRefund_Success(t *testing.T) {
	logger := zaptest.NewLogger(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/RMS/API/payment/RefundRequest", r.URL.Path)

		response := RefundResponse{
			RefundID:        "REF123",
			OrderID:         "ORD123",
			Amount:          50.00,
			Currency:        "MYR",
			RefundStatus:    "1",
			ErrorCode:       "",
			ErrorDescription: "",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient("test_merchant", "test_key", true, logger)
	client.baseURL = server.URL

	req := RefundRequest{
		OrderID:      "ORD123",
		Amount:       50.00,
		RefundID:     "REF123",
		RefundReason: "Customer requested refund",
	}

	ctx := context.Background()
	resp, err := client.CreateRefund(ctx, req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "REF123", resp.RefundID)
	assert.Equal(t, "ORD123", resp.OrderID)
	assert.Equal(t, 50.00, resp.Amount)
	assert.Equal(t, "MYR", resp.Currency)
	assert.Equal(t, "1", resp.RefundStatus)
}

func TestClient_GetPaymentStatus_Success(t *testing.T) {
	logger := zaptest.NewLogger(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/RMS/API/payment/QueryPaymentStatus", r.URL.Path)

		response := PaymentResponse{
			TransactionID:  "TX123456",
			OrderID:        "ORD123",
			Amount:         100.50,
			Currency:       "MYR",
			PaymentStatus:  "1", // Success
			ErrorCode:      "",
			ErrorDescription: "",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient("test_merchant", "test_key", true, logger)
	client.baseURL = server.URL

	ctx := context.Background()
	resp, err := client.GetPaymentStatus(ctx, "ORD123")

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "TX123456", resp.TransactionID)
	assert.Equal(t, "ORD123", resp.OrderID)
	assert.Equal(t, "1", resp.PaymentStatus)
}

func TestClient_GetSeamlessConfig(t *testing.T) {
	logger := zaptest.NewLogger(t)
	client := NewClient("test_merchant", "test_key", true, logger)

	config := client.GetSeamlessConfig(
		"ORD123",
		100.50,
		"John Doe",
		"john@example.com",
		"0123456789",
		"Test Payment",
		ChannelFPX,
	)

	assert.Equal(t, "test_merchant", config["mpsmerchantid"])
	assert.Equal(t, "FPX", config["mpschannel"])
	assert.Equal(t, "100.50", config["mpsamount"])
	assert.Equal(t, "ORD123", config["mpsorderid"])
	assert.Equal(t, "John Doe", config["mpsbill_name"])
	assert.Equal(t, "john@example.com", config["mpsbill_email"])
	assert.Equal(t, "0123456789", config["mpsbill_mobile"])
	assert.Equal(t, "Test Payment", config["mpsbill_desc"])
	assert.Equal(t, "MYR", config["mpscurrency"])
	assert.Equal(t, "en", config["mpslangcode"])
	assert.Equal(t, true, config["sandbox"])
	assert.NotEmpty(t, config["vcode"])
}

// Edge case tests
func TestClient_CreatePayment_InvalidResponse(t *testing.T) {
	logger := zaptest.NewLogger(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json")) // Invalid JSON response
	}))
	defer server.Close()

	client := NewClient("test_merchant", "test_key", true, logger)
	client.baseURL = server.URL

	req := PaymentRequest{
		Channel:     ChannelFPX,
		Amount:      100.50,
		OrderID:     "ORD123",
		BillName:    "John Doe",
		BillEmail:   "john@example.com",
		BillMobile:  "0123456789",
		BillDesc:    "Test Payment",
		Currency:    CurrencyMYR,
		LangCode:    "en",
	}

	ctx := context.Background()
	resp, err := client.CreatePayment(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "failed to parse response")
}

func TestClient_CreatePayment_EmptyResponse(t *testing.T) {
	logger := zaptest.NewLogger(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("")) // Empty response
	}))
	defer server.Close()

	client := NewClient("test_merchant", "test_key", true, logger)
	client.baseURL = server.URL

	req := PaymentRequest{
		Channel:     ChannelFPX,
		Amount:      100.50,
		OrderID:     "ORD123",
		BillName:    "John Doe",
		BillEmail:   "john@example.com",
		BillMobile:  "0123456789",
		BillDesc:    "Test Payment",
		Currency:    CurrencyMYR,
		LangCode:    "en",
	}

	ctx := context.Background()
	resp, err := client.CreatePayment(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "failed to parse response")
}

// Performance test
func TestClient_CreatePayment_Performance(t *testing.T) {
	logger := zaptest.NewLogger(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := PaymentResponse{
			TransactionID:  "TX123456",
			OrderID:        "ORD123",
			Amount:         100.50,
			Currency:       "MYR",
			PaymentStatus:  "1",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient("test_merchant", "test_key", true, logger)
	client.baseURL = server.URL

	req := PaymentRequest{
		Channel:     ChannelFPX,
		Amount:      100.50,
		OrderID:     "ORD123",
		BillName:    "John Doe",
		BillEmail:   "john@example.com",
		BillMobile:  "0123456789",
		BillDesc:    "Test Payment",
		Currency:    CurrencyMYR,
		LangCode:    "en",
	}

	// Measure performance of multiple concurrent requests
	const numRequests = 10
	start := time.Now()

	for i := 0; i < numRequests; i++ {
		req.OrderID = fmt.Sprintf("ORD%d", i)
		ctx := context.Background()
		resp, err := client.CreatePayment(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, resp)
	}

	duration := time.Since(start)
	avgDuration := duration / numRequests

	t.Logf("Processed %d requests in %v (avg: %v per request)", numRequests, duration, avgDuration)

	// Assert reasonable performance (each request should be under 1 second in test environment)
	assert.Less(t, avgDuration, 1*time.Second, "Average request duration should be under 1 second")
}