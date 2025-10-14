package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/blytz/shared/errors"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

type ErrorResponse struct {
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
			Error: ErrorResponse{
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
		Error: ErrorResponse{
			Type:    "INTERNAL_ERROR",
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		},
	})
}

func ValidationError(c *gin.Context, message string, details map[string]interface{}) {
	c.JSON(400, Response{
		Success: false,
		Error: ErrorResponse{
			Type:    "VALIDATION_ERROR",
			Code:    "VALIDATION_ERROR",
			Message: message,
			Details: details,
		},
	})
}
