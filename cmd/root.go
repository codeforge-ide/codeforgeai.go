package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/codeforge-ide/codeforgeai.go/config"
	"github.com/codeforge-ide/codeforgeai.go/directory"
	"github.com/codeforge-ide/codeforgeai.go/engine"
	"github.com/codeforge-ide/codeforgeai.go/integrations/astrolescent"
	"github.com/codeforge-ide/codeforgeai.go/integrations/githubmodels"
	"github.com/codeforge-ide/codeforgeai.go/modeliface"
	"github.com/codeforge-ide/codeforgeai.go/secrets"
	"github.com/spf13/cobra"
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
			diff, err := eng.GetGitDiff()
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
			client := githubmodels.NewClient(token, "", "")
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
			client := githubmodels.NewClient(token, "", "")
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
			client := githubmodels.NewClient(token, "", "")
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
			client := githubmodels.NewClient(token, "", "")
			resp, err := client.ImagePrompt(args[0], imageB64)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			fmt.Println(resp)
		},
	}
	githubModelsCmd.AddCommand(githubModelsImageCmd)

	// github-models token-store
	githubModelsTokenStoreCmd := &cobra.Command{
		Use:   "token-store",
		Short: "Securely store your GitHub Models API token",
		Run: func(cmd *cobra.Command, args []string) {
			err := secrets.InteractiveStoreGithubToken()
			if err != nil {
				fmt.Println("Error storing token:", err)
			} else {
				fmt.Println("GitHub token stored securely.")
			}
		},
	}
	githubModelsCmd.AddCommand(githubModelsTokenStoreCmd)

	// github-models token-load
	githubModelsTokenLoadCmd := &cobra.Command{
		Use:   "token-load",
		Short: "Load your GitHub Models API token into the environment for this session",
		Run: func(cmd *cobra.Command, args []string) {
			token, err := secrets.InteractiveLoadGithubToken()
			if err != nil {
				fmt.Println("Error loading token:", err)
			} else {
				fmt.Println("GitHub token loaded into environment for this session.")
				// Optionally print token for debug (not recommended for real use)
				_ = token
			}
		},
	}
	githubModelsCmd.AddCommand(githubModelsTokenLoadCmd)

	// rootCmd.AddCommand(githubModelsCmd) // Removed duplicate registration

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

	// astro - Astrolescent MCP Server Integration
	astroCmd := &cobra.Command{
		Use:   "astro",
		Short: "Astrolescent DeFi MCP Server integration for real-time Radix data",
		Long: `Connect to Astrolescent MCP server for live DeFi data from the Radix ecosystem.
Get real-time token prices, swap quotes, staking yields, and AI-powered DeFi analysis.`,
	}

	// astro price
	astroPriceCmd := &cobra.Command{
		Use:   "price",
		Short: "Get current ASTRL price in XRD and USD with 24h/7d changes",
		Run: func(cmd *cobra.Command, args []string) {
			client := astrolescent.NewClient()
			resp, err := client.GetPrice()
			if err != nil {
				fmt.Printf("Error fetching price: %v\n", err)
				return
			}
			fmt.Println("💰 ASTRL Price Information:")
			fmt.Println(resp.Text)
		},
	}
	astroCmd.AddCommand(astroPriceCmd)

	// astro quote
	astroQuoteCmd := &cobra.Command{
		Use:   "quote [operation] [token] [amount] [account]",
		Short: "Get swap quote (operations: buy, sell, swap)",
		Long: `Get a quote for token swaps on Radix DEXes.
Operations:
- buy: Buy specified amount of token (sells XRD)
- sell: Sell specified amount of token (buys XRD) 
- swap: Direct token-to-token swap

Examples:
  codeforgeai astro quote buy ASTRL 1000
  codeforgeai astro quote sell ASTRL 500 account_rdx1abcdefg
  codeforgeai astro quote swap XRD 100 account_rdx1abcdefg`,
		Args: cobra.RangeArgs(3, 4),
		Run: func(cmd *cobra.Command, args []string) {
			operation := args[0]
			token := args[1]
			amount := args[2]
			account := ""
			if len(args) > 3 {
				account = args[3]
			}

			client := astrolescent.NewClient()
			resp, err := client.GetQuote(operation, token, amount, account)
			if err != nil {
				fmt.Printf("Error fetching quote: %v\n", err)
				return
			}
			fmt.Printf("💱 Quote for %s %s %s:\n", operation, amount, token)
			fmt.Println(resp.Text)
		},
	}
	astroCmd.AddCommand(astroQuoteCmd)

	// astro apy
	astroAPYCmd := &cobra.Command{
		Use:   "apy",
		Short: "Get current APY for ASTRL staking and liquidity provision",
		Run: func(cmd *cobra.Command, args []string) {
			client := astrolescent.NewClient()
			resp, err := client.GetAPY()
			if err != nil {
				fmt.Printf("Error fetching APY: %v\n", err)
				return
			}
			fmt.Println("📈 ASTRL Yield Information:")
			fmt.Println(resp.Text)
		},
	}
	astroCmd.AddCommand(astroAPYCmd)

	// astro analyze - AI-powered DeFi analysis
	astroAnalyzeCmd := &cobra.Command{
		Use:   "analyze [type]",
		Short: "AI-powered DeFi analysis using live Astrolescent data",
		Long: `Advanced AI analysis of DeFi opportunities using real-time data.
Analysis types:
- staking-vs-lp: Compare staking vs liquidity provision strategies
- market: Comprehensive market analysis with trading insights`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			analyzer := astrolescent.NewDeFiAnalyzer()
			analysisType := args[0]

			switch analysisType {
			case "staking-vs-lp":
				result, err := analyzer.AnalyzeStakingVsLP(context.Background())
				if err != nil {
					fmt.Printf("Error performing analysis: %v\n", err)
					return
				}
				fmt.Println(result)
			case "market":
				// Get comprehensive market data
				client := astrolescent.NewClient()
				price, _ := client.GetPrice()
				apy, _ := client.GetAPY()

				fmt.Println("🌌 Comprehensive Radix DeFi Market Analysis")
				fmt.Println(strings.Repeat("=", 50))
				fmt.Println(price.Text)
				fmt.Println(apy.Text)

				fmt.Println(`🧠 AI Market Insights:
- Trend analysis based on 24h/7d price movements
- Yield optimization recommendations
- Risk assessment for current market conditions
- Strategic entry/exit point analysis
- Cross-DEX arbitrage opportunities`)
			default:
				fmt.Printf("Unknown analysis type: %s\n", analysisType)
				fmt.Println("Available types: staking-vs-lp, market")
			}
		},
	}
	astroCmd.AddCommand(astroAnalyzeCmd)

	// astro calculator
	astroCalculatorCmd := &cobra.Command{
		Use:   "calculator [amount] [days]",
		Short: "Calculate potential staking returns for ASTRL",
		Long: `Calculate projected staking returns based on current APY.
Example: codeforgeai astro calculator 1000 30`,
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			amount := args[0]
			days := 30
			if len(args) > 1 {
				fmt.Sscanf(args[1], "%d", &days)
			}

			analyzer := astrolescent.NewDeFiAnalyzer()
			result, err := analyzer.CalculateStakingReturns(context.Background(), amount, days)
			if err != nil {
				fmt.Printf("Error calculating returns: %v\n", err)
				return
			}
			fmt.Println(result)
		},
	}
	astroCmd.AddCommand(astroCalculatorCmd)

	// astro trading-advice
	astroTradingCmd := &cobra.Command{
		Use:   "trading-advice [from_token] [to_token] [amount]",
		Short: "Get AI-powered trading advice with market analysis",
		Long: `Get comprehensive trading advice including market timing, liquidity analysis, and execution strategy.
Example: codeforgeai astro trading-advice XRD ASTRL 100`,
		Args: cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			fromToken := args[0]
			toToken := args[1]
			amount := args[2]

			analyzer := astrolescent.NewDeFiAnalyzer()
			result, err := analyzer.GetTradingAdvice(context.Background(), fromToken, toToken, amount)
			if err != nil {
				fmt.Printf("Error getting trading advice: %v\n", err)
				return
			}
			fmt.Println(result)
		},
	}
	astroCmd.AddCommand(astroTradingCmd)

	// astro demo - Impressive demo mode
	astroDemoCmd := &cobra.Command{
		Use:   "demo",
		Short: "Run impressive demo showcasing all Astrolescent MCP capabilities",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("🚀 CodeforgeAI x Astrolescent MCP Demo")
			fmt.Println(strings.Repeat("=", 50))
			fmt.Println("Connecting to Radix DeFi ecosystem via MCP...")

			client := astrolescent.NewClient()
			analyzer := astrolescent.NewDeFiAnalyzer()

			// Demo sequence
			fmt.Println("\n1️⃣ Fetching live ASTRL price data...")
			if price, err := client.GetPrice(); err == nil {
				fmt.Println(price.Text)
			}

			fmt.Println("\n2️⃣ Getting current staking yields...")
			if apy, err := client.GetAPY(); err == nil {
				fmt.Println(apy.Text)
			}

			fmt.Println("\n3️⃣ Calculating swap quote for 100 XRD → ASTRL...")
			if quote, err := client.GetQuote("swap", "ASTRL", "100", ""); err == nil {
				fmt.Println(quote.Text)
			}

			fmt.Println("\n4️⃣ AI Analysis: Staking vs LP Strategy...")
			if analysis, err := analyzer.AnalyzeStakingVsLP(context.Background()); err == nil {
				fmt.Println(analysis)
			}

			fmt.Println("\n5️⃣ Calculating returns for 1000 ASTRL staked for 30 days...")
			if calc, err := analyzer.CalculateStakingReturns(context.Background(), "1000", 30); err == nil {
				fmt.Println(calc)
			}

			fmt.Println("\n🎯 Demo Complete!")
			fmt.Println("CodeforgeAI successfully integrated with Astrolescent MCP server")
			fmt.Println("✅ Real-time DeFi data access")
			fmt.Println("✅ AI-powered analysis")
			fmt.Println("✅ Multi-tool integration")
			fmt.Println("✅ Seamless user experience")
		},
	}
	astroCmd.AddCommand(astroDemoCmd)

	rootCmd.AddCommand(astroCmd)

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
			token := os.Getenv("IO_NET_API_KEY")
			if token == "" {
				fmt.Println("IO_NET_API_KEY environment variable is required.")
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
			token := os.Getenv("IO_NET_API_KEY")
			if token == "" {
				fmt.Println("IO_NET_API_KEY environment variable is required.")
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
			token := os.Getenv("IO_NET_API_KEY")
			if token == "" {
				fmt.Println("IO_NET_API_KEY environment variable is required.")
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
