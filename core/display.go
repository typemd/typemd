package core

import (
	"fmt"
	"strings"
	"time"
)

// BacklinksDisplayKey is the property key used for wiki-link backlinks.
const BacklinksDisplayKey = "backlinks"

// DisplayProperty represents a single property prepared for display.
type DisplayProperty struct {
	Key        string
	Value      any
	Type       string // property type from schema (string, date, checkbox, etc.)
	Emoji      string // property emoji from schema
	Pin        int    // pin order (0 = not pinned, positive = pinned with order)
	IsRelation bool
	IsReverse  bool
	IsBacklink bool
	FromID     string // populated for reverse relations and backlinks
}

// displayObjectID strips the ULID suffix from an object ID for human-readable display.
// "person/robert-martin-01kk39c30y47xb1dvbs8ywqv50" → "person/robert-martin"
func displayObjectID(id string) string {
	if parsed, err := ParseObjectID(id); err == nil {
		return parsed.DisplayID()
	}
	return StripULID(id)
}

// Format returns a human-readable string for this property.
func (p DisplayProperty) Format() string {
	if p.IsBacklink {
		return fmt.Sprintf("%s: ⟵ %s", p.Key, displayObjectID(p.FromID))
	}
	if p.IsReverse {
		return fmt.Sprintf("%s: ← %s", p.Key, displayObjectID(p.FromID))
	}
	if p.Value == nil {
		return fmt.Sprintf("%s: (null)", p.Key)
	}
	if p.IsRelation {
		return fmt.Sprintf("%s: → %s", p.Key, displayObjectID(fmt.Sprintf("%v", p.Value)))
	}

	switch p.Type {
	case "checkbox":
		if b, ok := p.Value.(bool); ok {
			if b {
				return fmt.Sprintf("%s: [x]", p.Key)
			}
			return fmt.Sprintf("%s: [ ]", p.Key)
		}
	case "date":
		if t, ok := p.Value.(time.Time); ok {
			return fmt.Sprintf("%s: %s", p.Key, t.Format("2006-01-02"))
		}
	case "datetime":
		if t, ok := p.Value.(time.Time); ok {
			return fmt.Sprintf("%s: %s", p.Key, t.Format("2006-01-02T15:04:05"))
		}
	case "multi_select":
		if arr, ok := p.Value.([]any); ok {
			items := make([]string, len(arr))
			for i, v := range arr {
				items[i] = fmt.Sprintf("%v", v)
			}
			return fmt.Sprintf("%s: [%s]", p.Key, strings.Join(items, ", "))
		}
	}

	return fmt.Sprintf("%s: %v", p.Key, p.Value)
}

// BuildDisplayProperties assembles display-ready properties. Delegates to QueryService.
func (v *Vault) BuildDisplayProperties(obj *Object) ([]DisplayProperty, error) {
	if v.Queries == nil {
		return nil, nil
	}
	return v.Queries.BuildDisplayProperties(obj)
}
