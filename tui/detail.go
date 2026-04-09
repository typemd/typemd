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

// titlePrefix returns the " emoji type" or " type" prefix used in title panels.
func titlePrefix(emoji, typeName string) string {
	if emoji != "" {
		return fmt.Sprintf(" %s %s", padEmoji(emoji), typeName)
	}
	return fmt.Sprintf(" %s", typeName)
}

// renderTitleContent builds the title string for the title panel.
// Format: "emoji type · name" or "type · name" when no emoji.
// Locked objects show a 🔒 badge right-aligned within the available width.
func renderTitleContent(obj *core.Object, typeName, emoji string, width int) string {
	if obj == nil {
		return ""
	}
	title := titlePrefix(emoji, typeName) + " · " + obj.GetName()
	if !obj.IsLocked() {
		if width > 0 {
			title = truncate(title, width)
		}
		return title
	}
	badge := " 🔒"
	badgeW := runewidth.StringWidth(badge)
	if width > 0 {
		title = truncate(title, width-badgeW)
		pad := width - runewidth.StringWidth(title) - badgeW
		if pad > 0 {
			title += strings.Repeat(" ", pad)
		}
	}
	return title + badge
}

// renderBody builds the body panel content: markdown body only.
func renderBody(obj *core.Object, width int) string {
	if obj == nil {
		return "  Select an object to view details."
	}

	var b strings.Builder
	body := strings.TrimSpace(obj.Body)

	if body == "" {
		b.WriteString(" (empty)\n")
	} else {
		body = renderMarkdown(body)
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
// When obj is nil but displayProps is provided (e.g., template preview), it still renders.
// Returns empty string only when both obj and displayProps are nil/empty.
func renderProperties(obj *core.Object, displayProps []core.DisplayProperty) string {
	if obj == nil && len(displayProps) == 0 {
		return ""
	}

	var b strings.Builder

	b.WriteString(" Properties\n")
	b.WriteString(" ──────────\n")

	// Pinned properties first (sorted by Pin), then unpinned; name shown in title panel
	pinned := pinnedProperties(displayProps)
	var unpinned []core.DisplayProperty
	for _, p := range displayProps {
		if p.Key == core.NameProperty || p.Pin > 0 {
			continue
		}
		unpinned = append(unpinned, p)
	}

	if len(pinned)+len(unpinned) == 0 {
		b.WriteString(" (none)\n")
	} else {
		for _, p := range pinned {
			b.WriteString(fmt.Sprintf(" %s\n", p.Format()))
		}
		if len(pinned) > 0 && len(unpinned) > 0 {
			b.WriteString(" ──────────\n")
		}
		for _, p := range unpinned {
			b.WriteString(fmt.Sprintf(" %s\n", p.Format()))
		}
	}

	return b.String()
}

