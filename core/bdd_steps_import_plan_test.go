package core

import (
	"fmt"

	"github.com/cucumber/godog"
)

func (dc *domainContext) aClassificationList(table *godog.Table) {
	headers := table.Rows[0].Cells
	dc.importClassifications = nil
	for _, row := range table.Rows[1:] {
		vals := make(map[string]string)
		for i, cell := range row.Cells {
			vals[headers[i].Value] = cell.Value
		}
		dc.importClassifications = append(dc.importClassifications, ObjectPlan{
			SourcePath: vals["source_path"],
			TypeName:   vals["type_name"],
			Name:       vals["name"],
			Conflict:   "none",
		})
	}
}

func (dc *domainContext) aClassificationListWithCircularDependencies() {
	dc.importClassifications = []ObjectPlan{
		{SourcePath: "a.md", TypeName: "page", Name: "A", DependsOn: []int{1}},
		{SourcePath: "b.md", TypeName: "page", Name: "B", DependsOn: []int{0}},
	}
}

func (dc *domainContext) iGenerateAnImportPlan() {
	dc.importPlan, dc.lastErr = dc.vault.GeneratePlan(dc.importClassifications)
}

func (dc *domainContext) thePlanShouldHaveNNewTypes(count int) error {
	if dc.importPlan == nil {
		return fmt.Errorf("import plan is nil")
	}
	if len(dc.importPlan.Types) != count {
		return fmt.Errorf("expected %d new types, got %d", count, len(dc.importPlan.Types))
	}
	return nil
}

func (dc *domainContext) thePlanOrderShouldPlaceTypeBeforeType(first, second string) error {
	if dc.importPlan == nil {
		return fmt.Errorf("import plan is nil")
	}
	firstIdx := -1
	secondIdx := -1
	for orderPos, objIdx := range dc.importPlan.Order {
		obj := dc.importPlan.Objects[objIdx]
		if obj.TypeName == first && firstIdx == -1 {
			firstIdx = orderPos
		}
		if obj.TypeName == second && secondIdx == -1 {
			secondIdx = orderPos
		}
	}
	if firstIdx == -1 {
		return fmt.Errorf("type %q not found in order", first)
	}
	if secondIdx == -1 {
		return fmt.Errorf("type %q not found in order", second)
	}
	if firstIdx >= secondIdx {
		return fmt.Errorf("expected %q (pos %d) before %q (pos %d)", first, firstIdx, second, secondIdx)
	}
	return nil
}

func (dc *domainContext) thePlanObjectNShouldHaveConflict(idx int, conflict string) error {
	if dc.importPlan == nil {
		return fmt.Errorf("import plan is nil")
	}
	if idx < 0 || idx >= len(dc.importPlan.Objects) {
		return fmt.Errorf("object index %d out of range (have %d objects)", idx, len(dc.importPlan.Objects))
	}
	if dc.importPlan.Objects[idx].Conflict != conflict {
		return fmt.Errorf("object %d: expected conflict %q, got %q", idx, conflict, dc.importPlan.Objects[idx].Conflict)
	}
	return nil
}

func (dc *domainContext) thePlanShouldIncludeAllObjectsInTheOrder() error {
	if dc.importPlan == nil {
		return fmt.Errorf("import plan is nil")
	}
	if len(dc.importPlan.Order) != len(dc.importPlan.Objects) {
		return fmt.Errorf("expected %d objects in order, got %d", len(dc.importPlan.Objects), len(dc.importPlan.Order))
	}
	return nil
}

func initImportPlanSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	ctx.Step(`^a classification list:$`, dc.aClassificationList)
	ctx.Step(`^a classification list with circular dependencies$`, dc.aClassificationListWithCircularDependencies)
	ctx.Step(`^I generate an import plan$`, dc.iGenerateAnImportPlan)
	ctx.Step(`^the plan should have (\d+) new types$`, dc.thePlanShouldHaveNNewTypes)
	ctx.Step(`^the plan order should place "([^"]*)" objects before "([^"]*)" objects$`, dc.thePlanOrderShouldPlaceTypeBeforeType)
	ctx.Step(`^the plan object (\d+) should have conflict "([^"]*)"$`, dc.thePlanObjectNShouldHaveConflict)
	ctx.Step(`^the plan should include all objects in the order$`, dc.thePlanShouldIncludeAllObjectsInTheOrder)
}
