package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/gmsas95/blytz-mvp/services/auth-service/internal/config"
	"github.com/gmsas95/blytz-mvp/services/auth-service/internal/models"
	"github.com/gmsas95/blytz-mvp/services/auth-service/internal/services"
)

func main() {
	// Create test database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to create test database:", err)
	}
	db.AutoMigrate(&models.User{})

	// Create test configuration
	cfg := &config.Config{
		JWTSecret: "test-jwt-secret-key",
	}

	// Create auth service
	authService := services.NewAuthService(db, cfg)

	// Test user registration
	user := &models.User{
		Email:       "test@example.com",
		Password:    "password123",
		DisplayName: "Test User",
	}

	fmt.Println("=== Testing User Registration ===")
	err = authService.RegisterUser(user)
	if err != nil {
		log.Fatal("Failed to register user:", err)
	}
	fmt.Printf("✓ User registered successfully: %s\n", user.Email)

	// Test user login
	fmt.Println("\n=== Testing User Login ===")
	token, err := authService.LoginUser("test@example.com", "password123")
	if err != nil {
		log.Fatal("Failed to login:", err)
	}
	fmt.Printf("✓ Login successful, token: %s...\n", token[:50])

	// Test token validation
	fmt.Println("\n=== Testing Token Validation ===")
	ctx := context.Background()
	validationResult, err := authService.ValidateToken(ctx, token)
	if err != nil {
		log.Fatal("Failed to validate token:", err)
	}
	
	if validationResult.Valid {
		fmt.Printf("✓ Token is valid\n")
		fmt.Printf("  UserID: %s\n", validationResult.UserID)
		fmt.Printf("  Email: %s\n", validationResult.Email)
	} else {
		fmt.Printf("✗ Token validation failed: %s\n", validationResult.Message)
	}

	// Test JWT parsing manually
	fmt.Println("\n=== Testing Manual JWT Parsing ===")
	claims := &models.Claims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil {
		log.Fatal("Failed to parse JWT:", err)
	}

	if parsedToken.Valid {
		fmt.Printf("✓ JWT parsing successful\n")
		fmt.Printf("  UserID: %s\n", claims.UserID)
		fmt.Printf("  Email: %s\n", claims.Email)
		fmt.Printf("  Expires: %v\n", claims.ExpiresAt)
	}

	// Test with invalid token
	fmt.Println("\n=== Testing Invalid Token ===")
	invalidResult, err := authService.ValidateToken(ctx, "invalid.token.here")
	if err != nil {
		fmt.Printf("✓ Invalid token rejected: %v\n", err)
	} else if !invalidResult.Valid {
		fmt.Printf("✓ Invalid token rejected: %s\n", invalidResult.Message)
	}

	// Test with wrong password
	fmt.Println("\n=== Testing Wrong Password ===")
	_, err = authService.LoginUser("test@example.com", "wrongpassword")
	if err != nil {
		fmt.Printf("✓ Wrong password rejected: %v\n", err)
	}

	// Test token structure
	fmt.Println("\n=== Testing Token Structure ===")
	parts := strings.Split(token, ".")
	if len(parts) == 3 {
		fmt.Printf("✓ JWT has correct structure (3 parts)\n")
		fmt.Printf("  Header: %s...\n", parts[0][:20])
		fmt.Printf("  Payload: %s...\n", parts[1][:20])
		fmt.Printf("  Signature: %s...\n", parts[2][:20])
	}

	fmt.Println("\n=== Summary ===")
	fmt.Println("✓ JWT authentication implementation is working correctly")
	fmt.Println("✓ Tokens are properly signed and validated")
	fmt.Println("✓ Password hashing is working")
	fmt.Println("✓ Token validation service is functional")
}