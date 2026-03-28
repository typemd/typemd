package tui

import (
	"github.com/typemd/typemd/tui/widget"
	"charm.land/lipgloss/v2"
)

// popupStyle returns the standard popup border style.
func popupStyle(popupWidth int) lipgloss.Style {
	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorFocusBorder).
		Padding(1, 2)

	if popupWidth > 0 {
		style = style.Width(popupWidth)
	}

	return style
}

// renderPopup renders content as a centered popup overlay with a rounded border.
// The border color uses the theme's colorFocusBorder. An optional width can be
// provided to constrain the popup; pass 0 to let lipgloss auto-size from content.
// The background is replaced with whitespace (legacy behavior).
func renderPopup(content string, termW, termH int, popupWidth int) string {
	return widget.CenteredPopup(content, popupStyle(popupWidth), termW, termH)
}

// renderOverlayPopup renders content as a centered popup on top of a background
// using lipgloss Layer/Compositor. The background remains visible outside the popup.
func renderOverlayPopup(background, content string, termW, termH int, popupWidth int) string {
	return widget.OverlayPopup(background, content, popupStyle(popupWidth), termW, termH)
}

// renderDatePickerPopup renders the date calendar overlay as a centered popup.
func renderDatePickerPopup(background, calendarContent string, termW, termH int) string {
	return widget.OverlayPopup(background, calendarContent, popupStyle(0), termW, termH)
}
