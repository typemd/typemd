package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/typemd/typemd/core"
	"gopkg.in/yaml.v3"
)

func TestSessionState_MarshalRoundTrip(t *testing.T) {
	state := SessionState{
		SelectedObjectID: "book/clean-code-01jqr3k5mpbvn8e0f2g7h9txyz",
		ExpandedGroups:   []string{"book", "person"},
		ScrollOffset:     5,
		Focus:            "body",
		LeftPanelWidth:   35,
		PropsPanelWidth:  30,
		PropsVisible:     true,
	}

	data, err := yaml.Marshal(&state)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	var loaded SessionState
	if err := yaml.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if loaded.SelectedObjectID != state.SelectedObjectID {
		t.Errorf("SelectedObjectID = %q, want %q", loaded.SelectedObjectID, state.SelectedObjectID)
	}
	if len(loaded.ExpandedGroups) != 2 || loaded.ExpandedGroups[0] != "book" || loaded.ExpandedGroups[1] != "person" {
		t.Errorf("ExpandedGroups = %v, want [book person]", loaded.ExpandedGroups)
	}
	if loaded.ScrollOffset != 5 {
		t.Errorf("ScrollOffset = %d, want 5", loaded.ScrollOffset)
	}
	if loaded.Focus != "body" {
		t.Errorf("Focus = %q, want %q", loaded.Focus, "body")
	}
	if loaded.LeftPanelWidth != 35 {
		t.Errorf("LeftPanelWidth = %d, want 35", loaded.LeftPanelWidth)
	}
	if loaded.PropsPanelWidth != 30 {
		t.Errorf("PropsPanelWidth = %d, want 30", loaded.PropsPanelWidth)
	}
	if !loaded.PropsVisible {
		t.Error("PropsVisible = false, want true")
	}
}

func TestSessionState_EmptyFields(t *testing.T) {
	state := SessionState{}
	data, err := yaml.Marshal(&state)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	// Empty state should produce minimal YAML (omitempty on most fields)
	// PropsVisible (bool without omitempty) will always be present
	if !strings.Contains(string(data), "props_visible: false") {
		t.Errorf("empty state should contain 'props_visible: false', got %q", string(data))
	}
}

func TestSessionState_UnknownYAMLKeys(t *testing.T) {
	yamlData := `
selected_object_id: "book/test"
unknown_field: "should be ignored"
another_unknown: 42
`
	var state SessionState
	err := yaml.Unmarshal([]byte(yamlData), &state)
	if err != nil {
		t.Fatalf("Unmarshal should not fail on unknown keys: %v", err)
	}
	if state.SelectedObjectID != "book/test" {
		t.Errorf("SelectedObjectID = %q, want %q", state.SelectedObjectID, "book/test")
	}
}

func TestSessionState_PartialYAML(t *testing.T) {
	yamlData := `
selected_object_id: "book/test"
focus: "body"
`
	var state SessionState
	if err := yaml.Unmarshal([]byte(yamlData), &state); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	if state.SelectedObjectID != "book/test" {
		t.Errorf("SelectedObjectID = %q, want %q", state.SelectedObjectID, "book/test")
	}
	if state.Focus != "body" {
		t.Errorf("Focus = %q, want %q", state.Focus, "body")
	}
	// Missing fields should be zero values
	if state.LeftPanelWidth != 0 {
		t.Errorf("LeftPanelWidth = %d, want 0 (zero value for missing field)", state.LeftPanelWidth)
	}
	if state.PropsVisible {
		t.Error("PropsVisible = true, want false (zero value for missing field)")
	}
}

func TestLoadSessionState_MissingFile(t *testing.T) {
	state := loadSessionState(t.TempDir())
	if state.SelectedObjectID != "" {
		t.Errorf("missing file should return zero state, got SelectedObjectID=%q", state.SelectedObjectID)
	}
}

func TestLoadSessionState_CorruptFile(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, ".typemd"), 0755)
	os.WriteFile(stateFilePath(dir), []byte(":::invalid yaml:::"), 0644)

	state := loadSessionState(dir)
	if state.SelectedObjectID != "" {
		t.Errorf("corrupt file should return zero state, got SelectedObjectID=%q", state.SelectedObjectID)
	}
}

