package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var lockCmd = &cobra.Command{
	Use:   "lock <id>",
	Short: "Lock an object to prevent editing",
	Long: `Lock an object so it cannot be modified.

Supports prefix matching — you can omit the ULID suffix if the prefix
uniquely identifies an object. If a prefix matches multiple objects,
an interactive picker is shown to select the intended one.

Examples:
  tmd object lock book/clean-code
  tmd object lock book/clean-code-01jqr3k5mpbvn8e0f2g7h9txyz`,
	Args: cobra.ExactArgs(1),
	ValidArgsFunction: func(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return completeObjectID(toComplete)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		vault, err := openVault(vaultPath, reindex)
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

		if obj.IsLocked() {
			fmt.Printf("Object %s is already locked\n", id)
			return nil
		}

		if err := vault.SetLocked(id, true); err != nil {
			return err
		}

		fmt.Printf("Locked %s\n", id)
		return nil
	},
}

var unlockCmd = &cobra.Command{
	Use:   "unlock <id>",
	Short: "Unlock an object to allow editing",
	Long: `Unlock a previously locked object so it can be modified again.

Supports prefix matching — you can omit the ULID suffix if the prefix
uniquely identifies an object. If a prefix matches multiple objects,
an interactive picker is shown to select the intended one.

Examples:
  tmd object unlock book/clean-code
  tmd object unlock book/clean-code-01jqr3k5mpbvn8e0f2g7h9txyz`,
	Args: cobra.ExactArgs(1),
	ValidArgsFunction: func(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return completeObjectID(toComplete)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		vault, err := openVault(vaultPath, reindex)
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

		if !obj.IsLocked() {
			fmt.Printf("Object %s is not locked\n", id)
			return nil
		}

		if err := vault.SetLocked(id, false); err != nil {
			return err
		}

		fmt.Printf("Unlocked %s\n", id)
		return nil
	},
}

func init() {
	objectCmd.AddCommand(lockCmd)
	objectCmd.AddCommand(unlockCmd)
}
