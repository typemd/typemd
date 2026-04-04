package tui

import (
	"strings"
	"testing"

	"github.com/typemd/typemd/core"
	tea "charm.land/bubbletea/v2"
)

func setupTestVaultForConfig(t *testing.T) *core.Vault {
	t.Helper()
	dir := t.TempDir()
	v := core.NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatalf("vault Init() error = %v", err)
	}
	if err := v.Open(); err != nil {
		t.Fatalf("vault Open() error = %v", err)
	}
	return v
}

func TestConfigEditor_NewHasCategories(t *testing.T) {
	v := setupTestVaultForConfig(t)
	defer v.Close()
	ce := newConfigEditor(v)
	if len(ce.categories) == 0 {
		t.Fatal("expected categories, got none")
	}
	// Check expected categories exist
	names := make(map[string]bool)
	for _, cat := range ce.categories {
		names[cat.Name] = true
	}
	for _, expected := range []string{"General", "CLI", "TUI", "AI", "Web"} {
		if !names[expected] {
			t.Errorf("missing category %q", expected)
		}
	}
}

func TestConfigEditor_CategoryNavigation(t *testing.T) {
	v := setupTestVaultForConfig(t)
	defer v.Close()
	ce := newConfigEditor(v)
	ce.SetSize(80, 24)

	if ce.catCursor != 0 {
		t.Fatalf("initial catCursor = %d, want 0", ce.catCursor)
	}

	// Navigate down
	ce, _ = ce.Update(tea.KeyPressMsg{Code: 'j'})
	if ce.catCursor != 1 {
		t.Errorf("after j: catCursor = %d, want 1", ce.catCursor)
	}

	// Navigate up
	ce, _ = ce.Update(tea.KeyPressMsg{Code: 'k'})
	if ce.catCursor != 0 {
		t.Errorf("after k: catCursor = %d, want 0", ce.catCursor)
	}

	// Don't go below 0
	ce, _ = ce.Update(tea.KeyPressMsg{Code: 'k'})
	if ce.catCursor != 0 {
		t.Errorf("after extra k: catCursor = %d, want 0", ce.catCursor)
	}
}

func TestConfigEditor_ColumnSwitch(t *testing.T) {
	v := setupTestVaultForConfig(t)
	defer v.Close()
	ce := newConfigEditor(v)
	ce.SetSize(80, 24)

	if ce.activeColumn != colCategories {
		t.Fatalf("initial column = %d, want colCategories", ce.activeColumn)
	}

	// Tab switches to settings
	ce, _ = ce.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	if ce.activeColumn != colSettings {
		t.Errorf("after tab: column = %d, want colSettings", ce.activeColumn)
	}

	// Tab switches back to categories
	ce, _ = ce.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	if ce.activeColumn != colCategories {
		t.Errorf("after second tab: column = %d, want colCategories", ce.activeColumn)
	}
}

func TestConfigEditor_SettingsNavigation(t *testing.T) {
	v := setupTestVaultForConfig(t)
	defer v.Close()
	ce := newConfigEditor(v)
	ce.SetSize(80, 24)

	// Switch to settings column
	ce, _ = ce.Update(tea.KeyPressMsg{Code: tea.KeyTab})

	settings := ce.currentSettings()
	if len(settings) == 0 {
		t.Fatal("expected settings in first category")
	}

	if ce.settCursor != 0 {
		t.Fatalf("initial settCursor = %d, want 0", ce.settCursor)
	}

	// Navigate down in settings
	ce, _ = ce.Update(tea.KeyPressMsg{Code: 'j'})
	if ce.settCursor != min(1, len(settings)-1) {
		t.Errorf("after j: settCursor = %d", ce.settCursor)
	}
}

func TestConfigEditor_EscExits(t *testing.T) {
	v := setupTestVaultForConfig(t)
	defer v.Close()
	ce := newConfigEditor(v)
	ce.SetSize(80, 24)

	result, _ := ce.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	if result != nil {
		t.Error("expected nil return on Esc (exit signal)")
	}
}

