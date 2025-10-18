package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"
)

// CORS middleware to allow frontend requests
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow requests from frontend
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// createReverseProxy creates a reverse proxy for a target service
func createReverseProxy(targetURL string) (*httputil.ReverseProxy, error) {
	target, err := url.Parse(targetURL)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	// Customize the director to modify requests
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		// Remove the service prefix from the path
		req.URL.Path = strings.TrimPrefix(req.URL.Path, strings.Split(req.URL.Path, "/")[1])
		if req.URL.Path == "" {
			req.URL.Path = "/"
		}
		// Set proper headers
		req.Header.Set("X-Forwarded-Host", req.Host)
		req.Header.Set("X-Forwarded-Proto", "http")
	}

	return proxy, nil
}

func main() {
	log.Println("🚀 Gateway service starting on port 8080")

	// Get service URLs from environment variables
	auctionServiceURL := os.Getenv("AUCTION_SERVICE_URL")
	if auctionServiceURL == "" {
		auctionServiceURL = "http://auction-service:8083"
	}

	authServiceURL := os.Getenv("AUTH_SERVICE_URL")
	if authServiceURL == "" {
		authServiceURL = "http://auth-service:8084"
	}

	productServiceURL := os.Getenv("PRODUCT_SERVICE_URL")
	if productServiceURL == "" {
		productServiceURL = "http://product-service:8082"
	}

	// Create reverse proxies
	auctionProxy, err := createReverseProxy(auctionServiceURL)
	if err != nil {
		log.Fatal("Failed to create auction proxy:", err)
	}

	authProxy, err := createReverseProxy(authServiceURL)
	if err != nil {
		log.Fatal("Failed to create auth proxy:", err)
	}

	productProxy, err := createReverseProxy(productServiceURL)
	if err != nil {
		log.Fatal("Failed to create product proxy:", err)
	}

	// Create router
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"gateway"}`))
	})

	// Route auction service requests
	mux.Handle("/api/v1/auctions/", http.StripPrefix("/api/v1/auctions", auctionProxy))
	mux.Handle("/api/v1/bids/", http.StripPrefix("/api/v1/bids", auctionProxy))

	// Route auth service requests
	mux.Handle("/api/v1/auth/", http.StripPrefix("/api/v1/auth", authProxy))
	mux.Handle("/api/v1/users/", http.StripPrefix("/api/v1/users", authProxy))

	// Route product service requests
	mux.Handle("/api/v1/products/", http.StripPrefix("/api/v1/products", productProxy))

	// Default handler for unmatched routes
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Unhandled request: %s %s", r.Method, r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"Route not found","path":"` + r.URL.Path + `"}`))
	})

	// Create server with CORS middleware
	server := &http.Server{
		Addr:         ":8080",
		Handler:      corsMiddleware(mux),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Println("✅ Gateway service ready!")
	log.Println("📍 Auction service proxy:", auctionServiceURL)
	log.Println("🔐 Auth service proxy:", authServiceURL)
	log.Println("📦 Product service proxy:", productServiceURL)

	log.Fatal(server.ListenAndServe())
}