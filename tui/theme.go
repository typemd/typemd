package tui

import (
	"image/color"
	"os"
	"path/filepath"

	"charm.land/lipgloss/v2"
	"gopkg.in/yaml.v3"
)

// themeConfig represents the theme section in .typemd/tui.yaml.
type themeConfig struct {
	FocusBorder string `yaml:"focus_border"`
	WikiLink    string `yaml:"wiki_link"`
	Heading     string `yaml:"heading"`
	Bold        string `yaml:"bold"`
	Italic      string `yaml:"italic"`
	InlineCode  string `yaml:"inline_code"`
	CodeBlock   string `yaml:"code_block"`
	Link        string `yaml:"link"`
	Blockquote  string `yaml:"blockquote"`
	HRule       string `yaml:"hrule"`
}

// tuiConfig represents the top-level structure of .typemd/tui.yaml.
type tuiConfig struct {
	Theme themeConfig `yaml:"theme"`
}

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

// loadTheme reads .typemd/tui.yaml from the vault root and overrides default
// colors if values are present. Missing file or missing fields are silently
// ignored, keeping the defaults.
func loadTheme(vaultRoot string) {
	data, err := os.ReadFile(filepath.Join(vaultRoot, ".typemd", "tui.yaml"))
	if err != nil {
		return
	}

	var cfg tuiConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return
	}

	if cfg.Theme.FocusBorder != "" {
		colorFocusBorder = lipgloss.Color(cfg.Theme.FocusBorder)
	}
	if cfg.Theme.WikiLink != "" {
		colorWikiLink = lipgloss.Color(cfg.Theme.WikiLink)
	}
	wikiLinkStyleBase = lipgloss.NewStyle().Foreground(colorWikiLink)

	if cfg.Theme.Heading != "" {
		mdHeadingStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(cfg.Theme.Heading)).Bold(true)
	}
	if cfg.Theme.Bold != "" {
		mdBoldStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(cfg.Theme.Bold)).Bold(true)
	}
	if cfg.Theme.Italic != "" {
		mdItalicStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(cfg.Theme.Italic)).Italic(true)
	}
	if cfg.Theme.InlineCode != "" {
		mdInlineCodeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(cfg.Theme.InlineCode))
	}
	if cfg.Theme.CodeBlock != "" {
		mdCodeBlockStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(cfg.Theme.CodeBlock))
	}
	if cfg.Theme.Link != "" {
		mdLinkStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(cfg.Theme.Link))
	}
	if cfg.Theme.Blockquote != "" {
		mdBlockquoteStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(cfg.Theme.Blockquote))
	}
	if cfg.Theme.HRule != "" {
		mdHRuleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(cfg.Theme.HRule))
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
