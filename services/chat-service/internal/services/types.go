package services

import "time"

type SendMessageRequest struct {
	Content string `json:"content" binding:"required"`
}

type MessageResponse struct {
	ID        string `json:"id"`
	RoomID    string `json:"room_id"`
	UserID    string `json:"user_id"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

type ChatRoomResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	CreatedAt string `json:"created_at"`
}

type GetMessagesResponse struct {
	Messages []MessageResponse `json:"messages"`
}

type GetRoomsResponse struct {
	Rooms []ChatRoomResponse `json:"rooms"`
}

func formatTime(t time.Time) string {
	return t.Format("2006-01-02T15:04:05Z")
}