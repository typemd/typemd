package core

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
)

// SQLiteObjectIndex implements ObjectIndex using SQLite with FTS5.
type SQLiteObjectIndex struct {
	db *sql.DB
}

// NewSQLiteObjectIndex creates a new SQLiteObjectIndex wrapping the given database connection.
func NewSQLiteObjectIndex(db *sql.DB) *SQLiteObjectIndex {
	return &SQLiteObjectIndex{db: db}
}

// EnsureSchema creates tables, indexes, and FTS triggers if they don't exist.
func (idx *SQLiteObjectIndex) EnsureSchema() error {
	schema := `
CREATE TABLE IF NOT EXISTS objects (
    id         TEXT PRIMARY KEY,
    type       TEXT NOT NULL,
    filename   TEXT NOT NULL,
    properties TEXT NOT NULL DEFAULT '{}',
    body       TEXT NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS relations (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    name       TEXT NOT NULL,
    from_id    TEXT NOT NULL,
    to_id      TEXT NOT NULL,
    FOREIGN KEY (from_id) REFERENCES objects(id),
    FOREIGN KEY (to_id)   REFERENCES objects(id)
);

CREATE INDEX IF NOT EXISTS idx_objects_type ON objects(type);
CREATE INDEX IF NOT EXISTS idx_relations_from ON relations(from_id);
CREATE INDEX IF NOT EXISTS idx_relations_to   ON relations(to_id);

CREATE TABLE IF NOT EXISTS wikilinks (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    from_id      TEXT NOT NULL,
    to_id        TEXT NOT NULL DEFAULT '',
    target       TEXT NOT NULL,
    display_text TEXT NOT NULL DEFAULT '',
    FOREIGN KEY (from_id) REFERENCES objects(id)
);

CREATE INDEX IF NOT EXISTS idx_wikilinks_from ON wikilinks(from_id);
CREATE INDEX IF NOT EXISTS idx_wikilinks_to   ON wikilinks(to_id);

CREATE VIRTUAL TABLE IF NOT EXISTS objects_fts USING fts5(
    filename,
    properties,
    body,
    content='objects',
    content_rowid='rowid'
);

CREATE TRIGGER IF NOT EXISTS objects_ai AFTER INSERT ON objects BEGIN
    INSERT INTO objects_fts(rowid, filename, properties, body)
    VALUES (new.rowid, new.filename, new.properties, new.body);
END;

CREATE TRIGGER IF NOT EXISTS objects_ad AFTER DELETE ON objects BEGIN
    INSERT INTO objects_fts(objects_fts, rowid, filename, properties, body)
    VALUES ('delete', old.rowid, old.filename, old.properties, old.body);
END;

CREATE TRIGGER IF NOT EXISTS objects_au AFTER UPDATE ON objects BEGIN
    INSERT INTO objects_fts(objects_fts, rowid, filename, properties, body)
    VALUES ('delete', old.rowid, old.filename, old.properties, old.body);
    INSERT INTO objects_fts(rowid, filename, properties, body)
    VALUES (new.rowid, new.filename, new.properties, new.body);
END;
`
	_, err := idx.db.Exec(schema)
	return err
}

// NeedsSync returns true if the objects table has zero rows.
func (idx *SQLiteObjectIndex) NeedsSync() (bool, error) {
	var count int
	if err := idx.db.QueryRow("SELECT COUNT(*) FROM objects").Scan(&count); err != nil {
		return false, err
	}
	return count == 0, nil
}

// scanObjectResults scans SQL rows into ObjectResult slices.
func scanObjectResults(rows *sql.Rows) ([]*ObjectResult, error) {
	var results []*ObjectResult
	for rows.Next() {
		var r ObjectResult
		var propsJSON string
		if err := rows.Scan(&r.ID, &r.Type, &r.Filename, &propsJSON, &r.Body); err != nil {
			return nil, fmt.Errorf("scan object result: %w", err)
		}
		if err := json.Unmarshal([]byte(propsJSON), &r.Properties); err != nil {
			return nil, fmt.Errorf("unmarshal properties: %w", err)
		}
		results = append(results, &r)
	}
	return results, rows.Err()
}

