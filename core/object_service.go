package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// ObjectService orchestrates command-side operations on objects.
// It coordinates entities, repositories, and the index, and collects domain events.
type ObjectService struct {
	repo       ObjectRepository
	index      ObjectIndex
	dispatcher *EventDispatcher
}

// NewObjectService creates an ObjectService.
func NewObjectService(repo ObjectRepository, index ObjectIndex, dispatcher *EventDispatcher) *ObjectService {
	return &ObjectService{repo: repo, index: index, dispatcher: dispatcher}
}

// Create creates a new object with the given type and filename.
// If templateName is non-empty, the specified template is loaded and applied.
func (s *ObjectService) Create(typeName, filename, templateName string) (*Object, error) {
	schema, err := s.repo.GetSchema(typeName)
	if err != nil {
		return nil, fmt.Errorf("load type: %w", err)
	}

	// Load template if specified
	var tmpl *Template
	if templateName != "" {
		tmpl, err = s.repo.GetTemplate(typeName, templateName)
		if err != nil {
			return nil, fmt.Errorf("load template: %w", err)
		}
	}

	now := time.Now()

	// Handle empty name: check template, then name template, then error
	if filename == "" {
		if tmpl != nil {
			if nameVal, ok := tmpl.Properties[NameProperty]; ok {
				if s, ok := nameVal.(string); ok && s != "" {
					filename = s
				}
			}
		}
		if filename == "" {
			if schema.NameTemplate != "" {
				filename = EvaluateNameTemplate(schema.NameTemplate, now)
			} else {
				return nil, fmt.Errorf("name is required (type %q has no name template)", typeName)
			}
		}
	}

	// Enforce name uniqueness for types with unique constraint
	if schema.Unique {
		if err := s.checkNameUnique(typeName, filename); err != nil {
			return nil, err
		}
	}

	// Generate ObjectID with ULID suffix
	slug := filename
	objID, err := NewObjectID(typeName, slug)
	if err != nil {
		return nil, err
	}

	// Create type directory
	if err := s.repo.EnsureDir(typeName); err != nil {
		return nil, fmt.Errorf("create directory: %w", err)
	}

	// Generate initial properties from schema defaults
	props := make(map[string]any)
	props[NameProperty] = slug
	nowStr := now.Format(time.RFC3339)
	props[CreatedAtProperty] = nowStr
	props[UpdatedAtProperty] = nowStr
	for _, p := range schema.Properties {
		if p.Default != nil {
			props[p.Name] = p.Default
		} else {
			props[p.Name] = nil
		}
	}

	// Build the new object entity
	newObj := &Object{
		ID:         objID.String(),
		Type:       typeName,
		Filename:   objID.Filename,
		Properties: props,
	}

	// Apply template (overrides schema defaults)
	if tmpl != nil {
		newObj.ApplyTemplate(tmpl, schema)
	}
	if err := s.repo.Create(newObj, OrderedPropKeys(props, schema)); err != nil {
		return nil, fmt.Errorf("create object file: %w", err)
	}

	// Insert into index
	propsJSON, err := json.Marshal(props)
	if err != nil {
		return nil, fmt.Errorf("marshal properties: %w", err)
	}
	if err := s.index.Upsert(newObj.ID, typeName, objID.Filename, string(propsJSON), newObj.Body); err != nil {
		return nil, fmt.Errorf("insert object: %w", err)
	}

	s.dispatcher.Dispatch([]DomainEvent{ObjectCreated{Object: newObj}})
	return newObj, nil
}

// Save persists an object's properties and body to file and index.
func (s *ObjectService) Save(obj *Object) error {
	obj.MarkUpdated()
	schema, _ := s.repo.GetSchema(obj.Type)
	keyOrder := OrderedPropKeys(obj.Properties, schema)

	if err := s.repo.Save(obj, keyOrder); err != nil {
		return fmt.Errorf("save object file: %w", err)
	}

	propsJSON, err := json.Marshal(obj.Properties)
	if err != nil {
		return fmt.Errorf("marshal properties: %w", err)
	}
	if err := s.index.Upsert(obj.ID, obj.Type, obj.Filename, string(propsJSON), obj.Body); err != nil {
		return fmt.Errorf("update index: %w", err)
	}

	s.dispatcher.Dispatch([]DomainEvent{ObjectSaved{Object: obj}})
	return nil
}

// SetProperty updates a single property on an object.
func (s *ObjectService) SetProperty(id, key string, value any) error {
	obj, err := s.repo.Get(id)
	if err != nil {
		return fmt.Errorf("get object: %w", err)
	}

	schema, err := s.repo.GetSchema(obj.Type)
	if err != nil {
		return fmt.Errorf("load type: %w", err)
	}

	event, err := obj.SetProperty(key, value, schema)
	if err != nil {
		return err
	}

	if err := s.Save(obj); err != nil {
		return err
	}

	s.dispatcher.Dispatch([]DomainEvent{event})
	return nil
}

