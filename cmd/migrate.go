package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/typemd/typemd/core"
)

var migrateDryRun bool
var migrateRenames []string

var migrateCmd = &cobra.Command{
	Use:   "migrate [type]",
	Short: "Update schemas and objects to match current conventions",
	Long: `Migrate type schemas and objects.

Without arguments, migrates all type schemas (e.g. enum → select).
With a type argument, migrates objects of that type to match the current schema.

Examples:
  tmd migrate                          # migrate schemas (enum → select)
  tmd migrate --dry-run                # preview schema migrations
  tmd migrate book                     # migrate book objects to match schema
  tmd migrate book --dry-run
  tmd migrate book --rename old_field:new_field`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		vault, err := openVault(vaultPath, reindex)
		if err != nil {
			return err
		}
		defer vault.Close()

		// No type argument: migrate schemas
		if len(args) == 0 {
			result, err := vault.MigrateSchemas(migrateDryRun)
			if err != nil {
				return err
			}

			if len(result.Changes) == 0 {
				fmt.Println("All schemas are up to date. No changes needed.")
				return nil
			}

			if migrateDryRun {
				fmt.Println("Dry run — no files modified:")
			}

			for _, c := range result.Changes {
				fmt.Printf("  %s: converted enum → select for %v\n", c.TypeName, c.Properties)
			}
			fmt.Printf("Schema migration complete: %d type(s) updated.\n", len(result.Changes))
			return nil
		}

		// Type argument: migrate objects
		opts := core.MigrateOptions{
			DryRun:  migrateDryRun,
			Renames: make(map[string]string),
		}

		for _, r := range migrateRenames {
			pair := strings.SplitN(r, ":", 2)
			if len(pair) != 2 || pair[0] == "" || pair[1] == "" {
				return fmt.Errorf("invalid rename format %q, expected old_name:new_name", r)
			}
			opts.Renames[pair[0]] = pair[1]
		}

		result, err := vault.MigrateObjects(args[0], opts)
		if err != nil {
			return err
		}

		if len(result.Changes) == 0 {
			fmt.Println("All objects already match the schema. No changes needed.")
			return nil
		}

		if migrateDryRun {
			fmt.Println("Dry run — no files modified:")
		}

		for _, c := range result.Changes {
			var summary []string
			if len(c.Added) > 0 {
				summary = append(summary, fmt.Sprintf("added %v", c.Added))
			}
			if len(c.Removed) > 0 {
				summary = append(summary, fmt.Sprintf("removed %v", c.Removed))
			}
			for old, new := range c.Renamed {
				summary = append(summary, fmt.Sprintf("renamed %s -> %s", old, new))
			}
			fmt.Printf("  %s: %s\n", c.ObjectID, strings.Join(summary, ", "))
		}

		fmt.Printf("Migration complete: %d object(s) updated.\n", len(result.Changes))
		return nil
	},
}

func init() {
	migrateCmd.Flags().BoolVar(&migrateDryRun, "dry-run", false, "Preview changes without modifying files")
	migrateCmd.Flags().StringArrayVar(&migrateRenames, "rename", nil, "Rename a property (format: old_name:new_name, repeatable)")
	rootCmd.AddCommand(migrateCmd)
}
