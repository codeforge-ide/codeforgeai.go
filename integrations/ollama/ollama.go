package ollama

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/codeforge-ide/codeforgeai.go/models"
)

// Ollama API endpoint (default)
const defaultOllamaEndpoint = "http://localhost:11434/api/generate"

// OllamaModel holds model name and endpoint.
type OllamaModel struct {
	Model    string
	Endpoint string
	Timeout  time.Duration
}

// Request/Response structs for Ollama API
type ollamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	// Add other fields as needed (e.g., stream, options)
}

type ollamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
	Error    string `json:"error,omitempty"`
}

// NewOllamaModel creates a new OllamaModel with optional endpoint and timeout.
func NewOllamaModel(model string, endpoint string, timeout time.Duration) *OllamaModel {
	if endpoint == "" {
		endpoint = os.Getenv("OLLAMA_API_ENDPOINT")
		if endpoint == "" {
			endpoint = defaultOllamaEndpoint
		}
	}
	if timeout == 0 {
		timeout = 60 * time.Second
	}
	return &OllamaModel{
		Model:    model,
		Endpoint: endpoint,
		Timeout:  timeout,
	}
}

// SendRequest sends a prompt to the Ollama API and returns the response.
// config can be nil or a map with additional options.
func (o *OllamaModel) SendRequest(prompt string, config interface{}) (string, error) {
	reqBody := ollamaRequest{
		Model:  o.Model,
		Prompt: prompt,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	client := &http.Client{Timeout: o.Timeout}
	resp, err := client.Post(o.Endpoint, "application/json", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama API error: %s", string(b))
	}

	// Ollama streams responses line by line (JSON per line)
	var result string
	decoder := json.NewDecoder(resp.Body)
	for {
		var r ollamaResponse
		if err := decoder.Decode(&r); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return "", err
		}
		if r.Error != "" {
			return "", errors.New(r.Error)
		}
		result += r.Response
		if r.Done {
			break
		}
	}
	return result, nil
}

var _ models.Model = (*OllamaModel)(nil)