// Query queries objects using key=value filter syntax.
// "type" filters on the type column; other keys filter on JSON properties.
// An empty filter returns all objects.
func (idx *SQLiteObjectIndex) Query(filter string) ([]*ObjectResult, error) {
	conditions, err := parseFilter(filter)
	if err != nil {
		return nil, fmt.Errorf("parse filter: %w", err)
	}

	query := "SELECT id, type, filename, properties, body FROM objects"
	var whereClauses []string
	var args []any

	for _, c := range conditions {
		if c.Key == "type" {
			whereClauses = append(whereClauses, "type = ?")
			args = append(args, c.Value)
		} else {
			whereClauses = append(whereClauses, "json_extract(properties, ?) = ?")
			args = append(args, "$."+c.Key, c.Value)
		}
	}

	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	rows, err := idx.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("query objects: %w", err)
	}
	defer rows.Close()

	return scanObjectResults(rows)
}

// Search performs full-text search using FTS5.
// Returns nil, nil for empty keyword.
func (idx *SQLiteObjectIndex) Search(keyword string) ([]*ObjectResult, error) {
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return nil, nil
	}

	query := `SELECT o.id, o.type, o.filename, o.properties, o.body
		FROM objects o
		JOIN objects_fts fts ON o.rowid = fts.rowid
		WHERE objects_fts MATCH ?`

	rows, err := idx.db.Query(query, keyword)
	if err != nil {
		return nil, fmt.Errorf("search objects: %w", err)
	}
	defer rows.Close()

	return scanObjectResults(rows)
}

// Rebuild rebuilds the FTS5 index from the objects table.
func (idx *SQLiteObjectIndex) Rebuild() error {
	_, err := idx.db.Exec("INSERT INTO objects_fts(objects_fts) VALUES('rebuild')")
	if err != nil {
		return fmt.Errorf("rebuild index: %w", err)
	}
	return nil
}

