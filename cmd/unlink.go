package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var unlinkCmd = &cobra.Command{
	Use:   "unlink <from-id> <relation> <to-id>",
	Short: "Remove a relation between two objects",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		both, _ := cmd.Flags().GetBool("both")

		vault := resolveVault(vaultPath)
		if err := vault.Open(); err != nil {
			return err
		}
		defer vault.Close()

		fromID, relName, toID := args[0], args[1], args[2]
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
