package cmd

import (
	mcpserver "github.com/typemd/typemd/mcp"
	"github.com/spf13/cobra"
)

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Start MCP server via stdio",
	Long:  `Start a Model Context Protocol (MCP) server using stdio transport. This allows AI clients to query the vault.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		vault, err := openVault(vaultPath, reindex)
		if err != nil {
			return err
		}
		defer vault.Close()

		return mcpserver.Serve(vault)
	},
}

func init() {
	rootCmd.AddCommand(mcpCmd)
}
