package cmd

import (
	"os"
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

	// Verify file was created
	objPath := filepath.Join(dir, "objects", "book", "clean-code.md")
	if _, err := os.Stat(objPath); os.IsNotExist(err) {
		t.Error("expected object file to exist")
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

func TestCreateCmd_Duplicate(t *testing.T) {
	dir := setupTestVaultDir(t)

	vaultPath = dir

	// Create first
	rootCmd.SetArgs([]string{"create", "book", "duplicate"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("first create error = %v", err)
	}

	// Create duplicate
	rootCmd.SetArgs([]string{"create", "book", "duplicate"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for duplicate object")
	}
}
