package cmd

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/typemd/typemd/web"
)

var servePort int

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start web UI server",
	Long:  `Starts an HTTP server that serves the web UI and REST API for the vault.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		vault, err := openVault(vaultPath)
		if err != nil {
			return err
		}
		defer vault.Close()

		frontend := web.FrontendFS()
		srv := web.NewServer(vault, frontend)

		addr := fmt.Sprintf(":%d", servePort)
		fmt.Printf("TypeMD web UI: http://localhost%s\n", addr)
		return http.ListenAndServe(addr, srv)
	},
}

func init() {
	serveCmd.Flags().IntVarP(&servePort, "port", "p", 3000, "port to listen on")
	rootCmd.AddCommand(serveCmd)
}
