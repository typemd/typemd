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
	}
	entries = append(entries,
		helpEntry{keys.Search.Help().Key, keys.Search.Help().Desc},
		helpEntry{keys.GrowPanel.Help().Key, keys.GrowPanel.Help().Desc},
		helpEntry{keys.ShrinkPanel.Help().Key, keys.ShrinkPanel.Help().Desc},
		helpEntry{keys.ToggleProps.Help().Key, keys.ToggleProps.Help().Desc},
		helpEntry{keys.ToggleWrap.Help().Key, keys.ToggleWrap.Help().Desc},
		helpEntry{keys.Help.Help().Key, keys.Help.Help().Desc},
		helpEntry{keys.Quit.Help().Key, keys.Quit.Help().Desc},
	)
	return entries
}

// helpContent builds the help popup text content.
func helpContent(width int, readOnly bool) (content string, popupW int) {
	entries := helpEntries(readOnly)

	// Find max key width for alignment
	maxKeyW := 0
	for _, e := range entries {
		if len(e.Key) > maxKeyW {
			maxKeyW = len(e.Key)
		}
	}

	contentW := maxKeyW + helpKeyPadding + helpDescWidth

	// Build lines
	var lines []string
	lines = append(lines, "Keybindings")
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
	content, popupW := helpContent(width, readOnly)
	return renderOverlayPopup(background, content, width, height, popupW)
}
