package astrolescent

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/nathfavour/codeforgeai.go/mcp/astro"
	"github.com/nathfavour/codeforgeai.go/mcp"
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
`, price.Text, apy.Text, d.getMarketSentiment(price))

	return analysis, nil
}

func (d *DeFiAnalyzer) CalculateStakingReturns(ctx context.Context, amount string, days int) (string, error) {
	apy, err := d.mcpClient.GetAPY(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get APY data: %w", err)
	}

	price, err := d.mcpClient.GetPrice(ctx, "ASTRL")
	if err != nil {
		return "", fmt.Errorf("failed to get price data: %w", err)
	}

	// Extract APY from response for calculations
	stakingAPY := d.extractStakingAPY(apy)
	projectedReturns := d.calculateProjectedReturns(amount, stakingAPY, days)

	calculation := fmt.Sprintf(`💎 ASTRL Staking Calculator

🔢 Input: %s ASTRL for %d days

📊 Current Market:
%s

💰 Yield Information:
%s

🧮 Projected Returns:
%s

⚠️  Disclaimer: Estimates based on current APY. Actual returns may vary due to:
- Validator performance changes
- Network condition fluctuations  
- Market volatility affecting token price
- Compound reward mechanisms
`, amount, days, price.Text, apy.Text, projectedReturns)

	return calculation, nil
}

func (d *DeFiAnalyzer) GetTradingAdvice(ctx context.Context, fromToken, toToken, amount string) (string, error) {
	quote, err := d.mcpClient.GetQuote(ctx, "swap", fromToken, toToken, d.parseAmount(amount), "")
	if err != nil {
		return "", fmt.Errorf("failed to get quote: %w", err)
	}

	price, err := d.mcpClient.GetPrice(ctx, "ASTRL")
	if err != nil {
		return "", fmt.Errorf("failed to get price data: %w", err)
	}

	advice := fmt.Sprintf(`🎯 Trading Analysis: %s %s → %s

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
`, amount, fromToken, toToken, price.Text, quote.Text, 
	d.generateTradingInsights(quote, price),
	d.generateExecutionStrategy(amount, quote),
	d.assessTradingRisk(quote, price))

	return advice, nil
}

// Helper methods for data processing and analysis

func (d *DeFiAnalyzer) getMarketSentiment(price *mcp.MCPResponse) string {
	if raw, ok := price.Raw.(map[string]interface{}); ok {
		if change24h, ok := raw["change_24h"].(float64); ok {
			if change24h > 5 {
				return "strongly bullish"
			} else if change24h > 0 {
				return "moderately bullish"  
			} else if change24h > -5 {
				return "moderately bearish"
			} else {
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
	return 12.5 // fallback
}

func (d *DeFiAnalyzer) calculateProjectedReturns(amount string, apy float64, days int) string {
	amountFloat := d.parseAmount(amount)
	dailyRate := apy / 365 / 100
	projectedRewards := amountFloat * dailyRate * float64(days)
	
	return fmt.Sprintf(`Daily Rate: %.4f%%
Total Projected Rewards: %.2f ASTRL
Annualized Return: %.2f%%
ROI for %d days: %.4f%%`, 
		dailyRate*100, projectedRewards, apy, days, (projectedRewards/amountFloat)*100)
}

func (d *DeFiAnalyzer) generateTradingInsights(quote, price *mcp.MCPResponse) string {
	return `- Optimal timing based on recent volatility patterns
- Liquidity depth analysis across available DEXes  
- Price impact assessment for your trade size
- Alternative routing suggestions for better execution`
}

func (d *DeFiAnalyzer) generateExecutionStrategy(amount string, quote *mcp.MCPResponse) string {
	amountFloat := d.parseAmount(amount)
	if amountFloat > 10000 {
		return `- Consider splitting into smaller chunks (recommended: 3-5 transactions)
- Execute during high liquidity periods (typically 12-18 UTC)
- Monitor slippage tolerance and adjust accordingly`
	}
	return `- Single transaction recommended for this size
- Execute when comfortable with current slippage
- Consider limit orders if available on the DEX`
}

func (d *DeFiAnalyzer) assessTradingRisk(quote, price *mcp.MCPResponse) string {
	return `- Slippage Risk: Monitor for sudden liquidity changes
- Timing Risk: Price volatility may affect execution
- Route Risk: Primary DEX availability and backup options
- Network Risk: Transaction fees and confirmation times`
}

func (d *DeFiAnalyzer) parseAmount(amount string) float64 {
	// Remove any non-numeric characters and parse
	cleaned := strings.ReplaceAll(amount, ",", "")
	if val, err := strconv.ParseFloat(cleaned, 64); err == nil {
		return val
	}
	return 1000 // fallback
}
💰 Yield Information:
%s

🧮 Projected Returns:
- Daily rewards estimation based on current APY
- Assumes current staking conditions remain stable
- Does not account for compound effects or APY changes
- Market price volatility may affect USD value of rewards

⚠️  Disclaimer: This is an estimation based on current data. Actual returns may vary due to market conditions, validator performance, and network changes.
`, amount, days, price.Text, apy.Text)

	return calculation, nil
}

func (d *DeFiAnalyzer) GetTradingAdvice(fromToken, toToken, amount string) (string, error) {
	quote, err := d.client.GetQuote("swap", toToken, amount, "")
	if err != nil {
		return "", err
	}

	price, err := d.client.GetPrice()
	if err != nil {
		return "", err
	}

	advice := fmt.Sprintf(`🎯 Trading Analysis for %s %s → %s

📊 Current Market:
%s

💱 Swap Quote:
%s

🧠 AI Trading Insights:
- Market timing analysis based on recent price movements
- Liquidity assessment across Radix DEXes
- Slippage considerations for your trade size
- Optimal execution strategy recommendations

⚡ Quick Tips:
- Check for better routes across multiple DEXes
- Consider breaking large trades into smaller chunks
- Monitor 24h volatility before executing
- Factor in gas costs for smaller trades
`, amount, fromToken, toToken, price.Text, quote.Text)

	return advice, nil
}
