package githubcopilot

import (
	"bytes"
	"fmt"
	"github.com/codeforge-ide/codeforgeai.go/config"
	"os/exec"
	"strings"
)

type CopilotModel struct {
	// ...fields...
}

// IsAuthenticated checks if the Copilot token is present in config
func (c *CopilotModel) IsAuthenticated() bool {
	cfg, err := config.LoadConfig("")
	if err != nil {
		return false
	}
	return cfg.CopilotToken != ""
}

func (c *CopilotModel) SendRequest(prompt string, config interface{}) (string, error) {
	// Check for gh CLI
	_, err := exec.LookPath("gh")
	if err != nil {
		return "", fmt.Errorf("GitHub CLI (gh) not found. Please install it from https://cli.github.com/")
	}
	// Check for copilot extension
	extCheck := exec.Command("gh", "extension", "list")
	out, err := extCheck.Output()
	if err != nil || !strings.Contains(string(out), "gh-copilot") {
		return "", fmt.Errorf("GitHub Copilot CLI extension not found. Install with: gh extension install github/gh-copilot")
	}
	// Run gh copilot suggest with the prompt
	cmd := exec.Command("gh", "copilot", "suggest", prompt)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("Copilot CLI error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

// Satisfy the Model interface
// (registration is handled in models/model.go)
