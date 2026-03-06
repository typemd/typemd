package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var reindexCmd = &cobra.Command{
	Use:   "reindex",
	Short: "Sync objects from disk and rebuild search index",
	Long:  `Scan the objects directory, sync all files to the database, clean up orphaned relations, and rebuild the FTS5 search index.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		vault := resolveVault(vaultPath)
		if err := vault.Open(); err != nil {
			return err
		}
		defer vault.Close()

		result, err := vault.SyncIndex()
		if err != nil {
			return err
		}

		fmt.Println("Index synced successfully.")
		if len(result.Orphaned) > 0 {
			fmt.Printf("Warning: Found %d orphaned relation(s):\n", len(result.Orphaned))
			for _, o := range result.Orphaned {
				fmt.Printf("  %s -> %s (relation: %q)\n", o.FromID, o.ToID, o.Name)
			}
			fmt.Println("Orphaned relations have been removed from the index.")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(reindexCmd)
}
