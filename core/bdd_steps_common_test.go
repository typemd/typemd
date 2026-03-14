package core

import (
	"fmt"

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

	// name uniqueness validation results
	nameUniquenessErrors []error

	// template results
	templateNames  []string
	loadedTemplate *Template
}

func newDomainContext() *domainContext {
	return &domainContext{
		objects: make(map[string]*Object),
	}
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

func initCommonSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	ctx.Step(`^an error should occur$`, dc.anErrorShouldOccur)
	ctx.Step(`^no error should occur$`, dc.noErrorShouldOccur)
}
