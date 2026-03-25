package tui

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/typemd/typemd/core"
)

// setupTestVaultForView creates a vault with a book type schema for testing cell editing.
func setupTestVaultForView(t *testing.T) *core.Vault {
	t.Helper()
	dir := t.TempDir()
	v := core.NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatalf("vault Init() error = %v", err)
	}
	schema := `name: book
properties:
  - name: title
    type: string
  - name: rating
    type: number
  - name: status
    type: select
    options:
      - value: draft
      - value: published
      - value: archived
  - name: genres
    type: multi_select
    options:
      - value: fiction
      - value: non-fiction
      - value: sci-fi
  - name: done
    type: checkbox
  - name: author
    type: relation
    target: person
`
	typesDir := v.TypesDir()
	os.MkdirAll(filepath.Join(typesDir, "book"), 0755)
	os.WriteFile(filepath.Join(typesDir, "book", "schema.yaml"), []byte(schema), 0644)
	if err := v.Open(); err != nil {
		t.Fatalf("vault Open() error = %v", err)
	}
	t.Cleanup(func() { v.Close() })
	return v
}

func testViewModeSchema() *core.TypeSchema {
	return &core.TypeSchema{
		Name: "book",
		Properties: []core.Property{
			{Name: "title", Type: "string"},
			{Name: "rating", Type: "number"},
			{Name: "status", Type: "select", Options: []core.Option{
				{Value: "draft"},
				{Value: "published"},
				{Value: "archived"},
			}},
			{Name: "genres", Type: "multi_select", Options: []core.Option{
				{Value: "fiction"},
				{Value: "non-fiction"},
				{Value: "sci-fi"},
			}},
			{Name: "done", Type: "checkbox"},
			{Name: "author", Type: "relation", Target: "person"},
		},
	}
}

func testViewMode() *viewMode {
	schema := testViewModeSchema()
	vm := &viewMode{
		typeName: "book",
		schema:   schema,
		view: &core.ViewConfig{
			Layout:  core.ViewLayoutTable,
			Columns: []string{"title", "rating", "status", "genres", "done", "author"},
		},
		objects: []*core.Object{
			{
				ID:       "book/test-01abc",
				Type:     "book",
				Filename: "test",
				Properties: map[string]any{
					"name":   "Test Book",
					"title":  "Hello World",
					"rating": 5,
					"status": "draft",
					"done":   false,
				},
			},
			{
				ID:       "book/other-02def",
				Type:     "book",
				Filename: "other",
				Properties: map[string]any{
					"name":   "Other Book",
					"title":  "Goodbye",
					"rating": 3,
					"status": "published",
				},
			},
		},
		width:  120,
		height: 30,
	}
	vm.groups = []viewGroup{{
		Label:    "",
		Objects:  vm.objects,
		Expanded: true,
	}}
	return vm
}

func TestIsCellReadOnly(t *testing.T) {
	vm := testViewMode()
	cols := vm.viewColumns()

	tests := []struct {
		name     string
		colIdx   int
		readOnly bool
	}{
		{"NAME column is editable", 0, false},
		{"title (string) is editable", 1, false},
		{"rating (number) is editable", 2, false},
		{"status (select) is editable", 3, false},
		{"genres (multi_select) is editable", 4, false},
		{"done (checkbox) is editable", 5, false},
		{"author (relation) is read-only", 6, true},
		{"out of bounds is read-only", 99, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := vm.isCellReadOnly(tt.colIdx, cols)
			if got != tt.readOnly {
				t.Errorf("isCellReadOnly(%d) = %v, want %v", tt.colIdx, got, tt.readOnly)
			}
		})
	}
}

func TestIsCellReadOnly_SystemProperties(t *testing.T) {
	vm := &viewMode{
		schema: &core.TypeSchema{Name: "test"},
		view: &core.ViewConfig{
			Layout:  core.ViewLayoutTable,
			Columns: []string{"created_at", "updated_at", "description"},
		},
	}
	cols := vm.viewColumns()

	tests := []struct {
		name     string
		colIdx   int
		readOnly bool
	}{
		{"created_at is read-only", 1, true},
		{"updated_at is read-only", 2, true},
		{"description is editable", 3, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := vm.isCellReadOnly(tt.colIdx, cols)
			if got != tt.readOnly {
				t.Errorf("isCellReadOnly(%d) = %v, want %v", tt.colIdx, got, tt.readOnly)
			}
		})
	}
}

