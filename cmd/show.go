package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show <object-id>",
	Short: "Show object detail (properties, relations, body)",
	Long: `Display an object's properties, relations, and body content.

Examples:
  tmd show book/clean-code-01jqr3k5mpbvn8e0f2g7h9txyz
  tmd show person/robert-martin-01jqr3k8yznw2a4dbx6t7c9fpq`,
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

		props, err := vault.BuildDisplayProperties(obj)
		if err != nil {
			return fmt.Errorf("build display properties: %w", err)
		}

		// Title
		fmt.Println(obj.ID)
		fmt.Println()

		// Properties & Relations
		fmt.Println("Properties")
		fmt.Println("──────────")
		if len(props) == 0 {
			fmt.Println("  (none)")
		} else {
			for _, p := range props {
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
