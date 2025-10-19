package betterauth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/gmsas95/blytz-mvp/services/auth-service/internal/config"
	"github.com/gmsas95/blytz-mvp/services/auth-service/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"go.uber.org/zap"
)

// BetterUser represents a user in Better Auth system
type BetterUser struct {
	ID          string    `json:"id"`
	Email       string    `json:"email"`
	DisplayName string    `json:"display_name"`
	Password    string    `json:"-"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Client handles Better Auth operations
type Client struct {
	config *config.Config
	logger *zap.Logger
	users  map[string]*BetterUser // In-memory storage for MVP
}

// NewClient creates a new Better Auth client
func NewClient(cfg *config.Config) (*Client, error) {
	return &Client{
		config: cfg,
		logger: zap.NewNop(), // Will be set by service
		users:  make(map[string]*BetterUser),
	}, nil
}

// SetLogger sets the logger for the client
func (c *Client) SetLogger(logger *zap.Logger) {
	c.logger = logger
}

// CreateUser creates a new user with Better Auth
func (c *Client) CreateUser(ctx context.Context, email, password, displayName string) (*BetterUser, error) {
	c.logger.Info("Creating user with Better Auth", zap.String("email", email))

	// Check if user already exists
	if _, exists := c.findUserByEmail(email); exists {
		return nil, fmt.Errorf("user with email %s already exists", email)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Generate user ID
	userID, err := c.generateUserID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate user ID: %w", err)
	}

	// Create user
	user := &BetterUser{
		ID:          userID,
		Email:       email,
		DisplayName: displayName,
		Password:    string(hashedPassword),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Store user
	c.users[userID] = user

	c.logger.Info("User created successfully", zap.String("user_id", userID))
	return user, nil
}

// AuthenticateUser authenticates a user with Better Auth
func (c *Client) AuthenticateUser(ctx context.Context, email, password string) (*BetterUser, error) {
	c.logger.Info("Authenticating user with Better Auth", zap.String("email", email))

	// Find user by email
	user, exists := c.findUserByEmail(email)
	if !exists {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	c.logger.Info("User authenticated successfully", zap.String("user_id", user.ID))
	return user, nil
}

// UpdateProfile updates user profile information
func (c *Client) UpdateProfile(ctx context.Context, userID string, req *models.UpdateProfileRequest) error {
	c.logger.Info("Updating user profile", zap.String("user_id", userID))

	user, exists := c.users[userID]
	if !exists {
		return fmt.Errorf("user not found")
	}

	// Update fields if provided
	if req.DisplayName != "" {
		user.DisplayName = req.DisplayName
	}

	user.UpdatedAt = time.Now()

	c.logger.Info("Profile updated successfully", zap.String("user_id", userID))
	return nil
}

// ChangePassword changes user password
func (c *Client) ChangePassword(ctx context.Context, userID, currentPassword, newPassword string) error {
	c.logger.Info("Changing user password", zap.String("user_id", userID))

	user, exists := c.users[userID]
	if !exists {
		return fmt.Errorf("user not found")
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword)); err != nil {
		return fmt.Errorf("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	// Update password
	user.Password = string(hashedPassword)
	user.UpdatedAt = time.Now()

	c.logger.Info("Password changed successfully", zap.String("user_id", userID))
	return nil
}

// GenerateJWT generates a JWT token for the user
func (c *Client) GenerateJWT(userID, email string) (string, error) {
	c.logger.Info("Generating JWT token", zap.String("user_id", userID))

	// Create claims
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // 24 hours
		"iat":     time.Now().Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token
	tokenString, err := token.SignedString([]byte(c.config.GetJWTSecret()))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	c.logger.Info("JWT token generated successfully", zap.String("user_id", userID))
	return tokenString, nil
}

// ValidateToken validates a JWT token
func (c *Client) ValidateToken(ctx context.Context, tokenString string) (string, string, error) {
	c.logger.Info("Validating JWT token")

	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(c.config.GetJWTSecret()), nil
	})

	if err != nil {
		return "", "", fmt.Errorf("invalid token: %w", err)
	}

	// Extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok1 := claims["user_id"].(string)
		email, ok2 := claims["email"].(string)

		if !ok1 || !ok2 {
			return "", "", fmt.Errorf("invalid token claims")
		}

		c.logger.Info("Token validated successfully", zap.String("user_id", userID))
		return userID, email, nil
	}

	return "", "", fmt.Errorf("invalid token")
}

// RefreshToken refreshes a JWT token
func (c *Client) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	c.logger.Info("Refreshing JWT token")

	// For MVP, we'll just validate the existing token and generate a new one
	// In production, you'd have a separate refresh token mechanism
	userID, email, err := c.ValidateToken(ctx, refreshToken)
	if err != nil {
		return "", fmt.Errorf("invalid refresh token: %w", err)
	}

	// Generate new token
	newToken, err := c.GenerateJWT(userID, email)
	if err != nil {
		return "", fmt.Errorf("failed to generate new token: %w", err)
	}

	c.logger.Info("Token refreshed successfully", zap.String("user_id", userID))
	return newToken, nil
}

// Helper methods

func (c *Client) findUserByEmail(email string) (*BetterUser, bool) {
	for _, user := range c.users {
		if user.Email == email {
			return user, true
		}
	}
	return nil, false
}

func (c *Client) generateUserID() (string, error) {
	// Generate random bytes
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	// Encode to base64
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// GetUserByID gets user by ID (for internal use)
func (c *Client) GetUserByID(ctx context.Context, userID string) (*BetterUser, error) {
	user, exists := c.users[userID]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}