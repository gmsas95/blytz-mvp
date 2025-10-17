package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/blytz/auth-service/internal/models"
	"go.uber.org/zap"
)

// Database interface defines the database operations
type Database interface {
	// User operations
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, userID string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUserProfile(ctx context.Context, userID string, req *models.UpdateProfileRequest) error
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, userID string) error

	// Health check
	Health(ctx context.Context) error
	Close() error
}

// InMemoryDatabase implements Database interface using in-memory storage
type InMemoryDatabase struct {
	users  map[string]*models.User
	emails map[string]string // email -> userID mapping
	mu     sync.RWMutex
	logger *zap.Logger
}

// NewDatabase creates a new database instance
func NewDatabase(databaseURL string, logger *zap.Logger) (Database, error) {
	logger.Info("Creating in-memory database for auth service")

	return &InMemoryDatabase{
		users:  make(map[string]*models.User),
		emails: make(map[string]string),
		logger: logger,
	}, nil
}

// CreateUser creates a new user in the database
func (db *InMemoryDatabase) CreateUser(ctx context.Context, user *models.User) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.logger.Info("Creating user in database", zap.String("user_id", user.ID), zap.String("email", user.Email))

	// Check if user already exists
	if _, exists := db.users[user.ID]; exists {
		return fmt.Errorf("user with ID %s already exists", user.ID)
	}

	// Check if email already exists
	if _, exists := db.emails[user.Email]; exists {
		return fmt.Errorf("user with email %s already exists", user.Email)
	}

	// Set timestamps
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	// Store user
	db.users[user.ID] = user
	db.emails[user.Email] = user.ID

	db.logger.Info("User created successfully", zap.String("user_id", user.ID))
	return nil
}

// GetUserByID retrieves a user by ID
func (db *InMemoryDatabase) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	db.logger.Debug("Getting user by ID", zap.String("user_id", userID))

	user, exists := db.users[userID]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}

	// Return a copy to prevent external modifications
	userCopy := *user
	return &userCopy, nil
}

// GetUserByEmail retrieves a user by email
func (db *InMemoryDatabase) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	db.logger.Debug("Getting user by email", zap.String("email", email))

	userID, exists := db.emails[email]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}

	user, exists := db.users[userID]
	if !exists {
		return nil, fmt.Errorf("user not found")
	}

	// Return a copy to prevent external modifications
	userCopy := *user
	return &userCopy, nil
}

// UpdateUserProfile updates user profile information
func (db *InMemoryDatabase) UpdateUserProfile(ctx context.Context, userID string, req *models.UpdateProfileRequest) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.logger.Info("Updating user profile", zap.String("user_id", userID))

	user, exists := db.users[userID]
	if !exists {
		return fmt.Errorf("user not found")
	}

	// Update fields if provided
	if req.DisplayName != "" {
		user.DisplayName = req.DisplayName
	}
	if req.PhoneNumber != "" {
		user.PhoneNumber = req.PhoneNumber
	}
	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}

	user.UpdatedAt = time.Now()

	db.logger.Info("User profile updated successfully", zap.String("user_id", userID))
	return nil
}

// UpdateUser updates a user in the database
func (db *InMemoryDatabase) UpdateUser(ctx context.Context, user *models.User) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.logger.Info("Updating user", zap.String("user_id", user.ID))

	_, exists := db.users[user.ID]
	if !exists {
		return fmt.Errorf("user not found")
	}

	user.UpdatedAt = time.Now()
	db.users[user.ID] = user

	db.logger.Info("User updated successfully", zap.String("user_id", user.ID))
	return nil
}

// DeleteUser deletes a user from the database
func (db *InMemoryDatabase) DeleteUser(ctx context.Context, userID string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.logger.Info("Deleting user", zap.String("user_id", userID))

	user, exists := db.users[userID]
	if !exists {
		return fmt.Errorf("user not found")
	}

	// Remove from email mapping
	delete(db.emails, user.Email)

	// Remove user
	delete(db.users, userID)

	db.logger.Info("User deleted successfully", zap.String("user_id", userID))
	return nil
}

// Health checks the database health
func (db *InMemoryDatabase) Health(ctx context.Context) error {
	db.mu.RLock()
	defer db.mu.RUnlock()

	// Simple health check - ensure we can access the maps
	if db.users == nil || db.emails == nil {
		return fmt.Errorf("database not initialized")
	}

	return nil
}

// Close closes the database connection
func (db *InMemoryDatabase) Close() error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.logger.Info("Closing in-memory database")

	// Clear all data
	db.users = make(map[string]*models.User)
	db.emails = make(map[string]string)

	return nil
}

// GetStats returns database statistics (for monitoring)
func (db *InMemoryDatabase) GetStats() map[string]interface{} {
	db.mu.RLock()
	defer db.mu.RUnlock()

	return map[string]interface{}{
		"total_users": len(db.users),
		"total_emails": len(db.emails),
	}
}

// PostgreSQLDatabase implements Database interface using PostgreSQL
type PostgreSQLDatabase struct {
	// This would contain PostgreSQL connection and methods
	// For MVP, we'll use in-memory database
	db     *InMemoryDatabase
	logger *zap.Logger
}

// NewPostgreSQLDatabase creates a new PostgreSQL database instance
func NewPostgreSQLDatabase(databaseURL string, logger *zap.Logger) (Database, error) {
	logger.Info("Creating PostgreSQL database connection", zap.String("url", databaseURL))

	// For MVP, fall back to in-memory database
	// In production, this would connect to actual PostgreSQL
	return NewDatabase(databaseURL, logger)
}