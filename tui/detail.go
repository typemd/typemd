package tui

import (
	"fmt"
	"strings"

	"github.com/MilesChou/typemd/core"
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
		b.WriteString("   (empty)\n")
	} else {
		for _, line := range strings.Split(body, "\n") {
			b.WriteString(fmt.Sprintf("   %s\n", line))
		}
	}

	return b.String()
}

// renderProperties builds the properties panel content: YAML frontmatter fields + relations.
func renderProperties(obj *core.Object, relations []core.Relation, schema *core.TypeSchema) string {
	if obj == nil {
		return ""
	}

	var b strings.Builder

	b.WriteString(" Properties\n")
	b.WriteString(" ──────────\n")

	// Build a set of property names that are relation types
	relProps := make(map[string]bool)
	if schema != nil {
		for _, p := range schema.Properties {
			if p.Type == "relation" {
				relProps[p.Name] = true
			}
		}
	}

	// Collect reverse relations
	var reverseRels []core.Relation
	for _, r := range relations {
		if r.ToID == obj.ID {
			reverseRels = append(reverseRels, r)
		}
	}

	if len(obj.Properties) == 0 && len(reverseRels) == 0 {
		b.WriteString("   (none)\n")
	} else {
		propKeys := core.OrderedPropKeys(obj.Properties, schema)
		for _, k := range propKeys {
			v := obj.Properties[k]
			if v == nil {
				b.WriteString(fmt.Sprintf("   %s: (null)\n", k))
			} else if relProps[k] {
				b.WriteString(fmt.Sprintf("   %s: → %v\n", k, v))
			} else {
				b.WriteString(fmt.Sprintf("   %s: %v\n", k, v))
			}
		}
		for _, r := range reverseRels {
			b.WriteString(fmt.Sprintf("   %s: ← %s\n", r.Name, r.FromID))
		}
	}

	return b.String()
}

