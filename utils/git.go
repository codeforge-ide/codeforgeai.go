package utils

import (
	"fmt"
	"os/exec"
)

// GetGitDiff gets the current git diff.
func GetGitDiff() (string, error) {
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
