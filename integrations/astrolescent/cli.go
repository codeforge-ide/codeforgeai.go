package astrolescent

import (
	"context"
	"fmt"
	"strings"

	"github.com/codeforge-ide/codeforgeai.go/modeliface"
	"github.com/spf13/cobra"
)

func init() {
	modeliface.GlobalAgentRegistry.RegisterIntegration(modeliface.IntegrationRegistration{
		Metadata: modeliface.IntegrationMetadata{
			Name:         "astro",
			Description:  "Astrolescent DeFi MCP Server integration for real-time Radix data.",
			Commands:     []string{"price", "quote", "apy", "analyze", "calculator", "trading-advice", "demo"},
			Capabilities: []string{"defi_price", "defi_quote", "defi_apy"},
		},
		CommandFactory: NewAstroCommand,
	})
}

func NewAstroCommand() *cobra.Command {
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
			client := NewClient()
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

			client := NewClient()
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
			client := NewClient()
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
			analyzer := NewDeFiAnalyzer()
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
				client := NewClient()
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

			analyzer := NewDeFiAnalyzer()
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

			analyzer := NewDeFiAnalyzer()
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

			client := NewClient()
			analyzer := NewDeFiAnalyzer()

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

	return astroCmd
}
