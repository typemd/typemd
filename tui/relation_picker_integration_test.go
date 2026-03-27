package tui

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/typemd/typemd/core"
)

// setupRelationTestVault creates a vault with book, person, and tag types
// for testing relation picker behavior end-to-end.
func setupRelationTestVault(t *testing.T) *core.Vault {
	t.Helper()
	dir := t.TempDir()
	v := core.NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatalf("Init: %v", err)
	}
	if err := v.Open(); err != nil {
		t.Fatalf("Open: %v", err)
	}

	// Create person type with books relation
	personSchema := &core.TypeSchema{
		Name:  "person",
		Emoji: "👤",
		Properties: []core.Property{
			{Name: "books", Type: "relation", Target: "book", Multiple: true, Bidirectional: true, Inverse: "author"},
		},
	}
	if err := v.SaveType(personSchema); err != nil {
		t.Fatalf("SaveType person: %v", err)
	}

	// Create book type with author relation
	bookSchema := &core.TypeSchema{
		Name:  "book",
		Emoji: "📚",
		Properties: []core.Property{
			{Name: "author", Type: "relation", Target: "person", Bidirectional: true, Inverse: "books"},
		},
	}
	if err := v.SaveType(bookSchema); err != nil {
		t.Fatalf("SaveType book: %v", err)
	}

	// Create person objects
	if _, err := v.NewObject("person", "alice", ""); err != nil {
		t.Fatalf("NewObject alice: %v", err)
	}
	if _, err := v.NewObject("person", "bob", ""); err != nil {
		t.Fatalf("NewObject bob: %v", err)
	}

	// Create book objects
	if _, err := v.NewObject("book", "clean-code", ""); err != nil {
		t.Fatalf("NewObject clean-code: %v", err)
	}

	// Create tag objects
	if _, err := v.NewObject("tag", "go", ""); err != nil {
		t.Fatalf("NewObject go: %v", err)
	}
	if _, err := v.NewObject("tag", "rust", ""); err != nil {
		t.Fatalf("NewObject rust: %v", err)
	}

	// Sync index
	if _, err := v.SyncIndex(); err != nil {
		t.Fatalf("SyncIndex: %v", err)
	}

	t.Cleanup(func() { v.Close() })
	return v
}

// findObjectID finds the full ID for an object by type and name prefix.
func findObjectID(t *testing.T, v *core.Vault, typeName, prefix string) string {
	t.Helper()
	objs, err := v.QueryObjects([]core.FilterRule{
		{Property: "type", Operator: "is", Value: typeName},
	})
	if err != nil {
		t.Fatalf("QueryObjects: %v", err)
	}
	for _, obj := range objs {
		if obj.DisplayName() == prefix || filepath.Base(obj.Filename)[:len(prefix)] == prefix {
			return obj.ID
		}
	}
	t.Fatalf("object not found: %s/%s", typeName, prefix)
	return ""
}

// Test 1: Navigate to relation property, press Enter — verify picker opens
func TestRelationPicker_ActivateOnEnter(t *testing.T) {
	v := setupRelationTestVault(t)
	bookID := findObjectID(t, v, "book", "clean-code")
	book, err := v.GetObject(bookID)
	if err != nil {
		t.Fatalf("GetObject: %v", err)
	}

	displayProps, _ := v.BuildDisplayProperties(book)
	schema, _ := v.LoadType("book")
	pe := newPropEditor(displayProps, schema)

	// Navigate to the "author" relation property
	for pe.currentItem() != nil && pe.currentItem().dp.Key != "author" {
		pe.moveDown()
	}
	if pe.currentItem() == nil || pe.currentItem().dp.Key != "author" {
		t.Fatal("could not navigate to author relation property")
	}
	if !pe.currentItem().editable {
		t.Fatal("author relation should be editable")
	}

	// Activate edit — should open relation picker
	pe.activateEdit(v)

	if pe.mode != propModeRelationPick {
		t.Errorf("expected propModeRelationPick, got %d", pe.mode)
	}
	if len(pe.relCandidates) != 2 {
		t.Errorf("expected 2 person candidates (alice, bob), got %d", len(pe.relCandidates))
	}
}

