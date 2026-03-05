package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/typemd/typemd/core"
)

var createCmd = &cobra.Command{
	Use:   "create <type> <name>",
	Short: "Create a new object from a type schema",
	Long: `Create a new object file (Markdown + YAML frontmatter) based on the type schema.

Examples:
  tmd create book clean-code
  tmd create person robert-martin`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		vault := core.NewVault(vaultPath)
		if vaultPath == "" {
			vault = core.NewVault(".")
		}
		if err := vault.Open(); err != nil {
			return err
		}
		defer vault.Close()

		obj, err := vault.NewObject(args[0], args[1])
		if err != nil {
			return err
		}

		fmt.Printf("Created %s\n", obj.ID)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
