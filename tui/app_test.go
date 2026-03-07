package tui

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/typemd/typemd/core"
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
		focus:        focusLeft,
		groups:       groups,
		cursor:       0,
		selected:     groups[0].Objects[0],
		propsVisible: false,
		searchInput:  initSearchInput(),
		width:        120,
		height:       24,
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
	if updated.focus != focusBody {
		t.Errorf("focus = %d, want focusBody(%d)", updated.focus, focusBody)
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

func TestAdjustScrollOffset(t *testing.T) {
	// cursor above viewport — scroll up
	if o := adjustScrollOffset(2, 5, 10); o != 2 {
		t.Errorf("adjustScrollOffset(2,5,10) = %d, want 2", o)
	}
	// cursor below viewport — scroll down
	if o := adjustScrollOffset(15, 5, 10); o != 6 {
		t.Errorf("adjustScrollOffset(15,5,10) = %d, want 6", o)
	}
	// cursor within viewport — no change
	if o := adjustScrollOffset(7, 5, 10); o != 5 {
		t.Errorf("adjustScrollOffset(7,5,10) = %d, want 5", o)
	}
}

func TestScrollOffset_CursorFollows(t *testing.T) {
	var objects []*core.Object
	for i := 0; i < 30; i++ {
		objects = append(objects, &core.Object{
			ID: fmt.Sprintf("note/%03d", i), Type: "note", Filename: fmt.Sprintf("%03d", i),
		})
	}
	groups := buildGroups(objects)
	groups[0].Expanded = true
	m := model{
		focus:       focusLeft,
		groups:      groups,
		cursor:      0,
		searchInput: initSearchInput(),
		width:       80,
		height:      10,
	}

	for i := 0; i < 8; i++ {
		msg := tea.KeyMsg{Type: tea.KeyDown}
		newM, _ := m.Update(msg)
		m = newM.(model)
	}

	contentH := m.height - 3
	if m.scrollOffset > m.cursor || m.scrollOffset+contentH <= m.cursor {
		t.Errorf("cursor %d not visible with scrollOffset %d and height %d", m.cursor, m.scrollOffset, contentH)
	}
}

func TestRenderBody_WithContent(t *testing.T) {
	obj := &core.Object{
		ID:         "book/test",
		Type:       "book",
		Filename:   "test",
		Properties: map[string]any{"title": "Test"},
		Body:       "# Hello\nWorld",
	}
	result := renderBody(obj, 80)
	if !strings.Contains(result, "book/test") {
		t.Error("renderBody should contain object ID as title")
	}
	if !strings.Contains(result, "# Hello") {
		t.Error("renderBody should contain body content")
	}
	if strings.Contains(result, "title:") {
		t.Error("renderBody should NOT contain properties")
	}
}

func TestRenderBody_Nil(t *testing.T) {
	result := renderBody(nil, 80)
	if !strings.Contains(result, "Select an object") {
		t.Error("renderBody(nil, 80) should show placeholder")
	}
}

func TestRenderBody_EmptyBody(t *testing.T) {
	obj := &core.Object{ID: "book/test", Body: ""}
	result := renderBody(obj, 80)
	if !strings.Contains(result, "(empty)") {
		t.Error("renderBody with empty body should show (empty)")
	}
}

func TestRenderProperties_WithSchema(t *testing.T) {
	obj := &core.Object{
		ID:         "book/test",
		Properties: map[string]any{"title": "Go", "status": "reading"},
	}
	props := []core.DisplayProperty{
		{Key: "title", Value: "Go"},
		{Key: "status", Value: "reading"},
	}
	result := renderProperties(obj, props)
	if !strings.Contains(result, "title: Go") {
		t.Error("renderProperties should contain title property")
	}
	if !strings.Contains(result, "status: reading") {
		t.Error("renderProperties should contain status property")
	}
}

func TestRenderProperties_Nil(t *testing.T) {
	result := renderProperties(nil, nil)
	if result != "" {
		t.Errorf("renderProperties(nil) should return empty string, got %q", result)
	}
}

func TestRenderProperties_WithRelation(t *testing.T) {
	obj := &core.Object{
		ID:         "book/test",
		Properties: map[string]any{"author": "person/alan"},
	}
	props := []core.DisplayProperty{
		{Key: "author", Value: "person/alan", IsRelation: true},
	}
	result := renderProperties(obj, props)
	if !strings.Contains(result, "→") {
		t.Error("renderProperties should show arrow for relation properties")
	}
}

func TestRenderProperties_ReverseRelation(t *testing.T) {
	obj := &core.Object{
		ID:         "person/alan",
		Properties: map[string]any{},
	}
	props := []core.DisplayProperty{
		{Key: "author", Value: "book/test", IsReverse: true, FromID: "book/test"},
	}
	result := renderProperties(obj, props)
	if !strings.Contains(result, "←") {
		t.Error("renderProperties should show reverse arrow for reverse relations")
	}
}

func TestBuildGroups_DefaultCollapse(t *testing.T) {
	objects := []*core.Object{
		{ID: "book/a", Type: "book", Filename: "a"},
		{ID: "note/b", Type: "note", Filename: "b"},
	}
	groups := buildGroups(objects)
	for _, g := range groups {
		if g.Expanded {
			t.Errorf("expected %q to be collapsed by default", g.Name)
		}
	}
}

func TestModel_TabCyclesThreePanels(t *testing.T) {
	m := setupTestModel(t)
	m.propsVisible = true // enable props for three-panel cycling
	tab := tea.KeyMsg{Type: tea.KeyTab}

	// Left → Body
	newM, _ := m.Update(tab)
	m = newM.(model)
	if m.focus != focusBody {
		t.Errorf("after 1st tab: focus = %d, want focusBody(%d)", m.focus, focusBody)
	}

	// Body → Props
	newM, _ = m.Update(tab)
	m = newM.(model)
	if m.focus != focusProps {
		t.Errorf("after 2nd tab: focus = %d, want focusProps(%d)", m.focus, focusProps)
	}

	// Props → Left
	newM, _ = m.Update(tab)
	m = newM.(model)
	if m.focus != focusLeft {
		t.Errorf("after 3rd tab: focus = %d, want focusLeft(%d)", m.focus, focusLeft)
	}
}

func TestModel_TabSkipsPropsWhenHidden(t *testing.T) {
	m := setupTestModel(t) // propsVisible defaults to false
	tab := tea.KeyMsg{Type: tea.KeyTab}

	// Left → Body
	newM, _ := m.Update(tab)
	m = newM.(model)
	if m.focus != focusBody {
		t.Errorf("after 1st tab: focus = %d, want focusBody(%d)", m.focus, focusBody)
	}

	// Body → Left (skip Props)
	newM, _ = m.Update(tab)
	m = newM.(model)
	if m.focus != focusLeft {
		t.Errorf("after 2nd tab: focus = %d, want focusLeft(%d)", m.focus, focusLeft)
	}
}

func TestModel_PropsWidthDefault(t *testing.T) {
	m := setupTestModel(t)
	msg := tea.WindowSizeMsg{Width: 120, Height: 40}
	newM, _ := m.Update(msg)
	updated := newM.(model)

	if updated.propsWidth < 20 || updated.propsWidth > 40 {
		t.Errorf("propsWidth = %d, want between 20 and 40", updated.propsWidth)
	}
}

func TestModel_ThreePanelView_NotEmpty(t *testing.T) {
	m := setupTestModel(t)
	msg := tea.WindowSizeMsg{Width: 120, Height: 24}
	newM, _ := m.Update(msg)
	m = newM.(model)

	view := m.View()
	if view == "Loading..." {
		t.Error("View should not be Loading after WindowSizeMsg")
	}
	if len(view) == 0 {
		t.Error("View should not be empty")
	}
}

func TestModel_AutoHidePropsNarrowTerminal(t *testing.T) {
	m := setupTestModel(t)
	m.propsVisible = true // start with props visible
	msg := tea.WindowSizeMsg{Width: 50, Height: 24}
	newM, _ := m.Update(msg)
	updated := newM.(model)

	if updated.propsVisible {
		t.Error("propsVisible should be auto-hidden on narrow terminal (width=50)")
	}
}

func TestModel_ResizeLeftPanel(t *testing.T) {
	m := setupTestModel(t)
	sizeMsg := tea.WindowSizeMsg{Width: 120, Height: 24}
	newM, _ := m.Update(sizeMsg)
	m = newM.(model)
	m.focus = focusLeft

	before := m.leftW
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{']'}}
	newM, _ = m.Update(msg)
	m = newM.(model)

	if m.leftW != before+2 {
		t.Errorf("leftW = %d, want %d (grew by 2)", m.leftW, before+2)
	}
}

