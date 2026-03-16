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
// E2E Test 1: n → title panel creation form → template selection → name → edit
// =============================================================================

func TestE2E_CreateAndEdit_WithTemplateSelection(t *testing.T) {
	m := setupCreateTestModelWithTemplates(t, map[string]string{
		"review":  "---\ntitle: Review Template\n---\n## Review\n",
		"summary": "---\ntitle: Summary Template\n---\n## Summary\n",
	})
	m = initSize(t, m)

	// Press n → creation form appears in title panel
	m = pressNamedKey(t, m, "n")
	if m.create == nil {
		t.Fatal("create should be active")
	}

	// Verify View shows title panel with creation form (not sidebar)
	view := m.View().Content
	if !strings.Contains(view, "📝") {
		t.Error("title panel should show template selector icon")
	}
	// Sidebar should NOT contain creation UI
	if strings.Contains(view, "Select template:") {
		t.Error("sidebar should NOT show template selection list")
	}

	// Verify help bar
	if !strings.Contains(view, "NEW OBJECT") {
		t.Error("help bar should show [NEW OBJECT]")
	}

	// Tab to template field, switch to "summary"
	m = pressKey(t, m, tea.KeyTab)
	if m.create.field != createFieldTemplate {
		t.Fatal("should be on template field after Tab")
	}

	m = pressKey(t, m, tea.KeyRight)
	if m.create.currentTemplateName() != "summary" {
		t.Errorf("template = %q, want 'summary'", m.create.currentTemplateName())
	}

	// Verify live preview updated in body viewport
	bodyContent := m.bodyViewport.View()
	if !strings.Contains(bodyContent, "Summary") {
		t.Error("body viewport should show summary template preview")
	}

	// Tab back to name, type name
	m = pressKey(t, m, tea.KeyTab)
	m = typeString(t, m, "my-summary-book")

	// Enter to create
	m = pressKey(t, m, tea.KeyEnter)

	if m.create != nil {
		t.Error("create should be nil after creation")
	}
	if m.selected == nil {
		t.Fatal("should have selected object")
	}
	if !strings.Contains(m.selected.ID, "book/my-summary-book") {
		t.Errorf("ID = %q, want to contain 'book/my-summary-book'", m.selected.ID)
	}
	if !m.editMode {
		t.Error("should be in edit mode")
	}
	if !strings.Contains(m.selected.Body, "## Summary") {
		t.Errorf("body = %q, want summary template body", m.selected.Body)
	}
}

func TestE2E_CreateAndEdit_NoTemplates(t *testing.T) {
	m := setupCreateTestModel(t)
	m = initSize(t, m)

	m = pressNamedKey(t, m, "n")
	if m.create == nil {
		t.Fatal("create should be active")
	}

	// No template selector visible
	view := m.View().Content
	if strings.Contains(view, "📝") {
		t.Error("should not show template icon when no templates")
	}

	m = typeString(t, m, "simple-book")
	m = pressKey(t, m, tea.KeyEnter)

	if !m.editMode {
		t.Error("should be in edit mode")
	}
	if m.selected == nil || !strings.Contains(m.selected.ID, "book/simple-book") {
		t.Error("should select new object")
	}
}

// =============================================================================
// E2E Test 2: N → batch creation in title panel
// =============================================================================

func TestE2E_QuickCreate_BatchFlow(t *testing.T) {
	m := setupCreateTestModel(t)
	m = initSize(t, m)

	m = pressNamedKey(t, m, "N")
	if m.create == nil || m.create.mode != createModeBatch {
		t.Fatal("should be in batch mode")
	}

	view := m.View().Content
	if !strings.Contains(view, "QUICK CREATE") {
		t.Error("help bar should show QUICK CREATE")
	}

	// Create first object
	m = typeString(t, m, "batch-one")
	m = pressKey(t, m, tea.KeyEnter)

	if m.create == nil {
		t.Fatal("should still be active")
	}
	if !strings.Contains(m.create.flash, "batch-one") {
		t.Errorf("flash = %q, should mention object name", m.create.flash)
	}

	// Flash in title panel view
	view = m.View().Content
	if !strings.Contains(view, "Created") {
		t.Error("View should show flash in title panel")
	}

	// Create second
	m = typeString(t, m, "batch-two")
	m = pressKey(t, m, tea.KeyEnter)

	lastID := m.create.lastObj.ID

	// Esc to exit
	m = pressKey(t, m, tea.KeyEscape)
	if m.create != nil {
		t.Error("should exit batch mode")
	}
	if m.selected == nil || m.selected.ID != lastID {
		t.Error("should select last object")
	}
	if m.editMode {
		t.Error("should NOT be in edit mode after batch")
	}
}

