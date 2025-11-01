package fiuu

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestCircuitBreaker_NewCircuitBreaker(t *testing.T) {
	logger := zaptest.NewLogger(t)
	cb := NewCircuitBreaker(5, 30*time.Second, logger)

	assert.Equal(t, 5, cb.maxFailures)
	assert.Equal(t, 30*time.Second, cb.resetTimeout)
	assert.Equal(t, CircuitClosed, cb.state)
	assert.Equal(t, 0, cb.failures)
}

func TestCircuitBreaker_ClosedState(t *testing.T) {
	logger := zaptest.NewLogger(t)
	cb := NewCircuitBreaker(2, 1*time.Second, logger)

	// Initially closed
	assert.True(t, cb.CanExecute())
	assert.Equal(t, CircuitClosed, cb.GetState())

	// Success should not change state
	cb.OnSuccess()
	assert.Equal(t, CircuitClosed, cb.state)
	assert.Equal(t, 0, cb.failures)
}

func TestCircuitBreaker_OpenState(t *testing.T) {
	logger := zaptest.NewLogger(t)
	cb := NewCircuitBreaker(2, 1*time.Second, logger)

	// Record failures to trigger open state
	cb.OnFailure()
	assert.Equal(t, CircuitClosed, cb.state)
	assert.Equal(t, 1, cb.failures)

	cb.OnFailure()
	assert.Equal(t, CircuitOpen, cb.state)
	assert.Equal(t, 2, cb.failures)

	// Should not allow execution when open
	assert.False(t, cb.CanExecute())
}

func TestCircuitBreaker_HalfOpenState(t *testing.T) {
	logger := zaptest.NewLogger(t)
	cb := NewCircuitBreaker(2, 100*time.Millisecond, logger)

	// Trigger open state
	cb.OnFailure()
	cb.OnFailure()
	assert.Equal(t, CircuitOpen, cb.state)

	// Wait for reset timeout
	time.Sleep(150 * time.Millisecond)

	// Should now be half-open and allow execution
	assert.True(t, cb.CanExecute())
	assert.Equal(t, CircuitHalfOpen, cb.state)

	// Success should close the circuit
	cb.OnSuccess()
	assert.Equal(t, CircuitClosed, cb.state)
	assert.Equal(t, 0, cb.failures)
}

func TestCircuitBreaker_HalfOpenFailure(t *testing.T) {
	logger := zaptest.NewLogger(t)
	cb := NewCircuitBreaker(2, 100*time.Millisecond, logger)

	// Trigger open state
	cb.OnFailure()
	cb.OnFailure()
	assert.Equal(t, CircuitOpen, cb.state)

	// Wait for reset timeout
	time.Sleep(150 * time.Millisecond)

	// Should be half-open
	assert.True(t, cb.CanExecute())
	assert.Equal(t, CircuitHalfOpen, cb.state)

	// Failure should open the circuit again
	cb.OnFailure()
	assert.Equal(t, CircuitOpen, cb.state)
	assert.Equal(t, 3, cb.failures) // Should not reset failures on half-open failure
}

func TestIsRetryable_NetworkErrors(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		retryable bool
	}{
		{
			name:      "connection refused",
			err:       errors.New("connection refused"),
			retryable: true,
		},
		{
			name:      "timeout error",
			err:       errors.New("operation timed out"),
			retryable: true,
		},
		{
			name:      "network unreachable",
			err:       errors.New("network is unreachable"),
			retryable: true,
		},
		{
			name:      "service unavailable",
			err:       errors.New("service temporarily unavailable"),
			retryable: true,
		},
		{
			name:      "invalid parameter",
			err:       errors.New("invalid parameter: amount"),
			retryable: false,
		},
		{
			name:      "authentication failed",
			err:       errors.New("authentication failed"),
			retryable: false,
		},
		{
			name:      "not found",
			err:       errors.New("order not found"),
			retryable: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retryableErr := IsRetryable(tt.err)
			assert.Equal(t, tt.retryable, retryableErr.Retryable)
		})
	}
}