func TestModel_ResizeBodyPanel(t *testing.T) {
	m := setupTestModel(t)
	sizeMsg := tea.WindowSizeMsg{Width: 120, Height: 24}
	newM, _ := m.Update(sizeMsg)
	m = newM.(model)
	m.focus = focusBody
	m.propsVisible = true
	m.propsWidth = 30

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{']'}}
	newM, _ = m.Update(msg)
	m = newM.(model)

	// Growing body shrinks props
	if m.propsWidth != 28 {
		t.Errorf("propsWidth = %d, want 28 (body grow shrinks props)", m.propsWidth)
	}
}

func TestModel_ResizePanelGrow(t *testing.T) {
	m := setupTestModel(t)
	sizeMsg := tea.WindowSizeMsg{Width: 120, Height: 24}
	newM, _ := m.Update(sizeMsg)
	m = newM.(model)
	m.focus = focusProps

	before := m.propsWidth
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{']'}}
	newM, _ = m.Update(msg)
	m = newM.(model)

	if m.propsWidth != before+2 {
		t.Errorf("propsWidth = %d, want %d (grew by 2)", m.propsWidth, before+2)
	}
}

func TestModel_ResizePanelShrink(t *testing.T) {
	m := setupTestModel(t)
	sizeMsg := tea.WindowSizeMsg{Width: 120, Height: 24}
	newM, _ := m.Update(sizeMsg)
	m = newM.(model)
	m.focus = focusProps
	// Ensure propsWidth is above minimum so shrink has room
	m.propsWidth = 30

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'['}}
	newM, _ = m.Update(msg)
	m = newM.(model)

	if m.propsWidth != 28 {
		t.Errorf("propsWidth = %d, want 28 (shrunk by 2)", m.propsWidth)
	}
}