func TestActivateCellEdit_ReadOnlyReturnsNil(t *testing.T) {
	vm := testViewMode()
	cols := vm.viewColumns()
	rows := vm.visibleRows()

	// author is a relation (col index 6) — should return nil
	cmd := vm.activateCellEdit(rows[0], 6, cols)
	if cmd != nil {
		t.Error("activateCellEdit should return nil for read-only cell")
	}
	if vm.cellEdit != nil {
		t.Error("cellEdit should remain nil for read-only cell")
	}
}

func TestActivateCellEdit_StringProperty(t *testing.T) {
	vm := testViewMode()
	cols := vm.viewColumns()
	rows := vm.visibleRows()

	// title is a string (col index 1)
	cmd := vm.activateCellEdit(rows[0], 1, cols)
	if cmd == nil {
		t.Error("activateCellEdit should return a blink cmd for textinput")
	}
	if vm.cellEdit == nil {
		t.Fatal("cellEdit should be set")
	}
	if vm.cellEdit.propName != "title" {
		t.Errorf("propName = %q, want %q", vm.cellEdit.propName, "title")
	}
	if vm.cellEdit.propType != "string" {
		t.Errorf("propType = %q, want %q", vm.cellEdit.propType, "string")
	}
	if vm.cellEdit.mode != cellModeTextInput {
		t.Errorf("mode = %d, want cellModeTextInput", vm.cellEdit.mode)
	}
	if vm.cellEdit.textInput.Value() != "Hello World" {
		t.Errorf("textInput value = %q, want %q", vm.cellEdit.textInput.Value(), "Hello World")
	}
}

func TestActivateCellEdit_NameColumn(t *testing.T) {
	vm := testViewMode()
	cols := vm.viewColumns()
	rows := vm.visibleRows()

	cmd := vm.activateCellEdit(rows[0], 0, cols)
	if cmd == nil {
		t.Error("activateCellEdit should return a blink cmd for NAME textinput")
	}
	if vm.cellEdit == nil {
		t.Fatal("cellEdit should be set")
	}
	if vm.cellEdit.propName != core.NameProperty {
		t.Errorf("propName = %q, want %q", vm.cellEdit.propName, core.NameProperty)
	}
	if vm.cellEdit.textInput.Value() != "Test Book" {
		t.Errorf("textInput value = %q, want %q", vm.cellEdit.textInput.Value(), "Test Book")
	}
}

func TestActivateCellEdit_SelectProperty(t *testing.T) {
	vm := testViewMode()
	cols := vm.viewColumns()
	rows := vm.visibleRows()

	// status is a select (col index 3)
	cmd := vm.activateCellEdit(rows[0], 3, cols)
	if cmd != nil {
		t.Error("activateCellEdit should return nil cmd for select picker (no blink)")
	}
	if vm.cellEdit == nil {
		t.Fatal("cellEdit should be set")
	}
	if vm.cellEdit.mode != cellModeSelectPick {
		t.Errorf("mode = %d, want cellModeSelectPick", vm.cellEdit.mode)
	}
	if len(vm.cellEdit.pickerOptions) != 3 {
		t.Errorf("pickerOptions len = %d, want 3", len(vm.cellEdit.pickerOptions))
	}
	// Current value is "draft" — should be at index 0
	if vm.cellEdit.pickerCursor != 0 {
		t.Errorf("pickerCursor = %d, want 0 (current value 'draft')", vm.cellEdit.pickerCursor)
	}
}

