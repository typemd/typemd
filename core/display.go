package core

import (
	"fmt"
	"strings"
	"time"
)

// BacklinksDisplayKey is the property key used for wiki-link backlinks.
const BacklinksDisplayKey = "backlinks"

// Default display format tokens for date and datetime properties.
const (
	DefaultDateFormat     = "YYYY-MM-DD"
	DefaultDatetimeFormat = "YYYY-MM-DD HH:mm:ss"
)

// DisplayProperty represents a single property prepared for display.
type DisplayProperty struct {
	Key            string
	Value          any
	Type           string // property type from schema (string, date, checkbox, etc.)
	Emoji          string // property emoji from schema
	Pin            int    // pin order (0 = not pinned, positive = pinned with order)
	IsRelation     bool
	IsReverse      bool
	IsBacklink     bool
	IsLocal        bool   // true when property exists in object but not in type schema or system properties
	FromID         string // populated for reverse relations and backlinks
	DateFormat     string
	DatetimeFormat string
}

// displayObjectID strips the ULID suffix from an object ID for human-readable display.
// "person/robert-martin-01kk39c30y47xb1dvbs8ywqv50" → "person/robert-martin"
func displayObjectID(id string) string {
	if parsed, err := ParseObjectID(id); err == nil {
		return parsed.DisplayID()
	}
	return StripULID(id)
}

// asTime tries to extract a time.Time from the property value.
// It handles both time.Time values and string values parsed with the given layouts.
func (p DisplayProperty) asTime(layouts ...string) (time.Time, bool) {
	if t, ok := p.Value.(time.Time); ok {
		return t, true
	}
	if s, ok := p.Value.(string); ok && s != "" {
		for _, layout := range layouts {
			if t, err := time.Parse(layout, s); err == nil {
				return t, true
			}
		}
	}
	return time.Time{}, false
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
		if t, ok := p.asTime("2006-01-02"); ok {
			f := p.DateFormat
			if f == "" {
				f = DefaultDateFormat
			}
			return t.Format(ConvertDateFormat(f))
		}
	case "datetime":
		if t, ok := p.asTime(time.RFC3339, "2006-01-02T15:04:05"); ok {
			f := p.DatetimeFormat
			if f == "" {
				f = DefaultDatetimeFormat
			}
			return t.Format(ConvertDateFormat(f))
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

// BuildDisplayProperties assembles display-ready properties. Delegates to QueryService
// and injects configured date/datetime display formats from vault config.
func (v *Vault) BuildDisplayProperties(obj *Object) ([]DisplayProperty, error) {
	if v.Queries == nil {
		return nil, nil
	}
	props, err := v.Queries.BuildDisplayProperties(obj)
	if err != nil {
		return nil, err
	}

	cfg := v.config
	if cfg != nil {
		for i := range props {
			switch props[i].Type {
			case "date":
				if cfg.DateFormat != "" {
					props[i].DateFormat = cfg.DateFormat
				}
			case "datetime":
				if cfg.DatetimeFormat != "" {
					props[i].DatetimeFormat = cfg.DatetimeFormat
				}
			}
		}
	}

	return props, nil
}
