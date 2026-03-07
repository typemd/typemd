package cmd

import (
	"path/filepath"
	"testing"

	"github.com/typemd/typemd/core"
)

func setupTestVaultDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	v := core.NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	if err := v.Open(); err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	v.Close()
	return dir
}

func TestCreateCmd_Success(t *testing.T) {
	dir := setupTestVaultDir(t)

	vaultPath = dir
	rootCmd.SetArgs([]string{"create", "book", "clean-code"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	// Verify file matching pattern was created (filename includes ULID suffix)
	matches, _ := filepath.Glob(filepath.Join(dir, "objects", "book", "clean-code-*.md"))
	if len(matches) != 1 {
		t.Errorf("expected 1 matching file, got %d", len(matches))
	}
}

func TestCreateCmd_UnknownType(t *testing.T) {
	dir := setupTestVaultDir(t)

	vaultPath = dir
	rootCmd.SetArgs([]string{"create", "nonexistent", "foo"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for unknown type")
	}
}

func TestCreateCmd_SameNameTwice(t *testing.T) {
	dir := setupTestVaultDir(t)

	vaultPath = dir

	// Create first
	rootCmd.SetArgs([]string{"create", "book", "duplicate"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("first create error = %v", err)
	}

	// Create with same name — should succeed because ULID guarantees uniqueness
	rootCmd.SetArgs([]string{"create", "book", "duplicate"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("second create should succeed with ULID, error = %v", err)
	}

	matches, _ := filepath.Glob(filepath.Join(dir, "objects", "book", "duplicate-*.md"))
	if len(matches) != 2 {
		t.Errorf("expected 2 matching files, got %d", len(matches))
	}
}
