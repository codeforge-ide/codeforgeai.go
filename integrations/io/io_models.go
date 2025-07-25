package io

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/codeforge-ide/codeforgeai.go/modeliface"
)

// Model represents an IO.net model.
type Model struct {
	ID         string `json:"id"`
	Object     string `json:"object"`
	OwnedBy    string `json:"owned_by"`
	Permission []any  `json:"permission"`
}

// Ensure Model implements modeliface.Model
func (m *Model) SendRequest(prompt string, config interface{}) (string, error) {
	// Placeholder implementation
	return "IOModel response", nil
}

// ListModels returns a list of available models.
func (i *IO) ListModels(ctx context.Context) ([]modeliface.Model, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/v1/models", i.client.BaseURL), nil)
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
		return nil, fmt.Errorf("failed to get models: %s", resp.Status)
	}

	var models struct {
		Data []Model `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&models); err != nil {
		return nil, err
	}

	var result []modeliface.Model
	for _, m := range models.Data {
		result = append(result, &Model{
			ID:      m.ID,
			OwnedBy: m.OwnedBy,
		})
	}

	return result, nil
}
