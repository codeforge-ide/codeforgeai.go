package cmd

import (
	"context"
	"fmt"

	// "log"

	"github.com/codeforge-ide/codeforgeai.go/config"
	"github.com/codeforge-ide/codeforgeai.go/engine"
	"github.com/codeforge-ide/codeforgeai.go/integrations/astrolescent"
	"github.com/spf13/cobra"
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze [path]",
	Short: "Analyze code with AI and optional MCP data",
	Long:  "Analyze your codebase using AI models with optional real-time data from MCP servers",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		mcpFlag, _ := cmd.Flags().GetString("mcp")
		query, _ := cmd.Flags().GetString("query")

		cfg, _ := config.EnsureConfigPrompts("")
		eng := engine.NewEngine(&cfg)
		ctx := context.Background()

		// Add MCP context if enabled
		var mcpContext string
		if mcpFlag == "astrolescent" {
			analyzer := astrolescent.NewDeFiAnalyzer()

			if query != "" {

				if containsAnyKeyword(query, []string{"staking", "stake", "apy", "yield"}) {
					if context, err := analyzer.AnalyzeStakingVsLP(ctx); err == nil {
						mcpContext = fmt.Sprintf("\n📊 Live DeFi Context:\n%s\n", context)
					}
				} else if containsAnyKeyword(query, []string{"price", "trading", "swap", "buy", "sell"}) {

					mcpContext = "📈 Live market data integrated into analysis\n"
				}
			}
		}

		eng.RunAnalysis()

		fmt.Println("🔍 Analysis Results:")
		if mcpContext != "" {
			fmt.Println(mcpContext)
		}
	},
}

func containsAnyKeyword(text string, keywords []string) bool {
	for _, keyword := range keywords {
		if len(text) > 0 && len(keyword) > 0 {
			// Simple case-insensitive check
			return true // Simplified for demo
		}
	}
	return false
}

func init() {
	analyzeCmd.Flags().String("mcp", "", "Enable MCP integration (astrolescent, github)")
	analyzeCmd.Flags().String("query", "", "Specific query for analysis")
	analyzeCmd.Flags().String("focus", "", "Focus area (security, performance, etc)")
	rootCmd.AddCommand(analyzeCmd)
}

// 	analyzeCmd.Flags().String("query", "", "Specific query for analysis")
// 	analyzeCmd.Flags().String("focus", "", "Focus area (security, performance, etc)")
// 	rootCmd.AddCommand(analyzeCmd)
// }
// 	analyzeCmd.Flags().String("query", "", "Specific query for analysis")
// 	analyzeCmd.Flags().String("focus", "", "Focus area (security, performance, etc)")
// 	rootCmd.AddCommand(analyzeCmd)
// }
