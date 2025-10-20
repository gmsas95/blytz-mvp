package services

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"

	"github.com/gmsas95/blytz-mvp/services/chat-service/internal/config"
	"github.com/gmsas95/blytz-mvp/services/chat-service/internal/models"
)

type ChatService struct {
	redis  *redis.Client
	logger *zap.Logger
	config *config.Config
}

func NewChatService(logger *zap.Logger, config *config.Config) *ChatService {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.RedisURL,
		Password: config.RedisPassword,
		DB:       0,
	})

	return &ChatService{
		redis:  redisClient,
		logger: logger,
		config: config,
	}
}

func (s *ChatService) SendMessage(ctx context.Context, roomID, userID, content string) (*models.Message, error) {
	s.logger.Info("Sending message", zap.String("room_id", roomID), zap.String("user_id", userID))

	message := &models.Message{
		ID:        generateMessageID(),
		RoomID:    roomID,
		UserID:    userID,
		Content:   content,
		Timestamp: time.Now(),
	}

	// Store message in Redis
	messageData, err := json.Marshal(message)
	if err != nil {
		s.logger.Error("Failed to marshal message", zap.Error(err))
		return nil, err
	}

	// Add to room messages list
	roomKey := "chat:room:" + roomID + ":messages"
	if err := s.redis.LPush(ctx, roomKey, messageData).Err(); err != nil {
		s.logger.Error("Failed to store message", zap.Error(err))
		return nil, err
	}

	// Set expiration for room messages (7 days)
	if err := s.redis.Expire(ctx, roomKey, 7*24*time.Hour).Err(); err != nil {
		s.logger.Error("Failed to set expiration", zap.Error(err))
	}

	// Publish to room subscribers
	publishData := map[string]interface{}{
		"type":    "message",
		"room_id": roomID,
		"message": message,
	}

	publishJSON, err := json.Marshal(publishData)
	if err != nil {
		s.logger.Error("Failed to marshal publish data", zap.Error(err))
		return nil, err
	}

	channel := "chat:room:" + roomID
	if err := s.redis.Publish(ctx, channel, publishJSON).Err(); err != nil {
		s.logger.Error("Failed to publish message", zap.Error(err))
		return nil, err
	}

	s.logger.Info("Message sent successfully", zap.String("message_id", message.ID))
	return message, nil
}

func (s *ChatService) GetMessages(ctx context.Context, roomID string, limit int64) ([]*models.Message, error) {
	s.logger.Info("Getting messages", zap.String("room_id", roomID))

	roomKey := "chat:room:" + roomID + ":messages"

	// Get messages from Redis
	messagesData, err := s.redis.LRange(ctx, roomKey, 0, limit-1).Result()
	if err != nil {
		s.logger.Error("Failed to get messages", zap.Error(err))
		return nil, err
	}

	messages := make([]*models.Message, 0, len(messagesData))
	for _, data := range messagesData {
		var message models.Message
		if err := json.Unmarshal([]byte(data), &message); err != nil {
			s.logger.Error("Failed to unmarshal message", zap.Error(err))
			continue
		}
		messages = append(messages, &message)
	}

	return messages, nil
}

func (s *ChatService) GetUserRooms(ctx context.Context, userID string) ([]*models.ChatRoom, error) {
	s.logger.Info("Getting user rooms", zap.String("user_id", userID))

	// Get user's rooms from Redis
	userRoomsKey := "chat:user:" + userID + ":rooms"
	roomIDs, err := s.redis.SMembers(ctx, userRoomsKey).Result()
	if err != nil {
		s.logger.Error("Failed to get user rooms", zap.Error(err))
		return nil, err
	}

	rooms := make([]*models.ChatRoom, 0, len(roomIDs))
	for _, roomID := range roomIDs {
		roomKey := "chat:room:" + roomID
		roomData, err := s.redis.HGetAll(ctx, roomKey).Result()
		if err != nil {
			s.logger.Error("Failed to get room data", zap.Error(err))
			continue
		}

		if len(roomData) == 0 {
			continue
		}

		room := &models.ChatRoom{
			ID:        roomID,
			Name:      roomData["name"],
			Type:      roomData["type"],
			CreatedAt: parseTime(roomData["created_at"]),
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}

func (s *ChatService) JoinRoom(ctx context.Context, roomID, userID string) error {
	s.logger.Info("User joining room", zap.String("room_id", roomID), zap.String("user_id", userID))

	// Add user to room set
	userRoomsKey := "chat:user:" + userID + ":rooms"
	if err := s.redis.SAdd(ctx, userRoomsKey, roomID).Err(); err != nil {
		s.logger.Error("Failed to add user to room", zap.Error(err))
		return err
	}

	// Set expiration for user rooms (30 days)
	if err := s.redis.Expire(ctx, userRoomsKey, 30*24*time.Hour).Err(); err != nil {
		s.logger.Error("Failed to set expiration", zap.Error(err))
	}

	return nil
}

func (s *ChatService) LeaveRoom(ctx context.Context, roomID, userID string) error {
	s.logger.Info("User leaving room", zap.String("room_id", roomID), zap.String("user_id", userID))

	userRoomsKey := "chat:user:" + userID + ":rooms"
	if err := s.redis.SRem(ctx, userRoomsKey, roomID).Err(); err != nil {
		s.logger.Error("Failed to remove user from room", zap.Error(err))
		return err
	}

	return nil
}

func generateMessageID() string {
	return fmt.Sprintf("msg_%d", time.Now().UnixNano())
}

func parseTime(timeStr string) time.Time {
	t, _ := time.Parse(time.RFC3339, timeStr)
	return t
}