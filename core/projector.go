package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Projector synchronizes the file-based source of truth (ObjectRepository)
// into the search index (ObjectIndex). It reads objects from the repository,
// applies migrations, and upserts them into the index.
type Projector struct {
	repo  ObjectRepository
	index ObjectIndex

	// createTag is a callback for auto-creating tag objects during sync.
	// Provided by Vault to reuse NewObject logic (ULID generation, uniqueness, etc.).
	createTag func(slug string) (*Object, error)
}

// NewProjector creates a Projector that syncs from repo to index.
// createTag is called when a tag reference needs auto-creation.
func NewProjector(repo ObjectRepository, index ObjectIndex, createTag func(slug string) (*Object, error)) *Projector {
	return &Projector{repo: repo, index: index, createTag: createTag}
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

// upsertObject processes a single object: migrates name if needed, filters properties,
// and upserts into the index. Returns the filtered properties JSON and any error.
func (s *objectSyncer) upsertObject(obj *Object, index ObjectIndex) error {
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
			return fmt.Errorf("write name migration for %s: %w", obj.ID, err)
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
		return nil // skip unserializable
	}

	return index.Upsert(obj.ID, obj.Type, obj.Filename, string(propsJSON), obj.Body)
}

// buildSyncContext builds a syncContext from a list of objects for wikilink/tag sync.
func buildSyncContext(objects []*Object) *syncContext {
	ctx := &syncContext{
		diskIDs:     make(map[string]bool),
		diskBodies:  make(map[string]string),
		diskTags:    make(map[string]*Object),
		diskTagRefs: make(map[string][]string),
	}
	for _, obj := range objects {
		ctx.diskIDs[obj.ID] = true
		ctx.diskBodies[obj.ID] = obj.Body
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
	return ctx
}

// syncWikiLinksAndTags syncs wikilinks and tag relations from a syncContext.
func (p *Projector) syncWikiLinksAndTags(ctx *syncContext) error {
	for id, body := range ctx.diskBodies {
		if err := p.syncWikiLinks(id, body, ctx.diskIDs); err != nil {
			return fmt.Errorf("sync wikilinks for %s: %w", id, err)
		}
	}
	return p.syncTagRelations(ctx)
}

// Sync scans the repository, upserts all objects into the index,
// removes stale entries, cleans up orphaned relations, syncs wikilinks
// and tag relations, and rebuilds the FTS index.
func (p *Projector) Sync() (*SyncResult, error) {
	result := &SyncResult{}

	// Walk all objects from repository
	objects, err := p.repo.Walk()
	if err != nil {
		return nil, fmt.Errorf("walk objects: %w", err)
	}

	// If no objects directory exists, Walk returns nil — clear the index
	if objects == nil {
		ids, err := p.index.ListIDs()
		if err != nil {
			return nil, fmt.Errorf("list indexed objects: %w", err)
		}
		for _, id := range ids {
			if err := p.index.Remove(id); err != nil {
				return nil, fmt.Errorf("clean object %s: %w", id, err)
			}
		}
		return result, p.index.Rebuild()
	}

	syncer := newObjectSyncer(p.repo)

	for _, obj := range objects {
		if err := syncer.upsertObject(obj, p.index); err != nil {
			return nil, err
		}
	}

	ctx := buildSyncContext(objects)

	// Delete stale objects from index
	deleted, err := p.deleteStaleObjects(ctx.diskIDs)
	if err != nil {
		return nil, err
	}
	result.Deleted = len(deleted)

	// Clean orphaned relations
	orphaned, err := p.index.CleanOrphanedRelations()
	if err != nil {
		return nil, err
	}
	result.Orphaned = orphaned

	// Clean up wikilinks for deleted objects
	for _, id := range deleted {
		if err := p.index.DeleteWikiLinks(id); err != nil {
			return nil, fmt.Errorf("delete wikilinks for %s: %w", id, err)
		}
	}

	// Sync wikilinks and tag relations
	if err := p.syncWikiLinksAndTags(ctx); err != nil {
		return nil, err
	}

	return result, p.index.Rebuild()
}

// objectPathToID converts a filesystem path under the objects directory to an object ID.
// Path format: .../objects/<type>/<filename>.md → ID: <type>/<filename>
// Returns the ID and true if valid, or empty string and false otherwise.
func objectPathToID(path, objectsDir string) (string, bool) {
	rel, err := filepath.Rel(objectsDir, path)
	if err != nil {
		return "", false
	}
	// Must be exactly <type>/<filename>.md (two path components)
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

// SyncFiles incrementally synchronizes specific files to the index.
// For existing files, it reads and upserts them. For deleted files, it removes them.
// After incremental object sync, it performs a full wikilink and tag sync
// using data already in the index (no filesystem walk).
func (p *Projector) SyncFiles(paths []string, objectsDir string) (*SyncResult, error) {
	result := &SyncResult{}
	syncer := newObjectSyncer(p.repo)

	for _, path := range paths {
		id, ok := objectPathToID(path, objectsDir)
		if !ok {
			continue // skip non-object files
		}

		// Try to read the object — if it doesn't exist, it was deleted
		obj, err := p.repo.Get(id)
		if err != nil {
			if os.IsNotExist(err) || strings.Contains(err.Error(), "read object") {
				// File deleted — remove from index
				if removeErr := p.index.Remove(id); removeErr != nil {
					return nil, fmt.Errorf("remove deleted object %s: %w", id, removeErr)
				}
				if err := p.index.DeleteWikiLinks(id); err != nil {
					return nil, fmt.Errorf("delete wikilinks for %s: %w", id, err)
				}
				result.Deleted++
				continue
			}
			return nil, fmt.Errorf("read object %s: %w", id, err)
		}

		if err := syncer.upsertObject(obj, p.index); err != nil {
			return nil, err
		}
	}

	// Full wikilink and tag sync using all objects from disk
	objects, err := p.repo.Walk()
	if err != nil {
		return nil, fmt.Errorf("walk objects for wikilink/tag sync: %w", err)
	}

	if objects != nil {
		ctx := buildSyncContext(objects)
		if err := p.syncWikiLinksAndTags(ctx); err != nil {
			return nil, err
		}
	}

	return result, nil
}

// deleteStaleObjects removes index entries for objects not found on disk.
func (p *Projector) deleteStaleObjects(diskIDs map[string]bool) ([]string, error) {
	indexedIDs, err := p.index.ListIDs()
	if err != nil {
		return nil, fmt.Errorf("list indexed objects: %w", err)
	}

	var toDelete []string
	for _, id := range indexedIDs {
		if !diskIDs[id] {
			toDelete = append(toDelete, id)
		}
	}

	for _, id := range toDelete {
		if err := p.index.Remove(id); err != nil {
			return nil, fmt.Errorf("delete stale object %s: %w", id, err)
		}
	}

	return toDelete, nil
}

// syncWikiLinks extracts wiki-links from body and stores them in the index.
func (p *Projector) syncWikiLinks(objectID, body string, knownIDs map[string]bool) error {
	links := ParseWikiLinks(body)
	if len(links) == 0 {
		return p.index.SyncWikiLinks(objectID, nil)
	}
	entries := make([]WikiLinkEntry, len(links))
	for i, link := range links {
		toID := ""
		if knownIDs[link.Target] {
			toID = link.Target
		}
		entries[i] = WikiLinkEntry{
			ToID:        toID,
			Target:      link.Target,
			DisplayText: link.DisplayText,
		}
	}
	return p.index.SyncWikiLinks(objectID, entries)
}

// syncTagRelations clears existing tag relations and rebuilds them from frontmatter.
func (p *Projector) syncTagRelations(ctx *syncContext) error {
	if err := p.index.DeleteRelationsByName(TagsProperty); err != nil {
		return fmt.Errorf("clear tag relations: %w", err)
	}

	tagNameIndex := make(map[string]string)
	for _, obj := range ctx.diskTags {
		if name, ok := obj.Properties[NameProperty].(string); ok {
			tagNameIndex[name] = obj.ID
		}
	}

	for objID, refs := range ctx.diskTagRefs {
		for _, ref := range refs {
			tagID, err := p.resolveOrCreateTag(ref, ctx, tagNameIndex)
			if err != nil {
				continue // skip unresolvable tag references
			}
			if err := p.index.InsertRelation(TagsProperty, objID, tagID); err != nil {
				return fmt.Errorf("insert tag relation: %w", err)
			}
		}
	}

	return nil
}

// resolveOrCreateTag resolves a tag reference to an object ID, auto-creating if needed.
func (p *Projector) resolveOrCreateTag(ref string, ctx *syncContext, tagNameIndex map[string]string) (string, error) {
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

	if p.createTag == nil {
		return "", fmt.Errorf("cannot auto-create tag %q: no createTag callback", slug)
	}

	newTag, err := p.createTag(slug)
	if err != nil {
		return "", fmt.Errorf("auto-create tag %q: %w", slug, err)
	}
	ctx.diskTags[newTag.ID] = newTag
	ctx.diskIDs[newTag.ID] = true
	tagNameIndex[slug] = newTag.ID
	return newTag.ID, nil
}
