package fiuu

import (
	"context"
	"fmt"
	"math"
	"net"
	"strings"
	"time"

	"go.uber.org/zap"
)

// RetryConfig defines the retry behavior
type RetryConfig struct {
	MaxAttempts int           // Maximum number of retry attempts
	BaseDelay   time.Duration // Base delay between retries
	MaxDelay    time.Duration // Maximum delay between retries
	Backoff     BackoffType   // Backoff strategy
}

// BackoffType defines the backoff strategy
type BackoffType int

const (
	BackoffLinear BackoffType = iota
	BackoffExponential
	BackoffExponentialWithJitter
)

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	maxFailures   int
	resetTimeout  time.Duration
	failures      int
	lastFailTime  time.Time
	state         CircuitState
	logger        *zap.Logger
}

// CircuitState represents the circuit breaker state
type CircuitState int

const (
	CircuitClosed CircuitState = iota
	CircuitOpen
	CircuitHalfOpen
)

// String returns the string representation of circuit state
func (s CircuitState) String() string {
	switch s {
	case CircuitClosed:
		return "CLOSED"
	case CircuitOpen:
		return "OPEN"
	case CircuitHalfOpen:
		return "HALF_OPEN"
	default:
		return "UNKNOWN"
	}
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(maxFailures int, resetTimeout time.Duration, logger *zap.Logger) *CircuitBreaker {
	return &CircuitBreaker{
		maxFailures:  maxFailures,
		resetTimeout: resetTimeout,
		state:        CircuitClosed,
		logger:       logger,
	}
}

// CanExecute checks if the circuit breaker allows execution
func (cb *CircuitBreaker) CanExecute() bool {
	switch cb.state {
	case CircuitClosed:
		return true
	case CircuitOpen:
		if time.Since(cb.lastFailTime) > cb.resetTimeout {
			cb.state = CircuitHalfOpen
			cb.logger.Info("Circuit breaker transitioning to HALF_OPEN")
			return true
		}
		return false
	case CircuitHalfOpen:
		return true
	default:
		return false
	}
}

// OnSuccess records a successful execution
func (cb *CircuitBreaker) OnSuccess() {
	cb.failures = 0
	if cb.state == CircuitHalfOpen {
		cb.state = CircuitClosed
		cb.logger.Info("Circuit breaker transitioning to CLOSED")
	}
}

// OnFailure records a failed execution
func (cb *CircuitBreaker) OnFailure() {
	cb.failures++
	cb.lastFailTime = time.Now()

	if cb.failures >= cb.maxFailures {
		cb.state = CircuitOpen
		cb.logger.Warn("Circuit breaker transitioning to OPEN",
			zap.Int("failures", cb.failures),
			zap.Duration("reset_timeout", cb.resetTimeout))
	}
}

// GetState returns the current circuit breaker state
func (cb *CircuitBreaker) GetState() CircuitState {
	return cb.state
}

// RetryableError represents an error that can be retried
type RetryableError struct {
	Err        error
	Retryable  bool
	RetryAfter time.Duration
}

func (e *RetryableError) Error() string {
	return e.Err.Error()
}

func (e *RetryableError) Unwrap() error {
	return e.Err
}

// IsRetryable determines if an error is retryable
func IsRetryable(err error) *RetryableError {
	if err == nil {
		return &RetryableError{Err: err, Retryable: false}
	}

	errStr := strings.ToLower(err.Error())

	// Network errors that are typically retryable
	if isNetworkError(err) {
		return &RetryableError{
			Err:       err,
			Retryable: true,
		}
	}

	// Fiuu specific errors that are retryable
	if isFiuuRetryableError(errStr) {
		return &RetryableError{
			Err:       err,
			Retryable: true,
		}
	}

	// Errors that should not be retried
	if isNonRetryableError(errStr) {
		return &RetryableError{
			Err:       err,
			Retryable: false,
		}
	}

	// Default to retryable for unknown errors
	return &RetryableError{
		Err:       err,
		Retryable: true,
	}
}

// isNetworkError checks if the error is a network-related error
func isNetworkError(err error) bool {
	if netErr, ok := err.(net.Error); ok {
		// Timeout errors are retryable
		if netErr.Timeout() {
			return true
		}
		// Temporary network errors are retryable
		if netErr.Temporary() {
			return true
		}
	}

	// Check for common network error strings
	errStr := strings.ToLower(err.Error())
	networkErrors := []string{
		"connection refused",
		"connection reset",
		"connection timed out",
		"timeout",
		"network is unreachable",
		"no such host",
		"temporary failure",
		"service unavailable",
		"bad gateway",
		"gateway timeout",
	}

	for _, networkErr := range networkErrors {
		if strings.Contains(errStr, networkErr) {
			return true
		}
	}

	return false
}

// isFiuuRetryableError checks if the error is a Fiuu-specific retryable error
func isFiuuRetryableError(errStr string) bool {
	retryableErrors := []string{
		"system busy",
		"temporarily unavailable",
		"timeout",
		"processing",
		"pending",
		"try again later",
		"service busy",
		"rate limit",
	}

	for _, retryableErr := range retryableErrors {
		if strings.Contains(errStr, retryableErr) {
			return true
		}
	}

	return false
}

// isNonRetryableError checks if the error should not be retried
func isNonRetryableError(errStr string) bool {
	nonRetryableErrors := []string{
		"invalid parameter",
		"invalid merchant",
		"invalid signature",
		"authentication failed",
		"authorization failed",
		"access denied",
		"forbidden",
		"not found",
		"invalid order",
		"duplicate transaction",
		"invalid amount",
		"invalid currency",
		"invalid channel",
	}

	for _, nonRetryableErr := range nonRetryableErrors {
		if strings.Contains(errStr, nonRetryableErr) {
			return true
		}
	}

	return false
}

// calculateDelay calculates the delay for a given retry attempt
func calculateDelay(attempt int, config RetryConfig) time.Duration {
	var delay time.Duration

	switch config.Backoff {
	case BackoffLinear:
		delay = config.BaseDelay * time.Duration(attempt)
	case BackoffExponential:
		delay = config.BaseDelay * time.Duration(math.Pow(2, float64(attempt-1)))
	case BackoffExponentialWithJitter:
		exponentialDelay := config.BaseDelay * time.Duration(math.Pow(2, float64(attempt-1)))
		// Add jitter: ±25% of the delay
		jitter := time.Duration(float64(exponentialDelay) * 0.25 * (2.0*float64(time.Now().UnixNano()%1000)/1000.0 - 1.0))
		delay = exponentialDelay + jitter
	default:
		delay = config.BaseDelay
	}

	// Ensure delay doesn't exceed maximum
	if delay > config.MaxDelay {
		delay = config.MaxDelay
	}

	return delay
}

// RetryableFunction represents a function that can be retried
type RetryableFunction func(ctx context.Context) error

// WithRetry executes a function with retry logic
func WithRetry(ctx context.Context, config RetryConfig, fn RetryableFunction, logger *zap.Logger) error {
	var lastErr error

	for attempt := 1; attempt <= config.MaxAttempts; attempt++ {
		if ctx.Err() != nil {
			return fmt.Errorf("context cancelled: %w", ctx.Err())
		}

		err := fn(ctx)
		if err == nil {
			logger.Debug("Operation succeeded",
				zap.Int("attempt", attempt))
			return nil
		}

		lastErr = err
		retryableErr := IsRetryable(err)

		logger.Warn("Operation failed",
			zap.Int("attempt", attempt),
			zap.Int("max_attempts", config.MaxAttempts),
			zap.Error(err),
			zap.Bool("retryable", retryableErr.Retryable))

		// If error is not retryable or this is the last attempt, return
		if !retryableErr.Retryable || attempt == config.MaxAttempts {
			if !retryableErr.Retryable {
				logger.Warn("Non-retryable error, not attempting retry",
					zap.Error(err))
			} else {
				logger.Error("All retry attempts exhausted",
					zap.Int("attempts", attempt),
					zap.Error(err))
			}
			return lastErr
		}

		// Calculate delay and wait
		delay := calculateDelay(attempt, config)
		if retryableErr.RetryAfter > 0 && retryableErr.RetryAfter < delay {
			delay = retryableErr.RetryAfter
		}

		logger.Info("Retrying operation",
			zap.Int("attempt", attempt),
			zap.Int("next_attempt", attempt+1),
			zap.Duration("delay", delay))

		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled during retry wait: %w", ctx.Err())
		case <-time.After(delay):
			continue
		}
	}

	return lastErr
}