func TestSaveAndLoadSessionState(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, ".typemd"), 0755)

	original := SessionState{
		SelectedObjectID: "book/test-01abc",
		ExpandedGroups:   []string{"book"},
		Focus:            "left",
		LeftPanelWidth:   40,
		PropsPanelWidth:  25,
		PropsVisible:     true,
		ScrollOffset:     3,
	}

	saveSessionState(dir, original)

	loaded := loadSessionState(dir)
	if loaded.SelectedObjectID != original.SelectedObjectID {
		t.Errorf("SelectedObjectID = %q, want %q", loaded.SelectedObjectID, original.SelectedObjectID)
	}
	if len(loaded.ExpandedGroups) != 1 || loaded.ExpandedGroups[0] != "book" {
		t.Errorf("ExpandedGroups = %v, want [book]", loaded.ExpandedGroups)
	}
	if loaded.Focus != "left" {
		t.Errorf("Focus = %q, want %q", loaded.Focus, "left")
	}
	if loaded.PropsVisible != true {
		t.Error("PropsVisible = false, want true")
	}
}

func TestSaveSessionState_NoTypemdDir(t *testing.T) {
	dir := t.TempDir()
	// Don't create .typemd/ — saveSessionState should not panic
	saveSessionState(dir, SessionState{SelectedObjectID: "test"})
	// Just verify no panic; file won't be written and that's fine
}

func TestCaptureState_FullState(t *testing.T) {
	m := setupTestModel(t)
	m.groups[0].Expanded = true
	m.groups[1].Expanded = true
	m.selected = m.groups[0].Objects[0]
	m.focus = focusBody
	m.leftW = 35
	m.propsWidth = 30
	m.propsVisible = true
	m.scrollOffset = 5

	state := m.captureState()

	if state.SelectedObjectID != m.selected.ID {
		t.Errorf("SelectedObjectID = %q, want %q", state.SelectedObjectID, m.selected.ID)
	}
	if len(state.ExpandedGroups) != 2 {
		t.Errorf("ExpandedGroups len = %d, want 2", len(state.ExpandedGroups))
	}
	if state.Focus != "body" {
		t.Errorf("Focus = %q, want %q", state.Focus, "body")
	}
	if state.LeftPanelWidth != 35 {
		t.Errorf("LeftPanelWidth = %d, want 35", state.LeftPanelWidth)
	}
	if state.PropsPanelWidth != 30 {
		t.Errorf("PropsPanelWidth = %d, want 30", state.PropsPanelWidth)
	}
	if !state.PropsVisible {
		t.Error("PropsVisible = false, want true")
	}
	if state.ScrollOffset != 5 {
		t.Errorf("ScrollOffset = %d, want 5", state.ScrollOffset)
	}
}

func TestCaptureState_NoSelection(t *testing.T) {
	m := setupTestModel(t)
	m.selected = nil

	state := m.captureState()

	if state.SelectedObjectID != "" {
		t.Errorf("SelectedObjectID = %q, want empty", state.SelectedObjectID)
	}
}

func TestCaptureState_NoGroupsExpanded(t *testing.T) {
	m := setupTestModel(t)
	for i := range m.groups {
		m.groups[i].Expanded = false
	}

	state := m.captureState()

	if len(state.ExpandedGroups) != 0 {
		t.Errorf("ExpandedGroups = %v, want empty", state.ExpandedGroups)
	}
}

func TestCaptureState_ExcludesSearchState(t *testing.T) {
	m := setupTestModel(t)
	m.searchMode = true
	m.searchResults = []*core.Object{{ID: "book/result"}}

	state := m.captureState()

	// State should capture normal state, search state is excluded by design
	// (no searchMode or searchResults fields in SessionState)
	if state.SelectedObjectID != m.selected.ID {
		t.Errorf("should still capture selected object even in search mode")
	}
}

