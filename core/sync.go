package core

import "fmt"

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

// syncContext holds intermediate state collected during sync.
type syncContext struct {
	diskIDs     map[string]bool
	diskBodies  map[string]string
	diskTags    map[string]*Object
	diskTagRefs map[string][]string
}

// SyncIndex scans the objects directory, upserts all found objects into the index,
// removes stale entries, cleans up orphaned relations, and rebuilds the FTS index.
func (v *Vault) SyncIndex() (*SyncResult, error) {
	if v.projector == nil {
		return nil, fmt.Errorf("vault not opened")
	}
	return v.projector.Sync()
}
