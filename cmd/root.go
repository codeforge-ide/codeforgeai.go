package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"encoding/json"
	"github.com/codeforge-ide/codeforgeai.go/config"
	"github.com/codeforge-ide/codeforgeai.go/directory"
	"github.com/codeforge-ide/codeforgeai.go/engine"
	_ "github.com/codeforge-ide/codeforgeai.go/integrations/astrolescent"
	_ "github.com/codeforge-ide/codeforgeai.go/integrations/coingecko"
	_ "github.com/codeforge-ide/codeforgeai.go/integrations/githubmodels"
	"github.com/codeforge-ide/codeforgeai.go/modeliface"
	"github.com/codeforge-ide/codeforgeai.go/utils"
	"github.com/spf13/cobra"
	"net/http"
)

import io_integration "github.com/codeforge-ide/codeforgeai.go/integrations/io"

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

// runGitCommand runs a git command with the given arguments.
func runGitCommand(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no git args provided")
	}
	gitArgs := append([]string{"git"}, args...)
	proc := exec.Command(gitArgs[0], gitArgs[1:]...)
	proc.Stdout = os.Stdout
	proc.Stderr = os.Stderr
	return proc.Run()
}

// ---- END UTILS ----

func init() {
	// Dynamically add all registered integrations as CLI commands
	for _, reg := range modeliface.GlobalAgentRegistry.ListIntegrations() {
		rootCmd.AddCommand(reg.CommandFactory())
	}

	// --- Integrations Command ---
	integrationsCmd := &cobra.Command{
		Use:   "integrations",
		Short: "Manage and inspect AI integrations and models",
	}

	integrationsListCmd := &cobra.Command{
		Use:   "list",
		Short: "List all available integrations/providers",
		Run: func(cmd *cobra.Command, args []string) {
			for _, reg := range modeliface.GlobalAgentRegistry.ListIntegrations() {
				fmt.Printf("- %s: %s\n", reg.Metadata.Name, reg.Metadata.Description)
			}
		},
	}
	integrationsCmd.AddCommand(integrationsListCmd)

	integrationsListModelsCmd := &cobra.Command{
		Use:   "[provider] list-models",
		Short: "List available models for a provider",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			provider := args[0]
			switch provider {
			case "githubcopilot", "copilot":
				out, err := exec.Command("gh", "copilot", "models", "list").CombinedOutput()
				if err != nil {
					fmt.Printf("Error listing Copilot models: %v\n%s\n", err, string(out))
					return
				}
				fmt.Println(strings.TrimSpace(string(out)))
			case "ollama":
				endpoint := "http://localhost:11434/api/tags"
				resp, err := http.Get(endpoint)
				if err != nil {
					fmt.Printf("Error connecting to Ollama: %v\n", err)
					return
				}
				defer resp.Body.Close()
				var tags struct {
					Models []string `json:"models"`
				}
				if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil {
					fmt.Printf("Error decoding Ollama response: %v\n", err)
					return
				}
				for _, m := range tags.Models {
					fmt.Println(m)
				}
			default:
				fmt.Printf("Model listing not implemented for provider: %s\n", provider)
			}
		},
	}
	integrationsCmd.AddCommand(integrationsListModelsCmd)

	integrationsSetModelCmd := &cobra.Command{
		Use:   "[provider] set-model [model]",
		Short: "Set the default model for a provider",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			provider := args[0]
			model := args[2]
			cfg, _ := config.LoadConfig("")
			switch provider {
			case "githubcopilot", "copilot":
				cfg.CodeModel = model
				cfg.OllamaModel = model // for compatibility if needed
				fmt.Printf("Set Copilot model to: %s\n", model)
			case "ollama":
				cfg.OllamaModel = model
				fmt.Printf("Set Ollama model to: %s\n", model)
			default:
				fmt.Printf("Model setting not implemented for provider: %s\n", provider)
				return
			}
			config.SaveConfig("", cfg)
		},
	}
	integrationsCmd.AddCommand(integrationsSetModelCmd)

	rootCmd.AddCommand(integrationsCmd)

	// push command
	pushCmd := &cobra.Command{
		Use:   "push",
		Short: "Stage, commit, and push all changes with a single-line message",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Staging all changes...")
			if err := runGitCommand([]string{"add", "."}); err != nil {
				fmt.Println("git add failed:", err)
				return
			}
			var msg string
			if len(args) > 0 {
				msg = strings.Join(args, " ")
			} else {
				fmt.Print("Enter a single-line commit message: ")
				fmt.Scanln(&msg)
			}
			if msg == "" {
				fmt.Println("Commit message required.")
				return
			}
			fmt.Println("Committing...")
			if err := runGitCommand([]string{"commit", "-m", msg}); err != nil {
				fmt.Println("git commit failed:", err)
				return
			}
			fmt.Println("Pushing to origin...")
			if err := runGitCommand([]string{"push"}); err != nil {
				fmt.Println("git push failed:", err)
				return
			}
			fmt.Println("Push complete.")
		},
	}
	rootCmd.AddCommand(pushCmd)

	// Global flags
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "set loglevel to INFO")
	rootCmd.PersistentFlags().BoolVarP(&veryVerbose, "very-verbose", "V", false, "set loglevel to DEBUG")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug mode (overrides other verbosity flags)")

	// analyze
	analyzeCmd := &cobra.Command{
		Use:   "analyze",
		Short: "Analyze current working directory",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, _ := config.EnsureConfigPrompts("")
			eng := engine.NewEngine(&cfg)
			eng.RunAnalysis()
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
			cfg, _ := config.EnsureConfigPrompts("")
			eng := engine.NewEngine(&cfg)
			resp := eng.ProcessPrompt(strings.Join(args, " "))
			fmt.Println(resp)
		},
	}
	rootCmd.AddCommand(promptCmd)

	// config
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Run configuration checkup",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, _ := config.EnsureConfigPrompts("")
			fmt.Println("Configuration checkup complete. Current configuration:")
			config.PrintConfig(cfg)
		},
	}
	rootCmd.AddCommand(configCmd)

	// strip
	stripCmd := &cobra.Command{
		Use:   "strip",
		Short: "Print tree structure after removing gitignored files",
		Run: func(cmd *cobra.Command, args []string) {
			directory.StripDirectory()
		},
	}
	rootCmd.AddCommand(stripCmd)

	// commit-message
	commitMsgCmd := &cobra.Command{
		Use:   "commit-message",
		Short: "Generate commit message with code changes and gitmoji",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, _ := config.EnsureConfigPrompts("")
			eng := engine.NewEngine(&cfg)

			// Get git diff
			diff, err := utils.GetGitDiff()
			if err != nil {
				fmt.Println("Error getting git diff:", err)
				return
			}
			if strings.TrimSpace(diff) == "" {
				fmt.Println("No changes detected in git")
				return
			}

			resp := eng.ProcessCommitMessage(diff)
			fmt.Println(resp)
		},
	}
	rootCmd.AddCommand(commitMsgCmd)

	// copilot command
	rootCmd.AddCommand(copilotCmd)

	// explain
	explainCmd := &cobra.Command{
		Use:   "explain [file_path]",
		Short: "Explain the code in the given file",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cfg, _ := config.EnsureConfigPrompts("")
			eng := engine.NewEngine(&cfg)
			resp := eng.ExplainCode(args[0])
			fmt.Println(resp)
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
			userPrompts, _ := cmd.Flags().GetStringSlice("user_prompt")
			allowIgnore, _ := cmd.Flags().GetBool("allow-ignore")

			if len(userPrompts) == 0 {
				fmt.Println("Error: --user_prompt is required")
				return
			}

			userPrompt := strings.Join(userPrompts, " ")
			paths := args
			if len(paths) == 0 {
				paths = []string{"."} // Default to current directory
			}

			cfg, _ := config.EnsureConfigPrompts("")
			eng := engine.NewEngine(&cfg)
			err := eng.EditFiles(paths, userPrompt, allowIgnore)
			if err != nil {
				fmt.Println("Error editing files:", err)
			} else {
				fmt.Println("Edit complete. Check .codeforgedit files for results.")
			}
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
			filePath, _ := cmd.Flags().GetString("file")
			line, _ := cmd.Flags().GetInt("line")
			snippets, _ := cmd.Flags().GetStringSlice("string")
			entire, _ := cmd.Flags().GetBool("entire")

			cfg, _ := config.EnsureConfigPrompts("")
			eng := engine.NewEngine(&cfg)
			resp := eng.ProvideSuggestion(filePath, line, snippets, entire)
			fmt.Println(resp)
		},
	}
	suggestionCmd.Flags().String("file", "", "File to read code from")
	suggestionCmd.Flags().Int("line", 0, "Line number to use for suggestion")
	suggestionCmd.Flags().StringSlice("string", nil, "User-provided code snippet for suggestion")
	suggestionCmd.Flags().BoolP("entire", "E", false, "Send entire file content for suggestion")
	rootCmd.AddCommand(suggestionCmd)

	// secret-ai, web3, and zerepy commands are stubs and removed for now.

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



	// --- Enable/Disable Integration/Extension Commands ---
	// Removed enable/disable extension commands and stubs (web3, zerepy, secret-ai) as they are not implemented or needed.

	// Opencode integration
	opencodeCmd := &cobra.Command{
		Use:   "opencode",
		Short: "Opencode AI integration commands",
	}
	rootCmd.AddCommand(opencodeCmd)

	ioCmd := &cobra.Command{
		Use:   "io",
		Short: "Interact with IO.net API",
	}
	rootCmd.AddCommand(ioCmd)

	ioListModelsCmd := &cobra.Command{
		Use:   "list-models",
		Short: "List available IO.net models",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, _ := config.LoadConfig("")
			token := cfg.IONetAPIKey
			if token == "" {
				token = os.Getenv("IO_NET_API_KEY")
			}
			if token == "" {
				fmt.Println("IO_NET_API_KEY environment variable or config is required.")
				return
			}
			client := io_integration.New(token)
			models, err := client.ListModels(context.Background())
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			for _, model := range models {
				fmt.Printf("- %v\n", model)
			}
		},
	}
	ioCmd.AddCommand(ioListModelsCmd)

	ioListAgentsCmd := &cobra.Command{
		Use:   "list-agents",
		Short: "List available IO.net agents",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, _ := config.LoadConfig("")
			token := cfg.IONetAPIKey
			if token == "" {
				token = os.Getenv("IO_NET_API_KEY")
			}
			if token == "" {
				fmt.Println("IO_NET_API_KEY environment variable or config is required.")
				return
			}
			client := io_integration.New(token)
			agents, err := client.ListAgents(context.Background())
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			for _, agent := range agents {
				fmt.Printf("- %s: %s\n", agent.ID, agent.Name)
			}
		},
	}
	ioCmd.AddCommand(ioListAgentsCmd)

	ioRunAgentCmd := &cobra.Command{
		Use:   "run-agent [agent-id] [prompt]",
		Short: "Run an IO.net agent",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			cfg, _ := config.LoadConfig("")
			token := cfg.IONetAPIKey
			if token == "" {
				token = os.Getenv("IO_NET_API_KEY")
			}
			if token == "" {
				fmt.Println("IO_NET_API_KEY environment variable or config is required.")
				return
			}
			client := io_integration.New(token)
			req := &io_integration.AgentRequest{
				Model:   "gpt-3.5-turbo",
				Content: args[1],
			}
			resp, err := client.RunAgent(context.Background(), args[0], req)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			fmt.Println(resp.Content)
		},
	}
	ioCmd.AddCommand(ioRunAgentCmd)
}