func TestE2E_QuickCreate_WithTemplates(t *testing.T) {
	m := setupCreateTestModelWithTemplates(t, map[string]string{
		"review":  "---\ntitle: Review\n---\n",
		"summary": "---\ntitle: Summary\n---\n",
	})
	m = initSize(t, m)

	m = pressNamedKey(t, m, "N")
	if m.create == nil {
		t.Fatal("should be active")
	}

	// Template defaults to first; create two objects
	m = typeString(t, m, "book-a")
	m = pressKey(t, m, tea.KeyEnter)
	m = typeString(t, m, "book-b")
	m = pressKey(t, m, tea.KeyEnter)
	m = pressKey(t, m, tea.KeyEscape)

	if m.selected == nil {
		t.Fatal("should have selected object")
	}
	obj, err := m.vault.GetObject(m.selected.ID)
	if err != nil {
		t.Fatalf("GetObject: %v", err)
	}
	if obj.Properties["title"] != "Review" {
		t.Errorf("props[title] = %v, want 'Review' from template", obj.Properties["title"])
	}
}

// =============================================================================
// E2E Test 3: n/N disabled in readonly
// =============================================================================

func TestE2E_ReadOnly_CreationDisabled(t *testing.T) {
	m := setupCreateTestModel(t)
	m = initSize(t, m)
	m.readOnly = true

	view := m.View().Content
	if !strings.Contains(view, "READONLY") {
		t.Error("should show READONLY")
	}

	m = pressNamedKey(t, m, "n")
	if m.create != nil {
		t.Error("n ignored in readonly")
	}
	m = pressNamedKey(t, m, "N")
	if m.create != nil {
		t.Error("N ignored in readonly")
	}

	entries := helpEntries(true)
	for _, e := range entries {
		if e.Key == "n" || e.Key == "N" {
			t.Errorf("help should not show %q in readonly", e.Key)
		}
	}
}

// =============================================================================
// E2E Test 4: Unique constraint inline error
// =============================================================================

func TestE2E_UniqueConstraint_InlineError(t *testing.T) {
	m := setupCreateTestModel(t)
	os.WriteFile(filepath.Join(m.vault.TypesDir(), "tag.yaml"), []byte("name: tag\nunique: true\nplural: tags\n"), 0644)
	objects, _ := m.vault.QueryObjects("")
	m.groups = buildGroups(objects, m.vault)
	m = initSize(t, m)

	tagIdx := -1
	for i, g := range m.groups {
		if g.Name == "tag" {
			tagIdx = i
			break
		}
	}
	if tagIdx < 0 {
		t.Fatal("tag not found")
	}

	rows := m.currentRows()
	for i, row := range rows {
		if row.Kind == rowHeader && row.GroupIndex == tagIdx {
			m.cursor = i
			break
		}
	}

	// Create first
	m = pressNamedKey(t, m, "n")
	m = typeString(t, m, "important")
	m = pressKey(t, m, tea.KeyEnter)

	if m.create != nil {
		if m.create.errMsg != "" {
			t.Fatalf("first creation failed: %s", m.create.errMsg)
		}
		t.Fatal("create should be nil after success")
	}

	// Exit edit, back to sidebar
	m = pressKey(t, m, tea.KeyEscape)
	m.focus = focusLeft

	rows = m.currentRows()
	for i, row := range rows {
		if row.Kind == rowHeader && row.GroupIndex == tagIdx {
			m.cursor = i
			break
		}
	}

	// Try duplicate
	m = pressNamedKey(t, m, "n")
	m = typeString(t, m, "important")
	m = pressKey(t, m, tea.KeyEnter)

	if m.create == nil || m.create.errMsg == "" {
		t.Fatal("should have duplicate error")
	}

	// Error visible in View (title panel)
	view := m.View().Content
	if !strings.Contains(view, "✗") {
		t.Error("View should show error indicator")
	}

	// Typing clears error
	m = typeString(t, m, "-v2")
	if m.create.errMsg != "" {
		t.Error("errMsg should clear after typing")
	}
}

// =============================================================================
// E2E Test 5: Tab + template cycling with live preview
// =============================================================================

func TestE2E_TabAndTemplateCycling_LivePreview(t *testing.T) {
	m := setupCreateTestModelWithTemplates(t, map[string]string{
		"review":  "---\ntitle: Review\n---\n## Review Notes\n",
		"summary": "---\ntitle: Summary\n---\n## Key Takeaways\n",
	})
	m = initSize(t, m)

	m = pressNamedKey(t, m, "n")

	// Initial preview should show first template (review)
	bodyContent := m.bodyViewport.View()
	if !strings.Contains(bodyContent, "Review Notes") {
		t.Error("initial preview should show review template body")
	}

	// Tab to template, switch to summary
	m = pressKey(t, m, tea.KeyTab)
	m = pressKey(t, m, tea.KeyRight) // review → summary

	bodyContent = m.bodyViewport.View()
	if !strings.Contains(bodyContent, "Key Takeaways") {
		t.Error("preview should update to summary template body")
	}

	// Switch to (none)
	m = pressKey(t, m, tea.KeyRight) // summary → (none)
	bodyContent = m.bodyViewport.View()
	if !strings.Contains(bodyContent, "(empty)") {
		t.Error("(none) preview should show (empty)")
	}

	// Tab back to name, type name, create
	m = pressKey(t, m, tea.KeyTab)
	m = typeString(t, m, "no-template-book")
	m = pressKey(t, m, tea.KeyEnter)

	if m.selected == nil {
		t.Fatal("should have created object")
	}
	// Created with (none) → empty body
	if m.selected.Body != "" {
		t.Errorf("body should be empty when (none) selected, got %q", m.selected.Body)
	}
}
