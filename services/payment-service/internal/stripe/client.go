package stripe

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/blytz/payment-service/internal/models"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/paymentintent"
	"github.com/stripe/stripe-go/v76/refund"
	"github.com/stripe/stripe-go/v76/webhook"
)

type Client struct {
	secretKey string
}

func NewClient(secretKey string) *Client {
	stripe.Key = secretKey
	return &Client{
		secretKey: secretKey,
	}
}

// CreatePaymentIntent creates a new payment intent
func (c *Client) CreatePaymentIntent(ctx context.Context, payment *models.Payment) (*stripe.PaymentIntent, error) {
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(payment.Amount),
		Currency: stripe.String(payment.Currency),
		Metadata: map[string]string{
			"payment_id": payment.PaymentID,
			"order_id":   payment.OrderID,
			"user_id":    payment.UserID,
			"auction_id": payment.AuctionID,
		},
		Description: stripe.String(payment.Description),
	}

	intent, err := paymentintent.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment intent: %w", err)
	}

	return intent, nil
}

// GetPaymentIntent retrieves a payment intent
func (c *Client) GetPaymentIntent(ctx context.Context, paymentIntentID string) (*stripe.PaymentIntent, error) {
	intent, err := paymentintent.Get(paymentIntentID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment intent: %w", err)
	}
	return intent, nil
}

// CreateRefund creates a refund
func (c *Client) CreateRefund(ctx context.Context, paymentIntentID string, amount int64, reason string) (*stripe.Refund, error) {
	params := &stripe.RefundParams{
		PaymentIntent: stripe.String(paymentIntentID),
		Amount:        stripe.Int64(amount),
		Reason:        stripe.String(reason),
	}

	refundObj, err := refund.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create refund: %w", err)
	}

	return refundObj, nil
}

// ConstructWebhookEvent constructs and validates a webhook event
func (c *Client) ConstructWebhookEvent(payload []byte, signature string, secret string) (*stripe.Event, error) {
	event, err := webhook.ConstructEvent(payload, signature, secret)
	if err != nil {
		return nil, fmt.Errorf("failed to construct webhook event: %w", err)
	}
	return &event, nil
}

// MapStripeStatus maps Stripe payment status to our internal status
func (c *Client) MapStripeStatus(stripeStatus stripe.PaymentIntentStatus) string {
	switch stripeStatus {
	case stripe.PaymentIntentStatusRequiresPaymentMethod:
		return models.PaymentStatusPending
	case stripe.PaymentIntentStatusRequiresConfirmation:
		return models.PaymentStatusPending
	case stripe.PaymentIntentStatusRequiresAction:
		return models.PaymentStatusProcessing
	case stripe.PaymentIntentStatusProcessing:
		return models.PaymentStatusProcessing
	case stripe.PaymentIntentStatusRequiresCapture:
		return models.PaymentStatusProcessing
	case stripe.PaymentIntentStatusSucceeded:
		return models.PaymentStatusSucceeded
	case stripe.PaymentIntentStatusCanceled:
		return models.PaymentStatusCanceled
	default:
		return models.PaymentStatusFailed
	}
}

// ExtractPaymentData extracts payment data from Stripe intent
func (c *Client) ExtractPaymentData(intent *stripe.PaymentIntent) map[string]interface{} {
	return map[string]interface{}{
		"id":                intent.ID,
		"status":            intent.Status,
		"amount":            intent.Amount,
		"currency":          intent.Currency,
		"payment_method":    intent.PaymentMethod,
		"client_secret":     intent.ClientSecret,
		"next_action":       intent.NextAction,
		"created":           intent.Created,
		"charges":           intent.Charges,
	}
}

// ValidateWebhookSignature validates webhook signature
func (c *Client) ValidateWebhookSignature(payload []byte, signature string, secret string) error {
	_, err := webhook.ConstructEvent(payload, signature, secret)
	return err
}

// ProcessWebhookEvent processes a webhook event
func (c *Client) ProcessWebhookEvent(event *stripe.Event) (*models.WebhookEvent, error) {
	var webhookEvent models.WebhookEvent
	webhookEvent.ID = event.ID
	webhookEvent.Type = string(event.Type)

	// Handle different event types
	switch event.Type {
	case stripe.EventTypePaymentIntentSucceeded:
		var intent stripe.PaymentIntent
		if err := json.Unmarshal(event.Data.Raw, &intent); err != nil {
			return nil, fmt.Errorf("failed to unmarshal payment intent: %w", err)
		}
		webhookEvent.PaymentID = intent.ID
		webhookEvent.Status = string(intent.Status)
		webhookEvent.Amount = intent.Amount
		webhookEvent.Currency = string(intent.Currency)
		webhookEvent.ProviderData = c.ExtractPaymentData(&intent)

	case stripe.EventTypePaymentIntentPaymentFailed:
		var intent stripe.PaymentIntent
		if err := json.Unmarshal(event.Data.Raw, &intent); err != nil {
			return nil, fmt.Errorf("failed to unmarshal payment intent: %w", err)
		}
		webhookEvent.PaymentID = intent.ID
		webhookEvent.Status = string(intent.Status)
		webhookEvent.Amount = intent.Amount
		webhookEvent.Currency = string(intent.Currency)
		webhookEvent.ProviderData = c.ExtractPaymentData(&intent)

	case stripe.EventTypeChargeRefunded:
		var charge stripe.Charge
		if err := json.Unmarshal(event.Data.Raw, &charge); err != nil {
			return nil, fmt.Errorf("failed to unmarshal charge: %w", err)
		}
		webhookEvent.PaymentID = charge.PaymentIntent.ID
		webhookEvent.Status = "refunded"
		webhookEvent.Amount = charge.AmountRefunded
		webhookEvent.Currency = string(charge.Currency)
	}

	return &webhookEvent, nil
}