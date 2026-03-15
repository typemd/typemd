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
// It serves as a facade and DI container, delegating to focused services.
type Vault struct {
	Root              string
	db                *sql.DB
	index             ObjectIndex
	repo              ObjectRepository
	projector         *Projector
	Objects           *ObjectService
	Queries           *QueryService
	Events            *EventDispatcher
	sharedProperties  []Property
	sharedPropsMap    map[string]Property
	sharedPropsLoaded bool
}

// NewVault creates a Vault rooted at the given directory.
func NewVault(root string) *Vault {
	return &Vault{
		Root: root,
		repo: NewLocalObjectRepository(root),
	}
}

// Dir returns the vault metadata directory path.
func (v *Vault) Dir() string {
	return filepath.Join(v.Root, vaultDir)
}

// TypesDir returns the types directory path.
func (v *Vault) TypesDir() string {
	return filepath.Join(v.Dir(), "types")
}

// SharedPropertiesPath returns the path to the shared properties file.
func (v *Vault) SharedPropertiesPath() string {
	return filepath.Join(v.Dir(), "properties.yaml")
}

// DBPath returns the SQLite database path.
func (v *Vault) DBPath() string {
	return filepath.Join(v.Dir(), "index.db")
}

// TemplatesDir returns the templates directory path.
func (v *Vault) TemplatesDir() string {
	return filepath.Join(v.Root, "templates")
}

// TypeTemplatesDir returns the directory path for a specific type's templates.
func (v *Vault) TypeTemplatesDir(typeName string) string {
	return filepath.Join(v.TemplatesDir(), typeName)
}

// TemplatePath returns the file path for a specific template.
func (v *Vault) TemplatePath(typeName, templateName string) string {
	return filepath.Join(v.TypeTemplatesDir(typeName), templateName+".md")
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
	v.index = NewSQLiteObjectIndex(db)
	v.repo = NewLocalObjectRepository(v.Root)
	v.Events = NewEventDispatcher()
	v.Objects = NewObjectService(v.repo, v.index, v.Events)
	v.Queries = NewQueryService(v.repo, v.index)
	v.projector = NewProjector(v.repo, v.index, func(slug string) (*Object, error) {
		return v.Objects.Create(TagTypeName, slug, "")
	})

	if err := v.index.EnsureSchema(); err != nil {
		v.closeInternal()
		return fmt.Errorf("ensure schema: %w", err)
	}

	sync, err := v.index.NeedsSync()
	if err != nil {
		v.closeInternal()
		return fmt.Errorf("check index: %w", err)
	}
	if sync {
		if _, err := v.SyncIndex(); err != nil {
			v.closeInternal()
			return fmt.Errorf("auto sync index: %w", err)
		}
	}

	return nil
}

// closeInternal releases all resources without checking state.
func (v *Vault) closeInternal() {
	if v.db != nil {
		v.db.Close()
	}
	v.db = nil
	v.index = nil
	v.repo = nil
	v.projector = nil
	v.Objects = nil
	v.Queries = nil
	v.Events = nil
}

// Close closes the SQLite database connection.
func (v *Vault) Close() error {
	if v.db == nil {
		return nil
	}
	err := v.db.Close()
	v.db = nil
	v.index = nil
	v.repo = nil
	v.projector = nil
	v.Objects = nil
	v.Queries = nil
	v.Events = nil
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

	idx := NewSQLiteObjectIndex(db)
	if err := idx.EnsureSchema(); err != nil {
		return fmt.Errorf("initialize schema: %w", err)
	}

	return nil
}

// WriteStarterTypes writes the named starter type schemas to the vault's types directory.
// Only names matching an available starter type are written; unknown names are ignored.
func (v *Vault) WriteStarterTypes(names []string) error {
	if len(names) == 0 {
		return nil
	}
	starters := StarterTypes()
	byName := make(map[string]StarterType, len(starters))
	for _, st := range starters {
		byName[st.Name] = st
	}
	for _, name := range names {
		st, ok := byName[name]
		if !ok {
			continue
		}
		path := filepath.Join(v.TypesDir(), name+".yaml")
		if err := os.WriteFile(path, st.YAML, 0644); err != nil {
			return fmt.Errorf("write starter type %s: %w", name, err)
		}
	}
	return nil
}

