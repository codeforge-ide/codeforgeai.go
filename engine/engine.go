package engine

import (
	"fmt"
	"os"

	"github.com/codeforge-ide/codeforgeai.go/config"
	"github.com/codeforge-ide/codeforgeai.go/models"
)

// Engine is the central orchestrator.
type Engine struct {
	Config       *config.Config
	GeneralModel models.Model
	CodeModel    models.Model
}

// NewEngine loads config and instantiates models.
func NewEngine(cfg *config.Config) *Engine {
	// Model selection is pluggable and based on config.
	var generalModel models.Model
	var codeModel models.Model

	// General model selection
	switch cfg.GeneralModel {
	case "ollama", "gemma3:1b", "qwen2.5-coder:1.5b":
		generalModel = models.NewGeneralModel(cfg.GeneralModel)
	// Add more cases for other integrations (e.g., openai, copilot) as needed.
	default:
		generalModel = models.NewGeneralModel(cfg.GeneralModel)
	}

	// Code model selection
	switch cfg.CodeModel {
	case "ollama", "qwen2.5-coder:1.5b":
		codeModel = models.NewCodeModel(cfg.CodeModel)
	// Add more cases for other integrations as needed.
	default:
		codeModel = models.NewCodeModel(cfg.CodeModel)
	}

	return &Engine{
		Config:       cfg,
		GeneralModel: generalModel,
		CodeModel:    codeModel,
	}
}

// RunAnalysis analyzes the current directory.
func (e *Engine) RunAnalysis() {
	// TODO: Implement directory analysis and classification.
	fmt.Println("Engine: Run analysis (not implemented yet).")
}

// ProcessPrompt finetunes and processes a user prompt.
func (e *Engine) ProcessPrompt(prompt string) string {
	resp, err := e.GeneralModel.SendRequest(prompt, map[string]interface{}{
		"general_model": e.Config.GeneralModel,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error processing prompt:", err)
		return ""
	}
	return resp
}

// ExplainCode explains code in a file.
func (e *Engine) ExplainCode(filePath string) string {
	// TODO: Read file content, send to code model with explain prompt.
	fmt.Println("Engine: Explain code (not implemented yet).")
	return ""
}

// ProcessCommitMessage generates a commit message.
func (e *Engine) ProcessCommitMessage(diff string) string {
	// Use code model to generate commit message.
	resp, err := e.CodeModel.SendRequest(e.Config.CommitMessagePrompt+"\n"+diff, map[string]interface{}{
		"code_model": e.Config.CodeModel,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error generating commit message:", err)
		return ""
	}
	return resp
}

// TODO: Implement more methods for edit, suggestion, etc.
