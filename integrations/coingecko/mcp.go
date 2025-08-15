package coingecko

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
)

// StartCoinGeckoMCP builds an *exec.Cmd to launch a CoinGecko MCP client.
// mode: "remote-keyless", "remote-byok", or "local".
// apiKey is required for local mode. usePro controls whether to set PRO vs DEMO env var for local.
// The returned *exec.Cmd is not started; the caller may call Start/Run and control lifecycle.
func StartCoinGeckoMCP(ctx context.Context, mode string, apiKey string, usePro bool) (*exec.Cmd, error) {
	switch mode {
	case "remote-keyless":
		return buildRemoteCmd("https://mcp.api.coingecko.com/sse")
	case "remote-byok":
		return buildRemoteCmd("https://mcp.pro-api.coingecko.com/sse")
	case "local":
		if apiKey == "" {
			return nil, errors.New("local mode requires apiKey (COINGECKO_PRO_API_KEY or COINGECKO_DEMO_API_KEY)")
		}
		// Ensure npx is present
		if _, err := exec.LookPath("npx"); err != nil {
			return nil, fmt.Errorf("npx not found in PATH: %w", err)
		}
		args := []string{"-y", "@coingecko/coingecko-mcp"}
		cmd := exec.CommandContext(ctx, "npx", args...)
		// Inherit environment and inject COINGECKO_* vars
		env := os.Environ()
		if usePro {
			env = append(env, "COINGECKO_PRO_API_KEY="+apiKey, "COINGECKO_ENVIRONMENT=pro")
		} else {
			env = append(env, "COINGECKO_DEMO_API_KEY="+apiKey, "COINGECKO_ENVIRONMENT=demo")
		}
		cmd.Env = env
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd, nil
	default:
		return nil, fmt.Errorf("unknown mode: %s", mode)
	}
}

func buildRemoteCmd(url string) (*exec.Cmd, error) {
	if _, err := exec.LookPath("npx"); err != nil {
		return nil, fmt.Errorf("npx not found in PATH: %w", err)
	}
	args := []string{"mcp-remote", url}
	cmd := exec.Command("npx", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd, nil
}
