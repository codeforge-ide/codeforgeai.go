// Package io provides integration with the IO Intelligence Models API using Go.
// This integration uses github.com/openai/openai-go for API compatibility.

package io

import (
	"context"

	"github.com/openai/openai-go"
)

// ModelConfig holds configuration for an IO model.
type ModelConfig struct {
	Model   string
	APIKey  string
	BaseURL string
}

// IOModel wraps the OpenAI client for IO Intelligence model tasks.
type IOModel struct {
	client *openai.Client
	config ModelConfig
}

// NewIOModel creates a new IOModel instance.
func NewIOModel(cfg ModelConfig) *IOModel {
	client := openai.NewClient(cfg.APIKey)
	client.BaseURL = cfg.BaseURL
	return &IOModel{client: client, config: cfg}
}

// ChatCompletion runs a chat completion using the model.
func (m *IOModel) ChatCompletion(ctx context.Context, messages []openai.ChatCompletionMessage, maxTokens int) (string, error) {
	resp, err := m.client.CreateChatCompletion(ctx, &openai.ChatCompletionRequest{
		Model:     m.config.Model,
		Messages:  messages,
		MaxTokens: maxTokens,
	})
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}

// Add more model API methods as needed (embeddings, etc.)
