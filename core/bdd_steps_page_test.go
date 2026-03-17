package core

import (
	"fmt"
	"slices"

	"github.com/cucumber/godog"
)

// ── Page type steps ─────────────────────────────────────────────────────────

func (dc *domainContext) iListAllTypes() {
	dc.typeList = dc.vault.ListTypes()
}

func (dc *domainContext) theTypeListShouldContain(name string) error {
	if !slices.Contains(dc.typeList, name) {
		return fmt.Errorf("expected type list to contain %q, got %v", name, dc.typeList)
	}
	return nil
}

// ── Init ────────────────────────────────────────────────────────────────────

func initPageSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	// When
	ctx.Step(`^I list all types$`, dc.iListAllTypes)

	// Then
	ctx.Step(`^the type list should contain "([^"]*)"$`, dc.theTypeListShouldContain)
}
