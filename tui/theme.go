package tui

import (
	"image/color"

	"charm.land/lipgloss/v2"

	"github.com/typemd/typemd/core"
)

// Default color values.
const (
	defaultColorFocusBorder = "63"
	defaultColorEditBorder  = "214"
	defaultColorWikiLink    = "33"
	defaultColorHeading     = "3"
	defaultColorInlineCode  = "245"
	defaultColorCodeBlock   = "245"
	defaultColorLink        = "33"
	defaultColorBlockquote  = "8"
	defaultColorHRule       = "8"
)

// Theme colors and styles.
var (
	colorFocusBorder  color.Color = lipgloss.Color(defaultColorFocusBorder)
	colorEditBorder   color.Color = lipgloss.Color(defaultColorEditBorder)
	colorWikiLink     color.Color = lipgloss.Color(defaultColorWikiLink)
	wikiLinkStyleBase             = lipgloss.NewStyle().Foreground(colorWikiLink)
	highlightStyle                = lipgloss.NewStyle().Bold(true).Reverse(true)
	dimStyle                      = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	activeStyle                   = lipgloss.NewStyle().Foreground(lipgloss.Color("0")).Background(lipgloss.Color("6"))
	boldStyle                     = lipgloss.NewStyle().Bold(true)

	// Markdown element styles.
	mdHeadingStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color(defaultColorHeading)).Bold(true)
	mdBoldStyle       = lipgloss.NewStyle().Bold(true)
	mdItalicStyle     = lipgloss.NewStyle().Italic(true)
	mdInlineCodeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(defaultColorInlineCode))
	mdCodeBlockStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color(defaultColorCodeBlock))
	mdLinkStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color(defaultColorLink))
	mdBlockquoteStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(defaultColorBlockquote))
	mdHRuleStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color(defaultColorHRule))
)

// loadTheme applies theme overrides from VaultConfig. Missing fields are
// silently ignored, keeping the defaults.
func loadTheme(cfg *core.ThemeConfig) {
	if cfg == nil {
		return
	}

	if cfg.FocusBorder != "" {
		colorFocusBorder = lipgloss.Color(cfg.FocusBorder)
	}
	if cfg.WikiLink != "" {
		colorWikiLink = lipgloss.Color(cfg.WikiLink)
	}
	wikiLinkStyleBase = lipgloss.NewStyle().Foreground(colorWikiLink)

	if cfg.Heading != "" {
		mdHeadingStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(cfg.Heading)).Bold(true)
	}
	if cfg.Bold != "" {
		mdBoldStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(cfg.Bold)).Bold(true)
	}
	if cfg.Italic != "" {
		mdItalicStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(cfg.Italic)).Italic(true)
	}
	if cfg.InlineCode != "" {
		mdInlineCodeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(cfg.InlineCode))
	}
	if cfg.CodeBlock != "" {
		mdCodeBlockStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(cfg.CodeBlock))
	}
	if cfg.Link != "" {
		mdLinkStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(cfg.Link))
	}
	if cfg.Blockquote != "" {
		mdBlockquoteStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(cfg.Blockquote))
	}
	if cfg.HRule != "" {
		mdHRuleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(cfg.HRule))
	}
}

// resetThemeDefaults restores all theme state to defaults. Used by tests.
func resetThemeDefaults() {
	colorFocusBorder = lipgloss.Color(defaultColorFocusBorder)
	colorEditBorder = lipgloss.Color(defaultColorEditBorder)
	colorWikiLink = lipgloss.Color(defaultColorWikiLink)
	wikiLinkStyleBase = lipgloss.NewStyle().Foreground(colorWikiLink)

	mdHeadingStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(defaultColorHeading)).Bold(true)
	mdBoldStyle = lipgloss.NewStyle().Bold(true)
	mdItalicStyle = lipgloss.NewStyle().Italic(true)
	mdInlineCodeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(defaultColorInlineCode))
	mdCodeBlockStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(defaultColorCodeBlock))
	mdLinkStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(defaultColorLink))
	mdBlockquoteStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(defaultColorBlockquote))
	mdHRuleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(defaultColorHRule))
}