func TestCalculateDelay_LinearBackoff(t *testing.T) {
	config := RetryConfig{
		BaseDelay: 1 * time.Second,
		MaxDelay:  10 * time.Second,
		Backoff:   BackoffLinear,
	}

	assert.Equal(t, 1*time.Second, calculateDelay(1, config))
	assert.Equal(t, 2*time.Second, calculateDelay(2, config))
	assert.Equal(t, 3*time.Second, calculateDelay(3, config))
	assert.Equal(t, 10*time.Second, calculateDelay(15, config)) // Capped at max
}

func TestCalculateDelay_ExponentialBackoff(t *testing.T) {
	config := RetryConfig{
		BaseDelay: 1 * time.Second,
		MaxDelay:  10 * time.Second,
		Backoff:   BackoffExponential,
	}

	assert.Equal(t, 1*time.Second, calculateDelay(1, config))
	assert.Equal(t, 2*time.Second, calculateDelay(2, config))
	assert.Equal(t, 4*time.Second, calculateDelay(3, config))
	assert.Equal(t, 8*time.Second, calculateDelay(4, config))
	assert.Equal(t, 10*time.Second, calculateDelay(5, config)) // Capped at max
}

func TestCalculateDelay_ExponentialWithJitter(t *testing.T) {
	config := RetryConfig{
		BaseDelay: 1 * time.Second,
		MaxDelay:  10 * time.Second,
		Backoff:   BackoffExponentialWithJitter,
	}

	delay := calculateDelay(3, config)
	// Should be around 4 seconds with jitter (±25%)
	assert.True(t, delay >= 3*time.Second)
	assert.True(t, delay <= 5*time.Second)
}

func TestWithRetry_Success(t *testing.T) {
	logger := zaptest.NewLogger(t)
	config := RetryConfig{
		MaxAttempts: 3,
		BaseDelay:   10 * time.Millisecond,
		MaxDelay:    100 * time.Millisecond,
		Backoff:     BackoffLinear,
	}

	attempts := 0
	fn := func(ctx context.Context) error {
		attempts++
		if attempts < 2 {
			return errors.New("temporary failure")
		}
		return nil
	}

	ctx := context.Background()
	err := WithRetry(ctx, config, fn, logger)

	assert.NoError(t, err)
	assert.Equal(t, 2, attempts)
}

func TestWithRetry_AllAttemptsFail(t *testing.T) {
	logger := zaptest.NewLogger(t)
	config := RetryConfig{
		MaxAttempts: 3,
		BaseDelay:   10 * time.Millisecond,
		MaxDelay:    100 * time.Millisecond,
		Backoff:     BackoffLinear,
	}

	attempts := 0
	fn := func(ctx context.Context) error {
		attempts++
		return errors.New("persistent failure")
	}

	ctx := context.Background()
	err := WithRetry(ctx, config, fn, logger)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "persistent failure")
	assert.Equal(t, 3, attempts)
}

func TestWithRetry_NonRetryableError(t *testing.T) {
	logger := zaptest.NewLogger(t)
	config := RetryConfig{
		MaxAttempts: 3,
		BaseDelay:   10 * time.Millisecond,
		MaxDelay:    100 * time.Millisecond,
		Backoff:     BackoffLinear,
	}

	attempts := 0
	fn := func(ctx context.Context) error {
		attempts++
		return errors.New("invalid parameter")
	}

	ctx := context.Background()
	err := WithRetry(ctx, config, fn, logger)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid parameter")
	assert.Equal(t, 1, attempts) // Should not retry non-retryable error
}

