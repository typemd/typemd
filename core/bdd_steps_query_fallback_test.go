package core

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cucumber/godog"
)

// errorIndex is an ObjectIndex that always returns errors for read operations.
// Used to simulate SQLite unavailability in fallback tests.
type errorIndex struct{}

func (e *errorIndex) Query([]FilterRule, ...SortRule) ([]*ObjectResult, error) {
	return nil, fmt.Errorf("database is unavailable")
}
func (e *errorIndex) Search(string) ([]*ObjectResult, error) {
	return nil, fmt.Errorf("database is unavailable")
}
func (e *errorIndex) FindRelations(string) ([]Relation, error) {
	return nil, fmt.Errorf("database is unavailable")
}
func (e *errorIndex) FindBacklinks(string) ([]StoredWikiLink, error) {
	return nil, fmt.Errorf("database is unavailable")
}
func (e *errorIndex) ListWikiLinks(string) ([]StoredWikiLink, error) {
	return nil, fmt.Errorf("database is unavailable")
}
func (e *errorIndex) Upsert(string, string, string, string, string) error { return nil }
func (e *errorIndex) Remove(string) error                                 { return nil }
func (e *errorIndex) ListIDs() ([]string, error)                          { return nil, nil }
func (e *errorIndex) InsertRelation(string, string, string) error         { return nil }
func (e *errorIndex) DeleteRelation(string, string, string) error         { return nil }
func (e *errorIndex) DeleteRelationsByName(string) error                  { return nil }
func (e *errorIndex) DeleteNonTagRelations() error                        { return nil }
func (e *errorIndex) DeleteRelationsByObject(string) error                { return nil }
func (e *errorIndex) FindOrphanedRelations() ([]OrphanedRelation, error)  { return nil, nil }
func (e *errorIndex) CleanOrphanedRelations() ([]OrphanedRelation, error) { return nil, nil }
func (e *errorIndex) SyncWikiLinks(string, []WikiLinkEntry) error         { return nil }
func (e *errorIndex) DeleteWikiLinks(string) error                        { return nil }
func (e *errorIndex) Rebuild() error                                      { return nil }
func (e *errorIndex) EnsureSchema() error                                 { return nil }

// fallbackContext holds state for query fallback BDD scenarios.
type fallbackContext struct {
	rootDir       string
	repo          ObjectRepository
	qs            *QueryService
	queryResults  []*Object
	searchResults []*Object
	lastErr       error
}

func newFallbackContext() *fallbackContext {
	return &fallbackContext{}
}

func (fc *fallbackContext) aVaultWithObjectsAndBrokenIndex() {
	fc.rootDir = filepath.Join(os.TempDir(), "typemd-fallback-"+mustULID())
	os.MkdirAll(fc.rootDir, 0755)

	// Create type schema for "book"
	typesDir := filepath.Join(fc.rootDir, ".typemd", "types", "book")
	os.MkdirAll(typesDir, 0755)
	schemaContent := `emoji: "📚"
plural: books
properties:
  - name: status
    type: string
  - name: rating
    type: number
`
	os.WriteFile(filepath.Join(typesDir, "schema.yaml"), []byte(schemaContent), 0644)

	// Create type schema for "article"
	articleDir := filepath.Join(fc.rootDir, ".typemd", "types", "article")
	os.MkdirAll(articleDir, 0755)
	articleSchema := `emoji: "📰"
plural: articles
properties:
  - name: topic
    type: string
`
	os.WriteFile(filepath.Join(articleDir, "schema.yaml"), []byte(articleSchema), 0644)

	// Create object files
	bookDir := filepath.Join(fc.rootDir, "objects", "book")
	os.MkdirAll(bookDir, 0755)

	book1 := `---
name: Alpha Book
status: reading
rating: 8
---
Some interesting content here.
`
	os.WriteFile(filepath.Join(bookDir, "alpha-book-01jtest00000000000000001.md"), []byte(book1), 0644)

	book2 := `---
name: Beta Book
status: finished
rating: 5
---
Another book body.
`
	os.WriteFile(filepath.Join(bookDir, "beta-book-01jtest00000000000000002.md"), []byte(book2), 0644)

	articleDir2 := filepath.Join(fc.rootDir, "objects", "article")
	os.MkdirAll(articleDir2, 0755)

	article1 := `---
name: Gamma Article
topic: golang
---
Article about Go programming.
`
	os.WriteFile(filepath.Join(articleDir2, "gamma-article-01jtest00000000000000003.md"), []byte(article1), 0644)

	// Create repo and QueryService with errorIndex
	fc.repo = NewLocalObjectRepository(fc.rootDir)
	fc.qs = NewQueryService(fc.repo, &errorIndex{})
}

