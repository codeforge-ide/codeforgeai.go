package cmd

import (
	"fmt"

	"github.com/codeforge-ide/codeforgeai.go/config"
	"github.com/spf13/cobra"
)

var copilotCmd = &cobra.Command{
	Use:   "copilot",
	Short: "GitHub Copilot integration commands",
}

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication for Copilot",
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to GitHub Copilot",
	Run: func(cmd *cobra.Command, args []string) {
		var token string
		fmt.Print("Enter your GitHub Copilot token: ")
		fmt.Scanln(&token)
		if token == "" {
			fmt.Println("No token provided. Aborting login.")
			return
		}
		err := config.Set("copilot_token", token)
		if err != nil {
			fmt.Println("Failed to save token:", err)
			return
		}
		fmt.Println("Copilot login successful.")
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from GitHub Copilot",
	Run: func(cmd *cobra.Command, args []string) {
		err := config.Delete("copilot_token")
		if err != nil {
			fmt.Println("Failed to remove token:", err)
			return
		}
		fmt.Println("Copilot logout successful.")
	},
}

func init() {
	authCmd.AddCommand(loginCmd)
	authCmd.AddCommand(logoutCmd)
	copilotCmd.AddCommand(authCmd)
	rootCmd.AddCommand(copilotCmd)
}
