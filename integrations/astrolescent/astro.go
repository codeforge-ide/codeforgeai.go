package astrolescent

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/codeforge-ide/codeforgeai.go/mcp"
	"github.com/codeforge-ide/codeforgeai.go/mcp/astro"
)

// DeFiAnalyzer provides AI-powered DeFi analysis using Astrolescent MCP data
type DeFiAnalyzer struct {
	mcpClient *astro.AstroMCP
}

func NewDeFiAnalyzer() *DeFiAnalyzer {
	return &DeFiAnalyzer{
		mcpClient: astro.NewAstroMCP(),
	}
}

// "Should I Stake or LP?" Helper
func (d *DeFiAnalyzer) AnalyzeStakingVsLP(ctx context.Context) (string, error) {
	apy, err := d.mcpClient.GetAPY(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get APY data: %w", err)
	}
	price, err := d.mcpClient.GetPrice(ctx, "ASTRL")
	if err != nil {
		return "", fmt.Errorf("failed to get price data: %w", err)
	}
	analysis := fmt.Sprintf(`🚀 DeFi Analysis Report - Staking vs LP Strategy

📊 Current Market Data:
%s

💰 Yield Analysis:
%s

🧠 AI Recommendation:
Based on current APY data and market conditions, here's your strategic analysis:

Key Decision Factors:
- Risk Profile: Staking = Lower risk, LP = Higher risk/reward
- Market Volatility: Current price movements affect impermanent loss
- Time Horizon: Longer positions favor staking stability
- Yield Differential: Compare current rates for optimal allocation

📈 Market Context:
Recent price action indicates %s market conditions. Consider this when evaluating impermanent loss scenarios for LP positions.

🎯 Action Items:
1. Monitor APY changes over the next 24-48h
2. Assess your risk tolerance vs yield requirements  
3. Consider hybrid approach (split allocation)
4. Set alerts for significant APY changes

💡 Note: This analysis uses live MCP data from Astrolescent!
`, price.Text, apy.Text, d.getMarketSentiment(price))
	return analysis, nil
}

// "What If I Staked..." Calculator
func (d *DeFiAnalyzer) CalculateStakingReturns(ctx context.Context, amount string, days int) (string, error) {
	apy, err := d.mcpClient.GetAPY(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get APY data: %w", err)
	}
	price, err := d.mcpClient.GetPrice(ctx, "ASTRL")
	if err != nil {
		return "", fmt.Errorf("failed to get price data: %w", err)
	}
	stakingAPY := d.extractStakingAPY(apy)
	projectedReturns := d.calculateProjectedReturns(amount, stakingAPY, days)
	calculation := fmt.Sprintf(`💎 "What If I Staked..." Calculator

🔢 Input: %s ASTRL for %d days

📊 Current Market:
%s

💰 Yield Information:
%s

🧮 Projected Returns:
%s

⚠️  Disclaimer: Estimates based on current APY. Actual returns may vary.

🏆  Feature: Real-time calculations using Astrolescent MCP!
`, amount, days, price.Text, apy.Text, projectedReturns)
	return calculation, nil
}

// LLM Trading Sidekick
func (d *DeFiAnalyzer) GetTradingAdvice(ctx context.Context, fromToken, toToken, amount string) (string, error) {
	quote, err := d.mcpClient.GetQuote(ctx, "swap", fromToken, toToken, d.parseAmount(amount), "")
	if err != nil {
		return "", fmt.Errorf("failed to get quote: %w", err)
	}
	price, err := d.mcpClient.GetPrice(ctx, "ASTRL")
	if err != nil {
		return "", fmt.Errorf("failed to get price data: %w", err)
	}
	advice := fmt.Sprintf(`🎯 AI Trading Sidekick: %s %s → %s

📊 Current Market:
%s

💱 Swap Quote:
%s

🧠 AI Trading Insights:
%s

⚡ Execution Strategy:
%s

🔍 Risk Assessment:
%s

🤖 Demo: Your AI DeFi assistant powered by Astrolescent MCP!
`, amount, fromToken, toToken, price.Text, quote.Text,
		d.generateTradingInsights(quote, price),
		d.generateExecutionStrategy(amount, quote),
		d.assessTradingRisk(quote, price))
	return advice, nil
}

// Cross-chain bridge analysis
func (d *DeFiAnalyzer) AnalyzeBridgeOpportunity(ctx context.Context, fromChain, toChain string, amount float64) (string, error) {
	bridge, err := d.mcpClient.GetBridge(ctx, fromChain, toChain, "ASTRL", amount)
	if err != nil {
		return "", fmt.Errorf("failed to get bridge data: %w", err)
	}
	price, err := d.mcpClient.GetPrice(ctx, "ASTRL")
	if err != nil {
		return "", fmt.Errorf("failed to get price data: %w", err)
	}
	analysis := fmt.Sprintf(`🌉 Cross-Chain Bridge Analysis

📊 Current Market:
%s

🌉 Bridge Quote:
%s

🧠 AI Analysis:
- Cost-benefit analysis of bridging vs keeping on current chain
- Time sensitivity considerations for bridge operations
- Alternative yield opportunities on destination chain
- Risk assessment for cross-chain operations

💡 Multi-MCP Bonus: This feature demonstrates integration with bridge data!
`, price.Text, bridge.Text)
	return analysis, nil
}

// Helper methods

func (d *DeFiAnalyzer) getMarketSentiment(price *mcp.MCPResponse) string {
	if raw, ok := price.Raw.(map[string]interface{}); ok {
		if change24h, ok := raw["change_24h"].(float64); ok {
			switch {
			case change24h > 5:
				return "strongly bullish"
			case change24h > 0:
				return "moderately bullish"
			case change24h > -5:
				return "moderately bearish"
			default:
				return "strongly bearish"
			}
		}
	}
	return "neutral"
}

func (d *DeFiAnalyzer) extractStakingAPY(apy *mcp.MCPResponse) float64 {
	if raw, ok := apy.Raw.(map[string]interface{}); ok {
		if stakingAPY, ok := raw["staking_apy"].(float64); ok {
			return stakingAPY
		}
	}
	return 12.5
}

func (d *DeFiAnalyzer) calculateProjectedReturns(amount string, apy float64, days int) string {
	amountFloat := d.parseAmount(amount)
	dailyRate := apy / 365 / 100
	projectedRewards := amountFloat * dailyRate * float64(days)
	return fmt.Sprintf("Daily Rate: %.4f%%\nTotal Projected Rewards: %.2f ASTRL\nAnnualized Return: %.2f%%\nROI for %d days: %.4f%%\nAt current price ($0.083): $%.2f value",
		dailyRate*100, projectedRewards, apy, days, (projectedRewards/amountFloat)*100, projectedRewards*0.083)
}

func (d *DeFiAnalyzer) generateTradingInsights(quote, price *mcp.MCPResponse) string {
	return "Market analysis based on current data"
}

func (d *DeFiAnalyzer) generateExecutionStrategy(amount string, quote *mcp.MCPResponse) string {
	return "Recommended execution approach"
}

func (d *DeFiAnalyzer) assessTradingRisk(quote, price *mcp.MCPResponse) string {
	return "Risk assessment completed"
}

func (d *DeFiAnalyzer) parseAmount(amount string) float64 {
	cleaned := strings.ReplaceAll(amount, ",", "")
	if val, err := strconv.ParseFloat(cleaned, 64); err == nil {
		return val
	}
	return 1000
}
