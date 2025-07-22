// Package io provides integration with the IO Intelligence Agents API using Go.
// This integration uses github.com/openai/openai-go for API compatibility.

package io

import (
	"context"

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

// Add more agent workflows as needed (reasoning, sentiment, translation, etc.)
