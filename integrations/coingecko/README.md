CoinGecko MCP integration

This package provides helpers to construct commands that launch CoinGecko's MCP client/server wrappers.

Modes

- remote-keyless: uses the hosted public server at https://mcp.api.coingecko.com/sse (no API key)
- remote-byok: uses the hosted BYOK server at https://mcp.pro-api.coingecko.com/sse (requires auth in client)
- local: runs the npm package via npx (-y @coingecko/coingecko-mcp). Requires `npx` and Node installed; set API key and environment.

Examples

Programmatic usage:

    ctx := context.Background()
    cmd, err := coingecko.StartCoinGeckoMCP(ctx, "remote-keyless", "", false)
    if err != nil { ... }
    // start the process when ready:
    if err := cmd.Start(); err != nil { ... }

Local mode example (programmatic):

    cmd, err := coingecko.StartCoinGeckoMCP(ctx, "local", "MY_PRO_KEY", true)
    // call cmd.Start() to run

CLI

A small CLI helper `coingecko-mcp` is available as part of the project's CLI. Run with `--help` for options.

Docs

- Official MCP docs: https://docs.coingecko.com/reference/mcp-server
