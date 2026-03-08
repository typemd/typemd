package cmd

import (
	"os"

	"github.com/typemd/typemd/tui"
	"github.com/spf13/cobra"
)

var (
	vaultPath string
	readOnly  bool
	reindex   bool

	// Version is set at build time via ldflags.
	Version = "dev"
)

var rootCmd = &cobra.Command{
	Use:     "tmd",
	Short:   "A local-first CLI knowledge management tool",
	Version: Version,
	// 不帶子指令時啟動 TUI
	RunE: func(cmd *cobra.Command, args []string) error {
		return tui.Start(vaultPath, readOnly, reindex)
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&vaultPath, "vault", "", "path to vault directory (default: current directory)")
	rootCmd.PersistentFlags().BoolVar(&reindex, "reindex", false, "force reindex before running")
	rootCmd.Flags().BoolVar(&readOnly, "readonly", false, "open vault in read-only mode")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
