CoinGecko MCP integration

This package provides helpers to configure and launch/use the CoinGecko MCP server (remote or local) per https://docs.coingecko.com/reference/mcp-server

Usage:
- call coingecko.NewConfigFromEnv() to read env vars
- use MCPConfig.RemoteURL() for remote server URL
- use MCPConfig.LocalCommand() to get npm command + env to run local server

This package intentionally does not introduce runtime dependencies on external MCP registries; it provides utilities for CLI wiring in the host project.
