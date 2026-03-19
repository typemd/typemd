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

// --- View mode persistence tests ---

func TestSessionState_ViewModeFields_MarshalRoundTrip(t *testing.T) {
	state := SessionState{
		SelectedObjectID:   "book/test-01abc",
		ViewTypeName:       "book",
		ViewName:           "by-rating",
		ViewCursor:         3,
		ViewScroll:         2,
		ViewExpandedGroups: []string{"5 stars", "4 stars"},
	}

	data, err := yaml.Marshal(&state)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	var loaded SessionState
	if err := yaml.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if loaded.ViewTypeName != "book" {
		t.Errorf("ViewTypeName = %q, want %q", loaded.ViewTypeName, "book")
	}
	if loaded.ViewName != "by-rating" {
		t.Errorf("ViewName = %q, want %q", loaded.ViewName, "by-rating")
	}
	if loaded.ViewCursor != 3 {
		t.Errorf("ViewCursor = %d, want 3", loaded.ViewCursor)
	}
	if loaded.ViewScroll != 2 {
		t.Errorf("ViewScroll = %d, want 2", loaded.ViewScroll)
	}
	if len(loaded.ViewExpandedGroups) != 2 || loaded.ViewExpandedGroups[0] != "5 stars" || loaded.ViewExpandedGroups[1] != "4 stars" {
		t.Errorf("ViewExpandedGroups = %v, want [5 stars, 4 stars]", loaded.ViewExpandedGroups)
	}
}

func TestSessionState_ViewModeFields_OmittedWhenEmpty(t *testing.T) {
	state := SessionState{
		SelectedObjectID: "book/test",
		PropsVisible:     true,
	}
	data, err := yaml.Marshal(&state)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}

	yamlStr := string(data)
	if strings.Contains(yamlStr, "view_type_name") {
		t.Errorf("empty ViewTypeName should be omitted, got:\n%s", yamlStr)
	}
	if strings.Contains(yamlStr, "view_name") {
		t.Errorf("empty ViewName should be omitted, got:\n%s", yamlStr)
	}
	if strings.Contains(yamlStr, "view_cursor") {
		t.Errorf("zero ViewCursor should be omitted, got:\n%s", yamlStr)
	}
}

func TestSessionState_ViewModeFields_BackwardCompatible(t *testing.T) {
	// Old YAML without view fields should deserialize cleanly
	yamlData := `
selected_object_id: "book/test"
expanded_groups:
  - book
focus: "left"
props_visible: true
`
	var state SessionState
	if err := yaml.Unmarshal([]byte(yamlData), &state); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	if state.SelectedObjectID != "book/test" {
		t.Errorf("SelectedObjectID = %q, want %q", state.SelectedObjectID, "book/test")
	}
	if state.ViewTypeName != "" {
		t.Errorf("ViewTypeName = %q, want empty (missing field)", state.ViewTypeName)
	}
	if state.ViewName != "" {
		t.Errorf("ViewName = %q, want empty (missing field)", state.ViewName)
	}
}

func TestCaptureState_InViewMode(t *testing.T) {
	m := setupTestModel(t)
	m.rightPanel = panelView
	m.viewMode = &viewMode{
		typeName: "book",
		viewName: "by-rating",
		cursor:   5,
		scroll:   2,
		groups: []viewGroup{
			{Label: "5 stars", Expanded: true},
			{Label: "4 stars", Expanded: false},
			{Label: "3 stars", Expanded: true},
		},
	}

	state := m.captureState()

	if state.ViewTypeName != "book" {
		t.Errorf("ViewTypeName = %q, want %q", state.ViewTypeName, "book")
	}
	if state.ViewName != "by-rating" {
		t.Errorf("ViewName = %q, want %q", state.ViewName, "by-rating")
	}
	if state.ViewCursor != 5 {
		t.Errorf("ViewCursor = %d, want 5", state.ViewCursor)
	}
	if state.ViewScroll != 2 {
		t.Errorf("ViewScroll = %d, want 2", state.ViewScroll)
	}
	if len(state.ViewExpandedGroups) != 2 {
		t.Fatalf("ViewExpandedGroups len = %d, want 2", len(state.ViewExpandedGroups))
	}
	if state.ViewExpandedGroups[0] != "5 stars" || state.ViewExpandedGroups[1] != "3 stars" {
		t.Errorf("ViewExpandedGroups = %v, want [5 stars, 3 stars]", state.ViewExpandedGroups)
	}
}

