package githubmodels

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	defaultEndpoint = "https://models.inference.ai.azure.com/chat/completions"
	defaultModel    = ""
)

// Message represents a chat message for the API.
type Message struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"`
}

// ChatRequest is the payload for the chat/completions endpoint.
type ChatRequest struct {
	Messages []Message `json:"messages"`
	Model    string    `json:"model"`
	Stream   bool      `json:"stream,omitempty"`
}

// ChatResponse is a minimal response struct for non-streaming.
type ChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// Client for GitHub Models API.
type Client struct {
	Endpoint string
	Model    string
	Token    string
	Timeout  time.Duration
}

// NewClient creates a new GitHub Models client.
// Accepts model name and endpoint as arguments, falling back to env or default if empty.
func NewClient(token string, model string, endpoint string) *Client {
	if endpoint == "" {
		endpoint = os.Getenv("GITHUB_MODELS_ENDPOINT")
		if endpoint == "" {
			endpoint = defaultEndpoint
		}
	}
	if model == "" {
		model = os.Getenv("GITHUB_MODELS_MODEL")
	}
	timeout := 60 * time.Second
	return &Client{
		Endpoint: endpoint,
		Model:    model,
		Token:    token,
		Timeout:  timeout,
	}
}

// Chat sends a chat completion request using Go's http client.
func (c *Client) Chat(messages []Message, stream bool) (string, error) {
	reqBody := ChatRequest{
		Messages: messages,
		Model:    c.Model,
		Stream:   stream,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", c.Endpoint, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.Token)

	client := &http.Client{Timeout: c.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("github models API error: %s", string(b))
	}

	if stream {
		// For streaming, just return the raw output for now
		var sb strings.Builder
		buf := make([]byte, 4096)
		for {
			n, err := resp.Body.Read(buf)
			if n > 0 {
				sb.Write(buf[:n])
			}
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				return "", err
			}
		}
		return sb.String(), nil
	}

	// Non-streaming: parse JSON
	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", err
	}
	if len(chatResp.Choices) == 0 {
		return "", errors.New("no choices in response")
	}
	return chatResp.Choices[0].Message.Content, nil
}

// ChatWithWget sends a chat completion request using wget as a fallback.
func (c *Client) ChatWithWget(messages []Message, stream bool) (string, error) {
	reqBody := ChatRequest{
		Messages: messages,
		Model:    c.Model,
		Stream:   stream,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	tmpFile, err := os.CreateTemp("", "githubmodels_payload_*.json")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.Write(body); err != nil {
		return "", err
	}
	tmpFile.Close()

	args := []string{
		"--header=Content-Type: application/json",
		"--header=Authorization: Bearer " + c.Token,
		"--post-file=" + tmpFile.Name(),
		"-q", "-O", "-", c.Endpoint,
	}
	cmd := exec.Command("wget", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("wget error: %v\n%s", err, string(out))
	}
	return string(out), nil
}

// ChatAuto tries Go http client, falls back to wget if needed.
func (c *Client) ChatAuto(messages []Message, stream bool) (string, error) {
	resp, err := c.Chat(messages, stream)
	if err == nil {
		return resp, nil
	}
	// fallback to wget
	return c.ChatWithWget(messages, stream)
}

// Helper for simple prompt.
func (c *Client) SimplePrompt(prompt string) (string, error) {
	msgs := []Message{
		{Role: "system", Content: "You are a helpful assistant."},
		{Role: "user", Content: prompt},
	}
	return c.ChatAuto(msgs, false)
}

// Helper for multi-turn.
func (c *Client) MultiTurn(history []Message) (string, error) {
	return c.ChatAuto(history, false)
}

// Helper for streaming.
func (c *Client) StreamPrompt(prompt string) (string, error) {
	msgs := []Message{
		{Role: "system", Content: "You are a helpful assistant."},
		{Role: "user", Content: prompt},
	}
	return c.ChatAuto(msgs, true)
}

// Helper for image prompt.
func (c *Client) ImagePrompt(prompt string, imageB64 string) (string, error) {
	msgs := []Message{
		{Role: "system", Content: "You are a helpful assistant that describes images in details."},
		{Role: "user", Content: []interface{}{
			map[string]interface{}{"text": prompt, "type": "text"},
			map[string]interface{}{"image_url": map[string]interface{}{
				"url":    "data:image/png;base64," + imageB64,
				"detail": "low",
			}, "type": "image_url"},
		}},
	}
	return c.ChatAuto(msgs, false)
}

// Ensure *Client implements models.Model interface
func (c *Client) SendRequest(prompt string, config interface{}) (string, error) {
	// Use SimplePrompt for general, or allow config to select other helpers
	return c.SimplePrompt(prompt)
}
