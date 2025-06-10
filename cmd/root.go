package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/codeforge-ide/codeforgeai.go/integrations/githubmodels" // Import the GitHub Models client package
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

	githubModelsCmd := &cobra.Command{
		Use:   "github-models",
		Short: "Interact with GitHub Models API",
	}

	// github-models prompt
	githubModelsPromptCmd := &cobra.Command{
		Use:   "prompt [prompt]",
		Short: "Send a simple prompt to GitHub Models",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			token := os.Getenv("GITHUB_TOKEN")
			if token == "" {
				fmt.Println("GITHUB_TOKEN environment variable is required.")
				return
			}
			client := githubmodels.NewClient(token)
			resp, err := client.SimplePrompt(strings.Join(args, " "))
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			fmt.Println(resp)
		},
	}
	githubModelsCmd.AddCommand(githubModelsPromptCmd)

	// github-models multi-turn
	githubModelsMultiTurnCmd := &cobra.Command{
		Use:   "multi-turn",
		Short: "Send a multi-turn conversation to GitHub Models",
		Run: func(cmd *cobra.Command, args []string) {
			token := os.Getenv("GITHUB_TOKEN")
			if token == "" {
				fmt.Println("GITHUB_TOKEN environment variable is required.")
				return
			}
			// For demo, hardcode a conversation; in real use, parse from args or file
			history := []githubmodels.Message{
				{Role: "system", Content: "You are a helpful assistant."},
				{Role: "user", Content: "What is the capital of France?"},
				{Role: "assistant", Content: "The capital of France is Paris."},
				{Role: "user", Content: "What about Spain?"},
			}
			client := githubmodels.NewClient(token)
			resp, err := client.MultiTurn(history)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			fmt.Println(resp)
		},
	}
	githubModelsCmd.AddCommand(githubModelsMultiTurnCmd)

	// github-models stream
	githubModelsStreamCmd := &cobra.Command{
		Use:   "stream [prompt]",
		Short: "Stream a prompt response from GitHub Models",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			token := os.Getenv("GITHUB_TOKEN")
			if token == "" {
				fmt.Println("GITHUB_TOKEN environment variable is required.")
				return
			}
			client := githubmodels.NewClient(token)
			resp, err := client.StreamPrompt(strings.Join(args, " "))
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			fmt.Println(resp)
		},
	}
	githubModelsCmd.AddCommand(githubModelsStreamCmd)

	// github-models image
	githubModelsImageCmd := &cobra.Command{
		Use:   "image [prompt] [image_path]",
		Short: "Send an image prompt to GitHub Models",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			token := os.Getenv("GITHUB_TOKEN")
			if token == "" {
				fmt.Println("GITHUB_TOKEN environment variable is required.")
				return
			}
			imagePath := args[1]
			imageData, err := os.ReadFile(imagePath)
			if err != nil {
				fmt.Println("Error reading image:", err)
				return
			}
			imageB64 := encodeToBase64(imageData)
			client := githubmodels.NewClient(token)
			resp, err := client.ImagePrompt(args[0], imageB64)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			fmt.Println(resp)
		},
	}
	githubModelsCmd.AddCommand(githubModelsImageCmd)

	rootCmd.AddCommand(githubModelsCmd)

	// explain
	explainCmd := &cobra.Command{
		Use:   "explain [file_path]",
		Short: "Explain the code in the given file",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Explaining code in file: %s (not implemented in Go yet).\n", args[0])
		},
	}
	rootCmd.AddCommand(explainCmd)

	// extract
	extractCmd := &cobra.Command{
		Use:   "extract",
		Short: "Extract code blocks from file or string",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Extracting code blocks (not implemented in Go yet).")
		},
	}
	extractCmd.Flags().String("file", "", "Path to the file to process")
	extractCmd.Flags().String("string", "", "Input string containing code blocks")
	rootCmd.AddCommand(extractCmd)

	// format
	formatCmd := &cobra.Command{
		Use:   "format",
		Short: "Format code blocks for readability",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Formatting code blocks (not implemented in Go yet).")
		},
	}
	formatCmd.Flags().String("file", "", "Path to the file to process")
	formatCmd.Flags().String("string", "", "Input string containing code blocks")
	rootCmd.AddCommand(formatCmd)

	// command
	commandCmd := &cobra.Command{
		Use:   "command [user_command]",
		Short: "Process a command request",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Processing command: %s (not implemented in Go yet).\n", strings.Join(args, " "))
		},
	}
	rootCmd.AddCommand(commandCmd)

	// edit
	editCmd := &cobra.Command{
		Use:   "edit [paths...] --user_prompt PROMPT",
		Short: "Edit code in specified files or folders",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Editing code (not implemented in Go yet).")
		},
	}
	editCmd.Flags().StringSlice("user_prompt", nil, "User prompt for editing")
	editCmd.Flags().Bool("allow-ignore", false, "Allow explicitly passed directories to be processed even if .gitignore ignores them")
	rootCmd.AddCommand(editCmd)

	// suggestion
	suggestionCmd := &cobra.Command{
		Use:   "suggestion",
		Short: "Short suggestions from code model at lightning speed",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Providing suggestion (not implemented in Go yet).")
		},
	}
	suggestionCmd.Flags().String("file", "", "File to read code from")
	suggestionCmd.Flags().Int("line", 0, "Line number to use for suggestion")
	suggestionCmd.Flags().StringSlice("string", nil, "User-provided code snippet for suggestion")
	suggestionCmd.Flags().BoolP("entire", "E", false, "Send entire file content for suggestion")
	rootCmd.AddCommand(suggestionCmd)

	// secret-ai
	secretAICmd := &cobra.Command{
		Use:   "secret-ai",
		Short: "Secret AI SDK integration commands",
	}
	secretAICmd.AddCommand(&cobra.Command{
		Use:   "list-models",
		Short: "List available Secret AI models",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Listing Secret AI models (not implemented in Go yet).")
		},
	})
	secretAICmd.AddCommand(&cobra.Command{
		Use:   "test-connection",
		Short: "Test Secret AI connection",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Testing Secret AI connection (not implemented in Go yet).")
		},
	})
	secretAICmd.AddCommand(&cobra.Command{
		Use:   "chat [message]",
		Short: "Chat with Secret AI",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Chatting with Secret AI: %s (not implemented in Go yet).\n", strings.Join(args, " "))
		},
	})
	rootCmd.AddCommand(secretAICmd)

	// web3
	web3Cmd := &cobra.Command{
		Use:   "web3",
		Short: "Web3 development commands",
	}
	web3Cmd.AddCommand(&cobra.Command{
		Use:   "scaffold [project_name]",
		Short: "Scaffold a new web3 project",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Scaffolding web3 project (not implemented in Go yet).")
		},
	})
	web3Cmd.AddCommand(&cobra.Command{
		Use:   "analyze-contract [contract_file]",
		Short: "Analyze a smart contract",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Analyzing smart contract (not implemented in Go yet).")
		},
	})
	web3Cmd.AddCommand(&cobra.Command{
		Use:   "estimate-gas [contract_file]",
		Short: "Estimate gas costs for a smart contract",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Estimating gas costs (not implemented in Go yet).")
		},
	})
	web3Cmd.AddCommand(&cobra.Command{
		Use:   "generate-tests [contract_file]",
		Short: "Generate tests for a smart contract",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Generating tests (not implemented in Go yet).")
		},
	})
	web3Cmd.AddCommand(&cobra.Command{
		Use:   "check-env",
		Short: "Check web3 development environment",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Checking web3 environment (not implemented in Go yet).")
		},
	})
	web3Cmd.AddCommand(&cobra.Command{
		Use:   "install-deps",
		Short: "Install web3 dependencies",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Installing web3 dependencies (not implemented in Go yet).")
		},
	})
	rootCmd.AddCommand(web3Cmd)

	// zerepy
	zerepyCmd := &cobra.Command{
		Use:   "zerepy",
		Short: "ZerePy integration commands",
	}
	zerepyCmd.AddCommand(&cobra.Command{
		Use:   "status",
		Short: "Check ZerePy server status",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Checking ZerePy server status (not implemented in Go yet).")
		},
	})
	zerepyCmd.AddCommand(&cobra.Command{
		Use:   "list-agents",
		Short: "List available ZerePy agents",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Listing ZerePy agents (not implemented in Go yet).")
		},
	})
	zerepyCmd.AddCommand(&cobra.Command{
		Use:   "load-agent [agent_name]",
		Short: "Load a ZerePy agent",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Loading ZerePy agent (not implemented in Go yet).")
		},
	})
	zerepyCmd.AddCommand(&cobra.Command{
		Use:   "action [connection] [action]",
		Short: "Execute a ZerePy action",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Executing ZerePy action (not implemented in Go yet).")
		},
	})
	zerepyCmd.AddCommand(&cobra.Command{
		Use:   "chat [message]",
		Short: "Chat with a ZerePy agent",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Chatting with ZerePy agent (not implemented in Go yet).")
		},
	})
	rootCmd.AddCommand(zerepyCmd)

	// solana
	solanaCmd := &cobra.Command{
		Use:   "solana",
		Short: "Solana blockchain commands",
	}
	solanaCmd.AddCommand(&cobra.Command{
		Use:   "status",
		Short: "Check Solana Agent status",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Checking Solana Agent status (not implemented in Go yet).")
		},
	})
	solanaCmd.AddCommand(&cobra.Command{
		Use:   "balance",
		Short: "Get wallet balance",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Getting wallet balance (not implemented in Go yet).")
		},
	})
	solanaCmd.AddCommand(&cobra.Command{
		Use:   "transfer [destination] [amount]",
		Short: "Transfer SOL to an address",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Transferring SOL (not implemented in Go yet).")
		},
	})
	mcpCmd := &cobra.Command{
		Use:   "mcp",
		Short: "Solana MCP commands",
	}
	mcpCmd.AddCommand(&cobra.Command{
		Use:   "interact [program_id] [action_type]",
		Short: "Interact with an MCP",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Interacting with MCP (not implemented in Go yet).")
		},
	})
	mcpCmd.AddCommand(&cobra.Command{
		Use:   "state [program_id] [account_address]",
		Short: "Get state from an MCP",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Getting MCP state (not implemented in Go yet).")
		},
	})
	mcpCmd.AddCommand(&cobra.Command{
		Use:   "init-account [program_id] [space]",
		Short: "Initialize a new MCP account",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Initializing MCP account (not implemented in Go yet).")
		},
	})
	solanaCmd.AddCommand(mcpCmd)
	rootCmd.AddCommand(solanaCmd)
}

// Helper function for base64 encoding
func encodeToBase64(data []byte) string {
	return strings.TrimRight(strings.ReplaceAll(fmt.Sprintf("%+q", data), "\\x", ""), "\"")
}
