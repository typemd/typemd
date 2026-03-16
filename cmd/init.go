package cmd

import (
	"fmt"
	"slices"

	tea "charm.land/bubbletea/v2"
	"github.com/spf13/cobra"
	"github.com/typemd/typemd/core"
)

var noStarters bool

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new vault with optional starter types",
	RunE: func(cmd *cobra.Command, args []string) error {
		vault := resolveVault(vaultPath)
		if err := vault.Init(); err != nil {
			return err
		}
		fmt.Printf("Initialized vault at %s\n", vault.Dir())

		if noStarters {
			return nil
		}

		// Run the starter type picker
		picker := newStarterPicker()
		p := tea.NewProgram(picker)
		finalModel, err := p.Run()
		if err != nil {
			return fmt.Errorf("starter picker: %w", err)
		}
		final, ok := finalModel.(starterPicker)
		if !ok {
			return fmt.Errorf("unexpected model type from starter picker")
		}
		selected := final.selectedItems()

		if len(selected) == 0 {
			return nil
		}

		names := make([]string, len(selected))
		for i, item := range selected {
			names[i] = item.name
		}
		if err := vault.WriteStarterTypes(names); err != nil {
			return err
		}

		// Write config.yaml with default_type if a quick-suitable type was selected
		if defaultType := resolveDefaultType(names); defaultType != "" {
			cfg := &core.VaultConfig{
				CLI: core.CLIConfig{DefaultType: defaultType},
			}
			if err := vault.WriteConfig(cfg); err != nil {
				return fmt.Errorf("write config: %w", err)
			}
		}

		fmt.Println()
		fmt.Printf("Created %d starter type(s):\n", len(selected))
		for _, item := range selected {
			fmt.Printf("  %s %s\n", item.emoji, item.name)
		}

		return nil
	},
}

// resolveDefaultType picks the best default type from selected starter names.
// Returns "idea" if selected, then "note" as fallback, or empty if neither.
func resolveDefaultType(names []string) string {
	if slices.Contains(names, "idea") {
		return "idea"
	}
	if slices.Contains(names, "note") {
		return "note"
	}
	return ""
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolVar(&noStarters, "no-starters", false, "skip starter type selection")
}
