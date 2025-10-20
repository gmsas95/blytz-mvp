package services

import (
	"context"
	"testing"
	"time"

	"github.com/gmsas95/blytz-mvp/services/auth-service/internal/config"
	"github.com/gmsas95/blytz-mvp/services/auth-service/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"go.uber.org/zap"
)

func TestAuthService(t *testing.T) {
	// Create test database
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Auto migrate the User model
	db.AutoMigrate(&models.User{})

	// Create test configuration
	cfg := &config.Config{
		ServicePort:      "8084",
		Environment:      "test",
		DatabaseURL:      "memory://test",
		BetterAuthSecret: "test-secret",
		JWTSecret:        "test-jwt-secret",
	}

	// Create auth service
	authService := NewAuthService(db, cfg)

	ctx := context.Background()

	t.Run("CreateUser", func(t *testing.T) {
		user := &models.User{
			Email:       "test@example.com",
			DisplayName: "Test User",
			Password:    "password123",
			PhoneNumber: "+1234567890",
		}

		err := authService.RegisterUser(user)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}

		// Verify user was created by getting it from the database
		createdUser, err := authService.GetUserByEmail(user.Email)
		if err != nil {
			t.Fatalf("Failed to get created user: %v", err)
		}

		if createdUser.Email != user.Email {
			t.Errorf("Expected email %s, got %s", user.Email, createdUser.Email)
		}

		if createdUser.DisplayName != user.DisplayName {
			t.Errorf("Expected display name %s, got %s", user.DisplayName, createdUser.DisplayName)
		}

		if createdUser.Role != "user" {
			t.Errorf("Expected role 'user', got %s", createdUser.Role)
		}
	})

	t.Run("CreateUserDuplicateEmail", func(t *testing.T) {
		user := &models.User{
			Email:       "test@example.com",
			DisplayName: "Test User 2",
			Password:    "password123",
		}

		err := authService.RegisterUser(user)
		if err == nil {
			t.Error("Expected error for duplicate email, got nil")
		}
	})

	t.Run("AuthenticateUser", func(t *testing.T) {
		// First register a user to authenticate
		user := &models.User{
			Email:       "test@example.com",
			DisplayName: "Test User",
			Password:    "password123",
		}
		err := authService.RegisterUser(user)
		if err != nil {
			t.Fatalf("Failed to register user: %v", err)
		}

		// Now try to login
		token, err := authService.LoginUser("test@example.com", "password123")
		if err != nil {
			t.Fatalf("Failed to authenticate user: %v", err)
		}

		if token == "" {
			t.Error("Expected non-empty token")
		}
	})

	t.Run("AuthenticateUserInvalidPassword", func(t *testing.T) {
		_, err := authService.LoginUser("test@example.com", "wrongpassword")
		if err == nil {
			t.Error("Expected error for invalid password, got nil")
		}
	})

	t.Run("ValidateToken", func(t *testing.T) {
		// First authenticate to get a token
		token, err := authService.LoginUser("test@example.com", "password123")
		if err != nil {
			t.Fatalf("Failed to authenticate user: %v", err)
		}

		// Validate the token
		response, err := authService.ValidateToken(ctx, token)
		if err != nil {
			t.Fatalf("Failed to validate token: %v", err)
		}

		if !response.Valid {
			t.Error("Expected token to be valid")
		}

		if response.Email != "test@example.com" {
			t.Errorf("Expected email test@example.com, got %s", response.Email)
		}
	})

	t.Run("ValidateTokenInvalid", func(t *testing.T) {
		invalidToken := "invalid.token.here"

		response, err := authService.ValidateToken(ctx, invalidToken)
		if err != nil {
			t.Fatalf("Failed to validate token: %v", err)
		}

		if response.Valid {
			t.Error("Expected token to be invalid")
		}
	})

}

func TestDatabase(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	db, err := NewDatabase("memory://test", logger)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	t.Run("HealthCheck", func(t *testing.T) {
		err := db.Health(ctx)
		if err != nil {
			t.Errorf("Expected healthy database, got error: %v", err)
		}
	})

	t.Run("CreateAndGetUser", func(t *testing.T) {
		user := &models.User{
			ID:          "test-user-123",
			Email:       "test@example.com",
			DisplayName: "Test User",
			PhoneNumber: "+1234567890",
			Role:        "user",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		// Create user
		err := db.CreateUser(ctx, user)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}

		// Get user by ID
		retrievedUser, err := db.GetUserByID(ctx, user.ID)
		if err != nil {
			t.Fatalf("Failed to get user by ID: %v", err)
		}

		if retrievedUser.Email != user.Email {
			t.Errorf("Expected email %s, got %s", user.Email, retrievedUser.Email)
		}

		// Get user by email
		retrievedUser, err = db.GetUserByEmail(ctx, user.Email)
		if err != nil {
			t.Fatalf("Failed to get user by email: %v", err)
		}

		if retrievedUser.ID != user.ID {
			t.Errorf("Expected ID %s, got %s", user.ID, retrievedUser.ID)
		}
	})

	t.Run("UpdateUserProfile", func(t *testing.T) {
		userID := "test-user-123"
		req := &models.UpdateProfileRequest{
			DisplayName: "Updated Name",
			PhoneNumber: "+0987654321",
			AvatarURL:   "https://example.com/avatar.jpg",
		}

		err := db.UpdateUserProfile(ctx, userID, req)
		if err != nil {
			t.Fatalf("Failed to update user profile: %v", err)
		}

		// Verify update
		updatedUser, err := db.GetUserByID(ctx, userID)
		if err != nil {
			t.Fatalf("Failed to get updated user: %v", err)
		}

		if updatedUser.DisplayName != req.DisplayName {
			t.Errorf("Expected display name %s, got %s", req.DisplayName, updatedUser.DisplayName)
		}

		if updatedUser.PhoneNumber != req.PhoneNumber {
			t.Errorf("Expected phone number %s, got %s", req.PhoneNumber, updatedUser.PhoneNumber)
		}

		if updatedUser.AvatarURL != req.AvatarURL {
			t.Errorf("Expected avatar URL %s, got %s", req.AvatarURL, updatedUser.AvatarURL)
		}
	})

	t.Run("DeleteUser", func(t *testing.T) {
		user := &models.User{
			ID:          "delete-test-user",
			Email:       "delete@test.com",
			DisplayName: "Delete Test User",
			Role:        "user",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		// Create user
		err := db.CreateUser(ctx, user)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}

		// Delete user
		err = db.DeleteUser(ctx, user.ID)
		if err != nil {
			t.Fatalf("Failed to delete user: %v", err)
		}

		// Verify deletion
		_, err = db.GetUserByID(ctx, user.ID)
		if err == nil {
			t.Error("Expected error when getting deleted user, got nil")
		}
	})
}