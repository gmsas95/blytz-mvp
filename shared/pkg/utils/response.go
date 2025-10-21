package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/gmsas95/blytz-mvp/shared/pkg/errors"
)

// SendSuccessResponse sends a success response with a status code and data
func SendSuccessResponse(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, gin.H{"status": "success", "data": data})
}

// SendErrorResponse sends an error response with a status code and error message
func SendErrorResponse(c *gin.Context, err error) {
	appErr, ok := err.(*errors.AppError)
	if !ok {
		appErr = errors.ErrInternalServer
	}

	c.JSON(appErr.StatusCode, gin.H{"status": "error", "message": appErr.Message})
}

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

type APIError struct {
	Type    string                 `json:"type"`
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(200, Response{
		Success: true,
		Data:    data,
	})
}

func ErrorResponse(c *gin.Context, err error) {
	if appErr, ok := errors.IsAppError(err); ok {
		c.JSON(appErr.StatusCode, Response{
			Success: false,
			Error: APIError{
				Type:    appErr.Type,
				Code:    appErr.Code,
				Message: appErr.Message,
				Details: appErr.Details,
			},
		})
		return
	}

	// Handle generic errors
	c.JSON(500, Response{
		Success: false,
		Error: APIError{
			Type:    "INTERNAL_ERROR",
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		},
	})
}

func ValidationError(c *gin.Context, message string, details map[string]interface{}) {
	c.JSON(400, Response{
		Success: false,
		Error: APIError{
			Type:    "VALIDATION_ERROR",
			Code:    "VALIDATION_ERROR",
			Message: message,
			Details: details,
		},
	})
}
