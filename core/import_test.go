package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseFrontmatter_ForImport(t *testing.T) {
	data := []byte("---\ntitle: Hello\nauthor: World\n---\nBody text\n")
	fm, body, err := parseFrontmatter(data)
	if err != nil {
		t.Fatal(err)
	}
	if fm["title"] != "Hello" {
		t.Errorf("title = %q, want %q", fm["title"], "Hello")
	}
	if fm["author"] != "World" {
		t.Errorf("author = %q, want %q", fm["author"], "World")
	}
	if body != "Body text\n" {
		t.Errorf("body = %q, want %q", body, "Body text\n")
	}
}

func TestParseFrontmatter_NoFrontmatter(t *testing.T) {
	data := []byte("Just plain text\nno frontmatter here\n")
	fm, _, err := parseFrontmatter(data)
	// parseFrontmatter returns empty map for no frontmatter
	if err != nil {
		t.Fatal(err)
	}
	if len(fm) != 0 {
		t.Errorf("expected empty map, got %v", fm)
	}
}

func TestScanSources_NonExistentPath(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatal(err)
	}

	_, err := v.ScanSources([]string{"does-not-exist"})
	if err == nil {
		t.Error("expected error for non-existent path")
	}
}

func TestScanSources_EmptyDirectory(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatal(err)
	}

	emptyDir := filepath.Join(dir, "empty")
	if err := os.MkdirAll(emptyDir, 0755); err != nil {
		t.Fatal(err)
	}

	result, err := v.ScanSources([]string{"empty"})
	if err != nil {
		t.Fatal(err)
	}
	if result.FileCount != 0 {
		t.Errorf("FileCount = %d, want 0", result.FileCount)
	}
}

func TestScanSources_MixedFileTypes(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatal(err)
	}

	srcDir := filepath.Join(dir, "mixed")
	if err := os.MkdirAll(srcDir, 0755); err != nil {
		t.Fatal(err)
	}
	os.WriteFile(filepath.Join(srcDir, "note.md"), []byte("---\ntitle: Note\n---\n"), 0644)
	os.WriteFile(filepath.Join(srcDir, "image.png"), []byte("fake"), 0644)
	os.WriteFile(filepath.Join(srcDir, "data.json"), []byte("{}"), 0644)

	result, err := v.ScanSources([]string{"mixed"})
	if err != nil {
		t.Fatal(err)
	}
	if result.FileCount != 1 {
		t.Errorf("FileCount = %d, want 1 (only .md files)", result.FileCount)
	}
}

func TestScanSources_FilesWithoutFrontmatter(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatal(err)
	}

	srcDir := filepath.Join(dir, "raw")
	if err := os.MkdirAll(srcDir, 0755); err != nil {
		t.Fatal(err)
	}
	os.WriteFile(filepath.Join(srcDir, "plain.md"), []byte("Just text"), 0644)
	os.WriteFile(filepath.Join(srcDir, "with-fm.md"), []byte("---\ntitle: yes\n---\n"), 0644)

	result, err := v.ScanSources([]string{"raw"})
	if err != nil {
		t.Fatal(err)
	}
	if result.NoFrontmatterCount != 1 {
		t.Errorf("NoFrontmatterCount = %d, want 1", result.NoFrontmatterCount)
	}
}

func TestComputeImportOrder_TagsFirst(t *testing.T) {
	v := NewVault(t.TempDir())
	objects := []ObjectPlan{
		{TypeName: "book", Name: "a"},
		{TypeName: "tag", Name: "go"},
		{TypeName: "page", Name: "c"},
	}
	order := v.computeImportOrder(objects)
	if len(order) != 3 {
		t.Fatalf("order length = %d, want 3", len(order))
	}
	// Tag (index 1) should be first
	if order[0] != 1 {
		t.Errorf("order[0] = %d, want 1 (tag index)", order[0])
	}
}

func TestComputeImportOrder_CircularDeps(t *testing.T) {
	v := NewVault(t.TempDir())
	objects := []ObjectPlan{
		{TypeName: "page", Name: "A", DependsOn: []int{1}},
		{TypeName: "page", Name: "B", DependsOn: []int{0}},
	}
	order := v.computeImportOrder(objects)
	if len(order) != 2 {
		t.Fatalf("order length = %d, want 2 (circular deps broken)", len(order))
	}
}

func TestGenerateSuggestions_WithFailures(t *testing.T) {
	v := NewVault(t.TempDir())
	report := &ImportReport{ObjectsFailed: 3}
	v.generateSuggestions(report)
	if len(report.Suggestions) == 0 {
		t.Error("expected suggestions for failed imports")
	}
}

func TestGenerateSuggestions_WithUnresolvedRefs(t *testing.T) {
	v := NewVault(t.TempDir())
	report := &ImportReport{
		UnresolvedRefs: []UnresolvedRef{{SourceObjectID: "a", Reference: "b"}},
	}
	v.generateSuggestions(report)
	if len(report.Suggestions) == 0 {
		t.Error("expected suggestions for unresolved refs")
	}
}

func TestOnboardingSkillDiscoverable(t *testing.T) {
	skills := ListSkills()
	found := false
	for _, s := range skills {
		if s.Name == "onboarding" {
			found = true
			break
		}
	}
	if !found {
		t.Error("onboarding skill not found in ListSkills()")
	}
}

func TestGetSkill_Onboarding(t *testing.T) {
	skill, err := GetSkill("onboarding")
	if err != nil {
		t.Fatal(err)
	}
	if skill.Name != "onboarding" {
		t.Errorf("skill.Name = %q, want %q", skill.Name, "onboarding")
	}
	if skill.Instructions == "" {
		t.Error("skill.Instructions is empty")
	}
}
