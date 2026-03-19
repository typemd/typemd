package core

import (
	"fmt"
	"strings"

	"github.com/cucumber/godog"
)

// ── Query sort step state ────────────────────────────────────────────────────

type querySortContext struct {
	dc            *domainContext
	sortedResults []*Object
}

func newQuerySortContext(dc *domainContext) *querySortContext {
	return &querySortContext{dc: dc}
}

// ── When steps ──────────────────────────────────────────────────────────────

func (qs *querySortContext) iQueryWithFilterSortedBy(filter, property, direction string) {
	results, err := qs.dc.vault.Queries.Query(parseTestFilter(filter), SortRule{
		Property:  property,
		Direction: direction,
	})
	qs.sortedResults = results
	qs.dc.lastErr = err
}

func (qs *querySortContext) iQueryWithFilterAndNoSort(filter string) {
	results, err := qs.dc.vault.Queries.Query(parseTestFilter(filter))
	qs.sortedResults = results
	qs.dc.lastErr = err
}

// ── Then steps ──────────────────────────────────────────────────────────────

func (qs *querySortContext) theSortedResultsShouldHaveNObjects(n int) error {
	if len(qs.sortedResults) != n {
		return fmt.Errorf("expected %d sorted results, got %d", n, len(qs.sortedResults))
	}
	return nil
}

func (qs *querySortContext) theFirstSortedResultNameShouldComeBeforeTheSecondAlphabetically() error {
	if len(qs.sortedResults) < 2 {
		return fmt.Errorf("need at least 2 results, got %d", len(qs.sortedResults))
	}
	name1 := qs.sortedResults[0].GetName()
	name2 := qs.sortedResults[1].GetName()
	if strings.Compare(name1, name2) >= 0 {
		return fmt.Errorf("expected %q < %q alphabetically", name1, name2)
	}
	return nil
}

// ── Init ────────────────────────────────────────────────────────────────────

func initQuerySortSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	qs := newQuerySortContext(dc)

	// When
	ctx.Step(`^I query with filter "([^"]*)" sorted by "([^"]*)" "([^"]*)"$`, qs.iQueryWithFilterSortedBy)
	ctx.Step(`^I query with filter "([^"]*)" and no sort$`, qs.iQueryWithFilterAndNoSort)

	// Then
	ctx.Step(`^the sorted results should have (\d+) objects?$`, qs.theSortedResultsShouldHaveNObjects)
	ctx.Step(`^the first sorted result name should come before the second alphabetically$`, qs.theFirstSortedResultNameShouldComeBeforeTheSecondAlphabetically)
}