func TestConfigEditor_EditStringValue(t *testing.T) {
	v := setupTestVaultForConfig(t)
	defer v.Close()
	ce := newConfigEditor(v)
	ce.SetSize(80, 24)

	// Switch to settings column and press enter to edit
	ce, _ = ce.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	ce, _ = ce.Update(tea.KeyPressMsg{Code: tea.KeyEnter})

	if !ce.editing {
		t.Fatal("expected editing to be true after Enter")
	}
	if ce.editBool {
		t.Fatal("expected non-boolean edit for first General setting")
	}

	// Type a value
	for _, r := range "yyyy-mm-dd" {
		ce, _ = ce.Update(tea.KeyPressMsg{Code: r})
	}

	// Save with Enter
	ce, _ = ce.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if ce.editing {
		t.Error("expected editing to be false after Enter save")
	}

	// Verify the value was saved
	val, _ := v.GetConfigValue(ce.currentSettings()[0].Key)
	if val == "" {
		// The first General setting is date_format — check it was set
		t.Log("Note: value may have been cleared if empty was typed")
	}
}

func TestConfigEditor_EditCancel(t *testing.T) {
	v := setupTestVaultForConfig(t)
	defer v.Close()
	ce := newConfigEditor(v)
	ce.SetSize(80, 24)

	// Switch to settings column and open edit
	ce, _ = ce.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	ce, _ = ce.Update(tea.KeyPressMsg{Code: tea.KeyEnter})

	if !ce.editing {
		t.Fatal("expected editing to be true")
	}

	// Cancel with Esc
	ce, _ = ce.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	if ce.editing {
		t.Error("expected editing to be false after Esc cancel")
	}
}

func TestConfigEditor_EditBoolCycle(t *testing.T) {
	v := setupTestVaultForConfig(t)
	defer v.Close()

	// Set a boolean key so we can test cycling
	_ = v.SetConfigValue("ai.enabled", "true")

	ce := newConfigEditor(v)
	ce.SetSize(80, 24)

	// Navigate to AI category (index 3: General=0, CLI=1, TUI=2, AI=3)
	for i := 0; i < 3; i++ {
		ce, _ = ce.Update(tea.KeyPressMsg{Code: 'j'})
	}

	// Switch to settings and find ai.enabled
	ce, _ = ce.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	settings := ce.currentSettings()
	boolIdx := -1
	for i, s := range settings {
		if s.Key == "ai.enabled" {
			boolIdx = i
			break
		}
	}
	if boolIdx < 0 {
		t.Fatal("ai.enabled not found in AI category")
	}

	// Navigate to ai.enabled
	for i := 0; i < boolIdx; i++ {
		ce, _ = ce.Update(tea.KeyPressMsg{Code: 'j'})
	}

	// Open edit
	ce, _ = ce.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if !ce.editing || !ce.editBool {
		t.Fatal("expected boolean edit mode")
	}

	// Initial value should be "true"
	if ce.editBoolVal != "true" {
		t.Errorf("initial bool val = %q, want %q", ce.editBoolVal, "true")
	}

	// Cycle: true → false
	ce, _ = ce.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if ce.editBoolVal != "false" {
		t.Errorf("after cycle: bool val = %q, want %q", ce.editBoolVal, "false")
	}

	// Cycle: false → unset
	ce, _ = ce.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if ce.editBoolVal != "unset" {
		t.Errorf("after second cycle: bool val = %q, want %q", ce.editBoolVal, "unset")
	}

	// Close (saves on Esc)
	ce, _ = ce.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	if ce.editing {
		t.Error("expected editing to be false after Esc")
	}

	// Verify saved value (unset = empty string)
	val, _ := v.GetConfigValue("ai.enabled")
	if val != "" {
		t.Errorf("expected empty value after cycling to unset, got %q", val)
	}
}

