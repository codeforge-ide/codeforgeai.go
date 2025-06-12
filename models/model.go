package models

import (
	"errors"
	"time"

	"github.com/codeforge-ide/codeforgeai.go/config"
	"github.com/codeforge-ide/codeforgeai.go/integrations/githubmodels"
	"github.com/codeforge-ide/codeforgeai.go/integrations/ollama"
	// ...add other integrations as needed...
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
		// Use Ollama model
		modelName := cfg.GeneralModel
		if modelType == "code" {
			modelName = cfg.CodeModel
		}
		return ollama.NewOllamaModel(modelName, "", 60*time.Second), nil
	case "githubmodels":
		token := "" // You may want to load from env or config
		modelName := cfg.GeneralModelGithub
		if modelType == "code" {
			modelName = cfg.CodeModelGithub
		}
		client := githubmodels.NewClient(token, modelName, "")
		return client, nil
	// Add more providers here as needed
	default:
		return nil, errors.New("unknown model provider: " + provider)
	}
}
