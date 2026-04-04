package tui

import (
	"fmt"
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/typemd/typemd/core"
	"github.com/typemd/typemd/tui/widget"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

// configColumn tracks which column has focus in the two-column layout.
type configColumn int

const (
	colCategories configColumn = iota
	colSettings
)

// configCategory groups config keys by their dot-notation prefix.
type configCategory struct {
	Name string
	Keys []core.ConfigKeyInfo
}

// configEditor is a sub-model for the full-width config settings page.
type configEditor struct {
	vault *core.Vault

	categories   []configCategory
	catCursor    int
	settCursor   int
	settScroll   int
	activeColumn configColumn

	// Edit popup state
	editing    bool
	editKey    string
	editDesc   string
	editDef    string
	editInput  textinput.Model
	editBool   bool   // true if current key is boolean-typed
	editBoolVal string // "true", "false", or "" (unset)

	width  int
	height int
}

// newConfigEditor creates a config editor and loads current config key info.
func newConfigEditor(vault *core.Vault) *configEditor {
	ce := &configEditor{
		vault: vault,
	}
	ce.loadKeys()
	return ce
}

// categoryForKey returns the category name for a dot-notation config key.
func categoryForKey(key string) string {
	if i := strings.IndexByte(key, '.'); i >= 0 {
		switch key[:i] {
		case "cli":
			return "CLI"
		case "tui":
			return "TUI"
		case "ai":
			return "AI"
		case "web":
			return "Web"
		}
	}
	return "General"
}

// categoryOrder defines the display order of categories.
var categoryOrder = []string{"General", "CLI", "TUI", "AI", "Web"}

// loadKeys fetches config key info from the vault and groups by category.
func (ce *configEditor) loadKeys() {
	infos := ce.vault.ConfigKeysInfo()

	grouped := make(map[string][]core.ConfigKeyInfo)
	for _, info := range infos {
		cat := categoryForKey(info.Key)
		grouped[cat] = append(grouped[cat], info)
	}

	ce.categories = nil
	for _, name := range categoryOrder {
		if keys, ok := grouped[name]; ok {
			ce.categories = append(ce.categories, configCategory{
				Name: name,
				Keys: keys,
			})
		}
	}
}

// currentSettings returns the settings for the currently selected category.
func (ce *configEditor) currentSettings() []core.ConfigKeyInfo {
	if ce.catCursor >= 0 && ce.catCursor < len(ce.categories) {
		return ce.categories[ce.catCursor].Keys
	}
	return nil
}

// SetSize updates the available rendering dimensions.
func (ce *configEditor) SetSize(w, h int) {
	ce.width = w
	ce.height = h
}

// Update handles key events for the config editor.
func (ce *configEditor) Update(msg tea.Msg) (*configEditor, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if ce.editing {
			return ce.updateEdit(msg)
		}
		return ce.updateBrowse(msg)
	}
	return ce, nil
}

