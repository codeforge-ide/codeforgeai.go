package mcp

import (
	"context"

	"github.com/codeforge-ide/codeforgeai.go/mcp/astro"
	"github.com/codeforge-ide/codeforgeai.go/mcp/githubcopilot"
	"github.com/codeforge-ide/codeforgeai.go/mcp/io"
	"github.com/codeforge-ide/codeforgeai.go/modeliface"
)

// Client is a client for the Multi-Model Communication Protocol (MCP).
type Client struct {
	Astro         *astro.Astro
	GithubCopilot *githubcopilot.GithubCopilot
	IO            *io.IO
}

// NewClient creates a new MCP client.
func NewClient(astroAPIKey, githubCopilotAPIKey, ioAPIKey string) *Client {
	return &Client{
		Astro:         astro.New(astroAPIKey),
		GithubCopilot: githubcopilot.New(githubCopilotAPIKey),
		IO:            io.New(ioAPIKey),
	}
}

// Models returns a list of available models from all integrations.
func (c *Client) Models(ctx context.Context) ([]modeliface.Model, error) {
	var models []modeliface.Model

	astroModels, err := c.Astro.Models(ctx)
	if err != nil {
		return nil, err
	}
	models = append(models, astroModels...)

	githubCopilotModels, err := c.GithubCopilot.Models(ctx)
	if err != nil {
		return nil, err
	}
	models = append(models, githubCopilotModels...)

	ioModels, err := c.IO.Models(ctx)
	if err != nil {
		return nil, err
	}
	models = append(models, ioModels...)

	return models, nil
}
