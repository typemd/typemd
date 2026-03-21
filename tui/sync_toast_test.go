package tui

import (
	"fmt"
	"strings"
	"testing"

	"github.com/typemd/typemd/core"
	"github.com/typemd/typemd/tui/widget"
)

func TestSyncUnresolvedToToast(t *testing.T) {
	// Simulate the conversion logic from refreshData:
	// given a SyncResult with Unresolved items, build ToastItems and call Show.
	unresolved := []core.UnresolvedRelation{
		{ObjectID: "book/foo-01abc", Property: "author", Value: "John", Reason: "not_found"},
		{ObjectID: "note/bar-02def", Property: "related", Value: "Baz", Reason: "ambiguous"},
	}

	items := make([]widget.ToastItem, len(unresolved))
	for i, u := range unresolved {
		items[i] = widget.ToastItem{
			Message: fmt.Sprintf("%s.%s: %s", u.ObjectID, u.Property, u.Value),
			Group:   "unresolved refs",
		}
	}

	toast := widget.NewToastModel()
	cmd := toast.Show(widget.ToastWarning, items)

	if !toast.Active() {
		t.Fatal("expected toast to be active after Show with unresolved items")
	}
	if cmd == nil {
		t.Fatal("expected non-nil cmd (tick) from Show")
	}

	view := toast.View()

	// Should show warning prefix
	if !strings.Contains(view, "⚠") {
		t.Errorf("expected warning prefix ⚠ in view, got: %s", view)
	}

	// Should show aggregated count (2 items with same group key)
	if !strings.Contains(view, "2") {
		t.Errorf("expected count '2' in view, got: %s", view)
	}

	// Should show group name
	if !strings.Contains(view, "unresolved refs") {
		t.Errorf("expected group name 'unresolved refs' in view, got: %s", view)
	}
}

func TestSyncNoUnresolved_NoToast(t *testing.T) {
	// When SyncResult has no Unresolved items, no toast should be shown.
	result := &core.SyncResult{
		Synced:     5,
		Unresolved: nil,
	}

	toast := widget.NewToastModel()

	// Simulate the condition from refreshData
	if result != nil && len(result.Unresolved) > 0 {
		t.Fatal("should not enter toast branch when Unresolved is empty")
	}

	if toast.Active() {
		t.Fatal("expected toast to remain inactive when no unresolved items")
	}
}

func TestSyncSingleUnresolved(t *testing.T) {
	// Single unresolved item should still produce a grouped toast with count 1.
	unresolved := []core.UnresolvedRelation{
		{ObjectID: "book/test-01abc", Property: "genre", Value: "sci-fi", Reason: "not_found"},
	}

	items := make([]widget.ToastItem, len(unresolved))
	for i, u := range unresolved {
		items[i] = widget.ToastItem{
			Message: fmt.Sprintf("%s.%s: %s", u.ObjectID, u.Property, u.Value),
			Group:   "unresolved refs",
		}
	}

	toast := widget.NewToastModel()
	cmd := toast.Show(widget.ToastWarning, items)

	if !toast.Active() {
		t.Fatal("expected toast to be active")
	}
	if cmd == nil {
		t.Fatal("expected non-nil cmd")
	}

	view := toast.View()
	if !strings.Contains(view, "1") {
		t.Errorf("expected count '1' in view, got: %s", view)
	}
	if !strings.Contains(view, "unresolved refs") {
		t.Errorf("expected group name in view, got: %s", view)
	}
}