func (ce *configEditor) updateBrowse(msg tea.KeyPressMsg) (*configEditor, tea.Cmd) {
	switch msg.String() {
	case "esc":
		return nil, nil // signal exit to caller

	case "tab":
		if ce.activeColumn == colCategories {
			ce.activeColumn = colSettings
		} else {
			ce.activeColumn = colCategories
		}
		return ce, nil

	case "j", "down":
		switch ce.activeColumn {
		case colCategories:
			if ce.catCursor < len(ce.categories)-1 {
				ce.catCursor++
				ce.settCursor = 0
				ce.settScroll = 0
			}
		case colSettings:
			settings := ce.currentSettings()
			if ce.settCursor < len(settings)-1 {
				ce.settCursor++
			}
			ce.settScroll = widget.AdjustScroll(ce.settCursor, ce.settScroll, ce.settingsVisibleH())
		}
		return ce, nil

	case "k", "up":
		switch ce.activeColumn {
		case colCategories:
			if ce.catCursor > 0 {
				ce.catCursor--
				ce.settCursor = 0
				ce.settScroll = 0
			}
		case colSettings:
			if ce.settCursor > 0 {
				ce.settCursor--
			}
			ce.settScroll = widget.AdjustScroll(ce.settCursor, ce.settScroll, ce.settingsVisibleH())
		}
		return ce, nil

	case "enter":
		if ce.activeColumn == colCategories {
			// Switch to settings column on enter
			ce.activeColumn = colSettings
			return ce, nil
		}
		// Open edit popup for current setting
		settings := ce.currentSettings()
		if ce.settCursor >= 0 && ce.settCursor < len(settings) {
			info := settings[ce.settCursor]
			ce.editKey = info.Key
			ce.editDesc = info.Description
			ce.editDef = info.Default

			// Detect boolean keys
			if isBoolKey(info.Key) {
				ce.editBool = true
				ce.editBoolVal = info.Value
				if ce.editBoolVal == "" {
					ce.editBoolVal = "unset"
				}
				ce.editing = true
				return ce, nil
			}

			ce.editBool = false
			ti := textinput.New()
			ti.Prompt = ""
			ti.Placeholder = info.Default
			ti.CharLimit = 200
			if info.Value != "" {
				ti.SetValue(info.Value)
				ti.CursorEnd()
			}
			ti.Focus()
			ce.editInput = ti
			ce.editing = true
			return ce, textinput.Blink
		}
		return ce, nil
	}
	return ce, nil
}

// isBoolKey returns true if the config key uses boolean values.
func isBoolKey(key string) bool {
	switch key {
	case "tui.toast.show_warnings", "tui.toast.show_success", "ai.enabled":
		return true
	}
	return false
}

func (ce *configEditor) updateEdit(msg tea.KeyPressMsg) (*configEditor, tea.Cmd) {
	switch {
	case ce.editBool:
		return ce.updateEditBool(msg)
	default:
		return ce.updateEditText(msg)
	}
}

func (ce *configEditor) updateEditBool(msg tea.KeyPressMsg) (*configEditor, tea.Cmd) {
	switch msg.String() {
	case "esc":
		// Save on close
		value := ce.editBoolVal
		if value == "unset" {
			value = ""
		}
		_ = ce.vault.SetConfigValue(ce.editKey, value)
		ce.editing = false
		ce.loadKeys()
		return ce, nil
	case "enter", "j", "down", " ":
		// Cycle: unset → true → false → unset
		switch ce.editBoolVal {
		case "unset":
			ce.editBoolVal = "true"
		case "true":
			ce.editBoolVal = "false"
		case "false":
			ce.editBoolVal = "unset"
		default:
			ce.editBoolVal = "true"
		}
		return ce, nil
	case "k", "up":
		// Cycle reverse: unset ← true ← false ← unset
		switch ce.editBoolVal {
		case "unset":
			ce.editBoolVal = "false"
		case "true":
			ce.editBoolVal = "unset"
		case "false":
			ce.editBoolVal = "true"
		default:
			ce.editBoolVal = "true"
		}
		return ce, nil
	}
	return ce, nil
}

func (ce *configEditor) updateEditText(msg tea.KeyPressMsg) (*configEditor, tea.Cmd) {
	switch msg.String() {
	case "esc":
		ce.editing = false
		return ce, nil
	case "enter":
		value := strings.TrimSpace(ce.editInput.Value())
		_ = ce.vault.SetConfigValue(ce.editKey, value)
		ce.editing = false
		ce.loadKeys()
		return ce, nil
	}
	var cmd tea.Cmd
	ce.editInput, cmd = ce.editInput.Update(msg)
	return ce, cmd
}

// settingsVisibleH returns the number of visible rows for settings.
func (ce *configEditor) settingsVisibleH() int {
	h := ce.height - 4 // title + spacing
	if h < 1 {
		h = 1
	}
	return h
}

