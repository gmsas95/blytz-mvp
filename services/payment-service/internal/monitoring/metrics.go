package monitoring

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
)

var (
	// Payment metrics
	paymentTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "payments_total",
		Help: "Total number of payment attempts",
	}, []string{"method", "currency", "status"})

	paymentAmount = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "payment_amount",
		Help:    "Distribution of payment amounts",
		Buckets: []float64{1, 5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000, 25000, 50000},
	}, []string{"method", "currency"})

	paymentDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "payment_duration_seconds",
		Help:    "Time taken to process payments",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "status"})

	paymentSuccessRate = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "payment_success_rate",
		Help: "Payment success rate by method and currency",
	}, []string{"method", "currency"})

	// Refund metrics
	refundTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "refunds_total",
		Help: "Total number of refund attempts",
	}, []string{"method", "currency", "status"})

	refundAmount = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "refund_amount",
		Help:    "Distribution of refund amounts",
		Buckets: []float64{1, 5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000},
	}, []string{"method", "currency"})

	// Fiuu specific metrics
	fiuuRequestTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "fiuu_requests_total",
		Help: "Total number of Fiuu API requests",
	}, []string{"endpoint", "method", "status"})

	fiuuRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "fiuu_request_duration_seconds",
		Help:    "Time taken for Fiuu API requests",
		Buckets: prometheus.DefBuckets,
	}, []string{"endpoint", "method", "status"})

	fiuuCircuitBreakerState = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "fiuu_circuit_breaker_state",
		Help: "Current state of Fiuu circuit breaker (0=CLOSED, 1=OPEN, 2=HALF_OPEN)",
	}, []string{"merchant_id"})

	fiuuRetryAttempts = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "fiuu_retry_attempts",
		Help:    "Number of retry attempts for Fiuu requests",
		Buckets: []float64{0, 1, 2, 3, 4, 5},
	}, []string{"endpoint", "error_type"})

	// Webhook metrics
	webhookReceivedTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "webhooks_received_total",
		Help: "Total number of webhooks received",
	}, []string{"source", "event_type", "status"})

	webhookProcessingDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "webhook_processing_duration_seconds",
		Help:    "Time taken to process webhooks",
		Buckets: prometheus.DefBuckets,
	}, []string{"source", "event_type", "status"})

	webhookRetryQueue = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "webhook_retry_queue_size",
		Help: "Number of webhooks currently in retry queue",
	})

	// Database metrics
	dbConnections = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "db_connections",
		Help: "Number of database connections",
	}, []string{"state"})

	dbQueryDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "db_query_duration_seconds",
		Help:    "Time taken for database queries",
		Buckets: prometheus.DefBuckets,
	}, []string{"operation", "table"})

	// Business metrics
	revenueTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "revenue_total",
		Help: "Total revenue processed",
	}, []string{"method", "currency"})

	activeUsers = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "active_users",
		Help: "Number of active users making payments",
	})

	topPaymentMethods = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "top_payment_methods",
		Help: "Most used payment methods",
	}, []string{"method"})

	// System metrics
	httpRequestTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	}, []string{"method", "endpoint", "status_code"})

	httpRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Time taken to process HTTP requests",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "endpoint", "status_code"})

	// Error metrics
	errorTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "errors_total",
		Help: "Total number of errors",
	}, []string{"type", "component", "severity"})

	panicTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "panics_total",
		Help: "Total number of panics",
	})
)

// PaymentMetrics provides methods to record payment-related metrics
type PaymentMetrics struct {
	logger *zap.Logger
}

// NewPaymentMetrics creates a new PaymentMetrics instance
func NewPaymentMetrics(logger *zap.Logger) *PaymentMetrics {
	return &PaymentMetrics{
		logger: logger,
	}
}

// RecordPaymentAttempt records a payment attempt
func (pm *PaymentMetrics) RecordPaymentAttempt(method, currency, status string, amount float64, duration time.Duration) {
	labels := prometheus.Labels{
		"method":   method,
		"currency": currency,
		"status":   status,
	}

	paymentTotal.With(labels).Inc()
	paymentAmount.With(prometheus.Labels{
		"method":   method,
		"currency": currency,
	}).Observe(amount)

	paymentDuration.With(prometheus.Labels{
		"method": method,
		"status": status,
	}).Observe(duration.Seconds())

	pm.logger.Info("Payment attempt recorded",
		zap.String("method", method),
		zap.String("currency", currency),
		zap.String("status", status),
		zap.Float64("amount", amount),
		zap.Duration("duration", duration))
}

