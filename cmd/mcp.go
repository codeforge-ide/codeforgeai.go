package cmd

import (
	"fmt"
	"log"

	"github.com/codeforge-ide/codeforgeai.go/config"
	"github.com/spf13/cobra"
)

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Manage MCP server integrations",
	Long:  "Enable, disable, and configure Model Context Protocol server integrations",
}

var mcpEnableCmd = &cobra.Command{
	Use:   "enable [server]",
	Short: "Enable an MCP server integration",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		server := args[0]
		cfg, _ := config.LoadConfig("")

		switch server {
		case "astrolescent":
			cfg.MCP.Astrolescent.Enabled = true
			fmt.Println("✅ Astrolescent MCP server enabled")
		case "github":
			cfg.MCP.GitHub.Enabled = true
			fmt.Println("✅ GitHub MCP server enabled")
		default:
			log.Fatalf("Unknown MCP server: %s", server)
		}

		if err := config.SaveConfig("", cfg); err != nil {
			log.Fatalf("Failed to save config: %v", err)
		}
	},
}

var mcpDisableCmd = &cobra.Command{
	Use:   "disable [server]",
	Short: "Disable an MCP server integration",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		server := args[0]
		cfg, _ := config.LoadConfig("")

		switch server {
		case "astrolescent":
			cfg.MCP.Astrolescent.Enabled = false
			fmt.Println("❌ Astrolescent MCP server disabled")
		case "github":
			cfg.MCP.GitHub.Enabled = false
			fmt.Println("❌ GitHub MCP server disabled")
		default:
			log.Fatalf("Unknown MCP server: %s", server)
		}

		if err := config.SaveConfig("", cfg); err != nil {
			log.Fatalf("Failed to save config: %v", err)
		}
	},
}

func init() {
	mcpCmd.AddCommand(mcpEnableCmd)
	mcpCmd.AddCommand(mcpDisableCmd)
	rootCmd.AddCommand(mcpCmd)
}
