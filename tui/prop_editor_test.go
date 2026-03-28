package tui

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/typemd/typemd/core"
	tea "charm.land/bubbletea/v2"
)

// testDisplayProps returns a typical set of display properties for testing.
func testDisplayProps() []core.DisplayProperty {
	return []core.DisplayProperty{
		{Key: "name", Value: "Test Object", Type: "string"},
		{Key: "description", Value: "A test object", Type: "string"},
		{Key: "title", Value: "Hello", Type: "string"},
		{Key: "rating", Value: 5, Type: "number"},
		{Key: "published", Value: false, Type: "checkbox"},
		{Key: "status", Value: "draft", Type: "select"},
		{Key: "created_at", Value: "2024-01-01", Type: "date"},
		{Key: "updated_at", Value: "2024-01-02", Type: "date"},
		{Key: "tags", Value: nil, IsRelation: true},
		{Key: "author", Value: "person/alice", IsRelation: true},
		{Key: "books", IsReverse: true, FromID: "book/clean-code"},
	}
}

func testPropSchema() *core.TypeSchema {
	return &core.TypeSchema{
		Name: "book",
		Properties: []core.Property{
			{Name: "title", Type: "string"},
			{Name: "rating", Type: "number"},
			{Name: "published", Type: "checkbox"},
			{Name: "status", Type: "select", Options: []core.Option{
				{Value: "draft"},
				{Value: "published"},
				{Value: "archived"},
			}},
			{Name: "author", Type: "relation", Target: "person"},
		},
	}
}

// --- Property cursor navigation tests ---

func TestPropEditor_CursorAppearsOnCreation(t *testing.T) {
	pe := newPropEditor(testDisplayProps(), testPropSchema())

	// Should have items (excluding name, pinned, etc.)
	if len(pe.items) == 0 {
		t.Fatal("expected non-empty items list")
	}

	// Cursor should be on first editable item
	item := pe.currentItem()
	if item == nil {
		t.Fatal("expected cursor on an item")
	}
	if !item.editable {
		t.Errorf("cursor should be on editable item, got key=%s editable=%v", item.dp.Key, item.editable)
	}
}

func TestPropEditor_CursorSkipsName(t *testing.T) {
	pe := newPropEditor(testDisplayProps(), testPropSchema())

	for _, item := range pe.items {
		if item.dp.Key == "name" {
			t.Error("name property should not appear in items list")
		}
	}
}

func TestPropEditor_CursorSkipsImmutableSystemProps(t *testing.T) {
	pe := newPropEditor(testDisplayProps(), testPropSchema())

	// created_at and updated_at should not be editable
	for _, item := range pe.items {
		if item.dp.Key == "created_at" || item.dp.Key == "updated_at" {
			if item.editable {
				t.Errorf("%s should not be editable", item.dp.Key)
			}
		}
	}
}

func TestPropEditor_ForwardRelationsAreEditable(t *testing.T) {
	pe := newPropEditor(testDisplayProps(), testPropSchema())

	for _, item := range pe.items {
		if item.dp.IsRelation && !item.dp.IsReverse && !item.dp.IsBacklink {
			if !item.editable {
				t.Errorf("forward relation property %s should be editable", item.dp.Key)
			}
		}
	}
}

func TestPropEditor_CursorSkipsReverseRelations(t *testing.T) {
	pe := newPropEditor(testDisplayProps(), testPropSchema())

	for _, item := range pe.items {
		if item.dp.IsReverse && item.editable {
			t.Errorf("reverse relation property %s should not be editable", item.dp.Key)
		}
	}
}

func TestPropEditor_TagsAreEditable(t *testing.T) {
	pe := newPropEditor(testDisplayProps(), testPropSchema())

	for _, item := range pe.items {
		if item.dp.Key == "tags" && !item.editable {
			t.Error("tags property should be editable")
		}
	}
}

func TestPropEditor_MoveDown(t *testing.T) {
	pe := newPropEditor(testDisplayProps(), testPropSchema())

	firstKey := pe.currentItem().dp.Key
	pe.moveDown()
	secondKey := pe.currentItem().dp.Key

	if firstKey == secondKey {
		t.Error("cursor should have moved to a different property")
	}
	if !pe.currentItem().editable {
		t.Error("cursor should be on an editable property after moveDown")
	}
}

