package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/blytz/auth-service/internal/models"
	"github.com/blytz/auth-service/internal/services"
	"github.com/blytz/shared/utils"
	"github.com/blytz/shared/errors"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// SignUp handles user registration
func (h *AuthHandler) SignUp(c *gin.Context) {
	var req models.AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, "Invalid request data", gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.SignUp(c.Request.Context(), &req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, user)
}

// Login handles user authentication
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, "Invalid request data", gin.H{"error": err.Error()})
		return
	}

	response, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, response)
}

// Verify handles JWT token verification
func (h *AuthHandler) Verify(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, errors.Unauthorized("NO_USER_ID", "User ID not found"))
		return
	}

	user, err := h.authService.GetUserByID(c.Request.Context(), userID.(string))
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, models.VerifyResponse{
		User:  *user,
		Valid: true,
	})
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	refreshToken := c.GetHeader("X-Refresh-Token")
	if refreshToken == "" {
		utils.ErrorResponse(c, errors.ValidationError("MISSING_REFRESH_TOKEN", "Refresh token is required"))
		return
	}

	response, err := h.authService.RefreshToken(c.Request.Context(), refreshToken)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, response)
}
