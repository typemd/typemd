package cmd

import (
	"github.com/spf13/cobra"
)

var typeCmd = &cobra.Command{
	Use:   "type",
	Short: "Manage type schemas (show, list, validate)",
}

func init() {
	rootCmd.AddCommand(typeCmd)
}
