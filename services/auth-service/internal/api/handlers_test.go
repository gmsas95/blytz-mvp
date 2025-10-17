package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/blytz/auth-service/internal/config"
	"github.com/blytz/auth-service/internal/models"
	"github.com/blytz/auth-service/internal/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func setupTestRouter() (*gin.Engine, *AuthHandler, *services.AuthService) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

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

	// Create auth service
	authService, _ := services.NewAuthService(cfg, logger)

	// Create auth handler
	authHandler := NewAuthHandler(authService, logger)

	// Create router
	router := gin.New()

	// Setup routes
	setupRoutes(router, authHandler, logger)

	return router, authHandler, authService
}

func setupRoutes(router *gin.Engine, authHandler *AuthHandler, logger *zap.Logger) {
	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"service":   "auth",
			"timestamp": time.Now().Unix(),
		})
	})

	// API v1 routes
	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.RegisterUser)
			auth.POST("/login", authHandler.LoginUser)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/validate", authHandler.ValidateToken)
		}
	}
}

func TestRegisterUser(t *testing.T) {
	router, _, _ := setupTestRouter()

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name: "Valid registration",
			requestBody: models.RegisterRequest{
				Email:       "test@example.com",
				Password:    "password123",
				DisplayName: "Test User",
				PhoneNumber: "+1234567890",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var response models.AuthResponse
				if err := json.Unmarshal(body, &response); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}

				if !response.Success {
					t.Error("Expected success to be true")
				}

				if response.User == nil {
					t.Error("Expected user in response")
				}

				if response.Token == "" {
					t.Error("Expected token in response")
				}
			},
		},
		{
			name: "Invalid email format",
			requestBody: models.RegisterRequest{
				Email:       "invalid-email",
				Password:    "password123",
				DisplayName: "Test User",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body []byte) {
				var response models.ErrorResponse
				if err := json.Unmarshal(body, &response); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}

				if response.Success {
					t.Error("Expected success to be false")
				}
			},
		},
		{
			name: "Missing required fields",
			requestBody: models.RegisterRequest{
				Email: "test@example.com",
				// Missing password and display name
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body []byte) {
				var response models.ErrorResponse
				if err := json.Unmarshal(body, &response); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}

				if response.Success {
					t.Error("Expected success to be false")
				}
			},
		},
		{
			name:           "Invalid JSON body",
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, body []byte) {
				var response models.ErrorResponse
				if err := json.Unmarshal(body, &response); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}

				if response.Success {
					t.Error("Expected success to be false")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal request body
			jsonBody, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}

			// Create request
			req, err := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(jsonBody))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Serve request
			router.ServeHTTP(w, req)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Check response
			tt.checkResponse(t, w.Body.Bytes())
		})
	}
}

func TestLoginUser(t *testing.T) {
	router, _, authService := setupTestRouter()

	// First create a user for testing login
	registerReq := &models.RegisterRequest{
		Email:       "login@example.com",
		Password:    "password123",
		DisplayName: "Login Test User",
		PhoneNumber: "+1234567890",
	}

	_, err := authService.CreateUser(context.Background(), registerReq)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name: "Valid login",
			requestBody: models.LoginRequest{
				Email:    "login@example.com",
				Password: "password123",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var response models.AuthResponse
				if err := json.Unmarshal(body, &response); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}

				if !response.Success {
					t.Error("Expected success to be true")
				}

				if response.User == nil {
					t.Error("Expected user in response")
				}

				if response.Token == "" {
					t.Error("Expected token in response")
				}
			},
		},
		{
			name: "Invalid password",
			requestBody: models.LoginRequest{
				Email:    "login@example.com",
				Password: "wrongpassword",
			},
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, body []byte) {
				var response models.ErrorResponse
				if err := json.Unmarshal(body, &response); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}

				if response.Success {
					t.Error("Expected success to be false")
				}
			},
		},
		{
			name: "Non-existent user",
			requestBody: models.LoginRequest{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			expectedStatus: http.StatusUnauthorized,
			checkResponse: func(t *testing.T, body []byte) {
				var response models.ErrorResponse
				if err := json.Unmarshal(body, &response); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}

				if response.Success {
					t.Error("Expected success to be false")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal request body
			jsonBody, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}

			// Create request
			req, err := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonBody))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Serve request
			router.ServeHTTP(w, req)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Check response
			tt.checkResponse(t, w.Body.Bytes())
		})
	}
}

func TestValidateToken(t *testing.T) {
	router, _, authService := setupTestRouter()

	// First create and authenticate a user to get a valid token
	registerReq := &models.RegisterRequest{
		Email:       "validate@example.com",
		Password:    "password123",
		DisplayName: "Validate Test User",
	}

	_, err := authService.CreateUser(context.Background(), registerReq)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	loginReq := &models.LoginRequest{
		Email:    "validate@example.com",
		Password: "password123",
	}

	_, validToken, err := authService.AuthenticateUser(context.Background(), loginReq)
	if err != nil {
		t.Fatalf("Failed to authenticate user: %v", err)
	}

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		checkResponse  func(t *testing.T, body []byte)
	}{
		{
			name: "Valid token",
			requestBody: models.ValidateTokenRequest{
				Token: validToken,
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var response models.ValidateTokenResponse
				if err := json.Unmarshal(body, &response); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}

				if !response.Valid {
					t.Error("Expected token to be valid")
				}

				if response.UserID == "" {
					t.Error("Expected user ID in response")
				}

				if response.Email == "" {
					t.Error("Expected email in response")
				}
			},
		},
		{
			name: "Invalid token",
			requestBody: models.ValidateTokenRequest{
				Token: "invalid.token.here",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, body []byte) {
				var response models.ValidateTokenResponse
				if err := json.Unmarshal(body, &response); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}

				if response.Valid {
					t.Error("Expected token to be invalid")
				}

				if response.Message == "" {
					t.Error("Expected message in response")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal request body
			jsonBody, err := json.Marshal(tt.requestBody)
			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}

			// Create request
			req, err := http.NewRequest("POST", "/api/v1/auth/validate", bytes.NewBuffer(jsonBody))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Serve request
			router.ServeHTTP(w, req)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Check response
			tt.checkResponse(t, w.Body.Bytes())
		})
	}
}

func TestHealthCheck(t *testing.T) {
	router, _, _ := setupTestRouter()

	// Create request
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create response recorder
	w := httptest.NewRecorder()

	// Serve request
	router.ServeHTTP(w, req)

	// Check status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Check response
	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["status"] != "healthy" {
		t.Errorf("Expected status 'healthy', got %v", response["status"])
	}

	if response["service"] != "auth" {
		t.Errorf("Expected service 'auth', got %v", response["service"])
	}
}