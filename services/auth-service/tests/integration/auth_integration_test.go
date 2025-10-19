// +build integration

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gmsas95/blytz-mvp/services/auth-service/internal/models"
)

const (
	baseURL = "http://localhost:8084"
	timeout = 10 * time.Second
)

func waitForService(t *testing.T) {
	client := &http.Client{Timeout: timeout}
	maxRetries := 30

	for i := 0; i < maxRetries; i++ {
		resp, err := client.Get(baseURL + "/health")
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(1 * time.Second)
	}

	t.Fatal("Service did not become healthy within timeout")
}

func TestAuthIntegration(t *testing.T) {
	// Wait for service to be ready
	waitForService(t)

	client := &http.Client{Timeout: timeout}

	t.Run("Complete auth flow", func(t *testing.T) {
		// Step 1: Register a new user
		registerReq := models.RegisterRequest{
			Email:       "integration@example.com",
			Password:    "password123",
			DisplayName: "Integration Test User",
			PhoneNumber: "+1234567890",
		}

		jsonBody, err := json.Marshal(registerReq)
		if err != nil {
			t.Fatalf("Failed to marshal register request: %v", err)
		}

		resp, err := client.Post(baseURL+"/api/v1/auth/register", "application/json", bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Fatalf("Failed to register user: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body := make([]byte, 1024)
			n, _ := resp.Body.Read(body)
			t.Fatalf("Expected status 200, got %d. Response: %s", resp.StatusCode, string(body[:n]))
		}

		var registerResp models.AuthResponse
		if err := json.NewDecoder(resp.Body).Decode(&registerResp); err != nil {
			t.Fatalf("Failed to decode register response: %v", err)
		}

		if !registerResp.Success {
			t.Fatal("Expected successful registration")
		}

		if registerResp.Token == "" {
			t.Fatal("Expected token in registration response")
		}

		// Step 2: Login with the same credentials
		loginReq := models.LoginRequest{
			Email:    "integration@example.com",
			Password: "password123",
		}

		jsonBody, err = json.Marshal(loginReq)
		if err != nil {
			t.Fatalf("Failed to marshal login request: %v", err)
		}

		resp, err = client.Post(baseURL+"/api/v1/auth/login", "application/json", bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Fatalf("Failed to login: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body := make([]byte, 1024)
			n, _ := resp.Body.Read(body)
			t.Fatalf("Expected status 200, got %d. Response: %s", resp.StatusCode, string(body[:n]))
		}

		var loginResp models.AuthResponse
		if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
			t.Fatalf("Failed to decode login response: %v", err)
		}

		if !loginResp.Success {
			t.Fatal("Expected successful login")
		}

		if loginResp.Token == "" {
			t.Fatal("Expected token in login response")
		}

		// Step 3: Validate the token
		validateReq := models.ValidateTokenRequest{
			Token: loginResp.Token,
		}

		jsonBody, err = json.Marshal(validateReq)
		if err != nil {
			t.Fatalf("Failed to marshal validate request: %v", err)
		}

		resp, err = client.Post(baseURL+"/api/v1/auth/validate", "application/json", bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Fatalf("Failed to validate token: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body := make([]byte, 1024)
			n, _ := resp.Body.Read(body)
			t.Fatalf("Expected status 200, got %d. Response: %s", resp.StatusCode, string(body[:n]))
		}

		var validateResp models.ValidateTokenResponse
		if err := json.NewDecoder(resp.Body).Decode(&validateResp); err != nil {
			t.Fatalf("Failed to decode validate response: %v", err)
		}

		if !validateResp.Valid {
			t.Fatal("Expected token to be valid")
		}

		if validateResp.Email != "integration@example.com" {
			t.Fatalf("Expected email integration@example.com, got %s", validateResp.Email)
		}

		// Step 4: Refresh the token
		refreshReq := models.RefreshTokenRequest{
			RefreshToken: loginResp.Token,
		}

		jsonBody, err = json.Marshal(refreshReq)
		if err != nil {
			t.Fatalf("Failed to marshal refresh request: %v", err)
		}

		resp, err = client.Post(baseURL+"/api/v1/auth/refresh", "application/json", bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Fatalf("Failed to refresh token: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body := make([]byte, 1024)
			n, _ := resp.Body.Read(body)
			t.Fatalf("Expected status 200, got %d. Response: %s", resp.StatusCode, string(body[:n]))
		}

		var refreshResp models.SuccessResponse
		if err := json.NewDecoder(resp.Body).Decode(&refreshResp); err != nil {
			t.Fatalf("Failed to decode refresh response: %v", err)
		}

		if !refreshResp.Success {
			t.Fatal("Expected successful token refresh")
		}
	})

	t.Run("Invalid token validation", func(t *testing.T) {
		validateReq := models.ValidateTokenRequest{
			Token: "invalid.token.here",
		}

		jsonBody, err := json.Marshal(validateReq)
		if err != nil {
			t.Fatalf("Failed to marshal validate request: %v", err)
		}

		resp, err := client.Post(baseURL+"/api/v1/auth/validate", "application/json", bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Fatalf("Failed to validate token: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}

		var validateResp models.ValidateTokenResponse
		if err := json.NewDecoder(resp.Body).Decode(&validateResp); err != nil {
			t.Fatalf("Failed to decode validate response: %v", err)
		}

		if validateResp.Valid {
			t.Fatal("Expected token to be invalid")
		}

		if validateResp.Message == "" {
			t.Fatal("Expected message in invalid token response")
		}
	})

	t.Run("Duplicate registration", func(t *testing.T) {
		registerReq := models.RegisterRequest{
			Email:       "duplicate@example.com",
			Password:    "password123",
			DisplayName: "Duplicate Test User",
		}

		jsonBody, err := json.Marshal(registerReq)
		if err != nil {
			t.Fatalf("Failed to marshal register request: %v", err)
		}

		// First registration should succeed
		resp, err := client.Post(baseURL+"/api/v1/auth/register", "application/json", bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Fatalf("Failed to register user: %v", err)
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200 for first registration, got %d", resp.StatusCode)
		}

		// Second registration should fail
		resp, err = client.Post(baseURL+"/api/v1/auth/register", "application/json", bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Fatalf("Failed to register user second time: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			body := make([]byte, 1024)
			n, _ := resp.Body.Read(body)
			t.Fatalf("Expected status 400, got %d. Response: %s", resp.StatusCode, string(body[:n]))
		}
	})

	t.Run("Invalid login credentials", func(t *testing.T) {
		loginReq := models.LoginRequest{
			Email:    "nonexistent@example.com",
			Password: "wrongpassword",
		}

		jsonBody, err := json.Marshal(loginReq)
		if err != nil {
			t.Fatalf("Failed to marshal login request: %v", err)
		}

		resp, err := client.Post(baseURL+"/api/v1/auth/login", "application/json", bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Fatalf("Failed to login: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			body := make([]byte, 1024)
			n, _ := resp.Body.Read(body)
			t.Fatalf("Expected status 401, got %d. Response: %s", resp.StatusCode, string(body[:n]))
		}
	})
}

func TestHealthCheck(t *testing.T) {
	client := &http.Client{Timeout: timeout}

	resp, err := client.Get(baseURL + "/health")
	if err != nil {
		t.Fatalf("Failed to check health: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	var healthResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&healthResp); err != nil {
		t.Fatalf("Failed to decode health response: %v", err)
	}

	if healthResp["status"] != "healthy" {
		t.Fatalf("Expected status 'healthy', got %v", healthResp["status"])
	}

	if healthResp["service"] != "auth" {
		t.Fatalf("Expected service 'auth', got %v", healthResp["service"])
	}

	if _, ok := healthResp["timestamp"].(float64); !ok {
		t.Fatal("Expected timestamp to be a number")
	}
}

func TestConcurrentRequests(t *testing.T) {
	client := &http.Client{Timeout: timeout}

	// Create multiple users concurrently
	numRequests := 10
	done := make(chan bool, numRequests)

	for i := 0; i < numRequests; i++ {
		go func(index int) {
			defer func() { done <- true }()

			registerReq := models.RegisterRequest{
				Email:       fmt.Sprintf("concurrent%d@example.com", index),
				Password:    "password123",
				DisplayName: fmt.Sprintf("Concurrent User %d", index),
			}

			jsonBody, err := json.Marshal(registerReq)
			if err != nil {
				t.Errorf("Failed to marshal request %d: %v", index, err)
				return
			}

			resp, err := client.Post(baseURL+"/api/v1/auth/register", "application/json", bytes.NewBuffer(jsonBody))
			if err != nil {
				t.Errorf("Failed to register user %d: %v", index, err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Errorf("Expected status 200 for request %d, got %d", index, resp.StatusCode)
			}
		}(i)
	}

	// Wait for all requests to complete
	for i := 0; i < numRequests; i++ {
		<-done
	}
}

func BenchmarkAuthEndpoints(b *testing.B) {
	client := &http.Client{Timeout: timeout}

	// First create a test user
	registerReq := models.RegisterRequest{
		Email:       "benchmark@example.com",
		Password:    "password123",
		DisplayName: "Benchmark User",
	}

	jsonBody, _ := json.Marshal(registerReq)
	resp, _ := client.Post(baseURL+"/api/v1/auth/register", "application/json", bytes.NewBuffer(jsonBody))
	if resp != nil {
		resp.Body.Close()
	}

	loginReq := models.LoginRequest{
		Email:    "benchmark@example.com",
		Password: "password123",
	}

	jsonBody, _ = json.Marshal(loginReq)
	resp, _ = client.Post(baseURL+"/api/v1/auth/login", "application/json", bytes.NewBuffer(jsonBody))
	if resp != nil {
		resp.Body.Close()
	}

	b.Run("Login", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			jsonBody, _ := json.Marshal(loginReq)
			resp, err := client.Post(baseURL+"/api/v1/auth/login", "application/json", bytes.NewBuffer(jsonBody))
			if err != nil {
				b.Fatalf("Failed to login: %v", err)
			}
			resp.Body.Close()
		}
	})

	b.Run("ValidateToken", func(b *testing.B) {
		// Get a valid token first
		jsonBody, _ := json.Marshal(loginReq)
		resp, _ := client.Post(baseURL+"/api/v1/auth/login", "application/json", bytes.NewBuffer(jsonBody))
		if resp != nil {
			defer resp.Body.Close()
			var loginResp models.AuthResponse
			json.NewDecoder(resp.Body).Decode(&loginResp)

			validateReq := models.ValidateTokenRequest{
				Token: loginResp.Token,
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				jsonBody, _ := json.Marshal(validateReq)
				resp, err := client.Post(baseURL+"/api/v1/auth/validate", "application/json", bytes.NewBuffer(jsonBody))
				if err != nil {
					b.Fatalf("Failed to validate token: %v", err)
				}
				resp.Body.Close()
			}
		}
	})

	b.Run("HealthCheck", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			resp, err := client.Get(baseURL + "/health")
			if err != nil {
				b.Fatalf("Failed to check health: %v", err)
			}
			resp.Body.Close()
		}
	})
}

