package tui

import (
	"fmt"
	"strings"
)

const (
	helpDescWidth    = 20 // width reserved for description column
	helpKeyPadding   = 4  // spacing between key and description columns
	helpPopupPadding = 4  // horizontal padding from popup border + Padding(1,2)
)

// helpEntry represents a single keybinding shown in the help popup.
type helpEntry struct {
	Key  string
	Desc string
}

// helpEntries returns the list of keybindings to display in the help popup
// using the caller-supplied keyMap. When readOnly is true, edit-related
// keybindings are hidden.
//
// Callers pass the model's resolved keyMap so the help popup reflects user
// overrides from `.typemd/config.yaml` under `tui.keybindings`.
func helpEntries(km keyMap, readOnly bool) []helpEntry {
	entries := []helpEntry{
		{km.Up.Help().Key, km.Up.Help().Desc},
		{km.Down.Help().Key, km.Down.Help().Desc},
		{km.Enter.Help().Key, km.Enter.Help().Desc},
		{km.Tab.Help().Key, km.Tab.Help().Desc},
	}
	if !readOnly {
		entries = append(entries, helpEntry{km.EnterEdit.Help().Key, km.EnterEdit.Help().Desc})
		entries = append(entries, helpEntry{km.Rename.Help().Key, km.Rename.Help().Desc})
		entries = append(entries, helpEntry{km.NewObject.Help().Key, km.NewObject.Help().Desc})
		entries = append(entries, helpEntry{km.QuickCreate.Help().Key, km.QuickCreate.Help().Desc})
		entries = append(entries, helpEntry{km.AIGenerate.Help().Key, km.AIGenerate.Help().Desc})
		entries = append(entries, helpEntry{km.SchemaExplore.Help().Key, km.SchemaExplore.Help().Desc})
	}
	entries = append(entries,
		helpEntry{km.Search.Help().Key, km.Search.Help().Desc},
		helpEntry{km.Stats.Help().Key, km.Stats.Help().Desc},
		helpEntry{km.Settings.Help().Key, km.Settings.Help().Desc},
		helpEntry{km.GrowPanel.Help().Key, km.GrowPanel.Help().Desc},
		helpEntry{km.ShrinkPanel.Help().Key, km.ShrinkPanel.Help().Desc},
		helpEntry{km.FocusMode.Help().Key, km.FocusMode.Help().Desc},
		helpEntry{km.ToggleProps.Help().Key, km.ToggleProps.Help().Desc},
		helpEntry{km.ToggleWrap.Help().Key, km.ToggleWrap.Help().Desc},
		helpEntry{km.Help.Help().Key, km.Help.Help().Desc},
		helpEntry{km.Quit.Help().Key, km.Quit.Help().Desc},
	)
	return entries
}

// statsHelpEntries returns keybindings specific to stats mode using the
// caller-supplied keyMap.
func statsHelpEntries(km keyMap, screen statsScreen) []helpEntry {
	entries := []helpEntry{
		{km.Up.Help().Key, "up"},
		{km.Down.Help().Key, "down"},
	}
	switch screen {
	case statsOverview:
		entries = append(entries,
			helpEntry{"enter", "type detail"},
			helpEntry{"r", "refresh"},
			helpEntry{"esc", "exit stats"},
		)
	case statsDetail:
		entries = append(entries,
			helpEntry{"t", "toggle layout"},
			helpEntry{"r", "refresh"},
			helpEntry{"esc", "back to overview"},
		)
	}
	entries = append(entries,
		helpEntry{km.Help.Help().Key, km.Help.Help().Desc},
		helpEntry{km.Quit.Help().Key, km.Quit.Help().Desc},
	)
	return entries
}

// buildHelpPopup formats help entries into a popup content string and width.
func buildHelpPopup(title string, entries []helpEntry, width int) (content string, popupW int) {
	maxKeyW := 0
	for _, e := range entries {
		if len(e.Key) > maxKeyW {
			maxKeyW = len(e.Key)
		}
	}

	contentW := maxKeyW + helpKeyPadding + helpDescWidth

	var lines []string
	lines = append(lines, title)
	lines = append(lines, strings.Repeat("─", contentW))
	for _, e := range entries {
		lines = append(lines, fmt.Sprintf("  %-*s   %s", maxKeyW, e.Key, e.Desc))
	}
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("Press Esc or %s to close", keys.Help.Help().Key))

	popupW = contentW + helpPopupPadding
	if popupW > width-4 {
		popupW = width - 4
	}
	popupW += 2 // lipgloss v2 Width includes border

	return strings.Join(lines, "\n"), popupW
}

// renderHelp builds the help overlay on top of a background screen
// using lipgloss Layer/Compositor. The background remains visible outside the popup.
// Uses the caller-supplied keyMap so the overlay reflects user overrides.
func renderHelp(km keyMap, background string, width, height int, readOnly bool) string {
	content, popupW := buildHelpPopup("Keybindings", helpEntries(km, readOnly), width)
	return renderOverlayPopup(background, content, width, height, popupW)
}

// renderStatsHelp builds the help overlay for stats mode.
func renderStatsHelp(km keyMap, background string, width, height int, screen statsScreen) string {
	content, popupW := buildHelpPopup("Stats Keybindings", statsHelpEntries(km, screen), width)
	return renderOverlayPopup(background, content, width, height, popupW)
}

// configHelpEntries returns keybindings specific to config settings mode
// using the caller-supplied keyMap.
func configHelpEntries(km keyMap) []helpEntry {
	return []helpEntry{
		{km.Up.Help().Key, "up"},
		{km.Down.Help().Key, "down"},
		{"tab", "switch column"},
		{"enter", "edit setting"},
		{"esc", "exit settings"},
		{km.Help.Help().Key, km.Help.Help().Desc},
		{km.Quit.Help().Key, km.Quit.Help().Desc},
	}
}

// renderConfigHelp builds the help overlay for config settings mode.
func renderConfigHelp(km keyMap, background string, width, height int) string {
	content, popupW := buildHelpPopup("Settings Keybindings", configHelpEntries(km), width)
	return renderOverlayPopup(background, content, width, height, popupW)
}
