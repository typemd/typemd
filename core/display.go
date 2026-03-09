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
	IsRelation bool
	IsReverse  bool
	IsBacklink bool
	FromID     string // populated for reverse relations and backlinks
}

// displayObjectID strips the ULID suffix from an object ID for human-readable display.
// "person/robert-martin-01kk39c30y47xb1dvbs8ywqv50" → "person/robert-martin"
func displayObjectID(id string) string {
	if i := strings.IndexByte(id, '/'); i >= 0 {
		return id[:i+1] + StripULID(id[i+1:])
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

// BuildDisplayProperties assembles an ordered list of display-ready properties
// for an object, including schema defaults and reverse relations.
func (v *Vault) BuildDisplayProperties(obj *Object) ([]DisplayProperty, error) {
	if obj == nil {
		return nil, nil
	}

	// LoadType may fail for unknown types; nil schema is safe here because
	// the code below skips schema-dependent logic when schema is nil.
	schema, _ := v.LoadType(obj.Type)

	// Build merged properties map (original + schema defaults) without mutating obj
	merged := make(map[string]any, len(obj.Properties))
	for k, v := range obj.Properties {
		merged[k] = v
	}

	// Single pass over schema: fill missing properties + build relation/type sets
	relProps := make(map[string]bool)
	propTypes := make(map[string]string)
	if schema != nil {
		for _, p := range schema.Properties {
			if _, ok := merged[p.Name]; !ok {
				merged[p.Name] = nil
			}
			if p.Type == "relation" {
				relProps[p.Name] = true
			}
			propTypes[p.Name] = p.Type
		}
	}

	// Get relations
	relations, err := v.ListRelations(obj.ID)
	if err != nil {
		return nil, fmt.Errorf("list relations: %w", err)
	}

	// Build ordered properties
	propKeys := OrderedPropKeys(merged, schema)
	var result []DisplayProperty
	for _, k := range propKeys {
		result = append(result, DisplayProperty{
			Key:        k,
			Value:      merged[k],
			Type:       propTypes[k],
			IsRelation: relProps[k],
		})
	}

	// Append reverse relations
	for _, r := range relations {
		if r.ToID == obj.ID {
			result = append(result, DisplayProperty{
				Key:       r.Name,
				Value:     r.FromID,
				IsReverse: true,
				FromID:    r.FromID,
			})
		}
	}

	// Append backlinks (wiki-links pointing to this object)
	backlinks, err := v.ListBacklinks(obj.ID)
	if err != nil {
		return nil, fmt.Errorf("list backlinks: %w", err)
	}
	for _, bl := range backlinks {
		result = append(result, DisplayProperty{
			Key:        BacklinksDisplayKey,
			Value:      bl.FromID,
			IsBacklink: true,
			FromID:     bl.FromID,
		})
	}

	return result, nil
}