// RecordRefundAttempt records a refund attempt
func (pm *PaymentMetrics) RecordRefundAttempt(method, currency, status string, amount float64) {
	labels := prometheus.Labels{
		"method":   method,
		"currency": currency,
		"status":   status,
	}

	refundTotal.With(labels).Inc()
	refundAmount.With(prometheus.Labels{
		"method":   method,
		"currency": currency,
	}).Observe(amount)

	pm.logger.Info("Refund attempt recorded",
		zap.String("method", method),
		zap.String("currency", currency),
		zap.String("status", status),
		zap.Float64("amount", amount))
}

// RecordFiuuRequest records a Fiuu API request
func (pm *PaymentMetrics) RecordFiuuRequest(endpoint, method, status string, duration time.Duration) {
	labels := prometheus.Labels{
		"endpoint": endpoint,
		"method":   method,
		"status":   status,
	}

	fiuuRequestTotal.With(labels).Inc()
	fiuuRequestDuration.With(labels).Observe(duration.Seconds())

	pm.logger.Debug("Fiuu request recorded",
		zap.String("endpoint", endpoint),
		zap.String("method", method),
		zap.String("status", status),
		zap.Duration("duration", duration))
}

// RecordFiuuCircuitBreakerState records the circuit breaker state
func (pm *PaymentMetrics) RecordFiuuCircuitBreakerState(merchantID string, state int) {
	fiuuCircuitBreakerState.With(prometheus.Labels{
		"merchant_id": merchantID,
	}).Set(float64(state))

	stateStr := "UNKNOWN"
	switch state {
	case 0:
		stateStr = "CLOSED"
	case 1:
		stateStr = "OPEN"
	case 2:
		stateStr = "HALF_OPEN"
	}

	pm.logger.Info("Fiuu circuit breaker state changed",
		zap.String("merchant_id", merchantID),
		zap.String("state", stateStr))
}

// RecordFiuuRetry records retry attempts for Fiuu requests
func (pm *PaymentMetrics) RecordFiuuRetry(endpoint, errorType string, attempts int) {
	fiuuRetryAttempts.With(prometheus.Labels{
		"endpoint":   endpoint,
		"error_type": errorType,
	}).Observe(float64(attempts))

	pm.logger.Info("Fiuu retry recorded",
		zap.String("endpoint", endpoint),
		zap.String("error_type", errorType),
		zap.Int("attempts", attempts))
}

// RecordWebhookReceived records a received webhook
func (pm *PaymentMetrics) RecordWebhookReceived(source, eventType, status string, duration time.Duration) {
	labels := prometheus.Labels{
		"source":     source,
		"event_type": eventType,
		"status":     status,
	}

	webhookReceivedTotal.With(labels).Inc()
	webhookProcessingDuration.With(labels).Observe(duration.Seconds())

	pm.logger.Info("Webhook received and processed",
		zap.String("source", source),
		zap.String("event_type", eventType),
		zap.String("status", status),
		zap.Duration("duration", duration))
}

// UpdateWebhookRetryQueue updates the webhook retry queue size
func (pm *PaymentMetrics) UpdateWebhookRetryQueue(size int) {
	webhookRetryQueue.Set(float64(size))
	pm.logger.Debug("Webhook retry queue size updated", zap.Int("size", size))
}

// RecordDatabaseQuery records a database query
func (pm *PaymentMetrics) RecordDatabaseQuery(operation, table string, duration time.Duration) {
	dbQueryDuration.With(prometheus.Labels{
		"operation": operation,
		"table":     table,
	}).Observe(duration.Seconds())

	pm.logger.Debug("Database query recorded",
		zap.String("operation", operation),
		zap.String("table", table),
		zap.Duration("duration", duration))
}

// RecordRevenue records revenue
func (pm *PaymentMetrics) RecordRevenue(method, currency string, amount float64) {
	revenueTotal.With(prometheus.Labels{
		"method":   method,
		"currency": currency,
	}).Add(amount)

	pm.logger.Info("Revenue recorded",
		zap.String("method", method),
		zap.String("currency", currency),
		zap.Float64("amount", amount))
}

// UpdateActiveUsers updates the active users count
func (pm *PaymentMetrics) UpdateActiveUsers(count int) {
	activeUsers.Set(float64(count))
	pm.logger.Debug("Active users updated", zap.Int("count", count))
}