// FindRelations returns all relations where objectID is either the source or target.
func (idx *SQLiteObjectIndex) FindRelations(objectID string) ([]Relation, error) {
	rows, err := idx.db.Query(
		"SELECT name, from_id, to_id FROM relations WHERE from_id = ? OR to_id = ?",
		objectID, objectID,
	)
	if err != nil {
		return nil, fmt.Errorf("find relations: %w", err)
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

// InsertRelation stores a new relation record in the index.
func (idx *SQLiteObjectIndex) InsertRelation(name, fromID, toID string) error {
	_, err := idx.db.Exec(
		"INSERT INTO relations (name, from_id, to_id) VALUES (?, ?, ?)",
		name, fromID, toID,
	)
	if err != nil {
		return fmt.Errorf("insert relation: %w", err)
	}
	return nil
}

// DeleteRelation removes a specific relation record from the index.
func (idx *SQLiteObjectIndex) DeleteRelation(name, fromID, toID string) error {
	_, err := idx.db.Exec(
		"DELETE FROM relations WHERE name = ? AND from_id = ? AND to_id = ?",
		name, fromID, toID,
	)
	if err != nil {
		return fmt.Errorf("delete relation: %w", err)
	}
	return nil
}

// DeleteRelationsByName removes all relation records with the given name.
func (idx *SQLiteObjectIndex) DeleteRelationsByName(name string) error {
	_, err := idx.db.Exec("DELETE FROM relations WHERE name = ?", name)
	if err != nil {
		return fmt.Errorf("delete relations by name: %w", err)
	}
	return nil
}

// CleanOrphanedRelations detects and removes relation records that reference
// non-existent objects. Returns the list of orphaned relations found.
func (idx *SQLiteObjectIndex) CleanOrphanedRelations() ([]OrphanedRelation, error) {
	rows, err := idx.db.Query(`
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
		_, err := idx.db.Exec(
			"DELETE FROM relations WHERE id IN ("+placeholders[:len(placeholders)-1]+")",
			orphanIDs...,
		)
		if err != nil {
			return nil, fmt.Errorf("delete orphaned relations: %w", err)
		}
	}

	return orphaned, nil
}

// ListWikiLinks returns all wiki-links from the given object.
func (idx *SQLiteObjectIndex) ListWikiLinks(objectID string) ([]StoredWikiLink, error) {
	return idx.queryWikiLinks("from_id", objectID)
}

// FindBacklinks returns all wiki-links pointing to the given object.
func (idx *SQLiteObjectIndex) FindBacklinks(objectID string) ([]StoredWikiLink, error) {
	return idx.queryWikiLinks("to_id", objectID)
}

// queryWikiLinks is a shared helper for querying wikilinks by a specific column.
func (idx *SQLiteObjectIndex) queryWikiLinks(column, objectID string) ([]StoredWikiLink, error) {
	query := fmt.Sprintf("SELECT from_id, to_id, target, display_text FROM wikilinks WHERE %s = ?", column)
	rows, err := idx.db.Query(query, objectID)
	if err != nil {
		return nil, fmt.Errorf("query wikilinks by %s: %w", column, err)
	}
	defer rows.Close()

	var links []StoredWikiLink
	for rows.Next() {
		var l StoredWikiLink
		if err := rows.Scan(&l.FromID, &l.ToID, &l.Target, &l.DisplayText); err != nil {
			return nil, fmt.Errorf("scan wikilink: %w", err)
		}
		links = append(links, l)
	}
	return links, rows.Err()
}

// SyncWikiLinks replaces all wikilinks for the given object with the provided entries.
func (idx *SQLiteObjectIndex) SyncWikiLinks(fromID string, links []WikiLinkEntry) error {
	if _, err := idx.db.Exec("DELETE FROM wikilinks WHERE from_id = ?", fromID); err != nil {
		return fmt.Errorf("delete old wikilinks: %w", err)
	}

	if len(links) == 0 {
		return nil
	}

	stmt, err := idx.db.Prepare("INSERT INTO wikilinks (from_id, to_id, target, display_text) VALUES (?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("prepare wikilink insert: %w", err)
	}
	defer stmt.Close()

	for _, link := range links {
		if _, err := stmt.Exec(fromID, link.ToID, link.Target, link.DisplayText); err != nil {
			return fmt.Errorf("insert wikilink: %w", err)
		}
	}
	return nil
}

// DeleteWikiLinks removes all wikilink records for the given object.
func (idx *SQLiteObjectIndex) DeleteWikiLinks(fromID string) error {
	_, err := idx.db.Exec("DELETE FROM wikilinks WHERE from_id = ?", fromID)
	if err != nil {
		return fmt.Errorf("delete wikilinks: %w", err)
	}
	return nil
}

// Upsert inserts or updates an object record in the index.
func (idx *SQLiteObjectIndex) Upsert(id, typeName, filename, propsJSON, body string) error {
	_, err := idx.db.Exec(`
		INSERT INTO objects (id, type, filename, properties, body)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			type = excluded.type,
			filename = excluded.filename,
			properties = excluded.properties,
			body = excluded.body
	`, id, typeName, filename, propsJSON, body)
	if err != nil {
		return fmt.Errorf("upsert object %s: %w", id, err)
	}
	return nil
}

// Remove deletes an object record from the index.
func (idx *SQLiteObjectIndex) Remove(id string) error {
	_, err := idx.db.Exec("DELETE FROM objects WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("remove object %s: %w", id, err)
	}
	return nil
}

// ListIDs returns all object IDs currently in the index.
func (idx *SQLiteObjectIndex) ListIDs() ([]string, error) {
	rows, err := idx.db.Query("SELECT id FROM objects")
	if err != nil {
		return nil, fmt.Errorf("list object IDs: %w", err)
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan id: %w", err)
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}
