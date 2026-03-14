package tui

import (
	"testing"

	"github.com/typemd/typemd/core"
	tea "charm.land/bubbletea/v2"
)

func testSchema() *core.TypeSchema {
	return &core.TypeSchema{
		Name:   "book",
		Plural: "books",
		Emoji:  "📖",
		Properties: []core.Property{
			{Name: "author", Type: "relation", Target: "person", Pin: 1, Emoji: "👤"},
			{Name: "genre", Type: "select", Pin: 2, Emoji: "🎭"},
			{Name: "rating", Type: "number"},
			{Name: "isbn", Type: "string"},
		},
	}
}

func TestTypeEditor_CursorNavigation(t *testing.T) {
	te := newTypeEditor(testSchema(), "book", false, nil)

	// Start at 0 (Name)
	if te.cursor != 0 {
		t.Errorf("initial cursor = %d, want 0", te.cursor)
	}

	// Move down through meta fields
	te.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	if te.cursor != 1 {
		t.Errorf("cursor after down = %d, want 1", te.cursor)
	}

	// Move to last item
	total := te.totalItems()
	for i := 0; i < total; i++ {
		te.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	}
	if te.cursor != total-1 {
		t.Errorf("cursor at end = %d, want %d", te.cursor, total-1)
	}

	// Can't go past end
	te.Update(tea.KeyPressMsg{Code: tea.KeyDown})
	if te.cursor != total-1 {
		t.Errorf("cursor past end = %d, want %d", te.cursor, total-1)
	}

	// Move back to start
	for i := 0; i < total; i++ {
		te.Update(tea.KeyPressMsg{Code: tea.KeyUp})
	}
	if te.cursor != 0 {
		t.Errorf("cursor at start = %d, want 0", te.cursor)
	}

	// Can't go before start
	te.Update(tea.KeyPressMsg{Code: tea.KeyUp})
	if te.cursor != 0 {
		t.Errorf("cursor before start = %d, want 0", te.cursor)
	}
}

func TestTypeEditor_TotalItems(t *testing.T) {
	te := newTypeEditor(testSchema(), "book", false, nil)
	// 4 meta fields + 4 properties + 1 add property = 9
	if te.totalItems() != 9 {
		t.Errorf("totalItems = %d, want 9", te.totalItems())
	}
}

func TestTypeEditor_DisplayItems(t *testing.T) {
	te := newTypeEditor(testSchema(), "book", false, nil)
	items := te.displayItems()

	// Should have 9 items: 4 meta + 2 pinned + 2 unpinned + 1 add property
	if len(items) != 9 {
		t.Fatalf("len(displayItems) = %d, want 9", len(items))
	}

	// First 4 are meta sentinels
	for i := 0; i < 4; i++ {
		if items[i] >= 0 {
			t.Errorf("items[%d] = %d, want negative sentinel", i, items[i])
		}
	}

	// Items 4-5 should be pinned properties (indices 0,1 — author pin:1, genre pin:2)
	if items[4] != 0 || items[5] != 1 {
		t.Errorf("pinned items = [%d, %d], want [0, 1]", items[4], items[5])
	}

	// Items 6-7 should be unpinned properties (indices 2,3 — rating, isbn)
	if items[6] != 2 || items[7] != 3 {
		t.Errorf("unpinned items = [%d, %d], want [2, 3]", items[6], items[7])
	}
}

func TestTypeEditor_NameNotEditable(t *testing.T) {
	te := newTypeEditor(testSchema(), "book", false, nil)
	te.cursor = 0 // Name field

	te.Update(keyMsg("e"))

	if te.mode != teModeView {
		t.Errorf("mode after e on Name = %d, want teModeView", te.mode)
	}
}

func TestTypeEditor_UniqueToggle(t *testing.T) {
	te := newTypeEditor(testSchema(), "book", false, nil)
	te.cursor = 3 // Unique field

	if te.schema.Unique {
		t.Error("Unique should start as false")
	}

	te.Update(keyMsg("e"))
	if !te.schema.Unique {
		t.Error("Unique should be true after toggle")
	}

	te.Update(keyMsg("e"))
	if te.schema.Unique {
		t.Error("Unique should be false after second toggle")
	}
}

