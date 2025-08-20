package cmd

import (
	"fmt"
	"os/exec"

	"github.com/codeforge-ide/codeforgeai.go/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(doctorCmd)
}

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check the configuration and environment for issues",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running doctor...")

		// Check for required tools
		checkCommand("git")
		checkCommand("npx")

		// Check for config file
		cfg, err := config.LoadConfig("")
		if err != nil {
			fmt.Println("❌ Config file not found or invalid.")
		} else {
			fmt.Println("✅ Config file found.")
			config.PrintConfig(cfg)
		}

		fmt.Println("\nDoctor check complete.")
	},
}

func checkCommand(name string) {
	_, err := exec.LookPath(name)
	if err != nil {
		fmt.Printf("❌ %s not found in PATH.\n", name)
	} else {
		fmt.Printf("✅ %s found.\n", name)
	}
}
