package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
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
func helpEntries() []helpEntry {
	return []helpEntry{
		{keys.Up.Help().Key, keys.Up.Help().Desc},
		{keys.Down.Help().Key, keys.Down.Help().Desc},
		{keys.Enter.Help().Key, keys.Enter.Help().Desc},
		{keys.Tab.Help().Key, keys.Tab.Help().Desc},
		{keys.EnterEdit.Help().Key, keys.EnterEdit.Help().Desc},
		{keys.Search.Help().Key, keys.Search.Help().Desc},
		{keys.GrowPanel.Help().Key, keys.GrowPanel.Help().Desc},
		{keys.ShrinkPanel.Help().Key, keys.ShrinkPanel.Help().Desc},
		{keys.ToggleProps.Help().Key, keys.ToggleProps.Help().Desc},
		{keys.ToggleWrap.Help().Key, keys.ToggleWrap.Help().Desc},
		{keys.Help.Help().Key, keys.Help.Help().Desc},
		{keys.Quit.Help().Key, keys.Quit.Help().Desc},
	}
}

// renderHelp builds the help overlay popup content.
func renderHelp(width, height int) string {
	entries := helpEntries()

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

	content := strings.Join(lines, "\n")

	popupW := contentW + helpPopupPadding
	if popupW > width-4 {
		popupW = width - 4
	}

	popup := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorFocusBorder).
		Width(popupW).
		Padding(1, 2).
		Render(content)

	// Center the popup
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, popup)
}
