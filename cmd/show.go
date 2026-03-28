package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/typemd/typemd/core"
)

var showCmd = &cobra.Command{
	Use:   "show <object-id>",
	Short: "Show object detail (properties, relations, body)",
	Long: `Display an object's properties, relations, and body content.

Supports prefix matching — you can omit the ULID suffix if the prefix
uniquely identifies an object. If a prefix matches multiple objects,
an interactive picker is shown to select the intended one.

Examples:
  tmd object show book/clean-code
  tmd object show book/clean-code-01jqr3k5mpbvn8e0f2g7h9txyz
  tmd object show person/robert-martin`,
	Args: cobra.ExactArgs(1),
	ValidArgsFunction: func(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return completeObjectID(toComplete)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		vault, err := openVault(vaultPath)
		if err != nil {
			return err
		}
		defer vault.Close()

		obj, err := resolveObjectInteractive(vault, args[0])
		if err != nil {
			return err
		}

		props, err := vault.BuildDisplayProperties(obj)
		if err != nil {
			return fmt.Errorf("build display properties: %w", err)
		}

		// Title
		fmt.Println(obj.DisplayID())
		fmt.Println()

		// Split properties into schema and local
		var schemaProps, localProps []core.DisplayProperty
		for _, p := range props {
			if p.IsLocal {
				localProps = append(localProps, p)
			} else {
				schemaProps = append(schemaProps, p)
			}
		}

		// Properties & Relations
		fmt.Println("Properties")
		fmt.Println("──────────")
		if len(schemaProps) == 0 {
			fmt.Println("  (none)")
		} else {
			for _, p := range schemaProps {
				fmt.Printf("  %s\n", p.Format())
			}
		}

		// Local Properties (not in type schema)
		if len(localProps) > 0 {
			fmt.Println()
			fmt.Println("Local Properties")
			fmt.Println("────────────────")
			for _, p := range localProps {
				fmt.Printf("  %s\n", p.Format())
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
			body = core.RenderWikiLinks(body)
			for _, line := range strings.Split(body, "\n") {
				fmt.Printf("  %s\n", line)
			}
		}

		return nil
	},
}

func init() {
	objectCmd.AddCommand(showCmd)
}
