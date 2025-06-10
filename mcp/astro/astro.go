package astro

import (
	"context"

	"github.com/codeforge-ide/codeforgeai.go/mcp"
)

const AstrolescentMCPURL = "https://mcp.astrolescent.com/sse"

type AstroMCP struct {
	client mcp.MCPInterface
}

func NewAstroMCP() *AstroMCP {
	return &AstroMCP{
		client: mcp.NewMCPClient(AstrolescentMCPURL),
	}
}

func (a *AstroMCP) GetPrice(ctx context.Context, token string) (*mcp.MCPResponse, error) {
	return a.client.CallTool(ctx, "price", map[string]interface{}{
		"token": token,
	})
}

func (a *AstroMCP) GetAPY(ctx context.Context) (*mcp.MCPResponse, error) {
	return a.client.CallTool(ctx, "apy", map[string]interface{}{})
}

func (a *AstroMCP) GetQuote(ctx context.Context, operation, fromToken, toToken string, amount float64, account string) (*mcp.MCPResponse, error) {
	args := map[string]interface{}{
		"operation":  operation,
		"from_token": fromToken,
		"to_token":   toToken,
		"amount":     amount,
	}
	if account != "" {
		args["account"] = account
	}

	return a.client.CallTool(ctx, "quote", args)
}

func (a *AstroMCP) GetBridge(ctx context.Context, fromChain, toChain, token string, amount float64) (*mcp.MCPResponse, error) {
	return a.client.CallTool(ctx, "bridge", map[string]interface{}{
		"from_chain": fromChain,
		"to_chain":   toChain,
		"token":      token,
		"amount":     amount,
	})
}

func (a *AstroMCP) GetAvailableTools() []string {
	return a.client.GetAvailableTools()
}
