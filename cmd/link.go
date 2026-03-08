package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var linkCmd = &cobra.Command{
	Use:   "link <from-id> <relation> <to-id>",
	Short: "Link two objects with a relation",
	Long: `Create a relation between two objects.

Supports prefix matching — you can omit the ULID suffix if the prefix
uniquely identifies an object.

Examples:
  tmd relation link book/clean-code author person/robert-martin
  tmd relation link book/clean-code-01jqr3k5mpbvn8e0f2g7h9txyz author person/robert-martin-01jqr3k8yznw2a4dbx6t7c9fpq`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		vault := resolveVault(vaultPath)
		if err := vault.Open(); err != nil {
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
