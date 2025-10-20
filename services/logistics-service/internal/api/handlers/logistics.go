package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/gmsas95/blytz-mvp/services/logistics-service/internal/models"
	"github.com/gmsas95/blytz-mvp/services/logistics-service/internal/services"
	"github.com/gmsas95/blytz-mvp/shared/pkg/errors"
	"github.com/gmsas95/blytz-mvp/shared/pkg/utils"
)

type LogisticsHandler struct {
	logisticsService *services.LogisticsService
	logger           *zap.Logger
}

func NewLogisticsHandler(logisticsService *services.LogisticsService, logger *zap.Logger) *LogisticsHandler {
	return &LogisticsHandler{
		logisticsService: logisticsService,
		logger:           logger,
	}
}

func (h *LogisticsHandler) CreateShipment(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		utils.RespondWithError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req services.CreateShipmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	shipment, err := h.logisticsService.CreateShipment(c.Request.Context(), userID, &req)
	if err != nil {
		h.logger.Error("Failed to create shipment", zap.Error(err))
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to create shipment")
		return
	}

	response := h.mapShipmentToResponse(shipment)
	utils.RespondWithJSON(c, http.StatusCreated, response)
}

func (h *LogisticsHandler) GetShipment(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		utils.RespondWithError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	shipmentID := c.Param("id")
	if shipmentID == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Shipment ID is required")
		return
	}

	shipment, err := h.logisticsService.GetShipment(c.Request.Context(), shipmentID, userID)
	if err != nil {
		if err == errors.ErrNotFound {
			utils.RespondWithError(c, http.StatusNotFound, "Shipment not found")
			return
		}
		h.logger.Error("Failed to get shipment", zap.Error(err))
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get shipment")
		return
	}

	response := h.mapShipmentToResponse(shipment)
	utils.RespondWithJSON(c, http.StatusOK, response)
}

func (h *LogisticsHandler) UpdateShipmentStatus(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		utils.RespondWithError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	shipmentID := c.Param("id")
	if shipmentID == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Shipment ID is required")
		return
	}

	var req services.UpdateShipmentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	shipment, err := h.logisticsService.UpdateShipmentStatus(c.Request.Context(), shipmentID, userID, models.ShipmentStatus(req.Status))
	if err != nil {
		if err == errors.ErrNotFound {
			utils.RespondWithError(c, http.StatusNotFound, "Shipment not found")
			return
		}
		h.logger.Error("Failed to update shipment status", zap.Error(err))
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to update shipment status")
		return
	}

	response := h.mapShipmentToResponse(shipment)
	utils.RespondWithJSON(c, http.StatusOK, response)
}

func (h *LogisticsHandler) GetShipmentByOrder(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		utils.RespondWithError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	orderID := c.Param("orderId")
	if orderID == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Order ID is required")
		return
	}

	shipment, err := h.logisticsService.GetShipmentByOrder(c.Request.Context(), orderID, userID)
	if err != nil {
		if err == errors.ErrNotFound {
			utils.RespondWithError(c, http.StatusNotFound, "Shipment not found for order")
			return
		}
		h.logger.Error("Failed to get shipment by order", zap.Error(err))
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to get shipment")
		return
	}

	response := h.mapShipmentToResponse(shipment)
	utils.RespondWithJSON(c, http.StatusOK, response)
}

func (h *LogisticsHandler) TrackShipment(c *gin.Context) {
	trackingNumber := c.Param("trackingNumber")
	if trackingNumber == "" {
		utils.RespondWithError(c, http.StatusBadRequest, "Tracking number is required")
		return
	}

	shipment, events, err := h.logisticsService.TrackShipment(c.Request.Context(), trackingNumber)
	if err != nil {
		if err == errors.ErrNotFound {
			utils.RespondWithError(c, http.StatusNotFound, "Shipment not found")
			return
		}
		h.logger.Error("Failed to track shipment", zap.Error(err))
		utils.RespondWithError(c, http.StatusInternalServerError, "Failed to track shipment")
		return
	}

	shipmentResponse := h.mapShipmentToResponse(shipment)
	eventResponses := make([]services.TrackingEventResponse, len(events))
	for i, event := range events {
		eventResponses[i] = services.TrackingEventResponse{
			ID:          event.ID,
			Status:      event.Status,
			Description: event.Description,
			Location:    event.Location,
			Timestamp:   event.Timestamp.Format("2006-01-02T15:04:05Z"),
		}
	}

	response := services.TrackingResponse{
		Shipment: shipmentResponse,
		Events:   eventResponses,
	}

	utils.RespondWithJSON(c, http.StatusOK, response)
}

func (h *LogisticsHandler) mapShipmentToResponse(shipment *models.Shipment) *services.ShipmentResponse {
	var estimatedDelivery, actualDelivery *string
	if shipment.EstimatedDelivery != nil {
		str := shipment.EstimatedDelivery.Format("2006-01-02T15:04:05Z")
		estimatedDelivery = &str
	}
	if shipment.ActualDelivery != nil {
		str := shipment.ActualDelivery.Format("2006-01-02T15:04:05Z")
		actualDelivery = &str
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
		Weight:            shipment.Weight,
		Dimensions: services.DimensionsResponse{
			Length: shipment.Dimensions.Length,
			Width:  shipment.Dimensions.Width,
			Height: shipment.Dimensions.Height,
		},
		EstimatedDelivery: estimatedDelivery,
		ActualDelivery:    actualDelivery,
		Cost:              shipment.Cost,
		Notes:             shipment.Notes,
		CreatedAt:         shipment.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:         shipment.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}