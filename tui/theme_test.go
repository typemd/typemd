package tui

import (
	"testing"

	"charm.land/lipgloss/v2"

	"github.com/typemd/typemd/core"
)

func TestLoadTheme_Defaults(t *testing.T) {
	resetThemeDefaults()
	t.Cleanup(resetThemeDefaults)

	loadTheme(nil)

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

	cfg := &core.ThemeConfig{
		FocusBorder: "196",
		WikiLink:    "82",
	}
	loadTheme(cfg)

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

	cfg := &core.ThemeConfig{
		WikiLink: "214",
	}
	loadTheme(cfg)

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

	loadTheme(nil)

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

	cfg := &core.ThemeConfig{
		Heading:    "196",
		InlineCode: "82",
		CodeBlock:  "120",
		Link:       "45",
		Blockquote: "241",
		HRule:      "240",
	}
	loadTheme(cfg)

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
