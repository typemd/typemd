package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDoctorReport_TotalIssues(t *testing.T) {
	tests := []struct {
		name       string
		categories []DoctorCategory
		want       int
	}{
		{
			name:       "empty report",
			categories: nil,
			want:       0,
		},
		{
			name: "single category no issues",
			categories: []DoctorCategory{
				{Name: "Schemas", Issues: nil},
			},
			want: 0,
		},
		{
			name: "single category with issues",
			categories: []DoctorCategory{
				{Name: "Files", Issues: []DoctorIssue{
					{Severity: SeverityError, Message: "corrupted file"},
					{Severity: SeverityWarning, Message: "orphan directory"},
				}},
			},
			want: 2,
		},
		{
			name: "multiple categories with issues",
			categories: []DoctorCategory{
				{Name: "Schemas", Issues: []DoctorIssue{
					{Severity: SeverityError, Message: "invalid schema"},
				}},
				{Name: "Files", Issues: []DoctorIssue{
					{Severity: SeverityWarning, Message: "orphan directory"},
				}},
				{Name: "Index", Issues: nil},
			},
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &DoctorReport{Categories: tt.categories}
			if got := r.TotalIssues(); got != tt.want {
				t.Errorf("TotalIssues() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestDoctorReport_TotalAutoFixed(t *testing.T) {
	tests := []struct {
		name       string
		categories []DoctorCategory
		want       int
	}{
		{
			name:       "empty report",
			categories: nil,
			want:       0,
		},
		{
			name: "no auto-fixes",
			categories: []DoctorCategory{
				{Name: "Schemas", AutoFixed: 0},
				{Name: "Files", AutoFixed: 0},
			},
			want: 0,
		},
		{
			name: "some auto-fixes",
			categories: []DoctorCategory{
				{Name: "Index", AutoFixed: 3},
				{Name: "Files", AutoFixed: 1},
			},
			want: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &DoctorReport{Categories: tt.categories}
			if got := r.TotalAutoFixed(); got != tt.want {
				t.Errorf("TotalAutoFixed() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestDoctorReport_HasErrors(t *testing.T) {
	tests := []struct {
		name       string
		categories []DoctorCategory
		want       bool
	}{
		{
			name:       "empty report",
			categories: nil,
			want:       false,
		},
		{
			name: "only warnings",
			categories: []DoctorCategory{
				{Name: "Files", Issues: []DoctorIssue{
					{Severity: SeverityWarning, Message: "orphan directory"},
				}},
			},
			want: false,
		},
		{
			name: "has errors",
			categories: []DoctorCategory{
				{Name: "Files", Issues: []DoctorIssue{
					{Severity: SeverityError, Message: "corrupted file"},
				}},
			},
			want: true,
		},
		{
			name: "mixed errors and warnings",
			categories: []DoctorCategory{
				{Name: "Files", Issues: []DoctorIssue{
					{Severity: SeverityWarning, Message: "orphan directory"},
					{Severity: SeverityError, Message: "corrupted file"},
				}},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &DoctorReport{Categories: tt.categories}
			if got := r.HasErrors(); got != tt.want {
				t.Errorf("HasErrors() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoctorReport_HasUnresolvedIssues(t *testing.T) {
	tests := []struct {
		name       string
		categories []DoctorCategory
		want       bool
	}{
		{
			name:       "empty report — no issues",
			categories: nil,
			want:       false,
		},
		{
			name: "only auto-fixed — no unresolved issues",
			categories: []DoctorCategory{
				{Name: "Index", AutoFixed: 2, Issues: nil},
			},
			want: false,
		},
		{
			name: "has real issues — unresolved",
			categories: []DoctorCategory{
				{Name: "Files", Issues: []DoctorIssue{
					{Severity: SeverityError, Message: "corrupted file"},
				}},
			},
			want: true,
		},
		{
			name: "auto-fixed plus real issues",
			categories: []DoctorCategory{
				{Name: "Index", AutoFixed: 1, Issues: nil},
				{Name: "Files", Issues: []DoctorIssue{
					{Severity: SeverityWarning, Message: "orphan directory"},
				}},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &DoctorReport{Categories: tt.categories}
			if got := r.HasUnresolvedIssues(); got != tt.want {
				t.Errorf("HasUnresolvedIssues() = %v, want %v", got, tt.want)
			}
		})
	}
}

func setupDoctorVault(t *testing.T) *Vault {
	t.Helper()
	dir := t.TempDir()
	v := NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatalf("init vault: %v", err)
	}
	os.WriteFile(filepath.Join(v.TypesDir(), "book.yaml"),
		[]byte("name: book\nproperties:\n  - name: title\n    type: string\n"), 0644)
	if err := v.Open(); err != nil {
		t.Fatalf("open vault: %v", err)
	}
	t.Cleanup(func() { v.Close() })
	return v
}

func TestRunDoctor_HealthyVault(t *testing.T) {
	v := setupDoctorVault(t)
	report := RunDoctor(v)
	if len(report.Categories) != 8 {
		t.Errorf("categories = %d, want 8", len(report.Categories))
	}
	if report.TotalIssues() != 0 {
		t.Errorf("total issues = %d, want 0", report.TotalIssues())
	}
	if report.HasUnresolvedIssues() {
		t.Error("expected no unresolved issues")
	}
}

func TestRunDoctor_CorruptedFile(t *testing.T) {
	v := setupDoctorVault(t)
	objDir := filepath.Join(v.ObjectsDir(), "book")
	os.MkdirAll(objDir, 0755)
	os.WriteFile(filepath.Join(objDir, "bad.md"),
		[]byte("---\n: invalid: [broken\n---\n"), 0644)

	report := RunDoctor(v)
	var filesCat *DoctorCategory
	for i := range report.Categories {
		if report.Categories[i].Name == "Files" {
			filesCat = &report.Categories[i]
			break
		}
	}
	if filesCat == nil {
		t.Fatal("Files category not found")
	}
	if len(filesCat.Issues) != 1 {
		t.Errorf("Files issues = %d, want 1", len(filesCat.Issues))
	}
}

func TestRunDoctor_OrphanDir(t *testing.T) {
	v := setupDoctorVault(t)
	os.MkdirAll(filepath.Join(v.ObjectsDir(), "ghost"), 0755)

	report := RunDoctor(v)
	var orphanCat *DoctorCategory
	for i := range report.Categories {
		if report.Categories[i].Name == "Orphans" {
			orphanCat = &report.Categories[i]
			break
		}
	}
	if orphanCat == nil {
		t.Fatal("Orphans category not found")
	}
	if len(orphanCat.Issues) != 1 {
		t.Errorf("Orphans issues = %d, want 1", len(orphanCat.Issues))
	}
	if orphanCat.Issues[0].Severity != SeverityWarning {
		t.Error("orphan should be a warning, not an error")
	}
}

func TestRunDoctor_ExitCodeOnlyAutoFixed(t *testing.T) {
	report := &DoctorReport{
		Categories: []DoctorCategory{
			{Name: "Index", AutoFixed: 1},
			{Name: "Schemas"},
		},
	}
	if report.HasUnresolvedIssues() {
		t.Error("expected no unresolved issues when only auto-fixed")
	}
}

func TestRunDoctor_ExitCodeWithIssues(t *testing.T) {
	report := &DoctorReport{
		Categories: []DoctorCategory{
			{Name: "Index", AutoFixed: 1},
			{Name: "Files", Issues: []DoctorIssue{
				{Severity: SeverityError, Message: "corrupted"},
			}},
		},
	}
	if !report.HasUnresolvedIssues() {
		t.Error("expected unresolved issues")
	}
}