func TestCaptureState_NotInViewMode(t *testing.T) {
	m := setupTestModel(t)
	m.rightPanel = panelObject
	m.viewMode = nil

	state := m.captureState()

	if state.ViewTypeName != "" {
		t.Errorf("ViewTypeName = %q, want empty (not in view mode)", state.ViewTypeName)
	}
	if state.ViewName != "" {
		t.Errorf("ViewName = %q, want empty (not in view mode)", state.ViewName)
	}
}

func TestCaptureState_ViewModeWithFlatList(t *testing.T) {
	// Flat list (no group_by) has groups with empty labels — should not appear in expanded list
	m := setupTestModel(t)
	m.rightPanel = panelView
	m.viewMode = &viewMode{
		typeName: "book",
		viewName: "default",
		cursor:   0,
		groups: []viewGroup{
			{Label: "", Expanded: true, Objects: []*core.Object{{ID: "book/a"}}},
		},
	}

	state := m.captureState()

	if len(state.ViewExpandedGroups) != 0 {
		t.Errorf("ViewExpandedGroups = %v, want empty (flat list has empty labels)", state.ViewExpandedGroups)
	}
}

func TestExpandedGroupLabels(t *testing.T) {
	vm := &viewMode{
		groups: []viewGroup{
			{Label: "fiction", Expanded: true},
			{Label: "non-fiction", Expanded: false},
			{Label: "science", Expanded: true},
			{Label: "", Expanded: true}, // empty label (flat list) should be excluded
		},
	}

	labels := vm.expandedGroupLabels()

	if len(labels) != 2 {
		t.Fatalf("expandedGroupLabels() len = %d, want 2", len(labels))
	}
	if labels[0] != "fiction" || labels[1] != "science" {
		t.Errorf("expandedGroupLabels() = %v, want [fiction science]", labels)
	}
}

func TestRestoreViewMode_NoViewState(t *testing.T) {
	state := SessionState{SelectedObjectID: "book/test"}
	vm := restoreViewMode(state, nil)
	if vm != nil {
		t.Error("restoreViewMode should return nil when no view state fields")
	}
}

func TestRestoreViewMode_EmptyViewName(t *testing.T) {
	state := SessionState{ViewTypeName: "book"} // ViewName is empty
	vm := restoreViewMode(state, nil)
	if vm != nil {
		t.Error("restoreViewMode should return nil when ViewName is empty")
	}
}

func TestRestoreViewMode_TypeDeleted(t *testing.T) {
	dir := t.TempDir()
	v := core.NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatalf("vault Init() error = %v", err)
	}
	if err := v.Open(); err != nil {
		t.Fatalf("vault Open() error = %v", err)
	}
	defer v.Close()

	state := SessionState{
		ViewTypeName: "nonexistent",
		ViewName:     "default",
	}
	vm := restoreViewMode(state, v)
	if vm != nil {
		t.Error("restoreViewMode should return nil when type doesn't exist")
	}
}

func TestRestoreViewMode_Success(t *testing.T) {
	dir := t.TempDir()
	v := core.NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatalf("vault Init() error = %v", err)
	}
	if err := v.Open(); err != nil {
		t.Fatalf("vault Open() error = %v", err)
	}
	defer v.Close()

	// Create a type with objects
	ts := &core.TypeSchema{Name: "book", Emoji: "📚"}
	if err := v.SaveType(ts); err != nil {
		t.Fatalf("SaveType() error = %v", err)
	}
	if _, err := v.NewObject("book", "test-a", ""); err != nil {
		t.Fatalf("NewObject() error = %v", err)
	}
	if _, err := v.NewObject("book", "test-b", ""); err != nil {
		t.Fatalf("NewObject() error = %v", err)
	}

	state := SessionState{
		ViewTypeName: "book",
		ViewName:     "default",
		ViewCursor:   1,
	}
	vm := restoreViewMode(state, v)
	if vm == nil {
		t.Fatal("restoreViewMode should return non-nil viewMode")
	}
	if vm.typeName != "book" {
		t.Errorf("typeName = %q, want %q", vm.typeName, "book")
	}
	if vm.viewName != "default" {
		t.Errorf("viewName = %q, want %q", vm.viewName, "default")
	}
	if vm.cursor != 1 {
		t.Errorf("cursor = %d, want 1", vm.cursor)
	}
	if len(vm.objects) != 2 {
		t.Errorf("objects len = %d, want 2", len(vm.objects))
	}
}

