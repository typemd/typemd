package core

import (
	"fmt"
	"log/slog"
)

// Projector consumes domain events and writes them to the search index.
// It never reads or writes files — it only knows about ObjectIndex.
type Projector struct {
	index ObjectIndex
}

// NewProjector creates a Projector that writes to the given index.
func NewProjector(index ObjectIndex) *Projector {
	return &Projector{index: index}
}

// Apply processes a batch of domain events, writing each to the index.
// Calls Rebuild() after processing all events to refresh the FTS index.
func (p *Projector) Apply(events []DomainEvent) error {
	for _, event := range events {
		if err := p.applyEvent(event); err != nil {
			return err
		}
	}
	return p.index.Rebuild()
}

func (p *Projector) applyEvent(event DomainEvent) error {
	switch e := event.(type) {
	case ObjectUpserted:
		return p.index.Upsert(e.ID, e.Type, e.Filename, e.PropsJSON, e.Body)

	case ObjectDeleted:
		if err := p.index.Remove(e.ID); err != nil {
			return fmt.Errorf("remove object %s: %w", e.ID, err)
		}
		if err := p.index.DeleteWikiLinks(e.ID); err != nil {
			return fmt.Errorf("delete wikilinks for %s: %w", e.ID, err)
		}
		return p.index.DeleteRelationsByObject(e.ID)

	case RelationsCleared:
		if e.NonTagOnly {
			return p.index.DeleteNonTagRelations()
		}
		if e.TagsOnly {
			return p.index.DeleteRelationsByName(TagsProperty)
		}
		return p.index.DeleteRelationsByObject(e.ObjectID)

	case RelationIndexed:
		return p.index.InsertRelation(e.Name, e.FromID, e.ToID)

	case WikiLinksSynced:
		return p.index.SyncWikiLinks(e.ObjectID, e.Links)

	default:
		slog.Debug("projector: ignoring unknown event", "type", fmt.Sprintf("%T", event))
		return nil
	}
}
