package widget

import "charm.land/lipgloss/v2"

// CenteredPopup renders content inside a styled popup, centered within termW × termH.
// The caller controls the popup style (border, width, padding, colors).
// The background is replaced with whitespace (legacy behavior).
func CenteredPopup(content string, style lipgloss.Style, termW, termH int) string {
	popup := style.Render(content)
	return lipgloss.Place(termW, termH, lipgloss.Center, lipgloss.Center, popup,
		lipgloss.WithWhitespaceChars(" "),
	)
}

// OverlayPopup renders a styled popup centered on top of a background string
// using lipgloss Layer/Compositor. The background remains visible outside the popup area.
func OverlayPopup(background, content string, style lipgloss.Style, termW, termH int) string {
	popup := style.Render(content)

	bgLayer := lipgloss.NewLayer(background).ID("bg").Z(0)

	// Center the popup by computing x/y offsets
	popupW := lipgloss.Width(popup)
	popupH := lipgloss.Height(popup)
	x := (termW - popupW) / 2
	y := (termH - popupH) / 2
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}

	popupLayer := lipgloss.NewLayer(popup).ID("popup").Z(10).X(x).Y(y)

	comp := lipgloss.NewCompositor(bgLayer, popupLayer)
	canvas := lipgloss.NewCanvas(termW, termH)
	return canvas.Compose(comp).Render()
}
