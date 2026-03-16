package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

// helper: send a WindowSizeMsg so View() doesn't return "Loading..."
func initSize(t *testing.T, m model) model {
	t.Helper()
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 24})
	return newM.(model)
}

// helper: type a string character by character
func typeString(t *testing.T, m model, s string) model {
	t.Helper()
	for _, ch := range s {
		newM, _ := m.Update(tea.KeyPressMsg{Code: ch, Text: string(ch)})
		m = newM.(model)
	}
	return m
}

// helper: press a special key
func pressKey(t *testing.T, m model, code rune) model {
	t.Helper()
	newM, _ := m.Update(tea.KeyPressMsg{Code: code})
	return newM.(model)
}

// helper: press a named key (e.g. "n", "N")
func pressNamedKey(t *testing.T, m model, name string) model {
	t.Helper()
	var code rune
	if len(name) == 1 {
		code = rune(name[0])
	}
	newM, _ := m.Update(tea.KeyPressMsg{Code: code, Text: name})
	return newM.(model)
}

// =============================================================================
// PR Manual Test 1: n → template selection → name input → object created, edit mode
// =============================================================================

func TestE2E_CreateAndEdit_WithTemplateSelection(t *testing.T) {
	m := setupCreateTestModelWithTemplates(t, map[string]string{
		"review":  "---\ntitle: Review Template\n---\n## Review\n",
		"summary": "---\ntitle: Summary Template\n---\n## Summary\n",
	})
	m = initSize(t, m)

	// Step 1: Press n on type header → template selection appears
	m = pressNamedKey(t, m, "n")

	if m.create == nil {
		t.Fatal("create should be active after pressing n")
	}
	if m.create.step != createStepTemplate {
		t.Fatalf("step = %d, want createStepTemplate (multiple templates)", m.create.step)
	}

	// Verify View shows template selection UI
	view := m.View().Content
	if !strings.Contains(view, "Select template:") {
		t.Error("View should show 'Select template:' during template step")
	}
	if !strings.Contains(view, "review") {
		t.Error("View should show 'review' template option")
	}
	if !strings.Contains(view, noneOption) {
		t.Error("View should show '(none)' option")
	}

	// Verify help bar shows template selection hints
	if !strings.Contains(view, "NEW OBJECT") {
		t.Error("help bar should show [NEW OBJECT] during template selection")
	}
	if !strings.Contains(view, "↑↓") {
		t.Error("help bar should show ↑↓ navigation hint")
	}

	// Step 2: Select "review" template (first option, press Enter)
	m = pressKey(t, m, tea.KeyEnter)

	if m.create == nil {
		t.Fatal("create should still be active after template selection")
	}
	if m.create.step != createStepName {
		t.Fatalf("step = %d, want createStepName after template selection", m.create.step)
	}
	if m.create.template != "review" {
		t.Errorf("template = %q, want 'review'", m.create.template)
	}

	// Verify View shows name input
	view = m.View().Content
	if !strings.Contains(view, "New book:") {
		t.Error("View should show 'New book:' name input prompt")
	}
	if !strings.Contains(view, "create & edit") {
		t.Error("help bar should show 'create & edit' hint for single mode")
	}

	// Step 3: Type name and press Enter
	m = typeString(t, m, "my-review-book")
	m = pressKey(t, m, tea.KeyEnter)

	// Verify: creation done, edit mode active
	if m.create != nil {
		t.Error("create should be nil after single mode creation")
	}
	if m.selected == nil {
		t.Fatal("an object should be selected after creation")
	}
	if !strings.Contains(m.selected.ID, "book/my-review-book") {
		t.Errorf("selected.ID = %q, want to contain 'book/my-review-book'", m.selected.ID)
	}
	if m.focus != focusBody {
		t.Errorf("focus = %d, want focusBody (%d) — should auto-focus body", m.focus, focusBody)
	}
	if !m.editMode {
		t.Error("editMode should be true — should auto-enter edit mode")
	}
	if m.rightPanel != panelObject {
		t.Errorf("rightPanel = %d, want panelObject", m.rightPanel)
	}

	// Verify the template was applied (body should contain template content)
	if !strings.Contains(m.selected.Body, "## Review") {
		t.Errorf("body = %q, want to contain '## Review' from template", m.selected.Body)
	}

	// Verify the object actually exists on disk
	reloaded, err := m.vault.GetObject(m.selected.ID)
	if err != nil {
		t.Fatalf("GetObject(%q) error = %v — object should exist on disk", m.selected.ID, err)
	}
	if !strings.Contains(reloaded.Body, "## Review") {
		t.Errorf("reloaded body = %q, want template body", reloaded.Body)
	}
}

