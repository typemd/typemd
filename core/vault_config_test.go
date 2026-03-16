package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadVaultConfig_UnknownKeysIgnored(t *testing.T) {
	dir := t.TempDir()
	content := "unknown_key: value\ncli:\n  default_type: idea\n"
	os.WriteFile(filepath.Join(dir, configFileName), []byte(content), 0644)

	cfg, err := loadVaultConfig(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.CLI.DefaultType != "idea" {
		t.Errorf("expected default_type=idea, got %q", cfg.CLI.DefaultType)
	}
}

func TestLoadVaultConfig_PartialConfig_OnlyCLI(t *testing.T) {
	dir := t.TempDir()
	content := "cli:\n  default_type: note\n"
	os.WriteFile(filepath.Join(dir, configFileName), []byte(content), 0644)

	cfg, err := loadVaultConfig(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.CLI.DefaultType != "note" {
		t.Errorf("expected default_type=note, got %q", cfg.CLI.DefaultType)
	}
}

func TestLoadVaultConfig_PartialConfig_CLIWithoutDefaultType(t *testing.T) {
	dir := t.TempDir()
	content := "cli:\n  other_field: value\n"
	os.WriteFile(filepath.Join(dir, configFileName), []byte(content), 0644)

	cfg, err := loadVaultConfig(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.CLI.DefaultType != "" {
		t.Errorf("expected empty default_type, got %q", cfg.CLI.DefaultType)
	}
}

func TestLoadVaultConfig_MissingFile(t *testing.T) {
	dir := t.TempDir()

	cfg, err := loadVaultConfig(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.CLI.DefaultType != "" {
		t.Errorf("expected empty default_type, got %q", cfg.CLI.DefaultType)
	}
}

func TestLoadVaultConfig_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, configFileName), []byte(""), 0644)

	cfg, err := loadVaultConfig(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.CLI.DefaultType != "" {
		t.Errorf("expected empty default_type, got %q", cfg.CLI.DefaultType)
	}
}

func TestLoadVaultConfig_InvalidYAML(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, configFileName), []byte("[invalid: yaml: content"), 0644)

	_, err := loadVaultConfig(dir)
	if err == nil {
		t.Fatal("expected error for invalid YAML, got nil")
	}
}
