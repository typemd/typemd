package core

import (
	"errors"
	"fmt"
)

// errDuplicateRelation is returned when appending a value that already exists.
var errDuplicateRelation = errors.New("duplicate relation value")

// Relation represents a relationship record between two objects.
type Relation struct {
	Name   string
	FromID string
	ToID   string
}

// ListRelations returns all relations where objectID is either the source or target.
func (v *Vault) ListRelations(objectID string) ([]Relation, error) {
	if v.db == nil {
		return nil, fmt.Errorf("vault not opened")
	}

	rows, err := v.db.Query(
		"SELECT name, from_id, to_id FROM relations WHERE from_id = ? OR to_id = ?",
		objectID, objectID,
	)
	if err != nil {
		return nil, fmt.Errorf("list relations: %w", err)
	}
	defer rows.Close()

	var rels []Relation
	for rows.Next() {
		var r Relation
		if err := rows.Scan(&r.Name, &r.FromID, &r.ToID); err != nil {
			return nil, fmt.Errorf("scan relation: %w", err)
		}
		rels = append(rels, r)
	}
	return rels, rows.Err()
}

// getRelationSlice extracts the existing []any from a property value.
func getRelationSlice(props map[string]any, name string) []any {
	existing, _ := props[name]
	if existing == nil {
		return nil
	}
	if arr, ok := existing.([]any); ok {
		return arr
	}
	return nil
}

// appendRelationValue adds a value to a relation property.
// For multiple relations, it appends to the existing array.
// Returns an error if the value already exists in a multiple relation.
func appendRelationValue(props map[string]any, name, value string, multiple bool) error {
	if !multiple {
		props[name] = value
		return nil
	}
	arr := getRelationSlice(props, name)
	for _, item := range arr {
		if item == value {
			return errDuplicateRelation
		}
	}
	props[name] = append(arr, any(value))
	return nil
}

// removeRelationValue removes a value from a relation property.
// For multiple relations, it filters the value out of the existing array.
func removeRelationValue(props map[string]any, name, value string, multiple bool) {
	if !multiple {
		if props[name] == value {
			props[name] = nil
		}
		return
	}
	arr := getRelationSlice(props, name)
	var newArr []any
	for _, item := range arr {
		if item != value {
			newArr = append(newArr, item)
		}
	}
	if len(newArr) == 0 {
		props[name] = nil
	} else {
		props[name] = newArr
	}
}

// findRelationProperty finds a relation property by name in a type schema.
func findRelationProperty(schema *TypeSchema, name string) *Property {
	for i, p := range schema.Properties {
		if p.Name == name && p.Type == "relation" {
			return &schema.Properties[i]
		}
	}
	return nil
}

// findSystemRelationProperty returns a Property representation of a system property relation.
func findSystemRelationProperty(name string) *Property {
	for _, sp := range systemProperties {
		if sp.Name == name && sp.Type == "relation" {
			return &Property{
				Name:     sp.Name,
				Type:     "relation",
				Target:   sp.Target,
				Multiple: sp.Multiple,
			}
		}
	}
	return nil
}

// resolveRelationProperty looks up a relation property by name, checking both
// the type schema and system properties.
func resolveRelationProperty(schema *TypeSchema, name string) *Property {
	if p := findRelationProperty(schema, name); p != nil {
		return p
	}
	return findSystemRelationProperty(name)
}

// loadObjectAndSchema loads an object and its type schema.
func (v *Vault) loadObjectAndSchema(id string) (*Object, *TypeSchema, error) {
	obj, err := v.GetObject(id)
	if err != nil {
		return nil, nil, fmt.Errorf("get object: %w", err)
	}
	schema, err := v.LoadType(obj.Type)
	if err != nil {
		return nil, nil, fmt.Errorf("load type: %w", err)
	}
	return obj, schema, nil
}

