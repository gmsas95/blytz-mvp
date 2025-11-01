package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/gmsas95/blytz-mvp/services/auth-service/internal/models"
	"github.com/gmsas95/blytz-mvp/services/auth-service/internal/services"
	shared_errors "github.com/gmsas95/blytz-mvp/shared/pkg/errors"
	"github.com/gmsas95/blytz-mvp/shared/pkg/utils"
)

// AuthHandler handles authentication requests
type AuthHandler struct {
	authService *services.AuthService
	logger      *zap.Logger
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	logger, _ := zap.NewDevelopment()
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

// SignUp handles user registration
func (h *AuthHandler) SignUp(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, shared_errors.ErrInvalidRequestBody)
		return
	}

	user := &models.User{
		Email:       req.Email,
		DisplayName: req.DisplayName,
		PhoneNumber: req.PhoneNumber,
		Password:    req.Password,
	}

	if err := h.authService.RegisterUser(user); err != nil {
		utils.SendErrorResponse(c, err)
		return
	}

	utils.SendSuccessResponse(c, http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, shared_errors.ErrInvalidRequestBody)
		return
	}

	token, err := h.authService.LoginUser(req.Email, req.Password)
	if err != nil {
		utils.SendErrorResponse(c, err)
		return
	}

	user, err := h.authService.GetUserByEmail(req.Email)
	if err != nil {
		utils.SendErrorResponse(c, err)
		return
	}

	// Remove password from response
	user.Password = ""

	utils.SendSuccessResponse(c, http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}

// GetProfile handles getting user profile
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		utils.SendErrorResponse(c, shared_errors.ErrUnauthorized)
		return
	}

	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		utils.SendErrorResponse(c, err)
		return
	}

	// Remove password from response
	user.Password = ""

	utils.SendSuccessResponse(c, http.StatusOK, user)
}

// Verify handles token verification
func (h *AuthHandler) Verify(c *gin.Context) {
	var req models.ValidateTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, shared_errors.ErrInvalidRequestBody)
		return
	}

	response, err := h.authService.ValidateToken(c.Request.Context(), req.Token)
	if err != nil {
		utils.SendErrorResponse(c, err)
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, response)
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, shared_errors.ErrInvalidRequestBody)
		return
	}

	// TODO: Implement refresh token logic
	utils.SendErrorResponse(c, shared_errors.ErrNotImplemented)
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *gin.Context) {
	// TODO: Implement logout logic (token blacklisting)
	utils.SendSuccessResponse(c, http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// UpdateProfile handles profile updates
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		utils.SendErrorResponse(c, shared_errors.ErrUnauthorized)
		return
	}

	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, shared_errors.ErrInvalidRequestBody)
		return
	}

	if err := h.authService.UpdateUserProfile(c.Request.Context(), userID, &req); err != nil {
		utils.SendErrorResponse(c, err)
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

// RegisterRoutes registers the API routes for the auth service
func RegisterRoutes(router *gin.Engine, authService *services.AuthService) {
	api := router.Group("/api/v1")
	{
		api.POST("/register", registerUser(authService))
		api.POST("/login", loginUser(authService))
	}
}

func registerUser(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			utils.SendErrorResponse(c, shared_errors.ErrInvalidRequestBody)
			return
		}

		if err := authService.RegisterUser(&user); err != nil {
			utils.SendErrorResponse(c, err)
			return
		}

		utils.SendSuccessResponse(c, http.StatusCreated, gin.H{"message": "User registered successfully"})
	}
}

func loginUser(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginDetails struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&loginDetails); err != nil {
			utils.SendErrorResponse(c, shared_errors.ErrInvalidRequestBody)
			return
		}

		token, err := authService.LoginUser(loginDetails.Email, loginDetails.Password)
		if err != nil {
			utils.SendErrorResponse(c, err)
			return
		}

		utils.SendSuccessResponse(c, http.StatusOK, gin.H{"token": token})
	}
}
