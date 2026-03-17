package core

import (
	"os"
	"path/filepath"
	"strings"
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

func TestSetConfigValue_EmptyValue(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)
	v.Init()
	v.Open()
	defer v.Close()

	// Set a value first
	if err := v.SetConfigValue("cli.default_type", "idea"); err != nil {
		t.Fatalf("SetConfigValue error: %v", err)
	}
	// Set to empty string
	if err := v.SetConfigValue("cli.default_type", ""); err != nil {
		t.Fatalf("SetConfigValue empty error: %v", err)
	}
	val, err := v.GetConfigValue("cli.default_type")
	if err != nil {
		t.Fatalf("GetConfigValue error: %v", err)
	}
	if val != "" {
		t.Errorf("expected empty value, got %q", val)
	}
}

func TestSetConfigValue_UnknownKey(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)
	v.Init()
	v.Open()
	defer v.Close()

	err := v.SetConfigValue("foo.bar", "baz")
	if err == nil {
		t.Fatal("expected error for unknown key")
	}
	if !strings.Contains(err.Error(), "unknown config key") {
		t.Errorf("error should mention 'unknown config key', got: %v", err)
	}
	if !strings.Contains(err.Error(), "cli.default_type") {
		t.Errorf("error should list known keys, got: %v", err)
	}
}

func TestGetConfigValue_UnknownKey(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)
	v.Init()
	v.Open()
	defer v.Close()

	_, err := v.GetConfigValue("nonexistent")
	if err == nil {
		t.Fatal("expected error for unknown key")
	}
	if !strings.Contains(err.Error(), "unknown config key") {
		t.Errorf("error should mention 'unknown config key', got: %v", err)
	}
}

func TestGetConfigValue_UnsetKnownKey(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)
	v.Init()
	v.Open()
	defer v.Close()

	val, err := v.GetConfigValue("cli.default_type")
	if err != nil {
		t.Fatalf("GetConfigValue error: %v", err)
	}
	if val != "" {
		t.Errorf("expected empty value for unset key, got %q", val)
	}
}

func TestConfigKeys_Sorted(t *testing.T) {
	keys := ConfigKeys()
	if len(keys) == 0 {
		t.Fatal("expected at least one config key")
	}
	for i := 1; i < len(keys); i++ {
		if keys[i] < keys[i-1] {
			t.Errorf("keys not sorted: %v", keys)
			break
		}
	}
}

func TestSetConfigValue_CreatesConfigFile(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)
	v.Init()
	v.Open()
	defer v.Close()

	configPath := filepath.Join(v.Dir(), configFileName)
	if _, err := os.Stat(configPath); !os.IsNotExist(err) {
		t.Fatal("config file should not exist before set")
	}

	if err := v.SetConfigValue("cli.default_type", "idea"); err != nil {
		t.Fatalf("SetConfigValue error: %v", err)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("config file should exist after set")
	}

	data, _ := os.ReadFile(configPath)
	if !strings.Contains(string(data), "default_type: idea") {
		t.Errorf("config file should contain default_type: idea, got:\n%s", data)
	}
}
