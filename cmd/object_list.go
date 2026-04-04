package cmd

import (
	"github.com/spf13/cobra"
	"github.com/typemd/typemd/core"
)

var objectListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all objects",
	Long: `List all objects in the vault.

Examples:
  tmd object list
  tmd object list --json
  tmd object list --include-archived`,
	RunE: func(cmd *cobra.Command, args []string) error {
		vault, err := openVault(vaultPath)
		if err != nil {
			return err
		}
		defer vault.Close()

		var opts []core.QueryOption
		includeArchived, _ := cmd.Flags().GetBool("include-archived")
		if includeArchived {
			opts = append(opts, core.QueryIncludeArchived())
		}

		results, err := vault.QueryObjects(nil, opts...)
		if err != nil {
			return err
		}

		jsonOutput, _ := cmd.Flags().GetBool("json")
		return printObjects(results, jsonOutput)
	},
}

func init() {
	objectListCmd.Flags().Bool("json", false, "Output results as JSON")
	objectListCmd.Flags().Bool("include-archived", false, "Include archived objects in results")
	objectCmd.AddCommand(objectListCmd)
}
