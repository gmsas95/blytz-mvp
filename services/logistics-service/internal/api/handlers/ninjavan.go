package handlers

import (
	"io"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/gmsas95/blytz-mvp/services/logistics-service/internal/models"
	"github.com/gmsas95/blytz-mvp/services/logistics-service/internal/services"
	"github.com/gmsas95/blytz-mvp/shared/pkg/errors"
	"github.com/gmsas95/blytz-mvp/shared/pkg/utils"
)

type NinjaVanHandler struct {
	ninjaVanService *services.NinjaVanService
	logger          *zap.Logger
}

func NewNinjaVanHandler(ninjaVanService *services.NinjaVanService, logger *zap.Logger) *NinjaVanHandler {
	return &NinjaVanHandler{
		ninjaVanService: ninjaVanService,
		logger:          logger,
	}
}

func (h *NinjaVanHandler) CreateNinjaVanShipment(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		utils.ErrorResponse(c, errors.ErrUnauthorized)
		return
	}

	var req services.CreateShipmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, errors.ErrInvalidRequestBody)
		return
	}

	// Force carrier to be ninja_van for this endpoint
	req.Carrier = "ninja_van"

	shipment, err := h.ninjaVanService.CreateNinjaVanShipment(c.Request.Context(), userID, &req)
	if err != nil {
		h.logger.Error("Failed to create Ninja Van shipment", zap.Error(err))
		utils.ErrorResponse(c, err)
		return
	}

	response := h.mapShipmentToResponse(shipment)
	utils.SuccessResponse(c, response)
}

func (h *NinjaVanHandler) CancelNinjaVanShipment(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		utils.ErrorResponse(c, errors.ErrUnauthorized)
		return
	}

	shipmentID := c.Param("id")
	if shipmentID == "" {
		utils.ErrorResponse(c, errors.ErrInvalidRequest)
		return
	}

	err := h.ninjaVanService.CancelNinjaVanShipment(c.Request.Context(), shipmentID, userID)
	if err != nil {
		if err == errors.ErrNotFound {
			utils.ErrorResponse(c, errors.ErrNotFound)
			return
		}
		h.logger.Error("Failed to cancel Ninja Van shipment", zap.Error(err))
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, gin.H{
		"success": true,
		"message": "Shipment cancelled successfully",
	})
}

func (h *NinjaVanHandler) GetShippingCost(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		utils.ErrorResponse(c, errors.ErrUnauthorized)
		return
	}

	var req services.CreateShipmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, errors.ErrInvalidRequestBody)
		return
	}

	tariff, err := h.ninjaVanService.GetShippingCost(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to get shipping cost", zap.Error(err))
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, tariff)
}

func (h *NinjaVanHandler) GetPUDOPoints(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		utils.ErrorResponse(c, errors.ErrUnauthorized)
		return
	}

	pudoPoints, err := h.ninjaVanService.GetPUDOPoints(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to get PUDO points", zap.Error(err))
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, gin.H{
		"success": true,
		"data":    pudoPoints,
	})
}

func (h *NinjaVanHandler) ProcessWebhook(c *gin.Context) {
	// Get webhook signature from header
	signature := c.GetHeader("X-Ninjavan-Hmac-Sha256")
	if signature == "" {
		utils.ErrorResponse(c, errors.ErrInvalidRequest)
		return
	}

	// Read webhook body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		h.logger.Error("Failed to read webhook body", zap.Error(err))
		utils.ErrorResponse(c, errors.ErrInvalidRequest)
		return
	}

	// Process webhook
	err = h.ninjaVanService.ProcessWebhook(c.Request.Context(), body, signature)
	if err != nil {
		h.logger.Error("Failed to process webhook", zap.Error(err))
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, gin.H{
		"success": true,
		"message": "Webhook processed successfully",
	})
}

func (h *NinjaVanHandler) mapShipmentToResponse(shipment *models.Shipment) *services.ShipmentResponse {
	// Map models.Shipment to services.ShipmentResponse
	var estimatedDelivery, actualDelivery *string
	if shipment.EstimatedDelivery != nil {
		ed := shipment.EstimatedDelivery.Format("2006-01-02T15:04:05Z07:00")
		estimatedDelivery = &ed
	}
	if shipment.ActualDelivery != nil {
		ad := shipment.ActualDelivery.Format("2006-01-02T15:04:05Z07:00")
		actualDelivery = &ad
	}

	return &services.ShipmentResponse{
		ID:             shipment.ID,
		OrderID:        shipment.OrderID,
		UserID:         shipment.UserID,
		TrackingNumber: shipment.TrackingNumber,
		Carrier:        shipment.Carrier,
		Service:        shipment.Service,
		Status:         shipment.Status,
		OriginAddress: services.AddressResponse{
			Name:        shipment.OriginAddress.Name,
			Street:      shipment.OriginAddress.Street,
			City:        shipment.OriginAddress.City,
			State:       shipment.OriginAddress.State,
			PostalCode:  shipment.OriginAddress.PostalCode,
			Country:     shipment.OriginAddress.Country,
			PhoneNumber: shipment.OriginAddress.PhoneNumber,
		},
		DestinationAddress: services.AddressResponse{
			Name:        shipment.DestinationAddress.Name,
			Street:      shipment.DestinationAddress.Street,
			City:        shipment.DestinationAddress.City,
			State:       shipment.DestinationAddress.State,
			PostalCode:  shipment.DestinationAddress.PostalCode,
			Country:     shipment.DestinationAddress.Country,
			PhoneNumber: shipment.DestinationAddress.PhoneNumber,
		},
		Weight: shipment.Weight,
		Dimensions: services.DimensionsResponse{
			Length: shipment.Dimensions.Length,
			Width:  shipment.Dimensions.Width,
			Height: shipment.Dimensions.Height,
		},
		EstimatedDelivery: estimatedDelivery,
		ActualDelivery:    actualDelivery,
		Cost:              shipment.Cost,
		Notes:             shipment.Notes,
		CreatedAt:         shipment.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:         shipment.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
