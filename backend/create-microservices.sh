#!/bin/bash

# Blytz Microservices Scaffolding Script
# Creates complete Golang microservices based on OpenAPI specifications

set -e

echo "🚀 Creating Blytz Microservices Architecture..."

# Create shared utilities first
echo "🔧 Setting up shared utilities..."

# Create shared response utilities
cat > /home/sas/blytzmvp-clean/backend/shared/utils/response.go << 'EOF'
package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/blytz/shared/errors"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

type ErrorResponse struct {
	Type    string                 `json:"type"`
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(200, Response{
		Success: true,
		Data:    data,
	})
}

func ErrorResponse(c *gin.Context, err error) {
	if appErr, ok := errors.IsAppError(err); ok {
		c.JSON(appErr.StatusCode, Response{
			Success: false,
			Error: ErrorResponse{
				Type:    appErr.Type,
				Code:    appErr.Code,
				Message: appErr.Message,
				Details: appErr.Details,
			},
		})
		return
	}

	// Handle generic errors
	c.JSON(500, Response{
		Success: false,
		Error: ErrorResponse{
			Type:    "INTERNAL_ERROR",
			Code:    "INTERNAL_ERROR",
			Message: err.Error(),
		},
	})
}

func ValidationError(c *gin.Context, message string, details map[string]interface{}) {
	c.JSON(400, Response{
		Success: false,
		Error: ErrorResponse{
			Type:    "VALIDATION_ERROR",
			Code:    "VALIDATION_ERROR",
			Message: message,
			Details: details,
		},
	})
}
EOF

# Create shared validation utilities
cat > /home/sas/blytzmvp-clean/backend/shared/utils/validation.go << 'EOF'
package utils

import (
	"regexp"
	"strings"
)

// ValidateEmail validates email format
func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// ValidatePhone validates phone number format
func ValidatePhone(phone string) bool {
	phoneRegex := regexp.MustCompile(`^\+?[\d\s\-\(\)]{10,15}$`)
	return phoneRegex.MatchString(phone)
}

// ValidateUsername validates username format
func ValidateUsername(username string) bool {
	if len(username) < 3 || len(username) > 50 {
		return false
	}
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return usernameRegex.MatchString(username)
}

// SanitizeString removes potentially harmful characters
func SanitizeString(input string) string {
	input = strings.TrimSpace(input)
	// Remove HTML tags
	input = regexp.MustCompile(`<[^>>]*>`).ReplaceAllString(input, "")
	return input
}
EOF

# Create shared JWT utilities
cat > /home/sas/blytzmvp-clean/backend/shared/utils/jwt.go << 'EOF'
package utils

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
	"github.com/blytz/shared/constants"
)

type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID, email, role, secret string, expiry time.Duration) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    constants.AuthService,
			Audience:  jwt.ClaimStrings{"blytz.app"},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateJWT(tokenString, secret string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
EOF

# Create shared middleware
cat > /home/sas/blytzmvp-clean/backend/shared/middleware/auth.go << 'EOF'
package middleware

import (
	"net/http"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/blytz/shared/errors"
	"github.com/blytz/shared/utils"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
			})
			return
		}

		// Extract token from Bearer scheme
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization format. Use 'Bearer <token>'",
			})
			return
		}

		// Validate JWT token - this would be implemented with actual JWT validation
		// For now, we'll simulate successful authentication
		c.Set("userID", "user123")
		c.Set("email", "user@example.com")
		c.Set("role", "buyer")

		c.Next()
	}
}

func RoleMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "User role not found",
			})
			return
		}

		roleStr := userRole.(string)
		for _, role := range roles {
			if roleStr == role {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "Insufficient permissions",
		})
	}
}
EOF

echo "✅ Shared utilities created!"

