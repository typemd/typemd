package core

import "time"

// ObjectRepository provides entity-level persistence for all file-backed
// entities: objects, type schemas, templates, and shared properties.
//
// All methods return domain entities, not raw bytes. Serialization and
// path conventions are encapsulated within each implementation.
//
// This is the write-side repository in CQRS — the source of truth.
// Use ObjectIndex for search and query operations.
type ObjectRepository interface {
	// Object entity operations
	Get(id string) (*Object, error)
	Save(obj *Object, keyOrder []string) error
	Create(obj *Object, keyOrder []string) error // fails if file already exists
	Walk() ([]*Object, error)
	GlobIDs(typeName, prefix string) ([]string, error)
	ModTime(id string) (time.Time, error)
	EnsureDir(typeName string) error

	// Type schema operations
	GetSchema(typeName string) (*TypeSchema, error)
	WriteSchema(typeName string, data []byte) error
	DeleteSchema(typeName string) error
	ListSchemas() ([]string, error)

	// Template operations
	GetTemplate(typeName, name string) (*Template, error)
	ListTemplates(typeName string) ([]string, error)
	SaveTemplate(typeName, name string, tmpl *Template) error
	DeleteTemplate(typeName, name string) error

	// Shared property operations
	GetSharedProperties() ([]Property, error)
}
