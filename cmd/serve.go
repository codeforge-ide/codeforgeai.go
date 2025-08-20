package cmd

import (
	"github.com/codeforge-ide/codeforgeai.go/api"
	"github.com/spf13/cobra"
)

var port int

func init() {
	serveCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to run the API server on")
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the API server",
	Run: func(cmd *cobra.Command, args []string) {
		server := api.NewServer()
		server.Start(port)
	},
}
