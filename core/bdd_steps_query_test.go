package core

import (
	"fmt"
	"strings"

	"github.com/cucumber/godog"
)

// parseTestFilter converts a BDD filter string like "type=book status=reading"
// into []FilterRule for the structured query API.
func parseTestFilter(filter string) []FilterRule {
	filter = strings.TrimSpace(filter)
	if filter == "" {
		return nil
	}
	parts := strings.Fields(filter)
	rules := make([]FilterRule, 0, len(parts))
	for _, part := range parts {
		idx := strings.Index(part, "=")
		if idx < 1 {
			continue
		}
		rules = append(rules, FilterRule{
			Property: part[:idx],
			Operator: "is",
			Value:    part[idx+1:],
		})
	}
	return rules
}

// ── Query steps ─────────────────────────────────────────────────────────────

func (dc *domainContext) iQueryObjectsWithFilter(filter string) {
	results, err := dc.vault.QueryObjects(parseTestFilter(filter))
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

func initQuerySteps(ctx *godog.ScenarioContext, dc *domainContext) {
	ctx.Step(`^I query objects with filter "([^"]*)"$`, dc.iQueryObjectsWithFilter)
	ctx.Step(`^the query should return (\d+) results?$`, dc.theQueryShouldReturnNResults)
	ctx.Step(`^all results should have type "([^"]*)"$`, dc.allResultsShouldHaveType)
	ctx.Step(`^I search objects for "([^"]*)"$`, dc.iSearchObjectsFor)
	ctx.Step(`^the search should return (\d+) results?$`, dc.theSearchShouldReturnNResults)
}
