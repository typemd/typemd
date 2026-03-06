package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var searchJSON bool

var searchCmd = &cobra.Command{
	Use:   "search <keyword>",
	Short: "Full-text search objects",
	Long: `Search objects using SQLite FTS5 full-text search across filename, properties, and body.

Examples:
  tmd search "concurrency"
  tmd search "golang channel" --json`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		vault := resolveVault(vaultPath)
		if err := vault.Open(); err != nil {
			return err
		}
		defer vault.Close()

		results, err := vault.SearchObjects(args[0])
		if err != nil {
			return err
		}

		if searchJSON {
			data, err := json.MarshalIndent(results, "", "  ")
			if err != nil {
				return fmt.Errorf("marshal JSON: %w", err)
			}
			fmt.Println(string(data))
		} else {
			for _, obj := range results {
				fmt.Println(obj.ID)
			}
		}

		return nil
	},
}

func init() {
	searchCmd.Flags().BoolVar(&searchJSON, "json", false, "Output results as JSON")
	rootCmd.AddCommand(searchCmd)
}
