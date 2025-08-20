package linter

import (
	"bufio"
	"os"
	"strings"
)

// Issue represents a single issue found by the linter.
type Issue struct {
	File     string
	Line     int
	Message  string
	Severity string // e.g., "warning", "error"
}

// Linter is the main linter struct.
type Linter struct {
	MaxLineLength    int
	MaxFunctionLines int
}

// NewLinter creates a new linter with default settings.
func NewLinter() *Linter {
	return &Linter{
		MaxLineLength:    120,
		MaxFunctionLines: 50,
	}
}

// LintFile analyzes a single file and returns a list of issues.
func (l *Linter) LintFile(filePath string) ([]Issue, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var issues []Issue
	scanner := bufio.NewScanner(file)
	var lineNum int
	var funcLineCount int
	var inFunction bool

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Check for long lines
		if len(line) > l.MaxLineLength {
			issues = append(issues, Issue{
				File:     filePath,
				Line:     lineNum,
				Message:  "Line is too long",
				Severity: "warning",
			})
		}

		// Check for TODO comments
		if strings.Contains(line, "TODO") {
			issues = append(issues, Issue{
				File:     filePath,
				Line:     lineNum,
				Message:  "Found TODO comment",
				Severity: "info",
			})
		}

		// Check for long functions (very basic implementation)
		if strings.HasPrefix(line, "func ") {
			inFunction = true
			funcLineCount = 0
		}
		if inFunction {
			funcLineCount++
		}
		if inFunction && strings.HasPrefix(line, "}") {
			if funcLineCount > l.MaxFunctionLines {
				issues = append(issues, Issue{
					File:     filePath,
					Line:     lineNum,
					Message:  "Function is too long",
					Severity: "warning",
				})
			}
			inFunction = false
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return issues, nil
}
