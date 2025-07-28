package openai

import (
	"github.com/codeforge-ide/codeforgeai.go/config"
	"github.com/codeforge-ide/codeforgeai.go/modeliface"
	"github.com/codeforge-ide/codeforgeai.go/models"
	"os"
)

type OpenAIModel struct {
	modeliface.BaseAgent
	// ...fields...
}

func (o *OpenAIModel) SendRequest(prompt string, config interface{}) (string, error) {
	// TODO: Implement OpenAI API call
	return "", nil
}

var _ models.Model = (*OpenAIModel)(nil)
var _ modeliface.Agent = (*OpenAIModel)(nil)

// IsAuthenticated checks if the OpenAI API key is set in config or env
func (o *OpenAIModel) IsAuthenticated() bool {
	cfg, err := config.LoadConfig("")
	if err == nil && cfg != (config.Config{}) {
		if cfgKey := cfg.OpenAIAPIKey; cfgKey != "" {
			return true
		}
	}
	return os.Getenv("OPENAI_API_KEY") != ""
}

// NewOpenAIModel creates a new OpenAIModel with the given name.
func NewOpenAIModel(name string) *OpenAIModel {
	return &OpenAIModel{
		BaseAgent: modeliface.BaseAgent{NameStr: name},
	}
}
