package core

import (
	"fmt"
	"strings"

	"github.com/cucumber/godog"
)

// ── Filter match step state ────────────────────────────────────────────────

type filterMatchContext struct {
	dc      *domainContext
	obj     *Object
	matched bool
}

func newFilterMatchContext(dc *domainContext) *filterMatchContext {
	return &filterMatchContext{dc: dc}
}

// ── Given steps ────────────────────────────────────────────────────────────

func (fm *filterMatchContext) anObjectWithTypeAndPropertySetTo(typeName, property, value string) {
	fm.obj = &Object{
		Type:     typeName,
		Properties: map[string]any{
			property: value,
		},
	}
}

func (fm *filterMatchContext) anObjectWithTypeAndNoProperty(typeName, property string) {
	fm.obj = &Object{
		Type:       typeName,
		Properties: map[string]any{},
	}
}

func (fm *filterMatchContext) theObjectAlsoHasPropertySetTo(property, value string) {
	if fm.obj.Properties == nil {
		fm.obj.Properties = make(map[string]any)
	}
	fm.obj.Properties[property] = value
}

// ── When steps ─────────────────────────────────────────────────────────────

func (fm *filterMatchContext) iMatchFilterPropertyOperatorValue(property, operator, value string) {
	rule := FilterRule{Property: property, Operator: operator, Value: value}
	fm.matched = MatchFilter(fm.obj, rule)
}

func (fm *filterMatchContext) iMatchFilters(filtersStr string) {
	// Parse "type=book status=reading" into FilterRules with "is" operator
	parts := strings.Fields(filtersStr)
	var rules []FilterRule
	for _, part := range parts {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			continue
		}
		rules = append(rules, FilterRule{
			Property: kv[0],
			Operator: "is",
			Value:    kv[1],
		})
	}
	fm.matched = MatchFilters(fm.obj, rules)
}

// ── Then steps ─────────────────────────────────────────────────────────────

func (fm *filterMatchContext) theFilterShouldMatch() error {
	if !fm.matched {
		return fmt.Errorf("expected filter to match, but it did not")
	}
	return nil
}

func (fm *filterMatchContext) theFilterShouldNotMatch() error {
	if fm.matched {
		return fmt.Errorf("expected filter not to match, but it did")
	}
	return nil
}

// ── Init ───────────────────────────────────────────────────────────────────

func initFilterMatchSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	fm := newFilterMatchContext(dc)

	// Given
	ctx.Step(`^an object with type "([^"]*)" and property "([^"]*)" set to "([^"]*)"$`, fm.anObjectWithTypeAndPropertySetTo)
	ctx.Step(`^an object with type "([^"]*)" and no property "([^"]*)"$`, fm.anObjectWithTypeAndNoProperty)
	ctx.Step(`^the object also has property "([^"]*)" set to "([^"]*)"$`, fm.theObjectAlsoHasPropertySetTo)

	// When
	ctx.Step(`^I match filter property "([^"]*)" operator "([^"]*)" value "([^"]*)"$`, fm.iMatchFilterPropertyOperatorValue)
	ctx.Step(`^I match filters "([^"]*)"$`, fm.iMatchFilters)

	// Then
	ctx.Step(`^the filter should match$`, fm.theFilterShouldMatch)
	ctx.Step(`^the filter should not match$`, fm.theFilterShouldNotMatch)
}
