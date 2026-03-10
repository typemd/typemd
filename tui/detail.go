package tui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/typemd/typemd/core"
)

func wikiLinkStyle(s string) string {
	return wikiLinkStyleBase.Render(s)
}

// renderTitleContent builds the title string for the title panel.
// Format: "emoji type · name" or "type · name" when no emoji.
func renderTitleContent(obj *core.Object, typeName, emoji string, width int) string {
	if obj == nil {
		return ""
	}
	var title string
	if emoji != "" {
		title = fmt.Sprintf(" %s %s · %s", emoji, typeName, obj.GetName())
	} else {
		title = fmt.Sprintf(" %s · %s", typeName, obj.GetName())
	}
	maxWidth := width
	if maxWidth > 0 {
		title = runewidth.Truncate(title, maxWidth, "…")
	}
	return title
}

// renderBody builds the body panel content: pinned properties at top, then markdown body.
func renderBody(obj *core.Object, width int, displayProps []core.DisplayProperty) string {
	if obj == nil {
		return "  Select an object to view details."
	}

	var b strings.Builder
	body := strings.TrimSpace(obj.Body)

	// Pinned properties section
	pinned := pinnedProperties(displayProps)
	if len(pinned) > 0 {
		for _, p := range pinned {
			if p.Emoji != "" {
				b.WriteString(fmt.Sprintf(" %s %s\n", p.Emoji, p.Format()))
			} else {
				b.WriteString(fmt.Sprintf(" %s\n", p.Format()))
			}
		}
		// Separator only if there is body content
		if body != "" {
			b.WriteString(" ────────────────────\n")
		}
	}

	// Body section
	if body == "" && len(pinned) == 0 {
		b.WriteString(" (empty)\n")
	} else if body != "" {
		body = core.RenderWikiLinksStyled(body, wikiLinkStyle)
		for _, line := range strings.Split(body, "\n") {
			b.WriteString(fmt.Sprintf(" %s\n", line))
		}
	}

	return b.String()
}

// pinnedProperties returns display properties with Pin > 0, sorted by Pin value.
func pinnedProperties(props []core.DisplayProperty) []core.DisplayProperty {
	var pinned []core.DisplayProperty
	for _, p := range props {
		if p.Pin > 0 {
			pinned = append(pinned, p)
		}
	}
	sort.Slice(pinned, func(i, j int) bool {
		return pinned[i].Pin < pinned[j].Pin
	})
	return pinned
}

// renderProperties builds the properties panel content using display properties.
func renderProperties(obj *core.Object, displayProps []core.DisplayProperty) string {
	if obj == nil {
		return ""
	}

	var b strings.Builder

	b.WriteString(" Properties\n")
	b.WriteString(" ──────────\n")

	// Filter out pinned properties
	var unpinned []core.DisplayProperty
	for _, p := range displayProps {
		if p.Pin == 0 {
			unpinned = append(unpinned, p)
		}
	}

	if len(unpinned) == 0 {
		b.WriteString(" (none)\n")
	} else {
		for _, p := range unpinned {
			b.WriteString(fmt.Sprintf(" %s\n", p.Format()))
		}
	}

	return b.String()
}