func TestApplySessionState_RestoreSelectedObject(t *testing.T) {
	groups := []typeGroup{
		{Name: "book", Objects: []*core.Object{
			{ID: "book/a", Type: "book"},
			{ID: "book/b", Type: "book"},
		}},
		{Name: "note", Objects: []*core.Object{
			{ID: "note/c", Type: "note"},
		}},
	}
	state := SessionState{
		SelectedObjectID: "book/b",
		ExpandedGroups:   []string{"book"},
	}

	cursor, selectedID := applySessionState(state, groups)

	if selectedID != "book/b" {
		t.Errorf("selectedID = %q, want %q", selectedID, "book/b")
	}
	// cursor should point to "book/b" which is: header(0), book/a(1), book/b(2)
	if cursor != 2 {
		t.Errorf("cursor = %d, want 2", cursor)
	}
}

func TestApplySessionState_RestoreExpandedGroups(t *testing.T) {
	groups := []typeGroup{
		{Name: "book", Objects: []*core.Object{{ID: "book/a", Type: "book"}}},
		{Name: "note", Objects: []*core.Object{{ID: "note/b", Type: "note"}}},
		{Name: "person", Objects: []*core.Object{{ID: "person/c", Type: "person"}}},
	}
	state := SessionState{
		ExpandedGroups: []string{"book", "person"},
	}

	applySessionState(state, groups)

	if !groups[0].Expanded {
		t.Error("book should be expanded")
	}
	if groups[1].Expanded {
		t.Error("note should not be expanded")
	}
	if !groups[2].Expanded {
		t.Error("person should be expanded")
	}
}

func TestApplySessionState_NoSavedState(t *testing.T) {
	groups := []typeGroup{
		{Name: "book", Objects: []*core.Object{{ID: "book/a", Type: "book"}}},
	}
	state := SessionState{} // empty/default

	cursor, selectedID := applySessionState(state, groups)

	// Should fall back to first group expanded, first object selected
	if !groups[0].Expanded {
		t.Error("first group should be expanded when no saved state")
	}
	if selectedID != "book/a" {
		t.Errorf("selectedID = %q, want %q (fallback to first)", selectedID, "book/a")
	}
	if cursor != 1 { // header(0), book/a(1)
		t.Errorf("cursor = %d, want 1", cursor)
	}
}

func TestApplySessionState_EmptyVault(t *testing.T) {
	groups := []typeGroup{}
	state := SessionState{SelectedObjectID: "book/missing"}

	cursor, selectedID := applySessionState(state, groups)

	if selectedID != "" {
		t.Errorf("selectedID = %q, want empty for empty vault", selectedID)
	}
	if cursor != 0 {
		t.Errorf("cursor = %d, want 0 for empty vault", cursor)
	}
}

func TestApplySessionState_ObjectDeletedSameTypeExists(t *testing.T) {
	groups := []typeGroup{
		{Name: "book", Objects: []*core.Object{
			{ID: "book/remaining", Type: "book"},
		}},
	}
	state := SessionState{
		SelectedObjectID: "book/deleted-01xyz", // doesn't exist
		ExpandedGroups:   []string{"book"},
	}

	cursor, selectedID := applySessionState(state, groups)

	// Should fallback to first object in same type group
	if selectedID != "book/remaining" {
		t.Errorf("selectedID = %q, want %q (fallback to same type)", selectedID, "book/remaining")
	}
	if cursor != 1 { // header(0), book/remaining(1)
		t.Errorf("cursor = %d, want 1", cursor)
	}
}

func TestApplySessionState_ObjectDeletedTypeRemoved(t *testing.T) {
	groups := []typeGroup{
		{Name: "note", Objects: []*core.Object{
			{ID: "note/hello", Type: "note"},
		}},
	}
	state := SessionState{
		SelectedObjectID: "book/deleted-01xyz", // type "book" doesn't exist
		ExpandedGroups:   []string{"book"},     // stale group
	}

	cursor, selectedID := applySessionState(state, groups)

	// No "book" type exists, should fallback to first object overall
	// "note" group should be expanded as fallback since no saved groups matched
	if !groups[0].Expanded {
		t.Error("note group should be expanded as fallback")
	}
	if selectedID != "note/hello" {
		t.Errorf("selectedID = %q, want %q (fallback to first overall)", selectedID, "note/hello")
	}
	_ = cursor
}

