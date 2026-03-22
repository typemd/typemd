package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/typemd/typemd/core"
)

var (
	instructionsJSON  bool
	instructionsSkill bool
)

var instructionsCmd = &cobra.Command{
	Use:   "instructions [skill]",
	Short: "Output skill instructions enriched with vault context",
	Long: `Output skill instructions enriched with vault context.

Without arguments, lists all available skills.
With a skill name, outputs the skill instructions as JSON with vault context.

Examples:
  tmd instructions
  tmd instructions --json
  tmd instructions explore
  tmd instructions explore --skill`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return runListSkills(instructionsJSON)
		}
		return runGetSkill(args[0], instructionsSkill)
	},
}

func runListSkills(asJSON bool) error {
	skills := core.ListSkills()

	if asJSON {
		return printJSON(skills)
	}

	// Find max name length for alignment
	maxLen := 0
	for _, s := range skills {
		if len(s.Name) > maxLen {
			maxLen = len(s.Name)
		}
	}

	for _, s := range skills {
		fmt.Printf("%-*s  %s\n", maxLen, s.Name, s.Description)
	}
	return nil
}

func runGetSkill(name string, rawSkill bool) error {
	if rawSkill {
		return runGetSkillRaw(name)
	}

	// Try to open vault for context (optional)
	vault, vaultErr := openVault(vaultPath, reindex)
	var vaultRoot string
	if vaultErr == nil {
		defer vault.Close()
		vaultRoot = vault.Root
	}

	skill, err := core.GetSkillWithOverride(name, vaultRoot)
	if err != nil {
		return err
	}

	output := core.SkillOutput{
		Name:         skill.Name,
		Description:  skill.Description,
		Instructions: skill.Instructions,
	}

	if vault != nil {
		ctx, err := core.BuildSkillContext(vault)
		if err == nil {
			output.Context = ctx
		}
	}

	return printJSON(output)
}

func runGetSkillRaw(name string) error {
	var vaultRoot string
	vault, err := openVault(vaultPath, reindex)
	if err == nil {
		defer vault.Close()
		vaultRoot = vault.Root
	}

	raw, err := core.GetSkillRawWithOverride(name, vaultRoot)
	if err != nil {
		return err
	}

	fmt.Print(string(raw))
	return nil
}

func init() {
	instructionsCmd.Flags().BoolVar(&instructionsJSON, "json", false, "Output as JSON")
	instructionsCmd.Flags().BoolVar(&instructionsSkill, "skill", false, "Output raw SKILL.md content with frontmatter")
	rootCmd.AddCommand(instructionsCmd)
}