func TestActivateCellEdit_MultiSelectProperty(t *testing.T) {
	vm := testViewMode()
	cols := vm.viewColumns()

	// Set genres on the object
	vm.objects[0].Properties["genres"] = []any{"fiction", "sci-fi"}
	rows := vm.visibleRows()

	// genres is multi_select (col index 4)
	cmd := vm.activateCellEdit(rows[0], 4, cols)
	if cmd != nil {
		t.Error("activateCellEdit should return nil cmd for multi_select picker")
	}
	if vm.cellEdit == nil {
		t.Fatal("cellEdit should be set")
	}
	if vm.cellEdit.mode != cellModeMultiPick {
		t.Errorf("mode = %d, want cellModeMultiPick", vm.cellEdit.mode)
	}
	// fiction=checked, non-fiction=unchecked, sci-fi=checked
	if !vm.cellEdit.pickerChecked[0] {
		t.Error("fiction should be checked")
	}
	if vm.cellEdit.pickerChecked[1] {
		t.Error("non-fiction should not be checked")
	}
	if !vm.cellEdit.pickerChecked[2] {
		t.Error("sci-fi should be checked")
	}
}

func TestActivateCellEdit_CheckboxToggles(t *testing.T) {
	v := setupTestVaultForView(t)
	vm := testViewMode()
	vm.vault = v

	// Create a real object in the vault so save works
	obj, err := v.NewObject("book", "toggle-test", "")
	if err != nil {
		t.Fatal(err)
	}
	obj.Properties["done"] = false
	if err := v.SaveObject(obj); err != nil {
		t.Fatal(err)
	}

	// Replace the test object with the real one
	vm.objects[0] = obj
	vm.groups[0].Objects[0] = obj

	cols := vm.viewColumns()
	rows := vm.visibleRows()

	// done is checkbox (col index 5), currently false
	before := obj.Properties["done"]
	if before != false {
		t.Fatalf("precondition: done should be false, got %v", before)
	}

	// activateCellEdit for checkbox toggles directly
	cmd := vm.activateCellEdit(rows[0], 5, cols)

	// cellEdit should NOT be set (checkbox toggles directly)
	if vm.cellEdit != nil {
		t.Error("cellEdit should be nil after checkbox toggle")
	}

	// Value should be toggled
	after := obj.Properties["done"]
	if after != true {
		t.Errorf("done should be true after toggle, got %v", after)
	}

	// cmd should be non-nil (save command)
	if cmd == nil {
		t.Error("checkbox toggle should return a save command")
	}
}

func TestActivateCellEdit_HeaderRow(t *testing.T) {
	vm := testViewMode()
	cols := vm.viewColumns()

	// Create a header row
	headerRow := viewRow{isHeader: true, groupIdx: 0, label: "Group A"}
	cmd := vm.activateCellEdit(headerRow, 0, cols)
	if cmd != nil {
		t.Error("activateCellEdit should return nil for header row")
	}
	if vm.cellEdit != nil {
		t.Error("cellEdit should be nil for header row")
	}
}

func TestCancelCellEdit(t *testing.T) {
	vm := testViewMode()
	vm.cellEdit = &cellEdit{propName: "title"}

	vm.cancelCellEdit()
	if vm.cellEdit != nil {
		t.Error("cancelCellEdit should set cellEdit to nil")
	}
}

func TestIsTableLayout(t *testing.T) {
	vm := testViewMode()
	if !vm.isTableLayout() {
		t.Error("isTableLayout should return true for table layout")
	}

	vm.view.Layout = core.ViewLayoutList
	if vm.isTableLayout() {
		t.Error("isTableLayout should return false for list layout")
	}
}

func TestNextEditableCell(t *testing.T) {
	vm := testViewMode()
	vm.colCursor = 0 // NAME
	vm.cursor = 0

	// NAME → title (skip nothing)
	if !vm.nextEditableCell() {
		t.Error("should move to next editable cell")
	}
	if vm.colCursor != 1 {
		t.Errorf("colCursor = %d, want 1 (title)", vm.colCursor)
	}
}

func TestNextEditableCell_SkipsReadOnly(t *testing.T) {
	vm := testViewMode()
	// Start at done (col 5), next should skip author (col 6, relation) and wrap to next row NAME (col 0)
	vm.colCursor = 5
	vm.cursor = 0

	if !vm.nextEditableCell() {
		t.Error("should move to next editable cell")
	}
	// Should wrap to next row col 0 (NAME), skipping author (relation)
	if vm.cursor != 1 || vm.colCursor != 0 {
		t.Errorf("cursor=%d colCursor=%d, want cursor=1 colCursor=0", vm.cursor, vm.colCursor)
	}
}

