package tui

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// tuiLog is the package-level logger for TUI debug output.
// Nil when logging is disabled (default until initLog is called).
var tuiLog *log.Logger

// initLog sets up file-based logging to .typemd/logs/{date}.log.
// Call once at TUI startup. Logging is best-effort; errors are silently ignored.
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

	tuiLog = log.New(f, "", log.LstdFlags)
	tuiLog.Println("=== TUI session started ===")

	return func() {
		tuiLog.Println("=== TUI session ended ===")
		f.Close()
		tuiLog = nil
	}
}

// logf writes a formatted log message if logging is enabled.
func logf(format string, args ...any) {
	if tuiLog != nil {
		tuiLog.Output(2, fmt.Sprintf(format, args...))
	}
}
