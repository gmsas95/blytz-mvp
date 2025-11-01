package monitoring

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LoggerConfig defines the configuration for the application logger
type LoggerConfig struct {
	Level      string `json:"level" yaml:"level"`
	Format     string `json:"format" yaml:"format"` // "json" or "console"
	Output     string `json:"output" yaml:"output"` // "stdout" or file path
	MaxSize    int    `json:"max_size" yaml:"max_size"`
	MaxBackups int    `json:"max_backups" yaml:"max_backups"`
	MaxAge     int    `json:"max_age" yaml:"max_age"`
	Compress   bool   `json:"compress" yaml:"compress"`
}

// DefaultLoggerConfig returns default logging configuration
func DefaultLoggerConfig() LoggerConfig {
	return LoggerConfig{
		Level:      "info",
		Format:     "json",
		Output:     "stdout",
		MaxSize:    100, // MB
		MaxBackups: 3,
		MaxAge:     28, // days
		Compress:   true,
	}
}

// NewLogger creates a new structured logger with the given configuration
func NewLogger(config LoggerConfig) (*zap.Logger, error) {
	zapConfig := zap.NewProductionConfig()

	// Parse log level
	level, err := zap.ParseAtomicLevel(config.Level)
	if err != nil {
		return nil, fmt.Errorf("invalid log level %s: %w", config.Level, err)
	}
	zapConfig.Level = level

	// Configure encoder
	if config.Format == "console" {
		zapConfig.Encoding = "console"
		zapConfig.EncoderConfig = zap.NewDevelopmentEncoderConfig()
		zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		zapConfig.Encoding = "json"
		zapConfig.EncoderConfig = zap.NewProductionEncoderConfig()
		zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		zapConfig.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	}

	// Configure output paths
	if config.Output == "stdout" {
		zapConfig.OutputPaths = []string{"stdout"}
		zapConfig.ErrorOutputPaths = []string{"stderr"}
	} else {
		zapConfig.OutputPaths = []string{config.Output}
		zapConfig.ErrorOutputPaths = []string{config.Output}
	}

	return zapConfig.Build()
}

// PaymentLogger provides structured logging for payment operations
type PaymentLogger struct {
	logger *zap.Logger
}

// NewPaymentLogger creates a new PaymentLogger
func NewPaymentLogger(logger *zap.Logger) *PaymentLogger {
	return &PaymentLogger{
		logger: logger,
	}
}

// LogPaymentRequest logs a payment request
func (pl *PaymentLogger) LogPaymentRequest(ctx context.Context, req PaymentRequestLog) {
	fields := []zap.Field{
		zap.String("event", "payment_request"),
		zap.String("payment_id", req.PaymentID),
		zap.String("order_id", req.OrderID),
		zap.String("user_id", req.UserID),
		zap.String("method", req.Method),
		zap.String("currency", req.Currency),
		zap.Float64("amount", req.Amount),
		zap.String("bill_name", req.BillName),
		zap.String("bill_email", req.BillEmail),
		zap.String("bill_mobile", req.BillMobile),
		zap.String("ip_address", req.IPAddress),
		zap.String("user_agent", req.UserAgent),
	}

	// Add metadata if present
	if req.Metadata != nil {
		fields = append(fields, zap.Any("metadata", req.Metadata))
	}

	// Add trace ID if present
	if traceID := GetTraceID(ctx); traceID != "" {
		fields = append(fields, zap.String("trace_id", traceID))
	}

	pl.logger.Info("Payment request received", fields...)
}

// LogPaymentResponse logs a payment response
func (pl *PaymentLogger) LogPaymentResponse(ctx context.Context, resp PaymentResponseLog) {
	fields := []zap.Field{
		zap.String("event", "payment_response"),
		zap.String("payment_id", resp.PaymentID),
		zap.String("order_id", resp.OrderID),
		zap.String("user_id", resp.UserID),
		zap.String("method", resp.Method),
		zap.String("status", resp.Status),
		zap.String("gateway_status", resp.GatewayStatus),
		zap.String("transaction_id", resp.TransactionID),
		zap.Float64("amount", resp.Amount),
		zap.String("currency", resp.Currency),
		zap.Duration("processing_time", resp.ProcessingTime),
		zap.String("gateway_response", resp.GatewayResponse),
	}

	// Add error if present
	if resp.Error != nil {
		fields = append(fields, zap.Error(resp.Error))
	}

	// Add trace ID if present
	if traceID := GetTraceID(ctx); traceID != "" {
		fields = append(fields, zap.String("trace_id", traceID))
	}

	if resp.Error != nil {
		pl.logger.Error("Payment response with error", fields...)
	} else {
		pl.logger.Info("Payment response sent", fields...)
	}
}

