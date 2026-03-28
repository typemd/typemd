package tui

import (
	"strings"
	"testing"

	"github.com/typemd/typemd/core"
	"github.com/typemd/typemd/tui/widget"
)

func TestSyncUnresolvedToToast_MixedReasons(t *testing.T) {
	// Two different reasons should produce two separate group lines.
	unresolved := []core.UnresolvedRelation{
		{ObjectID: "book/foo-01abc", Property: "author", Value: "John", Reason: "not_found"},
		{ObjectID: "note/bar-02def", Property: "related", Value: "Baz", Reason: "ambiguous"},
	}

	items := unresolvedToToastItems(unresolved)

	toast := widget.NewToastModel()
	cmd := toast.Show(widget.ToastWarning, items)

	if !toast.Active() {
		t.Fatal("expected toast to be active after Show with unresolved items")
	}
	if cmd == nil {
		t.Fatal("expected non-nil cmd (tick) from Show")
	}

	view := toast.View()

	if !strings.Contains(view, "⚠") {
		t.Errorf("expected warning prefix ⚠ in view, got: %s", view)
	}
	if !strings.Contains(view, "1 not found") {
		t.Errorf("expected '1 not found' in view, got: %s", view)
	}
	if !strings.Contains(view, "1 ambiguous") {
		t.Errorf("expected '1 ambiguous' in view, got: %s", view)
	}
}

func TestSyncUnresolvedToToast_OnlyNotFound(t *testing.T) {
	unresolved := []core.UnresolvedRelation{
		{ObjectID: "book/foo-01abc", Property: "author", Value: "John", Reason: "not_found"},
		{ObjectID: "book/bar-02def", Property: "author", Value: "Jane", Reason: "not_found"},
	}

	items := unresolvedToToastItems(unresolved)

	toast := widget.NewToastModel()
	toast.Show(widget.ToastWarning, items)

	view := toast.View()
	if !strings.Contains(view, "2 not found") {
		t.Errorf("expected '2 not found' in view, got: %s", view)
	}
	if strings.Contains(view, "ambiguous") {
		t.Errorf("should not contain 'ambiguous' when all reasons are not_found, got: %s", view)
	}
}

func TestSyncUnresolvedToToast_OnlyAmbiguous(t *testing.T) {
	unresolved := []core.UnresolvedRelation{
		{ObjectID: "book/foo-01abc", Property: "author", Value: "John", Reason: "ambiguous"},
		{ObjectID: "note/bar-02def", Property: "related", Value: "Baz", Reason: "ambiguous"},
		{ObjectID: "idea/qux-03ghi", Property: "source", Value: "Quux", Reason: "ambiguous"},
	}

	items := unresolvedToToastItems(unresolved)

	toast := widget.NewToastModel()
	toast.Show(widget.ToastWarning, items)

	view := toast.View()
	if !strings.Contains(view, "3 ambiguous") {
		t.Errorf("expected '3 ambiguous' in view, got: %s", view)
	}
}

func TestSyncNoUnresolved_NoToast(t *testing.T) {
	// When ReconcileResult has no Unresolved items, no toast should be shown.
	result := &core.ReconcileResult{
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

	items := unresolvedToToastItems(unresolved)

	toast := widget.NewToastModel()
	cmd := toast.Show(widget.ToastWarning, items)

	if !toast.Active() {
		t.Fatal("expected toast to be active")
	}
	if cmd == nil {
		t.Fatal("expected non-nil cmd")
	}

	view := toast.View()
	if !strings.Contains(view, "1 not found") {
		t.Errorf("expected '1 not found' in view, got: %s", view)
	}
}

func TestSyncUnresolvedToToast_UnknownReasonFallback(t *testing.T) {
	// Unknown reason should fall back to "unresolved" group.
	unresolved := []core.UnresolvedRelation{
		{ObjectID: "book/foo-01abc", Property: "author", Value: "John", Reason: "something_else"},
	}

	items := unresolvedToToastItems(unresolved)

	toast := widget.NewToastModel()
	toast.Show(widget.ToastWarning, items)

	view := toast.View()
	if !strings.Contains(view, "1 unresolved") {
		t.Errorf("expected '1 unresolved' in view for unknown reason, got: %s", view)
	}
}

func TestReasonToGroup(t *testing.T) {
	tests := []struct {
		reason string
		want   string
	}{
		{"not_found", "not found"},
		{"ambiguous", "ambiguous"},
		{"", "unresolved"},
		{"unknown", "unresolved"},
	}
	for _, tt := range tests {
		got := reasonToGroup(tt.reason)
		if got != tt.want {
			t.Errorf("reasonToGroup(%q) = %q, want %q", tt.reason, got, tt.want)
		}
	}
}
