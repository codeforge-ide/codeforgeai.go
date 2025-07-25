package githubcopilot

import (
	"context"

	"github.com/codeforge-ide/codeforgeai.go/modeliface"
)

// GithubCopilot represents the Github Copilot integration.
type GithubCopilot struct {
	// For now, this is a placeholder.
}

// CopilotModel implements modeliface.Model for Github Copilot
// It holds model metadata and implements SendRequest.
type CopilotModel struct {
	ID      string
	OwnedBy string
}

func (m *CopilotModel) SendRequest(prompt string, config interface{}) (string, error) {
	// Placeholder implementation
	return "CopilotModel response", nil
}

// New returns a new Github Copilot integration.
func New(apiKey string) *GithubCopilot {
	return &GithubCopilot{}
}

// Models returns a list of available models.
func (g *GithubCopilot) Models(ctx context.Context) ([]modeliface.Model, error) {
	// For now, this is a placeholder.
	return []modeliface.Model{
		&CopilotModel{
			ID:      "copilot",
			OwnedBy: "github",
		},
	}, nil
}
