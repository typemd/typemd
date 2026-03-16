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
	// Expand first group
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

	// Create templates directory
	tmplDir := filepath.Join(m.vault.Root, "templates", "book")
	os.MkdirAll(tmplDir, 0755)
	for name, content := range templates {
		os.WriteFile(filepath.Join(tmplDir, name+".md"), []byte(content), 0644)
	}

	return m
}

// --- Task 1.1: createState struct tests ---

func TestCreateState_TemplateOptions(t *testing.T) {
	cs := &createState{
		templates: []string{"review", "summary"},
	}
	opts := cs.templateOptions()
	if len(opts) != 3 {
		t.Fatalf("len(opts) = %d, want 3", len(opts))
	}
	if opts[0] != "review" {
		t.Errorf("opts[0] = %q, want %q", opts[0], "review")
	}
	if opts[2] != noneOption {
		t.Errorf("opts[2] = %q, want %q", opts[2], noneOption)
	}
}

func TestCreateState_SelectedTemplateName(t *testing.T) {
	cs := &createState{
		templates: []string{"review", "summary"},
		cursor:    0,
	}
	if name := cs.selectedTemplateName(); name != "review" {
		t.Errorf("selectedTemplateName() = %q, want %q", name, "review")
	}
	cs.cursor = 2 // "(none)"
	if name := cs.selectedTemplateName(); name != "" {
		t.Errorf("selectedTemplateName() for (none) = %q, want empty", name)
	}
}

func TestCreateState_StepTransitions(t *testing.T) {
	cs := &createState{
		mode:     createModeSingle,
		step:     createStepTemplate,
		typeName: "book",
	}
	// Template → Name
	cs.step = createStepName
	if cs.step != createStepName {
		t.Errorf("step = %d, want createStepName", cs.step)
	}
}

func TestCreateState_ModeValues(t *testing.T) {
	if createModeSingle != 1 {
		t.Errorf("createModeSingle = %d, want 1", createModeSingle)
	}
	if createModeBatch != 2 {
		t.Errorf("createModeBatch = %d, want 2", createModeBatch)
	}
}

// --- Task 1.4: startCreate step determination ---

func TestStartCreate_NoTemplates_GoesToNameStep(t *testing.T) {
	m := setupCreateTestModel(t)
	// "book" type has no templates
	m.startCreate(0, createModeSingle)
	if m.create == nil {
		t.Fatal("create should not be nil")
	}
	if m.create.step != createStepName {
		t.Errorf("step = %d, want createStepName", m.create.step)
	}
}

func TestStartCreate_SingleTemplate_AutoSelects(t *testing.T) {
	m := setupCreateTestModelWithTemplates(t, map[string]string{
		"default": "---\ntitle: default\n---\n",
	})
	m.startCreate(0, createModeSingle)
	if m.create == nil {
		t.Fatal("create should not be nil")
	}
	if m.create.template != "default" {
		t.Errorf("template = %q, want %q", m.create.template, "default")
	}
	if m.create.step != createStepName {
		t.Errorf("step = %d, want createStepName", m.create.step)
	}
}

func TestStartCreate_MultipleTemplates_GoesToTemplateStep(t *testing.T) {
	m := setupCreateTestModelWithTemplates(t, map[string]string{
		"review":  "---\ntitle: review\n---\n",
		"summary": "---\ntitle: summary\n---\n",
	})
	m.startCreate(0, createModeSingle)
	if m.create == nil {
		t.Fatal("create should not be nil")
	}
	if m.create.step != createStepTemplate {
		t.Errorf("step = %d, want createStepTemplate", m.create.step)
	}
	if len(m.create.templates) != 2 {
		t.Errorf("len(templates) = %d, want 2", len(m.create.templates))
	}
}

func TestStartCreate_ReadOnly_NoCreate(t *testing.T) {
	m := setupCreateTestModel(t)
	m.readOnly = true
	m.startCreate(0, createModeSingle)
	if m.create != nil {
		t.Error("create should be nil in read-only mode")
	}
}

// --- Task 2.1: Template selection key handling ---

