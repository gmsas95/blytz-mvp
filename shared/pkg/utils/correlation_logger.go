package utils

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// CorrelationIDKey is the context key for correlation ID
type CorrelationIDKey struct{}

// LoggerConfig holds logging configuration
type LoggerConfig struct {
	Level       string
	Format      string
	Service     string
	Version     string
	Environment string
}

// StructuredLogger wraps zap logger with correlation support
type StructuredLogger struct {
	logger *zap.Logger
	config LoggerConfig
}

// NewStructuredLogger creates a new structured logger
func NewStructuredLogger(config LoggerConfig) (*StructuredLogger, error) {
	var zapConfig zap.Config
	if config.Environment == "production" {
		zapConfig = zap.NewProductionConfig()
	} else {
		zapConfig = zap.NewDevelopmentConfig()
	}

	zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	zapConfig.OutputPaths = []string{"stdout"}
	zapConfig.ErrorOutputPaths = []string{"stderr"}

	logger, err := zapConfig.Build()
	if err != nil {
		return nil, err
	}

	return &StructuredLogger{
		logger: logger,
		config: config,
	}, nil
}

// GetLogger returns the underlying zap logger
func (sl *StructuredLogger) GetLogger() *zap.Logger {
	return sl.logger
}

// WithCorrelation adds correlation ID to log fields
func (sl *StructuredLogger) WithCorrelation(correlationID string) *zap.Logger {
	return sl.logger.With(
		zap.String("correlation_id", correlationID),
		zap.String("service", sl.config.Service),
		zap.String("version", sl.config.Version),
		zap.String("environment", sl.config.Environment),
	)
}

// LogRequest logs HTTP request with correlation
func (sl *StructuredLogger) LogRequest(c *gin.Context, correlationID string) {
	sl.WithCorrelation(correlationID).Info("HTTP Request",
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.String("query", c.Request.URL.RawQuery),
		zap.String("user_agent", c.Request.UserAgent()),
		zap.String("client_ip", c.ClientIP()),
		zap.Time("start_time", time.Now()),
	)
}

// LogResponse logs HTTP response with correlation
func (sl *StructuredLogger) LogResponse(c *gin.Context, correlationID string, statusCode int, duration time.Duration) {
	sl.WithCorrelation(correlationID).Info("HTTP Response",
		zap.Int("status_code", statusCode),
		zap.Duration("duration", duration),
		zap.String("client_ip", c.ClientIP()),
	)
}

// LogError logs error with correlation
func (sl *StructuredLogger) LogError(c *gin.Context, correlationID string, err error, message string) {
	sl.WithCorrelation(correlationID).Error(message,
		zap.Error(err),
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.String("client_ip", c.ClientIP()),
	)
}

// LogBusiness logs business events with correlation
func (sl *StructuredLogger) LogBusiness(correlationID string, event string, data map[string]interface{}) {
	fields := []zap.Field{zap.String("event", event)}
	for k, v := range data {
		fields = append(fields, zap.Any(k, v))
	}
	sl.WithCorrelation(correlationID).Info("Business Event", fields...)
}

// GenerateCorrelationID creates a new correlation ID
func GenerateCorrelationID() string {
	return uuid.New().String()
}

// GetCorrelationID extracts correlation ID from context
func GetCorrelationID(c context.Context) string {
	if correlationID, ok := c.Value(CorrelationIDKey{}).(string); ok {
		return correlationID
	}
	return ""
}

// CorrelationMiddleware adds correlation ID to requests
func CorrelationMiddleware(logger *StructuredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate or extract correlation ID
		correlationID := c.GetHeader("X-Correlation-ID")
		if correlationID == "" {
			correlationID = GenerateCorrelationID()
		}

		// Add to context
		c.Set("correlation_id", correlationID)
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), CorrelationIDKey{}, correlationID))

		// Add correlation ID to response headers
		c.Header("X-Correlation-ID", correlationID)

		// Log request
		startTime := time.Now()
		logger.LogRequest(c, correlationID)

		// Process request
		c.Next()

		// Log response
		duration := time.Since(startTime)
		logger.LogResponse(c, correlationID, c.Writer.Status(), duration)
	}
}
