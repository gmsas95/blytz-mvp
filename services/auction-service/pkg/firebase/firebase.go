package firebase

import "context"

// FirebaseApp represents the Firebase application interface
type FirebaseApp interface {
	// Auth methods
	CreateUser(ctx context.Context, email, password, displayName, phoneNumber string) (*CreateUserResponse, error)
	ValidateToken(ctx context.Context) (*ValidateTokenResponse, error)
	UpdateProfile(ctx context.Context, displayName, phoneNumber string) (*UpdateProfileResponse, error)

	// Auction methods
	CreateAuction(ctx context.Context, auctionData AuctionData) (*AuctionResponse, error)
	UpdateAuction(ctx context.Context, auctionID string, updateData interface{}) error
	GetAuction(ctx context.Context, auctionID string) (map[string]interface{}, error)

	// Payment methods
	ProcessPayment(ctx context.Context, paymentData interface{}) error

	// Notification methods
	SendNotification(ctx context.Context, userID, title, body string, data map[string]string) (*NotificationResponse, error)
}