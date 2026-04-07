package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/typemd/typemd/core"
)

var importOutputFlag string

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import external files into the vault",
}

var importScanCmd = &cobra.Command{
	Use:   "scan <paths...>",
	Short: "Scan source files and output a scan result",
	Long: `Scan source directories or files for markdown content, extracting
frontmatter patterns and collecting file statistics.

The scan result includes file metadata, frontmatter key patterns,
and existing vault types for AI-driven classification.

Examples:
  tmd import scan ~/notes
  tmd import scan ~/notes/blog ~/notes/ideas
  tmd import scan ~/notes/single-file.md`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		vault, err := openVault(vaultPath)
		if err != nil {
			return err
		}
		defer vault.Close()

		result, err := vault.ScanSources(args)
		if err != nil {
			return err
		}

		return printJSON(result)
	},
}

var importPlanCmd = &cobra.Command{
	Use:   "plan <classifications-file>",
	Short: "Generate an import plan from classifications",
	Long: `Generate an import plan from a JSON file containing object classifications.

The classifications file should contain a JSON array of ObjectPlan entries,
typically produced by an AI analyzing a scan result.

Examples:
  tmd import plan classifications.json
  tmd import plan classifications.json --output plan.json`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		vault, err := openVault(vaultPath)
		if err != nil {
			return err
		}
		defer vault.Close()

		data, err := os.ReadFile(args[0])
		if err != nil {
			return fmt.Errorf("reading classifications file: %w", err)
		}

		var classifications []core.ObjectPlan
		if err := json.Unmarshal(data, &classifications); err != nil {
			return fmt.Errorf("parsing classifications: %w", err)
		}

		plan, err := vault.GeneratePlan(classifications)
		if err != nil {
			return err
		}

		if importOutputFlag != "" {
			out, err := json.MarshalIndent(plan, "", "  ")
			if err != nil {
				return fmt.Errorf("marshal plan: %w", err)
			}
			if err := os.WriteFile(importOutputFlag, out, 0644); err != nil {
				return fmt.Errorf("writing plan file: %w", err)
			}
			fmt.Printf("Plan written to %s\n", importOutputFlag)
			return nil
		}

		return printJSON(plan)
	},
}

var importExecuteCmd = &cobra.Command{
	Use:   "execute <plan-file>",
	Short: "Execute an import plan",
	Long: `Execute an import plan to create types and objects in the vault.

The plan file should contain an ImportPlan JSON object, typically
produced by the plan command or constructed by an AI skill.

Examples:
  tmd import execute plan.json`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		vault, err := openVault(vaultPath)
		if err != nil {
			return err
		}
		defer vault.Close()

		data, err := os.ReadFile(args[0])
		if err != nil {
			return fmt.Errorf("reading plan file: %w", err)
		}

		var plan core.ImportPlan
		if err := json.Unmarshal(data, &plan); err != nil {
			return fmt.Errorf("parsing plan: %w", err)
		}

		report, err := vault.ExecutePlan(&plan)
		if err != nil {
			return err
		}

		return printJSON(report)
	},
}

func init() {
	importPlanCmd.Flags().StringVarP(&importOutputFlag, "output", "o", "", "write plan to file instead of stdout")
	importCmd.AddCommand(importScanCmd)
	importCmd.AddCommand(importPlanCmd)
	importCmd.AddCommand(importExecuteCmd)
	rootCmd.AddCommand(importCmd)
}
