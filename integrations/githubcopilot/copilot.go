package githubcopilot

import "github.com/codeforge-ide/codeforgeai.go/models"

type CopilotModel struct {
	// ...fields...
}

func (c *CopilotModel) SendRequest(prompt string, config interface{}) (string, error) {
	// TODO: Implement Copilot LSP protocol
	return "", nil
}
var _ models.Model = (*CopilotModel)(nil)
