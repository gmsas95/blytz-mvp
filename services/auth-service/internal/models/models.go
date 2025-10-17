package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Email       string    `json:"email" gorm:"uniqueIndex;not null"`
	DisplayName string    `json:"display_name"`
	PhoneNumber string    `json:"phone_number,omitempty"`
	AvatarURL   string    `json:"avatar_url,omitempty"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	Role        string    `json:"role" gorm:"default:user"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// RegisterRequest represents user registration request
type RegisterRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8"`
	DisplayName string `json:"display_name" binding:"required,min=2,max=50"`
	PhoneNumber string `json:"phone_number,omitempty" binding:"omitempty,e164"`
}

// LoginRequest represents user login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	User    *User  `json:"user,omitempty"`
	Token   string `json:"token,omitempty"`
}

// ValidateTokenRequest represents token validation request
type ValidateTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

// ValidateTokenResponse represents token validation response
type ValidateTokenResponse struct {
	Valid   bool   `json:"valid"`
	UserID  string `json:"user_id,omitempty"`
	Email   string `json:"email,omitempty"`
	Message string `json:"message,omitempty"`
}

// RefreshTokenRequest represents token refresh request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// UpdateProfileRequest represents profile update request
type UpdateProfileRequest struct {
	DisplayName string `json:"display_name,omitempty" binding:"omitempty,min=2,max=50"`
	PhoneNumber string `json:"phone_number,omitempty" binding:"omitempty,e164"`
	AvatarURL   string `json:"avatar_url,omitempty" binding:"omitempty,url"`
}

// ChangePasswordRequest represents password change request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}

// ErrorResponse represents error response
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

// SuccessResponse represents success response
type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
