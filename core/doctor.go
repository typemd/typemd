package core

import "fmt"

// IssueSeverity represents the severity level of a doctor issue.
type IssueSeverity int

const (
	SeverityError IssueSeverity = iota
	SeverityWarning
)

// DoctorIssue represents a single health check finding.
type DoctorIssue struct {
	Severity IssueSeverity
	Message  string
}

// DoctorCategory represents results for one category of health checks.
type DoctorCategory struct {
	Name      string
	Issues    []DoctorIssue
	AutoFixed int
}

// DoctorReport holds the complete results of a vault health check.
type DoctorReport struct {
	Categories []DoctorCategory
}

// TotalIssues returns the total number of issues across all categories.
func (r *DoctorReport) TotalIssues() int {
	total := 0
	for _, c := range r.Categories {
		total += len(c.Issues)
	}
	return total
}

// TotalAutoFixed returns the total number of auto-fixed items.
func (r *DoctorReport) TotalAutoFixed() int {
	total := 0
	for _, c := range r.Categories {
		total += c.AutoFixed
	}
	return total
}

// HasErrors returns true if any category has error-severity issues.
func (r *DoctorReport) HasErrors() bool {
	for _, c := range r.Categories {
		for _, issue := range c.Issues {
			if issue.Severity == SeverityError {
				return true
			}
		}
	}
	return false
}

// HasUnresolvedIssues returns true if there are issues that were NOT auto-fixed.
// Used for exit code: 0 if only auto-fixed, 1 if unresolved issues remain.
func (r *DoctorReport) HasUnresolvedIssues() bool {
	return r.TotalIssues() > 0
}

// ScanCorruptedFiles scans the vault for object files with unparseable frontmatter.
func ScanCorruptedFiles(v *Vault) []CorruptedFile {
	repo, ok := v.repo.(*LocalObjectRepository)
	if !ok {
		return nil
	}
	_, corrupted, err := repo.WalkAll()
	if err != nil {
		return nil
	}
	return corrupted
}

// RunDoctor performs a comprehensive vault health check across 8 categories.
func RunDoctor(v *Vault) *DoctorReport {
	report := &DoctorReport{}

	// Validation checks (reuse existing Validate* functions)
	report.Categories = append(report.Categories, mapErrorsToCategory("Schemas", ValidateAllSchemas(v)))
	report.Categories = append(report.Categories, mapErrorsToCategory("Objects", ValidateAllObjects(v)))
	report.Categories = append(report.Categories, errorsToCategory("Relations", ValidateRelations(v)))
	report.Categories = append(report.Categories, errorsToCategory("Wiki-links", ValidateWikiLinks(v)))
	report.Categories = append(report.Categories, errorsToCategory("Uniqueness", ValidateNameUniqueness(v)))

	// Structural integrity checks
	report.Categories = append(report.Categories, checkCorruptedFiles(v))
	report.Categories = append(report.Categories, checkIndexSync(v))
	report.Categories = append(report.Categories, checkOrphans(v))

	return report
}

// errorsToCategory converts a slice of errors into a DoctorCategory.
func errorsToCategory(name string, errs []error) DoctorCategory {
	cat := DoctorCategory{Name: name}
	for _, e := range errs {
		cat.Issues = append(cat.Issues, DoctorIssue{
			Severity: SeverityError,
			Message:  e.Error(),
		})
	}
	return cat
}

// mapErrorsToCategory converts a map of keyed errors into a DoctorCategory.
func mapErrorsToCategory(name string, errs map[string][]error) DoctorCategory {
	cat := DoctorCategory{Name: name}
	for key, keyErrs := range errs {
		for _, e := range keyErrs {
			cat.Issues = append(cat.Issues, DoctorIssue{
				Severity: SeverityError,
				Message:  fmt.Sprintf("%s: %s", key, e),
			})
		}
	}
	return cat
}

func checkCorruptedFiles(v *Vault) DoctorCategory {
	cat := DoctorCategory{Name: "Files"}
	for _, cf := range ScanCorruptedFiles(v) {
		cat.Issues = append(cat.Issues, DoctorIssue{
			Severity: SeverityError,
			Message:  fmt.Sprintf("%s: %s", cf.Path, cf.Error),
		})
	}
	return cat
}

func checkIndexSync(v *Vault) DoctorCategory {
	cat := DoctorCategory{Name: "Index"}
	needsSync, err := v.index.NeedsSync()
	if err != nil {
		cat.Issues = append(cat.Issues, DoctorIssue{
			Severity: SeverityError,
			Message:  fmt.Sprintf("check index: %s", err),
		})
		return cat
	}
	if needsSync {
		if _, err := v.SyncIndex(); err != nil {
			cat.Issues = append(cat.Issues, DoctorIssue{
				Severity: SeverityError,
				Message:  fmt.Sprintf("rebuild index: %s", err),
			})
		} else {
			cat.AutoFixed = 1
		}
	}
	return cat
}

func checkOrphans(v *Vault) DoctorCategory {
	cat := DoctorCategory{Name: "Orphans"}
	for _, o := range ScanOrphanDirs(v) {
		cat.Issues = append(cat.Issues, DoctorIssue{
			Severity: SeverityWarning,
			Message:  fmt.Sprintf("%s directory %s has no type schema", o.Kind, o.Path),
		})
	}
	return cat
}
