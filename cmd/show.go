package cmd

import (
	"fmt"
	"strings"

	"github.com/typemd/typemd/core"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show <object-id>",
	Short: "Show object detail (properties, relations, body)",
	Long: `Display an object's properties, relations, and body content.

Examples:
  tmd show book/clean-code
  tmd show person/robert-martin`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		vault := resolveVault(vaultPath)
		if err := vault.Open(); err != nil {
			return err
		}
		defer vault.Close()

		obj, err := vault.GetObject(args[0])
		if err != nil {
			return err
		}

		// Fill missing schema-defined properties
		schema, _ := vault.LoadType(obj.Type)
		if schema != nil {
			for _, p := range schema.Properties {
				if _, ok := obj.Properties[p.Name]; !ok {
					obj.Properties[p.Name] = nil
				}
			}
		}

		relations, _ := vault.ListRelations(obj.ID)

		// Title
		fmt.Println(obj.ID)
		fmt.Println()

		// Properties & Relations (merged)
		relProps := make(map[string]bool)
		if schema != nil {
			for _, p := range schema.Properties {
				if p.Type == "relation" {
					relProps[p.Name] = true
				}
			}
		}

		var reverseRels []core.Relation
		for _, r := range relations {
			if r.ToID == obj.ID {
				reverseRels = append(reverseRels, r)
			}
		}

		fmt.Println("Properties")
		fmt.Println("──────────")
		if len(obj.Properties) == 0 && len(reverseRels) == 0 {
			fmt.Println("  (none)")
		} else {
			propKeys := core.OrderedPropKeys(obj.Properties, schema)
			for _, k := range propKeys {
				v := obj.Properties[k]
				if v == nil {
					fmt.Printf("  %s: (null)\n", k)
				} else if relProps[k] {
					fmt.Printf("  %s: → %v\n", k, v)
				} else {
					fmt.Printf("  %s: %v\n", k, v)
				}
			}
			for _, r := range reverseRels {
				fmt.Printf("  %s: ← %s\n", r.Name, r.FromID)
			}
		}

		// Body
		fmt.Println()
		fmt.Println("Body")
		fmt.Println("────")
		body := strings.TrimSpace(obj.Body)
		if body == "" {
			fmt.Println("  (empty)")
		} else {
			for _, line := range strings.Split(body, "\n") {
				fmt.Printf("  %s\n", line)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
