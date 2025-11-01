package security

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	"go.uber.org/zap"
)

// ComplianceManager handles compliance-related functionality
type ComplianceManager struct {
	logger        *zap.Logger
	encryptionKey []byte
	piiMasker     *PIIMasker
	rateLimiter   *RateLimiter
	auditLogger   *AuditLogger
}

// NewComplianceManager creates a new compliance manager
func NewComplianceManager(encryptionKey string, logger *zap.Logger) (*ComplianceManager, error) {
	// Validate encryption key
	if len(encryptionKey) != 32 {
		return nil, fmt.Errorf("encryption key must be 32 characters long")
	}

	keyBytes := []byte(encryptionKey)

	return &ComplianceManager{
		logger:        logger,
		encryptionKey: keyBytes,
		piiMasker:     NewPIIMasker(),
		rateLimiter:   NewRateLimiter(),
		auditLogger:   NewAuditLogger(logger),
	}, nil
}

// PIIMasker handles Personally Identifiable Information masking
type PIIMasker struct {
	emailRegex   *regexp.Regexp
	phoneRegex   *regexp.Regexp
	cardRegex    *regexp.Regexp
	accountRegex *regexp.Regexp
}

// NewPIIMasker creates a new PII masker
func NewPIIMasker() *PIIMasker {
	return &PIIMasker{
		emailRegex:   regexp.MustCompile(`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`),
		phoneRegex:   regexp.MustCompile(`\b(?:\+?60|0)?[1-9]\d{1,2}\d{6,8}\b`),
		cardRegex:    regexp.MustCompile(`\b\d{4}[\s-]?\d{4}[\s-]?\d{4}[\s-]?\d{4}\b`),
		accountRegex: regexp.MustCompile(`\b\d{10,20}\b`),
	}
}

// MaskPII masks personally identifiable information in the input string
func (pm *PIIMasker) MaskPII(input string) string {
	result := input

	// Mask email addresses
	result = pm.emailRegex.ReplaceAllStringFunc(result, func(email string) string {
		return pm.maskEmail(email)
	})

	// Mask phone numbers
	result = pm.phoneRegex.ReplaceAllStringFunc(result, func(phone string) string {
		return pm.maskPhone(phone)
	})

	// Mask credit card numbers
	result = pm.cardRegex.ReplaceAllStringFunc(result, func(card string) string {
		return pm.maskCard(card)
	})

	// Mask account numbers (10+ digits)
	result = pm.accountRegex.ReplaceAllStringFunc(result, func(account string) string {
		return pm.maskAccount(account)
	})

	return result
}

// maskEmail masks email address showing only first character and domain
func (pm *PIIMasker) maskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "***MASKED***"
	}

	local := parts[0]
	domain := parts[1]

	if len(local) <= 1 {
		return fmt.Sprintf("*@%s", domain)
	}

	return fmt.Sprintf("%c***@%s", local[0], domain)
}

// maskPhone masks phone number showing only last 4 digits
func (pm *PIIMasker) maskPhone(phone string) string {
	// Remove non-digit characters
	digits := regexp.MustCompile(`\D`).ReplaceAllString(phone, "")

	if len(digits) <= 4 {
		return strings.Repeat("*", len(digits))
	}

	visibleDigits := digits[len(digits)-4:]
	maskedDigits := strings.Repeat("*", len(digits)-4)
	return maskedDigits + visibleDigits
}

// maskCard masks credit card number showing only last 4 digits
func (pm *PIIMasker) maskCard(card string) string {
	// Remove non-digit characters
	digits := regexp.MustCompile(`\D`).ReplaceAllString(card, "")

	if len(digits) != 16 {
		return strings.Repeat("*", len(card))
	}

	visibleDigits := digits[len(digits)-4:]
	maskedDigits := strings.Repeat("*", len(digits)-4)

	// Format with dashes like original
	return fmt.Sprintf("%s-%s", maskedDigits[0:4]+"-"+maskedDigits[4:8]+"-"+maskedDigits[8:12], visibleDigits)
}

// maskAccount masks account number showing only last 4 digits
func (pm *PIIMasker) maskAccount(account string) string {
	digits := regexp.MustCompile(`\D`).ReplaceAllString(account, "")

	if len(digits) <= 4 {
		return strings.Repeat("*", len(digits))
	}

	visibleDigits := digits[len(digits)-4:]
	maskedDigits := strings.Repeat("*", len(digits)-4)
	return maskedDigits + visibleDigits
}