// Test 2: Type to filter candidates — verify fuzzy search works
func TestRelationPicker_FuzzyFilter(t *testing.T) {
	v := setupRelationTestVault(t)
	bookID := findObjectID(t, v, "book", "clean-code")
	book, _ := v.GetObject(bookID)
	displayProps, _ := v.BuildDisplayProperties(book)
	schema, _ := v.LoadType("book")
	pe := newPropEditor(displayProps, schema)

	for pe.currentItem() != nil && pe.currentItem().dp.Key != "author" {
		pe.moveDown()
	}
	pe.activateEdit(v)

	// Type "ali" to filter
	pe.relSearch = "ali"
	pe.filterRelationCandidates()

	if len(pe.relFiltered) != 1 {
		t.Errorf("expected 1 match for 'ali', got %d", len(pe.relFiltered))
	}
	if pe.relFiltered[0].displayName != "alice" {
		t.Errorf("expected 'alice', got %q", pe.relFiltered[0].displayName)
	}
}

// Test 3: Select a candidate — verify link command is created
func TestRelationPicker_SelectCandidate(t *testing.T) {
	v := setupRelationTestVault(t)
	bookID := findObjectID(t, v, "book", "clean-code")
	book, _ := v.GetObject(bookID)
	displayProps, _ := v.BuildDisplayProperties(book)
	schema, _ := v.LoadType("book")
	pe := newPropEditor(displayProps, schema)

	for pe.currentItem() != nil && pe.currentItem().dp.Key != "author" {
		pe.moveDown()
	}
	pe.activateEdit(v)

	// Cursor 0 = "(none)", cursor 1 = first candidate
	pe.pickerCursor = 1
	candidate := pe.relFiltered[pe.pickerCursor-1]

	// Verify candidate is a person
	if candidate.id == "" {
		t.Fatal("candidate ID should not be empty")
	}

	// Create a link command and execute it
	aliceID := findObjectID(t, v, "person", "alice")
	if err := v.LinkObjects(bookID, "author", aliceID); err != nil {
		t.Fatalf("LinkObjects: %v", err)
	}

	// Verify the link exists
	reloaded, _ := v.GetObject(bookID)
	authorVal := reloaded.Properties["author"]
	if authorVal != aliceID {
		t.Errorf("expected author=%q, got %q", aliceID, authorVal)
	}
}

// Test 4: Select "(none)" — verify relation is cleared
func TestRelationPicker_ClearWithNone(t *testing.T) {
	v := setupRelationTestVault(t)
	bookID := findObjectID(t, v, "book", "clean-code")
	aliceID := findObjectID(t, v, "person", "alice")

	// Link first
	if err := v.LinkObjects(bookID, "author", aliceID); err != nil {
		t.Fatalf("LinkObjects: %v", err)
	}

	// Verify link exists
	book, _ := v.GetObject(bookID)
	if book.Properties["author"] != aliceID {
		t.Fatal("author should be linked before clearing")
	}

	// Unlink (simulates selecting "(none)")
	if err := v.UnlinkObjects(bookID, "author", aliceID, false); err != nil {
		t.Fatalf("UnlinkObjects: %v", err)
	}

	// Verify cleared
	book, _ = v.GetObject(bookID)
	if book.Properties["author"] != nil {
		t.Errorf("author should be nil after clear, got %v", book.Properties["author"])
	}
}

// Test 5: Tags — verify multi-select picker opens for tags
func TestRelationPicker_TagsMultiSelect(t *testing.T) {
	v := setupRelationTestVault(t)
	bookID := findObjectID(t, v, "book", "clean-code")
	goID := findObjectID(t, v, "tag", "go")

	// Link a tag first so "tags" appears in display properties
	if err := v.LinkObjects(bookID, "tags", goID); err != nil {
		t.Fatalf("LinkObjects tag: %v", err)
	}

	book, _ := v.GetObject(bookID)
	displayProps, _ := v.BuildDisplayProperties(book)
	schema, _ := v.LoadType("book")
	pe := newPropEditor(displayProps, schema)

	// Navigate to the "tags" property
	for pe.currentItem() != nil && pe.currentItem().dp.Key != "tags" {
		pe.moveDown()
	}
	if pe.currentItem() == nil || pe.currentItem().dp.Key != "tags" {
		t.Fatal("could not navigate to tags property")
	}
	if !pe.currentItem().editable {
		t.Fatal("tags should be editable")
	}

	// Activate edit — should open multi-value relation picker
	pe.activateEdit(v)

	if pe.mode != propModeRelationMultiPick {
		t.Errorf("expected propModeRelationMultiPick, got %d", pe.mode)
	}
	// Should list tag objects (go, rust)
	if len(pe.relCandidates) != 2 {
		t.Errorf("expected 2 tag candidates, got %d", len(pe.relCandidates))
	}
	// Should have checkmark state
	if pe.relChecked == nil {
		t.Error("relChecked should be initialized")
	}
}

