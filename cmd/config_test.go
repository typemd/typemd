package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/typemd/typemd/core"
)

func setupConfigVault(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	v := core.NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	v.Close()
	return dir
}

func TestConfigSetCmd_ValidKey(t *testing.T) {
	dir := setupConfigVault(t)
	vaultPath = dir

	rootCmd.SetArgs([]string{"config", "set", "cli.default_type", "idea"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify value was written
	data, err := os.ReadFile(filepath.Join(dir, ".typemd", "config.yaml"))
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	if !strings.Contains(string(data), "default_type: idea") {
		t.Errorf("config should contain default_type: idea, got:\n%s", data)
	}
}

func TestConfigSetCmd_UnknownKey(t *testing.T) {
	dir := setupConfigVault(t)
	vaultPath = dir

	rootCmd.SetArgs([]string{"config", "set", "unknown.key", "value"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for unknown key")
	}
	if !strings.Contains(err.Error(), "unknown config key") {
		t.Errorf("error should mention 'unknown config key', got: %v", err)
	}
}

func TestConfigGetCmd_SetKey(t *testing.T) {
	dir := setupConfigVault(t)
	vaultPath = dir

	// Set a value first
	configContent := "cli:\n  default_type: idea\n"
	os.WriteFile(filepath.Join(dir, ".typemd", "config.yaml"), []byte(configContent), 0644)

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rootCmd.SetArgs([]string{"config", "get", "cli.default_type"})
	err := rootCmd.Execute()

	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := strings.TrimSpace(buf.String())
	if output != "idea" {
		t.Errorf("expected 'idea', got %q", output)
	}
}

func TestConfigGetCmd_UnsetKey(t *testing.T) {
	dir := setupConfigVault(t)
	vaultPath = dir

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rootCmd.SetArgs([]string{"config", "get", "cli.default_type"})
	err := rootCmd.Execute()

	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := strings.TrimSpace(buf.String())
	if output != "" {
		t.Errorf("expected empty output, got %q", output)
	}
}

func TestConfigGetCmd_UnknownKey(t *testing.T) {
	dir := setupConfigVault(t)
	vaultPath = dir

	rootCmd.SetArgs([]string{"config", "get", "unknown.key"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for unknown key")
	}
	if !strings.Contains(err.Error(), "unknown config key") {
		t.Errorf("error should mention 'unknown config key', got: %v", err)
	}
}

func TestConfigListCmd_WithValues(t *testing.T) {
	dir := setupConfigVault(t)
	vaultPath = dir

	configContent := "cli:\n  default_type: idea\n"
	os.WriteFile(filepath.Join(dir, ".typemd", "config.yaml"), []byte(configContent), 0644)

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rootCmd.SetArgs([]string{"config", "list"})
	err := rootCmd.Execute()

	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := strings.TrimSpace(buf.String())
	if !strings.Contains(output, "cli.default_type: idea") {
		t.Errorf("expected 'cli.default_type: idea', got %q", output)
	}
}

func TestConfigListCmd_Empty(t *testing.T) {
	dir := setupConfigVault(t)
	vaultPath = dir

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rootCmd.SetArgs([]string{"config", "list"})
	err := rootCmd.Execute()

	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := strings.TrimSpace(buf.String())
	if output != "" {
		t.Errorf("expected empty output, got %q", output)
	}
}