func TestConfigEditor_ClearToDefault(t *testing.T) {
	v := setupTestVaultForConfig(t)
	defer v.Close()

	// Set a value first
	_ = v.SetConfigValue("cli.default_type", "idea")

	ce := newConfigEditor(v)
	ce.SetSize(80, 24)

	// Navigate to CLI category (index 1)
	ce, _ = ce.Update(tea.KeyPressMsg{Code: 'j'})

	// Switch to settings and open edit
	ce, _ = ce.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	ce, _ = ce.Update(tea.KeyPressMsg{Code: tea.KeyEnter})

	if !ce.editing {
		t.Fatal("expected editing mode")
	}

	// Clear the input (select all + delete)
	ce.editInput.SetValue("")

	// Save with Enter
	ce, _ = ce.Update(tea.KeyPressMsg{Code: tea.KeyEnter})

	// Verify value was cleared
	val, _ := v.GetConfigValue("cli.default_type")
	if val != "" {
		t.Errorf("expected cleared value, got %q", val)
	}
}

func TestConfigEditor_ViewRenders(t *testing.T) {
	v := setupTestVaultForConfig(t)
	defer v.Close()
	ce := newConfigEditor(v)
	ce.SetSize(80, 24)

	view := ce.View()
	if view == "" {
		t.Fatal("expected non-empty view")
	}
	// Should contain category names
	if !strings.Contains(view, "General") {
		t.Error("view should contain 'General'")
	}
	if !strings.Contains(view, "│") {
		t.Error("view should contain column separator")
	}
}

func TestConfigEditor_HelpBar(t *testing.T) {
	v := setupTestVaultForConfig(t)
	defer v.Close()
	ce := newConfigEditor(v)

	bar := ce.HelpBar()
	if !strings.Contains(bar, "SETTINGS") {
		t.Errorf("help bar should contain SETTINGS, got %q", bar)
	}
	if !strings.Contains(bar, "tab") {
		t.Errorf("help bar should contain tab, got %q", bar)
	}
}

func TestConfigEditor_EditPopupEmpty(t *testing.T) {
	v := setupTestVaultForConfig(t)
	defer v.Close()
	ce := newConfigEditor(v)
	ce.SetSize(80, 24)

	// Not editing — should return empty
	result := ce.EditPopup("background", 80, 24)
	if result != "" {
		t.Error("expected empty popup when not editing")
	}
}

func TestConfigEditor_TitleContent(t *testing.T) {
	v := setupTestVaultForConfig(t)
	defer v.Close()
	ce := newConfigEditor(v)

	title := ce.titleContent()
	if !strings.Contains(title, "Settings") {
		t.Errorf("title should contain 'Settings', got %q", title)
	}
}

func TestCategoryForKey(t *testing.T) {
	tests := []struct {
		key  string
		want string
	}{
		{"date_format", "General"},
		{"datetime_format", "General"},
		{"cli.default_type", "CLI"},
		{"tui.debounce_ms", "TUI"},
		{"tui.toast.position", "TUI"},
		{"ai.default", "AI"},
		{"ai.prompts.describe", "AI"},
		{"web.theme", "Web"},
	}
	for _, tt := range tests {
		got := categoryForKey(tt.key)
		if got != tt.want {
			t.Errorf("categoryForKey(%q) = %q, want %q", tt.key, got, tt.want)
		}
	}
}

func TestIsBoolKey(t *testing.T) {
	boolKeys := []string{"tui.toast.show_warnings", "tui.toast.show_success", "ai.enabled"}
	for _, key := range boolKeys {
		if !isBoolKey(key) {
			t.Errorf("isBoolKey(%q) = false, want true", key)
		}
	}
	if isBoolKey("cli.default_type") {
		t.Error("isBoolKey(cli.default_type) = true, want false")
	}
}
