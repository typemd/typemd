package cmd

import (
	"fmt"
	"strings"

	"github.com/typemd/typemd/core"
	"github.com/spf13/cobra"
)

var fixCmd = &cobra.Command{
	Use:   "fix",
	Short: "Fix common issues in the vault",
}

var fixWikiLinksCmd = &cobra.Command{
	Use:   "wikilinks",
	Short: "Expand shorthand wiki-links to full IDs",
	Long: `Walk all objects and replace shorthand wiki-link targets with their
resolved full IDs (type/name-ulid). Shorthand formats:

  [[type/name]]   — resolved within the specified type
  [[name]]        — resolved within the source object's type

Ambiguous targets (matching multiple objects) are reported but not expanded.`,
	Args:          cobra.NoArgs,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		vault, err := openVault(vaultPath, reindex)
		if err != nil {
			return err
		}
		defer vault.Close()

		result, err := vault.FixWikiLinks()
		if err != nil {
			return err
		}

		if result.Expanded == 0 && len(result.UnresolvedWikiLinks) == 0 {
			fmt.Println("All wiki-links are already full IDs. No changes needed.")
			return nil
		}

		if result.Expanded > 0 {
			fmt.Printf("Expanded %d wiki-link(s) to full IDs.\n", result.Expanded)
		}

		for _, u := range result.UnresolvedWikiLinks {
			if u.Reason == core.ReasonAmbiguous && len(u.Matches) > 0 {
				fmt.Printf("  %s: ambiguous [[%s]] — matches: %s\n",
					u.ObjectID, u.Target, strings.Join(u.Matches, ", "))
			} else {
				fmt.Printf("  %s: unresolved [[%s]]\n", u.ObjectID, u.Target)
			}
		}

		return nil
	},
}

func init() {
	fixCmd.AddCommand(fixWikiLinksCmd)
	rootCmd.AddCommand(fixCmd)
}
