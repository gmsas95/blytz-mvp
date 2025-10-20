package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	"github.com/gmsas95/blytz-mvp/services/chat-service/internal/services"
	"github.com/gmsas95/blytz-mvp/shared/pkg/errors"
	"github.com/gmsas95/blytz-mvp/shared/pkg/utils"
)

type ChatHandler struct {
	chatService *services.ChatService
	logger      *zap.Logger
	upgrader    websocket.Upgrader
}

func NewChatHandler(chatService *services.ChatService, logger *zap.Logger) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
		logger:      logger,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for now
			},
		},
	}
}

func (h *ChatHandler) HandleWebSocket(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		utils.ErrorResponse(c, errors.ErrUnauthorized)
		return
	}

	roomID := c.Query("roomId")
	if roomID == "" {
		utils.ValidationError(c, "Room ID is required", nil)
		return
	}

	// Upgrade connection to WebSocket
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		h.logger.Error("Failed to upgrade connection", zap.Error(err))
		return
	}
	defer conn.Close()

	h.logger.Info("WebSocket connection established", zap.String("user_id", userID), zap.String("room_id", roomID))

	// Join room
	if err := h.chatService.JoinRoom(c.Request.Context(), roomID, userID); err != nil {
		h.logger.Error("Failed to join room", zap.Error(err))
		return
	}

	// Handle WebSocket messages
	for {
		var msg services.SendMessageRequest
		if err := conn.ReadJSON(&msg); err != nil {
			h.logger.Error("Failed to read message", zap.Error(err))
			break
		}

		// Send message
		message, err := h.chatService.SendMessage(c.Request.Context(), roomID, userID, msg.Content)
		if err != nil {
			h.logger.Error("Failed to send message", zap.Error(err))
			continue
		}

		// Send response back
		response := services.MessageResponse{
			ID:        message.ID,
			RoomID:    message.RoomID,
			UserID:    message.UserID,
			Content:   message.Content,
			Timestamp: message.Timestamp.Format("2006-01-02T15:04:05Z"),
		}

		if err := conn.WriteJSON(response); err != nil {
			h.logger.Error("Failed to write response", zap.Error(err))
			break
		}
	}

	// Leave room when connection closes
	if err := h.chatService.LeaveRoom(c.Request.Context(), roomID, userID); err != nil {
		h.logger.Error("Failed to leave room", zap.Error(err))
	}
}

func (h *ChatHandler) SendMessage(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		utils.ErrorResponse(c, errors.ErrUnauthorized)
		return
	}

	roomID := c.Param("roomId")
	if roomID == "" {
		utils.ValidationError(c, "Room ID is required", nil)
		return
	}

	var req services.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, "Invalid request body", nil)
		return
	}

	message, err := h.chatService.SendMessage(c.Request.Context(), roomID, userID, req.Content)
	if err != nil {
		h.logger.Error("Failed to send message", zap.Error(err))
		utils.ErrorResponse(c, err)
		return
	}

	response := services.MessageResponse{
		ID:        message.ID,
		RoomID:    message.RoomID,
		UserID:    message.UserID,
		Content:   message.Content,
		Timestamp: message.Timestamp.Format("2006-01-02T15:04:05Z"),
	}

	utils.SuccessResponse(c, response)
}

func (h *ChatHandler) GetMessages(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		utils.ErrorResponse(c, errors.ErrUnauthorized)
		return
	}

	roomID := c.Param("roomId")
	if roomID == "" {
		utils.ValidationError(c, "Room ID is required", nil)
		return
	}

	limit := int64(50)
	messages, err := h.chatService.GetMessages(c.Request.Context(), roomID, limit)
	if err != nil {
		h.logger.Error("Failed to get messages", zap.Error(err))
		utils.ErrorResponse(c, err)
		return
	}

	response := make([]services.MessageResponse, len(messages))
	for i, msg := range messages {
		response[i] = services.MessageResponse{
			ID:        msg.ID,
			RoomID:    msg.RoomID,
			UserID:    msg.UserID,
			Content:   msg.Content,
			Timestamp: msg.Timestamp.Format("2006-01-02T15:04:05Z"),
		}
	}

	utils.SuccessResponse(c, services.GetMessagesResponse{Messages: response})
}

func (h *ChatHandler) GetUserRooms(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		utils.ErrorResponse(c, errors.ErrUnauthorized)
		return
	}

	rooms, err := h.chatService.GetUserRooms(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error("Failed to get user rooms", zap.Error(err))
		utils.ErrorResponse(c, err)
		return
	}

	response := make([]services.ChatRoomResponse, len(rooms))
	for i, room := range rooms {
		response[i] = services.ChatRoomResponse{
			ID:        room.ID,
			Name:      room.Name,
			Type:      room.Type,
			CreatedAt: room.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	utils.SuccessResponse(c, services.GetRoomsResponse{Rooms: response})
}