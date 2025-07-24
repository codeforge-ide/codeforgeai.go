package io

import (
	"context"

	"github.com/codeforgeai/codeforgeai.go/modeliface"
)

// IO represents the IO.net integration.
type IO struct {
	client *Client
}

// New returns a new IO.net integration.
func New(apiKey string) *IO {
	return &IO{
		client: NewClient(apiKey),
	}
}

// Models returns a list of available models.
func (i *IO) Models(ctx context.Context) ([]modeliface.Model, error) {
	return i.ListModels(ctx)
}
