package models

import (
	"time"

	"github.com/codeforge-ide/codeforgeai.go/integrations/ollama"
)

// CodeModel wraps OllamaModel for code-specific prompts.
type CodeModel struct {
	ollama *ollama.OllamaModel
}

// NewCodeModel creates a CodeModel with the given model name.
// In future, this can be extended to select other backends.
func NewCodeModel(modelName string) *CodeModel {
	return &CodeModel{
		ollama: ollama.NewOllamaModel(modelName, "", 60*time.Second),
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
