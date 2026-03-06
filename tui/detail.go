package tui

import (
	"fmt"
	"strings"

	"github.com/typemd/typemd/core"
)

// renderBody builds the body panel content: object title + markdown body.
func renderBody(obj *core.Object) string {
	if obj == nil {
		return "  Select an object to view details."
	}

	var b strings.Builder

	// Title
	b.WriteString(fmt.Sprintf(" %s\n\n", obj.ID))

	// Body section
	b.WriteString(" Body\n")
	b.WriteString(" ────\n")
	body := strings.TrimSpace(obj.Body)
	if body == "" {
		b.WriteString(" (empty)\n")
	} else {
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
			if p.IsReverse {
				b.WriteString(fmt.Sprintf(" %s: ← %s\n", p.Key, p.FromID))
			} else if p.Value == nil {
				b.WriteString(fmt.Sprintf(" %s: (null)\n", p.Key))
			} else if p.IsRelation {
				b.WriteString(fmt.Sprintf(" %s: → %v\n", p.Key, p.Value))
			} else {
				b.WriteString(fmt.Sprintf(" %s: %v\n", p.Key, p.Value))
			}
		}
	}

	return b.String()
}