func TestE2E_CreateAndEdit_NoTemplates(t *testing.T) {
	m := setupCreateTestModel(t)
	m = initSize(t, m)

	// Press n → should go directly to name step (no templates)
	m = pressNamedKey(t, m, "n")

	if m.create == nil {
		t.Fatal("create should be active")
	}
	if m.create.step != createStepName {
		t.Fatalf("step = %d, want createStepName (no templates to select)", m.create.step)
	}

	// Type name, press Enter
	m = typeString(t, m, "simple-book")
	m = pressKey(t, m, tea.KeyEnter)

	if m.create != nil {
		t.Error("create should be nil after creation")
	}
	if !m.editMode {
		t.Error("should be in edit mode")
	}
	if m.selected == nil || !strings.Contains(m.selected.ID, "book/simple-book") {
		t.Error("should have selected the new object")
	}
}

// =============================================================================
// PR Manual Test 2: N → name input → enter creates with flash → stays → esc exits
// =============================================================================

func TestE2E_QuickCreate_BatchFlow(t *testing.T) {
	m := setupCreateTestModel(t)
	m = initSize(t, m)

	// Step 1: Press N → Quick Create mode
	m = pressNamedKey(t, m, "N")

	if m.create == nil {
		t.Fatal("create should be active after pressing N")
	}
	if m.create.mode != createModeBatch {
		t.Fatalf("mode = %d, want createModeBatch", m.create.mode)
	}
	if m.create.step != createStepName {
		t.Fatalf("step = %d, want createStepName", m.create.step)
	}

	// Verify help bar shows QUICK CREATE
	view := m.View().Content
	if !strings.Contains(view, "QUICK CREATE") {
		t.Error("help bar should show [QUICK CREATE]")
	}
	if !strings.Contains(view, "esc: done") {
		t.Error("help bar should show 'esc: done'")
	}

	// Step 2: Create first object
	m = typeString(t, m, "batch-one")
	m = pressKey(t, m, tea.KeyEnter)

	if m.create == nil {
		t.Fatal("create should still be active in batch mode after creation")
	}
	if m.create.nameInput.Value() != "" {
		t.Errorf("input should be cleared, got %q", m.create.nameInput.Value())
	}
	if m.create.flash == "" {
		t.Error("flash message should be visible")
	}
	if !strings.Contains(m.create.flash, "batch-one") {
		t.Errorf("flash = %q, should mention created object name", m.create.flash)
	}

	// Verify flash appears in View
	view = m.View().Content
	if !strings.Contains(view, "Created") {
		t.Error("View should show creation success flash")
	}

	// Step 3: Create second object
	m = typeString(t, m, "batch-two")
	m = pressKey(t, m, tea.KeyEnter)

	if m.create == nil {
		t.Fatal("create should still be active")
	}
	if !strings.Contains(m.create.flash, "batch-two") {
		t.Errorf("flash = %q, should show latest creation", m.create.flash)
	}

	lastCreatedID := m.create.lastObj.ID

	// Step 4: Press Esc to exit
	m = pressKey(t, m, tea.KeyEscape)

	if m.create != nil {
		t.Error("create should be nil after Esc")
	}

	// Should select the last created object
	if m.selected == nil {
		t.Fatal("should have a selected object after batch exit")
	}
	if m.selected.ID != lastCreatedID {
		t.Errorf("selected.ID = %q, want %q (last created)", m.selected.ID, lastCreatedID)
	}

	// Verify both objects exist on disk
	if _, err := m.vault.GetObject(lastCreatedID); err != nil {
		t.Errorf("last created object should exist on disk: %v", err)
	}

	// Should NOT be in edit mode (batch mode doesn't auto-edit)
	if m.editMode {
		t.Error("should NOT be in edit mode after batch exit")
	}
}

func TestE2E_QuickCreate_WithTemplates(t *testing.T) {
	m := setupCreateTestModelWithTemplates(t, map[string]string{
		"review":  "---\ntitle: Review\n---\n",
		"summary": "---\ntitle: Summary\n---\n",
	})
	m = initSize(t, m)

	// Press N → template selection first
	m = pressNamedKey(t, m, "N")

	if m.create == nil || m.create.step != createStepTemplate {
		t.Fatal("should show template selection in batch mode too")
	}

	// Select "review" → name input
	m = pressKey(t, m, tea.KeyEnter)

	if m.create.step != createStepName {
		t.Fatal("should be at name step after template selection")
	}

	// Create two objects with the same template
	m = typeString(t, m, "book-a")
	m = pressKey(t, m, tea.KeyEnter)

	m = typeString(t, m, "book-b")
	m = pressKey(t, m, tea.KeyEnter)

	// Both should use the "review" template
	m = pressKey(t, m, tea.KeyEscape)

	// Verify the template was reused (check the last created object)
	if m.selected == nil {
		t.Fatal("should have selected object after exit")
	}
	obj, err := m.vault.GetObject(m.selected.ID)
	if err != nil {
		t.Fatalf("GetObject error: %v", err)
	}
	if obj.Properties["title"] != "Review" {
		t.Errorf("properties[title] = %v, want 'Review' from template", obj.Properties["title"])
	}
}

