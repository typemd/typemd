package cmd

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/typemd/typemd/core"
)

func TestValidateCmd_TagUniqueness(t *testing.T) {
	dir := t.TempDir()
	v := core.NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	if err := v.Open(); err != nil {
		t.Fatalf("Open() error = %v", err)
	}

	// Create one tag via normal API
	if _, err := v.NewObject("tag", "go", ""); err != nil {
		t.Fatalf("NewObject error = %v", err)
	}

	// Write a raw duplicate tag file on disk to bypass uniqueness check
	dupeDir := filepath.Join(dir, "objects", "tag")
	os.WriteFile(filepath.Join(dupeDir, "go-01jjjjjjjjjjjjjjjjjjjjjjjj.md"), []byte("---\nname: go\n---\n"), 0644)

	v.Close()

	// Run validate command — index is synced on vault open
	resetAllFlags()
	vaultPath = dir
	rootCmd.SetArgs([]string{"type", "validate"})
	err := rootCmd.Execute()

	if err == nil {
		t.Fatal("expected validation error for duplicate tag names")
	}
	if !strings.Contains(err.Error(), "validation error") {
		t.Errorf("error = %q, want it to mention validation error", err)
	}
}

func TestAddWatchPaths_MissingDirectories(t *testing.T) {
	dir := t.TempDir()
	v := core.NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	if err := v.Open(); err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer v.Close()

	// Remove objects dir to test missing directory handling
	os.RemoveAll(filepath.Join(dir, "objects"))

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		t.Fatalf("NewWatcher() error = %v", err)
	}
	defer watcher.Close()

	// Should not panic with missing directories
	addWatchPaths(watcher, v)
}

func TestRunWatchValidation_ContextCancel(t *testing.T) {
	dir := t.TempDir()
	v := core.NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	if err := v.Open(); err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer v.Close()

	// Create a context that cancels immediately
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := runWatchValidation(ctx, v)
	if err != nil {
		t.Errorf("runWatchValidation() error = %v, want nil on context cancel", err)
	}
}

func TestValidateCmd_NoDuplicateTags(t *testing.T) {
	dir := t.TempDir()
	v := core.NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	if err := v.Open(); err != nil {
		t.Fatalf("Open() error = %v", err)
	}

	v.NewObject("tag", "go", "")
	v.NewObject("tag", "rust", "")
	v.Close()

	resetAllFlags()
	vaultPath = dir
	rootCmd.SetArgs([]string{"type", "validate"})
	err := rootCmd.Execute()
	if err != nil {
		t.Errorf("expected no validation errors, got %v", err)
	}
}
