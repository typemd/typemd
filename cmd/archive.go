package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var archiveCmd = &cobra.Command{
	Use:   "archive <id>",
	Short: "Archive an object (soft delete)",
	Long: `Archive an object to hide it from default queries without deleting the file.

Supports prefix matching — you can omit the ULID suffix if the prefix
uniquely identifies an object. If a prefix matches multiple objects,
an interactive picker is shown to select the intended one.

Examples:
  tmd object archive book/clean-code
  tmd object archive book/clean-code-01jqr3k5mpbvn8e0f2g7h9txyz`,
	Args: cobra.ExactArgs(1),
	ValidArgsFunction: func(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return completeObjectID(toComplete)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		vault, err := openVault(vaultPath)
		if err != nil {
			return err
		}
		defer vault.Close()

		id, err := resolveIDInteractive(vault, args[0])
		if err != nil {
			return err
		}

		obj, err := vault.GetObject(id)
		if err != nil {
			return err
		}

		if obj.IsArchived() {
			fmt.Printf("Object %s is already archived\n", id)
			return nil
		}

		if err := vault.SetArchived(id, true); err != nil {
			return err
		}

		fmt.Printf("Archived %s\n", id)
		return nil
	},
}

var unarchiveCmd = &cobra.Command{
	Use:   "unarchive <id>",
	Short: "Unarchive an object",
	Long: `Unarchive a previously archived object so it appears in default queries again.

Supports prefix matching — you can omit the ULID suffix if the prefix
uniquely identifies an object. If a prefix matches multiple objects,
an interactive picker is shown to select the intended one.

Examples:
  tmd object unarchive book/clean-code
  tmd object unarchive book/clean-code-01jqr3k5mpbvn8e0f2g7h9txyz`,
	Args: cobra.ExactArgs(1),
	ValidArgsFunction: func(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return completeObjectID(toComplete)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		vault, err := openVault(vaultPath)
		if err != nil {
			return err
		}
		defer vault.Close()

		id, err := resolveIDInteractive(vault, args[0])
		if err != nil {
			return err
		}

		obj, err := vault.GetObject(id)
		if err != nil {
			return err
		}

		if !obj.IsArchived() {
			fmt.Printf("Object %s is not archived\n", id)
			return nil
		}

		if err := vault.SetArchived(id, false); err != nil {
			return err
		}

		fmt.Printf("Unarchived %s\n", id)
		return nil
	},
}

func init() {
	objectCmd.AddCommand(archiveCmd)
	objectCmd.AddCommand(unarchiveCmd)
}