func TestModel_ToggleProperties(t *testing.T) {
	m := setupTestModel(t)
	sizeMsg := tea.WindowSizeMsg{Width: 120, Height: 24}
	newM, _ := m.Update(sizeMsg)
	m = newM.(model)

	if m.propsVisible {
		t.Fatal("propsVisible should default to false")
	}

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}}
	newM, _ = m.Update(msg)
	m = newM.(model)

	if !m.propsVisible {
		t.Error("propsVisible should be true after first toggle")
	}

	newM, _ = m.Update(msg)
	m = newM.(model)

	if m.propsVisible {
		t.Error("propsVisible should be false after second toggle")
	}
}

func TestModel_ToggleProps_MovesFocusWhenHiding(t *testing.T) {
	m := setupTestModel(t)
	sizeMsg := tea.WindowSizeMsg{Width: 120, Height: 24}
	newM, _ := m.Update(sizeMsg)
	m = newM.(model)
	m.propsVisible = true
	m.focus = focusProps

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}}
	newM, _ = m.Update(msg)
	m = newM.(model)

	if m.focus == focusProps {
		t.Error("focus should move away from focusProps when Properties is hidden")
	}
}

func TestModel_EnterEditMode_Body(t *testing.T) {
	m := setupTestModel(t)
	m.focus = focusBody

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}
	newM, _ := m.Update(msg)
	updated := newM.(model)

	if !updated.editMode {
		t.Error("editMode should be true after pressing 'e' when focused on body")
	}
}

