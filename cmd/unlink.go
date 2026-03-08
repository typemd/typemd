package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var unlinkCmd = &cobra.Command{
	Use:   "unlink <from-id> <relation> <to-id>",
	Short: "Remove a relation between two objects",
	Long: `Remove a relation between two objects.

Supports prefix matching — you can omit the ULID suffix if the prefix
uniquely identifies an object.

Examples:
  tmd relation unlink book/clean-code author person/robert-martin
  tmd relation unlink book/clean-code author person/robert-martin --both`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		both, _ := cmd.Flags().GetBool("both")

		vault, err := openVault(vaultPath, reindex)
		if err != nil {
			return err
		}
		defer vault.Close()

		relName := args[1]
		fromID, err := vault.ResolveID(args[0])
		if err != nil {
			return err
		}
		toID, err := vault.ResolveID(args[2])
		if err != nil {
			return err
		}
		if err := vault.UnlinkObjects(fromID, relName, toID, both); err != nil {
			return err
		}

		fmt.Printf("Unlinked: %s -[%s]-> %s\n", fromID, relName, toID)
		return nil
	},
}

func init() {
	unlinkCmd.Flags().Bool("both", false, "Also remove the inverse relation on the target object")
	relationCmd.AddCommand(unlinkCmd)
}
