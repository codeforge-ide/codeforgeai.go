package coingecko

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/codeforge-ide/codeforgeai.go/modeliface"
	"github.com/spf13/cobra"
)

func init() {
	modeliface.GlobalAgentRegistry.RegisterIntegration(modeliface.IntegrationRegistration{
		Metadata: modeliface.IntegrationMetadata{
			Name:         "coingecko",
			Description:  "CoinGecko MCP server helpers.",
			Commands:     []string{"prepare", "start"},
			Capabilities: []string{"crypto_price"},
		},
		CommandFactory: NewCoinGeckoCommand,
	})
}

func NewCoinGeckoCommand() *cobra.Command {
	coingeckoCmd := &cobra.Command{
		Use:   "coingecko",
		Short: "CoinGecko MCP server helpers",
	}

	coingeckoPrepare := &cobra.Command{
		Use:   "prepare",
		Short: "Print prepared CoinGecko MCP command (no side effects)",
		Run: func(cmd *cobra.Command, args []string) {
			mode, _ := cmd.Flags().GetString("mode")
			apiKey, _ := cmd.Flags().GetString("api-key")
			usePro, _ := cmd.Flags().GetBool("pro")
			ctx := context.Background()
			c, err := StartCoinGeckoMCP(ctx, mode, apiKey, usePro)
			if err != nil {
				fmt.Println("Error preparing CoinGecko MCP:", err)
				return
			}
			fmt.Println("prepared:", c.Args)
			// If local mode, also print injected env vars for clarity
			if strings.HasPrefix(c.Args[len(c.Args)-1], "@coingecko/coingecko-mcp") || (len(c.Args) >= 2 && c.Args[1] == "@coingecko/coingecko-mcp") {
				// Determine which env vars would be set (based on flags)
				if apiKey != "" {
					if usePro {
						fmt.Println("env:", "COINGECKO_PRO_API_KEY=<redacted>", "COINGECKO_ENVIRONMENT=pro")
					} else {
						fmt.Println("env:", "COINGECKO_DEMO_API_KEY=<redacted>", "COINGECKO_ENVIRONMENT=demo")
					}
				}
			}
		},
	}
	coingeckoPrepare.Flags().String("mode", "remote-keyless", "mode: remote-keyless|remote-byok|local")
	coingeckoPrepare.Flags().String("api-key", "", "api key for local or BYOK usage")
	coingeckoPrepare.Flags().Bool("pro", false, "use PRO environment for local mode")
	coingeckoCmd.AddCommand(coingeckoPrepare)

	coingeckoStart := &cobra.Command{
		Use:   "start",
		Short: "Start CoinGecko MCP process and supervise until SIGINT/SIGTERM",
		Run: func(cmd *cobra.Command, args []string) {
			mode, _ := cmd.Flags().GetString("mode")
			apiKey, _ := cmd.Flags().GetString("api-key")
			usePro, _ := cmd.Flags().GetBool("pro")
			ctx := context.Background()
			c, err := StartCoinGeckoMCP(ctx, mode, apiKey, usePro)
			if err != nil {
				fmt.Println("Error starting CoinGecko MCP:", err)
				return
			}
			fmt.Println("Starting CoinGecko MCP process... (press Ctrl+C to stop)")
			if err := c.Start(); err != nil {
				fmt.Println("Failed to start process:", err)
				return
			}
			// Wait for signal
			sigc := make(chan os.Signal, 1)
			signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
			<-sigc
			fmt.Println("Shutting down CoinGecko MCP process...")
			_ = c.Process.Kill()
		},
	}
	coingeckoStart.Flags().String("mode", "remote-keyless", "mode: remote-keyless|remote-byok|local")
	coingeckoStart.Flags().String("api-key", "", "api key for local or BYOK usage")
	coingeckoStart.Flags().Bool("pro", false, "use PRO environment for local mode")
	coingeckoCmd.AddCommand(coingeckoStart)

	return coingeckoCmd
}
