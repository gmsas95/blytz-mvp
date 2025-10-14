package constants

// Service ports
const (
	AuthServicePort     = "8081"
	ProductServicePort  = "8082"
	AuctionServicePort  = "8083"
	ChatServicePort     = "8084"
	OrderServicePort    = "8085"
	PaymentServicePort  = "8086"
	LogisticsServicePort = "8087"
	GatewayPort         = "8080"
)

// Service names
const (
	AuthService     = "auth-service"
	ProductService  = "product-service"
	AuctionService  = "auction-service"
	ChatService     = "chat-service"
	OrderService    = "order-service"
	PaymentService  = "payment-service"
	LogisticsService = "logistics-service"
	GatewayService  = "gateway-service"
)

// Database names
const (
	AuthDB     = "auth"
	ProductDB  = "products"
	AuctionDB  = "auctions"
	ChatDB     = "chat"
	OrderDB    = "orders"
	PaymentDB  = "payments"
	LogisticsDB = "logistics"
)

// Redis prefixes
const (
	RedisAuthPrefix     = "auth:"
	RedisProductPrefix  = "product:"
	RedisAuctionPrefix  = "auction:"
	RedisChatPrefix     = "chat:"
	RedisOrderPrefix    = "order:"
	RedisPaymentPrefix  = "payment:"
	RedisLogisticsPrefix = "logistics:"
)

// JWT constants
const (
	JWTSecretEnv     = "JWT_SECRET"
	JWTExpiryEnv     = "JWT_EXPIRY"
	DefaultJWTExpiry = 24 * 60 * 60 // 24 hours in seconds
)

// Firebase constants
const (
	FirebaseProjectIDEnv = "FIREBASE_PROJECT_ID"
	FirebaseServiceAccountEnv = "FIREBASE_SERVICE_ACCOUNT"
)

// Stripe constants
const (
	StripeSecretKeyEnv   = "STRIPE_SECRET_KEY"
	StripePublishKeyEnv  = "STRIPE_PUBLISH_KEY"
	StripeWebhookSecretEnv = "STRIPE_WEBHOOK_SECRET"
)

// Ninja Van constants
const (
	NinjaVanAPIKeyEnv    = "NINJA_VAN_API_KEY"
	NinjaVanAPISecretEnv = "NINJA_VAN_API_SECRET"
	NinjaVanBaseURLEnv   = "NINJA_VAN_BASE_URL"
)

// Common HTTP headers
const (
	HeaderAuthorization = "Authorization"
	HeaderContentType   = "Content-Type"
	HeaderXRequestID    = "X-Request-ID"
	HeaderXUserID       = "X-User-ID"
	HeaderXRole         = "X-Role"
)

// Content types
const (
	ContentTypeJSON        = "application/json"
	ContentTypeForm        = "application/x-www-form-urlencoded"
	ContentTypeMultipart   = "multipart/form-data"
	ContentTypeText        = "text/plain"
	ContentTypeSSE         = "text/event-stream"
	ContentTypeWebSocket   = "application/websocket"
)

// User roles
const (
	RoleAdmin   = "admin"
	RoleSeller  = "seller"
	RoleBuyer   = "buyer"
	RoleUser    = "user"
)

// Order statuses
const (
	OrderStatusPending    = "pending"
	OrderStatusProcessing = "processing"
	OrderStatusConfirmed  = "confirmed"
	OrderStatusShipped    = "shipped"
	OrderStatusDelivered  = "delivered"
	OrderStatusCancelled  = "cancelled"
	OrderStatusRefunded   = "refunded"
)

// Payment statuses
const (
	PaymentStatusPending    = "pending"
	PaymentStatusProcessing = "processing"
	PaymentStatusSucceeded  = "succeeded"
	PaymentStatusFailed     = "failed"
	PaymentStatusCancelled  = "cancelled"
	PaymentStatusRefunded   = "refunded"
)

// Auction statuses
const (
	AuctionStatusDraft     = "draft"
	AuctionStatusScheduled = "scheduled"
	AuctionStatusActive    = "active"
	AuctionStatusEnded     = "ended"
	AuctionStatusCancelled = "cancelled"
)

// Chat message types
const (
	ChatTypeText     = "text"
	ChatTypeImage    = "image"
	ChatTypeBid      = "bid"
	ChatTypeSystem   = "system"
	ChatTypeJoin     = "join"
	ChatTypeLeave    = "leave"
)

// Logistics statuses
const (
	LogisticsStatusPending      = "pending"
	LogisticsStatusPickedUp     = "picked_up"
	LogisticsStatusInTransit    = "in_transit"
	LogisticsStatusOutForDelivery = "out_for_delivery"
	LogisticsStatusDelivered    = "delivered"
	LogisticsStatusFailed       = "failed"
)

// Environment variables
const (
	EnvDatabaseURL = "DATABASE_URL"
	EnvRedisURL    = "REDIS_URL"
	EnvPort        = "PORT"
	EnvEnvironment = "GO_ENV"
	EnvLogLevel    = "LOG_LEVEL"
)