// LogPaymentRetry logs a payment retry attempt
func (pl *PaymentLogger) LogPaymentRetry(ctx context.Context, retry PaymentRetryLog) {
	fields := []zap.Field{
		zap.String("event", "payment_retry"),
		zap.String("payment_id", retry.PaymentID),
		zap.String("order_id", retry.OrderID),
		zap.String("user_id", retry.UserID),
		zap.Int("attempt_number", retry.AttemptNumber),
		zap.Int("max_attempts", retry.MaxAttempts),
		zap.Duration("delay", retry.Delay),
		zap.String("retry_reason", retry.RetryReason),
		zap.String("last_error", retry.LastError),
		zap.Time("next_attempt_at", retry.NextAttemptAt),
	}

	// Add trace ID if present
	if traceID := GetTraceID(ctx); traceID != "" {
		fields = append(fields, zap.String("trace_id", traceID))
	}

	pl.logger.Warn("Payment retry scheduled", fields...)
}

// LogRefundRequest logs a refund request
func (pl *PaymentLogger) LogRefundRequest(ctx context.Context, req RefundRequestLog) {
	fields := []zap.Field{
		zap.String("event", "refund_request"),
		zap.String("refund_id", req.RefundID),
		zap.String("payment_id", req.PaymentID),
		zap.String("order_id", req.OrderID),
		zap.String("user_id", req.UserID),
		zap.Float64("amount", req.Amount),
		zap.String("currency", req.Currency),
		zap.String("reason", req.Reason),
		zap.String("processed_by", req.ProcessedBy),
	}

	// Add trace ID if present
	if traceID := GetTraceID(ctx); traceID != "" {
		fields = append(fields, zap.String("trace_id", traceID))
	}

	pl.logger.Info("Refund request received", fields...)
}

// LogRefundResponse logs a refund response
func (pl *PaymentLogger) LogRefundResponse(ctx context.Context, resp RefundResponseLog) {
	fields := []zap.Field{
		zap.String("event", "refund_response"),
		zap.String("refund_id", resp.RefundID),
		zap.String("payment_id", resp.PaymentID),
		zap.String("order_id", resp.OrderID),
		zap.String("user_id", resp.UserID),
		zap.String("status", resp.Status),
		zap.String("gateway_status", resp.GatewayStatus),
		zap.Float64("amount", resp.Amount),
		zap.String("currency", resp.Currency),
		zap.Duration("processing_time", resp.ProcessingTime),
	}

	// Add error if present
	if resp.Error != nil {
		fields = append(fields, zap.Error(resp.Error))
	}

	// Add trace ID if present
	if traceID := GetTraceID(ctx); traceID != "" {
		fields = append(fields, zap.String("trace_id", traceID))
	}

	if resp.Error != nil {
		pl.logger.Error("Refund response with error", fields...)
	} else {
		pl.logger.Info("Refund response sent", fields...)
	}
}

// LogWebhookReceived logs a received webhook
func (pl *PaymentLogger) LogWebhookReceived(ctx context.Context, webhook WebhookLog) {
	fields := []zap.Field{
		zap.String("event", "webhook_received"),
		zap.String("webhook_id", webhook.WebhookID),
		zap.String("source", webhook.Source),
		zap.String("event_type", webhook.EventType),
		zap.String("source_id", webhook.SourceID),
		zap.String("signature", webhook.Signature),
		zap.String("ip_address", webhook.IPAddress),
		zap.String("user_agent", webhook.UserAgent),
		zap.Time("received_at", webhook.ReceivedAt),
	}

	// Add payload if present (masked for sensitive data)
	if webhook.Payload != nil {
		// Mask sensitive fields in payload
		maskedPayload := maskSensitiveData(webhook.Payload)
		fields = append(fields, zap.Any("payload", maskedPayload))
	}

	// Add trace ID if present
	if traceID := GetTraceID(ctx); traceID != "" {
		fields = append(fields, zap.String("trace_id", traceID))
	}

	pl.logger.Info("Webhook received", fields...)
}

// LogWebhookProcessed logs a processed webhook
func (pl *PaymentLogger) LogWebhookProcessed(ctx context.Context, webhook WebhookProcessLog) {
	fields := []zap.Field{
		zap.String("event", "webhook_processed"),
		zap.String("webhook_id", webhook.WebhookID),
		zap.String("source", webhook.Source),
		zap.String("event_type", webhook.EventType),
		zap.String("status", webhook.Status),
		zap.Duration("processing_time", webhook.ProcessingTime),
		zap.Int("retry_count", webhook.RetryCount),
	}

	// Add error if present
	if webhook.Error != nil {
		fields = append(fields, zap.Error(webhook.Error))
	}

	// Add trace ID if present
	if traceID := GetTraceID(ctx); traceID != "" {
		fields = append(fields, zap.String("trace_id", traceID))
	}

	if webhook.Error != nil {
		pl.logger.Error("Webhook processing failed", fields...)
	} else {
		pl.logger.Info("Webhook processed successfully", fields...)
	}
}

