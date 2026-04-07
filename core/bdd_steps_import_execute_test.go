package core

import (
	"fmt"
	"strings"

	"github.com/cucumber/godog"
)

func (dc *domainContext) anImportPlanWithObjects(table *godog.Table) {
	headers := table.Rows[0].Cells
	dc.importPlan = &ImportPlan{}
	for i, row := range table.Rows[1:] {
		vals := make(map[string]string)
		for j, cell := range row.Cells {
			vals[headers[j].Value] = cell.Value
		}
		dc.importPlan.Objects = append(dc.importPlan.Objects, ObjectPlan{
			SourcePath: vals["source_path"],
			TypeName:   vals["type_name"],
			Name:       vals["name"],
			Conflict:   "none",
		})
		dc.importPlan.Order = append(dc.importPlan.Order, i)
	}
}

func (dc *domainContext) anImportPlanWithNewTypeAndObjects(typeName string, table *godog.Table) {
	dc.anImportPlanWithObjects(table)
	dc.importPlan.Types = append(dc.importPlan.Types, TypePlan{
		Name: typeName,
	})
}

func (dc *domainContext) anImportPlanWithASkippedObject(table *godog.Table) {
	headers := table.Rows[0].Cells
	dc.importPlan = &ImportPlan{}
	for i, row := range table.Rows[1:] {
		vals := make(map[string]string)
		for j, cell := range row.Cells {
			vals[headers[j].Value] = cell.Value
		}
		conflict := vals["conflict"]
		if conflict == "" {
			conflict = "none"
		}
		dc.importPlan.Objects = append(dc.importPlan.Objects, ObjectPlan{
			SourcePath: vals["source_path"],
			TypeName:   vals["type_name"],
			Name:       vals["name"],
			Conflict:   conflict,
		})
		dc.importPlan.Order = append(dc.importPlan.Order, i)
	}
}

func (dc *domainContext) anImportPlanWithObjectsForANonexistentType(table *godog.Table) {
	dc.anImportPlanWithObjects(table)
}

func (dc *domainContext) iExecuteTheImportPlan() {
	dc.importReport, dc.lastErr = dc.vault.ExecutePlan(dc.importPlan)
}

func (dc *domainContext) theImportReportShouldShowNCreated(count int) error {
	if dc.importReport == nil {
		return fmt.Errorf("import report is nil")
	}
	if dc.importReport.ObjectsCreated != count {
		return fmt.Errorf("expected %d created, got %d", count, dc.importReport.ObjectsCreated)
	}
	return nil
}

func (dc *domainContext) theImportReportShouldShowNSkipped(count int) error {
	if dc.importReport == nil {
		return fmt.Errorf("import report is nil")
	}
	if dc.importReport.ObjectsSkipped != count {
		return fmt.Errorf("expected %d skipped, got %d", count, dc.importReport.ObjectsSkipped)
	}
	return nil
}

func (dc *domainContext) theImportReportShouldShowNFailed(count int) error {
	if dc.importReport == nil {
		return fmt.Errorf("import report is nil")
	}
	if dc.importReport.ObjectsFailed != count {
		return fmt.Errorf("expected %d failed, got %d", count, dc.importReport.ObjectsFailed)
	}
	return nil
}

func (dc *domainContext) theImportReportShouldShowNTypesCreated(count int) error {
	if dc.importReport == nil {
		return fmt.Errorf("import report is nil")
	}
	if dc.importReport.TypesCreated != count {
		return fmt.Errorf("expected %d types created, got %d", count, dc.importReport.TypesCreated)
	}
	return nil
}

func (dc *domainContext) theImportReportShouldSuggestReviewingFailedFiles() error {
	if dc.importReport == nil {
		return fmt.Errorf("import report is nil")
	}
	for _, s := range dc.importReport.Suggestions {
		if strings.Contains(s, "failed") {
			return nil
		}
	}
	return fmt.Errorf("expected suggestion about failed files, got: %v", dc.importReport.Suggestions)
}

func initImportExecuteSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	ctx.Step(`^an import plan with objects:$`, dc.anImportPlanWithObjects)
	ctx.Step(`^an import plan with new type "([^"]*)" and objects:$`, dc.anImportPlanWithNewTypeAndObjects)
	ctx.Step(`^an import plan with a skipped object:$`, dc.anImportPlanWithASkippedObject)
	ctx.Step(`^an import plan with objects for a nonexistent type:$`, dc.anImportPlanWithObjectsForANonexistentType)
	ctx.Step(`^I execute the import plan$`, dc.iExecuteTheImportPlan)
	ctx.Step(`^the import report should show (\d+) created$`, dc.theImportReportShouldShowNCreated)
	ctx.Step(`^the import report should show (\d+) skipped$`, dc.theImportReportShouldShowNSkipped)
	ctx.Step(`^the import report should show (\d+) failed$`, dc.theImportReportShouldShowNFailed)
	ctx.Step(`^the import report should show (\d+) types created$`, dc.theImportReportShouldShowNTypesCreated)
	ctx.Step(`^the import report should suggest reviewing failed files$`, dc.theImportReportShouldSuggestReviewingFailedFiles)
}
