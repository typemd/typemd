package core

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cucumber/godog"
)

func mustULID() string {
	id, err := GenerateULID()
	if err != nil {
		panic(fmt.Sprintf("GenerateULID failed: %v", err))
	}
	return id
}

// domainContext holds shared state across steps within a single scenario.
type domainContext struct {
	rootDir       string
	vault         *Vault
	objects       map[string]*Object // keyed by slug (e.g. "golang-in-action")
	currentObject *Object            // the most recently created/referenced object
	retrieved     *Object            // result of GetObject-by-ID
	prevObject    *Object            // for "create another" pattern
	lastErr       error

	// query/search results
	queryResults  []*Object
	searchResults []*Object

	// validation results
	schemaErrors   map[string][]error
	relationErrors []error
	wikiLinkErrors []error

	// wikilink results
	wikiLinks []StoredWikiLink

	// resolve results
	resolvedID string

	// property type validation results
	objectValidationErrors []error
	schemaMigrateResult    *SchemaMigrateResult

	// shared properties results
	sharedProperties []Property
	loadedSchema     *TypeSchema

	// system property tracking
	createdAtSnapshot string // snapshot of created_at after object creation
}

func newDomainContext() *domainContext {
	return &domainContext{
		objects: make(map[string]*Object),
	}
}

// ── Vault steps ─────────────────────────────────────────────────────────────

func (dc *domainContext) setupVaultDir() {
	dc.rootDir = filepath.Join(os.TempDir(), "typemd-bdd-"+mustULID())
	os.MkdirAll(dc.rootDir, 0755)
	dc.vault = NewVault(dc.rootDir)
}

func (dc *domainContext) iInitializeANewVault() {
	dc.setupVaultDir()
	dc.lastErr = dc.vault.Init()
}

func (dc *domainContext) aVaultIsInitialized() {
	dc.setupVaultDir()
	if err := dc.vault.Init(); err != nil {
		panic(fmt.Sprintf("vault init failed: %v", err))
	}
}

func (dc *domainContext) theVaultDirectoryStructureShouldExist() error {
	for _, d := range []string{dc.vault.TypesDir(), dc.vault.ObjectsDir()} {
		if _, err := os.Stat(d); os.IsNotExist(err) {
			return fmt.Errorf("expected directory %s to exist", d)
		}
	}
	return nil
}

func (dc *domainContext) theSQLiteDatabaseShouldExist() error {
	if _, err := os.Stat(dc.vault.DBPath()); os.IsNotExist(err) {
		return fmt.Errorf("expected index.db to exist")
	}
	return nil
}

func (dc *domainContext) theGitignoreShouldContain(expected string) error {
	data, err := os.ReadFile(filepath.Join(dc.vault.Dir(), ".gitignore"))
	if err != nil {
		return fmt.Errorf("expected .gitignore to exist: %v", err)
	}
	if !strings.Contains(string(data), expected) {
		return fmt.Errorf(".gitignore content = %q, want to contain %q", string(data), expected)
	}
	return nil
}

func (dc *domainContext) iInitializeTheVaultAgain() {
	dc.lastErr = dc.vault.Init()
}

func (dc *domainContext) iOpenTheVault() {
	dc.lastErr = dc.vault.Open()
}

func (dc *domainContext) iCloseTheVault() {
	if dc.lastErr == nil {
		dc.lastErr = dc.vault.Close()
	}
}

func (dc *domainContext) iOpenAnUninitializedVault() {
	dc.setupVaultDir()
	dc.lastErr = dc.vault.Open()
}

func (dc *domainContext) anObjectFileExistsOnDisk(relPath, title string) {
	fullPath := filepath.Join(dc.rootDir, "objects", relPath)
	os.MkdirAll(filepath.Dir(fullPath), 0755)
	content := fmt.Sprintf("---\ntitle: %s\n---\nHello world\n", title)
	os.WriteFile(fullPath, []byte(content), 0644)
}

func (dc *domainContext) theIndexShouldContainNObjects(expected int) error {
	var count int
	if err := dc.vault.db.QueryRow("SELECT COUNT(*) FROM objects").Scan(&count); err != nil {
		return fmt.Errorf("count query error: %v", err)
	}
	if count != expected {
		return fmt.Errorf("objects count = %d, want %d", count, expected)
	}
	return nil
}

// ── Object steps ────────────────────────────────────────────────────────────

func (dc *domainContext) aVaultIsReady() {
	dc.aVaultIsInitialized()
	if err := dc.vault.Open(); err != nil {
		panic(fmt.Sprintf("vault open failed: %v", err))
	}
}

func (dc *domainContext) iCreateAObjectNamed(typeName, name string) {
	obj, err := dc.vault.NewObject(typeName, name)
	dc.lastErr = err
	if err == nil {
		dc.objects[name] = obj
		dc.currentObject = obj
	}
}

func (dc *domainContext) iCreateAnotherObjectNamed(typeName, name string) {
	dc.prevObject = dc.currentObject
	obj, err := dc.vault.NewObject(typeName, name)
	dc.lastErr = err
	if err == nil {
		dc.objects[name+"_2"] = obj
		dc.currentObject = obj
	}
}

func (dc *domainContext) theObjectFilenameShouldStartWith(prefix string) error {
	for _, obj := range dc.objects {
		if strings.HasPrefix(obj.Filename, prefix) {
			return nil
		}
	}
	return fmt.Errorf("no object filename starts with %q", prefix)
}

func (dc *domainContext) theObjectFilenameShouldHaveACharacterULIDSuffix(length int) error {
	for _, obj := range dc.objects {
		parts := strings.SplitN(obj.Filename, "-", 2)
		if len(parts) < 2 {
			continue
		}
		// ULID is the last 26 chars
		if len(obj.Filename) >= length {
			ulidPart := obj.Filename[len(obj.Filename)-length:]
			if len(ulidPart) == length {
				return nil
			}
		}
	}
	return fmt.Errorf("no object has a %d-character ULID suffix", length)
}

func (dc *domainContext) theObjectTypeShouldBe(expected string) error {
	for _, obj := range dc.objects {
		if obj.Type == expected {
			return nil
		}
	}
	return fmt.Errorf("no object has type %q", expected)
}

func (dc *domainContext) theObjectFileShouldExistOnDisk() error {
	for _, obj := range dc.objects {
		path := dc.vault.ObjectPath(obj.Type, obj.Filename)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return fmt.Errorf("expected file %s to exist", path)
		}
	}
	return nil
}

func (dc *domainContext) theTwoObjectsShouldHaveDifferentIDs() error {
	if dc.prevObject == nil {
		return fmt.Errorf("no previous object to compare")
	}
	for _, obj := range dc.objects {
		if obj != dc.prevObject && obj.ID == dc.prevObject.ID {
			return fmt.Errorf("expected different IDs, both got %q", obj.ID)
		}
	}
	return nil
}

func (dc *domainContext) aObjectNamedExists(typeName, name string) {
	obj, err := dc.vault.NewObject(typeName, name)
	if err != nil {
		panic(fmt.Sprintf("create object %s/%s failed: %v", typeName, name, err))
	}
	dc.objects[name] = obj
	dc.prevObject = dc.currentObject
	dc.currentObject = obj
	// Snapshot created_at for system property tests
	if val, ok := obj.Properties["created_at"]; ok {
		dc.createdAtSnapshot = fmt.Sprintf("%v", val)
	}
}

func (dc *domainContext) iGetTheObjectByItsID() {
	got, err := dc.vault.GetObject(dc.currentObject.ID)
	dc.lastErr = err
	if err == nil {
		dc.retrieved = got
	}
}

func (dc *domainContext) theRetrievedObjectShouldMatchTheCreatedOne() error {
	if dc.retrieved == nil {
		return fmt.Errorf("no retrieved object found")
	}
	if dc.retrieved.ID != dc.currentObject.ID {
		return fmt.Errorf("retrieved ID = %q, want %q", dc.retrieved.ID, dc.currentObject.ID)
	}
	return nil
}

func (dc *domainContext) iSetPropertyToOnTheObject(key, value string) {
	dc.lastErr = dc.vault.SetProperty(dc.currentObject.ID, key, value)
}

func (dc *domainContext) theObjectPropertyShouldBe(key, expected string) error {
	got, err := dc.vault.GetObject(dc.currentObject.ID)
	if err != nil {
		return fmt.Errorf("GetObject error: %v", err)
	}
	val := fmt.Sprintf("%v", got.Properties[key])
	if val != expected {
		return fmt.Errorf("property %q = %q, want %q", key, val, expected)
	}
	return nil
}

func (dc *domainContext) iUpdateTheObjectBodyTo(body string) {
	dc.currentObject.Body = body
}

func (dc *domainContext) iUpdateTheObjectTitleTo(title string) {
	dc.currentObject.Properties["title"] = title
}

func (dc *domainContext) iSaveTheObject() {
	dc.lastErr = dc.vault.SaveObject(dc.currentObject)
}

func (dc *domainContext) theObjectFileShouldContain(expected string) error {
	data, err := os.ReadFile(dc.vault.ObjectPath(dc.currentObject.Type, dc.currentObject.Filename))
	if err != nil {
		return fmt.Errorf("ReadFile error: %v", err)
	}
	if !strings.Contains(string(data), expected) {
		return fmt.Errorf("file does not contain %q", expected)
	}
	return nil
}

