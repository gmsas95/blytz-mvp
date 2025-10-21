package services

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/gmsas95/blytz-mvp/services/logistics-service/internal/config"
	"github.com/gmsas95/blytz-mvp/services/logistics-service/internal/models"
	"github.com/gmsas95/blytz-mvp/shared/pkg/errors"
)

type LogisticsService struct {
	db     *gorm.DB
	logger *zap.Logger
	config *config.Config
}

func NewLogisticsService(db *gorm.DB, logger *zap.Logger, config *config.Config) *LogisticsService {
	return &LogisticsService{
		db:     db,
		logger: logger,
		config: config,
	}
}

func (s *LogisticsService) CreateShipment(ctx context.Context, userID string, req *CreateShipmentRequest) (*models.Shipment, error) {
	s.logger.Info("Creating shipment", zap.String("user_id", userID), zap.String("order_id", req.OrderID))

	// Validate request
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Generate tracking number
	trackingNumber := s.generateTrackingNumber()

	// Create shipment
	shipment := &models.Shipment{
		OrderID:        req.OrderID,
		UserID:         userID,
		TrackingNumber: trackingNumber,
		Carrier:        req.Carrier,
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
		Weight:     req.Weight,
		Dimensions: models.Dimensions{
			Length: req.Dimensions.Length,
			Width:  req.Dimensions.Width,
			Height: req.Dimensions.Height,
		},
		Cost:  req.Cost,
		Notes: req.Notes,
	}

	if err := s.db.Create(shipment).Error; err != nil {
		s.logger.Error("Failed to create shipment", zap.Error(err))
		return nil, errors.ErrInternalServer
	}

	s.logger.Info("Shipment created successfully", zap.String("shipment_id", shipment.ID))
	return shipment, nil
}

func (s *LogisticsService) GetShipment(ctx context.Context, shipmentID string, userID string) (*models.Shipment, error) {
	s.logger.Info("Getting shipment", zap.String("shipment_id", shipmentID), zap.String("user_id", userID))

	var shipment models.Shipment
	if err := s.db.Where("id = ? AND user_id = ?", shipmentID, userID).First(&shipment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		s.logger.Error("Failed to get shipment", zap.Error(err))
		return nil, errors.ErrInternalServer
	}

	return &shipment, nil
}

func (s *LogisticsService) GetShipmentByOrder(ctx context.Context, orderID string, userID string) (*models.Shipment, error) {
	s.logger.Info("Getting shipment by order", zap.String("order_id", orderID), zap.String("user_id", userID))

	var shipment models.Shipment
	if err := s.db.Where("order_id = ? AND user_id = ?", orderID, userID).First(&shipment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		s.logger.Error("Failed to get shipment by order", zap.Error(err))
		return nil, errors.ErrInternalServer
	}

	return &shipment, nil
}

func (s *LogisticsService) UpdateShipmentStatus(ctx context.Context, shipmentID string, userID string, status models.ShipmentStatus) (*models.Shipment, error) {
	s.logger.Info("Updating shipment status", zap.String("shipment_id", shipmentID), zap.String("user_id", userID), zap.String("status", string(status)))

	var shipment models.Shipment
	if err := s.db.Where("id = ? AND user_id = ?", shipmentID, userID).First(&shipment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		s.logger.Error("Failed to get shipment for status update", zap.Error(err))
		return nil, errors.ErrInternalServer
	}

	shipment.Status = string(status)
	shipment.UpdatedAt = time.Now()

	// Set delivery time if delivered
	if status == models.ShipmentStatusDelivered {
		now := time.Now()
		shipment.ActualDelivery = &now
	}

	if err := s.db.Save(&shipment).Error; err != nil {
		s.logger.Error("Failed to update shipment status", zap.Error(err))
		return nil, errors.ErrInternalServer
	}

	s.logger.Info("Shipment status updated successfully", zap.String("shipment_id", shipment.ID))
	return &shipment, nil
}

func (s *LogisticsService) TrackShipment(ctx context.Context, trackingNumber string) (*models.Shipment, []*models.TrackingEvent, error) {
	s.logger.Info("Tracking shipment", zap.String("tracking_number", trackingNumber))

	var shipment models.Shipment
	if err := s.db.Where("tracking_number = ?", trackingNumber).First(&shipment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil, errors.ErrNotFound
		}
		s.logger.Error("Failed to get shipment for tracking", zap.Error(err))
		return nil, nil, errors.ErrInternalServer
	}

	// Get tracking events
	var events []*models.TrackingEvent
	if err := s.db.Where("shipment_id = ?", shipment.ID).Order("timestamp DESC").Find(&events).Error; err != nil {
		s.logger.Error("Failed to get tracking events", zap.Error(err))
		return nil, nil, errors.ErrInternalServer
	}

	return &shipment, events, nil
}

func (s *LogisticsService) generateTrackingNumber() string {
	return fmt.Sprintf("BLTZ%010d", time.Now().UnixNano()%10000000000)
}