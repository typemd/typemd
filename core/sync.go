package core

import "fmt"

// OrphanedRelation represents a relation record that references a non-existent object.
type OrphanedRelation struct {
	Name   string
	FromID string
	ToID   string
}

// Reason constants for UnresolvedRelation.
const (
	ReasonNotFound  = "not_found"
	ReasonAmbiguous = "ambiguous"
)

// UnresolvedRelation represents a relation reference that could not be resolved during reconciliation.
type UnresolvedRelation struct {
	ObjectID string
	Property string
	Value    string
	Reason   string
	Matches  []string
}

// UnresolvedWikiLink represents a wiki-link that could not be resolved during reconciliation.
type UnresolvedWikiLink struct {
	ObjectID string
	Target   string
	Reason   string
	Matches  []string
}

// ReconcileResult holds statistics from a Reconcile operation.
type ReconcileResult struct {
	Synced              int
	Deleted             int
	Expanded            int
	Unresolved          []UnresolvedRelation
	WikiLinksExpanded   int
	UnresolvedWikiLinks []UnresolvedWikiLink
}

// syncContext holds intermediate state collected during reconciliation.
type syncContext struct {
	diskIDs     map[string]bool
	diskBodies  map[string]string
	diskTags    map[string]*Object
	diskTagRefs map[string][]string
	diskObjects map[string]*Object
	schemas     map[string]*TypeSchema
	nameIndex   map[string]map[string][]string
}

// Reconcile normalizes all vault files and returns domain events describing what changed.
func (v *Vault) Reconcile() ([]DomainEvent, *ReconcileResult, error) {
	if v.reconciler == nil {
		return nil, nil, fmt.Errorf("vault not opened")
	}
	return v.reconciler.Reconcile()
}

// ReconcileFiles incrementally reconciles specific files and returns domain events.
func (v *Vault) ReconcileFiles(paths []string) ([]DomainEvent, *ReconcileResult, error) {
	if v.reconciler == nil {
		return nil, nil, fmt.Errorf("vault not opened")
	}
	if len(paths) == 0 {
		return v.reconciler.Reconcile()
	}
	return v.reconciler.ReconcileFiles(paths, v.ObjectsDir())
}

// Project applies domain events to the search index.
func (v *Vault) Project(events []DomainEvent) error {
	if v.projector == nil {
		return fmt.Errorf("vault not opened")
	}
	return v.projector.Apply(events)
}
