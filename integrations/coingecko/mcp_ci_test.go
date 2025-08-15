package coingecko

import (
	"context"
	"os"
	"testing"
)

func TestStartCoinGeckoMCP_RemoteKeylessAndByok(t *testing.T) {
	// stub out lookPathFunc to avoid requiring npx in PATH
	orig := lookPathFunc
	lookPathFunc = func(file string) (string, error) { return "/usr/bin/npx", nil }
	defer func() { lookPathFunc = orig }()

	ctx := context.Background()
	cmd, err := StartCoinGeckoMCP(ctx, "remote-keyless", "", false)
	if err != nil {
		t.Fatalf("remote-keyless failed: %v", err)
	}
	if cmd.Path != "npx" && cmd.Args[0] != "npx" {
		t.Fatalf("expected npx command, got %v %v", cmd.Path, cmd.Args)
	}

	cmd, err = StartCoinGeckoMCP(ctx, "remote-byok", "", false)
	if err != nil {
		t.Fatalf("remote-byok failed: %v", err)
	}
}

func TestStartCoinGeckoMCP_Local_Mode_RequiresAPIKey(t *testing.T) {
	orig := lookPathFunc
	lookPathFunc = func(file string) (string, error) { return "/usr/bin/npx", nil }
	defer func() { lookPathFunc = orig }()

	ctx := context.Background()
	_, err := StartCoinGeckoMCP(ctx, "local", "", true)
	if err == nil {
		t.Fatalf("expected error when apiKey missing for local mode")
	}

	cmd, err := StartCoinGeckoMCP(ctx, "local", "ABC123", true)
	if err != nil {
		t.Fatalf("local mode with apiKey failed: %v", err)
	}
	// verify env contains COINGECKO_PRO_API_KEY when usePro true
	found := false
	for _, e := range cmd.Env {
		if e == "COINGECKO_PRO_API_KEY=ABC123" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected COINGECKO_PRO_API_KEY in env, got %v", cmd.Env)
	}

	// verify DEMO when usePro false
	cmd, err = StartCoinGeckoMCP(ctx, "local", "DEMOKEY", false)
	if err != nil {
		t.Fatalf("local mode with demo key failed: %v", err)
	}
	found = false
	for _, e := range cmd.Env {
		if e == "COINGECKO_DEMO_API_KEY=DEMOKEY" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected COINGECKO_DEMO_API_KEY in env, got %v", cmd.Env)
	}

	// ensure environment variable COINGECKO_ENVIRONMENT is set
	hasEnv := false
	for _, e := range cmd.Env {
		if e == "COINGECKO_ENVIRONMENT=demo" || e == "COINGECKO_ENVIRONMENT=pro" {
			hasEnv = true
			break
		}
	}
	if !hasEnv {
		t.Fatalf("expected COINGECKO_ENVIRONMENT set in env, got %v", cmd.Env)
	}

	// calling unknown mode yields error
	_, err = StartCoinGeckoMCP(ctx, "unknown-mode", "", false)
	if err == nil {
		t.Fatalf("expected error for unknown mode")
	}
}

func TestMain(m *testing.M) {
	// ensure tests run with a predictable environment
	os.Setenv("PATH", "/usr/bin:/bin")
	os.Exit(m.Run())
}
