# AGENTS.md

This repository is Go-based. Follow these guidelines for all agentic coding tasks:

## Build, Lint, and Test
- Build all platforms: `./build.sh`
- Build (local): `go build ./...`
- Run all tests: `go test ./...`
- Run a single test file: `go test ./path/to/file_test.go`
- Run a single test function: `go test -run ^TestFunc$ ./path/to/file_test.go`
- Lint: `golangci-lint run` (install: `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`)
- Format: `gofmt -s -w .` or `goimports -w .`

## Code Style
- Imports: Standard, third-party, local (use `goimports` to sort)
- Naming: `CamelCase` for exports, `camelCase` for locals, interfaces end with `-er`, files use `snake_case.go`, tests: `*_test.go`
- Error handling: Always check errors, return early, wrap with `fmt.Errorf("context: %w", err)`, prefer sentinel errors
- Tests: Use `*_test.go`, table-driven tests, `t.Helper()` for helpers
- Comments: All exports need doc comments; use `//nolint:<linter> // reason` for linter suppression
- General: Avoid globals, keep functions focused, use context for cancellation/timeouts

notes:
check the flow.md for further instructions 

also be extra careful so has not to commit any build files or build outputs to github