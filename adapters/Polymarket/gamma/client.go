package polymarket

import (
	"net/http"
	"time"
)

type Client struct {
	BaseURL string
	HTTP    *http.Client
}

// NewClient returns a Gamma client with sane defaults.
// You can later add rate limiting, retry/backoff, headers, etc.
func NewClient() *Client {
	return &Client{
		BaseURL: "https://gamma-api.polymarket.com",
		HTTP: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// WithHTTP allows injecting a custom http.Client (useful for tests).
func (c *Client) WithHTTP(h *http.Client) *Client {
	if h != nil {
		c.HTTP = h
	}
	return c
}
