package widget

import "charm.land/lipgloss/v2"

// CenteredPopup renders content inside a styled popup, centered within termW × termH.
// The caller controls the popup style (border, width, padding, colors).
func CenteredPopup(content string, style lipgloss.Style, termW, termH int) string {
	popup := style.Render(content)
	return lipgloss.Place(termW, termH, lipgloss.Center, lipgloss.Center, popup,
		lipgloss.WithWhitespaceChars(" "),
	)
}
