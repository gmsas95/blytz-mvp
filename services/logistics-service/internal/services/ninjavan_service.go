package services

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/gmsas95/blytz-mvp/services/logistics-service/internal/config"
	"github.com/gmsas95/blytz-mvp/services/logistics-service/internal/models"
	"github.com/gmsas95/blytz-mvp/shared/pkg/errors"
)

type NinjaVanService struct {
	db             *gorm.DB
	logger         *zap.Logger
	config         *config.Config
	ninjaVanClient *NinjaVanClient
}

func NewNinjaVanService(db *gorm.DB, logger *zap.Logger, config *config.Config) *NinjaVanService {
	ninjaVanConfig := &NinjaVanConfig{
		ClientID:    config.NinjaVanClientID,
		ClientKey:   config.NinjaVanClientKey, // Client Key for OAuth authentication
		Environment: config.NinjaVanEnvironment,
		CountryCode: config.NinjaVanCountryCode,
	}

	return &NinjaVanService{
		db:             db,
		logger:         logger,
		config:         config,
		ninjaVanClient: NewNinjaVanClient(logger, ninjaVanConfig),
	}
}

func (s *NinjaVanService) CreateNinjaVanShipment(ctx context.Context, userID string, req *CreateShipmentRequest) (*models.Shipment, error) {
	s.logger.Info("Creating Ninja Van shipment", zap.String("user_id", userID), zap.String("order_id", req.OrderID))

	// Validate request
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Convert to Ninja Van order request
	ninjaVanOrder := &NinjaVanOrderRequest{
		TrackingNumber:     fmt.Sprintf("BLTZ%d", time.Now().Unix()),
		ServiceType:        "parcel",
		ServiceLevel:       req.Service,
		OriginAddress:      s.ninjaVanClient.ConvertToNinjaVanAddress(req.OriginAddress),
		DestinationAddress: s.ninjaVanClient.ConvertToNinjaVanAddress(req.DestinationAddress),
		Parcel:             s.ninjaVanClient.ConvertToNinjaVanParcel(req.Weight, req.Dimensions),
		IsCOD:              false,
		Remarks:            req.Notes,
	}

	// Create order with Ninja Van
	ninjaVanResp, err := s.ninjaVanClient.CreateOrder(ctx, ninjaVanOrder)
	if err != nil {
		s.logger.Error("Failed to create Ninja Van order", zap.Error(err))
		return nil, fmt.Errorf("failed to create Ninja Van order: %w", err)
	}

	// Create internal shipment record
	shipment := &models.Shipment{
		OrderID:        req.OrderID,
		UserID:         userID,
		TrackingNumber: ninjaVanResp.TrackingID,
		Carrier:        "ninja_van",
		Service:        req.Service,
		Status:         string(models.ShipmentStatusPending),
		OriginAddress: models.Address{
			Name:        req.OriginAddress.Name,
			Street:      req.OriginAddress.Street,
			City:        req.OriginAddress.City,
			State:       req.OriginAddress.State,
			PostalCode:  req.OriginAddress.PostalCode,
			Country:     req.OriginAddress.Country,
			PhoneNumber: req.OriginAddress.PhoneNumber,
		},
		DestinationAddress: models.Address{
			Name:        req.DestinationAddress.Name,
			Street:      req.DestinationAddress.Street,
			City:        req.DestinationAddress.City,
			State:       req.DestinationAddress.State,
			PostalCode:  req.DestinationAddress.PostalCode,
			Country:     req.DestinationAddress.Country,
			PhoneNumber: req.DestinationAddress.PhoneNumber,
		},
		Weight: req.Weight,
		Dimensions: models.Dimensions{
			Length: req.Dimensions.Length,
			Width:  req.Dimensions.Width,
			Height: req.Dimensions.Height,
		},
		Cost:  req.Cost,
		Notes: req.Notes,
	}

	if err := s.db.Create(shipment).Error; err != nil {
		s.logger.Error("Failed to create shipment record", zap.Error(err))
		return nil, errors.ErrInternalServer
	}

	// Create initial tracking event
	trackingEvent := &models.TrackingEvent{
		ShipmentID:  shipment.ID,
		Status:      string(models.ShipmentStatusPending),
		Description: "Order created and sent to Ninja Van",
		Location:    "System",
		Timestamp:   time.Now(),
	}

	if err := s.db.Create(trackingEvent).Error; err != nil {
		s.logger.Error("Failed to create tracking event", zap.Error(err))
		// Don't fail the whole operation if tracking event creation fails
	}

	s.logger.Info("Ninja Van shipment created successfully",
		zap.String("shipment_id", shipment.ID),
		zap.String("tracking_id", ninjaVanResp.TrackingID))

	return shipment, nil
}

