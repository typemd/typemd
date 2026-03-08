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
)

// Theme colors and styles.
var (
	colorFocusBorder  color.Color = lipgloss.Color(defaultColorFocusBorder)
	colorEditBorder   color.Color = lipgloss.Color(defaultColorEditBorder)
	colorWikiLink     color.Color = lipgloss.Color(defaultColorWikiLink)
	wikiLinkStyleBase             = lipgloss.NewStyle().Foreground(colorWikiLink)
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
}

// resetThemeDefaults restores all theme state to defaults. Used by tests.
func resetThemeDefaults() {
	colorFocusBorder = lipgloss.Color(defaultColorFocusBorder)
	colorEditBorder = lipgloss.Color(defaultColorEditBorder)
	colorWikiLink = lipgloss.Color(defaultColorWikiLink)
	wikiLinkStyleBase = lipgloss.NewStyle().Foreground(colorWikiLink)
}
