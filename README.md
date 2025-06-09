# codeforgeai.go

A modular, pluggable AI-powered code analysis and automation tool, rewritten in Go.

- Functionality matches the original Python `codeforgeai` project.
- Modular architecture: core engine, config, models, integrations, CLI, and utilities.
- Pluggable LLM integrations: Ollama, OpenAI, GitHub Copilot, and more.
- CLI interface for analysis, prompt processing, commit message generation, and more.

## Structure

- `cmd/` - CLI entrypoints
- `engine/` - Core engine logic
- `config/` - Configuration management
- `models/` - Model interface and adapters
- `integrations/` - Pluggable integrations (Ollama, OpenAI, GitHub, etc)
- `utils/` - Utility functions

## Migration

This Go project is a line-for-line functional rewrite of the original Python `codeforgeai`, with improved modularity and extensibility.
