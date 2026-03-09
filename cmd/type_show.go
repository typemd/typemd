package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var typeShowCmd = &cobra.Command{
	Use:   "show <type>",
	Short: "Show type schema details",
	Long: `Display the schema definition for a type, including all properties and their configurations.

Examples:
  tmd type show book
  tmd type show person`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		vault := resolveVault(vaultPath)

		schema, err := vault.LoadType(args[0])
		if err != nil {
			return err
		}

		if schema.Emoji != "" {
			fmt.Printf("Type: %s %s\n", schema.Emoji, schema.Name)
		} else {
			fmt.Printf("Type: %s\n", schema.Name)
		}
		fmt.Println()
		fmt.Println("Properties")
		fmt.Println("──────────")

		if len(schema.Properties) == 0 {
			fmt.Println("  (none)")
			return nil
		}

		for _, p := range schema.Properties {
			line := fmt.Sprintf("  %s (%s)", p.Name, p.Type)
			if (p.Type == "select" || p.Type == "multi_select") && len(p.Options) > 0 {
				line += fmt.Sprintf(" [%s]", strings.Join(p.OptionValues(), ", "))
			}
			if p.Type == "relation" && p.Target != "" {
				line += fmt.Sprintf(" -> %s", p.Target)
				if p.Multiple {
					line += " (multiple)"
				}
				if p.Bidirectional {
					line += " (bidirectional)"
					if p.Inverse != "" {
						line += fmt.Sprintf(" inverse=%s", p.Inverse)
					}
				}
			}
			if p.Default != nil {
				line += fmt.Sprintf(" default=%v", p.Default)
			}
			fmt.Println(line)
		}

		return nil
	},
}

func init() {
	typeCmd.AddCommand(typeShowCmd)
}
