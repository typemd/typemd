package cmd

import (
	"io"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/typemd/typemd/core"
	"github.com/typemd/typemd/tui"
)

var (
	vaultPath string
	readOnly  bool
	debug     bool

	// Version is set at build time via ldflags.
	Version = "dev"
)

var rootCmd = &cobra.Command{
	Use:     "tmd",
	Short:   "A local-first CLI knowledge management tool",
	Version: Version,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if debug {
			core.InitLogging(slog.LevelDebug, os.Stderr)
		} else {
			core.InitLogging(slog.LevelWarn, io.Discard)
		}
		return nil
	},
	// Launch TUI when no subcommand is given
	RunE: func(cmd *cobra.Command, args []string) error {
		return tui.Start(vaultPath, readOnly)
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&vaultPath, "vault", "", "path to vault directory (default: current directory)")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug logging to stderr")
	rootCmd.Flags().BoolVar(&readOnly, "readonly", false, "open vault in read-only mode")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
