package cmd

import (
	"fmt"
	"sort"

	"github.com/MilesChou/typemd/core"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate type schemas, objects, and relations",
	RunE: func(cmd *cobra.Command, args []string) error {
		vault := core.NewVault(vaultPath)
		if vaultPath == "" {
			vault = core.NewVault(".")
		}
		if err := vault.Open(); err != nil {
			return err
		}
		defer vault.Close()

		totalErrors := 0

		// Phase 1: Schema validation
		schemaErrs := core.ValidateAllSchemas(vault)
		if len(schemaErrs) > 0 {
			fmt.Println("Schema errors:")
			schemaNames := make([]string, 0, len(schemaErrs))
			for name := range schemaErrs {
				schemaNames = append(schemaNames, name)
			}
			sort.Strings(schemaNames)
			for _, name := range schemaNames {
				for _, e := range schemaErrs[name] {
					fmt.Printf("  %s.yaml: %s\n", name, e)
					totalErrors++
				}
			}
			fmt.Println()
		}

		// Phase 2: Object validation
		objectErrs := core.ValidateAllObjects(vault)
		if len(objectErrs) > 0 {
			fmt.Println("Object errors:")
			objectIDs := make([]string, 0, len(objectErrs))
			for id := range objectErrs {
				objectIDs = append(objectIDs, id)
			}
			sort.Strings(objectIDs)
			for _, id := range objectIDs {
				for _, e := range objectErrs[id] {
					fmt.Printf("  %s: %s\n", id, e)
					totalErrors++
				}
			}
			fmt.Println()
		}

		// Phase 3: Relation validation
		relationErrs := core.ValidateRelations(vault)
		if len(relationErrs) > 0 {
			fmt.Println("Relation errors:")
			for _, e := range relationErrs {
				fmt.Printf("  %s\n", e)
				totalErrors++
			}
			fmt.Println()
		}

		if totalErrors > 0 {
			return fmt.Errorf("found %d validation error(s)", totalErrors)
		}

		fmt.Println("Validation passed.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}
