package directory

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// Node represents a file or directory in the tree.
type Node struct {
	Type           string  `json:"type"` // "file" or "directory"
	Name           string  `json:"name"`
	Path           string  `json:"path"`
	Children       []*Node `json:"children,omitempty"`
	Classification string  `json:"classification,omitempty"`
	Size           int64   `json:"size,omitempty"`
}

// ParseGitignore reads .gitignore from root and returns patterns.
func ParseGitignore(root string) ([]string, error) {
	var patterns []string
	f, err := os.Open(filepath.Join(root, ".gitignore"))
	if err != nil {
		return patterns, nil // No .gitignore is not an error
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		patterns = append(patterns, line)
	}
	return patterns, nil
}

// ShouldIgnore checks if a path matches any ignore pattern (simple glob).
func ShouldIgnore(relPath string, patterns []string) bool {
	for _, pat := range patterns {
		// Simple: ignore trailing slashes, match prefix or suffix
		pat = strings.TrimSuffix(pat, "/")
		if pat == "" {
			continue
		}
		if strings.HasPrefix(relPath, pat) || strings.HasSuffix(relPath, pat) {
			return true
		}
		// TODO: Use filepath.Match for more robust glob support
	}
	return false
}

// BuildTree walks the directory and builds a tree, applying .gitignore.
func BuildTree(root string) (*Node, error) {
	patterns, _ := ParseGitignore(root)
	return buildTreeRec(root, root, patterns)
}

func buildTreeRec(root, curr string, patterns []string) (*Node, error) {
	rel, _ := filepath.Rel(root, curr)
	if rel == "." {
		rel = ""
	}
	info, err := os.Lstat(curr)
	if err != nil {
		return nil, err
	}
	if rel != "" && ShouldIgnore(rel, patterns) {
		return nil, nil
	}
	node := &Node{
		Name: filepath.Base(curr),
		Path: rel,
	}
	if info.IsDir() {
		node.Type = "directory"
		entries, err := os.ReadDir(curr)
		if err != nil {
			return nil, err
		}
		for _, entry := range entries {
			childPath := filepath.Join(curr, entry.Name())
			child, err := buildTreeRec(root, childPath, patterns)
			if err != nil {
				continue
			}
			if child != nil {
				node.Children = append(node.Children, child)
			}
		}
	} else {
		node.Type = "file"
		node.Size = info.Size()
	}
	return node, nil
}

// AnalyzeDirectory builds and prints the directory tree as JSON.
func AnalyzeDirectory() {
	root, _ := os.Getwd()
	tree, err := BuildTree(root)
	if err != nil {
		println("Error:", err.Error())
		return
	}
	b, _ := json.MarshalIndent(tree, "", "  ")
	println(string(b))
}

// StripDirectory prints the tree after removing gitignored files.
func StripDirectory() {
	AnalyzeDirectory()
}

// ReadFileContent reads a file and returns its content as string.
func ReadFileContent(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// ClassifyFile is a stub for file classification.
func ClassifyFile(path, content string) string {
	// TODO: Implement real classification logic
	return "useful"
}
