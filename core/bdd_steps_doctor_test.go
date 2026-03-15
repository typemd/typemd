package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cucumber/godog"
)

// ── Doctor steps ────────────────────────────────────────────────────────────

func (dc *domainContext) iRunDoctor() {
	dc.doctorReport = RunDoctor(dc.vault)
}

func (dc *domainContext) theDoctorReportShouldHaveNCategories(n int) error {
	if dc.doctorReport == nil {
		return fmt.Errorf("doctor report is nil")
	}
	got := len(dc.doctorReport.Categories)
	if got != n {
		return fmt.Errorf("expected %d categories, got %d", n, got)
	}
	return nil
}

func (dc *domainContext) findCategory(name string) (*DoctorCategory, error) {
	if dc.doctorReport == nil {
		return nil, fmt.Errorf("doctor report is nil")
	}
	for i := range dc.doctorReport.Categories {
		if dc.doctorReport.Categories[i].Name == name {
			return &dc.doctorReport.Categories[i], nil
		}
	}
	return nil, fmt.Errorf("category %q not found in report", name)
}

func (dc *domainContext) theCategoryShouldPass(name string) error {
	cat, err := dc.findCategory(name)
	if err != nil {
		return err
	}
	if len(cat.Issues) != 0 {
		return fmt.Errorf("expected category %q to pass (0 issues), got %d issues", name, len(cat.Issues))
	}
	return nil
}

func (dc *domainContext) theCategoryShouldHaveNIssues(name string, n int) error {
	cat, err := dc.findCategory(name)
	if err != nil {
		return err
	}
	if len(cat.Issues) != n {
		return fmt.Errorf("expected category %q to have %d issues, got %d", name, n, len(cat.Issues))
	}
	return nil
}

func (dc *domainContext) theDoctorReportShouldHaveNTotalIssues(n int) error {
	if dc.doctorReport == nil {
		return fmt.Errorf("doctor report is nil")
	}
	got := dc.doctorReport.TotalIssues()
	if got != n {
		return fmt.Errorf("expected %d total issues, got %d", n, got)
	}
	return nil
}

func (dc *domainContext) theDoctorReportShouldHaveNAutoFixed(n int) error {
	if dc.doctorReport == nil {
		return fmt.Errorf("doctor report is nil")
	}
	got := dc.doctorReport.TotalAutoFixed()
	if got != n {
		return fmt.Errorf("expected %d auto-fixed, got %d", n, got)
	}
	return nil
}

func (dc *domainContext) aCorruptedObjectFileExistsAt(relPath string) {
	fullPath := filepath.Join(dc.vault.ObjectsDir(), relPath)
	os.MkdirAll(filepath.Dir(fullPath), 0755)
	// Write a file with invalid YAML frontmatter
	os.WriteFile(fullPath, []byte("---\n: invalid: yaml: [broken\n---\nsome body\n"), 0644)
}

func (dc *domainContext) anOrphanObjectDirectoryExists(dirName string) {
	orphanDir := filepath.Join(dc.vault.ObjectsDir(), dirName)
	os.MkdirAll(orphanDir, 0755)
}

func (dc *domainContext) theIndexIsOutOfSync() {
	// Create a valid object file on disk that the index doesn't know about.
	// The doctor should detect and auto-fix this by re-syncing.
	typeName := "book"
	slug := "unindexed-" + mustULID()
	filename := slug
	objDir := filepath.Join(dc.vault.ObjectsDir(), typeName)
	os.MkdirAll(objDir, 0755)
	content := fmt.Sprintf("---\nname: %s\ntitle: Unindexed Book\n---\n", slug)
	os.WriteFile(filepath.Join(objDir, filename+".md"), []byte(content), 0644)
}

func initDoctorSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	ctx.Step(`^I run doctor$`, dc.iRunDoctor)
	ctx.Step(`^the doctor report should have (\d+) categories$`, dc.theDoctorReportShouldHaveNCategories)
	ctx.Step(`^the "([^"]*)" category should pass$`, dc.theCategoryShouldPass)
	ctx.Step(`^the "([^"]*)" category should have (\d+) issues?$`, dc.theCategoryShouldHaveNIssues)
	ctx.Step(`^the doctor report should have (\d+) total issues$`, dc.theDoctorReportShouldHaveNTotalIssues)
	ctx.Step(`^the doctor report should have (\d+) auto-fixed$`, dc.theDoctorReportShouldHaveNAutoFixed)
	ctx.Step(`^a corrupted object file exists at "([^"]*)"$`, dc.aCorruptedObjectFileExistsAt)
	ctx.Step(`^an orphan object directory "([^"]*)" exists$`, dc.anOrphanObjectDirectoryExists)
	ctx.Step(`^the index is out of sync$`, dc.theIndexIsOutOfSync)
}
