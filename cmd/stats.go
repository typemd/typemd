package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/typemd/typemd/core"
)

var statsTypeName string
var statsJSON bool

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show aggregate statistics about the vault",
	Long: `Show aggregate statistics about the vault.

Without flags, displays a per-type summary with object counts.
With --type, displays per-property aggregations for a specific type.

Examples:
  tmd stats
  tmd stats --type book
  tmd stats --type book --json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		vault, err := openVault(vaultPath, reindex)
		if err != nil {
			return err
		}
		defer vault.Close()

		if statsTypeName != "" {
			return runTypeStats(vault, statsTypeName, statsJSON)
		}
		return runVaultStats(vault, statsJSON)
	},
}

func runVaultStats(vault *core.Vault, asJSON bool) error {
	stats, err := vault.VaultStats()
	if err != nil {
		return err
	}

	if asJSON {
		return printJSON(stats)
	}

	if len(stats.Types) == 0 {
		fmt.Println("No objects in vault.")
		return nil
	}

	// Find max name length for alignment
	maxLen := 0
	for _, ts := range stats.Types {
		name := ts.DisplayName()
		if len(name) > maxLen {
			maxLen = len(name)
		}
	}

	for _, ts := range stats.Types {
		emoji := ts.Emoji
		if emoji == "" {
			emoji = " "
		}

		lastUpdated := ""
		if !ts.LastUpdated.IsZero() {
			lastUpdated = fmt.Sprintf("   last updated %s", ts.LastUpdated.Format("2006-01-02"))
		}

		fmt.Printf("%s %-*s  %3d%s\n", emoji, maxLen, ts.DisplayName(), ts.Count, lastUpdated)
	}

	fmt.Println(strings.Repeat("─", maxLen+20))
	fmt.Printf("  %-*s  %3d\n", maxLen, "Total", stats.Total)

	return nil
}

func runTypeStats(vault *core.Vault, typeName string, asJSON bool) error {
	stats, err := vault.TypeStats(typeName)
	if err != nil {
		return err
	}

	if asJSON {
		return printJSON(stats)
	}

	emoji := ""
	if stats.Emoji != "" {
		emoji = stats.Emoji + " "
	}

	fmt.Printf("%s%s (%d objects)\n", emoji, typeName, stats.Count)
	fmt.Println(strings.Repeat("─", 40))

	if len(stats.Properties) == 0 {
		fmt.Println("No properties defined.")
		return nil
	}

	for _, ps := range stats.Properties {
		fmt.Println()
		fmt.Printf("%s (%s)\n", ps.Name, ps.Type)
		fmt.Printf("  filled: %d/%d\n", ps.Filled, ps.Total)

		if ps.Stats == nil {
			continue
		}

		switch s := ps.Stats.(type) {
		case *core.NumberStats:
			fmt.Printf("  sum: %g  avg: %g  min: %g  max: %g\n", s.Sum, s.Avg, s.Min, s.Max)
		case *core.SelectStats:
			for option, count := range s.Distribution {
				bar := strings.Repeat("█", count)
				fmt.Printf("  %-15s %s %d\n", option, bar, count)
			}
		case *core.CheckboxStats:
			fmt.Printf("  true: %d  false: %d\n", s.TrueCount, s.FalseCount)
		case *core.DateStats:
			fmt.Printf("  earliest: %s  latest: %s\n", s.Earliest, s.Latest)
		case *core.RelationStats:
			fmt.Printf("  links: %d\n", s.Count)
		}
	}

	return nil
}

func init() {
	statsCmd.Flags().StringVar(&statsTypeName, "type", "", "Show stats for a specific type")
	statsCmd.Flags().BoolVar(&statsJSON, "json", false, "Output as JSON")
	rootCmd.AddCommand(statsCmd)
}
