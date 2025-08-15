package coingecko

import (
	"context"
	"os/exec"
	"strings"
	"testing"
)

func TestBuildRemoteCmd(t *testing.T) {
	// We can't rely on npx being present in CI; instead validate args when npx absent by temporarily
	// faking LookPath? Simpler: skip if npx not found.
	if _, err := exec.LookPath("npx"); err != nil {
		t.Skip("npx not found; skipping remote cmd build test")
	}
	cmd, err := buildRemoteCmd("https://mcp.api.coingecko.com/sse")
	if err != nil {
		t.Fatalf("buildRemoteCmd error: %v", err)
	}
	if cmd.Path == "" {
		t.Fatalf("expected cmd.Path to be set")
	}
	if len(cmd.Args) < 3 || !strings.Contains(cmd.Args[2], "coingecko.com") {
		t.Fatalf("unexpected args: %v", cmd.Args)
	}
}

func TestStartLocalCmdEnv(t *testing.T) {
	// validate construction only; skip if npx missing
	if _, err := exec.LookPath("npx"); err != nil {
		t.Skip("npx not found; skipping local cmd env test")
	}
	ctx := context.Background()
	cmd, err := StartCoinGeckoMCP(ctx, "local", "SOMEKEY", true)
	if err != nil {
		t.Fatalf("StartCoinGeckoMCP error: %v", err)
	}
	found := false
	for _, e := range cmd.Env {
		if strings.HasPrefix(e, "COINGECKO_PRO_API_KEY=") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected COINGECKO_PRO_API_KEY in env: %v", cmd.Env)
	}
}
