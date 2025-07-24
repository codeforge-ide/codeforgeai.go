// Package io provides integration with the IO Intelligence Models API using Go.
// This integration uses github.com/openai/openai-go for API compatibility.

package io

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

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

// EmbeddingRequest and EmbeddingResponse for IO API
type EmbeddingRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}
type EmbeddingResponse struct {
	Data []struct {
		Embedding []float64 `json:"embedding"`
		Index     int       `json:"index"`
	} `json:"data"`
}

// GetEmbeddings fetches embeddings for input texts.
func (m *IOModel) GetEmbeddings(ctx context.Context, input []string) ([][]float64, error) {
	reqBody, _ := json.Marshal(EmbeddingRequest{Model: m.config.Model, Input: input})
	req, err := http.NewRequestWithContext(ctx, "POST", m.config.BaseURL+"/embeddings", bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+m.config.APIKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result EmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	var embeddings [][]float64
	for _, d := range result.Data {
		embeddings = append(embeddings, d.Embedding)
	}
	return embeddings, nil
}

// Add more model API methods as needed (embeddings, etc.)
