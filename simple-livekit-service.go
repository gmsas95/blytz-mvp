package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type LiveKitTokenResponse struct {
	Token string `json:"token"`
	URL   string `json:"url"`
}

func main() {
	// Get environment variables
	apiKey := os.Getenv("LIVEKIT_API_KEY")
	apiSecret := os.Getenv("LIVEKIT_API_SECRET")

	if apiKey == "" || apiSecret == "" {
		log.Fatal("LIVEKIT_API_KEY and LIVEKIT_API_SECRET are required")
	}

	// Create simple HTTP server
	http.HandleFunc("/api/livekit/token", func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Get query parameters
		room := r.URL.Query().Get("room")
		if room == "" {
			room = "demo-room"
		}

		name := r.URL.Query().Get("name")
		if name == "" {
			name = "demo-user"
		}

		// Create JWT token for LiveKit
		claims := jwt.MapClaims{
			"iss": apiKey,
			"nbf": time.Now().Unix(),
			"exp": time.Now().Add(24 * time.Hour).Unix(),
			"sub": name,
			"video": map[string]interface{}{
				"room":     room,
				"roomJoin": true,
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(apiSecret))
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		// Return token
		response := LiveKitTokenResponse{
			Token: tokenString,
			URL:   "wss://livekit.blytz.app",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	// Health endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"service": "livekit-token-service",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8089"
	}

	log.Printf("Starting LiveKit token service on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
