package io

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Agent represents an IO.net agent.
type Agent struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ListAgents returns a list of available agents.
func (i *IO) ListAgents(ctx context.Context) ([]Agent, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/v1/agents", i.client.BaseURL), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", i.client.APIKey))

	resp, err := i.client.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get agents: %s", resp.Status)
	}

	var agents []Agent
	if err := json.NewDecoder(resp.Body).Decode(&agents); err != nil {
		return nil, err
	}

	return agents, nil
}

// AgentRequest represents a request to an IO.net agent.
type AgentRequest struct {
	Model   string `json:"model"`
	Content string `json:"content"`
}

// AgentResponse represents a response from an IO.net agent.
type AgentResponse struct {
	Content string `json:"content"`
}

// RunAgent runs an IO.net agent.
func (i *IO) RunAgent(ctx context.Context, agentID string, req *AgentRequest) (*AgentResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/v1/agents/%s/run", i.client.BaseURL, agentID), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", i.client.APIKey))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := i.client.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to run agent: %s", resp.Status)
	}

	var agentResp AgentResponse
	if err := json.NewDecoder(resp.Body).Decode(&agentResp); err != nil {
		return nil, err
	}

	return &agentResp, nil
}