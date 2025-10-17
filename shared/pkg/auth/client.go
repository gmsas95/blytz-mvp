package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/blytz/shared/errors"
	"github.com/gin-gonic/gin"
)

// AuthClient provides authentication client for microservices
type AuthClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewAuthClient creates a new authentication client
func NewAuthClient(baseURL string) *AuthClient {
	return &AuthClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// ValidateToken validates a JWT token with the auth service
func (c *AuthClient) ValidateToken(ctx context.Context, token string) (*ValidateTokenResponse, error) {
	// Prepare request
	reqBody := ValidateTokenRequest{
		Token: token,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/v1/auth/validate", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		var errorResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return nil, fmt.Errorf("authentication service error (status %d)", resp.StatusCode)
		}
		return nil, errors.AuthenticationError("AUTH_SERVICE_ERROR", errorResp.Message)
	}

	// Decode response
	var result ValidateTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// GetUserInfo gets user information from the auth service
func (c *AuthClient) GetUserInfo(ctx context.Context, token string) (*UserInfo, error) {
	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/api/v1/auth/me", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		var errorResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return nil, fmt.Errorf("authentication service error (status %d)", resp.StatusCode)
		}
		return nil, errors.AuthenticationError("AUTH_SERVICE_ERROR", errorResp.Message)
	}

	// Decode response
	var result struct {
		Success bool     `json:"success"`
		Message string   `json:"message"`
		Data    UserInfo `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result.Data, nil
}

// RefreshToken refreshes a JWT token
func (c *AuthClient) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	// Prepare request
	reqBody := RefreshTokenRequest{
		RefreshToken: refreshToken,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/v1/auth/refresh", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		var errorResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return "", fmt.Errorf("authentication service error (status %d)", resp.StatusCode)
		}
		return "", errors.AuthenticationError("AUTH_SERVICE_ERROR", errorResp.Message)
	}

	// Decode response
	var result struct {
		Success bool `json:"success"`
		Message string `json:"message"`
		Data    struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Data.Token, nil
}

// Request models

type ValidateTokenRequest struct {
	Token string `json:"token"`
}

type ValidateTokenResponse struct {
	Valid   bool   `json:"valid"`
	UserID  string `json:"user_id,omitempty"`
	Email   string `json:"email,omitempty"`
	Message string `json:"message,omitempty"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type UserInfo struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	PhoneNumber string `json:"phone_number,omitempty"`
	AvatarURL   string `json:"avatar_url,omitempty"`
	Role        string `json:"role"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

// Middleware helper for other services

// AuthMiddleware creates authentication middleware for other microservices
func AuthMiddleware(authClient *AuthClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "No authorization token provided", http.StatusUnauthorized)
				return
			}

			// Expected format: "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
				return
			}

			token := parts[1]

			// Validate token
			response, err := authClient.ValidateToken(r.Context(), token)
			if err != nil || !response.Valid {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			// Set user context in request
			ctx := context.WithValue(r.Context(), "userID", response.UserID)
			ctx = context.WithValue(ctx, "userEmail", response.Email)

			// Continue with authenticated request
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GinAuthMiddleware creates authentication middleware for Gin framework
func GinAuthMiddleware(authClient *AuthClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "No authorization token provided",
			})
			return
		}

		// Expected format: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization format",
			})
			return
		}

		token := parts[1]

		// Validate token
		response, err := authClient.ValidateToken(c.Request.Context(), token)
		if err != nil || !response.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			return
		}

		// Set user context
		c.Set("userID", response.UserID)
		c.Set("userEmail", response.Email)

		c.Next()
	}
}

// OptionalGinAuthMiddleware creates optional authentication middleware for Gin
func OptionalGinAuthMiddleware(authClient *AuthClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// No token provided, continue without authentication
			c.Set("isAuthenticated", false)
			c.Next()
			return
		}

		// Expected format: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			// Invalid format, mark as unauthenticated
			c.Set("isAuthenticated", false)
			c.Next()
			return
		}

		token := parts[1]

		// Try to validate token
		response, err := authClient.ValidateToken(c.Request.Context(), token)
		if err == nil && response.Valid {
			// Token is valid, set user context
			c.Set("userID", response.UserID)
			c.Set("userEmail", response.Email)
			c.Set("isAuthenticated", true)
		} else {
			// Token is invalid, mark as unauthenticated
			c.Set("isAuthenticated", false)
		}

		c.Next()
	}
}

// Context helpers

// GetUserID extracts user ID from context
func GetUserID(ctx context.Context) string {
	if userID, ok := ctx.Value("userID").(string); ok {
		return userID
	}
	return ""
}

// GetUserEmail extracts user email from context
func GetUserEmail(ctx context.Context) string {
	if userEmail, ok := ctx.Value("userEmail").(string); ok {
		return userEmail
	}
	return ""
}

// ContextWithUser creates a new context with user information
func ContextWithUser(ctx context.Context, userID, userEmail string) context.Context {
	ctx = context.WithValue(ctx, "userID", userID)
	ctx = context.WithValue(ctx, "userEmail", userEmail)
	return ctx
}