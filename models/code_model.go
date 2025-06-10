package models

import (
	"github.com/codeforge-ide/codeforgeai.go/integrations/ollama"
	"github.com/codeforge-ide/codeforgeai.go/modeliface"
	// import other integrations as needed
)

// CodeModel wraps OllamaModel for code-specific prompts.
type CodeModel struct {
	ollama *ollama.OllamaModel
}

// NewCodeModel creates a CodeModel with the given model name.
// In future, this can be extended to select other backends.
func NewCodeModel(modelName string) modeliface.Model {
	switch modelName {
	case "ollama", "qwen2.5-coder:1.5b":
		return ollama.NewOllamaModel(modelName, "", 0)
	// case "openai": return openai.NewOpenAIModel(modelName, ...)
	// case "copilot": return copilot.NewCopilotModel(modelName, ...)
	// Add more cases for other integrations as needed.
	default:
		return ollama.NewOllamaModel(modelName, "", 0)
	}
}

// SendRequest sends a prompt to the code model.
func (c *CodeModel) SendRequest(prompt string, config interface{}) (string, error) {
	// If config is a map with "code_model", use that model name.
	if cfg, ok := config.(map[string]interface{}); ok {
		if m, ok := cfg["code_model"].(string); ok && m != "" && m != c.ollama.Model {
			c.ollama.Model = m
		}
	}
	return c.ollama.SendRequest(prompt, config)
}

var _ Model = (*CodeModel)(nil)
