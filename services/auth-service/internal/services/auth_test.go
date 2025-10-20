package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/gmsas95/blytz-mvp/services/auth-service/internal/config"
	"github.com/gmsas95/blytz-mvp/services/auth-service/internal/models"
)

func TestAuthService(t *testing.T) {
	// Create test database
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	assert.NoError(t, err)
	db.AutoMigrate(&models.User{})

	// Create test configuration
	cfg := &config.Config{
		JWTSecret: "test-jwt-secret",
	}

	// Create auth service
	authService := NewAuthService(db, cfg)

	ctx := context.Background()

	t.Run("RegisterUser", func(t *testing.T) {
		user := &models.User{
			Email:       "test@example.com",
			Password:    "password123",
			DisplayName: "Test User",
		}

		err := authService.RegisterUser(user)
		assert.NoError(t, err)

		// Verify user was created by getting it from the database
		createdUser, err := authService.GetUserByEmail(user.Email)
		assert.NoError(t, err)
		assert.NotNil(t, createdUser)
		assert.Equal(t, user.Email, createdUser.Email)
		assert.Equal(t, user.DisplayName, createdUser.DisplayName)
		assert.Equal(t, "user", createdUser.Role)
	})

	t.Run("RegisterUserDuplicateEmail", func(t *testing.T) {
		user := &models.User{
			Email:       "test@example.com",
			Password:    "password123",
			DisplayName: "Test User 2",
		}

		err := authService.RegisterUser(user)
		assert.Error(t, err)
	})

	t.Run("LoginUser", func(t *testing.T) {
		// Now try to login
		token, err := authService.LoginUser("test@example.com", "password123")
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("LoginUserInvalidPassword", func(t *testing.T) {
		_, err := authService.LoginUser("test@example.com", "wrongpassword")
		assert.Error(t, err)
	})

	t.Run("ValidateToken", func(t *testing.T) {
		// First authenticate to get a token
		token, err := authService.LoginUser("test@example.com", "password123")
		assert.NoError(t, err)

		// Validate the token
		claims, err := authService.ValidateToken(token)
		assert.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, "test@example.com", claims.Email)
	})

	t.Run("ValidateTokenInvalid", func(t *testing.T) {
		invalidToken := "invalid.token.here"

		_, err := authService.ValidateToken(invalidToken)
		assert.Error(t, err)
	})
}