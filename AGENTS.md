# AGENTS.md

## Build, Lint, and Test Commands

- **Build all platforms:**  
  `./build.sh`
- **Build (local):**  
  `go build ./...`
- **Run all tests:**  
  `go test ./...`
- **Run a single test file:**  
  `go test ./path/to/file_test.go`
- **Run a single test function:**  
  `go test -run ^TestFuncName$ ./path/to/file_test.go`
- **Lint (recommended):**  
  `golangci-lint run`  
  (Install: `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`)
- **Format:**  
  `gofmt -s -w .`  
  or  
  `goimports -w .`

## Code Style Guidelines

- **Imports:**  
  - Standard, third-party, and local packages in separate blocks.
  - Use `goimports` to auto-sort.
- **Formatting:**  
  - Always run `gofmt` or `goimports` before committing.
- **Types & Naming:**  
  - Use `CamelCase` for exported names, `camelCase` for locals.
  - Interface names end with `-er` (e.g., `Reader`).
  - File names: `snake_case.go`, test files: `*_test.go`.
- **Error Handling:**  
  - Always check errors; return early on error.
  - Use `fmt.Errorf("context: %w", err)` for wrapping.
  - Prefer sentinel errors over string matching.
- **Tests:**  
  - Place in `*_test.go` files.
  - Use table-driven tests for variations.
  - Use `t.Helper()` for helpers.
- **Comments:**  
  - Exported functions/types must have doc comments.
  - Use `//nolint:<linter> // reason` to suppress linter warnings, with justification.
- **General:**  
  - Avoid global variables.
  - Keep functions short and focused.
  - Use context for cancellation/timeouts in long-running operations.