func (fc *fallbackContext) iQueryWithFallbackFilterSortedBy(filter, prop, dir string) {
	rules := parseTestFilter(filter)
	sort := []SortRule{{Property: prop, Direction: dir}}
	results, err := fc.qs.Query(rules, QuerySort(sort...))
	fc.lastErr = err
	fc.queryResults = results
}

func (fc *fallbackContext) theFallbackQueryShouldReturnNResults(expected int) error {
	if fc.lastErr != nil {
		return fmt.Errorf("unexpected error: %v", fc.lastErr)
	}
	if len(fc.queryResults) != expected {
		return fmt.Errorf("fallback query results = %d, want %d", len(fc.queryResults), expected)
	}
	return nil
}

func (fc *fallbackContext) allFallbackResultsShouldHaveType(expected string) error {
	for _, obj := range fc.queryResults {
		if obj.Type != expected {
			return fmt.Errorf("result type = %q, want %q", obj.Type, expected)
		}
	}
	return nil
}

func (fc *fallbackContext) theFirstFallbackResultNameShouldBe(expected string) error {
	if len(fc.queryResults) == 0 {
		return fmt.Errorf("no results")
	}
	name := fc.queryResults[0].GetName()
	if name != expected {
		return fmt.Errorf("first result name = %q, want %q", name, expected)
	}
	return nil
}

func (fc *fallbackContext) iSearchWithFallbackFor(keyword string) {
	results, err := fc.qs.Search(keyword)
	fc.lastErr = err
	fc.searchResults = results
}

func (fc *fallbackContext) theFallbackSearchShouldReturnNResults(expected int) error {
	if fc.lastErr != nil {
		return fmt.Errorf("unexpected error: %v", fc.lastErr)
	}
	if len(fc.searchResults) != expected {
		names := make([]string, len(fc.searchResults))
		for i, obj := range fc.searchResults {
			names[i] = obj.GetName()
		}
		return fmt.Errorf("fallback search results = %d (%s), want %d", len(fc.searchResults), strings.Join(names, ", "), expected)
	}
	return nil
}

func initQueryFallbackSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	fc := newFallbackContext()

	ctx.After(func(hookCtx context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
		if fc.rootDir != "" {
			os.RemoveAll(fc.rootDir)
		}
		return hookCtx, nil
	})

	// Given
	ctx.Step(`^a vault with objects and a broken index$`, fc.aVaultWithObjectsAndBrokenIndex)

	// When - Query
	ctx.Step(`^I query with fallback filter "([^"]*)" sorted by "([^"]*)" "([^"]*)"$`, fc.iQueryWithFallbackFilterSortedBy)

	// Then - Query
	ctx.Step(`^the fallback query should return (\d+) results?$`, fc.theFallbackQueryShouldReturnNResults)
	ctx.Step(`^all fallback results should have type "([^"]*)"$`, fc.allFallbackResultsShouldHaveType)
	ctx.Step(`^the first fallback result name should be "([^"]*)"$`, fc.theFirstFallbackResultNameShouldBe)

	// When - Search
	ctx.Step(`^I search with fallback for "([^"]*)"$`, fc.iSearchWithFallbackFor)

	// Then - Search
	ctx.Step(`^the fallback search should return (\d+) results?$`, fc.theFallbackSearchShouldReturnNResults)
}
