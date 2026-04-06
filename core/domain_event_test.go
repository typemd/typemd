package core

import "testing"

func TestNewDomainEventNames(t *testing.T) {
	tests := []struct {
		event DomainEvent
		want  string
	}{
		{ObjectDeleted{ID: "book/x"}, "object.deleted"},
		{ObjectUpserted{ID: "book/x"}, "object.upserted"},
		{TypeSaved{Schema: &TypeSchema{Name: "book"}}, "type.saved"},
		{TypeDeleted{Name: "book"}, "type.deleted"},
		{WikiLinksSynced{ObjectID: "book/x"}, "wikilinks.synced"},
		{RelationIndexed{Name: "author"}, "relation.indexed"},
		{RelationsCleared{NonTagOnly: true}, "relations.cleared"},
	}
	for _, tt := range tests {
		if got := tt.event.eventName(); got != tt.want {
			t.Errorf("%T.eventName() = %q, want %q", tt.event, got, tt.want)
		}
	}
}