func TestPropEditor_MoveUp(t *testing.T) {
	pe := newPropEditor(testDisplayProps(), testPropSchema())

	// Move down first, then up
	pe.moveDown()
	secondKey := pe.currentItem().dp.Key
	pe.moveUp()
	firstKey := pe.currentItem().dp.Key

	if firstKey == secondKey {
		t.Error("cursor should have moved back after moveUp")
	}
}

func TestPropEditor_MoveUpAtTop(t *testing.T) {
	pe := newPropEditor(testDisplayProps(), testPropSchema())

	firstKey := pe.currentItem().dp.Key
	pe.moveUp() // already at top
	sameKey := pe.currentItem().dp.Key

	if firstKey != sameKey {
		t.Error("cursor should stay at top when moveUp at boundary")
	}
}

func TestPropEditor_MoveDownAtBottom(t *testing.T) {
	pe := newPropEditor(testDisplayProps(), testPropSchema())

	// Move to bottom
	for range 20 {
		pe.moveDown()
	}
	bottomKey := pe.currentItem().dp.Key
	pe.moveDown() // already at bottom
	sameKey := pe.currentItem().dp.Key

	if bottomKey != sameKey {
		t.Error("cursor should stay at bottom when moveDown at boundary")
	}
}

// --- Textinput editing tests ---

func TestPropEditor_ActivateTextInput(t *testing.T) {
	pe := newPropEditor(testDisplayProps(), testPropSchema())

	// Move to a string property
	for pe.currentItem() != nil && pe.currentItem().dp.Key != "title" {
		pe.moveDown()
	}
	if pe.currentItem() == nil || pe.currentItem().dp.Key != "title" {
		t.Fatal("could not navigate to title property")
	}

	pe.activateEdit(nil)

	if pe.mode != propModeTextInput {
		t.Errorf("expected propModeTextInput, got %d", pe.mode)
	}
	if pe.textInput.Value() != "Hello" {
		t.Errorf("textinput should be pre-filled with 'Hello', got %q", pe.textInput.Value())
	}
}

func TestPropEditor_CancelTextInput(t *testing.T) {
	pe := newPropEditor(testDisplayProps(), testPropSchema())

	for pe.currentItem() != nil && pe.currentItem().dp.Key != "title" {
		pe.moveDown()
	}
	pe.activateEdit(nil)
	pe.cancelEdit()

	if pe.mode != propModeNavigate {
		t.Error("mode should be propModeNavigate after cancel")
	}
}

// --- Checkbox toggle tests ---

func TestPropEditor_CheckboxIsNotTextInput(t *testing.T) {
	pe := newPropEditor(testDisplayProps(), testPropSchema())

	// Move to checkbox property
	for pe.currentItem() != nil && pe.currentItem().dp.Key != "published" {
		pe.moveDown()
	}
	if pe.currentItem() == nil || pe.currentItem().dp.Key != "published" {
		t.Fatal("could not navigate to published property")
	}

	// activateEdit should return nil for checkbox (direct toggle handled elsewhere)
	cmd := pe.activateEdit(nil)
	if cmd != nil {
		t.Error("checkbox should not return a command from activateEdit")
	}
	if pe.mode != propModeNavigate {
		t.Error("checkbox should not enter textinput mode")
	}
}

// --- Select picker tests ---

func TestPropEditor_ActivateSelectPicker(t *testing.T) {
	pe := newPropEditor(testDisplayProps(), testPropSchema())

	// Move to status (select) property
	for pe.currentItem() != nil && pe.currentItem().dp.Key != "status" {
		pe.moveDown()
	}
	if pe.currentItem() == nil || pe.currentItem().dp.Key != "status" {
		t.Fatal("could not navigate to status property")
	}

	pe.activateEdit(nil)

	if pe.mode != propModeSelectPick {
		t.Errorf("expected propModeSelectPick, got %d", pe.mode)
	}
	if len(pe.pickerOptions) != 3 {
		t.Errorf("expected 3 options, got %d", len(pe.pickerOptions))
	}
	// Current value "draft" should be highlighted
	if pe.pickerCursor != 0 {
		t.Errorf("picker cursor should be on 'draft' (index 0), got %d", pe.pickerCursor)
	}
}

func TestPropEditor_SelectPickerNavigation(t *testing.T) {
	pe := newPropEditor(testDisplayProps(), testPropSchema())

	for pe.currentItem() != nil && pe.currentItem().dp.Key != "status" {
		pe.moveDown()
	}
	pe.activateEdit(nil)

	// Move down in picker
	if pe.pickerCursor != 0 {
		t.Fatalf("expected cursor at 0, got %d", pe.pickerCursor)
	}
	pe.pickerCursor++
	if pe.pickerCursor != 1 {
		t.Errorf("expected cursor at 1, got %d", pe.pickerCursor)
	}
}

