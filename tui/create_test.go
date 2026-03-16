package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/typemd/typemd/core"
	tea "charm.land/bubbletea/v2"
)

// setupCreateTestModel creates a model with a real vault for testing object creation.
func setupCreateTestModel(t *testing.T) model {
	t.Helper()
	dir := t.TempDir()
	v := core.NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatalf("vault Init() error = %v", err)
	}
	os.WriteFile(filepath.Join(v.TypesDir(), "book.yaml"), []byte("name: book\nproperties:\n  - name: title\n    type: string\n"), 0644)
	os.WriteFile(filepath.Join(v.TypesDir(), "note.yaml"), []byte("name: note\n"), 0644)
	if err := v.Open(); err != nil {
		t.Fatalf("vault Open() error = %v", err)
	}
	t.Cleanup(func() { v.Close() })

	objects, _ := v.QueryObjects("")
	groups := buildGroups(objects, v)
	if len(groups) > 0 {
		groups[0].Expanded = true
	}
	return model{
		vault:        v,
		focus:        focusLeft,
		groups:       groups,
		cursor:       0,
		bodyTextarea: newBodyTextarea(),
		searchInput:  initSearchInput(),
		width:        120,
		height:       24,
	}
}

// setupCreateTestModelWithTemplates creates a model with templates for the "book" type.
func setupCreateTestModelWithTemplates(t *testing.T, templates map[string]string) model {
	t.Helper()
	m := setupCreateTestModel(t)
	tmplDir := filepath.Join(m.vault.Root, "templates", "book")
	os.MkdirAll(tmplDir, 0755)
	for name, content := range templates {
		os.WriteFile(filepath.Join(tmplDir, name+".md"), []byte(content), 0644)
	}
	return m
}

// --- createState tests ---

func TestCreateState_SelectedTemplateName(t *testing.T) {
	cs := &createState{templates: []string{"review", "summary"}, cursor: 0}
	if name := cs.selectedTemplateName(); name != "review" {
		t.Errorf("selectedTemplateName() = %q, want %q", name, "review")
	}
	cs.cursor = 2
	if name := cs.selectedTemplateName(); name != "" {
		t.Errorf("selectedTemplateName() for (none) = %q, want empty", name)
	}
}

func TestCreateState_CurrentTemplateName(t *testing.T) {
	cs := &createState{templates: []string{"review"}, cursor: 0}
	if name := cs.currentTemplateName(); name != "review" {
		t.Errorf("currentTemplateName() = %q, want %q", name, "review")
	}
	cs.cursor = 1
	if name := cs.currentTemplateName(); name != noneOption {
		t.Errorf("currentTemplateName() at end = %q, want %q", name, noneOption)
	}
}

func TestCreateState_HasTemplateSelector(t *testing.T) {
	cs := &createState{templates: []string{"a", "b"}}
	if !cs.hasTemplateSelector() {
		t.Error("should have selector with 2+ templates")
	}
	cs.templates = []string{"a"}
	if cs.hasTemplateSelector() {
		t.Error("should not have interactive selector with 1 template")
	}
	cs.templates = nil
	if cs.hasTemplateSelector() {
		t.Error("should not have selector with 0 templates")
	}
}

// --- startCreate tests ---

func TestStartCreate_NoTemplates(t *testing.T) {
	m := setupCreateTestModel(t)
	m.startCreate(0, createModeSingle)
	if m.create == nil {
		t.Fatal("create should not be nil")
	}
	if m.create.field != createFieldName {
		t.Errorf("field = %d, want createFieldName", m.create.field)
	}
	if len(m.create.templates) != 0 {
		t.Errorf("templates should be empty, got %d", len(m.create.templates))
	}
}

func TestStartCreate_SingleTemplate(t *testing.T) {
	m := setupCreateTestModelWithTemplates(t, map[string]string{
		"default": "---\ntitle: default\n---\n",
	})
	m.startCreate(0, createModeSingle)
	if m.create == nil {
		t.Fatal("create should not be nil")
	}
	if m.create.selectedTemplateName() != "default" {
		t.Errorf("selectedTemplateName() = %q, want %q", m.create.selectedTemplateName(), "default")
	}
}