func (dc *domainContext) gettingTheObjectByIDShouldReturnBody(expected string) error {
	got, err := dc.vault.GetObject(dc.currentObject.ID)
	if err != nil {
		return fmt.Errorf("GetObject error: %v", err)
	}
	if got.Body != expected {
		return fmt.Errorf("body = %q, want %q", got.Body, expected)
	}
	return nil
}

// ── Object with property/body setup steps ───────────────────────────────────

func (dc *domainContext) aObjectNamedExistsWithPropertySetTo(typeName, name, prop, value string) {
	dc.aObjectNamedExists(typeName, name)
	if err := dc.vault.SetProperty(dc.currentObject.ID, prop, value); err != nil {
		panic(fmt.Sprintf("SetProperty failed: %v", err))
	}
}

func (dc *domainContext) aObjectNamedExistsWithBody(typeName, name, body string) {
	dc.aObjectNamedExists(typeName, name)
	dc.currentObject.Body = body
	if err := dc.vault.saveObjectFile(dc.currentObject); err != nil {
		panic(fmt.Sprintf("saveObjectFile failed: %v", err))
	}
}

// ── Relation steps ──────────────────────────────────────────────────────────

func (dc *domainContext) aVaultIsReadyWithRelationSchemas() {
	dc.aVaultIsReady()

	bookSchema := []byte(`name: book
properties:
  - name: title
    type: string
  - name: author
    type: relation
    target: person
    bidirectional: true
    inverse: books
`)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), "book.yaml"), bookSchema, 0644)

	personSchema := []byte(`name: person
properties:
  - name: name
    type: string
  - name: books
    type: relation
    target: book
    multiple: true
    bidirectional: true
    inverse: author
`)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), "person.yaml"), personSchema, 0644)
}

func (dc *domainContext) iLinkToVia(sourceName, targetName, relation string) {
	source := dc.objects[sourceName]
	target := dc.objects[targetName]
	if source == nil || target == nil {
		dc.lastErr = fmt.Errorf("object %q or %q not found", sourceName, targetName)
		return
	}
	dc.lastErr = dc.vault.LinkObjects(source.ID, relation, target.ID)
}

func (dc *domainContext) iLinkTheFirstBookToTheSecondBookVia(relation string) {
	if dc.prevObject == nil || dc.currentObject == nil {
		dc.lastErr = fmt.Errorf("need at least 2 objects")
		return
	}
	dc.lastErr = dc.vault.LinkObjects(dc.prevObject.ID, relation, dc.currentObject.ID)
}

func (dc *domainContext) thePropertyOfShouldReference(prop, ownerName, targetName string) error {
	owner := dc.objects[ownerName]
	target := dc.objects[targetName]
	if owner == nil || target == nil {
		return fmt.Errorf("object %q or %q not found", ownerName, targetName)
	}
	obj, err := dc.vault.GetObject(owner.ID)
	if err != nil {
		return fmt.Errorf("GetObject error: %v", err)
	}
	if obj.Properties[prop] != target.ID {
		return fmt.Errorf("%s.%s = %v, want %q", ownerName, prop, obj.Properties[prop], target.ID)
	}
	return nil
}

func (dc *domainContext) thePropertyOfShouldContain(prop, ownerName, targetName string) error {
	owner := dc.objects[ownerName]
	target := dc.objects[targetName]
	if owner == nil || target == nil {
		return fmt.Errorf("object %q or %q not found", ownerName, targetName)
	}
	obj, err := dc.vault.GetObject(owner.ID)
	if err != nil {
		return fmt.Errorf("GetObject error: %v", err)
	}
	items, ok := obj.Properties[prop].([]any)
	if !ok {
		return fmt.Errorf("%s.%s type = %T, want []any", ownerName, prop, obj.Properties[prop])
	}
	for _, item := range items {
		if item == target.ID {
			return nil
		}
	}
	return fmt.Errorf("%s.%s does not contain %q", ownerName, prop, target.ID)
}

func (dc *domainContext) thePropertyOfShouldBeEmpty(prop, ownerName string) error {
	owner := dc.objects[ownerName]
	if owner == nil {
		return fmt.Errorf("object %q not found", ownerName)
	}
	obj, err := dc.vault.GetObject(owner.ID)
	if err != nil {
		return fmt.Errorf("GetObject error: %v", err)
	}
	val := obj.Properties[prop]
	if val != nil {
		return fmt.Errorf("%s.%s = %v, want nil", ownerName, prop, val)
	}
	return nil
}

func (dc *domainContext) iUnlinkFromViaWithBothFlag(sourceName, targetName, relation string) {
	source := dc.objects[sourceName]
	target := dc.objects[targetName]
	if source == nil || target == nil {
		dc.lastErr = fmt.Errorf("object %q or %q not found", sourceName, targetName)
		return
	}
	dc.lastErr = dc.vault.UnlinkObjects(source.ID, relation, target.ID, true)
}

func (dc *domainContext) listingRelationsForShouldReturnNEntries(name string, expected int) error {
	obj := dc.objects[name]
	if obj == nil {
		return fmt.Errorf("object %q not found", name)
	}
	rels, err := dc.vault.ListRelations(obj.ID)
	if err != nil {
		return fmt.Errorf("ListRelations error: %v", err)
	}
	if len(rels) != expected {
		return fmt.Errorf("relations count = %d, want %d", len(rels), expected)
	}
	return nil
}

// ── Query steps ─────────────────────────────────────────────────────────────

func (dc *domainContext) iQueryObjectsWithFilter(filter string) {
	results, err := dc.vault.QueryObjects(filter)
	dc.lastErr = err
	dc.queryResults = results
}

func (dc *domainContext) theQueryShouldReturnNResults(expected int) error {
	if len(dc.queryResults) != expected {
		return fmt.Errorf("query results = %d, want %d", len(dc.queryResults), expected)
	}
	return nil
}

func (dc *domainContext) allResultsShouldHaveType(expected string) error {
	for _, obj := range dc.queryResults {
		if obj.Type != expected {
			return fmt.Errorf("result type = %q, want %q", obj.Type, expected)
		}
	}
	return nil
}

func (dc *domainContext) iSearchObjectsFor(keyword string) {
	results, err := dc.vault.SearchObjects(keyword)
	dc.lastErr = err
	dc.searchResults = results
}

func (dc *domainContext) theSearchShouldReturnNResults(expected int) error {
	if len(dc.searchResults) != expected {
		return fmt.Errorf("search results = %d, want %d", len(dc.searchResults), expected)
	}
	return nil
}

// ── Validate steps ──────────────────────────────────────────────────────────

func (dc *domainContext) aTypeSchemaWithAStringProperty(typeName, propName string) {
	schema := fmt.Sprintf("name: %s\nproperties:\n  - name: %s\n    type: string\n", typeName, propName)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), typeName+".yaml"), []byte(schema), 0644)
}

func (dc *domainContext) aTypeSchemaWithASelectPropertyMissingOptions(typeName string) {
	schema := fmt.Sprintf("name: %s\nproperties:\n  - name: status\n    type: select\n", typeName)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), typeName+".yaml"), []byte(schema), 0644)
}

func (dc *domainContext) iValidateAllSchemas() {
	dc.schemaErrors = ValidateAllSchemas(dc.vault)
}

func (dc *domainContext) schemaShouldHaveNoErrors(typeName string) error {
	if errs, ok := dc.schemaErrors[typeName]; ok && len(errs) > 0 {
		return fmt.Errorf("expected no errors for %q, got %v", typeName, errs)
	}
	return nil
}

func (dc *domainContext) schemaShouldHaveErrors(typeName string) error {
	errs, ok := dc.schemaErrors[typeName]
	if !ok || len(errs) == 0 {
		return fmt.Errorf("expected errors for %q, got none", typeName)
	}
	return nil
}

func (dc *domainContext) anOrphanedRelationExists(fromID, toID string) {
	dc.vault.db.Exec("INSERT INTO relations (name, from_id, to_id) VALUES (?, ?, ?)",
		"author", fromID, toID)
	dc.vault.db.Exec("INSERT INTO objects (id, type, filename, properties, body) VALUES (?, ?, ?, ?, ?)",
		fromID, "book", "test-book", "{}", "")
}

func (dc *domainContext) iValidateRelations() {
	dc.relationErrors = ValidateRelations(dc.vault)
}

func (dc *domainContext) thereShouldBeNRelationErrors(expected int) error {
	if len(dc.relationErrors) != expected {
		return fmt.Errorf("relation errors = %d, want %d: %v", len(dc.relationErrors), expected, dc.relationErrors)
	}
	return nil
}

func (dc *domainContext) twoLinkedNotesExist() {
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), "note.yaml"),
		[]byte("name: note\nproperties:\n  - name: title\n    type: string\n"), 0644)

	noteA, _ := dc.vault.NewObject("note", "alpha")
	noteB, _ := dc.vault.NewObject("note", "beta")
	dc.objects["alpha"] = noteA
	dc.objects["beta"] = noteB

	body := fmt.Sprintf("---\ntitle: Alpha\n---\n\nSee [[%s]].\n", noteB.ID)
	os.WriteFile(dc.vault.ObjectPath(noteA.Type, noteA.Filename), []byte(body), 0644)
	dc.vault.SyncIndex()
}

func (dc *domainContext) aNoteWithABrokenWikiLinkExists() {
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), "note.yaml"),
		[]byte("name: note\nproperties:\n  - name: title\n    type: string\n"), 0644)

	note, _ := dc.vault.NewObject("note", "alpha")
	dc.objects["alpha"] = note

	body := "---\ntitle: Alpha\n---\n\nSee [[note/nonexistent-01jjjjjjjjjjjjjjjjjjjjjjjj]].\n"
	os.WriteFile(dc.vault.ObjectPath(note.Type, note.Filename), []byte(body), 0644)
	dc.vault.SyncIndex()
}