func TestModel_EnterEditMode_Props(t *testing.T) {
	m := setupTestModel(t)
	m.focus = focusProps
	m.propsVisible = true

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}
	newM, _ := m.Update(msg)
	updated := newM.(model)

	if !updated.editMode {
		t.Error("editMode should be true after pressing 'e' when focused on props")
	}
}

func TestModel_NoEditMode_WhenFocusLeft(t *testing.T) {
	m := setupTestModel(t)
	m.focus = focusLeft

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}}
	newM, _ := m.Update(msg)
	updated := newM.(model)

	if updated.editMode {
		t.Error("editMode should NOT enter when focused on left panel")
	}
}

func TestModel_ExitEditMode_Esc(t *testing.T) {
	m := setupTestModel(t)
	m.focus = focusBody
	m.editMode = true

	msg := tea.KeyMsg{Type: tea.KeyEsc}
	newM, _ := m.Update(msg)
	updated := newM.(model)

	if updated.editMode {
		t.Error("editMode should be false after pressing Esc")
	}
}

func TestModel_EditMode_TabDoesNotSwitchPanel(t *testing.T) {
	m := setupTestModel(t)
	m.focus = focusBody
	m.editMode = true

	msg := tea.KeyMsg{Type: tea.KeyTab}
	newM, _ := m.Update(msg)
	updated := newM.(model)

	if updated.focus != focusBody {
		t.Errorf("focus should remain focusBody in edit mode, got %d", updated.focus)
	}
	if !updated.editMode {
		t.Error("editMode should remain true after pressing Tab in edit mode")
	}
}

func TestModel_EditMode_NavKeysDoNotScrollLeft(t *testing.T) {
	m := setupTestModel(t)
	m.focus = focusBody
	m.editMode = true
	initialCursor := m.cursor

	msg := tea.KeyMsg{Type: tea.KeyDown}
	newM, _ := m.Update(msg)
	updated := newM.(model)

	// In edit mode focused on body, 'j' should not move left panel cursor
	if updated.cursor != initialCursor {
		t.Errorf("cursor should not change in edit mode, got %d want %d", updated.cursor, initialCursor)
	}
}

func TestModel_View_ShowsEditModeIndicator(t *testing.T) {
	m := setupTestModel(t)
	sizeMsg := tea.WindowSizeMsg{Width: 120, Height: 24}
	newM, _ := m.Update(sizeMsg)
	m = newM.(model)
	m.focus = focusBody
	m.editMode = true

	view := m.View()
	if !strings.Contains(view, "EDIT") {
		t.Error("View should contain 'EDIT' mode indicator when editMode is true")
	}
}

func TestModel_View_ShowsViewModeIndicator(t *testing.T) {
	m := setupTestModel(t)
	sizeMsg := tea.WindowSizeMsg{Width: 120, Height: 24}
	newM, _ := m.Update(sizeMsg)
	m = newM.(model)
	m.editMode = false

	view := m.View()
	if !strings.Contains(view, "VIEW") {
		t.Error("View should contain 'VIEW' mode indicator when editMode is false")
	}
}

func TestModel_HelpToggle(t *testing.T) {
	m := setupTestModel(t)
	sizeMsg := tea.WindowSizeMsg{Width: 120, Height: 24}
	newM, _ := m.Update(sizeMsg)
	m = newM.(model)

	// Press ? to open help
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
	newM, _ = m.Update(msg)
	m = newM.(model)

	if !m.showHelp {
		t.Error("showHelp should be true after pressing ?")
	}

	// Press ? again to close
	newM, _ = m.Update(msg)
	m = newM.(model)

	if m.showHelp {
		t.Error("showHelp should be false after pressing ? again")
	}
}

func TestModel_HelpToggle_H(t *testing.T) {
	m := setupTestModel(t)

	// Press h to open help
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}}
	newM, _ := m.Update(msg)
	m = newM.(model)

	if !m.showHelp {
		t.Error("showHelp should be true after pressing h")
	}

	// Press h again to close
	newM, _ = m.Update(msg)
	m = newM.(model)

	if m.showHelp {
		t.Error("showHelp should be false after pressing h again")
	}
}

