package engine

import (
	"fmt"
	"os"

	"github.com/codeforge-ide/codeforgeai.go/config"
	"github.com/codeforge-ide/codeforgeai.go/models"
)

// Engine orchestrates operations using fresh config and models.
type Engine struct{}

// NewEngine is kept for compatibility but does nothing now.
func NewEngine(cfg *config.Config) *Engine {
	return &Engine{}
}

// loadFreshConfig loads the latest config from disk.
func loadFreshConfig() (config.Config, error) {
	return config.EnsureConfigPrompts("")
}

// getGeneralModel instantiates the general model based on config.
func getGeneralModel(cfg *config.Config) models.Model {
	return models.NewGeneralModel(cfg.GeneralModel)
}

// getCodeModel instantiates the code model based on config.
func getCodeModel(cfg *config.Config) models.Model {
	return models.NewCodeModel(cfg.CodeModel)
}

// RunAnalysis analyzes the current directory.
func (e *Engine) RunAnalysis() {
	cfg, err := loadFreshConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error loading config:", err)
		return
	}
	// TODO: Call directory analysis utilities and classification here.
	fmt.Println("Engine: Run analysis (not implemented yet).")
	_ = cfg // placeholder to avoid unused warning
}

// ProcessPrompt finetunes and processes a user prompt.
func (e *Engine) ProcessPrompt(prompt string) string {
	cfg, err := loadFreshConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error loading config:", err)
		return ""
	}
	model := getGeneralModel(&cfg)
	resp, err := model.SendRequest(prompt, map[string]interface{}{
		"general_model": cfg.GeneralModel,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error processing prompt:", err)
		return ""
	}
	return resp
}

// ExplainCode explains code in a file.
func (e *Engine) ExplainCode(filePath string) string {
	cfg, err := loadFreshConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error loading config:", err)
		return ""
	}
	model := getCodeModel(&cfg)
	// TODO: Read file content, send to code model with explain prompt.
	fmt.Println("Engine: Explain code (not implemented yet).")
	_ = model
	return ""
}

// ProcessCommitMessage generates a commit message.
func (e *Engine) ProcessCommitMessage(diff string) string {
	cfg, err := loadFreshConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error loading config:", err)
		return ""
	}
	model := getCodeModel(&cfg)
	resp, err := model.SendRequest(cfg.CommitMessagePrompt+"\n"+diff, map[string]interface{}{
		"code_model": cfg.CodeModel,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error generating commit message:", err)
		return ""
	}
	return resp
}

// TODO: Implement more methods for edit, suggestion, etc.
