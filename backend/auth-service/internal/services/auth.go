package services

import (
	"context"
	"time"

	"github.com/blytz/auth-service/internal/config"
	"github.com/blytz/auth-service/internal/models"
	"github.com/blytz/shared/utils"
	"go.uber.org/zap"
)

type AuthService struct {
	logger *zap.Logger
	config *config.Config
}

func NewAuthService(logger *zap.Logger, config *config.Config) *AuthService {
	return &AuthService{
		logger: logger,
		config: config,
	}
}

// SignUp creates a new user account
func (s *AuthService) SignUp(ctx context.Context, req *models.AuthRequest) (*models.User, error) {
	s.logger.Info("SignUp called", zap.String("email", req.Email))

	// Create user (simplified for now)
	user := &models.User{
		UserID:    generateUserID(),
		Email:     req.Email,
		Username:  req.Username,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      req.Role,
		IsActive:  true,
		IsVerified: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return user, nil
}

// Login authenticates a user
func (s *AuthService) Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error) {
	s.logger.Info("Login called", zap.String("email", req.Email))

	// Simplified login logic
	user := &models.User{
		UserID:    "user123",
		Email:     req.Email,
		Username:  "testuser",
		FirstName: "Test",
		LastName:  "User",
		Role:      "buyer",
		IsActive:  true,
		IsVerified: true,
	}

	// Generate JWT token
	accessToken, err := utils.GenerateJWT(user.UserID, user.Email, user.Role, s.config.JWTSecret, s.config.JWTExpiry)
	if err != nil {
		return nil, err
	}

	refreshToken, err := utils.GenerateJWT(user.UserID, user.Email, user.Role, s.config.JWTSecret, 7*24*time.Hour)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		User:         *user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.config.JWTExpiry.Seconds()),
	}, nil
}

// GetUserByID retrieves a user by ID
func (s *AuthService) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	s.logger.Info("GetUserByID called", zap.String("userID", userID))

	// Simplified user retrieval
	return &models.User{
		UserID:    userID,
		Email:     "user@example.com",
		Username:  "testuser",
		FirstName: "Test",
		LastName:  "User",
		Role:      "buyer",
		IsActive:  true,
		IsVerified: true,
	}, nil
}

// Logout handles user logout
func (s *AuthService) Logout(ctx context.Context, userID string) error {
	s.logger.Info("Logout called", zap.String("userID", userID))
	// Simplified logout - in real implementation, you'd invalidate tokens
	return nil
}

// RefreshToken handles token refresh
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*models.AuthResponse, error) {
	s.logger.Info("RefreshToken called")

	// Validate refresh token and generate new access token
	claims, err := utils.ValidateJWT(refreshToken, s.config.JWTSecret)
	if err != nil {
		return nil, err
	}

	// Generate new access token
	accessToken, err := utils.GenerateJWT(claims.UserID, claims.Email, claims.Role, s.config.JWTSecret, s.config.JWTExpiry)
	if err != nil {
		return nil, err
	}

	// Get user info
	user, err := s.GetUserByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		User:         *user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.config.JWTExpiry.Seconds()),
	}, nil
}

// UpdateProfile updates user profile
func (s *AuthService) UpdateProfile(ctx context.Context, userID string, req *models.UpdateProfileRequest) (*models.User, error) {
	s.logger.Info("UpdateProfile called", zap.String("userID", userID))

	// Get existing user
	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}

	user.UpdatedAt = time.Now()

	return user, nil
}

// Helper function to generate user ID
func generateUserID() string {
	return "user_" + time.Now().Format("20060102150405") + "_" + utils.GenerateRandomString(8)
}