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

// renderTitleContent builds the title string for the title panel.
// Format: "emoji type · DisplayName" or "type · DisplayName" when no emoji.
func renderTitleContent(obj *core.Object, typeName, emoji string, width int) string {
	if obj == nil {
		return ""
	}
	var title string
	if emoji != "" {
		title = fmt.Sprintf(" %s %s · %s", emoji, typeName, obj.DisplayName())
	} else {
		title = fmt.Sprintf(" %s · %s", typeName, obj.DisplayName())
	}
	maxWidth := width
	if maxWidth > 0 {
		title = runewidth.Truncate(title, maxWidth, "…")
	}
	return title
}

// renderBody builds the body panel content: markdown body only (no title header).
func renderBody(obj *core.Object, width int) string {
	if obj == nil {
		return "  Select an object to view details."
	}

	var b strings.Builder

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

