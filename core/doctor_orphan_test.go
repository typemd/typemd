package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanOrphanDirs_NoOrphans(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	// Create a custom type schema.
	schema := "name: book\nproperties:\n  - name: title\n    type: string\n"
	if err := os.WriteFile(filepath.Join(v.TypesDir(), "book.yaml"), []byte(schema), 0644); err != nil {
		t.Fatalf("write type schema: %v", err)
	}

	// Create matching directories.
	os.MkdirAll(filepath.Join(v.ObjectsDir(), "book"), 0755)
	os.MkdirAll(filepath.Join(v.ObjectsDir(), "tag"), 0755)
	os.MkdirAll(filepath.Join(v.TemplatesDir(), "book"), 0755)

	orphans := ScanOrphanDirs(v)
	if len(orphans) != 0 {
		t.Errorf("expected 0 orphans, got %d: %v", len(orphans), orphans)
	}
}

func TestScanOrphanDirs_OrphanObjectDir(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	// Create a directory for a type that has no schema.
	os.MkdirAll(filepath.Join(v.ObjectsDir(), "note"), 0755)

	orphans := ScanOrphanDirs(v)
	if len(orphans) != 1 {
		t.Fatalf("expected 1 orphan, got %d: %v", len(orphans), orphans)
	}
	if orphans[0].Path != filepath.Join("objects", "note") {
		t.Errorf("Path = %q, want %q", orphans[0].Path, filepath.Join("objects", "note"))
	}
	if orphans[0].Kind != "object" {
		t.Errorf("Kind = %q, want %q", orphans[0].Kind, "object")
	}
}

func TestScanOrphanDirs_OrphanTemplateDir(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	// Create a templates directory for a non-existent type.
	os.MkdirAll(filepath.Join(v.TemplatesDir(), "recipe"), 0755)

	orphans := ScanOrphanDirs(v)
	if len(orphans) != 1 {
		t.Fatalf("expected 1 orphan, got %d: %v", len(orphans), orphans)
	}
	if orphans[0].Path != filepath.Join("templates", "recipe") {
		t.Errorf("Path = %q, want %q", orphans[0].Path, filepath.Join("templates", "recipe"))
	}
	if orphans[0].Kind != "template" {
		t.Errorf("Kind = %q, want %q", orphans[0].Kind, "template")
	}
}

func TestScanOrphanDirs_NoObjectsDir(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	// Remove the objects directory that Init() creates.
	os.RemoveAll(v.ObjectsDir())

	orphans := ScanOrphanDirs(v)
	if len(orphans) != 0 {
		t.Errorf("expected 0 orphans when objects/ missing, got %d: %v", len(orphans), orphans)
	}
}

func TestScanOrphanDirs_NoTemplatesDir(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	// templates/ is not created by Init(), so it shouldn't exist.
	// Verify it doesn't exist, then scan.
	if _, err := os.Stat(v.TemplatesDir()); !os.IsNotExist(err) {
		t.Fatalf("expected templates/ to not exist, but it does")
	}

	orphans := ScanOrphanDirs(v)
	if len(orphans) != 0 {
		t.Errorf("expected 0 orphans when templates/ missing, got %d: %v", len(orphans), orphans)
	}
}

func TestScanOrphanDirs_BuiltInTypes(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	// "tag" is a built-in type — its directory should NOT be flagged as orphan.
	os.MkdirAll(filepath.Join(v.ObjectsDir(), "tag"), 0755)

	orphans := ScanOrphanDirs(v)
	if len(orphans) != 0 {
		t.Errorf("expected 0 orphans (tag is built-in), got %d: %v", len(orphans), orphans)
	}
}
