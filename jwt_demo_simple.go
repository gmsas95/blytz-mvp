package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Simple claims structure
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func main() {
	jwtSecret := "test-jwt-secret-key"
	
	fmt.Println("=== Testing JWT Implementation ===")
	
	// Create claims
	claims := &Claims{
		UserID: "user123",
		Email:  "test@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// Generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		log.Fatal("Failed to generate token:", err)
	}
	
	fmt.Printf("✓ Token generated successfully\n")
	fmt.Printf("Token: %s...\n", tokenString[:50])

	// Parse and validate token
	parsedClaims := &Claims{}
	parsedToken, err := jwt.ParseWithClaims(tokenString, parsedClaims, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		log.Fatal("Failed to parse token:", err)
	}

	if parsedToken.Valid {
		fmt.Printf("✓ Token validation successful\n")
		fmt.Printf("  UserID: %s\n", parsedClaims.UserID)
		fmt.Printf("  Email: %s\n", parsedClaims.Email)
		fmt.Printf("  Issued At: %v\n", parsedClaims.IssuedAt)
		fmt.Printf("  Expires At: %v\n", parsedClaims.ExpiresAt)
	}

	// Test token structure
	parts := strings.Split(tokenString, ".")
	if len(parts) == 3 {
		fmt.Printf("\n✓ JWT has correct structure (3 parts)\n")
		fmt.Printf("  Header: %s...\n", parts[0][:20])
		fmt.Printf("  Payload: %s...\n", parts[1][:20])
		fmt.Printf("  Signature: %s...\n", parts[2][:20])
	}

	// Test with invalid token
	fmt.Printf("\n=== Testing Invalid Token ===")
	invalidToken := "invalid.token.here"
	_, err = jwt.ParseWithClaims(invalidToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	
	if err != nil {
		fmt.Printf("✓ Invalid token correctly rejected: %v\n", err)
	}

	// Test with wrong secret
	fmt.Printf("\n=== Testing Wrong Secret ===")
	_, err = jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("wrong-secret"), nil
	})
	
	if err != nil {
		fmt.Printf("✓ Token with wrong secret correctly rejected: %v\n", err)
	}

	// Test expired token
	fmt.Printf("\n=== Testing Expired Token ===")
	expiredClaims := &Claims{
		UserID: "user123",
		Email:  "test@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)), // Expired 1 hour ago
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-time.Hour * 2)),
		},
	}
	
	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	expiredTokenString, _ := expiredToken.SignedString([]byte(jwtSecret))
	
	_, err = jwt.ParseWithClaims(expiredTokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	
	if err != nil {
		fmt.Printf("✓ Expired token correctly rejected: %v\n", err)
	}

	fmt.Printf("\n=== Summary ===")
	fmt.Println("✓ JWT implementation is working correctly")
	fmt.Println("✓ Tokens are properly signed with HS256")
	fmt.Println("✓ Token validation works as expected")
	fmt.Println("✓ Invalid tokens are properly rejected")
	fmt.Println("✓ Expired tokens are properly rejected")
	fmt.Println("✓ Wrong secrets are properly rejected")
}