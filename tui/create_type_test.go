package tui

import (
	"strings"
	"testing"

	"github.com/typemd/typemd/core"
	tea "charm.land/bubbletea/v2"
)

// --- createTypeState tests ---

func TestCreateTypeField_NextCycles(t *testing.T) {
	cts := &createTypeState{
		field:       createTypeFieldName,
		emojiInput:  initCreateTypeInput("emoji", 4),
		nameInput:   initCreateTypeInput("type name", 50),
		pluralInput: initCreateTypeInput("plural", 50),
	}

	cts.nextCreateTypeField()
	if cts.field != createTypeFieldPlural {
		t.Errorf("from name, next = %d, want createTypeFieldPlural", cts.field)
	}

	cts.nextCreateTypeField()
	if cts.field != createTypeFieldEmoji {
		t.Errorf("from plural, next = %d, want createTypeFieldEmoji", cts.field)
	}

	cts.nextCreateTypeField()
	if cts.field != createTypeFieldName {
		t.Errorf("from emoji, next = %d, want createTypeFieldName", cts.field)
	}
}

// --- startCreateType tests ---

func TestStartCreateType_InitializesState(t *testing.T) {
	m := setupCreateTestModel(t)
	m.startCreateType()

	if m.createType == nil {
		t.Fatal("createType should not be nil")
	}
	if m.createType.field != createTypeFieldName {
		t.Errorf("field = %d, want createTypeFieldName", m.createType.field)
	}
}

func TestStartCreateType_ReadOnlyIgnored(t *testing.T) {
	m := setupCreateTestModel(t)
	m.readOnly = true
	m.startCreateType()

	if m.createType != nil {
		t.Error("createType should be nil in read-only mode")
	}
}

// --- Tab switching ---

func TestCreateTypeTab_CyclesFields(t *testing.T) {
	m := setupCreateTestModel(t)
	m.startCreateType()

	tab := tea.KeyPressMsg{Code: tea.KeyTab}

	// Name → Plural
	newM, _ := m.Update(tab)
	m = newM.(model)
	if m.createType.field != createTypeFieldPlural {
		t.Errorf("Tab from name: field = %d, want plural", m.createType.field)
	}

	// Plural → Emoji
	newM, _ = m.Update(tab)
	m = newM.(model)
	if m.createType.field != createTypeFieldEmoji {
		t.Errorf("Tab from plural: field = %d, want emoji", m.createType.field)
	}

	// Emoji → Name
	newM, _ = m.Update(tab)
	m = newM.(model)
	if m.createType.field != createTypeFieldName {
		t.Errorf("Tab from emoji: field = %d, want name", m.createType.field)
	}
}

// --- Creation ---

func TestCreateType_ValidName_CreatesAndOpensEditor(t *testing.T) {
	m := setupCreateTestModel(t)
	m.startCreateType()

	// Type "meeting"
	for _, ch := range "meeting" {
		newM, _ := m.Update(tea.KeyPressMsg{Code: ch, Text: string(ch)})
		m = newM.(model)
	}

	newM, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = newM.(model)

	if m.createType != nil {
		t.Error("createType should be nil after creation")
	}
	if m.rightPanel != panelTypeEditor {
		t.Error("should open type editor")
	}
	if m.typeEditor == nil {
		t.Fatal("typeEditor should not be nil")
	}
	if m.typeEditor.typeName != "meeting" {
		t.Errorf("typeName = %q, want %q", m.typeEditor.typeName, "meeting")
	}
	if m.focus != focusBody {
		t.Error("focus should be on body (type editor)")
	}
}

func TestCreateType_WithEmoji_SetsEmoji(t *testing.T) {
	m := setupCreateTestModel(t)
	m.startCreateType()

	// Switch to emoji field
	tab := tea.KeyPressMsg{Code: tea.KeyTab}
	newM, _ := m.Update(tab) // name → plural
	m = newM.(model)
	newM, _ = m.Update(tab) // plural → emoji
	m = newM.(model)

	// Type emoji
	newM, _ = m.Update(tea.KeyPressMsg{Code: '📋', Text: "📋"})
	m = newM.(model)

	// Switch to name field
	newM, _ = m.Update(tab)
	m = newM.(model)

	// Type name
	for _, ch := range "meeting" {
		newM, _ = m.Update(tea.KeyPressMsg{Code: ch, Text: string(ch)})
		m = newM.(model)
	}

	newM, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = newM.(model)

	if m.createType != nil {
		t.Error("createType should be nil after creation")
	}
	if m.typeEditor == nil {
		t.Fatal("typeEditor should not be nil")
	}
	if m.typeEditor.schema.Emoji != "📋" {
		t.Errorf("emoji = %q, want %q", m.typeEditor.schema.Emoji, "📋")
	}
}