func (dc *domainContext) iValidateWikiLinks() {
	dc.wikiLinkErrors = ValidateWikiLinks(dc.vault)
}

func (dc *domainContext) thereShouldBeNoWikiLinkErrors() error {
	if len(dc.wikiLinkErrors) != 0 {
		return fmt.Errorf("expected no wiki-link errors, got %v", dc.wikiLinkErrors)
	}
	return nil
}

func (dc *domainContext) thereShouldBeNWikiLinkErrors(expected int) error {
	if len(dc.wikiLinkErrors) != expected {
		return fmt.Errorf("wiki-link errors = %d, want %d", len(dc.wikiLinkErrors), expected)
	}
	return nil
}

func (dc *domainContext) theErrorShouldMention(substr string) error {
	for _, err := range dc.wikiLinkErrors {
		if strings.Contains(err.Error(), substr) {
			return nil
		}
	}
	return fmt.Errorf("no error mentions %q", substr)
}

// ── Wiki-link steps ─────────────────────────────────────────────────────────

func (dc *domainContext) aVaultIsReadyWithNoteSchemas() {
	dc.aVaultIsReady()
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), "book.yaml"),
		[]byte("name: book\nproperties:\n  - name: title\n    type: string\n"), 0644)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), "person.yaml"),
		[]byte("name: person\nproperties:\n  - name: name\n    type: string\n"), 0644)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), "note.yaml"),
		[]byte("name: note\nproperties:\n  - name: title\n    type: string\n"), 0644)
}

func (dc *domainContext) bodyContainsAWikiLinkTo(sourceName, targetName string) {
	source := dc.objects[sourceName]
	if source == nil {
		panic(fmt.Sprintf("source object %q not found", sourceName))
	}
	// If target is a known object slug, use its ID; otherwise treat as raw ID
	targetID := targetName
	if target, ok := dc.objects[targetName]; ok {
		targetID = target.ID
	}
	body := fmt.Sprintf("---\ntitle: %s\n---\n\nSee [[%s]].\n", sourceName, targetID)
	os.WriteFile(dc.vault.ObjectPath(source.Type, source.Filename), []byte(body), 0644)
}

func (dc *domainContext) bodyContainsAWikiLinkToWithDisplayText(sourceName, targetName, displayText string) {
	source := dc.objects[sourceName]
	target := dc.objects[targetName]
	if source == nil || target == nil {
		panic(fmt.Sprintf("object %q or %q not found", sourceName, targetName))
	}
	body := fmt.Sprintf("---\ntitle: %s\n---\n\nBy [[%s|%s]].\n", sourceName, target.ID, displayText)
	os.WriteFile(dc.vault.ObjectPath(source.Type, source.Filename), []byte(body), 0644)
}

func (dc *domainContext) iSyncTheIndex() {
	_, err := dc.vault.SyncIndex()
	dc.lastErr = err
}

func (dc *domainContext) shouldHaveNWikiLinks(name string, expected int) error {
	obj := dc.objects[name]
	if obj == nil {
		return fmt.Errorf("object %q not found", name)
	}
	links, err := dc.vault.ListWikiLinks(obj.ID)
	if err != nil {
		return fmt.Errorf("ListWikiLinks error: %v", err)
	}
	dc.wikiLinks = links
	if len(links) != expected {
		return fmt.Errorf("wiki-links = %d, want %d", len(links), expected)
	}
	return nil
}

func (dc *domainContext) theWikiLinkTargetShouldBe(targetName string) error {
	target := dc.objects[targetName]
	if target == nil {
		return fmt.Errorf("object %q not found", targetName)
	}
	if len(dc.wikiLinks) == 0 {
		return fmt.Errorf("no wiki-links to check")
	}
	if dc.wikiLinks[0].ToID != target.ID {
		return fmt.Errorf("wiki-link ToID = %q, want %q", dc.wikiLinks[0].ToID, target.ID)
	}
	return nil
}

func (dc *domainContext) shouldHaveNBacklinksFrom(targetName string, expected int, sourceName string) error {
	target := dc.objects[targetName]
	if target == nil {
		return fmt.Errorf("object %q not found", targetName)
	}
	backlinks, err := dc.vault.ListBacklinks(target.ID)
	if err != nil {
		return fmt.Errorf("ListBacklinks error: %v", err)
	}
	if len(backlinks) != expected {
		return fmt.Errorf("backlinks = %d, want %d", len(backlinks), expected)
	}
	source := dc.objects[sourceName]
	if source == nil {
		return fmt.Errorf("source object %q not found", sourceName)
	}
	found := false
	for _, bl := range backlinks {
		if bl.FromID == source.ID {
			found = true
		}
	}
	if !found {
		return fmt.Errorf("no backlink from %q found", sourceName)
	}
	return nil
}

func (dc *domainContext) theWikiLinkShouldHaveAnEmptyResolvedID() error {
	if len(dc.wikiLinks) == 0 {
		return fmt.Errorf("no wiki-links to check")
	}
	if dc.wikiLinks[0].ToID != "" {
		return fmt.Errorf("ToID = %q, want empty", dc.wikiLinks[0].ToID)
	}
	return nil
}

func (dc *domainContext) iChangeWikiLinkTo(sourceName, targetName string) {
	dc.bodyContainsAWikiLinkTo(sourceName, targetName)
}

func (dc *domainContext) wikiLinkShouldPointTo(sourceName, targetName string) error {
	source := dc.objects[sourceName]
	target := dc.objects[targetName]
	if source == nil || target == nil {
		return fmt.Errorf("object %q or %q not found", sourceName, targetName)
	}
	links, err := dc.vault.ListWikiLinks(source.ID)
	if err != nil {
		return fmt.Errorf("ListWikiLinks error: %v", err)
	}
	if len(links) != 1 {
		return fmt.Errorf("wiki-links = %d, want 1", len(links))
	}
	if links[0].ToID != target.ID {
		return fmt.Errorf("wiki-link ToID = %q, want %q", links[0].ToID, target.ID)
	}
	return nil
}

func (dc *domainContext) shouldHaveNBacklinks(targetName string, expected int) error {
	target := dc.objects[targetName]
	if target == nil {
		return fmt.Errorf("object %q not found", targetName)
	}
	backlinks, err := dc.vault.ListBacklinks(target.ID)
	if err != nil {
		return fmt.Errorf("ListBacklinks error: %v", err)
	}
	if len(backlinks) != expected {
		return fmt.Errorf("backlinks = %d, want %d", len(backlinks), expected)
	}
	return nil
}

func (dc *domainContext) theWikiLinkDisplayTextShouldBe(expected string) error {
	if len(dc.wikiLinks) == 0 {
		return fmt.Errorf("no wiki-links to check")
	}
	if dc.wikiLinks[0].DisplayText != expected {
		return fmt.Errorf("DisplayText = %q, want %q", dc.wikiLinks[0].DisplayText, expected)
	}
	return nil
}

// ── Resolve steps ───────────────────────────────────────────────────────────

func (dc *domainContext) iResolveTheObjectByItsFullID() {
	if dc.currentObject == nil {
		dc.lastErr = fmt.Errorf("no current object")
		return
	}
	dc.resolvedID, dc.lastErr = dc.vault.ResolveID(dc.currentObject.ID)
}

func (dc *domainContext) theResolvedIDShouldMatchTheOriginal() error {
	if dc.resolvedID != dc.currentObject.ID {
		return fmt.Errorf("resolved ID = %q, want %q", dc.resolvedID, dc.currentObject.ID)
	}
	return nil
}

func (dc *domainContext) iResolveTheObjectByPrefix(prefix string) {
	obj, err := dc.vault.ResolveObject(prefix)
	dc.lastErr = err
	if err == nil {
		dc.retrieved = obj
	}
}

func (dc *domainContext) theResolvedObjectShouldMatchTheCreatedOne() error {
	if dc.retrieved == nil {
		return fmt.Errorf("no resolved object")
	}
	if dc.currentObject == nil {
		return fmt.Errorf("no current object to compare")
	}
	if dc.retrieved.ID != dc.currentObject.ID {
		return fmt.Errorf("resolved ID = %q, want %q", dc.retrieved.ID, dc.currentObject.ID)
	}
	return nil
}

func (dc *domainContext) iResolveTheObjectByAPartialULIDPrefix() {
	if dc.currentObject == nil {
		dc.lastErr = fmt.Errorf("no current object")
		return
	}
	// Use type + display name + first 4 chars of ULID as partial prefix
	displayName := dc.currentObject.DisplayName()
	ulidPart := strings.TrimPrefix(dc.currentObject.Filename, displayName+"-")
	partial := ulidPart[:4]
	prefix := dc.currentObject.Type + "/" + displayName + "-" + partial
	obj, err := dc.vault.ResolveObject(prefix)
	dc.lastErr = err
	if err == nil {
		dc.retrieved = obj
	}
}

func (dc *domainContext) anAmbiguousMatchErrorShouldOccurWithNCandidates(expected int) error {
	if dc.lastErr == nil {
		return fmt.Errorf("expected AmbiguousMatchError, got nil")
	}
	ambErr, ok := dc.lastErr.(*AmbiguousMatchError)
	if !ok {
		return fmt.Errorf("expected *AmbiguousMatchError, got %T: %v", dc.lastErr, dc.lastErr)
	}
	if len(ambErr.Matches) != expected {
		return fmt.Errorf("candidates = %d, want %d", len(ambErr.Matches), expected)
	}
	return nil
}

