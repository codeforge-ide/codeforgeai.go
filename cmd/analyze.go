package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/nathfavour/codeforgeai.go/engine"
	"github.com/nathfavour/codeforgeai.go/mcp/astro"
	"github.com/spf13/cobra"
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze [path]",
	Short: "Analyze code with AI and optional MCP data",
	Long:  "Analyze your codebase using AI models with optional real-time data from MCP servers",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := "."
		if len(args) > 0 {
			path = args[0]
		}

		mcpFlag, _ := cmd.Flags().GetString("mcp")
		query, _ := cmd.Flags().GetString("query")

		eng := engine.NewEngine()

		if mcpFlag == "astrolescent" && query != "" {
			astroMCP := astro.NewAstroMCP()
			// Add MCP context to the analysis
			ctx := context.Background()

			// Example: Get live data for DeFi analysis
			if query != "" {
				price, err := astroMCP.GetPrice(ctx, "ASTRL")
				if err != nil {
					log.Printf("Failed to get price data: %v", err)
				} else {
					fmt.Printf("📊 Live Data: %s\n", price)
				}
			}
		}

		result, err := eng.AnalyzeProject(path, query)
		if err != nil {
			log.Fatalf("Analysis failed: %v", err)
		}

		fmt.Println("🔍 Analysis Results:")
		fmt.Println(result)
	},
}

func init() {
	analyzeCmd.Flags().String("mcp", "", "Enable MCP integration (astrolescent, github)")
	analyzeCmd.Flags().String("query", "", "Specific query for analysis")
	analyzeCmd.Flags().String("focus", "", "Focus area (security, performance, etc)")
	rootCmd.AddCommand(analyzeCmd)
}