func TestPropEditor_CancelSelectPicker(t *testing.T) {
	pe := newPropEditor(testDisplayProps(), testPropSchema())

	for pe.currentItem() != nil && pe.currentItem().dp.Key != "status" {
		pe.moveDown()
	}
	pe.activateEdit(nil)
	pe.cancelEdit()

	if pe.mode != propModeNavigate {
		t.Error("mode should be propModeNavigate after cancel")
	}
}

// --- Render tests ---

func TestPropEditor_RenderWithCursor(t *testing.T) {
	pe := newPropEditor(testDisplayProps(), testPropSchema())
	output := pe.Render(true) // focused

	if output == "" {
		t.Error("render output should not be empty")
	}
	// Should contain Properties header
	if !strings.Contains(output, "Properties") {
		t.Error("output should contain 'Properties' header")
	}
}

func TestPropEditor_RenderWithoutCursor(t *testing.T) {
	pe := newPropEditor(testDisplayProps(), testPropSchema())
	output := pe.Render(false) // not focused

	if output == "" {
		t.Error("render output should not be empty")
	}
}

// --- Integration tests with model ---

func TestModel_PropNavigation_UpDown(t *testing.T) {
	m := setupTestModel(t)
	m.focus = focusProps
	m.propsVisible = true
	// Manually initialize propEdit since no vault
	m.propEdit = newPropEditor(testDisplayProps(), testPropSchema())

	startCursor := m.propEdit.cursor

	// Press down
	msg := tea.KeyPressMsg{Code: 'j', Text: "j"}
	newM, _ := m.Update(msg)
	updated := newM.(model)

	if updated.propEdit.cursor == startCursor && len(updated.propEdit.items) > 1 {
		hasMoreEditable := false
		for i := startCursor + 1; i < len(updated.propEdit.items); i++ {
			if updated.propEdit.items[i].editable {
				hasMoreEditable = true
				break
			}
		}
		if hasMoreEditable {
			t.Error("cursor should have moved down")
		}
	}
}

func TestModel_TabDoesNotEnterPropsForLockedObject(t *testing.T) {
	m := setupTestModel(t)
	m.focus = focusBody
	m.propsVisible = true
	m.rightPanel = panelObject
	// Set selected object as locked
	m.selected = &core.Object{
		ID:         "book/test-locked",
		Type:       "book",
		Filename:   "test-locked",
		Properties: map[string]any{core.LockedProperty: true},
	}
	m.propEdit = newPropEditor(testDisplayProps(), testPropSchema())

	msg := tea.KeyPressMsg{Code: tea.KeyTab}
	newM, _ := m.Update(msg)
	updated := newM.(model)

	if updated.focus == focusProps {
		t.Error("focus should NOT be focusProps for locked object after Tab")
	}
	if updated.focus != focusBody {
		t.Errorf("focus should remain focusBody for locked object, got %d", updated.focus)
	}
}

func TestModel_TabEntersPropsForUnlockedObject(t *testing.T) {
	m := setupTestModel(t)
	m.focus = focusBody
	m.propsVisible = true
	m.rightPanel = panelObject
	m.selected = &core.Object{
		ID:         "book/test-unlocked",
		Type:       "book",
		Filename:   "test-unlocked",
		Properties: map[string]any{},
	}
	m.propEdit = newPropEditor(testDisplayProps(), testPropSchema())

	msg := tea.KeyPressMsg{Code: tea.KeyTab}
	newM, _ := m.Update(msg)
	updated := newM.(model)

	if updated.focus != focusProps {
		t.Errorf("focus should be focusProps for unlocked object after Tab, got %d", updated.focus)
	}
}

func TestModel_PropEscReturnsToSidebar(t *testing.T) {
	m := setupTestModel(t)
	m.focus = focusProps
	m.propsVisible = true
	m.propEdit = newPropEditor(testDisplayProps(), testPropSchema())

	msg := tea.KeyPressMsg{Code: tea.KeyEscape, Text: "esc"}
	newM, _ := m.Update(msg)
	updated := newM.(model)

	if updated.focus != focusLeft {
		t.Errorf("focus should be focusLeft after Esc, got %d", updated.focus)
	}
}

// --- Date picker integration tests ---

