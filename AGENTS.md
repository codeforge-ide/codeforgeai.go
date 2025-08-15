# AGENTS: quick guide for automated contributors

Build / lint / test
- Build: `go build ./...`
- Lint: `gofmt -w .` and `golangci-lint run` (if installed)
- Run all tests: `go test ./...`
- Run a single test: `go test ./path/to/package -run TestName -v`

Code style (short)
- Formatting: run `gofmt` (project uses standard gofmt). Keep max line length ~120.
- Imports: use `goimports` or `gofmt` to group stdlib, third-party, and local imports.
- Types & naming: use Go exported/unexported naming rules (CamelCase for exported). Keep receiver names short (e.g., `r`, `s`).
- Error handling: return errors, wrap with `fmt.Errorf("...: %w", err)` when adding context; avoid panics in library code.
- Logging: prefer structured logging where available; keep messages clear and actionable.
- Tests: keep tests deterministic; use `t.Parallel()` only when safe; name tests `Test_<Function>_<Condition>`.

Repo-specific rules
- Use modules (go.mod present). Run `go mod tidy` after edits.
- Cursor / Copilot rules: include rules from .cursor/rules or .cursorrules and .github/copilot-instructions.md if present; follow any project-specific prompting constraints.

Notes
- Commit changes meaningfully; run linters and tests before committing.
- If adding env secrets, place placeholders in a `.env` and do not commit real secrets.