// Link creates a relation between two objects.
func (s *ObjectService) Link(fromID, relName, toID string) error {
	fromObj, fromSchema, err := s.loadObjectAndSchema(fromID)
	if err != nil {
		return fmt.Errorf("get source: %w", err)
	}

	relProp := fromSchema.FindRelation(relName)
	if relProp == nil {
		return fmt.Errorf("relation %q not found in type %q", relName, fromObj.Type)
	}

	toObj, err := s.repo.Get(toID)
	if err != nil {
		return fmt.Errorf("get target object: %w", err)
	}
	if relProp.Target != "" && toObj.Type != relProp.Target {
		return fmt.Errorf("target type mismatch: expected %q, got %q", relProp.Target, toObj.Type)
	}

	event, err := fromObj.LinkTo(relName, toID, relProp)
	if err != nil {
		return fmt.Errorf("relation already exists: %s -[%s]-> %s", fromID, relName, toID)
	}
	if err := s.Save(fromObj); err != nil {
		return fmt.Errorf("write source object: %w", err)
	}
	if err := s.index.InsertRelation(relName, fromID, toID); err != nil {
		return fmt.Errorf("insert relation: %w", err)
	}

	var events []DomainEvent
	events = append(events, event)

	// Handle bidirectional
	if relProp.Bidirectional && relProp.Inverse != "" {
		toSchema, err := s.repo.GetSchema(toObj.Type)
		if err != nil {
			return fmt.Errorf("load target type: %w", err)
		}

		inverseProp := toSchema.FindRelation(relProp.Inverse)
		if inverseProp == nil {
			return fmt.Errorf("inverse relation %q not found in type %q", relProp.Inverse, toObj.Type)
		}

		invEvent, err := toObj.LinkTo(relProp.Inverse, fromID, inverseProp)
		if err != nil && !errors.Is(err, errDuplicateRelation) {
			return fmt.Errorf("set inverse relation: %w", err)
		}
		if err := s.Save(toObj); err != nil {
			return fmt.Errorf("write target object: %w", err)
		}
		if err := s.index.InsertRelation(relProp.Inverse, toID, fromID); err != nil {
			return fmt.Errorf("insert inverse relation: %w", err)
		}
		if invEvent != nil {
			events = append(events, invEvent)
		}
	}

	s.dispatcher.Dispatch(events)
	return nil
}

// Unlink removes a relation between two objects.
// If both is true and the relation is bidirectional, also removes the inverse.
func (s *ObjectService) Unlink(fromID, relName, toID string, both bool) error {
	fromObj, fromSchema, err := s.loadObjectAndSchema(fromID)
	if err != nil {
		return fmt.Errorf("get source: %w", err)
	}

	relProp := fromSchema.FindRelation(relName)
	if relProp == nil {
		return fmt.Errorf("relation %q not found in type %q", relName, fromObj.Type)
	}

	event := fromObj.Unlink(relName, toID, relProp)
	if err := s.Save(fromObj); err != nil {
		return fmt.Errorf("write source object: %w", err)
	}
	if err := s.index.DeleteRelation(relName, fromID, toID); err != nil {
		return fmt.Errorf("delete relation: %w", err)
	}

	var events []DomainEvent
	events = append(events, event)

	// Handle --both with bidirectional
	if both && relProp.Bidirectional && relProp.Inverse != "" {
		toObj, toSchema, err := s.loadObjectAndSchema(toID)
		if err != nil {
			return fmt.Errorf("get target: %w", err)
		}

		inverseProp := toSchema.FindRelation(relProp.Inverse)
		if inverseProp == nil {
			return fmt.Errorf("inverse relation %q not found in type %q", relProp.Inverse, toObj.Type)
		}

		invEvent := toObj.Unlink(relProp.Inverse, fromID, inverseProp)
		if err := s.Save(toObj); err != nil {
			return fmt.Errorf("write target object: %w", err)
		}
		if err := s.index.DeleteRelation(relProp.Inverse, toID, fromID); err != nil {
			return fmt.Errorf("delete inverse relation: %w", err)
		}
		events = append(events, invEvent)
	}

	s.dispatcher.Dispatch(events)
	return nil
}

// loadObjectAndSchema loads an object and its type schema.
func (s *ObjectService) loadObjectAndSchema(id string) (*Object, *TypeSchema, error) {
	obj, err := s.repo.Get(id)
	if err != nil {
		return nil, nil, fmt.Errorf("get object: %w", err)
	}
	schema, err := s.repo.GetSchema(obj.Type)
	if err != nil {
		return nil, nil, fmt.Errorf("load type: %w", err)
	}
	return obj, schema, nil
}

// checkNameUnique returns an error if an object with the given name already exists.
func (s *ObjectService) checkNameUnique(typeName, name string) error {
	results, err := s.index.Query(fmt.Sprintf("type=%s name=%s", typeName, name))
	if err != nil {
		return fmt.Errorf("check name uniqueness: %w", err)
	}
	if len(results) > 0 {
		return fmt.Errorf("%s name %q already exists: %s", typeName, name, results[0].ID)
	}
	return nil
}
