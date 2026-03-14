package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// OrphanedRelation represents a relation record that references a non-existent object.
type OrphanedRelation struct {
	Name   string
	FromID string
	ToID   string
}

// SyncResult holds statistics from a SyncIndex operation.
type SyncResult struct {
	Deleted  int
	Orphaned []OrphanedRelation
}

// syncContext holds intermediate state collected during SyncIndex.
type syncContext struct {
	diskIDs     map[string]bool
	diskBodies  map[string]string
	diskTags    map[string]*Object
	diskTagRefs map[string][]string
}

// SyncIndex scans the objects directory, upserts all found objects into the DB,
// removes DB entries for deleted files, cleans up orphaned relations, and rebuilds the FTS index.
func (v *Vault) SyncIndex() (*SyncResult, error) {
	if v.db == nil {
		return nil, fmt.Errorf("vault not opened")
	}

	result := &SyncResult{}

	// Handle empty vault (no objects directory)
	objsDir := v.ObjectsDir()
	if _, err := os.Stat(objsDir); os.IsNotExist(err) {
		if _, err := v.db.Exec("DELETE FROM objects"); err != nil {
			return nil, fmt.Errorf("clean objects: %w", err)
		}
		return result, v.RebuildIndex()
	}

	ctx, err := v.walkAndUpsertObjects()
	if err != nil {
		return nil, err
	}

	deleted, err := v.deleteStaleObjects(ctx.diskIDs)
	if err != nil {
		return nil, err
	}
	result.Deleted = len(deleted)

	orphaned, err := v.cleanOrphanedRelations()
	if err != nil {
		return nil, err
	}
	result.Orphaned = orphaned

	// Clean up wikilinks for deleted objects and sync current ones
	for _, id := range deleted {
		if _, err := v.db.Exec("DELETE FROM wikilinks WHERE from_id = ?", id); err != nil {
			return nil, fmt.Errorf("delete wikilinks for %s: %w", id, err)
		}
	}
	for id, body := range ctx.diskBodies {
		if err := v.syncWikiLinks(id, body, ctx.diskIDs); err != nil {
			return nil, fmt.Errorf("sync wikilinks for %s: %w", id, err)
		}
	}

	if err := v.syncTagRelations(ctx); err != nil {
		return nil, err
	}

	return result, v.RebuildIndex()
}

