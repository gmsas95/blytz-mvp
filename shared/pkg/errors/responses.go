package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorCode represents standardized error codes
type ErrorCode string

const (
	// Validation errors
	ErrorCodeInvalidRequest ErrorCode = "invalid_request"
	ErrorCodeMissingField   ErrorCode = "missing_field"
	ErrorCodeInvalidFormat  ErrorCode = "invalid_format"

	// Authentication errors
	ErrorCodeUnauthorized ErrorCode = "unauthorized"
	ErrorCodeForbidden    ErrorCode = "forbidden"
	ErrorCodeTokenExpired ErrorCode = "token_expired"

	// Business logic errors
	ErrorCodeNotFound         ErrorCode = "not_found"
	ErrorCodeAlreadyExists    ErrorCode = "already_exists"
	ErrorCodeInvalidOperation ErrorCode = "invalid_operation"

	// System errors
	ErrorCodeInternalServer     ErrorCode = "internal_server_error"
	ErrorCodeServiceUnavailable ErrorCode = "service_unavailable"
	ErrorCodeRateLimited        ErrorCode = "rate_limited"
)

// APIError represents a standardized API error response
type APIError struct {
	Code    ErrorCode   `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// ErrorResponse represents the complete error response structure
type ErrorResponse struct {
	Success bool     `json:"success"`
	Error   APIError `json:"error"`
}

// SuccessResponse represents the standard success response
type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
}

// PaginatedResponse represents paginated data response
type PaginatedResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Page    int         `json:"page"`
	PerPage int         `json:"per_page"`
	Total   int64       `json:"total"`
}

// NewAPIError creates a new API error
func NewAPIError(code ErrorCode, message string, details interface{}) APIError {
	return APIError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// RespondWithError sends a standardized error response
func RespondWithError(c *gin.Context, statusCode int, err APIError) {
	response := ErrorResponse{
		Success: false,
		Error:   err,
	}
	c.JSON(statusCode, response)
}

// RespondWithSuccess sends a standardized success response
func RespondWithSuccess(c *gin.Context, data interface{}) {
	response := SuccessResponse{
		Success: true,
		Data:    data,
	}
	c.JSON(http.StatusOK, response)
}

// RespondWithSuccessMessage sends a success response with message
func RespondWithSuccessMessage(c *gin.Context, data interface{}, message string) {
	response := SuccessResponse{
		Success: true,
		Data:    data,
		Message: message,
	}
	c.JSON(http.StatusOK, response)
}

// RespondWithPagination sends a paginated response
func RespondWithPagination(c *gin.Context, data interface{}, page, perPage int, total int64) {
	response := PaginatedResponse{
		Success: true,
		Data:    data,
		Page:    page,
		PerPage: perPage,
		Total:   total,
	}
	c.JSON(http.StatusOK, response)
}

// Common error responses
var (
	ErrInvalidRequestResponse     = NewAPIError(ErrorCodeInvalidRequest, "Invalid request format", nil)
	ErrUnauthorizedResponse       = NewAPIError(ErrorCodeUnauthorized, "Authentication required", nil)
	ErrForbiddenResponse          = NewAPIError(ErrorCodeForbidden, "Access denied", nil)
	ErrNotFoundResponse           = NewAPIError(ErrorCodeNotFound, "Resource not found", nil)
	ErrRateLimitedResponse        = NewAPIError(ErrorCodeRateLimited, "Rate limit exceeded", nil)
	ErrInternalServerResponse     = NewAPIError(ErrorCodeInternalServer, "Internal server error", nil)
	ErrServiceUnavailableResponse = NewAPIError(ErrorCodeServiceUnavailable, "Service temporarily unavailable", nil)
)
