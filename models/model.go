package models

import (
	"errors"
	"time"

	"github.com/codeforge-ide/codeforgeai.go/config"
	"github.com/codeforge-ide/codeforgeai.go/integrations/githubcopilot" // Copilot integration
	"github.com/codeforge-ide/codeforgeai.go/integrations/githubmodels"
	"github.com/codeforge-ide/codeforgeai.go/integrations/ollama"
)

// Model interface for all models
type Model interface {
	SendRequest(prompt string, config interface{}) (string, error)
}

// GetModelFromConfig returns a Model implementation based on config.Integrations.Default
func GetModelFromConfig(cfg *config.Config, modelType string) (Model, error) {
	provider := cfg.Integrations.Default
	switch provider {
	case "ollama":
		modelName := cfg.GeneralModel
		if modelType == "code" {
			modelName = cfg.CodeModel
		}
		return ollama.NewOllamaModel(modelName, "", 60*time.Second, cfg), nil
	case "githubmodels":
		token := ""
		modelName := cfg.GeneralModelGithub
		if modelType == "code" {
			modelName = cfg.CodeModelGithub
		}
		client := githubmodels.NewClient(token, modelName, "")
		return client, nil
	case "githubcopilot":
		return &githubcopilot.CopilotModel{}, nil
	default:
		return nil, errors.New("unknown model provider: " + provider)
	}
}