func TestCreateType_WithPlural_SetsPlural(t *testing.T) {
	m := setupCreateTestModel(t)
	m.startCreateType()

	// Type name
	for _, ch := range "meeting" {
		newM, _ := m.Update(tea.KeyPressMsg{Code: ch, Text: string(ch)})
		m = newM.(model)
	}

	// Switch to plural field
	tab := tea.KeyPressMsg{Code: tea.KeyTab}
	newM, _ := m.Update(tab)
	m = newM.(model)

	// Type plural
	for _, ch := range "meetings" {
		newM, _ = m.Update(tea.KeyPressMsg{Code: ch, Text: string(ch)})
		m = newM.(model)
	}

	newM, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = newM.(model)

	if m.typeEditor == nil {
		t.Fatal("typeEditor should not be nil")
	}
	if m.typeEditor.schema.Plural != "meetings" {
		t.Errorf("plural = %q, want %q", m.typeEditor.schema.Plural, "meetings")
	}
}

func TestCreateType_EmptyName_Rejected(t *testing.T) {
	m := setupCreateTestModel(t)
	m.startCreateType()

	newM, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = newM.(model)

	if m.createType == nil {
		t.Error("createType should still be active")
	}
}

func TestCreateType_DuplicateName_ShowsError(t *testing.T) {
	m := setupCreateTestModel(t)
	m.startCreateType()

	// "book" already exists
	for _, ch := range "book" {
		newM, _ := m.Update(tea.KeyPressMsg{Code: ch, Text: string(ch)})
		m = newM.(model)
	}

	newM, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = newM.(model)

	if m.createType == nil {
		t.Fatal("createType should still be active")
	}
	if m.createType.errMsg == "" {
		t.Error("should have error for duplicate name")
	}
	if !strings.Contains(m.createType.errMsg, "already exists") {
		t.Errorf("errMsg = %q, want 'already exists'", m.createType.errMsg)
	}
}

func TestCreateType_ErrorClearsOnInput(t *testing.T) {
	m := setupCreateTestModel(t)
	m.createType = &createTypeState{
		field:       createTypeFieldName,
		emojiInput:  initCreateTypeInput("emoji", 4),
		nameInput:   initCreateTypeInput("type name", 50),
		pluralInput: initCreateTypeInput("plural", 50),
		errMsg:      "some error",
	}
	m.createType.nameInput.Focus()

	newM, _ := m.Update(tea.KeyPressMsg{Code: 'x', Text: "x"})
	m = newM.(model)

	if m.createType.errMsg != "" {
		t.Error("errMsg should be cleared on input")
	}
}

func TestCreateType_Escape_Cancels(t *testing.T) {
	m := setupCreateTestModel(t)
	m.startCreateType()

	newM, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	m = newM.(model)

	if m.createType != nil {
		t.Error("createType should be nil after Escape")
	}
}

// --- Rendering ---

func TestRenderCreateTypeTitleContent_WithEmoji(t *testing.T) {
	cts := &createTypeState{
		field:       createTypeFieldName,
		emojiInput:  initCreateTypeInput("emoji", 4),
		nameInput:   initCreateTypeInput("type name", 50),
		pluralInput: initCreateTypeInput("plural", 50),
	}
	cts.emojiInput.SetValue("📝")

	result := renderCreateTypeTitleContent(cts)
	if !strings.Contains(result, "📝") {
		t.Error("should contain emoji")
	}
	if !strings.Contains(result, "new type") {
		t.Error("should contain 'new type'")
	}
	if !strings.Contains(result, "plural:") {
		t.Error("should contain plural field")
	}
}

func TestRenderCreateTypeTitleContent_EmojiFieldFocused(t *testing.T) {
	cts := &createTypeState{
		field:       createTypeFieldEmoji,
		emojiInput:  initCreateTypeInput("emoji", 4),
		nameInput:   initCreateTypeInput("type name", 50),
		pluralInput: initCreateTypeInput("plural", 50),
	}
	cts.emojiInput.Focus()

	result := renderCreateTypeTitleContent(cts)
	if !strings.Contains(result, "[") {
		t.Error("focused emoji field should be in brackets")
	}
}

