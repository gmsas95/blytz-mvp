package handlers

import (
	"context"
	"errors"
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

	// Validate the token using the auth service
	response, err := h.authService.ValidateToken(context.Background(), req.Token)
	if err != nil {
		utils.SendErrorResponse(c, err)
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, response)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, shared_errors.ErrInvalidRequestBody)
		return
	}

	// First validate the current token
	validationResponse, err := h.authService.ValidateToken(context.Background(), req.RefreshToken)
	if err != nil || !validationResponse.Valid {
		utils.SendErrorResponse(c, errors.New("invalid refresh token"))
		return
	}

	// Get user by ID to generate new token
	user, err := h.authService.GetUserByID(validationResponse.UserID)
	if err != nil {
		utils.SendErrorResponse(c, errors.New("user not found"))
		return
	}

	// Generate new token
	newToken, err := h.authService.GenerateJWT(user)
	if err != nil {
		utils.SendErrorResponse(c, err)
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, gin.H{
		"token":        newToken,
		"refreshToken": req.RefreshToken, // In production, you might want to rotate refresh tokens too
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req struct {
		Token string `json:"token"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		// If no token in body, try to get from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" && len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			req.Token = authHeader[7:]
		} else {
			utils.SendErrorResponse(c, shared_errors.ErrInvalidRequestBody)
			return
		}
	}

	// Validate the token first
	validationResponse, err := h.authService.ValidateToken(c.Request.Context(), req.Token)
	if err != nil || !validationResponse.Valid {
		utils.SendErrorResponse(c, errors.New("invalid token"))
		return
	}

	// In a production system, you would add the token to a blacklist here
	// For now, we'll just return success since the client should remove the token
	utils.SendSuccessResponse(c, http.StatusOK, gin.H{"message": "Logout successful"})
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, shared_errors.ErrInvalidRequestBody)
		return
	}

	// Get user ID from context (should be set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		utils.SendErrorResponse(c, errors.New("user not authenticated"))
		return
	}

	// Update user profile
	if err := h.authService.UpdateUserProfile(context.Background(), userID.(string), &req); err != nil {
		utils.SendErrorResponse(c, err)
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	// Get user ID from context (should be set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		utils.SendErrorResponse(c, errors.New("user not authenticated"))
		return
	}

	// Get user profile
	user, err := h.authService.GetUserByID(userID.(string))
	if err != nil {
		utils.SendErrorResponse(c, err)
		return
	}

	// Return user profile (excluding sensitive data like password)
	profile := gin.H{
		"id":           user.ID,
		"email":        user.Email,
		"display_name": user.DisplayName,
		"phone_number": user.PhoneNumber,
		"avatar_url":   user.AvatarURL,
		"is_active":    user.IsActive,
		"role":         user.Role,
		"created_at":   user.CreatedAt,
		"updated_at":   user.UpdatedAt,
	}

	utils.SendSuccessResponse(c, http.StatusOK, profile)
}