package cmd

import (
	"github.com/spf13/cobra"
)

var objectListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all objects",
	Long: `List all objects in the vault.

Examples:
  tmd object list
  tmd object list --json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		vault := resolveVault(vaultPath)
		if err := vault.Open(); err != nil {
			return err
		}
		defer vault.Close()

		results, err := vault.QueryObjects("")
		if err != nil {
			return err
		}

		jsonOutput, _ := cmd.Flags().GetBool("json")
		return printObjects(results, jsonOutput)
	},
}

func init() {
	objectListCmd.Flags().Bool("json", false, "Output results as JSON")
	objectCmd.AddCommand(objectListCmd)
}
