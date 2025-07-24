package cmd

import (
	"fmt"

	"github.com/codeforge-ide/codeforgeai.go/integrations/opencode"
	"github.com/spf13/cobra"
)

var opencodeCmd = &cobra.Command{
	Use:   "opencode",
	Short: "Opencode AI integration commands",
}

var opencodeAnalyzeCmd = &cobra.Command{
	Use:   "analyze [target]",
	Short: "Analyze a project or file using Opencode",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		params := map[string]interface{}{"target": args[0]}
		result, err := opencode.FindFiles(params)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println(result)
	},
}

var opencodePromptCmd = &cobra.Command{
	Use:   "prompt [prompt]",
	Short: "Send a prompt to Opencode AI",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		params := map[string]interface{}{"prompt": args[0]}
		resp, err := opencode.AppendPromptTui(params)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println(resp)
	},
}

var opencodeEditCmd = &cobra.Command{
	Use:   "edit [file] [edit_prompt]",
	Short: "Edit code using Opencode AI",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		params := map[string]interface{}{"file": args[0], "edit_prompt": args[1]}
		result, err := opencode.ChatSession(args[0], params)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println(result)
	},
}

var opencodeConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Get Opencode config",
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := opencode.GetConfig()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println(resp)
	},
}

var opencodeSessionCmd = &cobra.Command{
	Use:   "session [action] [id]",
	Short: "Manage Opencode sessions",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		action := args[0]
		id := ""
		if len(args) > 1 {
			id = args[1]
		}
		var resp string
		var err error
		switch action {
		case "create":
			resp, err = opencode.CreateSession()
		case "list":
			resp, err = opencode.ListSessions()
		case "delete":
			resp, err = opencode.DeleteSession(id)
		case "abort":
			resp, err = opencode.AbortSession(id)
		case "messages":
			resp, err = opencode.GetSessionMessages(id)
		default:
			fmt.Println("Unknown session action")
			return
		}
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println(resp)
	},
}

func init() {
	opencodeCmd.AddCommand(opencodeAnalyzeCmd)
	opencodeCmd.AddCommand(opencodePromptCmd)
	opencodeCmd.AddCommand(opencodeEditCmd)
	opencodeCmd.AddCommand(opencodeConfigCmd)
	opencodeCmd.AddCommand(opencodeSessionCmd)
	rootCmd.AddCommand(opencodeCmd)
}
