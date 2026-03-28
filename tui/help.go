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

// helpEntries returns the list of keybindings to display in the help popup.
// When readOnly is true, edit-related keybindings are hidden.
func helpEntries(readOnly bool) []helpEntry {
	entries := []helpEntry{
		{keys.Up.Help().Key, keys.Up.Help().Desc},
		{keys.Down.Help().Key, keys.Down.Help().Desc},
		{keys.Enter.Help().Key, keys.Enter.Help().Desc},
		{keys.Tab.Help().Key, keys.Tab.Help().Desc},
	}
	if !readOnly {
		entries = append(entries, helpEntry{keys.EnterEdit.Help().Key, keys.EnterEdit.Help().Desc})
		entries = append(entries, helpEntry{keys.Rename.Help().Key, keys.Rename.Help().Desc})
		entries = append(entries, helpEntry{keys.NewObject.Help().Key, keys.NewObject.Help().Desc})
		entries = append(entries, helpEntry{keys.QuickCreate.Help().Key, keys.QuickCreate.Help().Desc})
		entries = append(entries, helpEntry{keys.AIGenerate.Help().Key, keys.AIGenerate.Help().Desc})
		entries = append(entries, helpEntry{keys.SchemaExplore.Help().Key, keys.SchemaExplore.Help().Desc})
	}
	entries = append(entries,
		helpEntry{keys.Search.Help().Key, keys.Search.Help().Desc},
		helpEntry{keys.Stats.Help().Key, keys.Stats.Help().Desc},
		helpEntry{keys.GrowPanel.Help().Key, keys.GrowPanel.Help().Desc},
		helpEntry{keys.ShrinkPanel.Help().Key, keys.ShrinkPanel.Help().Desc},
		helpEntry{keys.FocusMode.Help().Key, keys.FocusMode.Help().Desc},
		helpEntry{keys.ToggleProps.Help().Key, keys.ToggleProps.Help().Desc},
		helpEntry{keys.ToggleWrap.Help().Key, keys.ToggleWrap.Help().Desc},
		helpEntry{keys.Help.Help().Key, keys.Help.Help().Desc},
		helpEntry{keys.Quit.Help().Key, keys.Quit.Help().Desc},
	)
	return entries
}

// statsHelpEntries returns keybindings specific to stats mode.
func statsHelpEntries(screen statsScreen) []helpEntry {
	entries := []helpEntry{
		{keys.Up.Help().Key, "up"},
		{keys.Down.Help().Key, "down"},
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
		helpEntry{keys.Help.Help().Key, keys.Help.Help().Desc},
		helpEntry{keys.Quit.Help().Key, keys.Quit.Help().Desc},
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
func renderHelp(background string, width, height int, readOnly bool) string {
	content, popupW := buildHelpPopup("Keybindings", helpEntries(readOnly), width)
	return renderOverlayPopup(background, content, width, height, popupW)
}

// renderStatsHelp builds the help overlay for stats mode.
func renderStatsHelp(background string, width, height int, screen statsScreen) string {
	content, popupW := buildHelpPopup("Stats Keybindings", statsHelpEntries(screen), width)
	return renderOverlayPopup(background, content, width, height, popupW)
}
