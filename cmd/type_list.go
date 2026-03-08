package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var typeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available types",
	Long: `List all available type schemas (both custom and built-in defaults).

Examples:
  tmd type list`,
	RunE: func(cmd *cobra.Command, args []string) error {
		vault := resolveVault(vaultPath)

		types := vault.ListTypes()
		for _, name := range types {
			fmt.Println(name)
		}

		return nil
	},
}

func init() {
	typeCmd.AddCommand(typeListCmd)
}