func TestStartCreate_MultipleTemplates(t *testing.T) {
	m := setupCreateTestModelWithTemplates(t, map[string]string{
		"review": "", "summary": "",
	})
	m.startCreate(0, createModeSingle)
	if m.create == nil {
		t.Fatal("create should not be nil")
	}
	if !m.create.hasTemplateSelector() {
		t.Error("multiple templates should be interactive")
	}
}

func TestStartCreate_ReadOnly(t *testing.T) {
	m := setupCreateTestModel(t)
	m.readOnly = true
	m.startCreate(0, createModeSingle)
	if m.create != nil {
		t.Error("create should be nil in read-only mode")
	}
}

func TestStartCreate_NameTemplatePrefill(t *testing.T) {
	m := setupCreateTestModel(t)
	os.WriteFile(filepath.Join(m.vault.TypesDir(), "journal.yaml"),
		[]byte("name: journal\nproperties:\n  - name: name\n    template: \"{{ date:YYYY-MM-DD }}\"\n"), 0644)
	objects, _ := m.vault.QueryObjects("")
	m.groups = buildGroups(objects, m.vault)

	journalIdx := -1
	for i, g := range m.groups {
		if g.Name == "journal" {
			journalIdx = i
			break
		}
	}
	if journalIdx < 0 {
		t.Fatal("journal type not found")
	}

	m.startCreate(journalIdx, createModeSingle)
	if m.create == nil {
		t.Fatal("create should not be nil — name should be pre-filled, not auto-skipped")
	}
	if m.create.nameInput.Value() == "" {
		t.Error("name should be pre-filled from name template")
	}
}

// --- Tab switching ---

func TestCreateTab_SwitchesFields(t *testing.T) {
	m := setupCreateTestModelWithTemplates(t, map[string]string{"review": "", "summary": ""})
	m.startCreate(0, createModeSingle)

	tab := tea.KeyPressMsg{Code: tea.KeyTab}
	newM, _ := m.Update(tab)
	m = newM.(model)
	if m.create.field != createFieldTemplate {
		t.Error("Tab should switch to template field")
	}

	newM, _ = m.Update(tab)
	m = newM.(model)
	if m.create.field != createFieldName {
		t.Error("Tab should switch back to name field")
	}
}

func TestCreateTab_NoOpWithoutTemplates(t *testing.T) {
	m := setupCreateTestModel(t)
	m.startCreate(0, createModeSingle)

	tab := tea.KeyPressMsg{Code: tea.KeyTab}
	newM, _ := m.Update(tab)
	m = newM.(model)
	if m.create.field != createFieldName {
		t.Error("Tab should be no-op when no templates")
	}
}

// --- Template cycling ---

func TestCreateTemplate_Cycling(t *testing.T) {
	m := setupCreateTestModelWithTemplates(t, map[string]string{"review": "", "summary": ""})
	m.startCreate(0, createModeSingle)

	tab := tea.KeyPressMsg{Code: tea.KeyTab}
	newM, _ := m.Update(tab)
	m = newM.(model)

	right := tea.KeyPressMsg{Code: tea.KeyRight}
	newM, _ = m.Update(right)
	m = newM.(model)
	if m.create.cursor != 1 {
		t.Errorf("cursor = %d, want 1", m.create.cursor)
	}

	newM, _ = m.Update(right)
	m = newM.(model)
	if m.create.cursor != 2 {
		t.Errorf("cursor = %d, want 2 (none)", m.create.cursor)
	}

	newM, _ = m.Update(right)
	m = newM.(model)
	if m.create.cursor != 0 {
		t.Errorf("cursor = %d, want 0 (wrap)", m.create.cursor)
	}

	left := tea.KeyPressMsg{Code: tea.KeyLeft}
	newM, _ = m.Update(left)
	m = newM.(model)
	if m.create.cursor != 2 {
		t.Errorf("cursor = %d, want 2 (wrap left)", m.create.cursor)
	}
}