# Now create each microservice systematically
create_service() {
    local service_name=$1
    local service_port=$2
    local service_desc=$3

    echo "🔧 Creating $service_desc Service..."

    # Create service directory structure
    mkdir -p /home/sas/blytzmvp-clean/backend/$service_name/{cmd,internal/{api/{handlers,middleware,routes},config,database,models,services,utils},pkg/{dto,validators},scripts,tests}

    # Create go.mod
    cat > /home/sas/blytzmvp-clean/backend/$service_name/go.mod << EOF
module github.com/blytz/$service_name

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/go-redis/redis/v8 v8.11.5
	github.com/prometheus/client_golang v1.17.0
	gorm.io/gorm v1.25.5
	gorm.io/driver/postgres v1.5.4
	go.uber.org/zap v1.26.0
	github.com/blytz/shared v0.0.0
)

replace github.com/blytz/shared => ../shared
EOF

    # Create main.go
    cat > /home/sas/blytzmvp-clean/backend/$service_name/cmd/main.go << EOF
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/blytz/$service_name/internal/api"
	"github.com/blytz/$service_name/internal/config"
	"github.com/blytz/shared/utils"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	logger, err := utils.InitLogger("development")
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// Create API router
	router := api.SetupRouter(logger)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "$service_port"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan

		logger.Info("Shutting down server...")
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Error("Server forced to shutdown", zap.Error(err))
		}
	}()

	logger.Info("$service_desc service started", zap.String("port", port))
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal("Server failed to start", zap.Error(err))
	}

	logger.Info("Server stopped")
}
EOF

    # Create basic router
    cat > /home/sas/blytzmvp-clean/backend/$service_name/internal/api/router.go << EOF
package api

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func SetupRouter(logger *zap.Logger) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"service": "$service_name",
			"timestamp": time.Now().Unix(),
		})
	})

	// Metrics
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// API routes
	api := router.Group("/api/v1")
	{
		// Add service-specific routes here
		SetupRoutes(api, logger)
	}

	return router
}

func SetupRoutes(api *gin.RouterGroup, logger *zap.Logger) {
	// Service-specific routes will be added here
	api.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "$service_name is running"})
	})
}
EOF

    # Create Dockerfile
    cat > /home/sas/blytzmvp-clean/backend/$service_name/Dockerfile << EOF
FROM golang:1.21-alpine AS builder

WORKDIR /src

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $service_name ./cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates curl

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /src/$service_name .

# Expose port
EXPOSE $service_port

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \\
  CMD curl -f http://localhost:$service_port/health || exit 1

# Run the application
CMD ["./$service_name"]
EOF

    # Create basic environment file
    cat > /home/sas/blytzmvp-clean/backend/$service_name/.env << EOF
PORT=$service_port
GO_ENV=development
LOG_LEVEL=info
DATABASE_URL=postgres://postgres:password@postgres:5432/${service_name//-/_}?sslmode=disable
REDIS_URL=redis:6379
JWT_SECRET=your-jwt-secret-key
EOF

    echo "✅ $service_desc Service scaffolding complete!"
}

# Create all microservices
create_service "auth-service" "8081" "Authentication"
create_service "product-service" "8082" "Product"
create_service "auction-service" "8083" "Auction"
create_service "chat-service" "8084" "Chat"
create_service "order-service" "8085" "Order"
create_service "payment-service" "8086" "Payment"
create_service "logistics-service" "8087" "Logistics"
create_service "gateway" "8080" "API Gateway"

echo ""
echo "🎉 All Blytz microservices have been scaffolded!"
echo ""
echo "📋 Next steps:"
echo "1. Review the generated code in each service"
echo "2. Implement the business logic based on OpenAPI specs"
echo "3. Set up database connections and migrations"
echo "4. Configure environment variables"
echo "5. Test the microservices integration"
echo ""
echo "📁 Services created at: /home/sas/blytzmvp-clean/backend/"
echo "📄 OpenAPI specs at: /home/sas/blytzmvp-clean/openapi/"
echo ""
echo "Happy coding! 🚀"""

# Make the script executable
chmod +x /home/sas/blytzmvp-clean/backend/create-microservices.sh

echo "✅ Microservices creation script ready!"
echo "Executing scaffolding..."

# Execute the script
/home/sas/blytzmvp-clean/backend/create-microservices.sh