func TestRenderCreateTypeTitleContent_WithError(t *testing.T) {
	cts := &createTypeState{
		field:       createTypeFieldName,
		emojiInput:  initCreateTypeInput("emoji", 4),
		nameInput:   initCreateTypeInput("type name", 50),
		pluralInput: initCreateTypeInput("plural", 50),
		errMsg:      "type \"book\" already exists",
	}

	result := renderCreateTypeTitleContent(cts)
	if !strings.Contains(result, "✗") {
		t.Error("should contain error marker")
	}
	if !strings.Contains(result, "already exists") {
		t.Error("should contain error message")
	}
}

func TestRenderCreateTypePreview_ShowsFields(t *testing.T) {
	cts := &createTypeState{
		field:       createTypeFieldName,
		emojiInput:  initCreateTypeInput("emoji", 4),
		nameInput:   initCreateTypeInput("type name", 50),
		pluralInput: initCreateTypeInput("plural", 50),
	}
	cts.nameInput.SetValue("meeting")
	cts.pluralInput.SetValue("meetings")
	cts.emojiInput.SetValue("📋")

	result := renderCreateTypePreview(cts)
	if !strings.Contains(result, "meeting") {
		t.Error("should contain name")
	}
	if !strings.Contains(result, "meetings") {
		t.Error("should contain plural")
	}
	if !strings.Contains(result, "📋") {
		t.Error("should contain emoji")
	}
	if !strings.Contains(result, "Properties: (none)") {
		t.Error("should show empty properties")
	}
}

func TestRenderCreateTypePreview_EmptyFields(t *testing.T) {
	cts := &createTypeState{
		field:       createTypeFieldName,
		emojiInput:  initCreateTypeInput("emoji", 4),
		nameInput:   initCreateTypeInput("type name", 50),
		pluralInput: initCreateTypeInput("plural", 50),
	}

	result := renderCreateTypePreview(cts)
	if !strings.Contains(result, "(not set)") {
		t.Error("empty fields should show (not set)")
	}
}

func TestRenderCreateTypeHelpBar(t *testing.T) {
	result := renderCreateTypeHelpBar()
	if !strings.Contains(result, "NEW TYPE") {
		t.Error("should contain NEW TYPE")
	}
	if !strings.Contains(result, "tab: next field") {
		t.Error("should contain tab hint")
	}
	if !strings.Contains(result, "enter: create") {
		t.Error("should contain enter hint")
	}
	if !strings.Contains(result, "esc: cancel") {
		t.Error("should contain esc hint")
	}
}

// --- hasTitlePanel ---

func TestHasTitlePanel_CreateType(t *testing.T) {
	m := setupCreateTestModel(t)
	if m.hasTitlePanel() {
		t.Error("should be false initially")
	}

	m.startCreateType()
	if !m.hasTitlePanel() {
		t.Error("should be true when createType is active")
	}
}

// --- Cursor position after creation ---

func TestCreateType_CursorMovesToNewType(t *testing.T) {
	m := setupCreateTestModel(t)
	m.startCreateType()

	for _, ch := range "zebra" {
		newM, _ := m.Update(tea.KeyPressMsg{Code: ch, Text: string(ch)})
		m = newM.(model)
	}

	newM, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = newM.(model)

	// Cursor should be on the "zebra" type header (last alphabetically)
	rows := m.currentRows()
	if m.cursor < 0 || m.cursor >= len(rows) {
		t.Fatalf("cursor %d out of range", m.cursor)
	}
	row := rows[m.cursor]
	if row.Kind != rowHeader {
		t.Errorf("cursor should be on header row, got kind %d", row.Kind)
	}
	if m.groups[row.GroupIndex].Name != "zebra" {
		t.Errorf("cursor on %q, want %q", m.groups[row.GroupIndex].Name, "zebra")
	}
}

// --- View rendering priority ---

func TestCreateType_OverridesTypeEditor(t *testing.T) {
	m := setupCreateTestModel(t)
	m.width = 120
	m.height = 24

	// Simulate being on a type header with type editor open
	m.rightPanel = panelTypeEditor
	ts := &core.TypeSchema{Name: "book", Emoji: "📚"}
	m.typeEditor = newTypeEditor(ts, "book", false, m.vault)

	// Start type creation — should override type editor in View()
	m.startCreateType()

	v := m.View()
	if !strings.Contains(v.Content, "new type") {
		t.Error("View should show type creation form, not type editor")
	}
	if !strings.Contains(v.Content, "NEW TYPE") {
		t.Error("help bar should show NEW TYPE mode")
	}
}