// LogCircuitBreakerEvent logs circuit breaker events
func (pl *PaymentLogger) LogCircuitBreakerEvent(ctx context.Context, event CircuitBreakerLog) {
	fields := []zap.Field{
		zap.String("event", "circuit_breaker"),
		zap.String("service", event.Service),
		zap.String("old_state", event.OldState),
		zap.String("new_state", event.NewState),
		zap.Int("failure_count", event.FailureCount),
		zap.Duration("timeout", event.Timeout),
		zap.Time("changed_at", event.ChangedAt),
	}

	// Add trace ID if present
	if traceID := GetTraceID(ctx); traceID != "" {
		fields = append(fields, zap.String("trace_id", traceID))
	}

	pl.logger.Warn("Circuit breaker state changed", fields...)
}

// LogDatabaseOperation logs database operations
func (pl *PaymentLogger) LogDatabaseOperation(ctx context.Context, db DBLog) {
	fields := []zap.Field{
		zap.String("event", "database_operation"),
		zap.String("operation", db.Operation),
		zap.String("table", db.Table),
		zap.Duration("duration", db.Duration),
		zap.Bool("success", db.Success),
	}

	// Add row count if present
	if db.RowCount != nil {
		fields = append(fields, zap.Int("row_count", *db.RowCount))
	}

	// Add error if present
	if db.Error != nil {
		fields = append(fields, zap.Error(db.Error))
	}

	// Add trace ID if present
	if traceID := GetTraceID(ctx); traceID != "" {
		fields = append(fields, zap.String("trace_id", traceID))
	}

	if db.Error != nil {
		pl.logger.Error("Database operation failed", fields...)
	} else {
		pl.logger.Debug("Database operation completed", fields...)
	}
}

// LogBusinessMetric logs business metrics
func (pl *PaymentLogger) LogBusinessMetric(ctx context.Context, metric BusinessMetricLog) {
	fields := []zap.Field{
		zap.String("event", "business_metric"),
		zap.String("metric_name", metric.MetricName),
		zap.Float64("value", metric.Value),
		zap.String("unit", metric.Unit),
		zap.Time("measured_at", metric.MeasuredAt),
	}

	// Add dimensions if present
	if metric.Dimensions != nil {
		fields = append(fields, zap.Any("dimensions", metric.Dimensions))
	}

	// Add trace ID if present
	if traceID := GetTraceID(ctx); traceID != "" {
		fields = append(fields, zap.String("trace_id", traceID))
	}

	pl.logger.Info("Business metric recorded", fields...)
}

// LogSecurityEvent logs security events
func (pl *PaymentLogger) LogSecurityEvent(ctx context.Context, event SecurityEventLog) {
	fields := []zap.Field{
		zap.String("event", "security_event"),
		zap.String("event_type", event.EventType),
		zap.String("severity", event.Severity),
		zap.String("user_id", event.UserID),
		zap.String("ip_address", event.IPAddress),
		zap.String("user_agent", event.UserAgent),
		zap.String("description", event.Description),
		zap.Time("occurred_at", event.OccurredAt),
	}

	// Add additional details if present
	if event.Details != nil {
		fields = append(fields, zap.Any("details", event.Details))
	}

	// Add trace ID if present
	if traceID := GetTraceID(ctx); traceID != "" {
		fields = append(fields, zap.String("trace_id", traceID))
	}

	// Log based on severity
	switch event.Severity {
	case "critical", "high":
		pl.logger.Error("Security event detected", fields...)
	case "medium":
		pl.logger.Warn("Security event detected", fields...)
	case "low":
		pl.logger.Info("Security event detected", fields...)
	default:
		pl.logger.Info("Security event detected", fields...)
	}
}

// LogPerformanceAlert logs performance alerts
func (pl *PaymentLogger) LogPerformanceAlert(ctx context.Context, alert PerformanceAlertLog) {
	fields := []zap.Field{
		zap.String("event", "performance_alert"),
		zap.String("alert_type", alert.AlertType),
		zap.String("metric", alert.Metric),
		zap.Float64("value", alert.Value),
		zap.Float64("threshold", alert.Threshold),
		zap.String("severity", alert.Severity),
		zap.String("description", alert.Description),
		zap.Time("detected_at", alert.DetectedAt),
	}

	// Add additional dimensions if present
	if alert.Dimensions != nil {
		fields = append(fields, zap.Any("dimensions", alert.Dimensions))
	}

	// Add trace ID if present
	if traceID := GetTraceID(ctx); traceID != "" {
		fields = append(fields, zap.String("trace_id", traceID))
	}

	// Log based on severity
	switch alert.Severity {
	case "critical":
		pl.logger.Error("Performance alert triggered", fields...)
	case "warning":
		pl.logger.Warn("Performance alert triggered", fields...)
	default:
		pl.logger.Info("Performance alert triggered", fields...)
	}
}

