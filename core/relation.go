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
	if v.Queries == nil {
		return nil, fmt.Errorf("vault not opened")
	}
	return v.Queries.ListRelations(objectID)
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

// LinkObjects creates a relation between two objects. Delegates to ObjectService.
func (v *Vault) LinkObjects(fromID, relName, toID string) error {
	if v.Objects == nil {
		return fmt.Errorf("vault not opened")
	}
	return v.Objects.Link(fromID, relName, toID)
}

// UnlinkObjects removes a relation between two objects. Delegates to ObjectService.
func (v *Vault) UnlinkObjects(fromID, relName, toID string, both bool) error {
	if v.Objects == nil {
		return fmt.Errorf("vault not opened")
	}
	return v.Objects.Unlink(fromID, relName, toID, both)
}
