package models

import (
	"github.com/codeforge-ide/codeforgeai.go/integrations/ollama"
	"github.com/codeforge-ide/codeforgeai.go/modeliface"
	// import other integrations as needed
)

// GeneralModel wraps OllamaModel for general-purpose prompts.
type GeneralModel struct {
	ollama *ollama.OllamaModel
}

// NewGeneralModel creates a GeneralModel with the given model name.
// In future, this can be extended to select other backends.
func NewGeneralModel(modelName string) modeliface.Model {
	switch modelName {
	case "ollama", "gemma3:1b", "qwen2.5-coder:1.5b":
		return ollama.NewOllamaModel(modelName, "", 0)
	// case "openai": return openai.NewOpenAIModel(modelName, ...)
	// case "copilot": return copilot.NewCopilotModel(modelName, ...)
	// Add more cases for other integrations as needed.
	default:
		return ollama.NewOllamaModel(modelName, "", 0)
	}
}

// SendRequest sends a prompt to the general model.
func (g *GeneralModel) SendRequest(prompt string, config interface{}) (string, error) {
	// If config is a map with "general_model", use that model name.
	if cfg, ok := config.(map[string]interface{}); ok {
		if m, ok := cfg["general_model"].(string); ok && m != "" && m != g.ollama.Model {
			g.ollama.Model = m
		}
	}
	return g.ollama.SendRequest(prompt, config)
}

var _ Model = (*GeneralModel)(nil)