func TestCreateTemplate_CursorNavigation(t *testing.T) {
	m := setupCreateTestModelWithTemplates(t, map[string]string{
		"review":  "",
		"summary": "",
	})
	m.startCreate(0, createModeSingle)
	if m.create == nil || m.create.step != createStepTemplate {
		t.Fatal("expected template step")
	}

	// Down moves cursor
	msg := tea.KeyPressMsg{Code: tea.KeyDown}
	newM, _ := m.Update(msg)
	m = newM.(model)
	if m.create.cursor != 1 {
		t.Errorf("cursor after down = %d, want 1", m.create.cursor)
	}

	// Up wraps to end
	msg = tea.KeyPressMsg{Code: tea.KeyUp}
	newM, _ = m.Update(msg)
	m = newM.(model)
	if m.create.cursor != 0 {
		t.Errorf("cursor after up = %d, want 0", m.create.cursor)
	}

	// Up from 0 wraps to last
	msg = tea.KeyPressMsg{Code: tea.KeyUp}
	newM, _ = m.Update(msg)
	m = newM.(model)
	opts := m.create.templateOptions()
	if m.create.cursor != len(opts)-1 {
		t.Errorf("cursor after wrap up = %d, want %d", m.create.cursor, len(opts)-1)
	}
}

func TestCreateTemplate_EnterConfirms(t *testing.T) {
	m := setupCreateTestModelWithTemplates(t, map[string]string{
		"review":  "",
		"summary": "",
	})
	m.startCreate(0, createModeSingle)

	msg := tea.KeyPressMsg{Code: tea.KeyEnter}
	newM, _ := m.Update(msg)
	m = newM.(model)

	if m.create == nil {
		t.Fatal("create should not be nil after template confirm")
	}
	if m.create.step != createStepName {
		t.Errorf("step = %d, want createStepName", m.create.step)
	}
}

func TestCreateTemplate_EscCancels(t *testing.T) {
	m := setupCreateTestModelWithTemplates(t, map[string]string{
		"review":  "",
		"summary": "",
	})
	m.startCreate(0, createModeSingle)

	msg := tea.KeyPressMsg{Code: tea.KeyEscape}
	newM, _ := m.Update(msg)
	m = newM.(model)

	if m.create != nil {
		t.Error("create should be nil after esc")
	}
}

// --- Task 2.4: None selection ---

func TestCreateTemplate_NoneSelection(t *testing.T) {
	m := setupCreateTestModelWithTemplates(t, map[string]string{
		"review":  "",
		"summary": "",
	})
	m.startCreate(0, createModeSingle)

	// Move to last option (none)
	opts := m.create.templateOptions()
	for i := 0; i < len(opts)-1; i++ {
		msg := tea.KeyPressMsg{Code: tea.KeyDown}
		newM, _ := m.Update(msg)
		m = newM.(model)
	}

	msg := tea.KeyPressMsg{Code: tea.KeyEnter}
	newM, _ := m.Update(msg)
	m = newM.(model)

	if m.create == nil {
		t.Fatal("create should not be nil")
	}
	if m.create.template != "" {
		t.Errorf("template = %q, want empty (none selected)", m.create.template)
	}
	if m.create.step != createStepName {
		t.Errorf("step = %d, want createStepName", m.create.step)
	}
}

// --- Task 3.1: Name input in Create & Edit mode ---

func TestCreateName_SingleMode_EnterCreatesAndEdits(t *testing.T) {
	m := setupCreateTestModel(t)
	m.startCreate(0, createModeSingle)

	// Type a name
	for _, ch := range "my-book" {
		msg := tea.KeyPressMsg{Code: ch, Text: string(ch)}
		newM, _ := m.Update(msg)
		m = newM.(model)
	}

	// Enter to create
	msg := tea.KeyPressMsg{Code: tea.KeyEnter}
	newM, _ := m.Update(msg)
	m = newM.(model)

	if m.create != nil {
		t.Error("create should be nil after single mode creation")
	}
	if m.selected == nil {
		t.Fatal("selected should not be nil after creation")
	}
	if !strings.Contains(m.selected.ID, "book/my-book") {
		t.Errorf("selected.ID = %q, want to contain 'book/my-book'", m.selected.ID)
	}
	if m.focus != focusBody {
		t.Errorf("focus = %d, want focusBody(%d)", m.focus, focusBody)
	}
	if !m.editMode {
		t.Error("editMode should be true after single mode creation")
	}
}

func TestCreateName_SingleMode_EscCancels(t *testing.T) {
	m := setupCreateTestModel(t)
	m.startCreate(0, createModeSingle)

	msg := tea.KeyPressMsg{Code: tea.KeyEscape}
	newM, _ := m.Update(msg)
	m = newM.(model)

	if m.create != nil {
		t.Error("create should be nil after esc")
	}
}

