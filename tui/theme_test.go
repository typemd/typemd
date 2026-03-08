package tui

import (
	"os"
	"path/filepath"
	"testing"

	"charm.land/lipgloss/v2"
)

func TestLoadTheme_Defaults(t *testing.T) {
	resetThemeDefaults()
	t.Cleanup(resetThemeDefaults)

	loadTheme(t.TempDir())

	if colorFocusBorder != lipgloss.Color(defaultColorFocusBorder) {
		t.Errorf("colorFocusBorder = %v, want %s", colorFocusBorder, defaultColorFocusBorder)
	}
	if colorWikiLink != lipgloss.Color(defaultColorWikiLink) {
		t.Errorf("colorWikiLink = %v, want %s", colorWikiLink, defaultColorWikiLink)
	}
}

func TestLoadTheme_Override(t *testing.T) {
	resetThemeDefaults()
	t.Cleanup(resetThemeDefaults)

	root := t.TempDir()
	dir := filepath.Join(root, ".typemd")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	yaml := []byte("theme:\n  focus_border: \"196\"\n  wiki_link: \"82\"\n")
	if err := os.WriteFile(filepath.Join(dir, "tui.yaml"), yaml, 0o644); err != nil {
		t.Fatal(err)
	}

	loadTheme(root)

	if colorFocusBorder != lipgloss.Color("196") {
		t.Errorf("colorFocusBorder = %v, want 196", colorFocusBorder)
	}
	if colorWikiLink != lipgloss.Color("82") {
		t.Errorf("colorWikiLink = %v, want 82", colorWikiLink)
	}
}

func TestLoadTheme_PartialOverride(t *testing.T) {
	resetThemeDefaults()
	t.Cleanup(resetThemeDefaults)

	root := t.TempDir()
	dir := filepath.Join(root, ".typemd")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	yaml := []byte("theme:\n  wiki_link: \"214\"\n")
	if err := os.WriteFile(filepath.Join(dir, "tui.yaml"), yaml, 0o644); err != nil {
		t.Fatal(err)
	}

	loadTheme(root)

	if colorFocusBorder != lipgloss.Color(defaultColorFocusBorder) {
		t.Errorf("colorFocusBorder = %v, want %s (unchanged)", colorFocusBorder, defaultColorFocusBorder)
	}
	if colorWikiLink != lipgloss.Color("214") {
		t.Errorf("colorWikiLink = %v, want 214", colorWikiLink)
	}
}
