package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	verbose     bool
	veryVerbose bool
	debug       bool
	userPrompt  []string
	filePath    string
	loop        bool
)

var rootCmd = &cobra.Command{
	Use:   "codeforgeai",
	Short: "CodeforgeAI AI agent",
	Long:  "CodeforgeAI AI agent - Go CLI",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "set loglevel to INFO")
	rootCmd.PersistentFlags().BoolVarP(&veryVerbose, "very-verbose", "V", false, "set loglevel to DEBUG")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug mode (overrides other verbosity flags)")

	// analyze
	analyzeCmd := &cobra.Command{
		Use:   "analyze",
		Short: "Analyze current working directory",
		Run: func(cmd *cobra.Command, args []string) {
			if loop {
				fmt.Println("Running analysis loop (not implemented in Go yet).")
			} else {
				fmt.Println("Analyzing current working directory (not implemented in Go yet).")
			}
		},
	}
	analyzeCmd.Flags().BoolVar(&loop, "loop", false, "Enable adaptive feedback loop")
	rootCmd.AddCommand(analyzeCmd)

	// prompt
	promptCmd := &cobra.Command{
		Use:   "prompt [user_prompt]",
		Short: "Process a user prompt",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Processing prompt: %s\n", strings.Join(args, " "))
			fmt.Println("(Prompt processing not implemented in Go yet.)")
		},
	}
	rootCmd.AddCommand(promptCmd)

	// config
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Run configuration checkup",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Configuration checkup complete. (Not implemented in Go yet.)")
		},
	}
	rootCmd.AddCommand(configCmd)

	// strip
	stripCmd := &cobra.Command{
		Use:   "strip",
		Short: "Print tree structure after removing gitignored files",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Printing tree structure (not implemented in Go yet).")
		},
	}
	rootCmd.AddCommand(stripCmd)

	// commit-message
	commitMsgCmd := &cobra.Command{
		Use:   "commit-message",
		Short: "Generate commit message with code changes and gitmoji",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Generating commit message (not implemented in Go yet).")
		},
	}
	rootCmd.AddCommand(commitMsgCmd)

	// github (with copilot subcommands)
	githubCmd := &cobra.Command{
		Use:   "github",
		Short: "GitHub Copilot integration",
	}
	copilotCmd := &cobra.Command{
		Use:   "copilot",
		Short: "Github copilot integration",
	}
	copilotCmd.AddCommand(&cobra.Command{
		Use:   "login",
		Short: "Authenticate with GitHub Copilot",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("GitHub Copilot login (not implemented in Go yet).")
		},
	})
	copilotCmd.AddCommand(&cobra.Command{
		Use:   "logout",
		Short: "Logout from GitHub Copilot",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("GitHub Copilot logout (not implemented in Go yet).")
		},
	})
	copilotCmd.AddCommand(&cobra.Command{
		Use:   "status",
		Short: "Check GitHub Copilot status",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("GitHub Copilot status (not implemented in Go yet).")
		},
	})
	copilotCmd.AddCommand(&cobra.Command{
		Use:   "lsp",
		Short: "install copilot language server globally",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Installing Copilot language server (not implemented in Go yet).")
		},
	})
	githubCmd.AddCommand(copilotCmd)
	rootCmd.AddCommand(githubCmd)
}
