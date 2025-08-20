package githubmodels

import (
	"fmt"
	"os"
	"strings"

	"github.com/codeforge-ide/codeforgeai.go/config"
	"github.com/codeforge-ide/codeforgeai.go/modeliface"
	"github.com/codeforge-ide/codeforgeai.go/secrets"
	"github.com/spf13/cobra"
)

func init() {
	modeliface.GlobalAgentRegistry.RegisterIntegration(modeliface.IntegrationRegistration{
		Metadata: modeliface.IntegrationMetadata{
			Name:         "github-models",
			Description:  "Interact with GitHub Models API.",
			Commands:     []string{"prompt", "multi-turn", "stream", "image", "token-store", "token-load"},
			Capabilities: []string{"text_completion"},
		},
		CommandFactory: NewGithubModelsCommand,
	})
}

func NewGithubModelsCommand() *cobra.Command {
	githubModelsCmd := &cobra.Command{
		Use:   "github-models",
		Short: "Interact with GitHub Models API",
	}

	// github-models prompt
	githubModelsPromptCmd := &cobra.Command{
		Use:   "prompt [prompt]",
		Short: "Send a simple prompt to GitHub Models",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cfg, _ := config.LoadConfig("")
			token := cfg.GithubToken
			if token == "" {
				token = os.Getenv("GITHUB_TOKEN")
			}
			if token == "" {
				fmt.Println("GITHUB_TOKEN environment variable or config is required.")
				return
			}
			model := cfg.GithubModelsModel
			if model == "" {
				model = os.Getenv("GITHUB_MODELS_MODEL")
			}
			endpoint := cfg.GithubModelsEndpoint
			if endpoint == "" {
				endpoint = os.Getenv("GITHUB_MODELS_ENDPOINT")
			}
			client := NewClient(token, model, endpoint)
			resp, err := client.SimplePrompt(strings.Join(args, " "))
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			fmt.Println(resp)
		},
	}
	githubModelsCmd.AddCommand(githubModelsPromptCmd)

	// github-models multi-turn
	githubModelsMultiTurnCmd := &cobra.Command{
		Use:   "multi-turn",
		Short: "Send a multi-turn conversation to GitHub Models",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, _ := config.LoadConfig("")
			token := cfg.GithubToken
			if token == "" {
				token = os.Getenv("GITHUB_TOKEN")
			}
			if token == "" {
				fmt.Println("GITHUB_TOKEN environment variable or config is required.")
				return
			}
			model := cfg.GithubModelsModel
			if model == "" {
				model = os.Getenv("GITHUB_MODELS_MODEL")
			}
			endpoint := cfg.GithubModelsEndpoint
			if endpoint == "" {
				endpoint = os.Getenv("GITHUB_MODELS_ENDPOINT")
			}
			// For demo, hardcode a conversation; in real use, parse from args or file
			history := []Message{
				{Role: "system", Content: "You are a helpful assistant."},
				{Role: "user", Content: "What is the capital of France?"},
				{Role: "assistant", Content: "The capital of France is Paris."},
				{Role: "user", Content: "What about Spain?"},
			}
			client := NewClient(token, model, endpoint)
			resp, err := client.MultiTurn(history)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			fmt.Println(resp)
		},
	}
	githubModelsCmd.AddCommand(githubModelsMultiTurnCmd)

	// github-models stream
	githubModelsStreamCmd := &cobra.Command{
		Use:   "stream [prompt]",
		Short: "Stream a prompt response from GitHub Models",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cfg, _ := config.LoadConfig("")
			token := cfg.GithubToken
			if token == "" {
				token = os.Getenv("GITHUB_TOKEN")
			}
			if token == "" {
				fmt.Println("GITHUB_TOKEN environment variable or config is required.")
				return
			}
			model := cfg.GithubModelsModel
			if model == "" {
				model = os.Getenv("GITHUB_MODELS_MODEL")
			}
			endpoint := cfg.GithubModelsEndpoint
			if endpoint == "" {
				endpoint = os.Getenv("GITHUB_MODELS_ENDPOINT")
			}
			client := NewClient(token, model, endpoint)
			resp, err := client.StreamPrompt(strings.Join(args, " "))
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			fmt.Println(resp)
		},
	}
	githubModelsCmd.AddCommand(githubModelsStreamCmd)

	// github-models image
	githubModelsImageCmd := &cobra.Command{
		Use:   "image [prompt] [image_path]",
		Short: "Send an image prompt to GitHub Models",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			cfg, _ := config.LoadConfig("")
			token := cfg.GithubToken
			if token == "" {
				token = os.Getenv("GITHUB_TOKEN")
			}
			if token == "" {
				fmt.Println("GITHUB_TOKEN environment variable or config is required.")
				return
			}
			model := cfg.GithubModelsModel
			if model == "" {
				model = os.Getenv("GITHUB_MODELS_MODEL")
			}
			endpoint := cfg.GithubModelsEndpoint
			if endpoint == "" {
				endpoint = os.Getenv("GITHUB_MODELS_ENDPOINT")
			}
			imagePath := args[1]
			imageData, err := os.ReadFile(imagePath)
			if err != nil {
				fmt.Println("Error reading image:", err)
				return
			}
			imageB64 := encodeToBase64(imageData)
			client := NewClient(token, model, endpoint)
			resp, err := client.ImagePrompt(args[0], imageB64)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			fmt.Println(resp)
		},
	}
	githubModelsCmd.AddCommand(githubModelsImageCmd)

	// github-models token-store
	githubModelsTokenStoreCmd := &cobra.Command{
		Use:   "token-store",
		Short: "Securely store your GitHub Models API token",
		Run: func(cmd *cobra.Command, args []string) {
			err := secrets.InteractiveStoreGithubToken()
			if err != nil {
				fmt.Println("Error storing token:", err)
			} else {
				fmt.Println("GitHub token stored securely.")
			}
		},
	}
	githubModelsCmd.AddCommand(githubModelsTokenStoreCmd)

	// github-models token-load
	githubModelsTokenLoadCmd := &cobra.Command{
		Use:   "token-load",
		Short: "Load your GitHub Models API token into the environment for this session",
		Run: func(cmd *cobra.Command, args []string) {
			token, err := secrets.InteractiveLoadGithubToken()
			if err != nil {
				fmt.Println("Error loading token:", err)
			} else {
				fmt.Println("GitHub token loaded into environment for this session.")
				// Optionally print token for debug (not recommended for real use)
				_ = token
			}
		},
	}
	githubModelsCmd.AddCommand(githubModelsTokenLoadCmd)

	return githubModelsCmd
}

// Helper function for base64 encoding
func encodeToBase64(data []byte) string {
	return strings.TrimRight(strings.ReplaceAll(fmt.Sprintf("%+q", data), "\\x", ""), "\"")
}
