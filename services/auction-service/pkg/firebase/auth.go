package firebase

import (
	"context"
	"fmt"
)

// CreateUserData represents user creation request
type CreateUserData struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"displayName,omitempty"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
}

// CreateUserResponse represents user creation response
type CreateUserResponse struct {
	Success bool   `json:"success"`
	UID     string `json:"uid"`
	Message string `json:"message"`
}

// ValidateTokenResponse represents token validation response
type ValidateTokenResponse struct {
	Success bool                   `json:"success"`
	User    map[string]interface{} `json:"user"`
}

// UpdateProfileData represents profile update request
type UpdateProfileData struct {
	DisplayName string `json:"displayName,omitempty"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
}

// UpdateProfileResponse represents profile update response
type UpdateProfileResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// CreateUser creates a new user via Firebase
func (c *Client) CreateUser(ctx context.Context, email, password, displayName, phoneNumber string) (*CreateUserResponse, error) {
	data := CreateUserData{
		Email:       email,
		Password:    password,
		DisplayName: displayName,
		PhoneNumber: phoneNumber,
	}

	var result CreateUserResponse
	err := c.callFirebase(ctx, "createUser", data, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return &result, nil
}

// ValidateToken validates a Firebase ID token
func (c *Client) ValidateToken(ctx context.Context) (*ValidateTokenResponse, error) {
	// This would typically include the token in the request
	// For now, we'll call the function without additional data
	var result ValidateTokenResponse
	err := c.callFirebase(ctx, "validateToken", nil, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %w", err)
	}
	return &result, nil
}

// UpdateProfile updates user profile via Firebase
func (c *Client) UpdateProfile(ctx context.Context, displayName, phoneNumber string) (*UpdateProfileResponse, error) {
	data := UpdateProfileData{
		DisplayName: displayName,
		PhoneNumber: phoneNumber,
	}

	var result UpdateProfileResponse
	err := c.callFirebase(ctx, "updateProfile", data, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}
	return &result, nil
}