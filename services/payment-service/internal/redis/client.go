package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type Client struct {
	client *redis.Client
}

func NewClient(ctx context.Context, addr string) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping: %w", err)
	}

	return &Client{
		client: rdb,
	}, nil
}

func (c *Client) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

func (c *Client) Close() error {
	return c.client.Close()
}

func (c *Client) GetClient() *redis.Client {
	return c.client
}

// CachePayment caches payment data for quick lookup
func (c *Client) CachePayment(ctx context.Context, paymentID string, data map[string]interface{}) error {
	key := fmt.Sprintf("payment:%s", paymentID)
	return c.client.HMSet(ctx, key, data).Err()
}

// GetCachedPayment retrieves cached payment data
func (c *Client) GetCachedPayment(ctx context.Context, paymentID string) (map[string]string, error) {
	key := fmt.Sprintf("payment:%s", paymentID)
	result := c.client.HGetAll(ctx, key)
	return result.Val(), result.Err()
}

// SetPaymentStatus sets payment status in cache
func (c *Client) SetPaymentStatus(ctx context.Context, paymentID, status string) error {
	key := fmt.Sprintf("payment:%s", paymentID)
	return c.client.HSet(ctx, key, "status", status).Err()
}

// CacheUserPayments caches user's payment IDs
func (c *Client) CacheUserPayments(ctx context.Context, userID string, paymentIDs []string) error {
	key := fmt.Sprintf("user_payments:%s", userID)
	return c.client.SAdd(ctx, key, paymentIDs).Err()
}

// GetUserPayments retrieves user's payment IDs from cache
func (c *Client) GetUserPayments(ctx context.Context, userID string) ([]string, error) {
	key := fmt.Sprintf("user_payments:%s", userID)
	result := c.client.SMembers(ctx, key)
	return result.Val(), result.Err()
}

// PublishPaymentEvent publishes payment events for other services
func (c *Client) PublishPaymentEvent(ctx context.Context, eventType string, data map[string]interface{}) error {
	channel := fmt.Sprintf("payment_events:%s", eventType)
	return c.client.Publish(ctx, channel, data).Err()
}

// SubscribeToPaymentEvents subscribes to payment events
func (c *Client) SubscribeToPaymentEvents(ctx context.Context, eventType string) *redis.PubSub {
	channel := fmt.Sprintf("payment_events:%s", eventType)
	return c.client.Subscribe(ctx, channel)
}

// SetPaymentLock sets a lock for payment processing to prevent duplicates
func (c *Client) SetPaymentLock(ctx context.Context, paymentID string, ttl time.Duration) (bool, error) {
	key := fmt.Sprintf("payment_lock:%s", paymentID)
	return c.client.SetNX(ctx, key, "locked", ttl).Result()
}

// RemovePaymentLock removes payment processing lock
func (c *Client) RemovePaymentLock(ctx context.Context, paymentID string) error {
	key := fmt.Sprintf("payment_lock:%s", paymentID)
	return c.client.Del(ctx, key).Err()
}