// ResilientClient wraps the Fiuu client with retry and circuit breaker functionality
type ResilientClient struct {
	*Client
	circuitBreaker *CircuitBreaker
	retryConfig    RetryConfig
	logger         *zap.Logger
}

// NewResilientClient creates a new resilient Fiuu client
func NewResilientClient(merchantID, verifyKey string, sandbox bool, logger *zap.Logger) *ResilientClient {
	client := NewClient(merchantID, verifyKey, sandbox, logger)

	// Default retry configuration
	retryConfig := RetryConfig{
		MaxAttempts: 3,
		BaseDelay:   1 * time.Second,
		MaxDelay:    30 * time.Second,
		Backoff:     BackoffExponentialWithJitter,
	}

	// Circuit breaker configuration
	circuitBreaker := NewCircuitBreaker(
		5,              // max failures
		60*time.Second, // reset timeout
		logger,
	)

	return &ResilientClient{
		Client:         client,
		circuitBreaker: circuitBreaker,
		retryConfig:    retryConfig,
		logger:         logger,
	}
}

// WithRetryConfig sets custom retry configuration
func (rc *ResilientClient) WithRetryConfig(config RetryConfig) *ResilientClient {
	rc.retryConfig = config
	return rc
}

// WithCircuitBreaker sets custom circuit breaker configuration
func (rc *ResilientClient) WithCircuitBreaker(maxFailures int, resetTimeout time.Duration) *ResilientClient {
	rc.circuitBreaker = NewCircuitBreaker(maxFailures, resetTimeout, rc.logger)
	return rc
}

