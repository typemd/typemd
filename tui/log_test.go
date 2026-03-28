package tui

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestInitLog_CreatesLogFile(t *testing.T) {
	dir := t.TempDir()

	// Create the .typemd directory structure
	vaultRoot := dir
	logsDir := filepath.Join(vaultRoot, ".typemd", "logs")

	cleanup := initLog(vaultRoot)
	defer cleanup()

	expectedFile := filepath.Join(logsDir, time.Now().Format("2006-01-02")+".log")
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		t.Fatalf("expected log file %s to exist", expectedFile)
	}

	// Read the log file and verify it contains the session start message
	data, err := os.ReadFile(expectedFile)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}
	content := string(data)
	if len(content) == 0 {
		t.Fatal("expected log file to have content")
	}
}

func TestInitLog_InvalidPath(t *testing.T) {
	// initLog with invalid path should return a no-op cleanup
	cleanup := initLog("/nonexistent/path/that/cannot/exist")
	cleanup() // should not panic
}
