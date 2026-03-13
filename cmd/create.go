package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create <type> [name]",
	Short: "Create a new object from a type schema",
	Long: `Create a new object file (Markdown + YAML frontmatter) based on the type schema.
A ULID is automatically appended to the filename for uniqueness.
If the type has a name template, the name argument is optional.

Examples:
  tmd object create book clean-code
  tmd object create person robert-martin
  tmd object create journal`,
	Args: cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		vault, err := openVault(vaultPath, reindex)
		if err != nil {
			return err
		}
		defer vault.Close()

		name := ""
		if len(args) > 1 {
			name = args[1]
		}
		obj, err := vault.NewObject(args[0], name)
		if err != nil {
			return err
		}

		fmt.Printf("Created %s\n", obj.ID)
		return nil
	},
}

func init() {
	objectCmd.AddCommand(createCmd)
}
