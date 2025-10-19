package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/gmsas95/blytz-mvp/services/auth-service/internal/config"
	"github.com/gmsas95/blytz-mvp/services/auth-service/internal/models"
	"github.com/gmsas95/blytz-mvp/services/auth-service/internal/services"
)

func setupTestRouter() (*gin.Engine, *gorm.DB) {
	db, _ := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	db.AutoMigrate(&models.User{})

	cfg := &config.Config{JWTSecret: "test-secret"}
	authService := services.NewAuthService(db, cfg)

	router := gin.Default()
	RegisterRoutes(router, authService)

	return router, db
}

func TestRegisterUser(t *testing.T) {
	router, _ := setupTestRouter()

	user := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password",
	}

	jsonUser, _ := json.Marshal(user)

	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonUser))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestLoginUser(t *testing.T) {
	router, db := setupTestRouter()

	// Create a user to log in
	user := &models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password",
	}
	authService := services.NewAuthService(db, &config.Config{JWTSecret: "test-secret"})
	authService.RegisterUser(user)

	loginDetails := gin.H{
		"email":    "test@example.com",
		"password": "password",
	}

	jsonLogin, _ := json.Marshal(loginDetails)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonLogin))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Contains(t, response, "token")
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