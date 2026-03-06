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
				if p.IsReverse {
					fmt.Printf("  %s: ← %s\n", p.Key, p.FromID)
				} else if p.Value == nil {
					fmt.Printf("  %s: (null)\n", p.Key)
				} else if p.IsRelation {
					fmt.Printf("  %s: → %v\n", p.Key, p.Value)
				} else {
					fmt.Printf("  %s: %v\n", p.Key, p.Value)
				}
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
