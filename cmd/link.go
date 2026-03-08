package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var linkCmd = &cobra.Command{
	Use:   "link <from-id> <relation> <to-id>",
	Short: "Link two objects with a relation",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		vault := resolveVault(vaultPath)
		if err := vault.Open(); err != nil {
			return err
		}
		defer vault.Close()

		fromID, relName, toID := args[0], args[1], args[2]
		if err := vault.LinkObjects(fromID, relName, toID); err != nil {
			return err
		}

		fmt.Printf("Linked: %s -[%s]-> %s\n", fromID, relName, toID)
		return nil
	},
}

func init() {
	relationCmd.AddCommand(linkCmd)
}
