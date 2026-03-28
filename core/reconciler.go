package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

// Reconciler normalizes file content (name migration, prefix expansion,
// wiki-link expansion, tag auto-creation), detects changes by comparing
// disk state against the index, and emits domain events describing what changed.
type Reconciler struct {
	repo      ObjectRepository
	index     ObjectIndex // read-only: ListIDs() for stale detection
	createTag func(slug string) (*Object, error)
}

// NewReconciler creates a Reconciler.
func NewReconciler(repo ObjectRepository, index ObjectIndex, createTag func(slug string) (*Object, error)) *Reconciler {
	return &Reconciler{repo: repo, index: index, createTag: createTag}
}

// objectSyncer holds per-sync caches for processing objects.
type objectSyncer struct {
	repo              ObjectRepository
	schemaCache       map[string]*TypeSchema
	propertyNameCache map[string]map[string]bool
	sysNames          []string
}

func newObjectSyncer(repo ObjectRepository) *objectSyncer {
	return &objectSyncer{
		repo:              repo,
		schemaCache:       make(map[string]*TypeSchema),
		propertyNameCache: make(map[string]map[string]bool),
		sysNames:          SystemPropertyNames(),
	}
}

// processObject normalizes a single object: migrates name if needed, filters properties,
// and returns an ObjectUpserted event.
func (s *objectSyncer) processObject(obj *Object) (*ObjectUpserted, error) {
	// Populate schema cache
	if _, cached := s.schemaCache[obj.Type]; !cached {
		schema, err := s.repo.GetSchema(obj.Type)
		if err != nil {
			s.schemaCache[obj.Type] = nil
		} else {
			s.schemaCache[obj.Type] = schema
			s.propertyNameCache[obj.Type] = schema.PropertyNames()
		}
	}

	// Migrate: add NameProperty if missing
	nameVal, hasName := obj.Properties[NameProperty]
	if !hasName || nameVal == nil || nameVal == "" {
		obj.Properties[NameProperty] = StripULID(obj.Filename)
		keyOrder := OrderedPropKeys(obj.Properties, s.schemaCache[obj.Type])
		if err := s.repo.Save(obj, keyOrder); err != nil {
			return nil, fmt.Errorf("write name migration for %s: %w", obj.ID, err)
		}
	}

	// Filter properties by type schema (only index schema-defined keys + system properties)
	props := obj.Properties
	if allowed := s.propertyNameCache[obj.Type]; allowed != nil {
		filtered := make(map[string]any, len(allowed)+len(s.sysNames))
		for _, name := range s.sysNames {
			if val, ok := props[name]; ok {
				filtered[name] = val
			}
		}
		for k, val := range props {
			if allowed[k] {
				filtered[k] = val
			}
		}
		props = filtered
	}

	propsJSON, err := json.Marshal(props)
	if err != nil {
		slog.Warn("skip object: cannot marshal properties", "id", obj.ID, "error", err)
		return nil, nil // skip unserializable
	}

	return &ObjectUpserted{
		ID:        obj.ID,
		Type:      obj.Type,
		Filename:  obj.Filename,
		PropsJSON: string(propsJSON),
		Body:      obj.Body,
	}, nil
}

// buildSyncContext builds a syncContext from a list of objects for wikilink/tag/relation reconciliation.
func buildSyncContext(objects []*Object) *syncContext {
	ctx := &syncContext{
		diskIDs:     make(map[string]bool),
		diskBodies:  make(map[string]string),
		diskTags:    make(map[string]*Object),
		diskTagRefs: make(map[string][]string),
		diskObjects: make(map[string]*Object),
		schemas:     make(map[string]*TypeSchema),
		nameIndex:   make(map[string]map[string][]string),
	}
	for _, obj := range objects {
		ctx.diskIDs[obj.ID] = true
		ctx.diskBodies[obj.ID] = obj.Body
		ctx.diskObjects[obj.ID] = obj
		if obj.Type == TagTypeName {
			ctx.diskTags[obj.ID] = obj
		}
		if tagsVal, ok := obj.Properties[TagsProperty]; ok {
			if tagsArr, ok := tagsVal.([]any); ok {
				var refs []string
				for _, item := range tagsArr {
					if ref, ok := item.(string); ok {
						refs = append(refs, ref)
					}
				}
				if len(refs) > 0 {
					ctx.diskTagRefs[obj.ID] = refs
				}
			}
		}
	}
	buildNameIndex(ctx)
	return ctx
}

// objectPathToID converts a filesystem path under the objects directory to an object ID.
func objectPathToID(path, objectsDir string) (string, bool) {
	rel, err := filepath.Rel(objectsDir, path)
	if err != nil {
		return "", false
	}
	rel = filepath.ToSlash(rel)
	if !strings.HasSuffix(rel, ".md") {
		return "", false
	}
	parts := strings.SplitN(rel, "/", 3)
	if len(parts) != 2 {
		return "", false
	}
	id := parts[0] + "/" + strings.TrimSuffix(parts[1], ".md")
	return id, true
}

