package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var queryJSON bool

var queryCmd = &cobra.Command{
	Use:   "query <filter>",
	Short: "Query objects by key=value filters",
	Long: `Query objects using key=value filter syntax. Multiple conditions are combined with AND.

Examples:
  tmd query "type=book"
  tmd query "type=book status=reading"
  tmd query "type=book status=reading" --json`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		vault := resolveVault(vaultPath)
		if err := vault.Open(); err != nil {
			return err
		}
		defer vault.Close()

		results, err := vault.QueryObjects(args[0])
		if err != nil {
			return err
		}

		if queryJSON {
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
	queryCmd.Flags().BoolVar(&queryJSON, "json", false, "Output results as JSON")
	rootCmd.AddCommand(queryCmd)
}