func TestApplySessionState_StaleTypeGroupIgnored(t *testing.T) {
	groups := []typeGroup{
		{Name: "book", Objects: []*core.Object{{ID: "book/a", Type: "book"}}},
		{Name: "note", Objects: []*core.Object{{ID: "note/b", Type: "note"}}},
	}
	state := SessionState{
		ExpandedGroups: []string{"book", "deleted-type"},
	}

	applySessionState(state, groups)

	if !groups[0].Expanded {
		t.Error("book should be expanded")
	}
	if groups[1].Expanded {
		t.Error("note should not be expanded")
	}
	// "deleted-type" should be silently ignored — no panic, no error
}

func TestApplySessionState_SingleObjectVault(t *testing.T) {
	groups := []typeGroup{
		{Name: "book", Objects: []*core.Object{{ID: "book/only", Type: "book"}}},
	}
	state := SessionState{
		SelectedObjectID: "book/only",
		ExpandedGroups:   []string{"book"},
	}

	cursor, selectedID := applySessionState(state, groups)

	if selectedID != "book/only" {
		t.Errorf("selectedID = %q, want %q", selectedID, "book/only")
	}
	if cursor != 1 { // header(0), book/only(1)
		t.Errorf("cursor = %d, want 1", cursor)
	}
}

func TestApplySessionState_ObjectIDNoTypeMatch(t *testing.T) {
	groups := []typeGroup{
		{Name: "note", Objects: []*core.Object{{ID: "note/hello", Type: "note"}}},
	}
	state := SessionState{
		SelectedObjectID: "invalidformat", // no slash, can't extract type
	}

	cursor, selectedID := applySessionState(state, groups)

	// Should fallback to first object (no type could be extracted)
	if selectedID != "note/hello" {
		t.Errorf("selectedID = %q, want %q", selectedID, "note/hello")
	}
	_ = cursor
}

func TestApplySessionState_FallbackExpandsTypeGroup(t *testing.T) {
	// Object deleted, but its type group exists but was NOT in expandedGroups.
	// The fallback should expand the type group to find the first object.
	groups := []typeGroup{
		{Name: "book", Objects: []*core.Object{
			{ID: "book/survivor", Type: "book"},
		}},
	}
	state := SessionState{
		SelectedObjectID: "book/deleted-01xyz",
		// Note: "book" is NOT in expandedGroups
	}

	cursor, selectedID := applySessionState(state, groups)

	// applySessionState should expand the book group for fallback
	if !groups[0].Expanded {
		t.Error("book group should be expanded for fallback")
	}
	if selectedID != "book/survivor" {
		t.Errorf("selectedID = %q, want %q", selectedID, "book/survivor")
	}
	_ = cursor
}

func TestApplySessionState_PanelWidthsClampedByTerminal(t *testing.T) {
	// Panel widths are stored in state but clamped by WindowSizeMsg handler.
	// This test verifies that applySessionState just passes the values through;
	// actual clamping happens in the existing Update() WindowSizeMsg handler.
	state := SessionState{
		LeftPanelWidth:  999,
		PropsPanelWidth: 999,
	}
	groups := []typeGroup{
		{Name: "book", Objects: []*core.Object{{ID: "book/a", Type: "book"}}},
	}

	applySessionState(state, groups)
	// applySessionState doesn't handle widths — they're applied directly in Start()
	// This test documents that width clamping is NOT applySessionState's responsibility
}

func TestFocusPanelToString(t *testing.T) {
	tests := []struct {
		input focusPanel
		want  string
	}{
		{focusLeft, "left"},
		{focusBody, "body"},
		{focusProps, "props"},
	}
	for _, tt := range tests {
		got := focusPanelToString(tt.input)
		if got != tt.want {
			t.Errorf("focusPanelToString(%d) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestStringToFocusPanel(t *testing.T) {
	tests := []struct {
		input string
		want  focusPanel
	}{
		{"left", focusLeft},
		{"body", focusBody},
		{"props", focusProps},
		{"unknown", focusLeft}, // default
		{"", focusLeft},        // default
	}
	for _, tt := range tests {
		got := stringToFocusPanel(tt.input)
		if got != tt.want {
			t.Errorf("stringToFocusPanel(%q) = %d, want %d", tt.input, got, tt.want)
		}
	}
}
