package core

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"
)

func TestVault_Init(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)

	if err := v.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	// 目錄結構
	for _, d := range []string{v.TypesDir(), v.ObjectsDir()} {
		if _, err := os.Stat(d); os.IsNotExist(err) {
			t.Errorf("expected directory %s to exist", d)
		}
	}

	// index.db 存在
	if _, err := os.Stat(v.DBPath()); os.IsNotExist(err) {
		t.Error("expected index.db to exist")
	}
}

func TestVault_Init_Gitignore(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)

	if err := v.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	gitignore := filepath.Join(v.Dir(), ".gitignore")
	data, err := os.ReadFile(gitignore)
	if err != nil {
		t.Fatalf("expected .gitignore to exist: %v", err)
	}
	if string(data) != "index.db\n" {
		t.Errorf(".gitignore content = %q, want %q", string(data), "index.db\n")
	}
}

func TestVault_Init_AlreadyInitialized(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)

	if err := v.Init(); err != nil {
		t.Fatalf("first Init() error = %v", err)
	}

	err := v.Init()
	if err == nil {
		t.Fatal("expected error on second Init(), got nil")
	}
}

func TestVault_IsInitialized(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)

	if v.IsInitialized() {
		t.Error("expected IsInitialized() = false before Init()")
	}

	if err := v.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	if !v.IsInitialized() {
		t.Error("expected IsInitialized() = true after Init()")
	}
}

func TestVault_Init_DBSchema(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)

	if err := v.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	db, err := sql.Open("sqlite", v.DBPath())
	if err != nil {
		t.Fatalf("open db error = %v", err)
	}
	defer db.Close()

	// 驗證 objects table 存在
	var name string
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='objects'").Scan(&name)
	if err != nil || name != "objects" {
		t.Error("expected 'objects' table to exist in DB")
	}

	// 驗證 relations table 存在
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='relations'").Scan(&name)
	if err != nil || name != "relations" {
		t.Error("expected 'relations' table to exist in DB")
	}
}

func TestVault_Paths(t *testing.T) {
	v := NewVault("/some/root")

	if got := v.Dir(); got != filepath.Join("/some/root", ".typemd") {
		t.Errorf("Dir() = %q", got)
	}
	if got := v.TypesDir(); got != filepath.Join("/some/root", ".typemd", "types") {
		t.Errorf("TypesDir() = %q", got)
	}
	if got := v.DBPath(); got != filepath.Join("/some/root", ".typemd", "index.db") {
		t.Errorf("DBPath() = %q", got)
	}
	if got := v.ObjectsDir(); got != filepath.Join("/some/root", "objects") {
		t.Errorf("ObjectsDir() = %q", got)
	}
}

func TestVault_Open_Close(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	if err := v.Open(); err != nil {
		t.Fatalf("Open() error = %v", err)
	}

	if err := v.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
}

func TestVault_Open_NotInitialized(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)

	err := v.Open()
	if err == nil {
		t.Fatal("expected error on Open() without Init(), got nil")
	}
}

func TestVault_Init_FTSSchema(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)

	if err := v.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	db, err := sql.Open("sqlite", v.DBPath())
	if err != nil {
		t.Fatalf("open db error = %v", err)
	}
	defer db.Close()

	// verify objects_fts virtual table exists
	var name string
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='objects_fts'").Scan(&name)
	if err != nil || name != "objects_fts" {
		t.Error("expected 'objects_fts' virtual table to exist in DB")
	}
}

func TestVault_FTSTrigger_InsertSync(t *testing.T) {
	v := setupTestVault(t)

	// NewObject INSERT should auto-sync FTS via trigger
	if _, err := v.NewObject("book", "golang-in-action"); err != nil {
		t.Fatalf("NewObject() error = %v", err)
	}

	var count int
	// FTS5 requires double quotes around terms with special characters
	err := v.db.QueryRow(`SELECT count(*) FROM objects_fts WHERE objects_fts MATCH '"golang-in-action"'`).Scan(&count)
	if err != nil {
		t.Fatalf("FTS query error = %v", err)
	}
	if count != 1 {
		t.Errorf("FTS count = %d, want 1", count)
	}
}

// setupTestVault creates a temporary vault with Init + Open, auto-cleanup on test end.
func setupTestVault(t *testing.T) *Vault {
	t.Helper()
	dir := t.TempDir()
	v := NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	if err := v.Open(); err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	t.Cleanup(func() { v.Close() })
	return v
}