func TestTypeEditor_PinToggle(t *testing.T) {
	te := newTypeEditor(testSchema(), "book", false, nil)

	// Move to "rating" (unpinned) — it's at display index 6
	te.cursor = 6

	items := te.displayItems()
	propIdx := items[6]
	if te.schema.Properties[propIdx].Pin != 0 {
		t.Fatalf("rating should start unpinned")
	}

	te.Update(keyMsg("p"))
	if te.schema.Properties[propIdx].Pin != 3 {
		t.Errorf("rating pin after pin = %d, want 3 (max+1)", te.schema.Properties[propIdx].Pin)
	}

	// Unpin it
	// After pin toggle, cursor position might change due to reorder. Let's find rating again.
	items = te.displayItems()
	for i, item := range items {
		if item >= 0 && te.schema.Properties[item].Name == "rating" {
			te.cursor = i
			break
		}
	}
	te.Update(keyMsg("p"))
	if te.schema.Properties[propIdx].Pin != 0 {
		t.Errorf("rating pin after unpin = %d, want 0", te.schema.Properties[propIdx].Pin)
	}
}

func TestTypeEditor_EscReturnsFocusLeft(t *testing.T) {
	te := newTypeEditor(testSchema(), "book", false, nil)

	result, cmd := te.Update(keyMsg("esc"))
	if result == nil {
		t.Error("esc should not close editor (returns non-nil)")
	}
	if cmd == nil {
		t.Error("esc should return a command (focusLeftMsg)")
	}
}

func TestTypeEditor_DeletePropertyConfirm(t *testing.T) {
	te := newTypeEditor(testSchema(), "book", false, nil)
	te.cursor = 7 // last unpinned property (isbn, index 3)

	te.Update(keyMsg("d"))
	if te.mode != teModeDeleteProp {
		t.Fatalf("mode after d = %d, want teModeDeleteProp", te.mode)
	}

	te.Update(keyMsg("y"))
	if len(te.schema.Properties) != 3 {
		t.Errorf("properties after delete = %d, want 3", len(te.schema.Properties))
	}
	if te.mode != teModeView {
		t.Errorf("mode after confirm = %d, want teModeView", te.mode)
	}
}

func TestTypeEditor_DeletePropertyCancel(t *testing.T) {
	te := newTypeEditor(testSchema(), "book", false, nil)
	te.cursor = 7

	te.Update(keyMsg("d"))
	te.Update(keyMsg("n"))

	if len(te.schema.Properties) != 4 {
		t.Errorf("properties after cancel = %d, want 4", len(te.schema.Properties))
	}
}

func TestTypeEditor_MoveMode(t *testing.T) {
	te := newTypeEditor(testSchema(), "book", false, nil)

	// Move to first unpinned property (rating, display index 6)
	te.cursor = 6

	te.Update(keyMsg("m"))
	if te.mode != teModeMove {
		t.Fatalf("mode after m = %d, want teModeMove", te.mode)
	}

	// Move down (swap rating and isbn)
	te.Update(keyMsg("j"))

	// Exit move mode
	te.Update(keyMsg("enter"))
	if te.mode != teModeView {
		t.Errorf("mode after enter = %d, want teModeView", te.mode)
	}
}

func TestTypeEditor_HelpBar(t *testing.T) {
	te := newTypeEditor(testSchema(), "book", false, nil)

	bar := te.HelpBar()
	if bar == "" {
		t.Error("help bar should not be empty in view mode")
	}

	te.mode = teModeMove
	bar = te.HelpBar()
	if bar == "" {
		t.Error("help bar should not be empty in move mode")
	}
}

func TestTypeEditor_ViewRendersAllSections(t *testing.T) {
	te := newTypeEditor(testSchema(), "book", false, nil)
	view := te.View()

	for _, want := range []string{"Properties", "author", "isbn"} {
		if !containsStr(view, want) {
			t.Errorf("view missing %q", want)
		}
	}
}

