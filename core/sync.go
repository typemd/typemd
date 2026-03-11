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
	Created  int
	Updated  int
	Deleted  int
	Orphaned []OrphanedRelation
}

// SyncIndex scans the objects directory, upserts all found objects into the DB,
// removes DB entries for deleted files, cleans up orphaned relations, and rebuilds the FTS index.
func (v *Vault) SyncIndex() (*SyncResult, error) {
	if v.db == nil {
		return nil, fmt.Errorf("vault not opened")
	}

	result := &SyncResult{}

	// Collect all object IDs found on disk, along with their bodies for wiki-link extraction
	diskIDs := make(map[string]bool)
	diskBodies := make(map[string]string)

	// Cache loaded type schemas and their property name sets per type
	schemaCache := make(map[string]*TypeSchema)
	propertyNameCache := make(map[string]map[string]bool)
	sysNames := SystemPropertyNames()

	objsDir := v.ObjectsDir()
	if _, err := os.Stat(objsDir); os.IsNotExist(err) {
		// No objects directory — just clean DB and return
		_, err := v.db.Exec("DELETE FROM objects")
		if err != nil {
			return nil, fmt.Errorf("clean objects: %w", err)
		}
		return result, v.RebuildIndex()
	}

	// Walk objects/<type>/<name>.md
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

		// Read and parse file
		data, err := os.ReadFile(path)
		if err != nil {
			return nil // skip unreadable files
		}

		props, body, err := parseFrontmatter(data)
		if err != nil {
			return nil // skip unparseable files
		}

		// Populate schema cache (used by both name migration and property filtering)
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
			// Preserve all system properties
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

		// Upsert into DB
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

		diskIDs[id] = true
		diskBodies[id] = body
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk objects: %w", err)
	}

	// Remove DB entries that no longer exist on disk
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
	result.Deleted = len(toDelete)

	// Detect and clean up orphaned relations
	orphanRows, err := v.db.Query(`
		SELECT r.name, r.from_id, r.to_id FROM relations r
		LEFT JOIN objects o1 ON r.from_id = o1.id
		LEFT JOIN objects o2 ON r.to_id = o2.id
		WHERE o1.id IS NULL OR o2.id IS NULL
	`)
	if err != nil {
		return nil, fmt.Errorf("detect orphaned relations: %w", err)
	}
	defer orphanRows.Close()

	for orphanRows.Next() {
		var o OrphanedRelation
		if err := orphanRows.Scan(&o.Name, &o.FromID, &o.ToID); err != nil {
			return nil, fmt.Errorf("scan orphaned relation: %w", err)
		}
		result.Orphaned = append(result.Orphaned, o)
	}
	if err := orphanRows.Err(); err != nil {
		return nil, fmt.Errorf("iterate orphaned relations: %w", err)
	}

	// Delete orphaned relations from DB
	if len(result.Orphaned) > 0 {
		_, err := v.db.Exec(`
			DELETE FROM relations WHERE id IN (
				SELECT r.id FROM relations r
				LEFT JOIN objects o1 ON r.from_id = o1.id
				LEFT JOIN objects o2 ON r.to_id = o2.id
				WHERE o1.id IS NULL OR o2.id IS NULL
			)
		`)
		if err != nil {
			return nil, fmt.Errorf("delete orphaned relations: %w", err)
		}
	}

	// Clean up wikilinks for deleted objects
	for _, id := range toDelete {
		if _, err := v.db.Exec("DELETE FROM wikilinks WHERE from_id = ?", id); err != nil {
			return nil, fmt.Errorf("delete wikilinks for %s: %w", id, err)
		}
	}

	// Sync wiki-links using in-memory diskIDs for target resolution
	for id, body := range diskBodies {
		if err := v.syncWikiLinks(id, body, diskIDs); err != nil {
			return nil, fmt.Errorf("sync wikilinks for %s: %w", id, err)
		}
	}

	return result, v.RebuildIndex()
}
