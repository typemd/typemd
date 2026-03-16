package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/typemd/typemd/core"
)

func setupTestVaultDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	v := core.NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	// Write test type schemas (book, person, note are no longer built-in)
	os.WriteFile(filepath.Join(v.TypesDir(), "book.yaml"), []byte("name: book\nproperties:\n  - name: title\n    type: string\n"), 0644)
	os.WriteFile(filepath.Join(v.TypesDir(), "person.yaml"), []byte("name: person\nproperties:\n  - name: role\n    type: string\n"), 0644)
	if err := v.Open(); err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	v.Close()
	return dir
}

func TestCreateCmd_Success(t *testing.T) {
	dir := setupTestVaultDir(t)

	vaultPath = dir
	rootCmd.SetArgs([]string{"object", "create", "book", "clean-code"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	// Verify file matching pattern was created (filename includes ULID suffix)
	matches, _ := filepath.Glob(filepath.Join(dir, "objects", "book", "clean-code-*.md"))
	if len(matches) != 1 {
		t.Errorf("expected 1 matching file, got %d", len(matches))
	}
}

func TestCreateCmd_UnknownType(t *testing.T) {
	dir := setupTestVaultDir(t)

	vaultPath = dir
	rootCmd.SetArgs([]string{"object", "create", "nonexistent", "foo"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for unknown type")
	}
}

func TestCreateCmd_SameNameTwice(t *testing.T) {
	dir := setupTestVaultDir(t)

	vaultPath = dir

	// Create first
	rootCmd.SetArgs([]string{"object", "create", "book", "duplicate"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("first create error = %v", err)
	}

	// Create with same name — should succeed because ULID guarantees uniqueness
	rootCmd.SetArgs([]string{"object", "create", "book", "duplicate"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("second create should succeed with ULID, error = %v", err)
	}

	matches, _ := filepath.Glob(filepath.Join(dir, "objects", "book", "duplicate-*.md"))
	if len(matches) != 2 {
		t.Errorf("expected 2 matching files, got %d", len(matches))
	}
}

func TestCreateCmd_WithTemplateFlag(t *testing.T) {
	dir := setupTestVaultDir(t)

	// Create a template
	tmplDir := filepath.Join(dir, "templates", "book")
	os.MkdirAll(tmplDir, 0755)
	os.WriteFile(filepath.Join(tmplDir, "review.md"), []byte("---\ntitle: Review Template\n---\n\n## Review Notes\n"), 0644)

	vaultPath = dir
	templateFlag = "review"
	defer func() { templateFlag = "" }()

	rootCmd.SetArgs([]string{"object", "create", "book", "my-book", "-t", "review"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	matches, _ := filepath.Glob(filepath.Join(dir, "objects", "book", "my-book-*.md"))
	if len(matches) != 1 {
		t.Fatalf("expected 1 matching file, got %d", len(matches))
	}

	data, _ := os.ReadFile(matches[0])
	content := string(data)
	if !strings.Contains(content, "Review Notes") {
		t.Error("object should contain template body")
	}
	if !strings.Contains(content, "title: Review Template") {
		t.Error("object should contain template property")
	}
}

func TestCreateCmd_SingleTemplateAutoApply(t *testing.T) {
	dir := setupTestVaultDir(t)

	// Create a single template
	tmplDir := filepath.Join(dir, "templates", "book")
	os.MkdirAll(tmplDir, 0755)
	os.WriteFile(filepath.Join(tmplDir, "default.md"), []byte("## Default Body\n"), 0644)

	vaultPath = dir
	templateFlag = ""

	rootCmd.SetArgs([]string{"object", "create", "book", "auto-book"})
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	matches, _ := filepath.Glob(filepath.Join(dir, "objects", "book", "auto-book-*.md"))
	if len(matches) != 1 {
		t.Fatalf("expected 1 matching file, got %d", len(matches))
	}

	data, _ := os.ReadFile(matches[0])
	if !strings.Contains(string(data), "Default Body") {
		t.Error("single template should be auto-applied")
	}
}

func TestCreateCmd_InvalidTemplateFlag(t *testing.T) {
	dir := setupTestVaultDir(t)

	vaultPath = dir
	templateFlag = "nonexistent"
	defer func() { templateFlag = "" }()

	rootCmd.SetArgs([]string{"object", "create", "book", "test-book", "-t", "nonexistent"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for nonexistent template")
	}
}

func TestResolveTypeAndName(t *testing.T) {
	knownTypes := []string{"book", "idea", "note", "tag"}

	tests := []struct {
		name        string
		args        []string
		typeFlag    string
		defaultType string
		wantType    string
		wantName    string
		wantErr     bool
	}{
		{
			name:     "two args: type and name (backward compatible)",
			args:     []string{"book", "Clean Code"},
			wantType: "book",
			wantName: "Clean Code",
		},
		{
			name:     "one arg: known type",
			args:     []string{"book"},
			wantType: "book",
			wantName: "",
		},
		{
			name:        "one arg: not a type, has default",
			args:        []string{"Some Thought"},
			defaultType: "idea",
			wantType:    "idea",
			wantName:    "Some Thought",
		},
		{
			name:    "one arg: not a type, no default",
			args:    []string{"Some Thought"},
			wantErr: true,
		},
		{
			name:        "zero args: has default type",
			args:        []string{},
			defaultType: "idea",
			wantType:    "idea",
			wantName:    "",
		},
		{
			name:    "zero args: no default type",
			args:    []string{},
			wantErr: true,
		},
		{
			name:     "type flag: with name",
			args:     []string{"Meeting Notes"},
			typeFlag: "note",
			wantType: "note",
			wantName: "Meeting Notes",
		},
		{
			name:     "type flag: no name",
			args:     []string{},
			typeFlag: "idea",
			wantType: "idea",
			wantName: "",
		},
		{
			name:        "type flag overrides config default",
			args:        []string{"Meeting Notes"},
			typeFlag:    "note",
			defaultType: "idea",
			wantType:    "note",
			wantName:    "Meeting Notes",
		},
		{
			name:     "type flag: multiple name args joined",
			args:     []string{"Some", "Thought"},
			typeFlag: "idea",
			wantType: "idea",
			wantName: "Some Thought",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotType, gotName, err := resolveTypeAndName(tt.args, tt.typeFlag, tt.defaultType, func() []string { return knownTypes })
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if gotType != tt.wantType {
				t.Errorf("type = %q, want %q", gotType, tt.wantType)
			}
			if gotName != tt.wantName {
				t.Errorf("name = %q, want %q", gotName, tt.wantName)
			}
		})
	}
}

