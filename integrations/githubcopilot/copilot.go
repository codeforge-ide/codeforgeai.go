package githubcopilot

import (
	"github.com/codeforge-ide/codeforgeai.go/config"
	"github.com/codeforge-ide/codeforgeai.go/models"
)

type CopilotModel struct {
	// ...fields...
}

// IsAuthenticated checks if the Copilot token is present in config
func (c *CopilotModel) IsAuthenticated() bool {
	// Load config and check CopilotToken
	cfg, err := config.LoadConfig("")
	if err != nil {
		return false
	}
	return cfg.CopilotToken != ""
}

func (c *CopilotModel) SendRequest(prompt string, config interface{}) (string, error) {
	// TODO: Implement Copilot LSP protocol
	return "", nil
}

var _ models.Model = (*CopilotModel)(nil)