// walkAndUpsertObjects walks the objects directory, parses each object file,
// and upserts it into the database. Returns collected state for downstream steps.
func (v *Vault) walkAndUpsertObjects() (*syncContext, error) {
	ctx := &syncContext{
		diskIDs:     make(map[string]bool),
		diskBodies:  make(map[string]string),
		diskTags:    make(map[string]*Object),
		diskTagRefs: make(map[string][]string),
	}

	schemaCache := make(map[string]*TypeSchema)
	propertyNameCache := make(map[string]map[string]bool)
	sysNames := SystemPropertyNames()

	objsDir := v.ObjectsDir()
	err := filepath.Walk(objsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		rel, err := filepath.Rel(objsDir, path)
		if err != nil {
			return nil
		}

		parts := strings.SplitN(rel, string(os.PathSeparator), 2)
		if len(parts) != 2 {
			return nil
		}
		typeName := parts[0]
		filename := strings.TrimSuffix(parts[1], ".md")
		id := typeName + "/" + filename

		data, err := os.ReadFile(path)
		if err != nil {
			return nil // skip unreadable files
		}

		props, body, err := parseFrontmatter(data)
		if err != nil {
			return nil // skip unparseable files
		}

		// Populate schema cache
		if _, cached := schemaCache[typeName]; !cached {
			schema, err := v.LoadType(typeName)
			if err != nil {
				schemaCache[typeName] = nil
			} else {
				schemaCache[typeName] = schema
				propertyNameCache[typeName] = schema.PropertyNames()
			}
		}

		// Migrate: add NameProperty if missing
		nameVal, hasName := props[NameProperty]
		if !hasName || nameVal == nil || nameVal == "" {
			props[NameProperty] = StripULID(filename)
			if updated, err := writeFrontmatter(props, body, OrderedPropKeys(props, schemaCache[typeName])); err == nil {
				if writeErr := os.WriteFile(path, updated, 0644); writeErr != nil {
					return fmt.Errorf("write name migration for %s: %w", id, writeErr)
				}
			}
		}

		// Filter properties by type schema (only index schema-defined keys + system properties)
		if allowed := propertyNameCache[typeName]; allowed != nil {
			filtered := make(map[string]any, len(allowed)+len(sysNames))
			for _, name := range sysNames {
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
			return nil
		}

		_, err = v.db.Exec(`
			INSERT INTO objects (id, type, filename, properties, body)
			VALUES (?, ?, ?, ?, ?)
			ON CONFLICT(id) DO UPDATE SET
				type = excluded.type,
				filename = excluded.filename,
				properties = excluded.properties,
				body = excluded.body
		`, id, typeName, filename, string(propsJSON), body)
		if err != nil {
			return fmt.Errorf("upsert object %s: %w", id, err)
		}

		ctx.diskIDs[id] = true
		ctx.diskBodies[id] = body

		if typeName == TagTypeName {
			ctx.diskTags[id] = &Object{
				ID:         id,
				Type:       typeName,
				Filename:   filename,
				Properties: props,
			}
		}

		if tagsVal, ok := props[TagsProperty]; ok {
			if tagsArr, ok := tagsVal.([]any); ok {
				var refs []string
				for _, item := range tagsArr {
					if ref, ok := item.(string); ok {
						refs = append(refs, ref)
					}
				}
				if len(refs) > 0 {
					ctx.diskTagRefs[id] = refs
				}
			}
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk objects: %w", err)
	}

	return ctx, nil
}

// deleteStaleObjects removes DB entries for objects that no longer exist on disk.
// Returns the list of deleted object IDs.
func (v *Vault) deleteStaleObjects(diskIDs map[string]bool) ([]string, error) {
	rows, err := v.db.Query("SELECT id FROM objects")
	if err != nil {
		return nil, fmt.Errorf("list db objects: %w", err)
	}
	defer rows.Close()

	var toDelete []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan id: %w", err)
		}
		if !diskIDs[id] {
			toDelete = append(toDelete, id)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate db objects: %w", err)
	}

	for _, id := range toDelete {
		if _, err := v.db.Exec("DELETE FROM objects WHERE id = ?", id); err != nil {
			return nil, fmt.Errorf("delete stale object %s: %w", id, err)
		}
	}

	return toDelete, nil
}

// cleanOrphanedRelations detects and removes relation records that reference
// non-existent objects. Returns the list of orphaned relations found.
func (v *Vault) cleanOrphanedRelations() ([]OrphanedRelation, error) {
	rows, err := v.db.Query(`
		SELECT r.id, r.name, r.from_id, r.to_id FROM relations r
		LEFT JOIN objects o1 ON r.from_id = o1.id
		LEFT JOIN objects o2 ON r.to_id = o2.id
		WHERE o1.id IS NULL OR o2.id IS NULL
	`)
	if err != nil {
		return nil, fmt.Errorf("detect orphaned relations: %w", err)
	}
	defer rows.Close()

	var orphaned []OrphanedRelation
	var orphanIDs []any
	for rows.Next() {
		var rowID int
		var o OrphanedRelation
		if err := rows.Scan(&rowID, &o.Name, &o.FromID, &o.ToID); err != nil {
			return nil, fmt.Errorf("scan orphaned relation: %w", err)
		}
		orphaned = append(orphaned, o)
		orphanIDs = append(orphanIDs, rowID)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate orphaned relations: %w", err)
	}

	if len(orphanIDs) > 0 {
		placeholders := strings.Repeat("?,", len(orphanIDs))
		_, err := v.db.Exec(
			"DELETE FROM relations WHERE id IN ("+placeholders[:len(placeholders)-1]+")",
			orphanIDs...,
		)
		if err != nil {
			return nil, fmt.Errorf("delete orphaned relations: %w", err)
		}
	}

	return orphaned, nil
}

// syncTagRelations clears existing tag relations and rebuilds them from frontmatter references.
// It auto-creates tag objects for name references that don't match existing tags.
func (v *Vault) syncTagRelations(ctx *syncContext) error {
	if _, err := v.db.Exec("DELETE FROM relations WHERE name = ?", TagsProperty); err != nil {
		return fmt.Errorf("clear tag relations: %w", err)
	}

	// Build a name->ID index from diskTags for quick lookups
	tagNameIndex := make(map[string]string)
	for _, obj := range ctx.diskTags {
		if name, ok := obj.Properties[NameProperty].(string); ok {
			tagNameIndex[name] = obj.ID
		}
	}

	for objID, refs := range ctx.diskTagRefs {
		for _, ref := range refs {
			tagID, ok := v.resolveTagReference(ref, ctx.diskTags, tagNameIndex)
			if !ok {
				slug := strings.TrimPrefix(ref, "tag/")
				if !ulidSuffixPattern.MatchString(slug) {
					if existingID, exists := tagNameIndex[slug]; exists {
						tagID = existingID
					} else {
						newTag, err := v.NewObject(TagTypeName, slug, "")
						if err != nil {
							continue
						}
						ctx.diskTags[newTag.ID] = newTag
						ctx.diskIDs[newTag.ID] = true
						tagNameIndex[slug] = newTag.ID
						tagID = newTag.ID
					}
				} else {
					continue
				}
			}
			_, err := v.db.Exec(
				"INSERT INTO relations (name, from_id, to_id) VALUES (?, ?, ?)",
				TagsProperty, objID, tagID,
			)
			if err != nil {
				return fmt.Errorf("insert tag relation: %w", err)
			}
		}
	}

	return nil
}
