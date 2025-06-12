package engine

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/codeforge-ide/codeforgeai.go/config"
	"github.com/codeforge-ide/codeforgeai.go/directory"
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
	model, err := models.GetModelFromConfig(cfg, "general")
	if err != nil {
		panic(err)
	}
	return model
}

// getCodeModel instantiates the code model based on config.
func getCodeModel(cfg *config.Config) models.Model {
	model, err := models.GetModelFromConfig(cfg, "code")
	if err != nil {
		panic(err)
	}
	return model
}

// RunAnalysis analyzes the current directory and classifies files.
func (e *Engine) RunAnalysis() {
	cfg, err := loadFreshConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error loading config:", err)
		return
	}

	root, _ := os.Getwd()
	tree, err := directory.BuildTree(root)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error building directory tree:", err)
		return
	}

	// Use general model to classify the directory structure
	model := getGeneralModel(&cfg)
	treeJSON, _ := directory.SerializeTree(tree)

	prompt := cfg.DirectoryClassificationPrompt + "\n" + treeJSON
	resp, err := model.SendRequest(prompt, map[string]interface{}{
		"operation": "directory_classification",
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error classifying directory:", err)
		return
	}

	// Save classified result to .codeforge.json
	err = directory.SaveAnalysisResult(root, resp)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error saving analysis:", err)
		return
	}

	fmt.Println("Directory analysis complete. Results saved to .codeforge.json")
}

// ProcessPrompt finetunes and processes a user prompt.
func (e *Engine) ProcessPrompt(prompt string) string {
	cfg, err := loadFreshConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error loading config:", err)
		return ""
	}

	// Step 1: Finetune the prompt
	generalModel := getGeneralModel(&cfg)
	finetunePrompt := cfg.PromptFinetunePrompt + "\n" + prompt
	fineTunedPrompt, err := generalModel.SendRequest(finetunePrompt, map[string]interface{}{
		"operation": "prompt_finetune",
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error finetuning prompt:", err)
		return ""
	}

	// Step 2: Determine if response should be code or command
	codeOrCommandPrompt := cfg.CodeOrCommand + "\n" + fineTunedPrompt
	responseType, err := generalModel.SendRequest(codeOrCommandPrompt, map[string]interface{}{
		"operation": "classify_response_type",
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error classifying response type:", err)
		responseType = "code" // default to code
	}

	// Step 3: Process with appropriate model and prompt
	var finalPrompt, response string
	if strings.Contains(strings.ToLower(responseType), "command") {
		finalPrompt = cfg.CommandAgentPrompt + "\n" + fineTunedPrompt
		response, err = generalModel.SendRequest(finalPrompt, map[string]interface{}{
			"operation": "command_generation",
		})
	} else {
		codeModel := getCodeModel(&cfg)
		finalPrompt = cfg.CodePrompt + "\n" + fineTunedPrompt
		response, err = codeModel.SendRequest(finalPrompt, map[string]interface{}{
			"operation": "code_generation",
		})
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error processing prompt:", err)
		return ""
	}
	return response
}

// ExplainCode explains code in a file.
func (e *Engine) ExplainCode(filePath string) string {
	cfg, err := loadFreshConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error loading config:", err)
		return ""
	}

	// Read file content
	content, err := directory.ReadFileContent(filePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading file:", err)
		return ""
	}

	// Use code model to explain
	model := getCodeModel(&cfg)
	prompt := cfg.ExplainCodePrompt + "\n\nFile: " + filePath + "\n\n" + content

	resp, err := model.SendRequest(prompt, map[string]interface{}{
		"operation": "code_explanation",
		"file_path": filePath,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error explaining code:", err)
		return ""
	}
	return resp
}

// ProcessCommitMessage generates a commit message with gitmoji.
func (e *Engine) ProcessCommitMessage(diff string) string {
	cfg, err := loadFreshConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error loading config:", err)
		return ""
	}

	codeModel := getCodeModel(&cfg)

	// Step 1: Generate commit message
	commitPrompt := cfg.CommitMessagePrompt + "\n" + diff
	commitMsg, err := codeModel.SendRequest(commitPrompt, map[string]interface{}{
		"operation": "commit_message",
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error generating commit message:", err)
		return ""
	}

	// Step 2: Generate appropriate gitmoji
	generalModel := getGeneralModel(&cfg)
	gitmojiPrompt := cfg.GitmojiPrompt + "\n" + commitMsg
	gitmoji, err := generalModel.SendRequest(gitmojiPrompt, map[string]interface{}{
		"operation": "gitmoji_selection",
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error generating gitmoji:", err)
		gitmoji = "✨" // default gitmoji
	}

	return strings.TrimSpace(gitmoji) + " " + strings.TrimSpace(commitMsg)
}

// EditFiles edits files according to user prompt.
func (e *Engine) EditFiles(paths []string, userPrompt string, allowIgnore bool) error {
	cfg, err := loadFreshConfig()
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	root, _ := os.Getwd()
	patterns, _ := directory.ParseGitignore(root)

	for _, path := range paths {
		// Check if path should be ignored
		if !allowIgnore {
			relPath, _ := filepath.Rel(root, path)
			if directory.ShouldIgnore(relPath, patterns) {
				fmt.Printf("Skipping ignored path: %s\n", path)
				continue
			}
		}

		err := e.editSinglePath(path, userPrompt, &cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error editing %s: %v\n", path, err)
		}
	}
	return nil
}

func (e *Engine) editSinglePath(path, userPrompt string, cfg *config.Config) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return e.editDirectory(path, userPrompt, cfg)
	}
	return e.editFile(path, userPrompt, cfg)
}

