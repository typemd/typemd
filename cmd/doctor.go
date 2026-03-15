package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/typemd/typemd/core"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check vault health and report issues",
	Long: `Performs a comprehensive vault health check across 8 categories:

  Schemas      Type schema validation
  Objects      Object property validation
  Relations    Relation integrity
  Wiki-links   Broken link detection
  Uniqueness   Duplicate name detection
  Files        Corrupted file detection
  Index        SQLite sync status (auto-fixed)
  Orphans      Directories without type schemas

A superset of "tmd type validate" with additional structural integrity checks.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		vault, err := openVault(vaultPath, reindex)
		if err != nil {
			return err
		}
		defer vault.Close()

		report := core.RunDoctor(vault)

		printDoctorReport(report)

		if report.HasUnresolvedIssues() {
			return fmt.Errorf("found %d issue(s)", report.TotalIssues())
		}
		return nil
	},
}

func printDoctorReport(report *core.DoctorReport) {
	for _, cat := range report.Categories {
		if len(cat.Issues) == 0 && cat.AutoFixed == 0 {
			fmt.Printf("  ✓ %s\n", cat.Name)
		} else if len(cat.Issues) == 0 && cat.AutoFixed > 0 {
			fmt.Printf("  ✓ %s (auto-fixed)\n", cat.Name)
		} else {
			fmt.Printf("  ✗ %s\n", cat.Name)
			for _, issue := range cat.Issues {
				prefix := "error"
				if issue.Severity == core.SeverityWarning {
					prefix = "warn"
				}
				fmt.Printf("    [%s] %s\n", prefix, issue.Message)
			}
		}
	}
	fmt.Println()

	total := report.TotalIssues()
	fixed := report.TotalAutoFixed()
	if total == 0 && fixed == 0 {
		fmt.Println("No issues found.")
	} else if total == 0 && fixed > 0 {
		fmt.Printf("No issues found, %d auto-fixed.\n", fixed)
	} else if fixed > 0 {
		fmt.Printf("%d issue(s) found, %d auto-fixed.\n", total, fixed)
	} else {
		fmt.Printf("%d issue(s) found.\n", total)
	}
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
