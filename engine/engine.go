package engine

import (
	"github.com/codeforge-ide/codeforgeai.go/config"
	"github.com/codeforge-ide/codeforgeai.go/models"
)

type Engine struct {
	Config      *config.Config
	GeneralModel models.Model
	CodeModel    models.Model
}

func NewEngine(cfg *config.Config) *Engine {
	// TODO: Instantiate models based on config
	return &Engine{
		Config: cfg,
		// GeneralModel: ...,
		// CodeModel: ...,
	}
}

// TODO: Implement methods: RunAnalysis, ProcessPrompt, ExplainCode, GenerateCommitMessage, etc.
