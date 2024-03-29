// Package bigcache is a wrapper of the original "github.com/allegro/bigcache" library.
//
// Only contains basic get, set, delete, and close methods.
// Data will be encoded to JSON before saving to cache.
package bigcache

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/allegro/bigcache"
	"github.com/rl404/mal-plugin/cache"
)

// Bigcache client implements Cacher interface.
var _ cache.Cacher = &Client{}

// Client is bigcache client.
type Client struct {
	bc *bigcache.BigCache
}

// New to create new cache with default bigcache config.
func New(expiredTime time.Duration) (*Client, error) {
	cfg := bigcache.DefaultConfig(expiredTime)
	cfg.CleanWindow = expiredTime
	return NewWithConfig(cfg)
}

// NewWithConfig to create new cache with bigcache config.
func NewWithConfig(cfg bigcache.Config) (*Client, error) {
	c, err := bigcache.NewBigCache(cfg)
	if err != nil {
		return nil, err
	}
	return NewFromBigCache(c), nil
}

// NewFromBigCache to create new cache from bigcache.
func NewFromBigCache(bc *bigcache.BigCache) *Client {
	return &Client{bc: bc}
}

// Get to get data from cache.
func (c *Client) Get(key string, data interface{}) error {
	d, err := c.bc.Get(key)
	if err != nil {
		return err
	}
	return json.Unmarshal(d, &data)
}

// Set to save data to cache.
func (c *Client) Set(key string, data interface{}) error {
	d, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return c.bc.Set(key, d)
}

// Delete to delete data from cache.
func (c *Client) Delete(key string) error {
	err :=  c.bc.Delete(key)
	if errors.Is(err, bigcache.ErrEntryNotFound) {
		return nil
	}
	return err
}

// Close to close cache connection.
func (c *Client) Close() error {
	return c.bc.Close()
}
