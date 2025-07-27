package io

import (
	"context"
	"errors"
	"github.com/codeforge-ide/codeforgeai.go/modeliface"
)

// IO represents the IO.net integration.
type IO struct {
	client   *Client
	apiKey   string
	loggedIn bool
}

func (i *IO) Name() string {
	return "io"
}

func (i *IO) Login() error {
	if i.apiKey != "" {
		i.loggedIn = true
		return nil
	}
	return errors.New("missing IO.net API key")
}

func (i *IO) Logout() error {
	i.loggedIn = false
	return nil
}

func (i *IO) IsAuthenticated() bool {
	return i.loggedIn
}

// New returns a new IO.net integration.
func New(apiKey string) *IO {
	agent := &IO{
		client:   NewClient(apiKey),
		apiKey:   apiKey,
		loggedIn: false,
	}
	// Register with AgentManager if available
	// engine.GetAgentManager().Register(agent) // To be called from main/init
	return agent
}

// Models returns a list of available models.
func (i *IO) Models(ctx context.Context) ([]modeliface.Model, error) {
	return i.ListModels(ctx)
}