// --- Task 3.2: Name input in Quick Create mode ---

func TestCreateName_BatchMode_EnterCreatesAndClears(t *testing.T) {
	m := setupCreateTestModel(t)
	m.startCreate(0, createModeBatch)

	// Type and create first object
	for _, ch := range "book-one" {
		msg := tea.KeyPressMsg{Code: ch, Text: string(ch)}
		newM, _ := m.Update(msg)
		m = newM.(model)
	}
	msg := tea.KeyPressMsg{Code: tea.KeyEnter}
	newM, _ := m.Update(msg)
	m = newM.(model)

	if m.create == nil {
		t.Fatal("create should still be active in batch mode")
	}
	if m.create.nameInput.Value() != "" {
		t.Errorf("input should be cleared, got %q", m.create.nameInput.Value())
	}
	if m.create.flash == "" {
		t.Error("flash should show success message")
	}
	if m.create.lastObj == nil {
		t.Error("lastObj should be set")
	}
}

func TestCreateName_BatchMode_EscExits(t *testing.T) {
	m := setupCreateTestModel(t)
	m.startCreate(0, createModeBatch)

	// Create one object first
	for _, ch := range "book-one" {
		msg := tea.KeyPressMsg{Code: ch, Text: string(ch)}
		newM, _ := m.Update(msg)
		m = newM.(model)
	}
	msg := tea.KeyPressMsg{Code: tea.KeyEnter}
	newM, _ := m.Update(msg)
	m = newM.(model)

	lastObj := m.create.lastObj

	// Esc to exit
	msg = tea.KeyPressMsg{Code: tea.KeyEscape}
	newM, _ = m.Update(msg)
	m = newM.(model)

	if m.create != nil {
		t.Error("create should be nil after esc")
	}
	if m.selected == nil || m.selected.ID != lastObj.ID {
		t.Error("should select last created object after batch exit")
	}
}

// --- Task 4.1: Name template auto-skip ---

func TestCreateName_SingleMode_NameTemplateAutoSkip(t *testing.T) {
	m := setupCreateTestModel(t)

	// Write a type with name template
	os.WriteFile(
		filepath.Join(m.vault.TypesDir(), "journal.yaml"),
		[]byte("name: journal\nproperties:\n  - name: name\n    template: \"{{ date:YYYY-MM-DD }}\"\n"),
		0644,
	)
	// Rebuild groups to include journal
	objects, _ := m.vault.QueryObjects("")
	m.groups = buildGroups(objects, m.vault)

	// Find journal group index
	journalIdx := -1
	for i, g := range m.groups {
		if g.Name == "journal" {
			journalIdx = i
			break
		}
	}
	if journalIdx < 0 {
		t.Fatal("journal type not found in groups")
	}

	m.startCreate(journalIdx, createModeSingle)

	// Should auto-create (no create state left)
	if m.create != nil {
		t.Error("create should be nil after name template auto-create")
	}
	if m.selected == nil {
		t.Fatal("selected should not be nil after auto-create")
	}
	if !strings.Contains(m.selected.ID, "journal/") {
		t.Errorf("selected.ID = %q, want to contain 'journal/'", m.selected.ID)
	}
	if !m.editMode {
		t.Error("editMode should be true after single mode auto-create")
	}
}

// --- Task 4.3: Quick Create ignores name template ---

func TestCreateName_BatchMode_IgnoresNameTemplate(t *testing.T) {
	m := setupCreateTestModel(t)

	os.WriteFile(
		filepath.Join(m.vault.TypesDir(), "journal.yaml"),
		[]byte("name: journal\nproperties:\n  - name: name\n    template: \"{{ date:YYYY-MM-DD }}\"\n"),
		0644,
	)
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

	m.startCreate(journalIdx, createModeBatch)

	// Should stay in create mode (not auto-create)
	if m.create == nil {
		t.Error("create should not be nil — batch mode should require name input")
	}
	if m.create.step != createStepName {
		t.Errorf("step = %d, want createStepName", m.create.step)
	}
}

// --- Task 5.1: Post-creation in Create & Edit mode ---

