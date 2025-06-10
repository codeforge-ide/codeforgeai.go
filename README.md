# 🚀 CodeForgeAI.go

> **The Ultimate AI/MCP Powered Code Analysis & Automation Toolkit for Modern Developers**

A blazingly fast, modular, and extensible AI-powered code analysis engine built in Go. Whether you're building the next DeFi protocol, contributing to open source, or shipping production applications, CodeForgeAI.go supercharges your development workflow with intelligent automation and real-time insights.

## ✨ Why CodeForgeAI.go?

- 🧠 **Multi-LLM Intelligence**: Seamlessly integrate with Ollama, OpenAI, GitHub Copilot, and more
- ⚡ **Lightning Fast**: Go-native performance for instant analysis and feedback
- 🔌 **Modular Architecture**: Plug-and-play integrations that scale with your needs
- 🌐 **MCP-Powered**: Real-time connection to external data sources and blockchain networks
- 🛠️ **Developer-First**: Built by developers, for developers who demand excellence

## 🌟 Model Context Protocol (MCP) Integrations

CodeForgeAI.go leverages the cutting-edge **Model Context Protocol** to connect your AI workflows with real-world data and services. No more isolated AI—get live, actionable intelligence.

### 🌌 Astrolescent DeFi Integration

**Perfect for Blockchain & DeFi Developers**

Connect your AI directly to the Radix DeFi ecosystem through our Astrolescent MCP server integration:

- 💸 **Live Token Prices**: Get real-time $ASTRL, XRD, and token prices in your prompts
- 🔄 **Smart Swap Quotes**: AI-powered trading analysis with live DEX data
- 📈 **Yield Intelligence**: Real-time APY data for staking and liquidity provision
- 🌉 **Cross-Chain Insights**: Bridge data and multi-chain analytics
- 🤖 **AI Trading Assistant**: Build intelligent DeFi bots and analysis tools

```bash
# Example: AI-powered DeFi analysis
codeforgeai analyze --mcp astrolescent "What's the best yield strategy for 10k ASTRL today?"
```

### 🐙 GitHub Copilot Enhanced

Supercharge your GitHub Copilot experience with contextual project intelligence:

- 📊 **Project Context Awareness**: Feed your entire codebase context to Copilot
- 🔍 **Smart Code Analysis**: Enhanced suggestions based on project patterns
- 📝 **Intelligent Commit Messages**: Auto-generate meaningful commit descriptions
- 🧪 **Test Generation**: AI-powered test creation with project context

## 🎯 Perfect For

### 🏗️ Blockchain Developers
- **Smart Contract Analysis**: AI-powered security audits and optimization
- **DeFi Protocol Development**: Real-time market data integration
- **Cross-Chain Development**: Multi-network insights and analytics
- **Token Economics**: AI-assisted tokenomics modeling and analysis

### 🚀 Modern Development Teams
- **Code Quality Automation**: Intelligent analysis and suggestions
- **Documentation Generation**: Auto-generate docs that actually make sense
- **Legacy Code Migration**: AI-assisted modernization strategies
- **Performance Optimization**: Smart bottleneck detection and solutions

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
go install github.com/nathfavour/codeforgeai.go@latest

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

**Ready to forge the future of AI-powered development?** 🔥

Star ⭐ this repo and join thousands of developers building with intelligent automation!