// --- Creation ---

func TestCreate_SingleMode_CreatesAndEdits(t *testing.T) {
	m := setupCreateTestModel(t)
	m.startCreate(0, createModeSingle)

	for _, ch := range "my-book" {
		newM, _ := m.Update(tea.KeyPressMsg{Code: ch, Text: string(ch)})
		m = newM.(model)
	}

	newM, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = newM.(model)

	if m.create != nil {
		t.Error("create should be nil")
	}
	if m.selected == nil || !strings.Contains(m.selected.ID, "book/my-book") {
		t.Error("should select new object")
	}
	if m.focus != focusBody || !m.editMode {
		t.Error("should be in body edit mode")
	}
}

func TestCreate_BatchMode_CreatesAndClears(t *testing.T) {
	m := setupCreateTestModel(t)
	m.startCreate(0, createModeBatch)

	for _, ch := range "book-one" {
		newM, _ := m.Update(tea.KeyPressMsg{Code: ch, Text: string(ch)})
		m = newM.(model)
	}
	newM, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = newM.(model)

	if m.create == nil {
		t.Fatal("should still be active")
	}
	if m.create.nameInput.Value() != "" {
		t.Error("input should be cleared")
	}
	if m.create.flash == "" {
		t.Error("flash should show")
	}
}

func TestCreate_BatchMode_EscExits(t *testing.T) {
	m := setupCreateTestModel(t)
	m.startCreate(0, createModeBatch)

	for _, ch := range "book-one" {
		newM, _ := m.Update(tea.KeyPressMsg{Code: ch, Text: string(ch)})
		m = newM.(model)
	}
	newM, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = newM.(model)
	lastObj := m.create.lastObj

	newM, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	m = newM.(model)

	if m.create != nil {
		t.Error("create should be nil")
	}
	if m.selected == nil || m.selected.ID != lastObj.ID {
		t.Error("should select last object")
	}
}

func TestCreate_EmptyNameRejected(t *testing.T) {
	m := setupCreateTestModel(t)
	m.startCreate(0, createModeSingle)

	newM, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = newM.(model)
	if m.create == nil {
		t.Error("create should still be active")
	}
}

func TestCreate_UniqueError(t *testing.T) {
	m := setupCreateTestModel(t)
	os.WriteFile(filepath.Join(m.vault.TypesDir(), "tag.yaml"), []byte("name: tag\nunique: true\n"), 0644)
	objects, _ := m.vault.QueryObjects("")
	m.groups = buildGroups(objects, m.vault)

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

	m.startCreate(tagIdx, createModeSingle)
	for _, ch := range "my-tag" {
		newM, _ := m.Update(tea.KeyPressMsg{Code: ch, Text: string(ch)})
		m = newM.(model)
	}
	newM, _ := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = newM.(model)

	m.startCreate(tagIdx, createModeSingle)
	for _, ch := range "my-tag" {
		newM, _ := m.Update(tea.KeyPressMsg{Code: ch, Text: string(ch)})
		m = newM.(model)
	}
	newM, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = newM.(model)

	if m.create == nil || m.create.errMsg == "" {
		t.Error("should have error for duplicate")
	}
}

func TestCreate_ErrorClearsOnInput(t *testing.T) {
	m := setupCreateTestModel(t)
	m.create = &createState{
		mode: createModeSingle, field: createFieldName, typeName: "book",
		nameInput: initNameInput(), errMsg: "error", templateCache: make(map[string]*core.Template),
	}

	newM, _ := m.Update(tea.KeyPressMsg{Code: 'a', Text: "a"})
	m = newM.(model)
	if m.create.errMsg != "" {
		t.Error("errMsg should clear on input")
	}
}

