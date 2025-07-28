
# 🚀 CodeForgeAI.go

> **The Ultimate AI/MCP Powered Code Analysis & Automation Toolkit for Modern Developers**

A blazingly fast, modular, and extensible AI-powered code analysis engine built in Go. Whether you're building the next DeFi protocol, contributing to open source, or shipping production applications, CodeForgeAI.go supercharges your development workflow with intelligent automation and real-time insights.

## 📦 Usage

Install the CLI tool with:

```bash
go install github.com/codeforge-ide/codeforgeai.go@latest
```

For full command line usage and examples, see [`USAGE.md`](./USAGE.md).

## 🔌 Integrations

All integrations are now plug-and-play: each integration self-registers with a central registry at startup, and the CLI automatically discovers and wires up all available integrations. To add a new integration, simply implement the registration interface and provide a CLI command factory—no manual CLI wiring required.

### 🌐 io.net Integration
Connect to the IO Intelligence API for advanced agent workflows and model completions, supporting multi-agent orchestration, embeddings, and OpenAI-compatible endpoints. Use this for scalable, production-grade AI tasks with flexible model selection and quota management.

### 🌌 Astrolescent DeFi Integration
Get real-time DeFi data, including token prices, swap quotes, and APY analytics, directly in your AI workflows. Ideal for blockchain and DeFi developers needing live market intelligence and trading automation.

### 🐙 GitHub Copilot Integration
Enhance Copilot with full project context, smart code analysis, and auto-generated commit messages. Boost code review and test generation with AI-powered suggestions tailored to your codebase.

### 🦙 Ollama Integration
Run local LLMs with Ollama for fast, private code analysis and automation. Easily switch models and providers for flexible, developer-first AI workflows.

### 🧠 OpenAI Integration
Leverage OpenAI models for high-quality completions, code review, and documentation. Integrate seamlessly with other providers and MCP servers for context-rich automation.

## ✨ Why CodeForgeAI.go?

- 🧠 **Multi-LLM Intelligence**: Seamlessly integrate with Ollama, OpenAI, GitHub Copilot, and more
- ⚡ **Lightning Fast**: Go-native performance for instant analysis and feedback
- 🔌 **Modular Architecture**: Plug-and-play integrations that scale with your needs. Integrations self-register with a central registry, and the CLI auto-discovers and wires up all available integrations at startup for seamless extensibility.
- 🌐 **MCP-Powered**: Real-time connection to external data sources and blockchain networks
- 🛠️ **Developer-First**: Built by developers, for developers who demand excellence

## 🏗️ Architecture

```
codeforgeai.go/
├── 🎯 cmd/           # CLI entrypoints & commands
├── ⚙️  engine/       # Core AI analysis engine
├── 🔧 config/        # Configuration management
├── 🤖 models/        # LLM interfaces & adapters
├── 🔌 integrations/  # Pluggable AI services
│   ├── 🐙 githubcopilot/
│   ├── 🦙 ollama/
│   ├── 🧠 openai/
│   └── 📊 githubmodels/
├── 🌐 mcp/          # Model Context Protocol servers
│   ├── 🌌 astro/    # Astrolescent DeFi integration
│   └── 🐙 github/   # Enhanced GitHub integration
└── 🛠️  utils/       # Developer utilities
```

## 🚀 Quick Start

```bash
# Install
go install github.com/codeforge-ide/codeforgeai.go@latest

# Configure your favorite LLM
codeforgeai config set --provider ollama --model codellama

# Enable MCP integrations
codeforgeai mcp enable astrolescent
codeforgeai mcp enable github

# Start building the future
codeforgeai analyze ./my-defi-project
```

## 💡 Use Cases

### 🔥 For DeFi Developers
```bash
# Get live market analysis for your protocol
codeforgeai prompt "Analyze current ASTRL staking yields vs our protocol's APY"

# Smart contract optimization
codeforgeai analyze --focus security ./contracts/

# Generate DeFi-aware documentation
codeforgeai docs --include-market-data ./protocol/
```

### ⚡ For Any Developer
```bash
# Intelligent code review
codeforgeai review --pr-ready ./src/

# Generate context-aware tests
codeforgeai test generate --coverage-target 90 ./api/

# Smart commit messages
codeforgeai commit --stage-changes
```

## 🌟 What Makes It Special

- **🔗 Real-World Connected**: MCP integrations bring live data to your AI
- **⚡ Go Performance**: Native Go speed for enterprise-scale projects
- **🔧 Truly Modular**: Swap providers, add integrations, customize everything
- **🌍 Blockchain Native**: Built with Web3 and DeFi workflows in mind
- **🤝 Community Driven**: Open source, extensible, and growing

## 🤝 Contributing

Join the revolution! Whether you're building new MCP servers, adding LLM integrations, or improving the core engine, we welcome all contributors.

**Special Recognition**: Originally inspired by the Python `codeforgeai` project, now evolved into a next-generation Go-native powerhouse.

---

**For detailed usage of all commands, see [`USAGE.md`](./USAGE.md).**