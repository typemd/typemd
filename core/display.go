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

// FormatValue returns the formatted value without the key prefix.
// Use this for contexts where the key is displayed separately (e.g. table columns).
func (p DisplayProperty) FormatValue() string {
	if p.IsBacklink {
		return "⟵ " + displayObjectID(p.FromID)
	}
	if p.IsReverse {
		return "← " + displayObjectID(p.FromID)
	}
	if p.Type == "checkbox" {
		if b, ok := p.Value.(bool); ok && b {
			return "☑"
		}
		return "☐"
	}
	if p.Value == nil {
		return ""
	}
	if p.IsRelation {
		return "→ " + displayObjectID(fmt.Sprintf("%v", p.Value))
	}

	switch p.Type {
	case "date":
		if t, ok := p.Value.(time.Time); ok {
			return t.Format("2006-01-02")
		}
	case "datetime":
		if t, ok := p.Value.(time.Time); ok {
			return t.Format("2006-01-02T15:04:05")
		}
	case "multi_select":
		if arr, ok := p.Value.([]any); ok {
			items := make([]string, len(arr))
			for i, v := range arr {
				items[i] = fmt.Sprintf("%v", v)
			}
			return "[" + strings.Join(items, ", ") + "]"
		}
	}

	return fmt.Sprintf("%v", p.Value)
}

// Format returns a human-readable string for this property (key: value).
func (p DisplayProperty) Format() string {
	return p.Key + ": " + p.FormatValue()
}

// BuildDisplayProperties assembles display-ready properties. Delegates to QueryService.
func (v *Vault) BuildDisplayProperties(obj *Object) ([]DisplayProperty, error) {
	if v.Queries == nil {
		return nil, nil
	}
	return v.Queries.BuildDisplayProperties(obj)
}