func (e *Engine) editFile(filePath, userPrompt string, cfg *config.Config) error {
	content, err := directory.ReadFileContent(filePath)
	if err != nil {
		return err
	}

	model := getCodeModel(cfg)
	prompt := cfg.EditFinetunePrompt + "\n\nUser Request: " + userPrompt + "\n\nFile: " + filePath + "\n\n" + content

	editedContent, err := model.SendRequest(prompt, map[string]interface{}{
		"operation": "file_edit",
		"file_path": filePath,
	})
	if err != nil {
		return err
	}

	// Save to .codeforgedit file
	editFileName := filePath + ".codeforgedit"
	return os.WriteFile(editFileName, []byte(editedContent), 0644)
}

func (e *Engine) editDirectory(dirPath, userPrompt string, cfg *config.Config) error {
	// Build tree and get relevant files
	tree, err := directory.BuildTree(dirPath)
	if err != nil {
		return err
	}

	// Get all useful files from the tree
	files := directory.GetUsefulFiles(tree)

	for _, file := range files {
		fullPath := filepath.Join(dirPath, file)
		err := e.editFile(fullPath, userPrompt, cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error editing file %s: %v\n", fullPath, err)
		}
	}
	return nil
}

// ProvideSuggestion provides code suggestions.
func (e *Engine) ProvideSuggestion(filePath string, line int, snippet []string, entire bool) string {
	cfg, err := loadFreshConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error loading config:", err)
		return ""
	}

	var content string
	if snippet != nil && len(snippet) > 0 {
		content = strings.Join(snippet, "\n")
	} else if filePath != "" {
		fileContent, err := directory.ReadFileContent(filePath)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading file:", err)
			return ""
		}

		if entire {
			content = fileContent
		} else if line > 0 {
			lines := strings.Split(fileContent, "\n")
			if line <= len(lines) {
				// Get context around the line
				start := max(0, line-5)
				end := min(len(lines), line+5)
				content = strings.Join(lines[start:end], "\n")
			}
		}
	}

	if content == "" {
		return "No content provided for suggestion"
	}

	model := getCodeModel(&cfg)
	prompt := "Provide a code suggestion for the following:\n\n" + content

	resp, err := model.SendRequest(prompt, map[string]interface{}{
		"operation": "code_suggestion",
		"file_path": filePath,
		"line":      line,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error providing suggestion:", err)
		return ""
	}
	return resp
}

// GetGitDiff gets the current git diff.
func (e *Engine) GetGitDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--cached")
	output, err := cmd.Output()
	if err != nil {
		// Try unstaged diff if no staged changes
		cmd = exec.Command("git", "diff")
		output, err = cmd.Output()
		if err != nil {
			return "", fmt.Errorf("error getting git diff: %w", err)
		}
	}
	return string(output), nil
}

// Helper functions
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
