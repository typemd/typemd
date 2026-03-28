package tui

import (
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/typemd/typemd/core"
)

// initLog sets up structured logging to .typemd/logs/{date}.log.
// The TUI always logs at DEBUG level to file (never to stderr, which would
// corrupt terminal rendering). Returns a cleanup function that closes the file.
func initLog(vaultRoot string) func() {
	dir := filepath.Join(vaultRoot, ".typemd", "logs")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return func() {}
	}

	filename := time.Now().Format("2006-01-02") + ".log"
	path := filepath.Join(dir, filename)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return func() {}
	}

	core.InitLogging(slog.LevelDebug, f)
	slog.Info("tui session started")

	return func() {
		slog.Info("tui session ended")
		f.Close()
	}
}