// ── Add Property Wizard tests ────────────────────────────────────────────────

func TestTypeEditor_WizardStartsOnA(t *testing.T) {
	te := newTypeEditor(testSchema(), "book", false, nil)
	te.Update(keyMsg("a"))
	if te.mode != teModeAddWizard {
		t.Errorf("mode after a = %d, want teModeAddWizard", te.mode)
	}
	if te.wizard == nil {
		t.Fatal("wizard should not be nil")
	}
	if te.wizard.step != wizStepName {
		t.Errorf("wizard step = %d, want wizStepName", te.wizard.step)
	}
}

func TestTypeEditor_WizardCancelOnEsc(t *testing.T) {
	te := newTypeEditor(testSchema(), "book", false, nil)
	te.Update(keyMsg("a"))
	te.Update(keyMsg("esc"))
	if te.mode != teModeView {
		t.Errorf("mode after esc = %d, want teModeView", te.mode)
	}
	if te.wizard != nil {
		t.Error("wizard should be nil after cancel")
	}
}

func TestTypeEditor_WizardDuplicateName(t *testing.T) {
	te := newTypeEditor(testSchema(), "book", false, nil)
	te.Update(keyMsg("a"))

	// Set input value directly (textinput component handles char input internally)
	te.wizard.nameInput.SetValue("author")
	te.Update(keyMsg("enter"))

	// Should still be on step 1 with error
	if te.wizard.step != wizStepName {
		t.Errorf("step = %d, want wizStepName (duplicate rejected)", te.wizard.step)
	}
	if te.saveErr == "" {
		t.Error("expected error message for duplicate name")
	}
}

func TestTypeEditor_WizardTypeSelection(t *testing.T) {
	te := newTypeEditor(testSchema(), "book", false, nil)
	te.Update(keyMsg("a"))

	// Set name directly and advance
	te.wizard.nameInput.SetValue("pages")
	te.Update(keyMsg("enter"))

	if te.wizard.step != wizStepType {
		t.Fatalf("step = %d, want wizStepType", te.wizard.step)
	}

	// Select "number" (index 1)
	te.Update(keyMsg("j"))
	te.Update(keyMsg("enter"))

	// Should create the property (number doesn't need extra steps)
	if te.mode != teModeView {
		t.Errorf("mode = %d, want teModeView after simple type", te.mode)
	}

	// Check property was added
	found := false
	for _, p := range te.schema.Properties {
		if p.Name == "pages" && p.Type == "number" {
			found = true
			break
		}
	}
	if !found {
		t.Error("property 'pages' not found in schema after wizard")
	}
}

func TestTypeEditor_WizardViewRenders(t *testing.T) {
	te := newTypeEditor(testSchema(), "book", false, nil)
	te.Update(keyMsg("a"))

	view := te.View()
	if !containsStr(view, "Add Property") {
		t.Error("wizard view should contain 'Add Property'")
	}
	if !containsStr(view, "Step 1") {
		t.Error("wizard view should contain 'Step 1'")
	}
}

// ── Panel mode tests ────────────────────────────────────────────────────────

func TestVisibleRows_IncludesNewType(t *testing.T) {
	groups := []typeGroup{
		{Name: "book", Objects: nil, Expanded: false},
		{Name: "note", Objects: nil, Expanded: false},
	}
	rows := visibleRows(groups)
	last := rows[len(rows)-1]
	if last.Kind != rowNewType {
		t.Errorf("last row Kind = %d, want rowNewType", last.Kind)
	}
}

// ── Helpers ─────────────────────────────────────────────────────────────────

func keyMsg(key string) tea.KeyPressMsg {
	if len(key) == 1 {
		return tea.KeyPressMsg{Code: rune(key[0])}
	}
	switch key {
	case "esc":
		return tea.KeyPressMsg{Code: tea.KeyEscape}
	case "enter":
		return tea.KeyPressMsg{Code: tea.KeyEnter}
	}
	return tea.KeyPressMsg{}
}

func containsStr(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && (s == substr || len(s) > len(substr) && findStr(s, substr))
}

func findStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