// LinkObjects creates a relation between two objects.
func (v *Vault) LinkObjects(fromID, relName, toID string) error {
	if v.db == nil {
		return fmt.Errorf("vault not opened")
	}

	fromObj, fromSchema, err := v.loadObjectAndSchema(fromID)
	if err != nil {
		return fmt.Errorf("get source: %w", err)
	}

	relProp := resolveRelationProperty(fromSchema, relName)
	if relProp == nil {
		return fmt.Errorf("relation %q not found in type %q", relName, fromObj.Type)
	}

	toObj, err := v.GetObject(toID)
	if err != nil {
		return fmt.Errorf("get target object: %w", err)
	}
	if relProp.Target != "" && toObj.Type != relProp.Target {
		return fmt.Errorf("target type mismatch: expected %q, got %q", relProp.Target, toObj.Type)
	}

	if err := appendRelationValue(fromObj.Properties, relName, toID, relProp.Multiple); err != nil {
		return fmt.Errorf("relation already exists: %s -[%s]-> %s", fromID, relName, toID)
	}
	if err := v.saveObjectFile(fromObj); err != nil {
		return fmt.Errorf("write source object: %w", err)
	}
	if _, err := v.db.Exec(
		"INSERT INTO relations (name, from_id, to_id) VALUES (?, ?, ?)",
		relName, fromID, toID,
	); err != nil {
		return fmt.Errorf("insert relation: %w", err)
	}

	// Handle bidirectional
	if relProp.Bidirectional && relProp.Inverse != "" {
		toSchema, err := v.LoadType(toObj.Type)
		if err != nil {
			return fmt.Errorf("load target type: %w", err)
		}

		inverseProp := findRelationProperty(toSchema, relProp.Inverse)
		if inverseProp == nil {
			return fmt.Errorf("inverse relation %q not found in type %q", relProp.Inverse, toObj.Type)
		}

		if err := appendRelationValue(toObj.Properties, relProp.Inverse, fromID, inverseProp.Multiple); err != nil && !errors.Is(err, errDuplicateRelation) {
			return fmt.Errorf("set inverse relation: %w", err)
		}
		if err := v.saveObjectFile(toObj); err != nil {
			return fmt.Errorf("write target object: %w", err)
		}
		if _, err := v.db.Exec(
			"INSERT INTO relations (name, from_id, to_id) VALUES (?, ?, ?)",
			relProp.Inverse, toID, fromID,
		); err != nil {
			return fmt.Errorf("insert inverse relation: %w", err)
		}
	}

	return nil
}

// UnlinkObjects removes a relation between two objects.
// If both is true and the relation is bidirectional, also removes the inverse.
func (v *Vault) UnlinkObjects(fromID, relName, toID string, both bool) error {
	if v.db == nil {
		return fmt.Errorf("vault not opened")
	}

	fromObj, fromSchema, err := v.loadObjectAndSchema(fromID)
	if err != nil {
		return fmt.Errorf("get source: %w", err)
	}

	relProp := resolveRelationProperty(fromSchema, relName)
	if relProp == nil {
		return fmt.Errorf("relation %q not found in type %q", relName, fromObj.Type)
	}

	removeRelationValue(fromObj.Properties, relName, toID, relProp.Multiple)
	if err := v.saveObjectFile(fromObj); err != nil {
		return fmt.Errorf("write source object: %w", err)
	}
	if _, err := v.db.Exec(
		"DELETE FROM relations WHERE name = ? AND from_id = ? AND to_id = ?",
		relName, fromID, toID,
	); err != nil {
		return fmt.Errorf("delete relation: %w", err)
	}

	// Handle --both with bidirectional
	if both && relProp.Bidirectional && relProp.Inverse != "" {
		toObj, toSchema, err := v.loadObjectAndSchema(toID)
		if err != nil {
			return fmt.Errorf("get target: %w", err)
		}

		inverseProp := findRelationProperty(toSchema, relProp.Inverse)
		if inverseProp == nil {
			return fmt.Errorf("inverse relation %q not found in type %q", relProp.Inverse, toObj.Type)
		}

		removeRelationValue(toObj.Properties, relProp.Inverse, fromID, inverseProp.Multiple)
		if err := v.saveObjectFile(toObj); err != nil {
			return fmt.Errorf("write target object: %w", err)
		}
		if _, err := v.db.Exec(
			"DELETE FROM relations WHERE name = ? AND from_id = ? AND to_id = ?",
			relProp.Inverse, toID, fromID,
		); err != nil {
			return fmt.Errorf("delete inverse relation: %w", err)
		}
	}

	return nil
}
