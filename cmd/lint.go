package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/codeforge-ide/codeforgeai.go/linter"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(lintCmd)
}

var lintCmd = &cobra.Command{
	Use:   "lint [path]",
	Short: "Lint a file or directory for common issues",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		l := linter.NewLinter()

		info, err := os.Stat(path)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		if info.IsDir() {
			err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() {
					lintFile(l, filePath)
				}
				return nil
			})
			if err != nil {
				fmt.Println("Error:", err)
			}
		} else {
			lintFile(l, path)
		}
	},
}

func lintFile(l *linter.Linter, filePath string) {
	issues, err := l.LintFile(filePath)
	if err != nil {
		fmt.Printf("Error linting file %s: %v\n", filePath, err)
		return
	}

	if len(issues) > 0 {
		fmt.Printf("Found %d issues in %s:\n", len(issues), filePath)
		for _, issue := range issues {
			fmt.Printf("  - [%s] Line %d: %s\n", issue.Severity, issue.Line, issue.Message)
		}
	}
}
