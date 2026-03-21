package core

import "fmt"

// OrphanedRelation represents a relation record that references a non-existent object.
type OrphanedRelation struct {
	Name   string
	FromID string
	ToID   string
}

// UnresolvedRelation represents a relation reference that could not be resolved during sync.
type UnresolvedRelation struct {
	ObjectID string // the object containing the reference
	Property string // the relation property name
	Value    string // the unresolved value
	Reason   string // "not_found" or "ambiguous"
	Matches  []string // candidate IDs for ambiguous matches
}

// SyncResult holds statistics from a SyncIndex operation.
type SyncResult struct {
	Synced     int                  // number of objects upserted into the index
	Deleted    int
	Orphaned   []OrphanedRelation
	Expanded   int                  // number of relation prefixes auto-expanded to full IDs
	Unresolved []UnresolvedRelation // relation references that could not be resolved
}

// syncContext holds intermediate state collected during sync.
type syncContext struct {
	diskIDs     map[string]bool
	diskBodies  map[string]string
	diskTags    map[string]*Object
	diskTagRefs map[string][]string
	diskObjects map[string]*Object           // all objects by ID
	schemas     map[string]*TypeSchema        // cached type schemas
	nameIndex   map[string]map[string][]string // nameIndex[type][name] = []objectID
}

// SyncIndex scans the objects directory, upserts all found objects into the index,
// removes stale entries, cleans up orphaned relations, and rebuilds the FTS index.
func (v *Vault) SyncIndex() (*SyncResult, error) {
	if v.projector == nil {
		return nil, fmt.Errorf("vault not opened")
	}
	return v.projector.Sync()
}

// SyncFiles incrementally synchronizes specific files to the index.
// Falls back to full SyncIndex if paths is empty.
func (v *Vault) SyncFiles(paths []string) (*SyncResult, error) {
	if v.projector == nil {
		return nil, fmt.Errorf("vault not opened")
	}
	if len(paths) == 0 {
		return v.projector.Sync()
	}
	return v.projector.SyncFiles(paths, v.ObjectsDir())
}