func TestFlashDismiss(t *testing.T) {
	m := setupCreateTestModel(t)
	m.create = &createState{mode: createModeBatch, flash: "✓ test", flashSeq: 1, templateCache: make(map[string]*core.Template)}

	newM, _ := m.Update(flashDismissMsg{seq: 1})
	m = newM.(model)
	if m.create.flash != "" {
		t.Error("flash should be empty")
	}
}

func TestFlashDismiss_StaleSeqIgnored(t *testing.T) {
	m := setupCreateTestModel(t)
	m.create = &createState{mode: createModeBatch, flash: "newer", flashSeq: 2, templateCache: make(map[string]*core.Template)}

	newM, _ := m.Update(flashDismissMsg{seq: 1})
	m = newM.(model)
	if m.create.flash != "newer" {
		t.Error("stale seq should not dismiss")
	}
}

func TestCreateKeybindings_ReadOnlyIgnored(t *testing.T) {
	m := setupCreateTestModel(t)
	m.readOnly = true

	newM, _ := m.Update(tea.KeyPressMsg{Code: 'n', Text: "n"})
	m = newM.(model)
	if m.create != nil {
		t.Error("n ignored in readonly")
	}
	newM, _ = m.Update(tea.KeyPressMsg{Code: 'N', Text: "N"})
	m = newM.(model)
	if m.create != nil {
		t.Error("N ignored in readonly")
	}
}

// --- Rendering ---

func TestRenderCreateTitleContent_WithTemplates(t *testing.T) {
	cs := &createState{
		field: createFieldName, typeName: "book", emoji: "📚",
		templates: []string{"review", "summary"}, cursor: 0,
		nameInput: initNameInput(), templateCache: make(map[string]*core.Template),
	}
	result := renderCreateTitleContent(cs, 80)
	if !strings.Contains(result, "📚") || !strings.Contains(result, "book") {
		t.Error("should contain emoji and type")
	}
	if !strings.Contains(result, "📝 review") {
		t.Error("should contain template")
	}
}

func TestRenderCreateTitleContent_NoTemplates(t *testing.T) {
	cs := &createState{
		field: createFieldName, typeName: "note",
		nameInput: initNameInput(), templateCache: make(map[string]*core.Template),
	}
	result := renderCreateTitleContent(cs, 80)
	if strings.Contains(result, "📝") {
		t.Error("should not contain template icon")
	}
}

func TestRenderCreateTitleContent_TemplateFocused(t *testing.T) {
	cs := &createState{
		field: createFieldTemplate, typeName: "book",
		templates: []string{"review", "summary"}, cursor: 0,
		nameInput: initNameInput(), templateCache: make(map[string]*core.Template),
	}
	result := renderCreateTitleContent(cs, 80)
	if !strings.Contains(result, "[📝 review]") {
		t.Error("focused template should be in brackets")
	}
}

func TestRenderCreateHelpBar_NameSingle(t *testing.T) {
	cs := &createState{field: createFieldName, mode: createModeSingle, templates: []string{"a", "b"}}
	result := renderCreateHelpBar(cs)
	if !strings.Contains(result, "tab: template") {
		t.Error("should contain tab hint")
	}
	if !strings.Contains(result, "create & edit") {
		t.Error("should contain create & edit")
	}
}

func TestRenderCreateHelpBar_TemplateFocused(t *testing.T) {
	cs := &createState{field: createFieldTemplate, mode: createModeSingle, templates: []string{"a", "b"}}
	result := renderCreateHelpBar(cs)
	if !strings.Contains(result, "◀▶: switch") {
		t.Error("should contain arrow hint")
	}
}

func TestRenderCreateHelpBar_Batch(t *testing.T) {
	cs := &createState{field: createFieldName, mode: createModeBatch}
	result := renderCreateHelpBar(cs)
	if !strings.Contains(result, "QUICK CREATE") || !strings.Contains(result, "esc: done") {
		t.Error("should show batch hints")
	}
}
