package services

import (
	"context"
	"testing"
	"time"

	"github.com/blytz/auth-service/internal/config"
	"github.com/blytz/auth-service/internal/models"
	"go.uber.org/zap"
)

func TestAuthService(t *testing.T) {
	// Create test configuration
	cfg := &config.Config{
		ServicePort:      "8084",
		Environment:      "test",
		DatabaseURL:      "memory://test",
		BetterAuthSecret: "test-secret",
		JWTSecret:        "test-jwt-secret",
	}

	// Create logger
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	// Create auth service
	authService, err := NewAuthService(cfg, logger)
	if err != nil {
		t.Fatalf("Failed to create auth service: %v", err)
	}

	ctx := context.Background()

	t.Run("CreateUser", func(t *testing.T) {
		req := &models.RegisterRequest{
			Email:       "test@example.com",
			Password:    "password123",
			DisplayName: "Test User",
			PhoneNumber: "+1234567890",
		}

		user, err := authService.CreateUser(ctx, req)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}

		if user.Email != req.Email {
			t.Errorf("Expected email %s, got %s", req.Email, user.Email)
		}

		if user.DisplayName != req.DisplayName {
			t.Errorf("Expected display name %s, got %s", req.DisplayName, user.DisplayName)
		}

		if user.Role != "user" {
			t.Errorf("Expected role 'user', got %s", user.Role)
		}
	})

	t.Run("CreateUserDuplicateEmail", func(t *testing.T) {
		req := &models.RegisterRequest{
			Email:       "test@example.com",
			Password:    "password123",
			DisplayName: "Test User 2",
		}

		_, err := authService.CreateUser(ctx, req)
		if err == nil {
			t.Error("Expected error for duplicate email, got nil")
		}
	})

	t.Run("AuthenticateUser", func(t *testing.T) {
		req := &models.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}

		user, token, err := authService.AuthenticateUser(ctx, req)
		if err != nil {
			t.Fatalf("Failed to authenticate user: %v", err)
		}

		if user.Email != req.Email {
			t.Errorf("Expected email %s, got %s", req.Email, user.Email)
		}

		if token == "" {
			t.Error("Expected non-empty token")
		}
	})

	t.Run("AuthenticateUserInvalidPassword", func(t *testing.T) {
		req := &models.LoginRequest{
			Email:    "test@example.com",
			Password: "wrongpassword",
		}

		_, _, err := authService.AuthenticateUser(ctx, req)
		if err == nil {
			t.Error("Expected error for invalid password, got nil")
		}
	})

	t.Run("ValidateToken", func(t *testing.T) {
		// First authenticate to get a token
		loginReq := &models.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}

		_, token, err := authService.AuthenticateUser(ctx, loginReq)
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

	t.Run("RefreshToken", func(t *testing.T) {
		// First authenticate to get a token
		loginReq := &models.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}

		_, token, err := authService.AuthenticateUser(ctx, loginReq)
		if err != nil {
			t.Fatalf("Failed to authenticate user: %v", err)
		}

		// Refresh the token
		newToken, err := authService.RefreshToken(ctx, token)
		if err != nil {
			t.Fatalf("Failed to refresh token: %v", err)
		}

		if newToken == "" {
			t.Error("Expected non-empty new token")
		}

		if newToken == token {
			t.Error("Expected different token after refresh")
		}
	})

	t.Run("UpdateProfile", func(t *testing.T) {
		userID := "user123" // This would be the actual user ID from authentication

		req := &models.UpdateProfileRequest{
			DisplayName: "Updated Name",
			PhoneNumber: "+0987654321",
			AvatarURL:   "https://example.com/avatar.jpg",
		}

		err := authService.UpdateProfile(ctx, userID, req)
		if err != nil {
			t.Fatalf("Failed to update profile: %v", err)
		}
	})

	t.Run("ChangePassword", func(t *testing.T) {
		userID := "user123" // This would be the actual user ID from authentication

		req := &models.ChangePasswordRequest{
			CurrentPassword: "password123",
			NewPassword:     "newpassword123",
		}

		err := authService.ChangePassword(ctx, userID, req)
		if err != nil {
			t.Fatalf("Failed to change password: %v", err)
		}
	})

	t.Run("ChangePasswordWrongCurrent", func(t *testing.T) {
		userID := "user123" // This would be the actual user ID from authentication

		req := &models.ChangePasswordRequest{
			CurrentPassword: "wrongpassword",
			NewPassword:     "newpassword123",
		}

		err := authService.ChangePassword(ctx, userID, req)
		if err == nil {
			t.Error("Expected error for wrong current password, got nil")
		}
	})

	t.Run("GetUserByID", func(t *testing.T) {
		userID := "user123" // This would be the actual user ID

		user, err := authService.GetUserByID(ctx, userID)
		if err != nil {
			t.Fatalf("Failed to get user by ID: %v", err)
		}

		if user.ID != userID {
			t.Errorf("Expected user ID %s, got %s", userID, user.ID)
		}
	})

	t.Run("GetUserByEmail", func(t *testing.T) {
		email := "test@example.com"

		user, err := authService.GetUserByEmail(ctx, email)
		if err != nil {
			t.Fatalf("Failed to get user by email: %v", err)
		}

		if user.Email != email {
			t.Errorf("Expected email %s, got %s", email, user.Email)
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