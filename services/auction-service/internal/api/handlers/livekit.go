package handlers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	// "github.com/livekit/protocol/auth"  // Temporarily commented out
	"go.uber.org/zap"
)

type LiveKitHandler struct {
	logger    *zap.Logger
	apiKey    string
	apiSecret string
	serverURL string
}

func NewLiveKitHandler(logger *zap.Logger) *LiveKitHandler {
	return &LiveKitHandler{
		logger:    logger,
		apiKey:    getEnv("LIVEKIT_API_KEY", "devkey"),
		apiSecret: getEnv("LIVEKIT_API_SECRET", "secret"),
		serverURL: getEnv("LIVEKIT_SERVER_URL", "ws://localhost:7880"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GenerateToken generates a LiveKit token for connecting to a room
func (h *LiveKitHandler) GenerateToken(c *gin.Context) {
	room := c.Query("room")
	role := c.Query("role") // "viewer" or "broadcaster"

	if room == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "room parameter is required"})
		return
	}

	if role == "" {
		role = "viewer" // Default to viewer
	}

	// Validate role
	if role != "viewer" && role != "broadcaster" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "role must be 'viewer' or 'broadcaster'"})
		return
	}

	h.logger.Info("Generating mock LiveKit token",
		zap.String("room", room),
		zap.String("role", role),
		zap.String("user_id", c.GetString("userID")),
	)

	// Mock token for now - will implement real token generation after fixing metrics issue
	mockToken := "mock_token_" + room + "_" + role + "_" + c.GetString("userID")

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"token":     mockToken,
			"room":      room,
			"role":      role,
			"serverUrl": h.serverURL,
		},
	})
}
