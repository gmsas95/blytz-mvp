package services

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/gmsas95/blytz-mvp/services/auth-service/internal/config"
	"github.com/gmsas95/blytz-mvp/services/auth-service/internal/models"
)

// AuthService provides authentication related services

type AuthService struct {
	db     *gorm.DB
	config *config.Config
}

// GetConfig returns the service configuration
func (s *AuthService) GetConfig() *config.Config {
	return s.config
}

// NewAuthService creates a new AuthService
func NewAuthService(db *gorm.DB, config *config.Config) *AuthService {
	return &AuthService{db: db, config: config}
}

// RegisterUser registers a new user
func (s *AuthService) RegisterUser(user *models.User) error {
	// Check if user already exists
	if s.userExists(user.Email) {
		return errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	// Create user
	return s.db.Create(user).Error
}

// LoginUser logs in a user and returns a JWT token
func (s *AuthService) LoginUser(email, password string) (string, error) {
	// Get user by email
	user, err := s.GetUserByEmail(email)
	if err != nil {
		return "", err
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := s.generateJWT(user)
	if err != nil {
		return "", err
	}

	return token, nil
}

// GetUserByEmail gets a user by email
func (s *AuthService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// ValidateToken validates a JWT token and returns the claims
func (s *AuthService) ValidateToken(ctx context.Context, tokenString string) (*models.ValidateTokenResponse, error) {
	claims := &models.Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil {
		return &models.ValidateTokenResponse{
			Valid:   false,
			Message: "Invalid token: " + err.Error(),
		}, nil
	}

	if !token.Valid {
		return &models.ValidateTokenResponse{
			Valid:   false,
			Message: "Token is not valid",
		}, nil
	}

	return &models.ValidateTokenResponse{
		Valid:   true,
		UserID:  claims.UserID,
		Email:   claims.Email,
		Message: "Token is valid",
	}, nil
}

// userExists checks if a user exists by email
func (s *AuthService) userExists(email string) bool {
	var user models.User
	return s.db.Where("email = ?", email).First(&user).Error == nil
}

// generateJWT generates a JWT token for a user
func (s *AuthService) generateJWT(user *models.User) (string, error) {
	claims := &models.Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWTSecret))
}

