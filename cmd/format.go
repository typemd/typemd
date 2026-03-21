package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var formatDryRun bool
var formatType string

// ErrFormatNeeded is returned by dry-run when files need formatting.
// It signals exit code 1 without printing an error message.
var ErrFormatNeeded = errors.New("files need formatting")

var formatCmd = &cobra.Command{
	Use:   "format",
	Short: "Normalize frontmatter and schema formatting",
	Long: `Format all object and schema files with canonical property ordering
and YAML style. Similar to gofmt, it enforces a consistent file format.

Property ordering: system properties first (name, description, created_at,
updated_at, tags), then schema-defined properties in schema order, then
extra properties alphabetically.

Body content is preserved unchanged. The updated_at timestamp is not modified.`,
	Args:          cobra.NoArgs,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		vault, err := openVault(vaultPath, reindex)
		if err != nil {
			return err
		}
		defer vault.Close()

		result, err := vault.FormatAll(formatType, formatDryRun)
		if err != nil {
			return err
		}

		if len(result.Changed) == 0 {
			fmt.Println("All files are already formatted. No changes needed.")
			return nil
		}

		if formatDryRun {
			for _, path := range result.Changed {
				fmt.Println(path)
			}
			return ErrFormatNeeded
		}

		fmt.Printf("Formatted %d file(s).\n", len(result.Changed))
		return nil
	},
}

func init() {
	formatCmd.Flags().BoolVar(&formatDryRun, "dry-run", false, "List files that need formatting without modifying them (exit code 1 if any)")
	formatCmd.Flags().StringVar(&formatType, "type", "", "Format only objects and schemas of a specific type")
	formatCmd.RegisterFlagCompletionFunc("type", func(_ *cobra.Command, _ []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completeTypeName(toComplete)
	})
	rootCmd.AddCommand(formatCmd)
}
