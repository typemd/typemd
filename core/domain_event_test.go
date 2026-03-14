package core

import "testing"

func TestEventDispatcher_SubscribeAndDispatch(t *testing.T) {
	d := NewEventDispatcher()

	var received []DomainEvent
	d.Subscribe(func(e DomainEvent) {
		received = append(received, e)
	})

	events := []DomainEvent{
		ObjectCreated{Object: &Object{ID: "book/test"}},
		PropertyChanged{ObjectID: "book/test", Key: "title", Old: nil, New: "Go"},
	}
	d.Dispatch(events)

	if len(received) != 2 {
		t.Fatalf("received %d events, want 2", len(received))
	}
	if _, ok := received[0].(ObjectCreated); !ok {
		t.Errorf("first event should be ObjectCreated, got %T", received[0])
	}
	if _, ok := received[1].(PropertyChanged); !ok {
		t.Errorf("second event should be PropertyChanged, got %T", received[1])
	}
}

func TestEventDispatcher_MultipleHandlers(t *testing.T) {
	d := NewEventDispatcher()

	count1, count2 := 0, 0
	d.Subscribe(func(e DomainEvent) { count1++ })
	d.Subscribe(func(e DomainEvent) { count2++ })

	d.Dispatch([]DomainEvent{ObjectSaved{Object: &Object{ID: "note/test"}}})

	if count1 != 1 || count2 != 1 {
		t.Errorf("both handlers should fire: count1=%d, count2=%d", count1, count2)
	}
}

func TestEventDispatcher_EmptyDispatch(t *testing.T) {
	d := NewEventDispatcher()

	called := false
	d.Subscribe(func(e DomainEvent) { called = true })

	d.Dispatch(nil)
	if called {
		t.Error("handler should not be called for nil events")
	}

	d.Dispatch([]DomainEvent{})
	if called {
		t.Error("handler should not be called for empty events")
	}
}

func TestObjectSetProperty_ReturnsDomainEvent(t *testing.T) {
	obj := &Object{
		ID:         "book/test-01abc",
		Type:       "book",
		Properties: map[string]any{"title": "Old"},
	}
	schema := &TypeSchema{
		Properties: []Property{{Name: "title", Type: "string"}},
	}

	event, err := obj.SetProperty("title", "New", schema)
	if err != nil {
		t.Fatalf("SetProperty: %v", err)
	}

	pc, ok := event.(PropertyChanged)
	if !ok {
		t.Fatalf("expected PropertyChanged, got %T", event)
	}
	if pc.Key != "title" {
		t.Errorf("Key = %q, want %q", pc.Key, "title")
	}
	if pc.Old != "Old" {
		t.Errorf("Old = %v, want %q", pc.Old, "Old")
	}
	if pc.New != "New" {
		t.Errorf("New = %v, want %q", pc.New, "New")
	}
}

func TestObjectLinkTo_ReturnsDomainEvent(t *testing.T) {
	obj := &Object{
		ID:         "book/test-01abc",
		Type:       "book",
		Properties: map[string]any{},
	}
	prop := &Property{Name: "author", Type: "relation", Target: "person"}

	event, err := obj.LinkTo("author", "person/bob-01abc", prop)
	if err != nil {
		t.Fatalf("LinkTo: %v", err)
	}

	ol, ok := event.(ObjectLinked)
	if !ok {
		t.Fatalf("expected ObjectLinked, got %T", event)
	}
	if ol.FromID != "book/test-01abc" {
		t.Errorf("FromID = %q, want %q", ol.FromID, "book/test-01abc")
	}
	if ol.RelName != "author" {
		t.Errorf("RelName = %q, want %q", ol.RelName, "author")
	}
}

func TestObjectUnlink_ReturnsDomainEvent(t *testing.T) {
	obj := &Object{
		ID:         "book/test-01abc",
		Type:       "book",
		Properties: map[string]any{"author": "person/bob-01abc"},
	}
	prop := &Property{Name: "author", Type: "relation", Target: "person"}

	event := obj.Unlink("author", "person/bob-01abc", prop)

	ou, ok := event.(ObjectUnlinked)
	if !ok {
		t.Fatalf("expected ObjectUnlinked, got %T", event)
	}
	if ou.FromID != "book/test-01abc" {
		t.Errorf("FromID = %q, want %q", ou.FromID, "book/test-01abc")
	}
}