// ── Property type steps ─────────────────────────────────────────────────────

func (dc *domainContext) aTypeSchemaWithAll9PropertyTypes() {
	schema := `name: complete
properties:
  - name: title
    type: string
  - name: count
    type: number
  - name: published
    type: date
  - name: due_at
    type: datetime
  - name: homepage
    type: url
  - name: active
    type: checkbox
  - name: status
    type: select
    options:
      - value: draft
      - value: published
  - name: tags
    type: multi_select
    options:
      - value: go
      - value: rust
  - name: author
    type: relation
    target: person
`
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), "complete.yaml"), []byte(schema), 0644)
}

func (dc *domainContext) aTypeSchemaWithAnEnumProperty(typeName string) {
	schema := fmt.Sprintf(`name: %s
properties:
  - name: status
    type: enum
    values:
      - to-read
      - reading
      - done
`, typeName)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), typeName+".yaml"), []byte(schema), 0644)
}

func (dc *domainContext) aTypeSchemaWithADateProperty(typeName string) {
	schema := fmt.Sprintf("name: %s\nproperties:\n  - name: date\n    type: date\n", typeName)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), typeName+".yaml"), []byte(schema), 0644)
}

func (dc *domainContext) aTypeSchemaWithAURLProperty(typeName string) {
	schema := fmt.Sprintf("name: %s\nproperties:\n  - name: link\n    type: url\n", typeName)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), typeName+".yaml"), []byte(schema), 0644)
}

func (dc *domainContext) aTypeSchemaWithASelectStatusProperty(typeName string) {
	schema := fmt.Sprintf(`name: %s
properties:
  - name: status
    type: select
    options:
      - value: to-read
      - value: reading
      - value: done
`, typeName)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), typeName+".yaml"), []byte(schema), 0644)
}

// aObjectNamedExistsWithRawProperty creates an object and writes a property directly
// to the file, bypassing SetProperty validation. This is needed for negative test cases.
func (dc *domainContext) aObjectNamedExistsWithRawProperty(typeName, name, prop, value string) {
	dc.aObjectNamedExists(typeName, name)
	dc.currentObject.Properties[prop] = value
	if err := dc.vault.saveObjectFile(dc.currentObject); err != nil {
		panic(fmt.Sprintf("saveObjectFile failed: %v", err))
	}
	// Re-sync to update DB
	dc.vault.SyncIndex()
}

func (dc *domainContext) iValidateTheObjectAgainstItsSchema() {
	if dc.currentObject == nil {
		dc.lastErr = fmt.Errorf("no current object")
		return
	}
	schema, err := dc.vault.LoadType(dc.currentObject.Type)
	if err != nil {
		dc.lastErr = err
		return
	}
	obj, err := dc.vault.GetObject(dc.currentObject.ID)
	if err != nil {
		dc.lastErr = err
		return
	}
	dc.objectValidationErrors = ValidateObject(obj.Properties, schema)
}

func (dc *domainContext) theObjectShouldHaveNoValidationErrors() error {
	if len(dc.objectValidationErrors) != 0 {
		return fmt.Errorf("expected no validation errors, got %v", dc.objectValidationErrors)
	}
	return nil
}

func (dc *domainContext) theObjectShouldHaveValidationErrors() error {
	if len(dc.objectValidationErrors) == 0 {
		return fmt.Errorf("expected validation errors, got none")
	}
	return nil
}

func (dc *domainContext) iMigrateSchemas() {
	result, err := dc.vault.MigrateSchemas(false)
	dc.lastErr = err
	dc.schemaMigrateResult = result
}

func (dc *domainContext) theSchemaShouldUseSelectInsteadOfEnum(typeName string) error {
	schema, err := dc.vault.LoadType(typeName)
	if err != nil {
		return fmt.Errorf("LoadType(%q) error: %v", typeName, err)
	}
	for _, p := range schema.Properties {
		if p.Type == "enum" {
			return fmt.Errorf("property %q still uses type \"enum\"", p.Name)
		}
	}
	// Verify at least one select property exists (was converted)
	for _, p := range schema.Properties {
		if p.Type == "select" {
			return nil
		}
	}
	return fmt.Errorf("no select property found in schema %q", typeName)
}

// ── Property emoji steps ────────────────────────────────────────────────────

func (dc *domainContext) aTypeSchemaWithPropertyHavingEmoji(typeName, propName, emoji string) {
	schema := fmt.Sprintf("name: %s\nproperties:\n  - name: %s\n    type: string\n    emoji: %s\n", typeName, propName, emoji)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), typeName+".yaml"), []byte(schema), 0644)
}

func (dc *domainContext) aTypeSchemaWithPropertiesHavingUniqueEmojis(typeName string) {
	schema := fmt.Sprintf(`name: %s
properties:
  - name: title
    type: string
    emoji: 📖
  - name: rating
    type: number
    emoji: ⭐
`, typeName)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), typeName+".yaml"), []byte(schema), 0644)
}

func (dc *domainContext) aTypeSchemaWithPropertiesHavingDuplicateEmojis(typeName string) {
	schema := fmt.Sprintf(`name: %s
properties:
  - name: title
    type: string
    emoji: 👤
  - name: author
    type: string
    emoji: 👤
`, typeName)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), typeName+".yaml"), []byte(schema), 0644)
}

func (dc *domainContext) aTypeSchemaWithSomePropertiesMissingEmojis(typeName string) {
	schema := fmt.Sprintf(`name: %s
properties:
  - name: title
    type: string
  - name: author
    type: string
  - name: rating
    type: number
    emoji: ⭐
`, typeName)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), typeName+".yaml"), []byte(schema), 0644)
}

// ── Pinned property steps ───────────────────────────────────────────────────

func (dc *domainContext) aTypeSchemaWithPropertyHavingPin(typeName, propName string, pin int) {
	schema := fmt.Sprintf("name: %s\nproperties:\n  - name: %s\n    type: string\n    pin: %d\n", typeName, propName, pin)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), typeName+".yaml"), []byte(schema), 0644)
}

func (dc *domainContext) aTypeSchemaWithPropertiesHavingUniquePins(typeName string) {
	schema := fmt.Sprintf(`name: %s
properties:
  - name: status
    type: string
    pin: 1
  - name: rating
    type: number
    pin: 2
`, typeName)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), typeName+".yaml"), []byte(schema), 0644)
}

func (dc *domainContext) aTypeSchemaWithPropertiesHavingDuplicatePins(typeName string) {
	schema := fmt.Sprintf(`name: %s
properties:
  - name: status
    type: string
    pin: 1
  - name: rating
    type: number
    pin: 1
`, typeName)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), typeName+".yaml"), []byte(schema), 0644)
}

func (dc *domainContext) aTypeSchemaWithSomePropertiesUnpinned(typeName string) {
	schema := fmt.Sprintf(`name: %s
properties:
  - name: title
    type: string
  - name: author
    type: string
  - name: status
    type: string
    pin: 1
`, typeName)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), typeName+".yaml"), []byte(schema), 0644)
}

// ── Property filtering steps ────────────────────────────────────────────────

func (dc *domainContext) aTypeSchemaWithProperties(typeName, propList string) {
	props := strings.Split(propList, ",")
	var yamlProps string
	for _, p := range props {
		yamlProps += fmt.Sprintf("  - name: %s\n    type: string\n", strings.TrimSpace(p))
	}
	schema := fmt.Sprintf("name: %s\nproperties:\n%s", typeName, yamlProps)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), typeName+".yaml"), []byte(schema), 0644)
}

func (dc *domainContext) aRawObjectFileWithProperties(relPath string, table *godog.Table) {
	var yamlContent string
	for _, row := range table.Rows[1:] { // skip header
		yamlContent += fmt.Sprintf("%s: %s\n", row.Cells[0].Value, row.Cells[1].Value)
	}
	content := fmt.Sprintf("---\n%s---\nSome body content.\n", yamlContent)

	fullPath := filepath.Join(dc.vault.ObjectsDir(), relPath)
	os.MkdirAll(filepath.Dir(fullPath), 0755)
	os.WriteFile(fullPath, []byte(content), 0644)
}

func (dc *domainContext) getIndexedProperties(objectID string) (map[string]any, error) {
	var propsJSON string
	err := dc.vault.db.QueryRow("SELECT properties FROM objects WHERE id = ?", objectID).Scan(&propsJSON)
	if err != nil {
		return nil, fmt.Errorf("query properties for %s: %v", objectID, err)
	}
	var props map[string]any
	if err := json.Unmarshal([]byte(propsJSON), &props); err != nil {
		return nil, fmt.Errorf("unmarshal properties for %s: %v", objectID, err)
	}
	return props, nil
}

func (dc *domainContext) theIndexedPropertiesForShouldContain(objectID, key string) error {
	props, err := dc.getIndexedProperties(objectID)
	if err != nil {
		return err
	}
	if _, ok := props[key]; !ok {
		return fmt.Errorf("indexed properties for %s do not contain %q, got: %v", objectID, key, props)
	}
	return nil
}

func (dc *domainContext) theIndexedPropertiesForShouldNotContain(objectID, key string) error {
	props, err := dc.getIndexedProperties(objectID)
	if err != nil {
		return err
	}
	if _, ok := props[key]; ok {
		return fmt.Errorf("indexed properties for %s should not contain %q, got: %v", objectID, key, props)
	}
	return nil
}