// deleteStaleObjects finds index entries for objects not found on disk
// and returns ObjectDeleted events.
func (r *Reconciler) deleteStaleObjects(diskIDs map[string]bool) ([]DomainEvent, error) {
	indexedIDs, err := r.index.ListIDs()
	if err != nil {
		return nil, fmt.Errorf("list indexed objects: %w", err)
	}

	var events []DomainEvent
	for _, id := range indexedIDs {
		if !diskIDs[id] {
			events = append(events, ObjectDeleted{ID: id})
		}
	}
	return events, nil
}

// Reconcile performs a full reconciliation: walks all objects, normalizes files,
// detects changes, and returns domain events for Projector to consume.
func (r *Reconciler) Reconcile() ([]DomainEvent, *ReconcileResult, error) {
	slog.Debug("reconcile started")
	result := &ReconcileResult{}
	var events []DomainEvent

	objects, err := r.repo.Walk()
	if err != nil {
		return nil, nil, fmt.Errorf("walk objects: %w", err)
	}

	// If no objects directory exists, emit delete events for all indexed objects
	if objects == nil {
		indexedIDs, err := r.index.ListIDs()
		if err != nil {
			return nil, nil, fmt.Errorf("list indexed objects: %w", err)
		}
		for _, id := range indexedIDs {
			events = append(events, ObjectDeleted{ID: id})
		}
		return events, result, nil
	}

	syncer := newObjectSyncer(r.repo)

	for _, obj := range objects {
		ev, err := syncer.processObject(obj)
		if err != nil {
			return nil, nil, err
		}
		if ev != nil {
			events = append(events, *ev)
			result.Synced++
		}
	}

	ctx := buildSyncContext(objects)
	ctx.schemas = syncer.schemaCache

	// Stale detection
	staleEvents, err := r.deleteStaleObjects(ctx.diskIDs)
	if err != nil {
		return nil, nil, err
	}
	events = append(events, staleEvents...)
	result.Deleted = len(staleEvents)

	// Relation reconciliation
	relationEvents, err := r.reconcileRelations(ctx, result)
	if err != nil {
		return nil, nil, err
	}
	events = append(events, relationEvents...)

	// WikiLink + tag reconciliation
	wikiTagEvents, err := r.reconcileWikiLinksAndTags(ctx, result)
	if err != nil {
		return nil, nil, err
	}
	events = append(events, wikiTagEvents...)

	slog.Debug("reconcile completed", "synced", result.Synced, "deleted", result.Deleted)
	return events, result, nil
}

// ReconcileFiles incrementally reconciles specific files.
func (r *Reconciler) ReconcileFiles(paths []string, objectsDir string) ([]DomainEvent, *ReconcileResult, error) {
	slog.Debug("incremental reconcile started", "files", len(paths))
	result := &ReconcileResult{}
	var events []DomainEvent
	syncer := newObjectSyncer(r.repo)

	for _, path := range paths {
		id, ok := objectPathToID(path, objectsDir)
		if !ok {
			continue
		}

		obj, err := r.repo.Get(id)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) || strings.Contains(err.Error(), "read object") {
				events = append(events, ObjectDeleted{ID: id})
				result.Deleted++
				continue
			}
			return nil, nil, fmt.Errorf("read object %s: %w", id, err)
		}

		ev, err := syncer.processObject(obj)
		if err != nil {
			return nil, nil, err
		}
		if ev != nil {
			events = append(events, *ev)
			result.Synced++
		}
	}

	// Full walk for relation/wikilink/tag reconciliation (needs global name index)
	objects, err := r.repo.Walk()
	if err != nil {
		return nil, nil, fmt.Errorf("walk objects for reconciliation: %w", err)
	}

	if objects != nil {
		ctx := buildSyncContext(objects)
		ctx.schemas = syncer.schemaCache

		// Incremental relation reconciliation: only for changed objects
		var changedIDs []string
		for _, path := range paths {
			id, ok := objectPathToID(path, objectsDir)
			if !ok {
				continue
			}
			changedIDs = append(changedIDs, id)
		}
		relationEvents, err := r.reconcileRelationsForObjects(ctx, result, changedIDs)
		if err != nil {
			return nil, nil, err
		}
		events = append(events, relationEvents...)

		wikiTagEvents, err := r.reconcileWikiLinksAndTags(ctx, result)
		if err != nil {
			return nil, nil, err
		}
		events = append(events, wikiTagEvents...)
	}

	return events, result, nil
}
