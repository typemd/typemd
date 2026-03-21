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

func TestLoadVaultConfig_TUIDebounceMs(t *testing.T) {
	dir := t.TempDir()
	content := "tui:\n  debounce_ms: 500\n"
	os.WriteFile(filepath.Join(dir, configFileName), []byte(content), 0644)

	cfg, err := loadVaultConfig(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.TUI.DebounceMs != 500 {
		t.Errorf("expected debounce_ms=500, got %d", cfg.TUI.DebounceMs)
	}
}

func TestLoadVaultConfig_TUIDebounceMs_NotSet(t *testing.T) {
	dir := t.TempDir()
	content := "cli:\n  default_type: page\n"
	os.WriteFile(filepath.Join(dir, configFileName), []byte(content), 0644)

	cfg, err := loadVaultConfig(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.TUI.DebounceMs != 0 {
		t.Errorf("expected debounce_ms=0, got %d", cfg.TUI.DebounceMs)
	}
}

func TestGetSetConfigValue_TUIDebounceMs(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)
	v.Init()
	v.Open()
	defer v.Close()

	// Set debounce_ms
	if err := v.SetConfigValue("tui.debounce_ms", "300"); err != nil {
		t.Fatalf("SetConfigValue error: %v", err)
	}
	val, err := v.GetConfigValue("tui.debounce_ms")
	if err != nil {
		t.Fatalf("GetConfigValue error: %v", err)
	}
	if val != "300" {
		t.Errorf("expected '300', got %q", val)
	}

	// Verify struct field
	if v.config.TUI.DebounceMs != 300 {
		t.Errorf("expected DebounceMs=300, got %d", v.config.TUI.DebounceMs)
	}
}

func TestGetConfigValue_TUIDebounceMs_Unset(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)
	v.Init()
	v.Open()
	defer v.Close()

	val, err := v.GetConfigValue("tui.debounce_ms")
	if err != nil {
		t.Fatalf("GetConfigValue error: %v", err)
	}
	if val != "" {
		t.Errorf("expected empty value for unset debounce_ms, got %q", val)
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

func TestStatsTypeLayoutDefault(t *testing.T) {
	cfg := &VaultConfig{}
	if cfg.TUI.StatsTypeLayout != "" {
		t.Errorf("expected empty default, got %q", cfg.TUI.StatsTypeLayout)
	}
}

func TestStatsTypeLayoutConfigKey(t *testing.T) {
	entry, ok := configKeyRegistry["tui.stats_type_layout"]
	if !ok {
		t.Fatal("tui.stats_type_layout not in config registry")
	}
	cfg := &VaultConfig{}
	// Default should return empty (consumer applies "fullscreen" default)
	if got := entry.Get(cfg); got != "" {
		t.Errorf("expected empty default, got %q", got)
	}
	// Set popup
	entry.Set(cfg, "popup")
	if cfg.TUI.StatsTypeLayout != "popup" {
		t.Errorf("expected 'popup', got %q", cfg.TUI.StatsTypeLayout)
	}
	if got := entry.Get(cfg); got != "popup" {
		t.Errorf("expected 'popup' from getter, got %q", got)
	}
}

// --- Toast config tests ---

func TestLoadVaultConfig_ToastConfig_AllFields(t *testing.T) {
	dir := t.TempDir()
	content := `tui:
  toast:
    position: help-bar
    duration_ms: 5000
    dismiss_key: q
    show_warnings: false
    show_success: true
`
	os.WriteFile(filepath.Join(dir, configFileName), []byte(content), 0644)

	cfg, err := loadVaultConfig(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.TUI.Toast.Position != "help-bar" {
		t.Errorf("expected position=help-bar, got %q", cfg.TUI.Toast.Position)
	}
	if cfg.TUI.Toast.DurationMs != 5000 {
		t.Errorf("expected duration_ms=5000, got %d", cfg.TUI.Toast.DurationMs)
	}
	if cfg.TUI.Toast.DismissKey != "q" {
		t.Errorf("expected dismiss_key=q, got %q", cfg.TUI.Toast.DismissKey)
	}
	if cfg.TUI.Toast.ShowWarnings == nil {
		t.Fatal("expected show_warnings to be non-nil")
	}
	if *cfg.TUI.Toast.ShowWarnings != false {
		t.Errorf("expected show_warnings=false, got %v", *cfg.TUI.Toast.ShowWarnings)
	}
	if cfg.TUI.Toast.ShowSuccess == nil {
		t.Fatal("expected show_success to be non-nil")
	}
	if *cfg.TUI.Toast.ShowSuccess != true {
		t.Errorf("expected show_success=true, got %v", *cfg.TUI.Toast.ShowSuccess)
	}
}

func TestLoadVaultConfig_ToastConfig_BoolNilWhenAbsent(t *testing.T) {
	dir := t.TempDir()
	content := `tui:
  toast:
    position: bottom-right
`
	os.WriteFile(filepath.Join(dir, configFileName), []byte(content), 0644)

	cfg, err := loadVaultConfig(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.TUI.Toast.Position != "bottom-right" {
		t.Errorf("expected position=bottom-right, got %q", cfg.TUI.Toast.Position)
	}
	if cfg.TUI.Toast.ShowWarnings != nil {
		t.Errorf("expected show_warnings=nil when absent, got %v", *cfg.TUI.Toast.ShowWarnings)
	}
	if cfg.TUI.Toast.ShowSuccess != nil {
		t.Errorf("expected show_success=nil when absent, got %v", *cfg.TUI.Toast.ShowSuccess)
	}
}

func TestLoadVaultConfig_ToastConfig_BoolExplicitTrue(t *testing.T) {
	dir := t.TempDir()
	content := `tui:
  toast:
    show_warnings: true
    show_success: true
`
	os.WriteFile(filepath.Join(dir, configFileName), []byte(content), 0644)

	cfg, err := loadVaultConfig(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.TUI.Toast.ShowWarnings == nil || *cfg.TUI.Toast.ShowWarnings != true {
		t.Errorf("expected show_warnings=true")
	}
	if cfg.TUI.Toast.ShowSuccess == nil || *cfg.TUI.Toast.ShowSuccess != true {
		t.Errorf("expected show_success=true")
	}
}

func TestLoadVaultConfig_ToastConfig_NotSet(t *testing.T) {
	dir := t.TempDir()
	content := "cli:\n  default_type: page\n"
	os.WriteFile(filepath.Join(dir, configFileName), []byte(content), 0644)

	cfg, err := loadVaultConfig(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.TUI.Toast.Position != "" {
		t.Errorf("expected empty position, got %q", cfg.TUI.Toast.Position)
	}
	if cfg.TUI.Toast.DurationMs != 0 {
		t.Errorf("expected duration_ms=0, got %d", cfg.TUI.Toast.DurationMs)
	}
	if cfg.TUI.Toast.DismissKey != "" {
		t.Errorf("expected empty dismiss_key, got %q", cfg.TUI.Toast.DismissKey)
	}
	if cfg.TUI.Toast.ShowWarnings != nil {
		t.Errorf("expected show_warnings=nil, got %v", *cfg.TUI.Toast.ShowWarnings)
	}
	if cfg.TUI.Toast.ShowSuccess != nil {
		t.Errorf("expected show_success=nil, got %v", *cfg.TUI.Toast.ShowSuccess)
	}
}

func TestToastConfigRegistry_Position(t *testing.T) {
	entry, ok := configKeyRegistry["tui.toast.position"]
	if !ok {
		t.Fatal("tui.toast.position not in config registry")
	}
	if entry.Default != "bottom-right" {
		t.Errorf("expected default 'bottom-right', got %q", entry.Default)
	}

	cfg := &VaultConfig{}
	if got := entry.Get(cfg); got != "" {
		t.Errorf("expected empty for unset, got %q", got)
	}
	entry.Set(cfg, "help-bar")
	if cfg.TUI.Toast.Position != "help-bar" {
		t.Errorf("expected 'help-bar', got %q", cfg.TUI.Toast.Position)
	}
	if got := entry.Get(cfg); got != "help-bar" {
		t.Errorf("expected 'help-bar' from getter, got %q", got)
	}
}

func TestToastConfigRegistry_DurationMs(t *testing.T) {
	entry, ok := configKeyRegistry["tui.toast.duration_ms"]
	if !ok {
		t.Fatal("tui.toast.duration_ms not in config registry")
	}
	if entry.Default != "3000" {
		t.Errorf("expected default '3000', got %q", entry.Default)
	}

	cfg := &VaultConfig{}
	if got := entry.Get(cfg); got != "" {
		t.Errorf("expected empty for unset, got %q", got)
	}
	entry.Set(cfg, "5000")
	if cfg.TUI.Toast.DurationMs != 5000 {
		t.Errorf("expected 5000, got %d", cfg.TUI.Toast.DurationMs)
	}
	if got := entry.Get(cfg); got != "5000" {
		t.Errorf("expected '5000' from getter, got %q", got)
	}
}

func TestToastConfigRegistry_DismissKey(t *testing.T) {
	entry, ok := configKeyRegistry["tui.toast.dismiss_key"]
	if !ok {
		t.Fatal("tui.toast.dismiss_key not in config registry")
	}
	if entry.Default != "esc" {
		t.Errorf("expected default 'esc', got %q", entry.Default)
	}

	cfg := &VaultConfig{}
	if got := entry.Get(cfg); got != "" {
		t.Errorf("expected empty for unset, got %q", got)
	}
	entry.Set(cfg, "q")
	if cfg.TUI.Toast.DismissKey != "q" {
		t.Errorf("expected 'q', got %q", cfg.TUI.Toast.DismissKey)
	}
	if got := entry.Get(cfg); got != "q" {
		t.Errorf("expected 'q' from getter, got %q", got)
	}
}

func TestToastConfigRegistry_ShowWarnings(t *testing.T) {
	entry, ok := configKeyRegistry["tui.toast.show_warnings"]
	if !ok {
		t.Fatal("tui.toast.show_warnings not in config registry")
	}
	if entry.Default != "true" {
		t.Errorf("expected default 'true', got %q", entry.Default)
	}

	cfg := &VaultConfig{}

	// nil → empty string
	if got := entry.Get(cfg); got != "" {
		t.Errorf("expected empty for nil, got %q", got)
	}

	// Set true
	entry.Set(cfg, "true")
	if cfg.TUI.Toast.ShowWarnings == nil {
		t.Fatal("expected non-nil after set true")
	}
	if *cfg.TUI.Toast.ShowWarnings != true {
		t.Errorf("expected true, got %v", *cfg.TUI.Toast.ShowWarnings)
	}
	if got := entry.Get(cfg); got != "true" {
		t.Errorf("expected 'true', got %q", got)
	}

	// Set false
	entry.Set(cfg, "false")
	if cfg.TUI.Toast.ShowWarnings == nil {
		t.Fatal("expected non-nil after set false")
	}
	if *cfg.TUI.Toast.ShowWarnings != false {
		t.Errorf("expected false, got %v", *cfg.TUI.Toast.ShowWarnings)
	}
	if got := entry.Get(cfg); got != "false" {
		t.Errorf("expected 'false', got %q", got)
	}
}

func TestToastConfigRegistry_ShowSuccess(t *testing.T) {
	entry, ok := configKeyRegistry["tui.toast.show_success"]
	if !ok {
		t.Fatal("tui.toast.show_success not in config registry")
	}
	if entry.Default != "false" {
		t.Errorf("expected default 'false', got %q", entry.Default)
	}

	cfg := &VaultConfig{}

	// nil → empty string
	if got := entry.Get(cfg); got != "" {
		t.Errorf("expected empty for nil, got %q", got)
	}

	// Set true
	entry.Set(cfg, "true")
	if cfg.TUI.Toast.ShowSuccess == nil {
		t.Fatal("expected non-nil after set true")
	}
	if *cfg.TUI.Toast.ShowSuccess != true {
		t.Errorf("expected true, got %v", *cfg.TUI.Toast.ShowSuccess)
	}
	if got := entry.Get(cfg); got != "true" {
		t.Errorf("expected 'true', got %q", got)
	}

	// Set false
	entry.Set(cfg, "false")
	if cfg.TUI.Toast.ShowSuccess == nil {
		t.Fatal("expected non-nil after set false")
	}
	if *cfg.TUI.Toast.ShowSuccess != false {
		t.Errorf("expected false, got %v", *cfg.TUI.Toast.ShowSuccess)
	}
	if got := entry.Get(cfg); got != "false" {
		t.Errorf("expected 'false', got %q", got)
	}
}

func TestToastConfigKeys_AllRegistered(t *testing.T) {
	expected := []string{
		"tui.toast.position",
		"tui.toast.duration_ms",
		"tui.toast.dismiss_key",
		"tui.toast.show_warnings",
		"tui.toast.show_success",
	}
	keys := ConfigKeys()
	for _, exp := range expected {
		found := false
		for _, k := range keys {
			if k == exp {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected config key %q to be registered, but not found in %v", exp, keys)
		}
	}
}

func TestToastConfigDefaults(t *testing.T) {
	tests := []struct {
		key     string
		want    string
	}{
		{"tui.toast.position", "bottom-right"},
		{"tui.toast.duration_ms", "3000"},
		{"tui.toast.dismiss_key", "esc"},
		{"tui.toast.show_warnings", "true"},
		{"tui.toast.show_success", "false"},
	}
	for _, tt := range tests {
		got := ConfigKeyDefault(tt.key)
		if got != tt.want {
			t.Errorf("ConfigKeyDefault(%q) = %q, want %q", tt.key, got, tt.want)
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
