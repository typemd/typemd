package core

// DomainEvent is the marker interface for all domain events.
// Entity methods return []DomainEvent to signal what happened.
// Use Case layer collects and dispatches events after successful operations.
type DomainEvent interface {
	eventName() string
}

// ObjectCreated is emitted when a new object is created.
type ObjectCreated struct{ Object *Object }

func (e ObjectCreated) eventName() string { return "object.created" }

// ObjectSaved is emitted when an existing object is saved.
type ObjectSaved struct{ Object *Object }

func (e ObjectSaved) eventName() string { return "object.saved" }

// PropertyChanged is emitted when a single property value changes.
type PropertyChanged struct {
	ObjectID string
	Key      string
	Old, New any
}

func (e PropertyChanged) eventName() string { return "object.property_changed" }

// ObjectLinked is emitted when a relation is created between two objects.
type ObjectLinked struct {
	FromID  string
	ToID    string
	RelName string
}

func (e ObjectLinked) eventName() string { return "object.linked" }

// ObjectUnlinked is emitted when a relation is removed between two objects.
type ObjectUnlinked struct {
	FromID  string
	ToID    string
	RelName string
}

func (e ObjectUnlinked) eventName() string { return "object.unlinked" }

// TagAutoCreated is emitted when a tag object is auto-created during sync.
type TagAutoCreated struct {
	Tag          *Object
	ReferencedBy string
}

func (e TagAutoCreated) eventName() string { return "tag.auto_created" }

// ObjectDeleted is emitted when an object is removed from disk (stale cleanup).
type ObjectDeleted struct{ ID string }

func (e ObjectDeleted) eventName() string { return "object.deleted" }

// ObjectUpserted is emitted when an object needs to be written to the index.
type ObjectUpserted struct {
	ID        string
	Type      string
	Filename  string
	PropsJSON string
	Body      string
}

func (e ObjectUpserted) eventName() string { return "object.upserted" }

// WikiLinksSynced is emitted when wiki-links for an object have been resolved.
type WikiLinksSynced struct {
	ObjectID string
	Links    []WikiLinkEntry
}

func (e WikiLinksSynced) eventName() string { return "wikilinks.synced" }

// RelationIndexed is emitted when a single relation needs to be indexed.
type RelationIndexed struct {
	Name   string
	FromID string
	ToID   string
}

func (e RelationIndexed) eventName() string { return "relation.indexed" }

// RelationsCleared is emitted to clear relations for a scope before rebuilding.
type RelationsCleared struct {
	ObjectID   string // if set, clear only this object's relations
	NonTagOnly bool   // if true, clear all non-tag relations (full sync)
	TagsOnly   bool   // if true, clear tag relations only
}

func (e RelationsCleared) eventName() string { return "relations.cleared" }

// TypeSaved is emitted when a type schema is created or updated.
type TypeSaved struct{ Schema *TypeSchema }

func (e TypeSaved) eventName() string { return "type.saved" }

// TypeDeleted is emitted when a type schema is deleted.
type TypeDeleted struct{ Name string }

func (e TypeDeleted) eventName() string { return "type.deleted" }

// EventDispatcher collects and dispatches domain events to subscribers.
type EventDispatcher struct {
	handlers []EventHandler
}

// EventHandler is a function that handles a domain event.
type EventHandler func(DomainEvent)

// NewEventDispatcher creates a new EventDispatcher.
func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{}
}

// Subscribe registers an event handler.
func (d *EventDispatcher) Subscribe(handler EventHandler) {
	d.handlers = append(d.handlers, handler)
}

// Dispatch sends events to all registered handlers.
func (d *EventDispatcher) Dispatch(events []DomainEvent) {
	for _, event := range events {
		for _, handler := range d.handlers {
			handler(event)
		}
	}
}
