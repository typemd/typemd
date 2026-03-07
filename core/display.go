package core

import "fmt"

// BacklinksDisplayKey is the property key used for wiki-link backlinks.
const BacklinksDisplayKey = "backlinks"

// DisplayProperty represents a single property prepared for display.
type DisplayProperty struct {
	Key        string
	Value      any
	IsRelation bool
	IsReverse  bool
	IsBacklink bool
	FromID     string // populated for reverse relations and backlinks
}

// Format returns a human-readable string for this property.
func (p DisplayProperty) Format() string {
	if p.IsBacklink {
		return fmt.Sprintf("%s: ⟵ %s", p.Key, p.FromID)
	}
	if p.IsReverse {
		return fmt.Sprintf("%s: ← %s", p.Key, p.FromID)
	}
	if p.Value == nil {
		return fmt.Sprintf("%s: (null)", p.Key)
	}
	if p.IsRelation {
		return fmt.Sprintf("%s: → %v", p.Key, p.Value)
	}
	return fmt.Sprintf("%s: %v", p.Key, p.Value)
}

// BuildDisplayProperties assembles an ordered list of display-ready properties
// for an object, including schema defaults and reverse relations.
func (v *Vault) BuildDisplayProperties(obj *Object) ([]DisplayProperty, error) {
	if obj == nil {
		return nil, nil
	}

	schema, _ := v.LoadType(obj.Type)

	// Build merged properties map (original + schema defaults) without mutating obj
	merged := make(map[string]any, len(obj.Properties))
	for k, v := range obj.Properties {
		merged[k] = v
	}

	// Single pass over schema: fill missing properties + build relation set
	relProps := make(map[string]bool)
	if schema != nil {
		for _, p := range schema.Properties {
			if _, ok := merged[p.Name]; !ok {
				merged[p.Name] = nil
			}
			if p.Type == "relation" {
				relProps[p.Name] = true
			}
		}
	}

	// Get relations
	relations, _ := v.ListRelations(obj.ID)

	// Build ordered properties
	propKeys := OrderedPropKeys(merged, schema)
	var result []DisplayProperty
	for _, k := range propKeys {
		result = append(result, DisplayProperty{
			Key:        k,
			Value:      merged[k],
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
	backlinks, _ := v.ListBacklinks(obj.ID)
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