// Helper function for base64 encoding
func encodeToBase64(data []byte) string {
	return strings.TrimRight(strings.ReplaceAll(fmt.Sprintf("%+q", data), "\\x", ""), "\"")
}

// setIntegrationEnabled enables/disables an integration by name (case-insensitive).
func setIntegrationEnabled(cfg *config.Config, name string, enabled bool) (bool, error) {
	switch name {
	case "ollama":
		if cfg.Integrations.Ollama.Enabled != enabled {
			cfg.Integrations.Ollama.Enabled = enabled
			return true, nil
		}
	case "githubmodels":
		if cfg.Integrations.GithubModels.Enabled != enabled {
			cfg.Integrations.GithubModels.Enabled = enabled
			return true, nil
		}
	case "openapi":
		if cfg.Integrations.OpenAPI.Enabled != enabled {
			cfg.Integrations.OpenAPI.Enabled = enabled
			return true, nil
		}
	case "githubcopilot":
		if cfg.Integrations.GithubCopilot.Enabled != enabled {
			cfg.Integrations.GithubCopilot.Enabled = enabled
			return true, nil
		}
	case "io":
		if cfg.Integrations.IO.Enabled != enabled {
			cfg.Integrations.IO.Enabled = enabled
			return true, nil
		}
	default:
		return false, errors.New("unknown integration: " + name)
	}
	return false, nil
}
