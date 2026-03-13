package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

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
	if _, err := v.NewObject("tag", "go"); err != nil {
		t.Fatalf("NewObject error = %v", err)
	}

	// Write a raw duplicate tag file on disk to bypass uniqueness check
	dupeDir := filepath.Join(dir, "objects", "tag")
	os.WriteFile(filepath.Join(dupeDir, "go-01jjjjjjjjjjjjjjjjjjjjjjjj.md"), []byte("---\nname: go\n---\n"), 0644)

	v.Close()

	// Run validate command with reindex to pick up the raw file
	vaultPath = dir
	reindex = true
	defer func() { reindex = false }()
	rootCmd.SetArgs([]string{"type", "validate"})
	err := rootCmd.Execute()

	if err == nil {
		t.Fatal("expected validation error for duplicate tag names")
	}
	if !strings.Contains(err.Error(), "validation error") {
		t.Errorf("error = %q, want it to mention validation error", err)
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

	v.NewObject("tag", "go")
	v.NewObject("tag", "rust")
	v.Close()

	vaultPath = dir
	rootCmd.SetArgs([]string{"type", "validate"})
	err := rootCmd.Execute()
	if err != nil {
		t.Errorf("expected no validation errors, got %v", err)
	}
}
