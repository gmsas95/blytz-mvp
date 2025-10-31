package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/gmsas95/blytz-mvp/services/auction-service/internal/models"
	"github.com/gmsas95/blytz-mvp/services/auction-service/internal/services"
	"github.com/gmsas95/blytz-mvp/services/auction-service/pkg/firebase"
	"github.com/gmsas95/blytz-mvp/shared/pkg/utils"
	shared_errors "github.com/gmsas95/blytz-mvp/shared/pkg/errors"
)


type AuctionHandler struct {
	auctionService *services.AuctionService
	logger         *zap.Logger
	firebaseApp    firebase.FirebaseApp
}

func NewAuctionHandler(auctionService *services.AuctionService, logger *zap.Logger, firebaseApp firebase.FirebaseApp) *AuctionHandler {
	return &AuctionHandler{auctionService: auctionService, logger: logger, firebaseApp: firebaseApp}
}

func (h *AuctionHandler) CreateAuction(c *gin.Context) {
	var auction models.Auction
	if err := c.ShouldBindJSON(&auction); err != nil {
		utils.SendErrorResponse(c, shared_errors.ErrInvalidRequestBody)
		return
	}

	if err := h.auctionService.CreateAuction(c.Request.Context(), &auction); err != nil {
		utils.SendErrorResponse(c, err)
		return
	}

	utils.SendSuccessResponse(c, http.StatusCreated, auction)
}

func (h *AuctionHandler) GetAuction(c *gin.Context) {
	id := c.Param("id")
	auction, err := h.auctionService.GetAuction(c.Request.Context(), id)
	if err != nil {
		utils.SendErrorResponse(c, err)
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, auction)
}

func (h *AuctionHandler) PlaceBid(c *gin.Context) {
	var bid models.Bid
	if err := c.ShouldBindJSON(&bid); err != nil {
		utils.SendErrorResponse(c, shared_errors.ErrInvalidRequestBody)
		return
	}

	bid.AuctionID = c.Param("id")

	if err := h.auctionService.PlaceBid(c.Request.Context(), &bid); err != nil {
		utils.SendErrorResponse(c, err)
		return
	}

	// Notify via Firebase asynchronously
	go func() {
		// Use a background context for the async operation
		if err := h.firebaseApp.SendBidNotification(context.Background(), &bid); err != nil {
			h.logger.Error("Failed to send bid notification", zap.Error(err))
		}
	}()

	utils.SendSuccessResponse(c, http.StatusCreated, bid)
}