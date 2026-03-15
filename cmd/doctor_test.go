package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/typemd/typemd/core"
)

func TestDoctorCmd_HealthyVault(t *testing.T) {
	dir := t.TempDir()
	v := core.NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	v.Close()

	vaultPath = dir
	rootCmd.SetArgs([]string{"doctor"})
	err := rootCmd.Execute()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestDoctorCmd_WithIssues(t *testing.T) {
	dir := t.TempDir()
	v := core.NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	// Create a corrupted file
	objDir := filepath.Join(dir, "objects", "book")
	os.MkdirAll(objDir, 0755)
	os.WriteFile(filepath.Join(objDir, "bad.md"),
		[]byte("---\n: invalid: [broken\n---\n"), 0644)

	v.Close()

	vaultPath = dir
	rootCmd.SetArgs([]string{"doctor"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for unhealthy vault")
	}
	if !strings.Contains(err.Error(), "issue(s)") {
		t.Errorf("error = %q, want it to mention issue(s)", err)
	}
}

func TestPrintDoctorReport_AllPassing(t *testing.T) {
	report := &core.DoctorReport{
		Categories: []core.DoctorCategory{
			{Name: "Schemas"},
			{Name: "Objects"},
		},
	}

	var buf bytes.Buffer
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	printDoctorReport(report)

	w.Close()
	os.Stdout = old
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "✓ Schemas") {
		t.Errorf("expected ✓ Schemas in output, got:\n%s", output)
	}
	if !strings.Contains(output, "No issues found.") {
		t.Errorf("expected 'No issues found.' in output, got:\n%s", output)
	}
}

func TestPrintDoctorReport_WithIssues(t *testing.T) {
	report := &core.DoctorReport{
		Categories: []core.DoctorCategory{
			{Name: "Schemas"},
			{Name: "Files", Issues: []core.DoctorIssue{
				{Severity: core.SeverityError, Message: "book/bad.md: corrupted"},
			}},
		},
	}

	var buf bytes.Buffer
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	printDoctorReport(report)

	w.Close()
	os.Stdout = old
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "✓ Schemas") {
		t.Errorf("expected ✓ Schemas in output, got:\n%s", output)
	}
	if !strings.Contains(output, "✗ Files") {
		t.Errorf("expected ✗ Files in output, got:\n%s", output)
	}
	if !strings.Contains(output, "[error]") {
		t.Errorf("expected [error] prefix in output, got:\n%s", output)
	}
	if !strings.Contains(output, "1 issue(s) found.") {
		t.Errorf("expected '1 issue(s) found.' in output, got:\n%s", output)
	}
}

func TestPrintDoctorReport_AutoFixed(t *testing.T) {
	report := &core.DoctorReport{
		Categories: []core.DoctorCategory{
			{Name: "Index", AutoFixed: 1},
		},
	}

	var buf bytes.Buffer
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	printDoctorReport(report)

	w.Close()
	os.Stdout = old
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "✓ Index (auto-fixed)") {
		t.Errorf("expected '✓ Index (auto-fixed)' in output, got:\n%s", output)
	}
	if !strings.Contains(output, "1 auto-fixed") {
		t.Errorf("expected '1 auto-fixed' in output, got:\n%s", output)
	}
}