func TestModel_HelpEscClose(t *testing.T) {
	m := setupTestModel(t)

	// Open help
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
	newM, _ := m.Update(msg)
	m = newM.(model)

	// Press Esc to close
	esc := tea.KeyMsg{Type: tea.KeyEscape}
	newM, _ = m.Update(esc)
	m = newM.(model)

	if m.showHelp {
		t.Error("showHelp should be false after pressing Esc")
	}
}

func TestModel_HelpInterceptsKeys(t *testing.T) {
	m := setupTestModel(t)
	m.showHelp = true
	originalCursor := m.cursor

	// Press j (should be intercepted, cursor should not move)
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
	newM, _ := m.Update(msg)
	m = newM.(model)

	if m.cursor != originalCursor {
		t.Errorf("cursor = %d, want %d (help should intercept navigation keys)", m.cursor, originalCursor)
	}
	if !m.showHelp {
		t.Error("showHelp should still be true after pressing j")
	}
}

func TestModel_HelpView(t *testing.T) {
	m := setupTestModel(t)
	sizeMsg := tea.WindowSizeMsg{Width: 120, Height: 24}
	newM, _ := m.Update(sizeMsg)
	m = newM.(model)
	m.showHelp = true

	view := m.View()
	if !strings.Contains(view, "Keybindings") {
		t.Error("help view should contain 'Keybindings' title")
	}
	if !strings.Contains(view, "help") {
		t.Error("help view should contain help entry")
	}
}

func TestHelpEntries_NotEmpty(t *testing.T) {
	entries := helpEntries()
	if len(entries) == 0 {
		t.Error("helpEntries should return at least one entry")
	}
	for _, e := range entries {
		if e.Key == "" {
			t.Error("helpEntry Key should not be empty")
		}
		if e.Desc == "" {
			t.Error("helpEntry Desc should not be empty")
		}
	}
}

// setupTestModelWithVault creates a model backed by a real temporary vault with one book object.
func setupTestModelWithVault(t *testing.T) (model, *core.Object) {
	t.Helper()
	dir := t.TempDir()
	v := core.NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatalf("vault Init() error = %v", err)
	}
	if err := v.Open(); err != nil {
		t.Fatalf("vault Open() error = %v", err)
	}
	t.Cleanup(func() { v.Close() })

	obj, err := v.NewObject("book", "test-save")
	if err != nil {
		t.Fatalf("NewObject() error = %v", err)
	}
	obj.Properties["title"] = "Original Title"
	obj.Body = "Original body"
	if err := v.SaveObject(obj); err != nil {
		t.Fatalf("SaveObject() error = %v", err)
	}

	// Get file mtime so loadedModTime is correctly initialized
	var loadedMod time.Time
	if info, err := os.Stat(v.ObjectPath(obj.Type, obj.Filename)); err == nil {
		loadedMod = info.ModTime()
	}

	groups := buildGroups([]*core.Object{obj})
	groups[0].Expanded = true
	m := model{
		vault:         v,
		focus:         focusBody,
		groups:        groups,
		cursor:        1,
		selected:      obj,
		propsVisible:  false,
		searchInput:   initSearchInput(),
		width:         120,
		height:        24,
		loadedModTime: loadedMod,
	}
	return m, obj
}

func TestModel_ExitEditMode_NoSaveWhenNotDirty(t *testing.T) {
	m, _ := setupTestModelWithVault(t)
	m.editMode = true
	m.dirty = false

	msg := tea.KeyMsg{Type: tea.KeyEsc}
	newM, _ := m.Update(msg)
	updated := newM.(model)

	if updated.editMode {
		t.Error("editMode should be false after esc")
	}
	if updated.dirty {
		t.Error("dirty should remain false when nothing changed")
	}
}

