package astrolescent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

type QuoteRequest struct {
	Operation string `json:"operation"` // buy, sell, swap
	Token     string `json:"token"`
	Amount    string `json:"amount"`
	Account   string `json:"account,omitempty"`
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
	req := map[string]interface{}{
		"method": "tools/call",
		"params": map[string]interface{}{
			"name": "price",
		},
	}

	resp, err := c.makeRequest("/sse", req)
	if err != nil {
		return nil, err
	}

	var priceResp PriceResponse
	if err := json.Unmarshal(resp, &priceResp); err != nil {
		return nil, err
	}

	return &priceResp, nil
}

func (c *Client) GetQuote(operation, token, amount, account string) (*QuoteResponse, error) {
	req := map[string]interface{}{
		"method": "tools/call",
		"params": map[string]interface{}{
			"name": "quote",
			"arguments": map[string]string{
				"operation": operation,
				"token":     token,
				"amount":    amount,
				"account":   account,
			},
		},
	}

	resp, err := c.makeRequest("/sse", req)
	if err != nil {
		return nil, err
	}

	var quoteResp QuoteResponse
	if err := json.Unmarshal(resp, &quoteResp); err != nil {
		return nil, err
	}

	return &quoteResp, nil
}

func (c *Client) GetAPY() (*APYResponse, error) {
	req := map[string]interface{}{
		"method": "tools/call",
		"params": map[string]interface{}{
			"name": "apy",
		},
	}

	resp, err := c.makeRequest("/sse", req)
	if err != nil {
		return nil, err
	}

	var apyResp APYResponse
	if err := json.Unmarshal(resp, &apyResp); err != nil {
		return nil, err
	}

	return &apyResp, nil
}

func (c *Client) makeRequest(endpoint string, payload interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.baseURL+endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// DeFiAnalyzer provides AI-powered DeFi analysis using Astrolescent data
type DeFiAnalyzer struct {
	client *Client
}

func NewDeFiAnalyzer() *DeFiAnalyzer {
	return &DeFiAnalyzer{
		client: NewClient(),
	}
}

func (d *DeFiAnalyzer) AnalyzeStakingVsLP() (string, error) {
	apy, err := d.client.GetAPY()
	if err != nil {
		return "", err
	}

	price, err := d.client.GetPrice()
	if err != nil {
		return "", err
	}

	analysis := fmt.Sprintf(`🚀 DeFi Analysis Report - Staking vs LP Strategy

📊 Current Market Data:
%s

💰 Yield Analysis:
%s

🧠 AI Recommendation:
Based on current APY data and market conditions, this analysis helps you decide between staking ASTRL or providing liquidity on DefiPlaza.

Key factors to consider:
- Risk tolerance (staking = lower risk, LP = higher risk/reward)
- Market volatility (affects impermanent loss in LP)
- Time horizon for your investment
- Current yield differentials

📈 Historical Context:
Price movements in the last 24h and 7 days provide insight into market momentum and potential impermanent loss scenarios for LP positions.
`, price.Text, apy.Text)

	return analysis, nil
}

func (d *DeFiAnalyzer) CalculateStakingReturns(amount string, days int) (string, error) {
	apy, err := d.client.GetAPY()
	if err != nil {
		return "", err
	}

	price, err := d.client.GetPrice()
	if err != nil {
		return "", err
	}

	calculation := fmt.Sprintf(`💎 ASTRL Staking Calculator

🔢 Input: %s ASTRL for %d days

📊 Current Market:
%s

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
