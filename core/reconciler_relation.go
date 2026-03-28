package core

import (
	"fmt"
	"strings"
)

// reconcileRelations resolves name references, writes back expanded IDs, and emits
// relation index events for schema-defined relation properties (excluding tags).
func (r *Reconciler) reconcileRelations(ctx *syncContext, result *ReconcileResult) ([]DomainEvent, error) {
	var events []DomainEvent
	events = append(events, RelationsCleared{NonTagOnly: true})

	modified := make(map[string]bool)
	for _, obj := range ctx.diskObjects {
		objEvents, err := r.processObjectRelations(obj, ctx, result, modified)
		if err != nil {
			return nil, err
		}
		events = append(events, objEvents...)
	}
	if err := r.writeBackModified(modified, ctx); err != nil {
		return nil, err
	}
	return events, nil
}

// reconcileRelationsForObjects resolves and emits relation events for specific changed objects.
func (r *Reconciler) reconcileRelationsForObjects(ctx *syncContext, result *ReconcileResult, objectIDs []string) ([]DomainEvent, error) {
	var events []DomainEvent

	// Emit per-object clear events for changed objects
	for _, id := range objectIDs {
		events = append(events, RelationsCleared{ObjectID: id})
	}

	modified := make(map[string]bool)
	for _, id := range objectIDs {
		obj := ctx.diskObjects[id]
		if obj == nil {
			continue
		}
		objEvents, err := r.processObjectRelations(obj, ctx, result, modified)
		if err != nil {
			return nil, err
		}
		events = append(events, objEvents...)
	}
	if err := r.writeBackModified(modified, ctx); err != nil {
		return nil, err
	}
	return events, nil
}

// resolveAndEmitRelation resolves a single relation reference and emits a RelationIndexed event.
func (r *Reconciler) resolveAndEmitRelation(ref, objectID, propName string, ctx *syncContext, result *ReconcileResult) (string, bool, []DomainEvent, error) {
	resolved, changed, err := resolveRelationValue(ref, ctx.nameIndex)
	if err != nil {
		result.Unresolved = append(result.Unresolved, UnresolvedRelation{
			ObjectID: objectID, Property: propName, Value: ref,
			Reason: resolveErrorReason(err), Matches: resolveErrorMatches(err),
		})
		return ref, false, nil, nil
	}
	if changed {
		result.Expanded++
	}
	var events []DomainEvent
	if ctx.diskIDs[resolved] {
		events = append(events, RelationIndexed{Name: propName, FromID: objectID, ToID: resolved})
	}
	return resolved, changed, events, nil
}

// processObjectRelations resolves relation properties for a single object and emits events.
func (r *Reconciler) processObjectRelations(obj *Object, ctx *syncContext, result *ReconcileResult, modified map[string]bool) ([]DomainEvent, error) {
	schema := ctx.schemas[obj.Type]
	if schema == nil {
		s, err := r.repo.GetSchema(obj.Type)
		if err != nil {
			return nil, nil // skip objects with unknown type
		}
		schema = s
		ctx.schemas[obj.Type] = schema
	}

	var events []DomainEvent
	for _, prop := range schema.Properties {
		if prop.Type != "relation" {
			continue
		}

		val, ok := obj.Properties[prop.Name]
		if !ok || val == nil {
			continue
		}

		if prop.Multiple {
			items, ok := val.([]any)
			if !ok {
				continue
			}
			anyChanged := false
			newItems := make([]any, len(items))
			for i, item := range items {
				ref, ok := item.(string)
				if !ok {
					newItems[i] = item
					continue
				}
				resolved, changed, evts, err := r.resolveAndEmitRelation(ref, obj.ID, prop.Name, ctx, result)
				if err != nil {
					return nil, err
				}
				events = append(events, evts...)
				if changed {
					anyChanged = true
				}
				newItems[i] = resolved
			}
			if anyChanged {
				obj.Properties[prop.Name] = newItems
				modified[obj.ID] = true
			}
		} else {
			ref, ok := val.(string)
			if !ok {
				continue
			}
			resolved, changed, evts, err := r.resolveAndEmitRelation(ref, obj.ID, prop.Name, ctx, result)
			if err != nil {
				return nil, err
			}
			events = append(events, evts...)
			if changed {
				obj.Properties[prop.Name] = resolved
				modified[obj.ID] = true
			}
		}
	}
	return events, nil
}

// writeBackModified saves objects whose relation properties were expanded during reconciliation.
func (r *Reconciler) writeBackModified(modified map[string]bool, ctx *syncContext) error {
	for id := range modified {
		obj := ctx.diskObjects[id]
		schema := ctx.schemas[obj.Type]
		keyOrder := OrderedPropKeys(obj.Properties, schema)
		if err := r.repo.Save(obj, keyOrder); err != nil {
			return fmt.Errorf("write-back expanded relations for %s: %w", id, err)
		}
	}
	return nil
}

// resolveErrorReason extracts the reason from a resolve error.
func resolveErrorReason(err error) string {
	if _, ok := err.(*AmbiguousMatchError); ok {
		return ReasonAmbiguous
	}
	return ReasonNotFound
}

// resolveErrorMatches extracts match candidates from an AmbiguousMatchError.
func resolveErrorMatches(err error) []string {
	if ame, ok := err.(*AmbiguousMatchError); ok {
		return ame.Matches
	}
	return nil
}

// reconcileTagRelations clears existing tag relations and rebuilds them, emitting events.
func (r *Reconciler) reconcileTagRelations(ctx *syncContext) ([]DomainEvent, error) {
	var events []DomainEvent
	events = append(events, RelationsCleared{TagsOnly: true})

	tagNameIndex := make(map[string]string)
	for _, obj := range ctx.diskTags {
		if name, ok := obj.Properties[NameProperty].(string); ok {
			tagNameIndex[name] = obj.ID
		}
	}

	for objID, refs := range ctx.diskTagRefs {
		for _, ref := range refs {
			tagID, err := r.resolveOrCreateTag(ref, ctx, tagNameIndex)
			if err != nil {
				continue
			}
			events = append(events, RelationIndexed{Name: TagsProperty, FromID: objID, ToID: tagID})
		}
	}
	return events, nil
}

// resolveOrCreateTag resolves a tag reference to an object ID, auto-creating if needed.
func (r *Reconciler) resolveOrCreateTag(ref string, ctx *syncContext, tagNameIndex map[string]string) (string, error) {
	if tagID, ok := resolveTagReference(ref, ctx.diskTags, tagNameIndex); ok {
		return tagID, nil
	}

	slug := strings.TrimPrefix(ref, "tag/")

	if ulidSuffixPattern.MatchString(slug) {
		return "", fmt.Errorf("broken tag reference: %s", ref)
	}

	if existingID, exists := tagNameIndex[slug]; exists {
		return existingID, nil
	}

	if r.createTag == nil {
		return "", fmt.Errorf("cannot auto-create tag %q: no createTag callback", slug)
	}

	newTag, err := r.createTag(slug)
	if err != nil {
		return "", fmt.Errorf("auto-create tag %q: %w", slug, err)
	}
	ctx.diskTags[newTag.ID] = newTag
	ctx.diskIDs[newTag.ID] = true
	tagNameIndex[slug] = newTag.ID
	return newTag.ID, nil
}