// =============================================================================
// PR Manual Test 3: n/N disabled in --readonly mode
// =============================================================================

func TestE2E_ReadOnly_CreationDisabled(t *testing.T) {
	m := setupCreateTestModel(t)
	m = initSize(t, m)
	m.readOnly = true

	// Verify [READONLY] in status bar
	view := m.View().Content
	if !strings.Contains(view, "READONLY") {
		t.Error("View should show READONLY indicator")
	}

	// Press n — should do nothing
	m = pressNamedKey(t, m, "n")
	if m.create != nil {
		t.Error("n should be ignored in readonly mode")
	}

	// Press N — should do nothing
	m = pressNamedKey(t, m, "N")
	if m.create != nil {
		t.Error("N should be ignored in readonly mode")
	}

	// Verify help popup hides n/N keybindings
	entries := helpEntries(true)
	for _, e := range entries {
		if e.Key == "n" || e.Key == "N" {
			t.Errorf("helpEntries(readOnly=true) should NOT contain %q keybinding", e.Key)
		}
	}
}

// =============================================================================
// PR Manual Test 4: Unique constraint inline error below input
// =============================================================================

func TestE2E_UniqueConstraint_InlineError(t *testing.T) {
	m := setupCreateTestModel(t)

	// Add unique type
	os.WriteFile(
		filepath.Join(m.vault.TypesDir(), "tag.yaml"),
		[]byte("name: tag\nunique: true\nplural: tags\n"),
		0644,
	)
	objects, _ := m.vault.QueryObjects("")
	m.groups = buildGroups(objects, m.vault)
	m = initSize(t, m)

	// Find tag group index
	tagIdx := -1
	for i, g := range m.groups {
		if g.Name == "tag" {
			tagIdx = i
			break
		}
	}
	if tagIdx < 0 {
		t.Fatal("tag type not found")
	}

	// Move cursor to tag header
	rows := m.currentRows()
	for i, row := range rows {
		if row.Kind == rowHeader && row.GroupIndex == tagIdx {
			m.cursor = i
			break
		}
	}

	// Create first tag
	m = pressNamedKey(t, m, "n")
	m = typeString(t, m, "important")
	m = pressKey(t, m, tea.KeyEnter)

	// First creation should succeed → edit mode
	if m.create != nil {
		if m.create.errMsg != "" {
			t.Fatalf("first creation should succeed, got error: %s", m.create.errMsg)
		}
		t.Fatal("create should be nil after successful single-mode creation")
	}
	if !m.editMode {
		t.Fatal("should be in edit mode after first creation")
	}

	// Exit edit mode first
	m = pressKey(t, m, tea.KeyEscape)
	m.focus = focusLeft

	// Move cursor back to tag header
	rows = m.currentRows()
	for i, row := range rows {
		if row.Kind == rowHeader && row.GroupIndex == tagIdx {
			m.cursor = i
			break
		}
	}

	// Try to create duplicate
	m = pressNamedKey(t, m, "n")
	if m.create == nil {
		t.Fatal("create should be active")
	}
	m = typeString(t, m, "important")
	m = pressKey(t, m, tea.KeyEnter)

	// Should show error, NOT create object
	if m.create == nil {
		t.Fatal("create should still be active on duplicate error")
	}
	if m.create.errMsg == "" {
		t.Fatal("errMsg should contain duplicate error")
	}
	if !strings.Contains(m.create.errMsg, "already exists") {
		t.Errorf("errMsg = %q, want to contain 'already exists'", m.create.errMsg)
	}

	// Verify error shows in View
	view := m.View().Content
	if !strings.Contains(view, "✗") {
		t.Error("View should show ✗ error indicator")
	}
	if !strings.Contains(view, "already exists") {
		t.Error("View should show 'already exists' error text")
	}

	// Typing clears the error
	m = typeString(t, m, "-v2")
	if m.create.errMsg != "" {
		t.Errorf("errMsg should be cleared after typing, got %q", m.create.errMsg)
	}

	// Verify error gone from View
	view = m.View().Content
	if strings.Contains(view, "✗") {
		t.Error("View should no longer show ✗ after typing")
	}
}