func TestRestoreViewMode_CursorClamped(t *testing.T) {
	dir := t.TempDir()
	v := core.NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatalf("vault Init() error = %v", err)
	}
	if err := v.Open(); err != nil {
		t.Fatalf("vault Open() error = %v", err)
	}
	defer v.Close()

	ts := &core.TypeSchema{Name: "book"}
	if err := v.SaveType(ts); err != nil {
		t.Fatalf("SaveType() error = %v", err)
	}
	if _, err := v.NewObject("book", "only-one", ""); err != nil {
		t.Fatalf("NewObject() error = %v", err)
	}

	state := SessionState{
		ViewTypeName: "book",
		ViewName:     "default",
		ViewCursor:   50, // way beyond actual items
		ViewScroll:   50,
	}
	vm := restoreViewMode(state, v)
	if vm == nil {
		t.Fatal("restoreViewMode should return non-nil viewMode")
	}
	totalRows := len(vm.visibleRows())
	if vm.cursor >= totalRows {
		t.Errorf("cursor = %d should be clamped below %d", vm.cursor, totalRows)
	}
	if vm.scroll >= totalRows {
		t.Errorf("scroll = %d should be clamped below %d", vm.scroll, totalRows)
	}
}

func TestRestoreViewMode_StaleExpandedGroups(t *testing.T) {
	dir := t.TempDir()
	v := core.NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatalf("vault Init() error = %v", err)
	}
	if err := v.Open(); err != nil {
		t.Fatalf("vault Open() error = %v", err)
	}
	defer v.Close()

	ts := &core.TypeSchema{
		Name: "book",
		Properties: []core.Property{
			{Name: "genre", Type: "string"},
		},
	}
	if err := v.SaveType(ts); err != nil {
		t.Fatalf("SaveType() error = %v", err)
	}
	obj, err := v.NewObject("book", "a", "")
	if err != nil {
		t.Fatalf("NewObject() error = %v", err)
	}
	obj.Properties["genre"] = "fiction"
	if err := v.SaveObject(obj); err != nil {
		t.Fatalf("SaveObject() error = %v", err)
	}

	// Save a view with group_by
	vc := &core.ViewConfig{Name: "by-genre", GroupBy: "genre"}
	if err := v.SaveView("book", vc); err != nil {
		t.Fatalf("SaveView() error = %v", err)
	}

	state := SessionState{
		ViewTypeName:       "book",
		ViewName:           "by-genre",
		ViewExpandedGroups: []string{"fiction", "deleted-group"},
	}
	vm := restoreViewMode(state, v)
	if vm == nil {
		t.Fatal("restoreViewMode should return non-nil viewMode")
	}

	// "fiction" group should exist and be expanded
	foundFiction := false
	for _, g := range vm.groups {
		if g.Label == "fiction" {
			if !g.Expanded {
				t.Error("fiction group should be expanded")
			}
			foundFiction = true
		}
	}
	if !foundFiction {
		t.Error("fiction group should exist in view groups")
	}
	// "deleted-group" should be silently ignored — no panic
}

func TestSaveAndLoad_ViewModeState(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, ".typemd"), 0755)

	original := SessionState{
		SelectedObjectID:   "book/test-01abc",
		ExpandedGroups:     []string{"book"},
		ViewTypeName:       "book",
		ViewName:           "by-rating",
		ViewCursor:         3,
		ViewScroll:         1,
		ViewExpandedGroups: []string{"5 stars"},
		PropsVisible:       true,
	}

	saveSessionState(dir, original)
	loaded := loadSessionState(dir)

	if loaded.ViewTypeName != "book" {
		t.Errorf("ViewTypeName = %q, want %q", loaded.ViewTypeName, "book")
	}
	if loaded.ViewName != "by-rating" {
		t.Errorf("ViewName = %q, want %q", loaded.ViewName, "by-rating")
	}
	if loaded.ViewCursor != 3 {
		t.Errorf("ViewCursor = %d, want 3", loaded.ViewCursor)
	}
	if loaded.ViewScroll != 1 {
		t.Errorf("ViewScroll = %d, want 1", loaded.ViewScroll)
	}
	if len(loaded.ViewExpandedGroups) != 1 || loaded.ViewExpandedGroups[0] != "5 stars" {
		t.Errorf("ViewExpandedGroups = %v, want [5 stars]", loaded.ViewExpandedGroups)
	}
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
