# CodeForgeAI.go Usage Guide

This document provides detailed usage instructions for all available commands and options in the CodeForgeAI.go CLI.

---

## Global Options

- `-v`, `--verbose`         Set loglevel to INFO
- `-V`, `--very-verbose`    Set loglevel to DEBUG
- `--debug`                 Enable debug mode (overrides other verbosity flags)

---

## Core Commands

### `analyze`

Analyze the current working directory or a specified path.

```bash
codeforgeai analyze [path] [--mcp astrolescent|github] [--query QUERY] [--focus AREA] [--loop]
```

- `[path]` (optional): Directory to analyze (default: current directory)
- `--mcp`: Enable MCP integration (`astrolescent`, `github`)
- `--query`: Specific query for analysis (e.g., "staking", "security")
- `--focus`: Focus area (e.g., "security", "performance")
- `--loop`: Enable adaptive feedback loop

---

### `prompt`

Process a user prompt and get an AI-powered response.

```bash
codeforgeai prompt "Your prompt here"
```

---

### `config`

Check and display the current configuration.

```bash
codeforgeai config
```

---

### `strip`

Print the directory tree after removing gitignored files.

```bash
codeforgeai strip
```

---

### `commit-message`

Generate a commit message with code changes and gitmoji.

```bash
codeforgeai commit-message
```

---

### `explain`

Explain the code in a given file.

```bash
codeforgeai explain [file_path]
```

---

### `extract`

Extract code blocks from a file or string.

```bash
codeforgeai extract --file [file_path]
codeforgeai extract --string "code block here"
```

---

### `format`

Format code blocks for readability.

```bash
codeforgeai format --file [file_path]
codeforgeai format --string "code block here"
```

---

### `command`

Process a command request.

```bash
codeforgeai command "Describe the command you want"
```

---

### `edit`

Edit code in specified files or folders.

```bash
codeforgeai edit [paths...] --user_prompt "Edit prompt here" [--allow-ignore]
```

- `[paths...]`: Files or directories to edit (default: current directory)
- `--user_prompt`: User prompt for editing (required)
- `--allow-ignore`: Allow explicitly passed directories to be processed even if .gitignore ignores them

---

### `suggestion`

Get code suggestions from the code model.

```bash
codeforgeai suggestion --file [file_path] --line [line_number] --string "code snippet" --entire
```

- `--file`: File to read code from
- `--line`: Line number to use for suggestion
- `--string`: User-provided code snippet for suggestion (can be repeated)
- `--entire`, `-E`: Send entire file content for suggestion

---

## Integration Commands

### `github`

GitHub Copilot integration.

```bash
codeforgeai github copilot login
codeforgeai github copilot logout
codeforgeai github copilot status
codeforgeai github copilot lsp
```

---

### `github-models`

Interact with GitHub Models API.

```bash
codeforgeai github-models prompt "Prompt here"
codeforgeai github-models multi-turn
codeforgeai github-models stream "Prompt here"
codeforgeai github-models image "Prompt here" [image_path]
```

---

### `secret-ai`

Secret AI SDK integration.

```bash
codeforgeai secret-ai list-models
codeforgeai secret-ai test-connection
codeforgeai secret-ai chat "Message here"
```

---

### `web3`

Web3 development tools.

```bash
codeforgeai web3 scaffold [project_name]
codeforgeai web3 analyze-contract [contract_file]
codeforgeai web3 estimate-gas [contract_file]
codeforgeai web3 generate-tests [contract_file]
codeforgeai web3 check-env
codeforgeai web3 install-deps
```

---

### `zerepy`

ZerePy agent integration.

```bash
codeforgeai zerepy status
codeforgeai zerepy list-agents
codeforgeai zerepy load-agent [agent_name]
codeforgeai zerepy action [connection] [action]
codeforgeai zerepy chat "Message here"
```

---

### `solana`

Solana blockchain and MCP integration.

```bash
codeforgeai solana status
codeforgeai solana balance
codeforgeai solana transfer [destination] [amount]
codeforgeai solana mcp interact [program_id] [action_type]
codeforgeai solana mcp state [program_id] [account_address]
codeforgeai solana mcp init-account [program_id] [space]
```

---

### `astro`

Astrolescent MCP Server integration for Radix DeFi data.

```bash
codeforgeai astro price
codeforgeai astro quote [operation] [token] [amount] [account]
codeforgeai astro apy
codeforgeai astro analyze [type]
codeforgeai astro calculator [amount] [days]
codeforgeai astro trading-advice [from_token] [to_token] [amount]
codeforgeai astro demo
```

- `operation`: `buy`, `sell`, or `swap`
- `type`: `staking-vs-lp`, `market`

---

### `mcp`

Manage MCP server integrations.

```bash
codeforgeai mcp enable [server]
codeforgeai mcp disable [server]
```

- `[server]`: `astrolescent`, `github`

---

## Additional Notes

- All commands support `-v`, `-V`, and `--debug` for logging and debugging.
- For integration commands, ensure required environment variables (e.g., `GITHUB_TOKEN`) are set.
- For more details on each command, use `codeforgeai [command] --help`.

