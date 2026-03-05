package core

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

const vaultDir = ".typemd"

// Vault represents a typemd vault.
type Vault struct {
	Root string
	db   *sql.DB
}

// NewVault creates a Vault rooted at the given directory.
func NewVault(root string) *Vault {
	return &Vault{Root: root}
}

// Dir returns the vault metadata directory path.
func (v *Vault) Dir() string {
	return filepath.Join(v.Root, vaultDir)
}

// TypesDir returns the types directory path.
func (v *Vault) TypesDir() string {
	return filepath.Join(v.Dir(), "types")
}

// DBPath returns the SQLite database path.
func (v *Vault) DBPath() string {
	return filepath.Join(v.Dir(), "index.db")
}

// ObjectsDir returns the objects directory path.
func (v *Vault) ObjectsDir() string {
	return filepath.Join(v.Root, "objects")
}

// ObjectDir returns the directory path for a specific type's objects.
func (v *Vault) ObjectDir(typeName string) string {
	return filepath.Join(v.ObjectsDir(), typeName)
}

// ObjectPath returns the file path for a specific object.
func (v *Vault) ObjectPath(typeName, filename string) string {
	return filepath.Join(v.ObjectDir(typeName), filename+".md")
}

// Open opens the SQLite database connection and ensures the schema exists.
func (v *Vault) Open() error {
	if v.db != nil {
		return fmt.Errorf("vault already opened")
	}
	if !v.IsInitialized() {
		return fmt.Errorf("vault not initialized at %s", v.Dir())
	}

	db, err := sql.Open("sqlite", v.DBPath())
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	v.db = db

	if err := v.ensureSchema(); err != nil {
		v.db.Close()
		v.db = nil
		return fmt.Errorf("ensure schema: %w", err)
	}

	sync, err := v.needsSync()
	if err != nil {
		v.db.Close()
		v.db = nil
		return fmt.Errorf("check index: %w", err)
	}
	if sync {
		if err := v.SyncIndex(); err != nil {
			v.db.Close()
			v.db = nil
			return fmt.Errorf("auto sync index: %w", err)
		}
	}

	return nil
}

// needsSync returns true if the objects table has zero rows.
func (v *Vault) needsSync() (bool, error) {
	var count int
	if err := v.db.QueryRow("SELECT COUNT(*) FROM objects").Scan(&count); err != nil {
		return false, err
	}
	return count == 0, nil
}

// Close closes the SQLite database connection.
func (v *Vault) Close() error {
	if v.db == nil {
		return nil
	}
	err := v.db.Close()
	v.db = nil
	return err
}

// IsInitialized returns true if the vault has already been initialized.
func (v *Vault) IsInitialized() bool {
	_, err := os.Stat(v.Dir())
	return !errors.Is(err, os.ErrNotExist)
}

// Init initializes the vault: creates directories and SQLite schema.
func (v *Vault) Init() error {
	if v.IsInitialized() {
		return fmt.Errorf("vault already initialized at %s", v.Dir())
	}

	// Create directories
	for _, dir := range []string{v.TypesDir(), v.ObjectsDir()} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("create directory %s: %w", dir, err)
		}
	}

	// Create .gitignore to exclude index.db
	gitignore := filepath.Join(v.Dir(), ".gitignore")
	if err := os.WriteFile(gitignore, []byte("index.db\n"), 0644); err != nil {
		return fmt.Errorf("create .gitignore: %w", err)
	}

	// Initialize SQLite DB
	db, err := sql.Open("sqlite", v.DBPath())
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer db.Close()

	v.db = db
	err = v.ensureSchema()
	v.db = nil
	if err != nil {
		return fmt.Errorf("initialize schema: %w", err)
	}

	return nil
}

// ensureSchema creates tables and indexes if they don't exist.
func (v *Vault) ensureSchema() error {
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
	_, err := v.db.Exec(schema)
	return err
}
