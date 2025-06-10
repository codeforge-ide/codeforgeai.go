package models

import (
	"time"

	"github.com/codeforge-ide/codeforgeai.go/integrations/ollama"
)

// GeneralModel wraps OllamaModel for general-purpose prompts.
type GeneralModel struct {
	ollama *ollama.OllamaModel
}

// NewGeneralModel creates a GeneralModel with the given model name.
// In future, this can be extended to select other backends.
func NewGeneralModel(modelName string) *GeneralModel {
	return &GeneralModel{
		ollama: ollama.NewOllamaModel(modelName, "", 60*time.Second),
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
