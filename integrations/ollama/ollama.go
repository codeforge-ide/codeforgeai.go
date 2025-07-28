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

	"github.com/codeforge-ide/codeforgeai.go/config"
	"github.com/codeforge-ide/codeforgeai.go/modeliface"
	"github.com/spf13/cobra"
)

// Ollama API endpoint (default)
const defaultOllamaEndpoint = "http://localhost:11434/api/generate"

// OllamaModel holds model name and endpoint.
type OllamaModel struct {
	modeliface.BaseAgent
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
// NewOllamaModel creates a new OllamaModel with optional endpoint and timeout.
// If model or endpoint are empty, tries config first, then env, then default.
func NewOllamaModel(model string, endpoint string, timeout time.Duration, cfg *config.Config) *OllamaModel {
	if model == "" && cfg != nil && cfg.OllamaModel != "" {
		model = cfg.OllamaModel
	}
	if model == "" {
		model = os.Getenv("OLLAMA_MODEL")
		if model == "" && cfg != nil && cfg.CodeModel != "" {
			model = cfg.CodeModel
		}
		if model == "" {
			model = "qwen2.5-coder:1.5b"
		}
	}
	if endpoint == "" && cfg != nil && cfg.OllamaEndpoint != "" {
		endpoint = cfg.OllamaEndpoint
	}
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
		Model:     model,
		Endpoint:  endpoint,
		Timeout:   timeout,
		BaseAgent: modeliface.BaseAgent{NameStr: model},
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

var _ modeliface.Model = (*OllamaModel)(nil)

// --- Agent interface implementation ---
// Name, Login, Logout are provided by BaseAgent.
// IsAuthenticated always returns true (no auth required for Ollama)
func (o *OllamaModel) IsAuthenticated() bool {
	return true
}

// --- Registry and CLI wiring ---

func ollamaRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ollama",
		Short: "Interact with Ollama local LLMs",
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "prompt [prompt]",
		Short: "Send a prompt to Ollama",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cfg, _ := config.LoadConfig("")
			model := os.Getenv("OLLAMA_MODEL")
			if model == "" && cfg.OllamaModel != "" {
				model = cfg.OllamaModel
			}
			if model == "" {
				model = "qwen2.5-coder:1.5b"
			}
			ollama := NewOllamaModel(model, "", 0, &cfg)
			resp, err := ollama.SendRequest(args[0], nil)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			fmt.Println(resp)
		},
	})
	return cmd
}

func init() {
	modeliface.GlobalAgentRegistry.RegisterIntegration(modeliface.IntegrationRegistration{
		Metadata: modeliface.IntegrationMetadata{
			Name:        "ollama",
			Description: "Run local LLMs with Ollama for fast, private code analysis and automation.",
			ConfigKeys:  []string{"OLLAMA_MODEL", "OLLAMA_API_ENDPOINT"},
			Secrets:     []string{},
			Commands:    []string{"prompt"},
		},
		CommandFactory: ollamaRootCommand,
	})
}
