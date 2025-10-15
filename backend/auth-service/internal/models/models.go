package models

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// Identity
	UserID       string `gorm:"uniqueIndex;not null" json:"user_id"`
	Email        string `gorm:"uniqueIndex;not null" json:"email"`
	Username     string `gorm:"uniqueIndex;not null" json:"username"`
	FirebaseUID  string `gorm:"uniqueIndex" json:"firebase_uid"`

	// Profile
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Phone       string `json:"phone"`
	AvatarURL   string `json:"avatar_url"`

	// Status
	IsVerified  bool   `gorm:"default:false" json:"is_verified"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
	Role        string `gorm:"default:user" json:"role"`

	// Metadata
	Metadata    string `gorm:"type:text" json:"metadata,omitempty"`
}

type AuthRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Username  string `json:"username" binding:"required,min=3,max=50"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Role      string `json:"role" binding:"required,oneof=buyer seller admin"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	User         User   `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type VerifyResponse struct {
	User  User  `json:"user"`
	Valid bool  `json:"valid"`
}

type UpdateProfileRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	AvatarURL string `json:"avatar_url"`
}