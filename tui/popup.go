package tui

import "charm.land/lipgloss/v2"

// renderPopup renders content as a centered popup overlay with a rounded border.
// The border color uses the theme's colorFocusBorder. An optional width can be
// provided to constrain the popup; pass 0 to let lipgloss auto-size from content.
func renderPopup(content string, termW, termH int, popupWidth int) string {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorFocusBorder).
		Padding(1, 2)

	if popupWidth > 0 {
		style = style.Width(popupWidth)
	}

	popup := style.Render(content)
	return lipgloss.Place(termW, termH, lipgloss.Center, lipgloss.Center, popup,
		lipgloss.WithWhitespaceChars(" "),
	)
}