func TestCreateSingle_PostCreation_FocusBodyEditMode(t *testing.T) {
	m := setupCreateTestModel(t)
	m.startCreate(0, createModeSingle)

	for _, ch := range "test-obj" {
		msg := tea.KeyPressMsg{Code: ch, Text: string(ch)}
		newM, _ := m.Update(msg)
		m = newM.(model)
	}

	msg := tea.KeyPressMsg{Code: tea.KeyEnter}
	newM, _ := m.Update(msg)
	m = newM.(model)

	if m.focus != focusBody {
		t.Errorf("focus = %d, want focusBody", m.focus)
	}
	if !m.editMode {
		t.Error("editMode should be true")
	}
	if m.rightPanel != panelObject {
		t.Errorf("rightPanel = %d, want panelObject", m.rightPanel)
	}
}

// --- Task 5.3: Post-creation in Quick Create mode ---

func TestCreateBatch_PostCreation_FlashAndClear(t *testing.T) {
	m := setupCreateTestModel(t)
	m.startCreate(0, createModeBatch)

	for _, ch := range "batch-obj" {
		msg := tea.KeyPressMsg{Code: ch, Text: string(ch)}
		newM, _ := m.Update(msg)
		m = newM.(model)
	}

	msg := tea.KeyPressMsg{Code: tea.KeyEnter}
	newM, _ := m.Update(msg)
	m = newM.(model)

	if m.create == nil {
		t.Fatal("create should still be active")
	}
	if !strings.Contains(m.create.flash, "Created") {
		t.Errorf("flash = %q, want to contain 'Created'", m.create.flash)
	}
	if m.create.nameInput.Value() != "" {
		t.Error("input should be cleared after creation")
	}
}

// --- Task 5.5: Flash dismiss ---

func TestFlashDismiss(t *testing.T) {
	m := setupCreateTestModel(t)
	m.create = &createState{
		mode:     createModeBatch,
		step:     createStepName,
		flash:    "✓ Created: test",
		flashSeq: 1,
	}

	// Matching seq dismisses the flash
	newM, _ := m.Update(flashDismissMsg{seq: 1})
	m = newM.(model)

	if m.create.flash != "" {
		t.Errorf("flash should be empty after dismiss, got %q", m.create.flash)
	}
}

func TestFlashDismiss_StaleSeqIgnored(t *testing.T) {
	m := setupCreateTestModel(t)
	m.create = &createState{
		mode:     createModeBatch,
		step:     createStepName,
		flash:    "✓ Created: newer",
		flashSeq: 2,
	}

	// Stale seq (from earlier creation) should not dismiss current flash
	newM, _ := m.Update(flashDismissMsg{seq: 1})
	m = newM.(model)

	if m.create.flash != "✓ Created: newer" {
		t.Errorf("flash should not be dismissed by stale seq, got %q", m.create.flash)
	}
}

// --- Task 6.1: Unique constraint error ---

func TestCreateName_UniqueError(t *testing.T) {
	m := setupCreateTestModel(t)

	// Write unique type
	os.WriteFile(
		filepath.Join(m.vault.TypesDir(), "tag.yaml"),
		[]byte("name: tag\nunique: true\n"),
		0644,
	)
	objects, _ := m.vault.QueryObjects("")
	m.groups = buildGroups(objects, m.vault)

	// Find tag group
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

	// Create first object
	m.startCreate(tagIdx, createModeSingle)
	for _, ch := range "my-tag" {
		msg := tea.KeyPressMsg{Code: ch, Text: string(ch)}
		newM, _ := m.Update(msg)
		m = newM.(model)
	}
	msg := tea.KeyPressMsg{Code: tea.KeyEnter}
	newM, _ := m.Update(msg)
	m = newM.(model)

	// Try to create duplicate
	m.startCreate(tagIdx, createModeSingle)
	for _, ch := range "my-tag" {
		msg := tea.KeyPressMsg{Code: ch, Text: string(ch)}
		newM, _ := m.Update(msg)
		m = newM.(model)
	}
	msg = tea.KeyPressMsg{Code: tea.KeyEnter}
	newM, _ = m.Update(msg)
	m = newM.(model)

	if m.create == nil {
		t.Fatal("create should still be active on error")
	}
	if m.create.errMsg == "" {
		t.Error("errMsg should contain duplicate error")
	}
}

// --- Task 6.2: Empty name rejection ---

