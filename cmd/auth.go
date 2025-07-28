package cmd

import (
	"fmt"
	"github.com/codeforge-ide/codeforgeai.go/config"
	"github.com/codeforge-ide/codeforgeai.go/integrations/astrolescent"
	"github.com/codeforge-ide/codeforgeai.go/integrations/githubcopilot"
	"github.com/codeforge-ide/codeforgeai.go/integrations/githubmodels"
	"github.com/codeforge-ide/codeforgeai.go/integrations/io"
	"github.com/codeforge-ide/codeforgeai.go/integrations/ollama"
	"github.com/codeforge-ide/codeforgeai.go/integrations/openai"
	"github.com/spf13/cobra"
	"strings"
	"time"
)

// authCmd is the root for all authentication-related commands
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Unified authentication for all providers (Copilot, Ollama, OpenAI, IO, etc.)",
}

// loginCmd handles login for a specified provider
var loginCmd = &cobra.Command{
	Use:   "login [provider]",
	Short: "Login to a provider (copilot, ollama, openai, io, etc.)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		provider := args[0]
		if provider == "copilot" {
			copilotLogin()
			return
		}
		fmt.Printf("[stub] Logging in to provider: %s\n", provider)
		// TODO: Call provider-specific login logic via agent manager
	},
}

// Copilot device flow authentication
func copilotLogin() {
	resp, err := githubcopilot.StartDeviceAuth()
	if err != nil {
		fmt.Printf("Error starting device auth: %v\n", err)
		return
	}
	fmt.Printf("\nTo authenticate with GitHub Copilot, visit: %s\nAnd enter the code: %s\n\n", resp.VerificationURI, resp.UserCode)
	fmt.Println("Waiting for authorization...")
	interval := resp.Interval
	expiry := time.Now().Add(time.Duration(resp.ExpiresIn) * time.Second)
	for time.Now().Before(expiry) {
		time.Sleep(time.Duration(interval) * time.Second)
		tokenResp, err := githubcopilot.PollForAccessToken(resp.DeviceCode)
		if err != nil {
			fmt.Printf("Error polling for access token: %v\n", err)
			return
		}
		if tokenResp.AccessToken != "" {
			copilotToken, err := githubcopilot.GetCopilotToken(tokenResp.AccessToken)
			if err != nil {
				fmt.Printf("Error getting Copilot token: %v\n", err)
				return
			}
			err = config.Set("copilot_token", copilotToken.Token)
			if err != nil {
				fmt.Printf("Error saving Copilot token: %v\n", err)
				return
			}
			fmt.Println("\nGitHub Copilot authentication successful! Token saved.")
			return
		}
		if tokenResp.Error != "authorization_pending" && tokenResp.Error != "" {
			fmt.Printf("Authentication failed: %s\n", tokenResp.ErrorDescription)
			return
		}
	}
	fmt.Println("Device code expired. Please try again.")
}

// statusCmd shows authentication status for all providers
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show authentication status for all providers",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%-18s %-20s\n", "Provider", "Status")
		fmt.Println(strings.Repeat("-", 38))

		statuses := []struct {
			Name   string
			Status string
		}{
			{"Copilot", func() string {
				c := &githubcopilot.CopilotModel{}
				if c.IsAuthenticated() {
					return "✅ Authenticated"
				} else {
					return "❌ Not Authenticated"
				}
			}()},
			{"GitHub Models", func() string {
				if githubmodels.NewClient("", "", "").IsAuthenticated() {
					return "✅ Authenticated"
				} else {
					return "❌ Not Authenticated"
				}
			}()},
			{"IO.net", func() string {
				if io.New("").IsAuthenticated() {
					return "✅ Authenticated"
				} else {
					return "❌ Not Authenticated"
				}
			}()},
			{"Ollama", func() string {
				if ollama.NewOllamaModel("", "", 0, nil).IsAuthenticated() {
					return "✅ Authenticated"
				} else {
					return "❌ Not Authenticated"
				}
			}()},
			{"OpenAI", func() string {
				if openai.NewOpenAIModel("").IsAuthenticated() {
					return "✅ Authenticated"
				} else {
					return "❌ Not Authenticated"
				}
			}()},
			{"Astrolescent", func() string {
				if astrolescent.NewDeFiAnalyzer().IsAuthenticated() {
					return "✅ Authenticated"
				} else {
					return "❌ Not Authenticated"
				}
			}()},
		}

		for _, s := range statuses {
			fmt.Printf("%-18s %-20s\n", s.Name, s.Status)
		}
	}, // fixed syntax error
}

// switchCmd sets the active provider
var switchCmd = &cobra.Command{
	Use:   "switch [provider]",
	Short: "Switch the active provider for all CLI features",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		provider := args[0]
		fmt.Printf("[stub] Switching active provider to: %s\n", provider)
		// TODO: Set the active provider in config/agent manager
	},
}

func init() {
	authCmd.AddCommand(loginCmd)
	authCmd.AddCommand(statusCmd)
	authCmd.AddCommand(switchCmd)
	rootCmd.AddCommand(authCmd)
}
