package errors

import (
	"fmt"
	"net/http"
)

// Error types
const (
	ErrTypeValidation     = "VALIDATION_ERROR"
	ErrTypeAuthentication = "AUTHENTICATION_ERROR"
	ErrTypeAuthorization  = "AUTHORIZATION_ERROR"
	ErrTypeNotFound       = "NOT_FOUND_ERROR"
	ErrTypeConflict       = "CONFLICT_ERROR"
	ErrTypeInternal       = "INTERNAL_ERROR"
	ErrTypeService        = "SERVICE_ERROR"
	ErrTypeDatabase       = "DATABASE_ERROR"
	ErrTypeRateLimit      = "RATE_LIMIT_ERROR"
)

// AppError represents a standardized application error
type AppError struct {
	Type       string                 `json:"type"`
	Code       string                 `json:"code"`
	Message    string                 `json:"message"`
	Details    map[string]interface{} `json:"details,omitempty"`
	StatusCode int                    `json:"status_code"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s: %s", e.Type, e.Code, e.Message)
}

// New creates a new application error
func New(errType, code, message string, statusCode int) *AppError {
	return &AppError{
		Type:       errType,
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
		Details:    make(map[string]interface{}),
	}
}

// WithDetails adds details to the error
func (e *AppError) WithDetails(details map[string]interface{}) *AppError {
	e.Details = details
	return e
}

// Common error constructors

// ValidationError creates a validation error
func ValidationError(code, message string) *AppError {
	return New(ErrTypeValidation, code, message, http.StatusBadRequest)
}

// AuthenticationError creates an authentication error
func AuthenticationError(code, message string) *AppError {
	return New(ErrTypeAuthentication, code, message, http.StatusUnauthorized)
}

// AuthorizationError creates an authorization error
func AuthorizationError(code, message string) *AppError {
	return New(ErrTypeAuthorization, code, message, http.StatusForbidden)
}

// NotFoundError creates a not found error
func NotFoundError(code, message string) *AppError {
	return New(ErrTypeNotFound, code, message, http.StatusNotFound)
}

// ConflictError creates a conflict error
func ConflictError(code, message string) *AppError {
	return New(ErrTypeConflict, code, message, http.StatusConflict)
}

// InternalError creates an internal server error
func InternalError(code, message string) *AppError {
	return New(ErrTypeInternal, code, message, http.StatusInternalServerError)
}

// ServiceError creates a service error
func ServiceError(code, message string) *AppError {
	return New(ErrTypeService, code, message, http.StatusServiceUnavailable)
}

// DatabaseError creates a database error
func DatabaseError(code, message string) *AppError {
	return New(ErrTypeDatabase, code, message, http.StatusInternalServerError)
}

// RateLimitError creates a rate limit error
func RateLimitError(code, message string) *AppError {
	return New(ErrTypeRateLimit, code, message, http.StatusTooManyRequests)
}

// Common error instances
var (
	ErrInvalidRequest     = ValidationError("INVALID_REQUEST", "Invalid request data")
	ErrUnauthorized       = AuthenticationError("UNAUTHORIZED", "Authentication required")
	ErrForbidden          = AuthorizationError("FORBIDDEN", "Insufficient permissions")
	ErrNotFound           = NotFoundError("NOT_FOUND", "Resource not found")
	ErrConflict           = ConflictError("CONFLICT", "Resource conflict")
	ErrInternalServer     = InternalError("INTERNAL_ERROR", "Internal server error")
	ErrServiceUnavailable = ServiceError("SERVICE_UNAVAILABLE", "Service temporarily unavailable")
	ErrDatabaseError      = DatabaseError("DATABASE_ERROR", "Database operation failed")
	ErrRateLimitExceeded  = RateLimitError("RATE_LIMIT_EXCEEDED", "Rate limit exceeded")
	// Auction-specific errors
	ErrAuctionEnded = ConflictError("AUCTION_ENDED", "Auction has already ended")
	ErrBidTooLow    = ValidationError("BID_TOO_LOW", "Bid amount is too low")
	// Auth-specific errors
	ErrInvalidRequestBody = ValidationError("INVALID_REQUEST_BODY", "Invalid request body")
	ErrNotImplemented     = ServiceError("NOT_IMPLEMENTED", "Feature not implemented")
	// Product-specific errors
	ErrInsufficientStock = ConflictError("INSUFFICIENT_STOCK", "Insufficient stock available")
)

// WrapError wraps an existing error with additional context
func WrapError(err error, errType, code, message string) *AppError {
	appErr := New(errType, code, message, http.StatusInternalServerError)
	if appErr.Details == nil {
		appErr.Details = make(map[string]interface{})
	}
	appErr.Details["original_error"] = err.Error()
	return appErr
}

// IsAppError checks if an error is an AppError
func IsAppError(err error) (*AppError, bool) {
	appErr, ok := err.(*AppError)
	return appErr, ok
}