func (dc *domainContext) theIndexedPropertiesForShouldBeEmpty(objectID string) error {
	props, err := dc.getIndexedProperties(objectID)
	if err != nil {
		return err
	}
	if len(props) != 0 {
		return fmt.Errorf("expected empty properties for %s, got: %v", objectID, props)
	}
	return nil
}

func (dc *domainContext) theFileShouldStillContainInFrontmatter(relPath, expected string) error {
	fullPath := filepath.Join(dc.vault.ObjectsDir(), relPath)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return fmt.Errorf("read file %s: %v", relPath, err)
	}
	if !strings.Contains(string(data), expected) {
		return fmt.Errorf("file %s does not contain %q", relPath, expected)
	}
	return nil
}

// ── Name property steps ─────────────────────────────────────────────────────

func (dc *domainContext) iSetTheObjectNameTo(name string) {
	dc.currentObject.Properties["name"] = name
}

func (dc *domainContext) iRemoveTheNamePropertyFromTheObject() {
	delete(dc.currentObject.Properties, "name")
}

func (dc *domainContext) getNameShouldReturn(expected string) error {
	got := dc.currentObject.GetName()
	if got != expected {
		return fmt.Errorf("GetName() = %q, want %q", got, expected)
	}
	return nil
}

func (dc *domainContext) getNameShouldReturnTheDisplayName() error {
	got := dc.currentObject.GetName()
	expected := dc.currentObject.DisplayName()
	if got != expected {
		return fmt.Errorf("GetName() = %q, want DisplayName() = %q", got, expected)
	}
	return nil
}

func (dc *domainContext) theSyncedObjectShouldHaveNameMatchingItsDisplayName() error {
	obj, err := dc.vault.GetObject(dc.currentObject.ID)
	if err != nil {
		return fmt.Errorf("get object: %v", err)
	}
	name, _ := obj.Properties["name"].(string)
	expected := obj.DisplayName()
	if name != expected {
		return fmt.Errorf("synced name = %q, want DisplayName() = %q", name, expected)
	}
	return nil
}

func (dc *domainContext) theSyncedObjectShouldHaveName(expected string) error {
	obj, err := dc.vault.GetObject(dc.currentObject.ID)
	if err != nil {
		return fmt.Errorf("get object: %v", err)
	}
	name, _ := obj.Properties["name"].(string)
	if name != expected {
		return fmt.Errorf("synced name = %q, want %q", name, expected)
	}
	return nil
}

// ── Common steps ────────────────────────────────────────────────────────────

func (dc *domainContext) anErrorShouldOccur() error {
	if dc.lastErr == nil {
		return fmt.Errorf("expected an error, got nil")
	}
	return nil
}

func (dc *domainContext) noErrorShouldOccur() error {
	if dc.lastErr != nil {
		return fmt.Errorf("expected no error, got %v", dc.lastErr)
	}
	return nil
}

// ── Registration ────────────────────────────────────────────────────────────

// ── Shared properties steps ──────────────────────────────────────────────

func (dc *domainContext) aSharedPropertiesFileWithDateAndSelectProperties(prop1, prop2 string) {
	content := fmt.Sprintf(`properties:
  - name: %s
    type: date
    emoji: 📅
  - name: %s
    type: select
    options:
      - value: high
      - value: medium
      - value: low
`, prop1, prop2)
	os.WriteFile(dc.vault.SharedPropertiesPath(), []byte(content), 0644)
}

func (dc *domainContext) anEmptySharedPropertiesFile() {
	os.WriteFile(dc.vault.SharedPropertiesPath(), []byte(""), 0644)
}

func (dc *domainContext) iLoadSharedProperties() {
	dc.sharedProperties, dc.lastErr = dc.vault.LoadSharedProperties()
}

func (dc *domainContext) sharedPropertiesShouldContainNEntries(expected int) error {
	got := len(dc.sharedProperties)
	if got != expected {
		return fmt.Errorf("shared properties count = %d, want %d", got, expected)
	}
	return nil
}

func (dc *domainContext) sharedPropertyShouldHaveType(name, expectedType string) error {
	for _, p := range dc.sharedProperties {
		if p.Name == name {
			if p.Type != expectedType {
				return fmt.Errorf("shared property %q type = %q, want %q", name, p.Type, expectedType)
			}
			return nil
		}
	}
	return fmt.Errorf("shared property %q not found", name)
}

func (dc *domainContext) aSharedPropertiesFileWithDuplicateProperties(name string) {
	content := fmt.Sprintf(`properties:
  - name: %s
    type: date
  - name: %s
    type: string
`, name, name)
	os.WriteFile(dc.vault.SharedPropertiesPath(), []byte(content), 0644)
}

func (dc *domainContext) aSharedPropertiesFileWithAnInvalidPropertyType() {
	content := `properties:
  - name: bad_prop
    type: invalid
`
	os.WriteFile(dc.vault.SharedPropertiesPath(), []byte(content), 0644)
}

func (dc *domainContext) aSharedPropertiesFileWithAPropertyNamedName() {
	content := `properties:
  - name: name
    type: string
`
	os.WriteFile(dc.vault.SharedPropertiesPath(), []byte(content), 0644)
}

func (dc *domainContext) aSharedPropertiesFileWithASelectPropertyMissingOptions() {
	content := `properties:
  - name: status
    type: select
`
	os.WriteFile(dc.vault.SharedPropertiesPath(), []byte(content), 0644)
}

func (dc *domainContext) sharedPropertiesShouldHaveNoErrors() error {
	if errs, ok := dc.schemaErrors["_shared_properties"]; ok && len(errs) > 0 {
		return fmt.Errorf("expected no shared properties errors, got %v", errs)
	}
	return nil
}

func (dc *domainContext) sharedPropertiesShouldHaveErrors() error {
	errs, ok := dc.schemaErrors["_shared_properties"]
	if !ok || len(errs) == 0 {
		return fmt.Errorf("expected shared properties errors, got none")
	}
	return nil
}

func (dc *domainContext) aTypeSchemaWithUse(typeName, useName string) {
	content := fmt.Sprintf(`name: %s
properties:
  - use: %s
`, typeName, useName)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), typeName+".yaml"), []byte(content), 0644)
}

func (dc *domainContext) aTypeSchemaWithUseAndPin(typeName, useName string, pin int) {
	content := fmt.Sprintf(`name: %s
properties:
  - use: %s
    pin: %d
`, typeName, useName, pin)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), typeName+".yaml"), []byte(content), 0644)
}

func (dc *domainContext) aTypeSchemaWithUseAndEmoji(typeName, useName, emoji string) {
	content := fmt.Sprintf(`name: %s
properties:
  - use: %s
    emoji: %s
`, typeName, useName, emoji)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), typeName+".yaml"), []byte(content), 0644)
}

func (dc *domainContext) aTypeSchemaWithUseAndDisallowedTypeOverride(typeName, useName string) {
	content := fmt.Sprintf(`name: %s
properties:
  - use: %s
    type: string
`, typeName, useName)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), typeName+".yaml"), []byte(content), 0644)
}

func (dc *domainContext) aTypeSchemaWithLocalProperty(typeName, propName string) {
	content := fmt.Sprintf(`name: %s
properties:
  - name: %s
    type: string
`, typeName, propName)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), typeName+".yaml"), []byte(content), 0644)
}

func (dc *domainContext) aTypeSchemaWithDuplicateUse(typeName, useName string) {
	content := fmt.Sprintf(`name: %s
properties:
  - use: %s
  - use: %s
`, typeName, useName, useName)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), typeName+".yaml"), []byte(content), 0644)
}

func (dc *domainContext) aTypeSchemaWithBothUseAndNameOnSameEntry(typeName string) {
	content := fmt.Sprintf(`name: %s
properties:
  - use: due_date
    name: my_date
`, typeName)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), typeName+".yaml"), []byte(content), 0644)
}

func (dc *domainContext) iLoadType(typeName string) {
	dc.loadedSchema, dc.lastErr = dc.vault.LoadType(typeName)
}

func (dc *domainContext) theLoadedTypeShouldHaveNProperties(expected int) error {
	if dc.loadedSchema == nil {
		return fmt.Errorf("loaded schema is nil")
	}
	got := len(dc.loadedSchema.Properties)
	if got != expected {
		return fmt.Errorf("loaded type properties = %d, want %d", got, expected)
	}
	return nil
}

func (dc *domainContext) theLoadedPropertyShouldHaveType(propName, expectedType string) error {
	if dc.loadedSchema == nil {
		return fmt.Errorf("loaded schema is nil")
	}
	for _, p := range dc.loadedSchema.Properties {
		if p.Name == propName {
			if p.Type != expectedType {
				return fmt.Errorf("loaded property %q type = %q, want %q", propName, p.Type, expectedType)
			}
			return nil
		}
	}
	return fmt.Errorf("loaded property %q not found", propName)
}

func (dc *domainContext) theLoadedPropertyShouldHaveEmoji(propName, expectedEmoji string) error {
	if dc.loadedSchema == nil {
		return fmt.Errorf("loaded schema is nil")
	}
	for _, p := range dc.loadedSchema.Properties {
		if p.Name == propName {
			if p.Emoji != expectedEmoji {
				return fmt.Errorf("loaded property %q emoji = %q, want %q", propName, p.Emoji, expectedEmoji)
			}
			return nil
		}
	}
	return fmt.Errorf("loaded property %q not found", propName)
}

