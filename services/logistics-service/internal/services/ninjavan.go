package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type NinjaVanClient struct {
	httpClient *http.Client
	logger     *zap.Logger
	config     *NinjaVanConfig
}

type NinjaVanConfig struct {
	ClientID    string // Client ID from Ninja Dashboard
	ClientKey   string // Client Key from Ninja Dashboard (used for OAuth)
	Environment string // "sandbox" or "production"
	CountryCode string // e.g., "sg", "my", "id"
}

type NinjaVanTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type NinjaVanOrderRequest struct {
	TrackingNumber     string          `json:"requested_tracking_number"`
	ServiceType        string          `json:"service_type"`
	ServiceLevel       string          `json:"service_level"`
	OriginAddress      NinjaVanAddress `json:"origin"`
	DestinationAddress NinjaVanAddress `json:"destination"`
	Parcel             NinjaVanParcel  `json:"parcel"`
	IsCOD              bool            `json:"is_cod"`
	CODAmount          *float64        `json:"cod_amount,omitempty"`
	InsuranceAmount    *float64        `json:"insurance_amount,omitempty"`
	Remarks            string          `json:"remarks,omitempty"`
}

type NinjaVanAddress struct {
	Name       string `json:"name"`
	Phone      string `json:"phone"`
	Email      string `json:"email,omitempty"`
	Address1   string `json:"address1"`
	Address2   string `json:"address2,omitempty"`
	Area       string `json:"area,omitempty"`
	City       string `json:"city"`
	State      string `json:"state,omitempty"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}

type NinjaVanParcel struct {
	Weight      float64            `json:"weight"`
	Dimensions  NinjaVanDimensions `json:"dimensions"`
	Description string             `json:"description,omitempty"`
}

type NinjaVanDimensions struct {
	Length float64 `json:"length"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

type NinjaVanOrderResponse struct {
	TrackingID string `json:"tracking_id"`
	Status     string `json:"status"`
}

type NinjaVanTariffRequest struct {
	OriginAddress      NinjaVanAddress `json:"origin"`
	DestinationAddress NinjaVanAddress `json:"destination"`
	Parcel             NinjaVanParcel  `json:"parcel"`
	ServiceType        string          `json:"service_type"`
	ServiceLevel       string          `json:"service_level"`
}

type NinjaVanTariffResponse struct {
	Price        float64 `json:"price"`
	Currency     string  `json:"currency"`
	ServiceType  string  `json:"service_type"`
	ServiceLevel string  `json:"service_level"`
}

type NinjaVanPUDOPoint struct {
	ID       string           `json:"id"`
	Name     string           `json:"name"`
	Address  NinjaVanAddress  `json:"address"`
	Location NinjaVanLocation `json:"location"`
	IsActive bool             `json:"is_active"`
}

type NinjaVanLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func NewNinjaVanClient(logger *zap.Logger, config *NinjaVanConfig) *NinjaVanClient {
	return &NinjaVanClient{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		logger:     logger,
		config:     config,
	}
}

func (c *NinjaVanClient) getBaseURL() string {
	if c.config.Environment == "sandbox" {
		return "https://api-sandbox.ninjavan.co/sg"
	}
	return fmt.Sprintf("https://api.ninjavan.co/%s", c.config.CountryCode)
}

func (c *NinjaVanClient) getAccessToken(ctx context.Context) (string, error) {
	url := fmt.Sprintf("%s/2.0/oauth/access_token", c.getBaseURL())

	payload := map[string]string{
		"client_id":     c.config.ClientID,
		"client_secret": c.config.ClientKey, // Client Key is used as client_secret for OAuth
		"grant_type":    "client_credentials",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal token request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get access token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("token request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp NinjaVanTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("failed to decode token response: %w", err)
	}

	return tokenResp.AccessToken, nil
}

func (c *NinjaVanClient) CreateOrder(ctx context.Context, order *NinjaVanOrderRequest) (*NinjaVanOrderResponse, error) {
	token, err := c.getAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	url := fmt.Sprintf("%s/4.2/orders", c.getBaseURL())

	jsonData, err := json.Marshal(order)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal order request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create order request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("order creation failed with status %d: %s", resp.StatusCode, string(body))
	}

	var orderResp NinjaVanOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&orderResp); err != nil {
		return nil, fmt.Errorf("failed to decode order response: %w", err)
	}

	c.logger.Info("Ninja Van order created successfully",
		zap.String("tracking_id", orderResp.TrackingID))

	return &orderResp, nil
}

func (c *NinjaVanClient) CancelOrder(ctx context.Context, trackingID string) error {
	token, err := c.getAccessToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get access token: %w", err)
	}

	url := fmt.Sprintf("%s/2.2/orders/%s", c.getBaseURL(), trackingID)

	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create cancel request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to cancel order: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("order cancellation failed with status %d: %s", resp.StatusCode, string(body))
	}

	c.logger.Info("Ninja Van order cancelled successfully",
		zap.String("tracking_id", trackingID))

	return nil
}

func (c *NinjaVanClient) GetTariff(ctx context.Context, tariffReq *NinjaVanTariffRequest) (*NinjaVanTariffResponse, error) {
	token, err := c.getAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	url := fmt.Sprintf("%s/1.0/public/price", c.getBaseURL())

	jsonData, err := json.Marshal(tariffReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal tariff request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create tariff request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get tariff: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("tariff request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tariffResp NinjaVanTariffResponse
	if err := json.NewDecoder(resp.Body).Decode(&tariffResp); err != nil {
		return nil, fmt.Errorf("failed to decode tariff response: %w", err)
	}

	return &tariffResp, nil
}

func (c *NinjaVanClient) GetPUDOPoints(ctx context.Context) ([]NinjaVanPUDOPoint, error) {
	url := fmt.Sprintf("%s/2.0/pudos", c.getBaseURL())

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create PUDO request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get PUDO points: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("PUDO request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var pudoPoints []NinjaVanPUDOPoint
	if err := json.NewDecoder(resp.Body).Decode(&pudoPoints); err != nil {
		return nil, fmt.Errorf("failed to decode PUDO response: %w", err)
	}

	return pudoPoints, nil
}

func (c *NinjaVanClient) ConvertToNinjaVanAddress(addr AddressRequest) NinjaVanAddress {
	return NinjaVanAddress{
		Name:       addr.Name,
		Phone:      addr.PhoneNumber,
		Address1:   addr.Street,
		City:       addr.City,
		State:      addr.State,
		PostalCode: addr.PostalCode,
		Country:    addr.Country,
	}
}

func (c *NinjaVanClient) ConvertToNinjaVanParcel(weight float64, dims DimensionsRequest) NinjaVanParcel {
	return NinjaVanParcel{
		Weight: weight,
		Dimensions: NinjaVanDimensions{
			Length: dims.Length,
			Width:  dims.Width,
			Height: dims.Height,
		},
	}
}