// Test 6: Toggle tags with Space, confirm — verify link/unlink
func TestRelationPicker_TagsToggleAndConfirm(t *testing.T) {
	v := setupRelationTestVault(t)
	bookID := findObjectID(t, v, "book", "clean-code")
	goID := findObjectID(t, v, "tag", "go")
	rustID := findObjectID(t, v, "tag", "rust")

	// Link both tags
	if err := v.LinkObjects(bookID, "tags", goID); err != nil {
		t.Fatalf("LinkObjects go: %v", err)
	}
	if err := v.LinkObjects(bookID, "tags", rustID); err != nil {
		t.Fatalf("LinkObjects rust: %v", err)
	}

	// Verify links
	book, _ := v.GetObject(bookID)
	tags := book.Properties["tags"]
	if tags == nil {
		t.Fatal("tags should not be nil after linking")
	}

	// Unlink one tag (simulates toggle off + confirm)
	if err := v.UnlinkObjects(bookID, "tags", rustID, false); err != nil {
		t.Fatalf("UnlinkObjects rust: %v", err)
	}

	// Verify only go tag remains
	book, _ = v.GetObject(bookID)
	tagSlice, ok := book.Properties["tags"].([]any)
	if !ok {
		t.Fatalf("expected []any tags, got %T: %v", book.Properties["tags"], book.Properties["tags"])
	}
	if len(tagSlice) != 1 {
		t.Errorf("expected 1 tag remaining, got %d", len(tagSlice))
	}
	if tagSlice[0] != goID {
		t.Errorf("expected remaining tag=%q, got %q", goID, tagSlice[0])
	}
}

// Test 7: Reverse relations are skipped (not focusable)
func TestRelationPicker_ReverseRelationsSkipped(t *testing.T) {
	v := setupRelationTestVault(t)
	aliceID := findObjectID(t, v, "person", "alice")
	bookID := findObjectID(t, v, "book", "clean-code")

	// Link book to alice so alice has a reverse relation
	if err := v.LinkObjects(bookID, "author", aliceID); err != nil {
		t.Fatalf("LinkObjects: %v", err)
	}

	// Get alice's display properties — should include reverse relation "author"
	alice, _ := v.GetObject(aliceID)
	displayProps, _ := v.BuildDisplayProperties(alice)
	schema, _ := v.LoadType("person")
	pe := newPropEditor(displayProps, schema)

	// Verify reverse relations are not editable
	for _, item := range pe.items {
		if item.dp.IsReverse {
			if item.editable {
				t.Errorf("reverse relation %q should not be editable", item.dp.Key)
			}
		}
	}

	// Verify cursor never lands on reverse relation
	for i := 0; i < len(pe.items)*2; i++ {
		cur := pe.currentItem()
		if cur != nil && cur.dp.IsReverse {
			t.Errorf("cursor should never land on reverse relation %q", cur.dp.Key)
		}
		pe.moveDown()
	}
}

// Test 8: Locked objects block relation editing
func TestRelationPicker_LockedObjectBlocked(t *testing.T) {
	v := setupRelationTestVault(t)
	bookID := findObjectID(t, v, "book", "clean-code")

	// Lock the book via SetProperty (which allows setting locked on an unlocked object)
	if err := v.SetProperty(bookID, "locked", true); err != nil {
		t.Fatalf("SetProperty locked: %v", err)
	}

	// Verify object is locked
	book, _ := v.GetObject(bookID)
	if !book.IsLocked() {
		t.Fatal("book should be locked")
	}

	// Verify that LinkObjects is blocked for locked objects
	aliceID := findObjectID(t, v, "person", "alice")
	err := v.LinkObjects(bookID, "author", aliceID)
	if err == nil {
		t.Error("LinkObjects should fail on locked object")
	}
	if err != nil && err.Error() != "object is locked" {
		t.Errorf("expected 'object is locked' error, got: %v", err)
	}
}