func (dc *domainContext) theLoadedPropertyShouldHavePin(propName string, expectedPin int) error {
	if dc.loadedSchema == nil {
		return fmt.Errorf("loaded schema is nil")
	}
	for _, p := range dc.loadedSchema.Properties {
		if p.Name == propName {
			if p.Pin != expectedPin {
				return fmt.Errorf("loaded property %q pin = %d, want %d", propName, p.Pin, expectedPin)
			}
			return nil
		}
	}
	return fmt.Errorf("loaded property %q not found", propName)
}

func (dc *domainContext) aTypeSchemaWithMixedUseAndNameProperties(typeName string) {
	content := fmt.Sprintf(`name: %s
properties:
  - name: title
    type: string
  - use: due_date
  - name: budget
    type: number
`, typeName)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), typeName+".yaml"), []byte(content), 0644)
}

func (dc *domainContext) theLoadedPropertyAtIndexShouldBe(index int, expectedName string) error {
	if dc.loadedSchema == nil {
		return fmt.Errorf("loaded schema is nil")
	}
	if index >= len(dc.loadedSchema.Properties) {
		return fmt.Errorf("index %d out of range (have %d properties)", index, len(dc.loadedSchema.Properties))
	}
	got := dc.loadedSchema.Properties[index].Name
	if got != expectedName {
		return fmt.Errorf("property at index %d = %q, want %q", index, got, expectedName)
	}
	return nil
}

// ── System property steps ────────────────────────────────────────────────

func (dc *domainContext) theSystemPropertyRegistryShouldContain(nameList string) error {
	expected := strings.Split(nameList, ", ")
	for i, s := range expected {
		expected[i] = strings.TrimSpace(s)
	}
	got := SystemPropertyNames()
	if len(got) != len(expected) {
		return fmt.Errorf("registry has %d entries, want %d: %v", len(got), len(expected), got)
	}
	for i, name := range expected {
		if got[i] != name {
			return fmt.Errorf("registry[%d] = %q, want %q", i, got[i], name)
		}
	}
	return nil
}

func (dc *domainContext) shouldBeASystemProperty(name string) error {
	if !IsSystemProperty(name) {
		return fmt.Errorf("%q should be a system property", name)
	}
	return nil
}

func (dc *domainContext) shouldNotBeASystemProperty(name string) error {
	if IsSystemProperty(name) {
		return fmt.Errorf("%q should not be a system property", name)
	}
	return nil
}

func (dc *domainContext) aTypeSchemaWithASystemProperty(typeName, propName string) {
	content := fmt.Sprintf(`name: %s
properties:
  - name: %s
    type: datetime
`, typeName, propName)
	os.WriteFile(filepath.Join(dc.vault.TypesDir(), typeName+".yaml"), []byte(content), 0644)
}

func (dc *domainContext) aSharedPropertiesFileWithASystemProperty(propName string) {
	content := fmt.Sprintf(`properties:
  - name: %s
    type: datetime
`, propName)
	os.WriteFile(dc.vault.SharedPropertiesPath(), []byte(content), 0644)
}

func (dc *domainContext) theObjectShouldHaveATimestamp(propName string) error {
	got, err := dc.vault.GetObject(dc.currentObject.ID)
	if err != nil {
		return fmt.Errorf("GetObject error: %v", err)
	}
	val, ok := got.Properties[propName]
	if !ok || val == nil || val == "" {
		return fmt.Errorf("expected %q to be set, got %v", propName, val)
	}
	return nil
}

func (dc *domainContext) theObjectTimestampShouldNotHaveChanged(propName string) error {
	got, err := dc.vault.GetObject(dc.currentObject.ID)
	if err != nil {
		return fmt.Errorf("GetObject error: %v", err)
	}
	val := fmt.Sprintf("%v", got.Properties[propName])
	if dc.createdAtSnapshot == "" {
		return fmt.Errorf("no snapshot for %q", propName)
	}
	if val != dc.createdAtSnapshot {
		return fmt.Errorf("%q changed: was %q, now %q", propName, dc.createdAtSnapshot, val)
	}
	return nil
}

func (dc *domainContext) theObjectTimestampShouldBeRecent(propName string) error {
	got, err := dc.vault.GetObject(dc.currentObject.ID)
	if err != nil {
		return fmt.Errorf("GetObject error: %v", err)
	}
	val, ok := got.Properties[propName]
	if !ok || val == nil || val == "" {
		return fmt.Errorf("expected %q to be set", propName)
	}
	s := fmt.Sprintf("%v", val)
	parsed, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return fmt.Errorf("%q value %q is not valid RFC 3339: %v", propName, s, err)
	}
	if time.Since(parsed) > 5*time.Second {
		return fmt.Errorf("%q value %q is not recent (older than 5s)", propName, s)
	}
	return nil
}

func (dc *domainContext) theFrontmatterShouldHaveSystemPropertiesBeforeSchemaProperties() error {
	pairs := [][2]string{{"name", "created_at"}, {"created_at", "updated_at"}, {"updated_at", "title"}}
	for _, p := range pairs {
		if err := dc.theFrontmatterShouldHaveBefore(p[0], p[1]); err != nil {
			return err
		}
	}
	return nil
}

func (dc *domainContext) theIndexedPropertiesForTheObjectShouldContain(propName string) error {
	var propsJSON string
	err := dc.vault.db.QueryRow("SELECT properties FROM objects WHERE id = ?", dc.currentObject.ID).Scan(&propsJSON)
	if err != nil {
		return fmt.Errorf("query error: %v", err)
	}
	var props map[string]any
	if err := json.Unmarshal([]byte(propsJSON), &props); err != nil {
		return fmt.Errorf("unmarshal error: %v", err)
	}
	if _, ok := props[propName]; !ok {
		return fmt.Errorf("indexed properties do not contain %q: %v", propName, props)
	}
	return nil
}

func (dc *domainContext) createRawObjectFile(prefix, frontmatter string) {
	typeName := "book"
	filename := prefix + mustULID()
	objPath := dc.vault.ObjectPath(typeName, filename)
	os.MkdirAll(filepath.Dir(objPath), 0755)
	os.WriteFile(objPath, []byte("---\n"+frontmatter+"---\n"), 0644)
	dc.currentObject = &Object{
		ID:       typeName + "/" + filename,
		Type:     typeName,
		Filename: filename,
	}
}

func (dc *domainContext) aRawObjectFileWithoutTimestampsExists() {
	dc.createRawObjectFile("legacy-book-", "name: legacy-book\ntitle: Legacy\n")
}

func (dc *domainContext) rawObjectFileShouldNotContain(propName string) error {
	data, err := os.ReadFile(dc.vault.ObjectPath(dc.currentObject.Type, dc.currentObject.Filename))
	if err != nil {
		return fmt.Errorf("ReadFile error: %v", err)
	}
	content := string(data)
	if strings.Contains(content, propName+":") {
		return fmt.Errorf("%s was added to existing object:\n%s", propName, content)
	}
	return nil
}

func (dc *domainContext) theRawObjectFileShouldNotHaveTimestampsAdded() error {
	if err := dc.rawObjectFileShouldNotContain("created_at"); err != nil {
		return err
	}
	return dc.rawObjectFileShouldNotContain("updated_at")
}

func (dc *domainContext) theObjectShouldNotHaveProperty(propName string) error {
	got, err := dc.vault.GetObject(dc.currentObject.ID)
	if err != nil {
		return fmt.Errorf("GetObject error: %v", err)
	}
	if _, ok := got.Properties[propName]; ok {
		return fmt.Errorf("expected object to not have property %q, but it does", propName)
	}
	return nil
}

func (dc *domainContext) aRawObjectFileWithDescriptionExists() {
	dc.createRawObjectFile("desc-raw-book-", "name: desc-raw-book\ndescription: A raw book with description\ntitle: Raw Book\n")
}

func (dc *domainContext) theRawObjectFileShouldNotHaveDescriptionAdded() error {
	return dc.rawObjectFileShouldNotContain("description")
}

func (dc *domainContext) theFrontmatterShouldHaveBefore(first, second string) error {
	data, err := os.ReadFile(dc.vault.ObjectPath(dc.currentObject.Type, dc.currentObject.Filename))
	if err != nil {
		return fmt.Errorf("ReadFile error: %v", err)
	}
	content := string(data)
	firstIdx := strings.Index(content, first+":")
	secondIdx := strings.Index(content, second+":")
	if firstIdx == -1 {
		return fmt.Errorf("%q not found in frontmatter:\n%s", first, content)
	}
	if secondIdx == -1 {
		return fmt.Errorf("%q not found in frontmatter:\n%s", second, content)
	}
	if firstIdx > secondIdx {
		return fmt.Errorf("%q should come before %q in frontmatter", first, second)
	}
	return nil
}

