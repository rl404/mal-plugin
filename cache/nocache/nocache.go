// Package nocache is a mock of caching.
package nocache

import "errors"

// Client is nocache client.
type Client struct{}

// New to create fake cache.
func New() (*Client, error) {
	return &Client{}, nil
}

// Set will just return nil.
func (c *Client) Set(key string, data interface{}) error {
	return nil
}

// Get will just return error to simulate as if data is not
// in cache.
func (c *Client) Get(key string, data interface{}) error {
	return errors.New("not using cache")
}

// Delete will just return nil.
func (c *Client) Delete(key string) error {
	return nil
}

// Close will just return nil.
func (c *Client) Close() error {
	return nil
}
