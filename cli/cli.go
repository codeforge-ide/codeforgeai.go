package cli

import (
	"fmt"
	"os"

	"github.com/codeforge-ide/codeforgeai.go/cmd"
	"github.com/codeforge-ide/codeforgeai.go/config"
)

func Main() {
	// Optionally, handle config command here for direct CLI entrypoint
	if len(os.Args) > 1 && os.Args[1] == "config" {
		cfg, err := config.EnsureConfigPrompts("")
		if err != nil {
			fmt.Println("Error loading config:", err)
			os.Exit(1)
		}
		fmt.Println("Configuration checkup complete. Current configuration:")
		config.PrintConfig(cfg)
		return
	}
	cmd.Execute()
}