func TestCreateName_EmptyNameRejected(t *testing.T) {
	m := setupCreateTestModel(t)
	m.startCreate(0, createModeSingle)

	// Press enter without typing
	msg := tea.KeyPressMsg{Code: tea.KeyEnter}
	newM, _ := m.Update(msg)
	m = newM.(model)

	if m.create == nil {
		t.Error("create should still be active (empty name should not create)")
	}
}

// --- Task 6.1: Error clears on input change ---

func TestCreateName_ErrorClearsOnInput(t *testing.T) {
	m := setupCreateTestModel(t)
	m.create = &createState{
		mode:      createModeSingle,
		step:      createStepName,
		typeName:  "book",
		nameInput: initNameInput(),
		errMsg:    "some error",
	}

	msg := tea.KeyPressMsg{Code: 'a', Text: "a"}
	newM, _ := m.Update(msg)
	m = newM.(model)

	if m.create.errMsg != "" {
		t.Errorf("errMsg should be cleared on input, got %q", m.create.errMsg)
	}
}

// --- Task 7.4: Read-only mode ---

func TestCreateKeybindings_ReadOnlyIgnored(t *testing.T) {
	m := setupCreateTestModel(t)
	m.readOnly = true

	// n key
	msg := tea.KeyPressMsg{Code: 'n', Text: "n"}
	newM, _ := m.Update(msg)
	m = newM.(model)
	if m.create != nil {
		t.Error("n should be ignored in read-only mode")
	}

	// N key
	msg = tea.KeyPressMsg{Code: 'N', Text: "N"}
	newM, _ = m.Update(msg)
	m = newM.(model)
	if m.create != nil {
		t.Error("N should be ignored in read-only mode")
	}
}

// --- Rendering tests ---

func TestRenderCreateUI_TemplateStep(t *testing.T) {
	cs := &createState{
		step:      createStepTemplate,
		typeName:  "book",
		templates: []string{"review", "summary"},
		cursor:    0,
	}
	result := renderCreateUI(cs)
	if !strings.Contains(result, "Select template:") {
		t.Error("should contain 'Select template:' header")
	}
	if !strings.Contains(result, "> review") {
		t.Error("should contain cursor at review")
	}
	if !strings.Contains(result, noneOption) {
		t.Errorf("should contain %q option", noneOption)
	}
}

func TestRenderCreateUI_NameStep(t *testing.T) {
	cs := &createState{
		step:      createStepName,
		typeName:  "book",
		nameInput: initNameInput(),
	}
	result := renderCreateUI(cs)
	if !strings.Contains(result, "New book:") {
		t.Error("should contain 'New book:' prompt")
	}
}

func TestRenderCreateUI_NameStepWithError(t *testing.T) {
	cs := &createState{
		step:      createStepName,
		typeName:  "book",
		nameInput: initNameInput(),
		errMsg:    "name already exists",
	}
	result := renderCreateUI(cs)
	if !strings.Contains(result, "✗ name already exists") {
		t.Error("should contain error message with ✗ prefix")
	}
}

func TestRenderCreateUI_WithFlash(t *testing.T) {
	cs := &createState{
		step:      createStepName,
		typeName:  "book",
		nameInput: initNameInput(),
		flash:     "✓ Created: my-book",
	}
	result := renderCreateUI(cs)
	if !strings.Contains(result, "✓ Created: my-book") {
		t.Error("should contain flash message")
	}
}

func TestRenderCreateHelpBar_TemplateStep(t *testing.T) {
	cs := &createState{step: createStepTemplate}
	result := renderCreateHelpBar(cs)
	if !strings.Contains(result, "↑↓") {
		t.Error("should contain navigation hint")
	}
}

func TestRenderCreateHelpBar_NameStepSingle(t *testing.T) {
	cs := &createState{step: createStepName, mode: createModeSingle}
	result := renderCreateHelpBar(cs)
	if !strings.Contains(result, "create & edit") {
		t.Error("should contain 'create & edit' in single mode")
	}
}

func TestRenderCreateHelpBar_NameStepBatch(t *testing.T) {
	cs := &createState{step: createStepName, mode: createModeBatch}
	result := renderCreateHelpBar(cs)
	if !strings.Contains(result, "QUICK CREATE") {
		t.Error("should contain 'QUICK CREATE' in batch mode")
	}
	if !strings.Contains(result, "esc: done") {
		t.Error("should contain 'esc: done' in batch mode")
	}
}
