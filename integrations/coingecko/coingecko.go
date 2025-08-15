package coingecko

import (
	"fmt"
	"os"
)

// MCPConfig holds configuration for launching/coordinating with CoinGecko MCP servers.
type MCPConfig struct {
	Environment string // "pro" or "demo" or empty
	ProKey     string
	DemoKey    string
	RemoteAuth bool // use pro remote (BYOK) when true
}

// NewConfigFromEnv reads environment variables per CoinGecko docs.
func NewConfigFromEnv() *MCPConfig {
	return &MCPConfig{
		Environment: os.Getenv("COINGECKO_ENVIRONMENT"),
		ProKey:      os.Getenv("COINGECKO_PRO_API_KEY"),
		DemoKey:     os.Getenv("COINGECKO_DEMO_API_KEY"),
	}
}

// RemoteURL returns the appropriate remote MCP SSE endpoint.
func (c *MCPConfig) RemoteURL() string {
	if c.Environment == "pro" {
		return "https://mcp.pro-api.coingecko.com/sse"
	}
	return "https://mcp.api.coingecko.com/sse"
}

// LocalCommand returns the npm-based local server command and env map when applicable.
func (c *MCPConfig) LocalCommand() (cmd string, args []string, env map[string]string) {
	cmd = "npx"
	args = []string{"-y", "@coingecko/coingecko-mcp"}
	env = map[string]string{}
	if c.Environment == "pro" && c.ProKey != "" {
		env["COINGECKO_PRO_API_KEY"] = c.ProKey
		env["COINGECKO_ENVIRONMENT"] = "pro"
	} else if c.DemoKey != "" {
		env["COINGECKO_DEMO_API_KEY"] = c.DemoKey
		env["COINGECKO_ENVIRONMENT"] = "demo"
	}
	return
}

// Register integration (no-op placeholder for init-based registration used elsewhere)
func Register() error {
	// If the project has a central registry for MCP integrations, hook here.
	// We'll check at runtime whether registry is available.
	if os.Getenv("COINGECKO_ENVIRONMENT") == "" && os.Getenv("COINGECKO_PRO_API_KEY") == "" && os.Getenv("COINGECKO_DEMO_API_KEY") == "" {
		// no config present; still valid — remote keyless can be used
		fmt.Println("coingecko: no env keys set; remote keyless usage allowed")
	}
	return nil
}
