package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/your/module/integrations/coingecko"
)

func main() {}

// This command is intended to be wired into the project's CLI tree. It's provided as a standalone file
// so maintainers can hook it into Cobra/Kingpin etc. Example usage from main package:
// go run ./cmd -subcommand coingecko-mcp --mode local --api-key XYZ --pro

func runCoinGeckoMCP() error {
	mode := flag.String("mode", "remote-keyless", "mcp mode: remote-keyless|remote-byok|local")
	apiKey := flag.String("api-key", "", "api key for local mode or BYOK")
	usePro := flag.Bool("pro", false, "set to use PRO environment (local mode)")
	startNow := flag.Bool("start", false, "start the command immediately")
	flag.Parse()
	ctx := context.Background()
	cmd, err := coingecko.StartCoinGeckoMCP(ctx, *mode, *apiKey, *usePro)
	if err != nil {
		return err
	}
	fmt.Printf("Prepared command: %v\n", cmd.Args)
	if *startNow {
		if err := cmd.Start(); err != nil {
			return err
		}
		// Wait for interrupt to stop
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		_ = cmd.Process.Kill()
	}
	return nil
}
