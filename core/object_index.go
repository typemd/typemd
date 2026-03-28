package core

// ObjectResult is a projection returned by index queries.
// Use ObjectRepository.Get(id) when the source-of-truth entity is needed.
type ObjectResult struct {
	ID         string
	Type       string
	Filename   string
	Properties map[string]any
	Body       string
}

// WikiLinkEntry represents a resolved wiki-link ready for index storage.
type WikiLinkEntry struct {
	ToID        string // Resolved target ID (empty if broken)
	Target      string // Original target text
	DisplayText string
}

// SortRule defines a single sort criterion for query results.
type SortRule struct {
	Property  string `yaml:"property"`
	Direction string `yaml:"direction"` // "asc" or "desc"
}

// ObjectIndex provides search, query, and discovery operations backed by
// an indexed read model. Query methods return lightweight ObjectResult
// projections, not full domain entities.
//
// It also exposes write methods used by the Projector to maintain the index.
type ObjectIndex interface {
	// Query operations — the read side of CQRS
	Query(filter []FilterRule, sort ...SortRule) ([]*ObjectResult, error)
	Search(keyword string) ([]*ObjectResult, error)
	FindRelations(objectID string) ([]Relation, error)
	FindBacklinks(objectID string) ([]StoredWikiLink, error)
	ListWikiLinks(objectID string) ([]StoredWikiLink, error)

	// Object index maintenance — used by Projector and command-side dual writes
	Upsert(id, typeName, filename, propsJSON, body string) error
	Remove(id string) error
	ListIDs() ([]string, error)
	// Relation index maintenance
	InsertRelation(name, fromID, toID string) error
	DeleteRelation(name, fromID, toID string) error
	DeleteRelationsByName(name string) error
	DeleteNonTagRelations() error
	DeleteRelationsByObject(objectID string) error
	FindOrphanedRelations() ([]OrphanedRelation, error)
	CleanOrphanedRelations() ([]OrphanedRelation, error)

	// WikiLink index maintenance
	SyncWikiLinks(fromID string, links []WikiLinkEntry) error
	DeleteWikiLinks(fromID string) error

	// Index lifecycle
	Rebuild() error
	EnsureSchema() error
}
