package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/typemd/typemd/core"
)

func setupStatsVault(t *testing.T) (*core.Vault, string) {
	t.Helper()
	dir := t.TempDir()
	v := core.NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	// Write book type schema (directory format)
	typesDir := filepath.Join(dir, ".typemd", "types")
	os.MkdirAll(filepath.Join(typesDir, "book"), 0755)
	os.WriteFile(filepath.Join(typesDir, "book", "schema.yaml"), []byte(`name: book
emoji: "📚"
plural: books
properties:
  - name: rating
    type: number
  - name: status
    type: select
    options:
      - value: reading
      - value: done
`), 0644)

	if err := v.Open(); err != nil {
		t.Fatalf("Open() error = %v", err)
	}

	return v, dir
}

func resetStatsFlags() {
	statsTypeName = ""
	statsJSON = false
}

func TestStatsCmd_VaultWide(t *testing.T) {
	resetStatsFlags()
	v, dir := setupStatsVault(t)

	// Create some objects
	if _, err := v.NewObject("book", "book1", ""); err != nil {
		t.Fatalf("NewObject error = %v", err)
	}
	if _, err := v.NewObject("book", "book2", ""); err != nil {
		t.Fatalf("NewObject error = %v", err)
	}
	v.Close()

	vaultPath = dir
	rootCmd.SetArgs([]string{"stats"})

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := rootCmd.Execute()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !strings.Contains(output, "books") {
		t.Errorf("expected 'books' in output, got:\n%s", output)
	}
	if !strings.Contains(output, "2") {
		t.Errorf("expected count '2' in output, got:\n%s", output)
	}
	if !strings.Contains(output, "Total") {
		t.Errorf("expected 'Total' in output, got:\n%s", output)
	}
}

func TestStatsCmd_SingleType(t *testing.T) {
	resetStatsFlags()
	v, dir := setupStatsVault(t)

	obj, err := v.NewObject("book", "book1", "")
	if err != nil {
		t.Fatalf("NewObject error = %v", err)
	}
	if err := v.SetProperty(obj.ID, "rating", float64(4)); err != nil {
		t.Fatalf("SetProperty error = %v", err)
	}
	v.Close()

	vaultPath = dir
	rootCmd.SetArgs([]string{"stats", "--type", "book"})

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err = rootCmd.Execute()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !strings.Contains(output, "book (1 objects)") {
		t.Errorf("expected 'book (1 objects)' in output, got:\n%s", output)
	}
	if !strings.Contains(output, "rating") {
		t.Errorf("expected 'rating' in output, got:\n%s", output)
	}
}

func TestStatsCmd_JSON(t *testing.T) {
	resetStatsFlags()
	v, dir := setupStatsVault(t)
	if _, err := v.NewObject("book", "book1", ""); err != nil {
		t.Fatalf("NewObject error = %v", err)
	}
	v.Close()

	vaultPath = dir
	rootCmd.SetArgs([]string{"stats", "--json"})

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := rootCmd.Execute()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !strings.Contains(output, `"Total"`) && !strings.Contains(output, `"total"`) {
		// JSON output should have structured data
		if !strings.Contains(output, "{") {
			t.Errorf("expected JSON output, got:\n%s", output)
		}
	}
}

func TestStatsCmd_NonExistentType(t *testing.T) {
	resetStatsFlags()
	v, dir := setupStatsVault(t)
	v.Close()

	vaultPath = dir
	rootCmd.SetArgs([]string{"stats", "--type", "nonexistent"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for non-existent type")
	}
	if !strings.Contains(err.Error(), "nonexistent") {
		t.Errorf("error = %q, expected to mention 'nonexistent'", err)
	}
}

func TestStatsCmd_EmptyVault(t *testing.T) {
	resetStatsFlags()
	dir := t.TempDir()
	v := core.NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	v.Close()

	vaultPath = dir
	rootCmd.SetArgs([]string{"stats"})

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := rootCmd.Execute()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !strings.Contains(output, "No objects") {
		t.Errorf("expected 'No objects' message, got:\n%s", output)
	}
}