// Log types
type PaymentRequestLog struct {
	PaymentID  string
	OrderID    string
	UserID     string
	Method     string
	Currency   string
	Amount     float64
	BillName   string
	BillEmail  string
	BillMobile string
	IPAddress  string
	UserAgent  string
	Metadata   map[string]interface{}
}

type PaymentResponseLog struct {
	PaymentID       string
	OrderID         string
	UserID          string
	Method          string
	Status          string
	GatewayStatus   string
	TransactionID   string
	Amount          float64
	Currency        string
	ProcessingTime  time.Duration
	GatewayResponse string
	Error           error
}

type PaymentRetryLog struct {
	PaymentID     string
	OrderID       string
	UserID        string
	AttemptNumber int
	MaxAttempts   int
	Delay         time.Duration
	RetryReason   string
	LastError     string
	NextAttemptAt time.Time
}

type RefundRequestLog struct {
	RefundID    string
	PaymentID   string
	OrderID     string
	UserID      string
	Amount      float64
	Currency    string
	Reason      string
	ProcessedBy string
}

type RefundResponseLog struct {
	RefundID       string
	PaymentID      string
	OrderID        string
	UserID         string
	Status         string
	GatewayStatus  string
	Amount         float64
	Currency       string
	ProcessingTime time.Duration
	Error          error
}

type WebhookLog struct {
	WebhookID string
	Source    string
	EventType string
	SourceID  string
	Signature string
	Payload   map[string]interface{}
	IPAddress string
	UserAgent string
	ReceivedAt time.Time
}

type WebhookProcessLog struct {
	WebhookID     string
	Source        string
	EventType     string
	Status        string
	ProcessingTime time.Duration
	RetryCount    int
	Error         error
}

type CircuitBreakerLog struct {
	Service      string
	OldState     string
	NewState     string
	FailureCount int
	Timeout      time.Duration
	ChangedAt    time.Time
}

type DBLog struct {
	Operation string
	Table     string
	Duration  time.Duration
	Success   bool
	RowCount  *int
	Error     error
}

type BusinessMetricLog struct {
	MetricName string
	Value      float64
	Unit       string
	MeasuredAt time.Time
	Dimensions map[string]string
}

type SecurityEventLog struct {
	EventType   string
	Severity    string
	UserID      string
	IPAddress   string
	UserAgent   string
	Description string
	OccurredAt  time.Time
	Details     map[string]interface{}
}

type PerformanceAlertLog struct {
	AlertType  string
	Metric     string
	Value      float64
	Threshold  float64
	Severity   string
	Description string
	DetectedAt time.Time
	Dimensions map[string]string
}

// Helper functions

// TraceIDKey is the context key for trace ID
type TraceIDKey struct{}

// GetTraceID extracts trace ID from context
func GetTraceID(ctx context.Context) string {
	if traceID, ok := ctx.Value(TraceIDKey{}).(string); ok {
		return traceID
	}
	return ""
}

// maskSensitiveData masks sensitive data in webhook payloads
func maskSensitiveData(data map[string]interface{}) map[string]interface{} {
	masked := make(map[string]interface{})

	sensitiveFields := []string{
		"card_number", "cvv", "expiry", "pan", "cvc",
		"account_number", "routing_number", "ssn",
		"password", "token", "secret", "key",
		"bill_email", "bill_mobile", "user_id",
	}

	for key, value := range data {
		// Check if this is a sensitive field
		isSensitive := false
		for _, sensitive := range sensitiveFields {
			if containsIgnoreCase(key, sensitive) {
				isSensitive = true
				break
			}
		}

		if isSensitive {
			masked[key] = "***MASKED***"
		} else if nested, ok := value.(map[string]interface{}); ok {
			masked[key] = maskSensitiveData(nested)
		} else {
			masked[key] = value
		}
	}

	return masked
}

// containsIgnoreCase checks if a string contains a substring ignoring case
func containsIgnoreCase(s, substr string) bool {
	s = strings.ToLower(s)
	substr = strings.ToLower(substr)
	return strings.Contains(s, substr)
}

// WithTraceID adds trace ID to context
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, TraceIDKey{}, traceID)
}