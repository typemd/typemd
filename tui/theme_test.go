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

func TestLoadTheme_MarkdownDefaults(t *testing.T) {
	resetThemeDefaults()
	t.Cleanup(resetThemeDefaults)

	loadTheme(t.TempDir())

	// Verify markdown styles use defaults when no config is present.
	got := mdHeadingStyle.GetForeground()
	want := lipgloss.Color(defaultColorHeading)
	if got != want {
		t.Errorf("mdHeadingStyle foreground = %v, want %v", got, want)
	}
	got = mdInlineCodeStyle.GetForeground()
	want = lipgloss.Color(defaultColorInlineCode)
	if got != want {
		t.Errorf("mdInlineCodeStyle foreground = %v, want %v", got, want)
	}
	got = mdCodeBlockStyle.GetForeground()
	want = lipgloss.Color(defaultColorCodeBlock)
	if got != want {
		t.Errorf("mdCodeBlockStyle foreground = %v, want %v", got, want)
	}
	got = mdLinkStyle.GetForeground()
	want = lipgloss.Color(defaultColorLink)
	if got != want {
		t.Errorf("mdLinkStyle foreground = %v, want %v", got, want)
	}
}

func TestLoadTheme_MarkdownOverride(t *testing.T) {
	resetThemeDefaults()
	t.Cleanup(resetThemeDefaults)

	root := t.TempDir()
	dir := filepath.Join(root, ".typemd")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	yamlData := []byte("theme:\n  heading: \"196\"\n  inline_code: \"82\"\n  code_block: \"120\"\n  link: \"45\"\n  blockquote: \"241\"\n  hrule: \"240\"\n")
	if err := os.WriteFile(filepath.Join(dir, "tui.yaml"), yamlData, 0o644); err != nil {
		t.Fatal(err)
	}

	loadTheme(root)

	tests := []struct {
		name string
		got  lipgloss.Style
		want string
	}{
		{"heading", mdHeadingStyle, "196"},
		{"inline_code", mdInlineCodeStyle, "82"},
		{"code_block", mdCodeBlockStyle, "120"},
		{"link", mdLinkStyle, "45"},
		{"blockquote", mdBlockquoteStyle, "241"},
		{"hrule", mdHRuleStyle, "240"},
	}
	for _, tt := range tests {
		if tt.got.GetForeground() != lipgloss.Color(tt.want) {
			t.Errorf("%s foreground = %v, want %s", tt.name, tt.got.GetForeground(), tt.want)
		}
	}
}

func TestResetThemeDefaults_MarkdownStyles(t *testing.T) {
	// Override styles, then reset and verify defaults are restored.
	mdHeadingStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	mdInlineCodeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("82"))

	resetThemeDefaults()

	if mdHeadingStyle.GetForeground() != lipgloss.Color(defaultColorHeading) {
		t.Errorf("after reset, mdHeadingStyle foreground = %v, want %s", mdHeadingStyle.GetForeground(), defaultColorHeading)
	}
	if mdInlineCodeStyle.GetForeground() != lipgloss.Color(defaultColorInlineCode) {
		t.Errorf("after reset, mdInlineCodeStyle foreground = %v, want %s", mdInlineCodeStyle.GetForeground(), defaultColorInlineCode)
	}
}