func (s *NinjaVanService) CancelNinjaVanShipment(ctx context.Context, shipmentID string, userID string) error {
	s.logger.Info("Cancelling Ninja Van shipment", zap.String("shipment_id", shipmentID), zap.String("user_id", userID))

	// Get shipment
	var shipment models.Shipment
	if err := s.db.Where("id = ? AND user_id = ?", shipmentID, userID).First(&shipment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.ErrNotFound
		}
		return errors.ErrInternalServer
	}

	// Cancel with Ninja Van
	if err := s.ninjaVanClient.CancelOrder(ctx, shipment.TrackingNumber); err != nil {
		s.logger.Error("Failed to cancel Ninja Van order", zap.Error(err))
		return fmt.Errorf("failed to cancel Ninja Van order: %w", err)
	}

	// Update internal status
	shipment.Status = string(models.ShipmentStatusReturned)
	if err := s.db.Save(&shipment).Error; err != nil {
		s.logger.Error("Failed to update shipment status", zap.Error(err))
		return errors.ErrInternalServer
	}

	// Create tracking event
	trackingEvent := &models.TrackingEvent{
		ShipmentID:  shipment.ID,
		Status:      string(models.ShipmentStatusReturned),
		Description: "Order cancelled with Ninja Van",
		Location:    "System",
		Timestamp:   time.Now(),
	}

	if err := s.db.Create(trackingEvent).Error; err != nil {
		s.logger.Error("Failed to create tracking event", zap.Error(err))
	}

	s.logger.Info("Ninja Van shipment cancelled successfully", zap.String("shipment_id", shipmentID))
	return nil
}

func (s *NinjaVanService) GetShippingCost(ctx context.Context, req *CreateShipmentRequest) (*NinjaVanTariffResponse, error) {
	s.logger.Info("Getting Ninja Van shipping cost", zap.String("order_id", req.OrderID))

	tariffReq := &NinjaVanTariffRequest{
		OriginAddress:      s.ninjaVanClient.ConvertToNinjaVanAddress(req.OriginAddress),
		DestinationAddress: s.ninjaVanClient.ConvertToNinjaVanAddress(req.DestinationAddress),
		Parcel:             s.ninjaVanClient.ConvertToNinjaVanParcel(req.Weight, req.Dimensions),
		ServiceType:        "parcel",
		ServiceLevel:       req.Service,
	}

	tariffResp, err := s.ninjaVanClient.GetTariff(ctx, tariffReq)
	if err != nil {
		s.logger.Error("Failed to get Ninja Van tariff", zap.Error(err))
		return nil, fmt.Errorf("failed to get shipping cost: %w", err)
	}

	return tariffResp, nil
}

func (s *NinjaVanService) GetPUDOPoints(ctx context.Context) ([]NinjaVanPUDOPoint, error) {
	s.logger.Info("Getting Ninja Van PUDO points")

	pudoPoints, err := s.ninjaVanClient.GetPUDOPoints(ctx)
	if err != nil {
		s.logger.Error("Failed to get PUDO points", zap.Error(err))
		return nil, fmt.Errorf("failed to get PUDO points: %w", err)
	}

	return pudoPoints, nil
}