func TestPrevEditableCell(t *testing.T) {
	vm := testViewMode()
	vm.colCursor = 1 // title
	vm.cursor = 0

	if !vm.prevEditableCell() {
		t.Error("should move to previous editable cell")
	}
	if vm.colCursor != 0 {
		t.Errorf("colCursor = %d, want 0 (NAME)", vm.colCursor)
	}
}

func TestColCursorClampedInViewTable(t *testing.T) {
	vm := testViewMode()
	vm.colCursor = 100

	// viewTable clamps colCursor internally
	rows := vm.visibleRows()
	vm.viewTable(rows)

	cols := vm.viewColumns()
	if vm.colCursor != len(cols) {
		t.Errorf("colCursor = %d, want %d after viewTable clamp", vm.colCursor, len(cols))
	}
}

func TestCanQuit_CellEditBlocks(t *testing.T) {
	vm := testViewMode()
	if !vm.CanQuit() {
		t.Error("should be able to quit with no cell edit")
	}

	vm.cellEdit = &cellEdit{}
	if vm.CanQuit() {
		t.Error("should not be able to quit with active cell edit")
	}
}

func TestHelpBar_TableMode(t *testing.T) {
	vm := testViewMode()
	help := vm.HelpBar()
	if help == "" {
		t.Error("HelpBar should not be empty")
	}
	// Table mode help should mention cell editing keys
	if !containsAll(help, "enter", "edit cell", "←/→") {
		t.Errorf("HelpBar = %q, should mention cell editing keys", help)
	}
}

func TestHelpBar_CellEditMode(t *testing.T) {
	vm := testViewMode()
	vm.cellEdit = &cellEdit{mode: cellModeTextInput}
	help := vm.HelpBar()
	if !containsAll(help, "EDIT", "confirm", "cancel") {
		t.Errorf("HelpBar = %q, should show EDIT mode hints", help)
	}
}

func TestHelpBar_PickerMode(t *testing.T) {
	vm := testViewMode()
	vm.cellEdit = &cellEdit{mode: cellModeSelectPick}
	help := vm.HelpBar()
	if !containsAll(help, "PICK", "select", "cancel") {
		t.Errorf("HelpBar = %q, should show PICK mode hints", help)
	}
}

func containsAll(s string, substrs ...string) bool {
	for _, sub := range substrs {
		found := false
		for i := 0; i <= len(s)-len(sub); i++ {
			if s[i:i+len(sub)] == sub {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func TestViewTable_CrosshairRendering(t *testing.T) {
	vm := testViewMode()
	vm.cursor = 0
	vm.colCursor = 1 // title column

	rows := vm.visibleRows()
	output := vm.viewTable(rows)

	// Should contain object data
	if output == "" {
		t.Error("viewTable should produce non-empty output")
	}
	// Header should be present
	if !containsAll(output, "NAME", "TITLE") {
		t.Errorf("viewTable should contain column headers, got: %q", output[:min(len(output), 200)])
	}
}

func TestViewTable_EmptySchema(t *testing.T) {
	vm := &viewMode{
		typeName: "test",
		schema:   &core.TypeSchema{Name: "test"},
		view: &core.ViewConfig{
			Layout: core.ViewLayoutTable,
		},
		objects: []*core.Object{
			{
				ID:       "test/a-01",
				Type:     "test",
				Filename: "a",
				Properties: map[string]any{
					"name": "Test A",
				},
			},
		},
		width:  80,
		height: 20,
	}
	vm.groups = []viewGroup{{Objects: vm.objects, Expanded: true}}

	rows := vm.visibleRows()
	output := vm.viewTable(rows)
	if output == "" {
		t.Error("viewTable should produce output even with no property columns")
	}
	if !containsAll(output, "NAME") {
		t.Error("viewTable should have NAME column header")
	}
}
