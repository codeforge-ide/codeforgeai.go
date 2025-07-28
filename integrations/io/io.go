package io

import (
	"context"
	"github.com/codeforge-ide/codeforgeai.go/modeliface"
)

// IO represents the IO.net integration.
type IO struct {
	modeliface.BaseAgent
	client *Client
	apiKey string
}

// Name, Login, Logout, IsAuthenticated are now provided by BaseAgent.

// New returns a new IO.net integration.
func New(apiKey string) *IO {
	agent := &IO{
		client:    NewClient(apiKey),
		apiKey:    apiKey,
		BaseAgent: modeliface.BaseAgent{NameStr: "io"},
	}
	// Register with AgentManager if available
	// engine.GetAgentManager().Register(agent) // To be called from main/init
	return agent
}

// Models returns a list of available models.
func (i *IO) Models(ctx context.Context) ([]modeliface.Model, error) {
	return i.ListModels(ctx)
}
