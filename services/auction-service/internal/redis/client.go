package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type Client struct {
	client *redis.Client
	scripts map[string]string
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
		scripts: make(map[string]string),
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

func (c *Client) LoadScripts(ctx context.Context) error {
	// Anti-snipe auction bid Lua script
	antiSnipeScript := `
local auction_key = KEYS[1]
local user_id = ARGV[1]
local bid_amount = tonumber(ARGV[2])
local current_time = tonumber(ARGV[3])
local extend_seconds = tonumber(ARGV[4])

local auction = redis.call('HMGET', auction_key, 'current_bid', 'end_time', 'winner_id', 'reserve_price')
local current_bid = tonumber(auction[1]) or 0
local end_time = tonumber(auction[2]) or 0
local winner_id = auction[3] or ''
local reserve_price = tonumber(auction[4]) or 0

if current_time > end_time then
    return {0, 'Auction ended'}
end

if bid_amount <= current_bid then
    return {0, 'Bid too low'}
end

redis.call('HMSET', auction_key,
    'current_bid', bid_amount,
    'winner_id', user_id,
    'last_bid_time', current_time
)

redis.call('PUBLISH', 'auction:' .. auction_key, user_id .. ':' .. bid_amount)

if (end_time - current_time) < 60 then
    redis.call('HINCRBY', auction_key, 'end_time', extend_seconds)
    redis.call('PUBLISH', 'auction:' .. auction_key, 'extended')
end

return {1, bid_amount}
`

	c.scripts["anti_snipe_bid"] = antiSnipeScript
	return nil
}

func (c *Client) ProcessBid(ctx context.Context, auctionID, userID string, amount int64) (bool, error) {
	auctionKey := fmt.Sprintf("auction:%s", auctionID)

	result, err := c.client.Eval(ctx, c.scripts["anti_snipe_bid"], []string{auctionKey},
		userID, amount, time.Now().Unix(), 60).Result()

	if err != nil {
		return false, fmt.Errorf("lua script error: %w", err)
	}

	values := result.([]interface{})
	success := values[0].(int64) == 1

	return success, nil
}

func (c *Client) GetAuction(ctx context.Context, auctionID string) (map[string]interface{}, error) {
	auctionKey := fmt.Sprintf("auction:%s", auctionID)
	result := c.client.HGetAll(ctx, auctionKey)
	return result.Val(), result.Err()
}

func (c *Client) CreateAuction(ctx context.Context, auctionID string, startTime, endTime int64, reservePrice int64) error {
	auctionKey := fmt.Sprintf("auction:%s", auctionID)

	return c.client.HMSet(ctx, auctionKey, map[string]interface{}{
		"start_time": startTime,
		"end_time": endTime,
		"reserve_price": reservePrice,
		"current_bid": 0,
		"winner_id": "",
		"status": "active",
	}).Err()
}

func (c *Client) EndAuction(ctx context.Context, auctionID string) error {
	auctionKey := fmt.Sprintf("auction:%s", auctionID)
	return c.client.HSet(ctx, auctionKey, "status", "ended").Err()
}

func (c *Client) SubscribeToAuction(ctx context.Context, auctionID string) *redis.PubSub {
	return c.client.Subscribe(ctx, fmt.Sprintf("auction:%s", auctionID))
}