func TestModel_ExitEditMode_SavesWhenDirty(t *testing.T) {
	m, obj := setupTestModelWithVault(t)
	m.editMode = true
	m.dirty = true
	m.selected.Body = "Updated body"

	msg := tea.KeyMsg{Type: tea.KeyEsc}
	newM, _ := m.Update(msg)
	updated := newM.(model)

	if updated.editMode {
		t.Error("editMode should be false after esc")
	}
	if updated.dirty {
		t.Error("dirty should be false after successful save")
	}
	if updated.saveErr != "" {
		t.Errorf("saveErr should be empty after successful save, got %q", updated.saveErr)
	}

	// Verify file was actually updated
	reloaded, err := m.vault.GetObject(obj.ID)
	if err != nil {
		t.Fatalf("GetObject() error = %v", err)
	}
	if reloaded.Body != "Updated body" {
		t.Errorf("body = %q, want %q", reloaded.Body, "Updated body")
	}
}

func TestModel_SaveError_ShowsInStatusBar(t *testing.T) {
	m, _ := setupTestModelWithVault(t)
	m.saveErr = "Save failed: disk full"
	m.width = 120
	m.height = 24

	view := m.View()
	if !strings.Contains(view, "Save failed") {
		t.Error("View should contain save error message when saveErr is set")
	}
}

func TestModel_SkipNextReload_SuppressesRefresh(t *testing.T) {
	m, _ := setupTestModelWithVault(t)
	m.skipNextReload = true

	originalSelected := m.selected

	newM, cmd := m.Update(fileChangedMsg{})
	updated := newM.(model)

	if updated.skipNextReload {
		t.Error("skipNextReload should be cleared after fileChangedMsg")
	}
	if updated.selected != originalSelected {
		t.Error("selected should not change when skipNextReload suppresses reload")
	}
	if cmd == nil {
		t.Error("should restart watcher (non-nil cmd) after suppressed reload")
	}
}

func TestModel_ConcurrentEdit_DetectedBeforeSave(t *testing.T) {
	m, obj := setupTestModelWithVault(t)

	// Simulate: file was modified externally after we loaded it
	// Set loadedModTime to the past so the file's real mtime is "newer"
	m.loadedModTime = m.loadedModTime.Add(-2 * time.Second)
	m.editMode = true
	m.dirty = true
	m.selected.Body = "Local changes"

	msg := tea.KeyMsg{Type: tea.KeyEsc}
	newM, _ := m.Update(msg)
	updated := newM.(model)

	if !updated.saveConflict {
		t.Error("saveConflict should be true when external edit detected")
	}
	if updated.saveErr == "" {
		t.Error("saveErr should explain the conflict")
	}
	if updated.dirty {
		// dirty stays true because save was blocked
	}

	// File should NOT be overwritten
	reloaded, _ := m.vault.GetObject(obj.ID)
	if reloaded.Body == "Local changes" {
		t.Error("file should NOT be overwritten when conflict detected")
	}
}

func TestModel_ConcurrentEdit_OverwriteWithY(t *testing.T) {
	m, obj := setupTestModelWithVault(t)
	m.loadedModTime = m.loadedModTime.Add(-2 * time.Second)
	m.saveConflict = true
	m.dirty = true
	m.selected.Body = "Force saved body"

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}}
	newM, _ := m.Update(msg)
	updated := newM.(model)

	if updated.saveConflict {
		t.Error("saveConflict should be cleared after y")
	}
	if updated.dirty {
		t.Error("dirty should be false after force save")
	}

	reloaded, _ := m.vault.GetObject(obj.ID)
	if reloaded.Body != "Force saved body" {
		t.Errorf("body = %q, want %q", reloaded.Body, "Force saved body")
	}
}

func TestModel_ConcurrentEdit_ReloadWithN(t *testing.T) {
	m, obj := setupTestModelWithVault(t)
	m.saveConflict = true
	m.dirty = true
	m.selected.Body = "Discarded local changes"

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
	newM, _ := m.Update(msg)
	updated := newM.(model)

	if updated.saveConflict {
		t.Error("saveConflict should be cleared after n")
	}
	if updated.dirty {
		t.Error("dirty should be false after reload")
	}

	// selected should reflect what's on disk (original body, not local changes)
	if updated.selected != nil && updated.selected.Body == "Discarded local changes" {
		t.Error("selected body should be reloaded from disk, not local changes")
	}
	_ = obj
}
