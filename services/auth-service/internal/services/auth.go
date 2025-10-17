package services

import (
	"context"
	"fmt"
	"time"

	"github.com/blytz/auth-service/internal/config"
	"github.com/blytz/auth-service/internal/models"
	"github.com/blytz/auth-service/pkg/betterauth"
	"go.uber.org/zap"
)

// AuthService handles authentication business logic
type AuthService struct {
	config        *config.Config
	betterAuth    *betterauth.Client
	db            *Database  // Database interface
	logger        *zap.Logger
}

// NewAuthService creates a new authentication service
func NewAuthService(cfg *config.Config, logger *zap.Logger) (*AuthService, error) {
	// Initialize Better Auth client
	betterAuthClient, err := betterauth.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Better Auth client: %w", err)
	}

	// Initialize database
	db, err := NewDatabase(cfg.GetDatabaseURL(), logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	return &AuthService{
		config:     cfg,
		betterAuth: betterAuthClient,
		db:         db,
		logger:     logger,
	}, nil
}

// CreateUser creates a new user with Better Auth
func (s *AuthService) CreateUser(ctx context.Context, req *models.RegisterRequest) (*models.User, error) {
	s.logger.Info("Creating new user", zap.String("email", req.Email))

	// Use Better Auth to create user
	betterUser, err := s.betterAuth.CreateUser(ctx, req.Email, req.Password, req.DisplayName)
	if err != nil {
		s.logger.Error("Failed to create user with Better Auth", zap.Error(err))
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Create user in our database
	user := &models.User{
		ID:          betterUser.ID,
		Email:       betterUser.Email,
		DisplayName: betterUser.DisplayName,
		PhoneNumber: req.PhoneNumber,
		Role:        "user",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.db.CreateUser(ctx, user); err != nil {
		s.logger.Error("Failed to create user in database", zap.Error(err))
		return nil, fmt.Errorf("failed to create user in database: %w", err)
	}

	s.logger.Info("User created successfully", zap.String("user_id", user.ID))
	return user, nil
}

// AuthenticateUser authenticates a user with Better Auth
func (s *AuthService) AuthenticateUser(ctx context.Context, req *models.LoginRequest) (*models.User, string, error) {
	s.logger.Info("Authenticating user", zap.String("email", req.Email))

	// Use Better Auth to authenticate
	betterUser, token, err := s.betterAuth.AuthenticateUser(ctx, req.Email, req.Password)
	if err != nil {
		s.logger.Error("Authentication failed", zap.Error(err))
		return nil, "", fmt.Errorf("authentication failed: %w", err)
	}

	// Get user from database
	user, err := s.db.GetUserByID(ctx, betterUser.ID)
	if err != nil {
		s.logger.Error("Failed to get user from database", zap.Error(err))
		return nil, "", fmt.Errorf("failed to get user: %w", err)
	}

	s.logger.Info("User authenticated successfully", zap.String("user_id", user.ID))
	return user, token, nil
}

// ValidateToken validates a JWT token
func (s *AuthService) ValidateToken(ctx context.Context, token string) (*models.ValidateTokenResponse, error) {
	// Use Better Auth to validate token
	userID, email, err := s.betterAuth.ValidateToken(ctx, token)
	if err != nil {
		return &models.ValidateTokenResponse{
			Valid:   false,
			Message: "Invalid or expired token",
		}, nil
	}

	return &models.ValidateTokenResponse{
		Valid:  true,
		UserID: userID,
		Email:  email,
	}, nil
}

// RefreshToken refreshes a JWT token
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	// Use Better Auth to refresh token
	newToken, err := s.betterAuth.RefreshToken(ctx, refreshToken)
	if err != nil {
		s.logger.Error("Failed to refresh token", zap.Error(err))
		return "", fmt.Errorf("failed to refresh token: %w", err)
	}

	return newToken, nil
}

// UpdateProfile updates user profile
func (s *AuthService) UpdateProfile(ctx context.Context, userID string, req *models.UpdateProfileRequest) error {
	s.logger.Info("Updating user profile", zap.String("user_id", userID))

	// Update in Better Auth
	if err := s.betterAuth.UpdateProfile(ctx, userID, req); err != nil {
		s.logger.Error("Failed to update profile in Better Auth", zap.Error(err))
		return fmt.Errorf("failed to update profile: %w", err)
	}

	// Update in database
	if err := s.db.UpdateUserProfile(ctx, userID, req); err != nil {
		s.logger.Error("Failed to update profile in database", zap.Error(err))
		return fmt.Errorf("failed to update profile in database: %w", err)
	}

	s.logger.Info("Profile updated successfully", zap.String("user_id", userID))
	return nil
}

// ChangePassword changes user password
func (s *AuthService) ChangePassword(ctx context.Context, userID string, req *models.ChangePasswordRequest) error {
	s.logger.Info("Changing user password", zap.String("user_id", userID))

	// Use Better Auth to change password
	if err := s.betterAuth.ChangePassword(ctx, userID, req.CurrentPassword, req.NewPassword); err != nil {
		s.logger.Error("Failed to change password", zap.Error(err))
		return fmt.Errorf("failed to change password: %w", err)
	}

	s.logger.Info("Password changed successfully", zap.String("user_id", userID))
	return nil
}

// GetUserByID gets user by ID
func (s *AuthService) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	return s.db.GetUserByID(ctx, userID)
}

// GetUserByEmail gets user by email
func (s *AuthService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.db.GetUserByEmail(ctx, email)
}
