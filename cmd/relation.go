package cmd

import (
	"github.com/spf13/cobra"
)

var relationCmd = &cobra.Command{
	Use:   "relation",
	Short: "Manage relations between objects (link, unlink)",
}

func init() {
	rootCmd.AddCommand(relationCmd)
}