// EncryptSensitiveData encrypts sensitive data using AES-GCM
func (cm *ComplianceManager) EncryptSensitiveData(plaintext string) (string, error) {
	block, err := aes.NewCipher(cm.encryptionKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptSensitiveData decrypts sensitive data using AES-GCM
func (cm *ComplianceManager) DecryptSensitiveData(ciphertext string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	block, err := aes.NewCipher(cm.encryptionKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext_bytes := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext_bytes, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %w", err)
	}

	return string(plaintext), nil
}

// PCIComplianceValidator validates PCI DSS compliance
type PCIComplianceValidator struct {
	logger *zap.Logger
}

// NewPCIComplianceValidator creates a new PCI compliance validator
func NewPCIComplianceValidator(logger *zap.Logger) *PCIComplianceValidator {
	return &PCIComplianceValidator{
		logger: logger,
	}
}

// ValidatePaymentData validates payment data against PCI DSS requirements
func (pci *PCIComplianceValidator) ValidatePaymentData(data PaymentData) *ComplianceResult {
	result := &ComplianceResult{
		Compliant:  true,
		Violations: []string{},
		Category:   "PCI_DSS",
	}

	// Validate credit card data
	if data.CardNumber != "" {
		if !pci.isValidCardNumber(data.CardNumber) {
			result.Compliant = false
			result.Violations = append(result.Violations, "Invalid credit card number format")
		}

		if pci.isFullCardNumberStored(data.CardNumber) {
			result.Compliant = false
			result.Violations = append(result.Violations, "Full credit card number cannot be stored")
		}
	}

	// Validate CVV handling
	if data.CVV != "" {
		if len(data.CVV) < 3 || len(data.CVV) > 4 {
			result.Compliant = false
			result.Violations = append(result.Violations, "Invalid CVV length")
		}

		// CVV should never be stored after authorization
		if data.AfterAuthorization && data.CVV != "" {
			result.Compliant = false
			result.Violations = append(result.Violations, "CVV cannot be stored after authorization")
		}
	}

	// Validate expiry date
	if data.ExpiryMonth < 1 || data.ExpiryMonth > 12 {
		result.Compliant = false
		result.Violations = append(result.Violations, "Invalid expiry month")
	}

	if data.ExpiryYear < time.Now().Year() {
		result.Compliant = false
		result.Violations = append(result.Violations, "Card has expired")
	}

	// Validate encryption
	if data.RequiresEncryption && !data.IsEncrypted {
		result.Compliant = false
		result.Violations = append(result.Violations, "Sensitive data must be encrypted")
	}

	return result
}

// isValidCardNumber checks if the card number passes Luhn algorithm
func (pci *PCIComplianceValidator) isValidCardNumber(number string) bool {
	// Remove non-digit characters
	digits := regexp.MustCompile(`\D`).ReplaceAllString(number, "")

	if len(digits) < 13 || len(digits) > 19 {
		return false
	}

	// Luhn algorithm
	sum := 0
	alternate := false

	for i := len(digits) - 1; i >= 0; i-- {
		n := int(digits[i] - '0')

		if alternate {
			n *= 2
			if n > 9 {
				n = (n % 10) + 1
			}
		}

		sum += n
		alternate = !alternate
	}

	return sum%10 == 0
}

// isFullCardNumberStored checks if full card number is being stored
func (pci *PCIComplianceValidator) isFullCardNumberStored(number string) bool {
	// Remove non-digit characters
	digits := regexp.MustCompile(`\D`).ReplaceAllString(number, "")

	// If more than last 4 digits are visible, it's a violation
	return len(digits) > 4 && !strings.Contains(number, "*")
}

// DataRetentionManager handles data retention policies
type DataRetentionManager struct {
	logger *zap.Logger
}

// NewDataRetentionManager creates a new data retention manager
func NewDataRetentionManager(logger *zap.Logger) *DataRetentionManager {
	return &DataRetentionManager{
		logger: logger,
	}
}

// GetRetentionPeriod returns the retention period for different data types
func (drm *DataRetentionManager) GetRetentionPeriod(dataType string) time.Duration {
	switch dataType {
	case "payment_data":
		return 365 * 24 * time.Hour // 1 year
	case "personal_data":
		return 2555 * 24 * time.Hour // 7 years
	case "audit_logs":
		return 2555 * 24 * time.Hour // 7 years
	case "transaction_data":
		return 2555 * 24 * time.Hour // 7 years
	case "piidata":
		return 365 * 24 * time.Hour // 1 year
	default:
		return 365 * 24 * time.Hour // Default 1 year
	}
}

// ShouldRetain checks if data should still be retained
func (drm *DataRetentionManager) ShouldRetain(dataType string, createdAt time.Time) bool {
	retentionPeriod := drm.GetRetentionPeriod(dataType)
	expiryTime := createdAt.Add(retentionPeriod)
	return time.Now().Before(expiryTime)
}

// RateLimiter implements rate limiting for security
type RateLimiter struct {
	// Implementation would depend on the chosen rate limiting strategy
	// This is a placeholder for demonstration
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{}
}

// CheckLimit checks if the request is within rate limits
func (rl *RateLimiter) CheckLimit(userID, operation string) bool {
	// Implementation would check against Redis or in-memory store
	// For now, return true (no limiting)
	return true
}

// AuditLogger logs security events for compliance
type AuditLogger struct {
	logger *zap.Logger
}

// NewAuditLogger creates a new audit logger
func NewAuditLogger(logger *zap.Logger) *AuditLogger {
	return &AuditLogger{
		logger: logger,
	}
}

// LogAccess logs data access events
func (al *AuditLogger) LogAccess(ctx context.Context, event AuditEvent) {
	fields := []zap.Field{
		zap.String("event_type", "data_access"),
		zap.String("user_id", event.UserID),
		zap.String("resource", event.Resource),
		zap.String("action", event.Action),
		zap.String("ip_address", event.IPAddress),
		zap.Time("timestamp", event.Timestamp),
		zap.Bool("success", event.Success),
	}

	if event.Error != nil {
		fields = append(fields, zap.Error(event.Error))
	}

	if event.Metadata != nil {
		fields = append(fields, zap.Any("metadata", event.Metadata))
	}

	al.logger.Info("Audit log: Data access", fields...)
}

// LogModification logs data modification events
func (al *AuditLogger) LogModification(ctx context.Context, event AuditEvent) {
	fields := []zap.Field{
		zap.String("event_type", "data_modification"),
		zap.String("user_id", event.UserID),
		zap.String("resource", event.Resource),
		zap.String("action", event.Action),
		zap.String("ip_address", event.IPAddress),
		zap.Time("timestamp", event.Timestamp),
		zap.Bool("success", event.Success),
	}

	if event.BeforeState != nil {
		fields = append(fields, zap.Any("before_state", event.BeforeState))
	}

	if event.AfterState != nil {
		fields = append(fields, zap.Any("after_state", event.AfterState))
	}

	if event.Error != nil {
		fields = append(fields, zap.Error(event.Error))
	}

	al.logger.Info("Audit log: Data modification", fields...)
}

// LogSecurityEvent logs security-related events
func (al *AuditLogger) LogSecurityEvent(ctx context.Context, event SecurityEvent) {
	fields := []zap.Field{
		zap.String("event_type", "security_event"),
		zap.String("event_category", event.Category),
		zap.String("severity", event.Severity),
		zap.String("description", event.Description),
		zap.String("source_ip", event.SourceIP),
		zap.String("user_agent", event.UserAgent),
		zap.Time("timestamp", event.Timestamp),
	}

	if event.UserID != "" {
		fields = append(fields, zap.String("user_id", event.UserID))
	}

	if event.Details != nil {
		fields = append(fields, zap.Any("details", event.Details))
	}

	al.logger.Warn("Audit log: Security event", fields...)
}

// Data structures

type ComplianceResult struct {
	Compliant  bool     `json:"compliant"`
	Violations []string `json:"violations"`
	Category   string   `json:"category"`
}

type PaymentData struct {
	CardNumber         string `json:"card_number"`
	CVV                string `json:"cvv"`
	ExpiryMonth        int    `json:"expiry_month"`
	ExpiryYear         int    `json:"expiry_year"`
	AfterAuthorization bool   `json:"after_authorization"`
	RequiresEncryption bool   `json:"requires_encryption"`
	IsEncrypted        bool   `json:"is_encrypted"`
}

type AuditEvent struct {
	UserID      string      `json:"user_id"`
	Resource    string      `json:"resource"`
	Action      string      `json:"action"`
	IPAddress   string      `json:"ip_address"`
	Timestamp   time.Time   `json:"timestamp"`
	Success     bool        `json:"success"`
	Error       error       `json:"error,omitempty"`
	Metadata    interface{} `json:"metadata,omitempty"`
	BeforeState interface{} `json:"before_state,omitempty"`
	AfterState  interface{} `json:"after_state,omitempty"`
}

type SecurityEvent struct {
	Category    string      `json:"category"`
	Severity    string      `json:"severity"`
	Description string      `json:"description"`
	SourceIP    string      `json:"source_ip"`
	UserAgent   string      `json:"user_agent"`
	UserID      string      `json:"user_id,omitempty"`
	Details     interface{} `json:"details,omitempty"`
	Timestamp   time.Time   `json:"timestamp"`
}

// ComplianceChecker interface for dependency injection
type ComplianceChecker interface {
	ValidatePaymentData(data PaymentData) *ComplianceResult
}

// Ensure PCIComplianceValidator implements ComplianceChecker
var _ ComplianceChecker = (*PCIComplianceValidator)(nil)