// CreatePaymentWithRetry creates a payment with retry and circuit breaker
func (rc *ResilientClient) CreatePaymentWithRetry(ctx context.Context, req PaymentRequest) (*PaymentResponse, error) {
	if !rc.circuitBreaker.CanExecute() {
		return nil, fmt.Errorf("circuit breaker is %s", rc.circuitBreaker.GetState())
	}

	var result *PaymentResponse
	var err error

	retryErr := WithRetry(ctx, rc.retryConfig, func(ctx context.Context) error {
		result, err = rc.Client.CreatePayment(ctx, req)
		return err
	}, rc.logger)

	if retryErr != nil {
		rc.circuitBreaker.OnFailure()
		return nil, retryErr
	}

	rc.circuitBreaker.OnSuccess()
	return result, nil
}

// CreateRefundWithRetry creates a refund with retry and circuit breaker
func (rc *ResilientClient) CreateRefundWithRetry(ctx context.Context, req RefundRequest) (*RefundResponse, error) {
	if !rc.circuitBreaker.CanExecute() {
		return nil, fmt.Errorf("circuit breaker is %s", rc.circuitBreaker.GetState())
	}

	var result *RefundResponse
	var err error

	retryErr := WithRetry(ctx, rc.retryConfig, func(ctx context.Context) error {
		result, err = rc.Client.CreateRefund(ctx, req)
		return err
	}, rc.logger)

	if retryErr != nil {
		rc.circuitBreaker.OnFailure()
		return nil, retryErr
	}

	rc.circuitBreaker.OnSuccess()
	return result, nil
}

// GetPaymentStatusWithRetry gets payment status with retry and circuit breaker
func (rc *ResilientClient) GetPaymentStatusWithRetry(ctx context.Context, orderID string) (*PaymentResponse, error) {
	if !rc.circuitBreaker.CanExecute() {
		return nil, fmt.Errorf("circuit breaker is %s", rc.circuitBreaker.GetState())
	}

	var result *PaymentResponse
	var err error

	retryErr := WithRetry(ctx, rc.retryConfig, func(ctx context.Context) error {
		result, err = rc.Client.GetPaymentStatus(ctx, orderID)
		return err
	}, rc.logger)

	if retryErr != nil {
		rc.circuitBreaker.OnFailure()
		return nil, retryErr
	}

	rc.circuitBreaker.OnSuccess()
	return result, nil
}

// GetCircuitBreakerState returns the current circuit breaker state
func (rc *ResilientClient) GetCircuitBreakerState() CircuitState {
	return rc.circuitBreaker.GetState()
}

// GetRetryConfig returns the current retry configuration
func (rc *ResilientClient) GetRetryConfig() RetryConfig {
	return rc.retryConfig
}