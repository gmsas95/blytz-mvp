package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/blytz/auth-service/internal/models"
	"github.com/blytz/auth-service/internal/services"
	"github.com/blytz/shared/errors"
	"github.com/blytz/shared/utils"
	"go.uber.org/zap"
)

// AuthHandler handles authentication HTTP requests
type AuthHandler struct {
	authService *services.AuthService
	logger      *zap.Logger
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(authService *services.AuthService, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

// RegisterUser handles user registration
func (h *AuthHandler) RegisterUser(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, "Invalid request data", gin.H{"error": err.Error()})
		return
	}

	// Validate request
	if err := h.validateRegisterRequest(&req); err != nil {
		utils.ValidationError(c, "Validation failed", gin.H{"error": err.Error()})
		return
	}

	// Create user with Better Auth
	user, err := h.authService.CreateUser(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create user", zap.Error(err))
		utils.ErrorResponse(c, err)
		return
	}

	// Generate JWT token
	token, err := h.generateJWT(user.ID, user.Email)
	if err != nil {
		h.logger.Error("Failed to generate JWT", zap.Error(err))
		utils.ErrorResponse(c, err)
		return
	}

	h.logger.Info("User registered successfully", zap.String("user_id", user.ID), zap.String("email", user.Email))
	
	utils.SuccessResponse(c, models.AuthResponse{
		Success: true,
		Message: "User registered successfully",
		User:    user,
		Token:   token,
	})
}

// LoginUser handles user login
func (h *AuthHandler) LoginUser(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, "Invalid request data", gin.H{"error": err.Error()})
		return
	}

	// Authenticate user with Better Auth
	user, token, err := h.authService.AuthenticateUser(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Authentication failed", zap.Error(err))
		utils.ErrorResponse(c, err)
		return
	}

	h.logger.Info("User logged in successfully", zap.String("user_id", user.ID), zap.String("email", user.Email))
	
	utils.SuccessResponse(c, models.AuthResponse{
		Success: true,
		Message: "Login successful",
		User:    user,
		Token:   token,
	})
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, "Invalid request data", gin.H{"error": err.Error()})
		return
	}

	// Refresh token with Better Auth
	token, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		h.logger.Error("Token refresh failed", zap.Error(err))
		utils.ErrorResponse(c, err)
		return
	}

	h.logger.Info("Token refreshed successfully")
	
	utils.SuccessResponse(c, models.SuccessResponse{
		Success: true,
		Message: "Token refreshed successfully",
		Data: gin.H{
			"token": token,
		},
	})
}

// ValidateToken handles token validation (internal endpoint)
func (h *AuthHandler) ValidateToken(c *gin.Context) {
	var req models.ValidateTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, "Invalid request data", gin.H{"error": err.Error()})
		return
	}

	// Validate token with Better Auth
	response, err := h.authService.ValidateToken(c.Request.Context(), req.Token)
	if err != nil {
		h.logger.Error("Token validation failed", zap.Error(err))
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, response)
}

// GetCurrentUser gets the current authenticated user
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, errors.AuthenticationError("NO_AUTH", "User not authenticated"))
		return
	}

	// Get user by ID
	user, err := h.authService.GetUserByID(c.Request.Context(), userID.(string))
	if err != nil {
		h.logger.Error("Failed to get user", zap.Error(err))
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, models.SuccessResponse{
		Success: true,
		Message: "User retrieved successfully",
		Data:    user,
	})
}

// UpdateProfile handles profile updates
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, errors.AuthenticationError("NO_AUTH", "User not authenticated"))
		return
	}

	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, "Invalid request data", gin.H{"error": err.Error()})
		return
	}

	// Update profile
	if err := h.authService.UpdateProfile(c.Request.Context(), userID.(string), &req); err != nil {
		h.logger.Error("Failed to update profile", zap.Error(err))
		utils.ErrorResponse(c, err)
		return
	}

	h.logger.Info("Profile updated successfully", zap.String("user_id", userID.(string)))
	
	utils.SuccessResponse(c, models.SuccessResponse{
		Success: true,
		Message: "Profile updated successfully",
	})
}

// ChangePassword handles password changes
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, errors.AuthenticationError("NO_AUTH", "User not authenticated"))
		return
	}

	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, "Invalid request data", gin.H{"error": err.Error()})
		return
	}

	// Change password
	if err := h.authService.ChangePassword(c.Request.Context(), userID.(string), &req); err != nil {
		h.logger.Error("Failed to change password", zap.Error(err))
		utils.ErrorResponse(c, err)
		return
	}

	h.logger.Info("Password changed successfully", zap.String("user_id", userID.(string)))
	
	utils.SuccessResponse(c, models.SuccessResponse{
		Success: true,
		Message: "Password changed successfully",
	})
}

// validateRegisterRequest validates registration request
func (h *AuthHandler) validateRegisterRequest(req *models.RegisterRequest) error {
	// Basic validation (Better Auth will do more comprehensive validation)
	if req.Email == "" || req.Password == "" || req.DisplayName == "" {
		return errors.ValidationError("MISSING_FIELDS", "Email, password, and display name are required")
	}
	return nil
}

// generateJWT generates JWT token for user
func (h *AuthHandler) generateJWT(userID, email string) (string, error) {
	// This would use Better Auth's JWT generation
	// For now, we'll use a simple implementation
	// In production, this would use Better Auth's built-in JWT generation
	
	// For this implementation, we'll use Better Auth's token generation
	// The actual JWT generation would be handled by Better Auth
	return h.authService.betterAuth.GenerateJWT(userID, email)
}
