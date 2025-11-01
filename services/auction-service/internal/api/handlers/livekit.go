package handlers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/livekit/protocol/auth"
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

	h.logger.Info("Generating LiveKit token",
		zap.String("room", room),
		zap.String("role", role),
		zap.String("user_id", c.GetString("userID")),
	)

	// Create access token
	at := auth.NewAccessToken(h.apiKey, h.apiSecret)

	// Set identity
	userID := c.GetString("userID")
	if userID == "" {
		userID = "anonymous"
	}
	at.SetIdentity(userID)

	// Set name (optional)
	name := c.Query("name")
	if name != "" {
		at.SetName(name)
	}

	// Set video grant based on role
	grant := &auth.VideoGrant{
		RoomJoin:     true,
		Room:         room,
		CanPublish:   role == "broadcaster",
		CanSubscribe: true,
	}

	at.AddGrant(grant)

	// Set token validity (24 hours)
	at.SetValidFor(24 * 3600)

	// Generate JWT token
	token, err := at.ToJWT()
	if err != nil {
		h.logger.Error("Failed to generate LiveKit token", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"token":     token,
			"room":      room,
			"role":      role,
			"serverUrl": h.serverURL,
		},
	})
}