func (s *NinjaVanService) ProcessWebhook(ctx context.Context, webhookData []byte, signature string) error {
	s.logger.Info("Processing Ninja Van webhook")

	// Verify webhook signature
	if !s.verifyWebhookSignature(webhookData, signature) {
		s.logger.Error("Invalid webhook signature")
		return fmt.Errorf("invalid webhook signature")
	}

	// Parse webhook data
	var webhook NinjaVanWebhookPayload
	if err := json.Unmarshal(webhookData, &webhook); err != nil {
		s.logger.Error("Failed to parse webhook", zap.Error(err))
		return fmt.Errorf("failed to parse webhook: %w", err)
	}

	// Find shipment by tracking number
	var shipment models.Shipment
	if err := s.db.Where("tracking_number = ?", webhook.TrackingID).First(&shipment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			s.logger.Warn("Shipment not found for webhook", zap.String("tracking_id", webhook.TrackingID))
			return nil // Don't fail, just log warning
		}
		return fmt.Errorf("failed to find shipment: %w", err)
	}

	// Map Ninja Van status to internal status
	internalStatus := s.mapNinjaVanStatus(webhook.Status)

	// Update shipment status
	shipment.Status = internalStatus
	if internalStatus == string(models.ShipmentStatusDelivered) {
		now := time.Now()
		shipment.ActualDelivery = &now
	}

	if err := s.db.Save(&shipment).Error; err != nil {
		s.logger.Error("Failed to update shipment from webhook", zap.Error(err))
		return fmt.Errorf("failed to update shipment: %w", err)
	}

	// Create tracking event
	trackingEvent := &models.TrackingEvent{
		ShipmentID:  shipment.ID,
		Status:      webhook.Status,
		Description: s.getWebhookDescription(webhook),
		Location:    s.getWebhookLocation(webhook),
		Timestamp:   s.parseWebhookTimestamp(webhook.Timestamp),
	}

	if err := s.db.Create(trackingEvent).Error; err != nil {
		s.logger.Error("Failed to create tracking event from webhook", zap.Error(err))
		// Don't fail the whole operation
	}

	s.logger.Info("Webhook processed successfully",
		zap.String("tracking_id", webhook.TrackingID),
		zap.String("status", webhook.Status))

	return nil
}

func (s *NinjaVanService) verifyWebhookSignature(data []byte, signature string) bool {
	if s.config.NinjaVanClientKey == "" {
		s.logger.Warn("Ninja Van client key not configured, skipping signature verification")
		return true
	}

	h := hmac.New(sha256.New, []byte(s.config.NinjaVanClientKey))
	h.Write(data)
	expectedSignature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return signature == expectedSignature
}

func (s *NinjaVanService) mapNinjaVanStatus(ninjaVanStatus string) string {
	// Map Ninja Van statuses to internal statuses
	switch ninjaVanStatus {
	case "Pending Pickup":
		return string(models.ShipmentStatusPending)
	case "Picked Up", "Picked Up, In Transit to Origin Hub":
		return string(models.ShipmentStatusPickedUp)
	case "In Transit", "In Transit to Next Sorting Hub", "Arrived at Origin Hub", "Arrived at Transit Hub", "Arrived at Destination Hub":
		return string(models.ShipmentStatusInTransit)
	case "On Vehicle for Delivery":
		return string(models.ShipmentStatusOutForDelivery)
	case "Delivered", "Delivered, Received by Customer", "Delivered, Left at Doorstep", "Delivered, Collected by Customer":
		return string(models.ShipmentStatusDelivered)
	case "Cancelled":
		return string(models.ShipmentStatusReturned)
	case "Pickup Exception", "Delivery Exception":
		return string(models.ShipmentStatusFailed)
	default:
		return string(models.ShipmentStatusPending)
	}
}

func (s *NinjaVanService) getWebhookDescription(webhook NinjaVanWebhookPayload) string {
	// Generate description based on webhook event
	switch webhook.Status {
	case "Pending Pickup":
		return "Order is ready for pickup"
	case "Picked Up, In Transit to Origin Hub":
		return "Package picked up and in transit to origin hub"
	case "On Vehicle for Delivery":
		return "Package is out for delivery"
	case "Delivered, Received by Customer":
		return "Package delivered to customer"
	case "Delivered, Left at Doorstep":
		return "Package left at doorstep"
	case "Cancelled":
		return "Order cancelled"
	default:
		return webhook.Status
	}
}

func (s *NinjaVanService) getWebhookLocation(webhook NinjaVanWebhookPayload) string {
	// Extract location from webhook if available
	if webhook.ArrivedAtOriginHubInfo != nil {
		return fmt.Sprintf("%s, %s", webhook.ArrivedAtOriginHubInfo.City, webhook.ArrivedAtOriginHubInfo.Hub)
	}
	if webhook.ArrivedAtDestinationHubInfo != nil {
		return fmt.Sprintf("%s, %s", webhook.ArrivedAtDestinationHubInfo.City, webhook.ArrivedAtDestinationHubInfo.Hub)
	}
	return "Ninja Van Network"
}

func (s *NinjaVanService) parseWebhookTimestamp(timestamp string) time.Time {
	// Parse ISO 8601 timestamp
	parsedTime, err := time.Parse("2006-01-02T15:04:05-0700", timestamp)
	if err != nil {
		s.logger.Error("Failed to parse webhook timestamp", zap.Error(err), zap.String("timestamp", timestamp))
		return time.Now()
	}
	return parsedTime
}
