package astrolescent

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	baseURL string
	client  *http.Client
}

type PriceResponse struct {
	Raw  map[string]interface{} `json:"raw"`
	Text string                 `json:"text"`
}

type QuoteResponse struct {
	Raw  map[string]interface{} `json:"raw"`
	Text string                 `json:"text"`
}

type APYResponse struct {
	Raw  map[string]interface{} `json:"raw"`
	Text string                 `json:"text"`
}

func NewClient() *Client {
	return &Client{
		baseURL: "https://mcp.astrolescent.com",
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) GetPrice() (*PriceResponse, error) {
	// For demo purposes, simulate the response
	return &PriceResponse{
		Raw: map[string]interface{}{
			"price_xrd":  0.083,
			"price_usd":  0.083,
			"change_24h": 2.5,
			"change_7d":  -1.2,
			"volume_24h": 125000,
			"market_cap": 8300000,
		},
		Text: "🚀 ASTRL Price: 0.083 XRD ($0.083 USD) | 24h: +2.5% ⬆️ | 7d: -1.2% ⬇️ | Vol: 125K XRD",
	}, nil
}

func (c *Client) GetQuote(operation, token, amount, account string) (*QuoteResponse, error) {
	// For demo purposes, simulate the response
	var text string
	switch operation {
	case "buy":
		text = fmt.Sprintf("💸 Buy Quote: %s XRD → 1,204 ASTRL | Price Impact: 0.1%% | Route: DefiPlaza | Gas: ~0.5 XRD", amount)
	case "sell":
		text = fmt.Sprintf("💰 Sell Quote: %s ASTRL → 998.2 XRD | Price Impact: 0.2%% | Route: DefiPlaza | Gas: ~0.5 XRD", amount)
	case "swap":
		text = fmt.Sprintf("🔄 Swap Quote: %s input → 1,200.5 output | Slippage: 0.1%% | Route: DefiPlaza", amount)
	}

	return &QuoteResponse{
		Raw: map[string]interface{}{
			"operation":    operation,
			"amount_in":    amount,
			"amount_out":   1200.5,
			"slippage":     0.1,
			"price_impact": 0.1,
			"route":        "DefiPlaza",
		},
		Text: text,
	}, nil
}

func (c *Client) GetAPY() (*APYResponse, error) {
	// For demo purposes, simulate the response
	return &APYResponse{
		Raw: map[string]interface{}{
			"staking_apy":       12.5,
			"lp_apy_defiplaza":  15.8,
			"validator_rewards": 3.2,
			"lp_rewards":        2.1,
			"total_staked":      "2.5M ASTRL",
			"total_lp":          "850K ASTRL",
		},
		Text: "💰 ASTRL Yields: Staking 12.5% APY | DefiPlaza LP 15.8% APY | Validator Rewards 3.2% | Total Staked: 2.5M ASTRL",
	}, nil
}

// Legacy methods for backwards compatibility with root.go commands
func (d *DeFiAnalyzer) AnalyzeStakingVsLP() (string, error) {
	// Wrapper for the context version
	return d.AnalyzeStakingVsLP(context.Background())
}

func (d *DeFiAnalyzer) CalculateStakingReturns(amount string, days int) (string, error) {
	// Wrapper for the context version
	return d.CalculateStakingReturns(context.Background(), amount, days)
}

func (d *DeFiAnalyzer) GetTradingAdvice(fromToken, toToken, amount string) (string, error) {
	// Wrapper for the context version
	return d.GetTradingAdvice(context.Background(), fromToken, toToken, amount)
}
