// Package redis is a wrapper of the original "github.com/go-redis/redis/v8" library.
//
// Only contains basic publish, subscribe, and close methods.
// Data will be encoded to JSON before publishing the message.
package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rl404/mal-plugin/pubsub"
)

// Redis client implements PubSub interface.
var _ pubsub.PubSub = &Client{}

// Redis channels implements Channel interface.
var _ pubsub.Channel = &Channel{}

// Client is redis pubsub client.
type Client struct {
	client *redis.Client
	ctx    context.Context
}

// Channel is redis pubsub channel.
type Channel struct {
	channel *redis.PubSub
}

// New to create new redis pubsub client.
func New(address, password string) (*Client, error) {
	return NewWithConfig(redis.Options{
		Addr:     address,
		Password: password,
	})
}

// NewWithConfig to create pubsub from go-redis options.
func NewWithConfig(option redis.Options) (*Client, error) {
	client := redis.NewClient(&option)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Ping test.
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return NewFromGoRedis(client), nil
}

// NewFromGoRedis to create pubsub from go-redis client.
func NewFromGoRedis(client *redis.Client) *Client {
	return &Client{
		client: client,
		ctx:    context.Background(),
	}
}

// Publish to publish message.
func (c *Client) Publish(channel string, data interface{}) error {
	j, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return c.client.Publish(c.ctx, channel, j).Err()
}

// Subscribe to subscribe channel.
func (c *Client) Subscribe(channel string) (interface{}, error) {
	return &Channel{
		channel: c.client.Subscribe(c.ctx, channel),
	}, nil
}

// Close to close redis pubsub client.
func (c *Client) Close() error {
	return c.client.Close()
}

// Read to read incoming message.
func (c *Channel) Read(model interface{}) (<-chan interface{}, <-chan error) {
	msgChan, errChan := make(chan interface{}), make(chan error)
	go func() {
		for msg := range c.channel.Channel() {
			if err := json.Unmarshal([]byte(msg.Payload), &model); err != nil {
				errChan <- err
			} else {
				msgChan <- model
			}
		}
	}()
	return (<-chan interface{})(msgChan), (<-chan error)(errChan)
}

// Close to close subscription.
func (c *Channel) Close() error {
	return c.channel.Close()
}
