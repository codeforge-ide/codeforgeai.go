package io

import (
	"net/http"
	"time"
)

const (
	IONetAPI = "https://api.io.net"
)

// Client represents a client for interacting with the IO.net API.
type Client struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient creates a new IO.net API client.
func NewClient(apiKey string) *Client {
	return &Client{
		APIKey:  apiKey,
		BaseURL: IONetAPI,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}