func initDomainSteps(ctx *godog.ScenarioContext) {
	dc := newDomainContext()

	// Cleanup after each scenario
	ctx.After(func(hookCtx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		if dc.vault != nil {
			dc.vault.Close()
		}
		if dc.rootDir != "" {
			os.RemoveAll(dc.rootDir)
		}
		return hookCtx, nil
	})

	// Vault steps
	ctx.Step(`^I initialize a new vault$`, dc.iInitializeANewVault)
	ctx.Step(`^a vault is initialized$`, dc.aVaultIsInitialized)
	ctx.Step(`^the vault directory structure should exist$`, dc.theVaultDirectoryStructureShouldExist)
	ctx.Step(`^the SQLite database should exist$`, dc.theSQLiteDatabaseShouldExist)
	ctx.Step(`^the \.gitignore should contain "([^"]*)"$`, dc.theGitignoreShouldContain)
	ctx.Step(`^I initialize the vault again$`, dc.iInitializeTheVaultAgain)
	ctx.Step(`^I open the vault$`, dc.iOpenTheVault)
	ctx.Step(`^I close the vault$`, dc.iCloseTheVault)
	ctx.Step(`^I open an uninitialized vault$`, dc.iOpenAnUninitializedVault)
	ctx.Step(`^an object file exists on disk at "([^"]*)" with title "([^"]*)"$`, dc.anObjectFileExistsOnDisk)
	ctx.Step(`^the index should contain (\d+) objects?$`, dc.theIndexShouldContainNObjects)

	// Object steps
	ctx.Step(`^a vault is ready$`, dc.aVaultIsReady)
	ctx.Step(`^I create a "([^"]*)" object named "([^"]*)"$`, dc.iCreateAObjectNamed)
	ctx.Step(`^I create another "([^"]*)" object named "([^"]*)"$`, dc.iCreateAnotherObjectNamed)
	ctx.Step(`^the object filename should start with "([^"]*)"$`, dc.theObjectFilenameShouldStartWith)
	ctx.Step(`^the object filename should have a (\d+)-character ULID suffix$`, dc.theObjectFilenameShouldHaveACharacterULIDSuffix)
	ctx.Step(`^the object type should be "([^"]*)"$`, dc.theObjectTypeShouldBe)
	ctx.Step(`^the object file should exist on disk$`, dc.theObjectFileShouldExistOnDisk)
	ctx.Step(`^the two objects should have different IDs$`, dc.theTwoObjectsShouldHaveDifferentIDs)
	ctx.Step(`^a "([^"]*)" object named "([^"]*)" exists$`, dc.aObjectNamedExists)
	ctx.Step(`^I get the object by its ID$`, dc.iGetTheObjectByItsID)
	ctx.Step(`^the retrieved object should match the created one$`, dc.theRetrievedObjectShouldMatchTheCreatedOne)
	ctx.Step(`^I set property "([^"]*)" to "([^"]*)" on the object$`, dc.iSetPropertyToOnTheObject)
	ctx.Step(`^the object property "([^"]*)" should be "([^"]*)"$`, dc.theObjectPropertyShouldBe)
	ctx.Step(`^I update the object body to "([^"]*)"$`, dc.iUpdateTheObjectBodyTo)
	ctx.Step(`^I update the object title to "([^"]*)"$`, dc.iUpdateTheObjectTitleTo)
	ctx.Step(`^I save the object$`, dc.iSaveTheObject)
	ctx.Step(`^the object file should contain "([^"]*)"$`, dc.theObjectFileShouldContain)
	ctx.Step(`^getting the object by ID should return body "([^"]*)"$`, dc.gettingTheObjectByIDShouldReturnBody)
	ctx.Step(`^a "([^"]*)" object named "([^"]*)" exists with property "([^"]*)" set to "([^"]*)"$`, dc.aObjectNamedExistsWithPropertySetTo)
	ctx.Step(`^a "([^"]*)" object named "([^"]*)" exists with body "([^"]*)"$`, dc.aObjectNamedExistsWithBody)

	// Relation steps
	ctx.Step(`^a vault is ready with relation schemas$`, dc.aVaultIsReadyWithRelationSchemas)
	ctx.Step(`^I link "([^"]*)" to "([^"]*)" via "([^"]*)"$`, dc.iLinkToVia)
	ctx.Step(`^I link the first book to the second book via "([^"]*)"$`, dc.iLinkTheFirstBookToTheSecondBookVia)
	ctx.Step(`^the "([^"]*)" property of "([^"]*)" should reference "([^"]*)"$`, dc.thePropertyOfShouldReference)
	ctx.Step(`^the "([^"]*)" property of "([^"]*)" should contain "([^"]*)"$`, dc.thePropertyOfShouldContain)
	ctx.Step(`^the "([^"]*)" property of "([^"]*)" should be empty$`, dc.thePropertyOfShouldBeEmpty)
	ctx.Step(`^I unlink "([^"]*)" from "([^"]*)" via "([^"]*)" with both flag$`, dc.iUnlinkFromViaWithBothFlag)
	ctx.Step(`^listing relations for "([^"]*)" should return (\d+) entries$`, dc.listingRelationsForShouldReturnNEntries)

	// Query steps
	ctx.Step(`^I query objects with filter "([^"]*)"$`, dc.iQueryObjectsWithFilter)
	ctx.Step(`^the query should return (\d+) results?$`, dc.theQueryShouldReturnNResults)
	ctx.Step(`^all results should have type "([^"]*)"$`, dc.allResultsShouldHaveType)
	ctx.Step(`^I search objects for "([^"]*)"$`, dc.iSearchObjectsFor)
	ctx.Step(`^the search should return (\d+) results?$`, dc.theSearchShouldReturnNResults)

	// Validate steps
	ctx.Step(`^a type schema "([^"]*)" with a "([^"]*)" string property$`, dc.aTypeSchemaWithAStringProperty)
	ctx.Step(`^a type schema "([^"]*)" with a select property missing options$`, dc.aTypeSchemaWithASelectPropertyMissingOptions)
	ctx.Step(`^I validate all schemas$`, dc.iValidateAllSchemas)
	ctx.Step(`^schema "([^"]*)" should have no errors$`, dc.schemaShouldHaveNoErrors)
	ctx.Step(`^schema "([^"]*)" should have errors$`, dc.schemaShouldHaveErrors)
	ctx.Step(`^an orphaned relation from "([^"]*)" to "([^"]*)" exists$`, dc.anOrphanedRelationExists)
	ctx.Step(`^I validate relations$`, dc.iValidateRelations)
	ctx.Step(`^there should be (\d+) relation errors?$`, dc.thereShouldBeNRelationErrors)
	ctx.Step(`^two linked notes exist$`, dc.twoLinkedNotesExist)
	ctx.Step(`^a note with a broken wiki-link exists$`, dc.aNoteWithABrokenWikiLinkExists)
	ctx.Step(`^I validate wiki-links$`, dc.iValidateWikiLinks)
	ctx.Step(`^there should be no wiki-link errors$`, dc.thereShouldBeNoWikiLinkErrors)
	ctx.Step(`^there should be (\d+) wiki-link errors?$`, dc.thereShouldBeNWikiLinkErrors)
	ctx.Step(`^the error should mention "([^"]*)"$`, dc.theErrorShouldMention)

	// Wiki-link steps
	ctx.Step(`^a vault is ready with note schemas$`, dc.aVaultIsReadyWithNoteSchemas)
	ctx.Step(`^"([^"]*)" body contains a wiki-link to "([^"]*)"$`, dc.bodyContainsAWikiLinkTo)
	ctx.Step(`^"([^"]*)" body contains a wiki-link to "([^"]*)" with display text "([^"]*)"$`, dc.bodyContainsAWikiLinkToWithDisplayText)
	ctx.Step(`^I sync the index$`, dc.iSyncTheIndex)
	ctx.Step(`^"([^"]*)" should have (\d+) wiki-links?$`, dc.shouldHaveNWikiLinks)
	ctx.Step(`^the wiki-link target should be "([^"]*)"$`, dc.theWikiLinkTargetShouldBe)
	ctx.Step(`^"([^"]*)" should have (\d+) backlinks? from "([^"]*)"$`, dc.shouldHaveNBacklinksFrom)
	ctx.Step(`^the wiki-link should have an empty resolved ID$`, dc.theWikiLinkShouldHaveAnEmptyResolvedID)
	ctx.Step(`^I change "([^"]*)" wiki-link to "([^"]*)"$`, dc.iChangeWikiLinkTo)
	ctx.Step(`^"([^"]*)" wiki-link should point to "([^"]*)"$`, dc.wikiLinkShouldPointTo)
	ctx.Step(`^"([^"]*)" should have (\d+) backlinks$`, dc.shouldHaveNBacklinks)
	ctx.Step(`^the wiki-link display text should be "([^"]*)"$`, dc.theWikiLinkDisplayTextShouldBe)

	// Resolve steps
	ctx.Step(`^I resolve the object by its full ID$`, dc.iResolveTheObjectByItsFullID)
	ctx.Step(`^the resolved ID should match the original$`, dc.theResolvedIDShouldMatchTheOriginal)
	ctx.Step(`^I resolve the object by prefix "([^"]*)"$`, dc.iResolveTheObjectByPrefix)
	ctx.Step(`^the resolved object should match the created one$`, dc.theResolvedObjectShouldMatchTheCreatedOne)
	ctx.Step(`^I resolve the object by a partial ULID prefix$`, dc.iResolveTheObjectByAPartialULIDPrefix)
	ctx.Step(`^an ambiguous match error should occur with (\d+) candidates$`, dc.anAmbiguousMatchErrorShouldOccurWithNCandidates)

	// Property type steps
	ctx.Step(`^an? "([^"]*)" object named "([^"]*)" exists with raw property "([^"]*)" set to "([^"]*)"$`, dc.aObjectNamedExistsWithRawProperty)
	ctx.Step(`^a type schema with all 9 property types$`, dc.aTypeSchemaWithAll9PropertyTypes)
	ctx.Step(`^a type schema "([^"]*)" with an enum property$`, dc.aTypeSchemaWithAnEnumProperty)
	ctx.Step(`^a type schema "([^"]*)" with a date property$`, dc.aTypeSchemaWithADateProperty)
	ctx.Step(`^a type schema "([^"]*)" with a url property$`, dc.aTypeSchemaWithAURLProperty)
	ctx.Step(`^a type schema "([^"]*)" with a select status property$`, dc.aTypeSchemaWithASelectStatusProperty)
	ctx.Step(`^I validate the object against its schema$`, dc.iValidateTheObjectAgainstItsSchema)
	ctx.Step(`^the object should have no validation errors$`, dc.theObjectShouldHaveNoValidationErrors)
	ctx.Step(`^the object should have validation errors$`, dc.theObjectShouldHaveValidationErrors)
	ctx.Step(`^I migrate schemas$`, dc.iMigrateSchemas)
	ctx.Step(`^the "([^"]*)" schema should use select instead of enum$`, dc.theSchemaShouldUseSelectInsteadOfEnum)

	// Property emoji steps
	ctx.Step(`^a type schema "([^"]*)" with property "([^"]*)" having emoji "([^"]*)"$`, dc.aTypeSchemaWithPropertyHavingEmoji)
	ctx.Step(`^a type schema "([^"]*)" with properties having unique emojis$`, dc.aTypeSchemaWithPropertiesHavingUniqueEmojis)
	ctx.Step(`^a type schema "([^"]*)" with properties having duplicate emojis$`, dc.aTypeSchemaWithPropertiesHavingDuplicateEmojis)
	ctx.Step(`^a type schema "([^"]*)" with some properties missing emojis$`, dc.aTypeSchemaWithSomePropertiesMissingEmojis)

	// Pinned property steps
	ctx.Step(`^a type schema "([^"]*)" with property "([^"]*)" having pin (-?\d+)$`, dc.aTypeSchemaWithPropertyHavingPin)
	ctx.Step(`^a type schema "([^"]*)" with properties having unique pins$`, dc.aTypeSchemaWithPropertiesHavingUniquePins)
	ctx.Step(`^a type schema "([^"]*)" with properties having duplicate pins$`, dc.aTypeSchemaWithPropertiesHavingDuplicatePins)
	ctx.Step(`^a type schema "([^"]*)" with some properties unpinned$`, dc.aTypeSchemaWithSomePropertiesUnpinned)

	// Property filtering steps
	ctx.Step(`^a type schema "([^"]*)" with properties "([^"]*)"$`, dc.aTypeSchemaWithProperties)
	ctx.Step(`^a raw object file "([^"]*)" with properties:$`, dc.aRawObjectFileWithProperties)
	ctx.Step(`^the indexed properties for "([^"]*)" should contain "([^"]*)"$`, dc.theIndexedPropertiesForShouldContain)
	ctx.Step(`^the indexed properties for "([^"]*)" should not contain "([^"]*)"$`, dc.theIndexedPropertiesForShouldNotContain)
	ctx.Step(`^the indexed properties for "([^"]*)" should be empty$`, dc.theIndexedPropertiesForShouldBeEmpty)
	ctx.Step(`^the file "([^"]*)" should still contain "([^"]*)" in frontmatter$`, dc.theFileShouldStillContainInFrontmatter)

	// Name property steps
	ctx.Step(`^I set the object name to "([^"]*)"$`, dc.iSetTheObjectNameTo)
	ctx.Step(`^I remove the name property from the object$`, dc.iRemoveTheNamePropertyFromTheObject)
	ctx.Step(`^GetName should return "([^"]*)"$`, dc.getNameShouldReturn)
	ctx.Step(`^GetName should return the DisplayName$`, dc.getNameShouldReturnTheDisplayName)
	ctx.Step(`^the synced object should have name matching its DisplayName$`, dc.theSyncedObjectShouldHaveNameMatchingItsDisplayName)
	ctx.Step(`^the synced object should have name "([^"]*)"$`, dc.theSyncedObjectShouldHaveName)

	// Shared properties steps
	ctx.Step(`^a shared properties file with "([^"]*)" date and "([^"]*)" select properties$`, dc.aSharedPropertiesFileWithDateAndSelectProperties)
	ctx.Step(`^an empty shared properties file$`, dc.anEmptySharedPropertiesFile)
	ctx.Step(`^I load shared properties$`, dc.iLoadSharedProperties)
	ctx.Step(`^shared properties should contain (\d+) entries?$`, dc.sharedPropertiesShouldContainNEntries)
	ctx.Step(`^shared property "([^"]*)" should have type "([^"]*)"$`, dc.sharedPropertyShouldHaveType)
	ctx.Step(`^a shared properties file with duplicate "([^"]*)" properties$`, dc.aSharedPropertiesFileWithDuplicateProperties)
	ctx.Step(`^a shared properties file with an invalid property type$`, dc.aSharedPropertiesFileWithAnInvalidPropertyType)
	ctx.Step(`^a shared properties file with a property named "([^"]*)"$`, dc.aSharedPropertiesFileWithAPropertyNamedName)
	ctx.Step(`^a shared properties file with a select property missing options$`, dc.aSharedPropertiesFileWithASelectPropertyMissingOptions)
	ctx.Step(`^shared properties should have no errors$`, dc.sharedPropertiesShouldHaveNoErrors)
	ctx.Step(`^shared properties should have errors$`, dc.sharedPropertiesShouldHaveErrors)
	ctx.Step(`^a type schema "([^"]*)" with use "([^"]*)"$`, dc.aTypeSchemaWithUse)
	ctx.Step(`^a type schema "([^"]*)" with use "([^"]*)" and pin (\d+)$`, dc.aTypeSchemaWithUseAndPin)
	ctx.Step(`^a type schema "([^"]*)" with use "([^"]*)" and emoji "([^"]*)"$`, dc.aTypeSchemaWithUseAndEmoji)
	ctx.Step(`^a type schema "([^"]*)" with use "([^"]*)" and disallowed type override$`, dc.aTypeSchemaWithUseAndDisallowedTypeOverride)
	ctx.Step(`^a type schema "([^"]*)" with local property "([^"]*)"$`, dc.aTypeSchemaWithLocalProperty)
	ctx.Step(`^a type schema "([^"]*)" with duplicate use "([^"]*)"$`, dc.aTypeSchemaWithDuplicateUse)
	ctx.Step(`^a type schema "([^"]*)" with both use and name on same entry$`, dc.aTypeSchemaWithBothUseAndNameOnSameEntry)
	ctx.Step(`^I load type "([^"]*)"$`, dc.iLoadType)
	ctx.Step(`^the loaded type should have (\d+) propert(?:y|ies)$`, dc.theLoadedTypeShouldHaveNProperties)
	ctx.Step(`^the loaded property "([^"]*)" should have type "([^"]*)"$`, dc.theLoadedPropertyShouldHaveType)
	ctx.Step(`^the loaded property "([^"]*)" should have emoji "([^"]*)"$`, dc.theLoadedPropertyShouldHaveEmoji)
	ctx.Step(`^the loaded property "([^"]*)" should have pin (\d+)$`, dc.theLoadedPropertyShouldHavePin)
	ctx.Step(`^a type schema "([^"]*)" with mixed use and name properties$`, dc.aTypeSchemaWithMixedUseAndNameProperties)
	ctx.Step(`^the loaded property at index (\d+) should be "([^"]*)"$`, dc.theLoadedPropertyAtIndexShouldBe)

	// System property steps
	ctx.Step(`^the system property registry should contain "([^"]*)"$`, dc.theSystemPropertyRegistryShouldContain)
	ctx.Step(`^"([^"]*)" should be a system property$`, dc.shouldBeASystemProperty)
	ctx.Step(`^"([^"]*)" should not be a system property$`, dc.shouldNotBeASystemProperty)
	ctx.Step(`^a type schema "([^"]*)" with a system property "([^"]*)"$`, dc.aTypeSchemaWithASystemProperty)
	ctx.Step(`^a shared properties file with a system property "([^"]*)"$`, dc.aSharedPropertiesFileWithASystemProperty)
	ctx.Step(`^the object should have an? "([^"]*)" timestamp$`, dc.theObjectShouldHaveATimestamp)
	ctx.Step(`^the object "([^"]*)" should not have changed$`, dc.theObjectTimestampShouldNotHaveChanged)
	ctx.Step(`^the object "([^"]*)" should be recent$`, dc.theObjectTimestampShouldBeRecent)
	ctx.Step(`^the frontmatter should have system properties before schema properties$`, dc.theFrontmatterShouldHaveSystemPropertiesBeforeSchemaProperties)
	ctx.Step(`^the indexed properties for the object should contain "([^"]*)"$`, dc.theIndexedPropertiesForTheObjectShouldContain)
	ctx.Step(`^a raw object file without timestamps exists$`, dc.aRawObjectFileWithoutTimestampsExists)
	ctx.Step(`^the raw object file should not have timestamps added$`, dc.theRawObjectFileShouldNotHaveTimestampsAdded)
	ctx.Step(`^the object should not have property "([^"]*)"$`, dc.theObjectShouldNotHaveProperty)
	ctx.Step(`^a raw object file with description exists$`, dc.aRawObjectFileWithDescriptionExists)
	ctx.Step(`^the raw object file should not have description added$`, dc.theRawObjectFileShouldNotHaveDescriptionAdded)
	ctx.Step(`^the frontmatter should have "([^"]*)" before "([^"]*)"$`, dc.theFrontmatterShouldHaveBefore)

	// Common steps
	ctx.Step(`^an error should occur$`, dc.anErrorShouldOccur)
	ctx.Step(`^no error should occur$`, dc.noErrorShouldOccur)
}
