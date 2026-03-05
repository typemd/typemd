package tui

import (
	"fmt"
	"testing"

	"github.com/MilesChou/typemd/core"
	tea "github.com/charmbracelet/bubbletea"
)

func setupTestModel(t *testing.T) model {
	t.Helper()
	objects := []*core.Object{
		{ID: "book/test-a", Type: "book", Filename: "test-a", Properties: map[string]any{"title": "A"}},
		{ID: "book/test-b", Type: "book", Filename: "test-b", Properties: map[string]any{"title": "B"}},
		{ID: "note/hello", Type: "note", Filename: "hello", Properties: map[string]any{"title": "Hello"}},
	}
	groups := buildGroups(objects)
	return model{
		focus:       focusLeft,
		groups:      groups,
		cursor:      0,
		selected:    groups[0].Objects[0],
		searchInput: initSearchInput(),
		width:       80,
		height:      24,
	}
}

func TestModel_CursorDown(t *testing.T) {
	m := setupTestModel(t)
	// cursor starts at 0 (first group header)
	msg := tea.KeyMsg{Type: tea.KeyDown}
	newM, _ := m.Update(msg)
	updated := newM.(model)
	if updated.cursor != 1 {
		t.Errorf("cursor = %d, want 1", updated.cursor)
	}
}

func TestModel_TabSwitchFocus(t *testing.T) {
	m := setupTestModel(t)
	msg := tea.KeyMsg{Type: tea.KeyTab}
	newM, _ := m.Update(msg)
	updated := newM.(model)
	if updated.focus != focusRight {
		t.Errorf("focus = %d, want focusRight", updated.focus)
	}
}

func TestModel_WindowSize(t *testing.T) {
	m := setupTestModel(t)
	msg := tea.WindowSizeMsg{Width: 120, Height: 40}
	newM, _ := m.Update(msg)
	updated := newM.(model)
	if updated.width != 120 || updated.height != 40 {
		t.Errorf("size = %dx%d, want 120x40", updated.width, updated.height)
	}
}

func TestBuildGroups_SortedByType(t *testing.T) {
	objects := []*core.Object{
		{ID: "note/a", Type: "note", Filename: "a"},
		{ID: "book/b", Type: "book", Filename: "b"},
	}
	groups := buildGroups(objects)
	if len(groups) != 2 {
		t.Fatalf("len(groups) = %d, want 2", len(groups))
	}
	if groups[0].Name != "book" {
		t.Errorf("groups[0].Name = %q, want %q", groups[0].Name, "book")
	}
	if groups[1].Name != "note" {
		t.Errorf("groups[1].Name = %q, want %q", groups[1].Name, "note")
	}
}

func TestVisibleRows_Collapsed(t *testing.T) {
	groups := []typeGroup{
		{Name: "book", Objects: []*core.Object{{ID: "book/a"}}, Expanded: false},
	}
	rows := visibleRows(groups)
	if len(rows) != 1 {
		t.Errorf("len(rows) = %d, want 1 (header only)", len(rows))
	}
}

func TestClampCursor(t *testing.T) {
	if c := clampCursor(-1, 5); c != 0 {
		t.Errorf("clampCursor(-1, 5) = %d, want 0", c)
	}
	if c := clampCursor(10, 5); c != 4 {
		t.Errorf("clampCursor(10, 5) = %d, want 4", c)
	}
	if c := clampCursor(3, 5); c != 3 {
		t.Errorf("clampCursor(3, 5) = %d, want 3", c)
	}
}

func TestBuildGroups_DefaultCollapse(t *testing.T) {
	var objects []*core.Object
	for i := 0; i < 25; i++ {
		objects = append(objects, &core.Object{
			ID: fmt.Sprintf("journal/%03d", i), Type: "journal", Filename: fmt.Sprintf("%03d", i),
		})
	}
	objects = append(objects, &core.Object{ID: "book/x", Type: "book", Filename: "x"})
	groups := buildGroups(objects)
	for _, g := range groups {
		if g.Name == "journal" && g.Expanded {
			t.Error("expected journal (25 objects) to be collapsed by default")
		}
		if g.Name == "book" && !g.Expanded {
			t.Error("expected book (1 object) to be expanded by default")
		}
	}
}
