package cmd

import (
	"fmt"
	"github.com/codeforge-ide/codeforgeai.go/config"
	"github.com/codeforge-ide/codeforgeai.go/engine"
	"github.com/spf13/cobra"
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze [path]",
	Short: "Analyze code with AI and optional MCP data",
	Long:  "Analyze your codebase using AI models with optional real-time data from MCP servers",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		cfg, _ := config.EnsureConfigPrompts("")
		eng := engine.NewEngine(&cfg)
		eng.RunAnalysis()
		fmt.Println("🔍 Analysis Results:")

	},
}

// containsAnyKeyword was a demo stub and is now removed.

func init() {
	analyzeCmd.Flags().String("mcp", "", "Enable MCP integration (astrolescent, github)")
	analyzeCmd.Flags().String("query", "", "Specific query for analysis")
	analyzeCmd.Flags().String("focus", "", "Focus area (security, performance, etc)")
	// rootCmd.AddCommand(analyzeCmd) // Registered in root.go only
}

// 	analyzeCmd.Flags().String("query", "", "Specific query for analysis")
// 	analyzeCmd.Flags().String("focus", "", "Focus area (security, performance, etc)")
// 	rootCmd.AddCommand(analyzeCmd)
// }
// 	analyzeCmd.Flags().String("query", "", "Specific query for analysis")
// 	analyzeCmd.Flags().String("focus", "", "Focus area (security, performance, etc)")
// 	rootCmd.AddCommand(analyzeCmd)
// }