// UpdateTopPaymentMethods updates the top payment methods
func (pm *PaymentMetrics) UpdateTopPaymentMethods(methods map[string]int) {
	// Reset existing metrics
	topPaymentMethods.Reset()

	for method, count := range methods {
		topPaymentMethods.With(prometheus.Labels{
			"method": method,
		}).Set(float64(count))
	}

	pm.logger.Debug("Top payment methods updated", zap.Any("methods", methods))
}

// RecordHTTPRequest records an HTTP request
func (pm *PaymentMetrics) RecordHTTPRequest(method, endpoint, statusCode string, duration time.Duration) {
	labels := prometheus.Labels{
		"method":      method,
		"endpoint":    endpoint,
		"status_code": statusCode,
	}

	httpRequestTotal.With(labels).Inc()
	httpRequestDuration.With(labels).Observe(duration.Seconds())

	pm.logger.Debug("HTTP request recorded",
		zap.String("method", method),
		zap.String("endpoint", endpoint),
		zap.String("status_code", statusCode),
		zap.Duration("duration", duration))
}

// RecordError records an error
func (pm *PaymentMetrics) RecordError(errorType, component, severity string) {
	labels := prometheus.Labels{
		"type":      errorType,
		"component": component,
		"severity":  severity,
	}

	errorTotal.With(labels).Inc()

	pm.logger.Error("Error recorded",
		zap.String("type", errorType),
		zap.String("component", component),
		zap.String("severity", severity))
}

// RecordPanic records a panic
func (pm *PaymentMetrics) RecordPanic() {
	panicTotal.Inc()
	pm.logger.Error("Panic recorded")
}

// UpdateDBConnections updates database connection metrics
func (pm *PaymentMetrics) UpdateDBConnections(state string, count int) {
	dbConnections.With(prometheus.Labels{
		"state": state,
	}).Set(float64(count))
}

// CalculateSuccessRate calculates and updates payment success rates
func (pm *PaymentMetrics) CalculateSuccessRate(method, currency string, successCount, totalCount int) {
	if totalCount > 0 {
		successRate := float64(successCount) / float64(totalCount)
		paymentSuccessRate.With(prometheus.Labels{
			"method":   method,
			"currency": currency,
		}).Set(successRate)

		pm.logger.Info("Payment success rate calculated",
			zap.String("method", method),
			zap.String("currency", currency),
			zap.Int("success_count", successCount),
			zap.Int("total_count", totalCount),
			zap.Float64("success_rate", successRate))
	}
}

// MetricsCollector interface for dependency injection
type MetricsCollector interface {
	RecordPaymentAttempt(method, currency, status string, amount float64, duration time.Duration)
	RecordRefundAttempt(method, currency, status string, amount float64)
	RecordFiuuRequest(endpoint, method, status string, duration time.Duration)
	RecordWebhookReceived(source, eventType, status string, duration time.Duration)
	RecordError(errorType, component, severity string)
}

// Ensure PaymentMetrics implements MetricsCollector
var _ MetricsCollector = (*PaymentMetrics)(nil)

// MetricsMiddleware provides middleware for HTTP metrics
type MetricsMiddleware struct {
	metrics *PaymentMetrics
	logger  *zap.Logger
}

// NewMetricsMiddleware creates a new MetricsMiddleware
func NewMetricsMiddleware(metrics *PaymentMetrics, logger *zap.Logger) *MetricsMiddleware {
	return &MetricsMiddleware{
		metrics: metrics,
		logger:  logger,
	}
}

// WrapHandler wraps an HTTP handler to record metrics
func (mm *MetricsMiddleware) WrapHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a response writer wrapper to capture status code
		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		defer func() {
			duration := time.Since(start)
			statusCode := fmt.Sprintf("%d", wrapped.statusCode)

			mm.metrics.RecordHTTPRequest(
				r.Method,
				r.URL.Path,
				statusCode,
				duration,
			)

			// Record panic if it occurred
			if err := recover(); err != nil {
				mm.metrics.RecordPanic()
				mm.logger.Error("HTTP handler panic",
					zap.Any("panic", err),
					zap.String("method", r.Method),
					zap.String("path", r.URL.Path))
				panic(err) // Re-panic after recording metrics
			}
		}()

		handler.ServeHTTP(wrapped, r)
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}
