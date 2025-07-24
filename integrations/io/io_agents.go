// Package io provides integration with the IO Intelligence Agents API using Go.
// This integration uses github.com/openai/openai-go for API compatibility.

package io

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/openai/openai-go"
)

// AgentConfig holds configuration for an IO agent.
type AgentConfig struct {
	Name         string
	Instructions string
	Model        string
	APIKey       string
	BaseURL      string
}

// AgentInfo represents metadata for an agent.
type AgentInfo struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Persona     *string `json:"persona"`
	Metadata    struct {
		ImageURL *string  `json:"image_url"`
		Tags     []string `json:"tags"`
	} `json:"metadata"`
}

// IOAgent wraps the OpenAI client for IO Intelligence agent tasks.
type IOAgent struct {
	client *openai.Client
	config AgentConfig
}

// NewIOAgent creates a new IOAgent instance.
func NewIOAgent(cfg AgentConfig) *IOAgent {
	client := openai.NewClient(cfg.APIKey)
	client.BaseURL = cfg.BaseURL
	return &IOAgent{client: client, config: cfg}
}

// SummarizeText runs a summarization workflow using the agent.
func (a *IOAgent) SummarizeText(ctx context.Context, text string, maxWords int) (string, error) {
	resp, err := a.client.CreateChatCompletion(ctx, &openai.ChatCompletionRequest{
		Model: a.config.Model,
		Messages: []openai.ChatCompletionMessage{
			{Role: "system", Content: a.config.Instructions},
			{Role: "user", Content: text},
		},
		MaxTokens: maxWords,
	})
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}

// ListAgents fetches all available agents from the IO API.
func (a *IOAgent) ListAgents(ctx context.Context) (map[string]AgentInfo, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", a.config.BaseURL+"/agents", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+a.config.APIKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result struct {
		Agents map[string]AgentInfo `json:"agents"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result.Agents, nil
}

// Workflow represents a multi-agent workflow.
type Workflow struct {
	Objective string
	Agents    []*IOAgent
}

// SummarizeText runs summarization with all agents and returns their results.
func (w *Workflow) SummarizeText(ctx context.Context, maxWords int) ([]string, error) {
	var results []string
	for _, agent := range w.Agents {
		res, err := agent.SummarizeText(ctx, w.Objective, maxWords)
		if err != nil {
			results = append(results, "error: "+err.Error())
		} else {
			results = append(results, res)
		}
	}
	return results, nil
}

// Add more agent workflows as needed (reasoning, sentiment, translation, etc.)
