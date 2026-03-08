package cmd

import (
	"github.com/spf13/cobra"
)

var objectCmd = &cobra.Command{
	Use:   "object",
	Short: "Manage objects (create, show, list)",
}

func init() {
	rootCmd.AddCommand(objectCmd)
}
