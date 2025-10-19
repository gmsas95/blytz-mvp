package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gmsas95/blytz-mvp/services/auth-service/internal/models"
	"github.com/gmsas95/blytz-mvp/services/auth-service/internal/services"
	"github.com/gmsas95/blytz-mvp/shared/pkg/utils"
	shared_errors "github.com/gmsas95/blytz-mvp/shared/pkg/errors"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.SendErrorResponse(c, shared_errors.ErrInvalidRequestBody)
		return
	}

	if err := h.authService.RegisterUser(&user); err != nil {
		utils.SendErrorResponse(c, err)
		return
	}

	utils.SendSuccessResponse(c, http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var loginDetails struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&loginDetails); err != nil {
		utils.SendErrorResponse(c, shared_errors.ErrInvalidRequestBody)
		return
	}

	token, err := h.authService.LoginUser(loginDetails.Email, loginDetails.Password)
	if err != nil {
		utils.SendErrorResponse(c, err)
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, gin.H{"token": token})
}

func (h *AuthHandler) SignUp(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, shared_errors.ErrInvalidRequestBody)
		return
	}

	user := &models.User{
		Email:       req.Email,
		Password:    req.Password,
		DisplayName: req.DisplayName,
		PhoneNumber: req.PhoneNumber,
	}

	if err := h.authService.RegisterUser(user); err != nil {
		utils.SendErrorResponse(c, err)
		return
	}

	utils.SendSuccessResponse(c, http.StatusCreated, user)
}

func (h *AuthHandler) Verify(c *gin.Context) {
	var req models.ValidateTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, shared_errors.ErrInvalidRequestBody)
		return
	}

	// For now, just return a simple validation response
	// In a real implementation, you'd validate the JWT token
	utils.SendSuccessResponse(c, http.StatusOK, models.ValidateTokenResponse{
		Valid:   true,
		Message: "Token is valid",
	})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, shared_errors.ErrInvalidRequestBody)
		return
	}

	// For now, just return the same token
	// In a real implementation, you'd generate a new token
	utils.SendSuccessResponse(c, http.StatusOK, gin.H{
		"token":        req.RefreshToken, // This should be a new token
		"refreshToken": req.RefreshToken,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// For now, just return success
	// In a real implementation, you'd invalidate the token
	utils.SendSuccessResponse(c, http.StatusOK, gin.H{"message": "Logout successful"})
}