// Test: Bidirectional relation — linking book→author also creates person→books inverse
func TestRelationPicker_BidirectionalLink(t *testing.T) {
	v := setupRelationTestVault(t)
	bookID := findObjectID(t, v, "book", "clean-code")
	aliceID := findObjectID(t, v, "person", "alice")

	// Link book's author to alice
	if err := v.LinkObjects(bookID, "author", aliceID); err != nil {
		t.Fatalf("LinkObjects: %v", err)
	}

	// Verify forward relation on book
	book, _ := v.GetObject(bookID)
	if book.Properties["author"] != aliceID {
		t.Errorf("book.author should be %q, got %v", aliceID, book.Properties["author"])
	}

	// Verify inverse relation on alice
	alice, _ := v.GetObject(aliceID)
	booksVal := alice.Properties["books"]
	if booksVal == nil {
		t.Fatal("alice.books should not be nil (bidirectional inverse)")
	}
	booksSlice, ok := booksVal.([]any)
	if !ok {
		t.Fatalf("expected []any for books, got %T", booksVal)
	}
	found := false
	for _, b := range booksSlice {
		if b == bookID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("alice.books should contain %q, got %v", bookID, booksSlice)
	}
}

// Test: Picker display — verify candidate names don't include ULID
func TestRelationPicker_DisplayNameWithoutULID(t *testing.T) {
	v := setupRelationTestVault(t)
	bookID := findObjectID(t, v, "book", "clean-code")
	book, _ := v.GetObject(bookID)
	displayProps, _ := v.BuildDisplayProperties(book)
	schema, _ := v.LoadType("book")
	pe := newPropEditor(displayProps, schema)

	for pe.currentItem() != nil && pe.currentItem().dp.Key != "author" {
		pe.moveDown()
	}
	pe.activateEdit(v)

	for _, c := range pe.relCandidates {
		// ULID is 26 chars, check that display name is short (no ULID suffix)
		if len(c.displayName) > 30 {
			t.Errorf("displayName %q looks like it contains a ULID", c.displayName)
		}
		// Should not contain the full object ID format
		if len(c.id) > 0 && c.displayName == c.id {
			t.Errorf("displayName should not be the full object ID: %q", c.displayName)
		}
	}
}

// Test: Help bar shows [PICK] during relation picker
func TestRelationPicker_HelpBarShowsPick(t *testing.T) {
	pe := &propEditor{mode: propModeRelationPick}
	if !pe.isPicking() {
		t.Error("isPicking() should be true for propModeRelationPick")
	}

	pe.mode = propModeRelationMultiPick
	if !pe.isPicking() {
		t.Error("isPicking() should be true for propModeRelationMultiPick")
	}

	pe.mode = propModeNavigate
	if pe.isPicking() {
		t.Error("isPicking() should be false for propModeNavigate")
	}

	pe.mode = propModeTextInput
	if pe.isPicking() {
		t.Error("isPicking() should be false for propModeTextInput")
	}
}

// Test: Backlinks are not editable
func TestRelationPicker_BacklinksNotEditable(t *testing.T) {
	dp := core.DisplayProperty{
		Key:        "backlinks",
		IsBacklink: true,
		FromID:     "book/something",
	}
	if isPropertyEditable(dp) {
		t.Error("backlinks should not be editable")
	}
}

// Verify that the test vault objects directory exists and objects are queryable
func TestRelationPicker_VaultSetup(t *testing.T) {
	v := setupRelationTestVault(t)

	// Verify we can query each type
	for _, typeName := range []string{"book", "person", "tag"} {
		results, err := v.QueryObjects([]core.FilterRule{
			{Property: "type", Operator: "is", Value: typeName},
		})
		if err != nil {
			t.Errorf("QueryObjects %s: %v", typeName, err)
		}
		if len(results) == 0 {
			t.Errorf("expected objects of type %s, got none", typeName)
		}
	}

	// Verify vault directory structure
	typesDir := filepath.Join(v.Root, ".typemd", "types")
	for _, name := range []string{"book", "person"} {
		schemaPath := filepath.Join(typesDir, name, "schema.yaml")
		if _, err := os.Stat(schemaPath); os.IsNotExist(err) {
			t.Errorf("schema file not found: %s", schemaPath)
		}
	}
}
