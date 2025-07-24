package io

import (
	"context"

	"github.com/codeforge-ide/codeforgeai.go/integrations/io"
	"github.com/codeforge-ide/codeforgeai.go/modeliface"
)

// IO represents the IO.net integration.
type IO struct {
	client *io.IO
}

// New returns a new IO.net integration.
func New(apiKey string) *IO {
	return &IO{
		client: io.New(apiKey),
	}
}

// Models returns a list of available models.
func (i *IO) Models(ctx context.Context) ([]modeliface.Model, error) {
	return i.client.Models(ctx)
}