func TestWithRetry_ContextCancellation(t *testing.T) {
	logger := zaptest.NewLogger(t)
	config := RetryConfig{
		MaxAttempts: 3,
		BaseDelay:   100 * time.Millisecond,
		MaxDelay:    1 * time.Second,
		Backoff:     BackoffLinear,
	}

	attempts := 0
	fn := func(ctx context.Context) error {
		attempts++
		return errors.New("temporary failure")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := WithRetry(ctx, config, fn, logger)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context cancelled")
}

func TestResilientClient_CreatePaymentWithRetry_Success(t *testing.T) {
	logger := zaptest.NewLogger(t)

	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := PaymentResponse{
			TransactionID: "TX123456",
			OrderID:       "ORD123",
			Amount:        100.50,
			Currency:      "MYR",
			PaymentStatus: "1",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewResilientClient("test_merchant", "test_key", true, logger)
	client.Client.baseURL = server.URL

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
	resp, err := client.CreatePaymentWithRetry(ctx, req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "TX123456", resp.TransactionID)
	assert.Equal(t, CircuitClosed, client.GetCircuitBreakerState())
}

func TestResilientClient_CreatePaymentWithRetry_RetrySuccess(t *testing.T) {
	logger := zaptest.NewLogger(t)

	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		response := PaymentResponse{
			TransactionID: "TX123456",
			OrderID:       "ORD123",
			Amount:        100.50,
			Currency:      "MYR",
			PaymentStatus: "1",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewResilientClient("test_merchant", "test_key", true, logger)
	client.WithRetryConfig(RetryConfig{
		MaxAttempts: 5,
		BaseDelay:   10 * time.Millisecond,
		MaxDelay:    100 * time.Millisecond,
		Backoff:     BackoffLinear,
	})
	client.Client.baseURL = server.URL

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
	resp, err := client.CreatePaymentWithRetry(ctx, req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "TX123456", resp.TransactionID)
	assert.Equal(t, 3, attempts)
	assert.Equal(t, CircuitClosed, client.GetCircuitBreakerState())
}

func TestResilientClient_CreatePaymentWithRetry_CircuitBreakerOpen(t *testing.T) {
	logger := zaptest.NewLogger(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewResilientClient("test_merchant", "test_key", true, logger)
	client.WithCircuitBreaker(2, 1*time.Second) // Low threshold for testing
	client.Client.baseURL = server.URL

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

	// First two failures should trigger circuit breaker
	_, err1 := client.CreatePaymentWithRetry(ctx, req)
	assert.Error(t, err1)

	_, err2 := client.CreatePaymentWithRetry(ctx, req)
	assert.Error(t, err2)

	// Third attempt should be blocked by circuit breaker
	_, err3 := client.CreatePaymentWithRetry(ctx, req)
	assert.Error(t, err3)
	assert.Contains(t, err3.Error(), "circuit breaker is OPEN")
	assert.Equal(t, CircuitOpen, client.GetCircuitBreakerState())
}

func TestResilientClient_CreatePaymentWithRetry_NonRetryableError(t *testing.T) {
	logger := zaptest.NewLogger(t)

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

	client := NewResilientClient("invalid_merchant", "test_key", true, logger)
	client.Client.baseURL = server.URL

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
	resp, err := client.CreatePaymentWithRetry(ctx, req)

	assert.Error(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "INVALID_PARAMETER", resp.ErrorCode)
	assert.Equal(t, CircuitClosed, client.GetCircuitBreakerState()) // Should not trigger circuit breaker for non-retryable errors
}

// Performance test for resilient client
func TestResilientClient_Performance(t *testing.T) {
	logger := zaptest.NewLogger(t)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := PaymentResponse{
			TransactionID: "TX123456",
			OrderID:       "ORD123",
			Amount:        100.50,
			Currency:      "MYR",
			PaymentStatus: "1",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewResilientClient("test_merchant", "test_key", true, logger)
	client.Client.baseURL = server.URL

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
	const numRequests := 20
	start := time.Now()

	for i := 0; i < numRequests; i++ {
		req.OrderID = fmt.Sprintf("ORD%d", i)
		ctx := context.Background()
		resp, err := client.CreatePaymentWithRetry(ctx, req)

		require.NoError(t, err)
		assert.NotNil(t, resp)
	}

	duration := time.Since(start)
	avgDuration := duration / numRequests

	t.Logf("Processed %d resilient requests in %v (avg: %v per request)", numRequests, duration, avgDuration)

	// Assert reasonable performance (resilient client should still be fast)
	assert.Less(t, avgDuration, 50*time.Millisecond, "Average resilient request duration should be under 50ms")
}