// View renders the config editor content.
func (ce *configEditor) View() string {
	if len(ce.categories) == 0 {
		return "  No config keys found."
	}

	catW := 14 // fixed width for category column
	settW := ce.width - catW - 3 // separator + spacing
	if settW < 20 {
		settW = 20
	}

	// Render categories column
	var catLines []string
	for i, cat := range ce.categories {
		count := len(cat.Keys)
		line := fmt.Sprintf("  %-10s (%d)", cat.Name, count)
		if i == ce.catCursor {
			if ce.activeColumn == colCategories {
				line = highlightStyle.Render(line)
			} else {
				line = boldStyle.Render(line)
			}
		}
		catLines = append(catLines, line)
	}

	// Render settings column
	settings := ce.currentSettings()
	var settLines []string
	visibleH := ce.settingsVisibleH()

	for i, info := range settings {
		if i < ce.settScroll {
			continue
		}
		if i >= ce.settScroll+visibleH {
			break
		}

		displayVal := info.Value
		if displayVal == "" {
			displayVal = dimStyle.Render(info.Default)
		}

		keyDisplay := info.Key
		maxKeyW := settW/2 - 2
		if maxKeyW > 0 && runewidth.StringWidth(keyDisplay) > maxKeyW {
			keyDisplay = truncate(keyDisplay, maxKeyW)
		}

		line := fmt.Sprintf("  %-*s  %s", maxKeyW, keyDisplay, displayVal)

		if i == ce.settCursor && ce.activeColumn == colSettings {
			line = highlightStyle.Render(line)
		}

		settLines = append(settLines, line)

		// Show description on the line below the selected item
		if i == ce.settCursor && ce.activeColumn == colSettings && info.Description != "" {
			descLine := "  " + dimStyle.Render(info.Description)
			settLines = append(settLines, descLine)
		}
	}

	// Pad columns to same height
	maxH := max(len(catLines), len(settLines))
	for len(catLines) < maxH {
		catLines = append(catLines, "")
	}
	for len(settLines) < maxH {
		settLines = append(settLines, "")
	}

	// Join columns with a vertical separator
	var lines []string
	for i := 0; i < maxH; i++ {
		catPart := fmt.Sprintf("%-*s", catW, catLines[i])
		line := catPart + " │ " + settLines[i]
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

// EditPopup renders the edit popup overlay if editing.
// Returns "" if not editing.
func (ce *configEditor) EditPopup(background string, termW, termH int) string {
	if !ce.editing {
		return ""
	}

	var lines []string
	lines = append(lines, boldStyle.Render(ce.editKey))
	if ce.editDesc != "" {
		lines = append(lines, dimStyle.Render(ce.editDesc))
	}
	lines = append(lines, "")

	if ce.editBool {
		lines = append(lines, fmt.Sprintf("Default: %s", dimStyle.Render(ce.editDef)))
		lines = append(lines, "")
		options := []string{"true", "false", "unset"}
		for _, opt := range options {
			prefix := "  "
			if opt == ce.editBoolVal {
				prefix = "▸ "
			}
			label := opt
			if opt == "unset" {
				label = fmt.Sprintf("unset (default: %s)", ce.editDef)
			}
			if opt == ce.editBoolVal {
				lines = append(lines, highlightStyle.Render(prefix+label))
			} else {
				lines = append(lines, "  "+label)
			}
		}
		lines = append(lines, "")
		lines = append(lines, dimStyle.Render("↑↓/enter: cycle  esc: close"))
	} else {
		lines = append(lines, fmt.Sprintf("Default: %s", dimStyle.Render(ce.editDef)))
		lines = append(lines, "")
		lines = append(lines, ce.editInput.View())
		lines = append(lines, "")
		lines = append(lines, dimStyle.Render("enter: save  esc: cancel"))
	}

	content := strings.Join(lines, "\n")
	popupW := 50
	if popupW > termW-4 {
		popupW = termW - 4
	}
	return renderOverlayPopup(background, content, termW, termH, popupW)
}

// HelpBar returns the help bar text for the config editor.
func (ce *configEditor) HelpBar() string {
	if ce.editing {
		if ce.editBool {
			return "  [SETTINGS]  ↑↓/enter: cycle value  |  esc: close"
		}
		return "  [SETTINGS]  enter: save  |  esc: cancel"
	}
	return "  [SETTINGS]  ↑↓: navigate  |  tab: switch column  |  enter: edit  |  esc: back"
}

// titleContent returns the title bar text for the config editor.
func (ce *configEditor) titleContent() string {
	return " ⚙ Settings"
}
