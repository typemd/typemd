package core

// DisplayProperty represents a single property prepared for display.
type DisplayProperty struct {
	Key        string
	Value      any
	IsRelation bool
	IsReverse  bool
	FromID     string // populated only for reverse relations
}

// BuildDisplayProperties assembles an ordered list of display-ready properties
// for an object, including schema defaults and reverse relations.
func (v *Vault) BuildDisplayProperties(obj *Object) ([]DisplayProperty, error) {
	if obj == nil {
		return nil, nil
	}

	schema, _ := v.LoadType(obj.Type)

	// Fill missing schema-defined properties
	if schema != nil {
		for _, p := range schema.Properties {
			if _, ok := obj.Properties[p.Name]; !ok {
				obj.Properties[p.Name] = nil
			}
		}
	}

	// Build relation property set
	relProps := make(map[string]bool)
	if schema != nil {
		for _, p := range schema.Properties {
			if p.Type == "relation" {
				relProps[p.Name] = true
			}
		}
	}

	// Get relations
	relations, _ := v.ListRelations(obj.ID)

	// Build ordered properties
	propKeys := OrderedPropKeys(obj.Properties, schema)
	var result []DisplayProperty
	for _, k := range propKeys {
		result = append(result, DisplayProperty{
			Key:        k,
			Value:      obj.Properties[k],
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

	return result, nil
}