func testDateDisplayProps() []core.DisplayProperty {
	return []core.DisplayProperty{
		{Key: "name", Value: "Test Object", Type: "string"},
		{Key: "published_date", Value: time.Date(2025, 6, 15, 0, 0, 0, 0, time.Local), Type: "date"},
		{Key: "title", Value: "Hello", Type: "string"},
	}
}

func testDatePropSchema() *core.TypeSchema {
	return &core.TypeSchema{
		Name: "event",
		Properties: []core.Property{
			{Name: "published_date", Type: "date"},
			{Name: "title", Type: "string"},
		},
	}
}

func TestPropEditor_ActivateDatePicker(t *testing.T) {
	pe := newPropEditor(testDateDisplayProps(), testDatePropSchema())

	// Navigate to the date property
	for pe.currentItem() != nil && pe.currentItem().dp.Key != "published_date" {
		pe.moveDown()
	}
	if pe.currentItem() == nil || pe.currentItem().dp.Key != "published_date" {
		t.Fatal("could not navigate to published_date property")
	}

	pe.activateEdit(nil)

	if pe.mode != propModeDateSegment {
		t.Errorf("expected propModeDateSegment, got %d", pe.mode)
	}
	if pe.dateEdit == nil {
		t.Fatal("dateEdit should not be nil")
	}
	if pe.dateEdit.Value() != "2025-06-15" {
		t.Errorf("expected pre-fill 2025-06-15, got %s", pe.dateEdit.Value())
	}
}

func TestPropEditor_DatePickerCancelRestores(t *testing.T) {
	pe := newPropEditor(testDateDisplayProps(), testDatePropSchema())

	for pe.currentItem() != nil && pe.currentItem().dp.Key != "published_date" {
		pe.moveDown()
	}
	pe.activateEdit(nil)

	// Modify the date
	pe.dateEdit.incrementSegment(1) // year + 1

	// Cancel
	pe.cancelEdit()

	if pe.mode != propModeNavigate {
		t.Error("mode should be propModeNavigate after cancel")
	}
	if pe.dateEdit != nil {
		t.Error("dateEdit should be nil after cancel")
	}
}

func TestPropEditor_DatePickerEmptyPreFillsToday(t *testing.T) {
	props := []core.DisplayProperty{
		{Key: "name", Value: "Test", Type: "string"},
		{Key: "due_date", Value: nil, Type: "date"},
	}
	schema := &core.TypeSchema{
		Name: "task",
		Properties: []core.Property{
			{Name: "due_date", Type: "date"},
		},
	}
	pe := newPropEditor(props, schema)

	for pe.currentItem() != nil && pe.currentItem().dp.Key != "due_date" {
		pe.moveDown()
	}
	pe.activateEdit(nil)

	if pe.dateEdit == nil {
		t.Fatal("dateEdit should not be nil")
	}

	today := time.Now()
	val := pe.dateEdit.Value()
	expected := fmt.Sprintf("%04d-%02d-%02d", today.Year(), today.Month(), today.Day())
	if val != expected {
		t.Errorf("expected today's date %s for empty value, got %s", expected, val)
	}
}

func TestPropEditor_DatePickerModeToggle(t *testing.T) {
	pe := newPropEditor(testDateDisplayProps(), testDatePropSchema())

	for pe.currentItem() != nil && pe.currentItem().dp.Key != "published_date" {
		pe.moveDown()
	}
	pe.activateEdit(nil)

	if pe.mode != propModeDateSegment {
		t.Error("expected segment mode initially")
	}

	// Toggle to calendar
	pe.dateEdit.Update(tea.KeyPressMsg{Code: 'c', Text: "c"})
	// Sync mode
	if pe.dateEdit.Mode() == dateCalendarMode {
		pe.mode = propModeDateCalendar
	}
	if pe.mode != propModeDateCalendar {
		t.Error("expected calendar mode after toggle")
	}
}

func TestPropEditor_DatePickerRender(t *testing.T) {
	pe := newPropEditor(testDateDisplayProps(), testDatePropSchema())

	for pe.currentItem() != nil && pe.currentItem().dp.Key != "published_date" {
		pe.moveDown()
	}
	pe.activateEdit(nil)

	output := pe.Render(true)
	if !strings.Contains(output, "published_date") {
		t.Error("output should contain the property name")
	}
	// Should contain day of week
	if !strings.Contains(output, "Sun") {
		t.Errorf("output should contain day of week, got: %s", output)
	}
}

