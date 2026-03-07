package tui

import (
	"fmt"
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/typemd/typemd/core"
)

func wikiLinkStyle(s string) string {
	return wikiLinkStyleBase.Render(s)
}

// renderBodyHeader builds the title + separator for the body panel (2 lines).
// Used in both view mode (via renderBody) and edit mode (above the textarea).
func renderBodyHeader(obj *core.Object, width int) string {
	if obj == nil {
		return ""
	}
	title := obj.DisplayID()
	maxWidth := width - 1 // 1 for leading space
	if maxWidth > 0 {
		title = runewidth.Truncate(title, maxWidth, "…")
	}
	titleWidth := runewidth.StringWidth(title)
	return fmt.Sprintf(" %s\n %s\n", title, strings.Repeat("─", titleWidth))
}

// renderBody builds the body panel content: object title + markdown body.
func renderBody(obj *core.Object, width int) string {
	if obj == nil {
		return "  Select an object to view details."
	}

	var b strings.Builder

	b.WriteString(renderBodyHeader(obj, width))

	// Body section
	body := strings.TrimSpace(obj.Body)
	if body == "" {
		b.WriteString(" (empty)\n")
	} else {
		body = core.RenderWikiLinksStyled(body, wikiLinkStyle)
		for _, line := range strings.Split(body, "\n") {
			b.WriteString(fmt.Sprintf(" %s\n", line))
		}
	}

	return b.String()
}

// renderProperties builds the properties panel content using display properties.
func renderProperties(obj *core.Object, displayProps []core.DisplayProperty) string {
	if obj == nil {
		return ""
	}

	var b strings.Builder

	b.WriteString(" Properties\n")
	b.WriteString(" ──────────\n")

	if len(displayProps) == 0 {
		b.WriteString(" (none)\n")
	} else {
		for _, p := range displayProps {
			b.WriteString(fmt.Sprintf(" %s\n", p.Format()))
		}
	}

	return b.String()
}

