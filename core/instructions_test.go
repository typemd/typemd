package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseSkillMD_NoFrontmatter(t *testing.T) {
	content := []byte("# Just a title\n\nSome body content.\n")
	fm, body := parseSkillMD(content)

	if fm.Name != "" {
		t.Errorf("expected empty name, got %q", fm.Name)
	}
	if fm.Description != "" {
		t.Errorf("expected empty description, got %q", fm.Description)
	}
	if body != string(content) {
		t.Errorf("expected body to be entire content, got %q", body)
	}
}

func TestParseSkillMD_WithFrontmatter(t *testing.T) {
	content := []byte("---\nname: test-skill\ndescription: A test skill\n---\n\n# Test\n\nBody here.\n")
	fm, body := parseSkillMD(content)

	if fm.Name != "test-skill" {
		t.Errorf("expected name %q, got %q", "test-skill", fm.Name)
	}
	if fm.Description != "A test skill" {
		t.Errorf("expected description %q, got %q", "A test skill", fm.Description)
	}
	if !strings.Contains(body, "# Test") {
		t.Errorf("expected body to contain '# Test', got %q", body)
	}
	if strings.Contains(body, "---") {
		t.Errorf("body should not contain frontmatter delimiters")
	}
}

func TestParseSkillMD_MalformedYAML(t *testing.T) {
	content := []byte("---\nname: [invalid yaml\n---\n\nBody.\n")
	fm, body := parseSkillMD(content)

	// Malformed YAML falls back to treating entire content as body
	if fm.Name != "" {
		t.Errorf("expected empty name for malformed YAML, got %q", fm.Name)
	}
	if body != string(content) {
		t.Errorf("expected body to be entire content for malformed YAML")
	}
}

func TestParseSkillMD_UnclosedFrontmatter(t *testing.T) {
	content := []byte("---\nname: test\nno closing delimiter\n")
	fm, body := parseSkillMD(content)

	if fm.Name != "" {
		t.Errorf("expected empty name for unclosed frontmatter, got %q", fm.Name)
	}
	if body != string(content) {
		t.Errorf("expected body to be entire content for unclosed frontmatter")
	}
}

func TestGetSkill_Unknown(t *testing.T) {
	_, err := GetSkill("nonexistent")
	if err == nil {
		t.Fatal("expected error for unknown skill")
	}
	if !strings.Contains(err.Error(), "unknown skill") {
		t.Errorf("expected 'unknown skill' in error, got: %v", err)
	}
	if !strings.Contains(err.Error(), "explore") {
		t.Errorf("expected available skills listed in error, got: %v", err)
	}
}

func TestGetSkillRaw_Unknown(t *testing.T) {
	_, err := GetSkillRaw("nonexistent")
	if err == nil {
		t.Fatal("expected error for unknown skill")
	}
	if !strings.Contains(err.Error(), "unknown skill") {
		t.Errorf("expected 'unknown skill' in error, got: %v", err)
	}
}

func TestListSkills_Count(t *testing.T) {
	skills := ListSkills()
	if len(skills) != 3 {
		t.Errorf("expected 3 skills, got %d", len(skills))
	}
}

func TestGetSkill_EachSkillValid(t *testing.T) {
	for _, name := range []string{"explore", "importer"} {
		skill, err := GetSkill(name)
		if err != nil {
			t.Errorf("GetSkill(%q): %v", name, err)
			continue
		}
		if skill.Name != name {
			t.Errorf("expected name %q, got %q", name, skill.Name)
		}
		if skill.Description == "" {
			t.Errorf("skill %q has empty description", name)
		}
		if skill.Instructions == "" {
			t.Errorf("skill %q has empty instructions", name)
		}
	}
}

func TestGetSkillRaw_ContainsFrontmatter(t *testing.T) {
	raw, err := GetSkillRaw("explore")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(string(raw), "---\n") {
		t.Error("expected raw content to start with ---")
	}
	if !strings.Contains(string(raw), "name: explore") {
		t.Error("expected raw content to contain 'name: explore'")
	}
}

func TestGetSkillWithOverride_NoVaultRoot(t *testing.T) {
	skill, err := GetSkillWithOverride("explore", "")
	if err != nil {
		t.Fatal(err)
	}
	if skill.Name != "explore" {
		t.Errorf("expected name 'explore', got %q", skill.Name)
	}
}

func TestGetSkillWithOverride_NonexistentOverride(t *testing.T) {
	skill, err := GetSkillWithOverride("explore", "/tmp/nonexistent-vault-path")
	if err != nil {
		t.Fatal(err)
	}
	if skill.Name != "explore" {
		t.Errorf("expected name 'explore', got %q", skill.Name)
	}
}

func TestGetSkillWithOverride_WithBody(t *testing.T) {
	dir := t.TempDir()
	overrideDir := filepath.Join(dir, ".typemd", "instructions")
	os.MkdirAll(overrideDir, 0755)
	content := "---\nname: explore\ndescription: Custom override\n---\n\nCustom instructions here\n"
	os.WriteFile(filepath.Join(overrideDir, "explore.md"), []byte(content), 0644)

	skill, err := GetSkillWithOverride("explore", dir)
	if err != nil {
		t.Fatal(err)
	}
	if skill.Name != "explore" {
		t.Errorf("name = %q, want %q", skill.Name, "explore")
	}
	if !strings.Contains(skill.Instructions, "Custom instructions here") {
		t.Errorf("expected override body in instructions, got %q", skill.Instructions)
	}
}

func TestGetSkillWithOverride_WithoutFrontmatter(t *testing.T) {
	dir := t.TempDir()
	overrideDir := filepath.Join(dir, ".typemd", "instructions")
	os.MkdirAll(overrideDir, 0755)
	os.WriteFile(filepath.Join(overrideDir, "explore.md"), []byte("Override content without frontmatter\n"), 0644)

	skill, err := GetSkillWithOverride("explore", dir)
	if err != nil {
		t.Fatal(err)
	}
	// Should fall back to embedded metadata
	if skill.Name != "explore" {
		t.Errorf("name = %q, want %q (should fall back to embedded)", skill.Name, "explore")
	}
	if skill.Description == "" {
		t.Error("description should fall back to embedded, got empty")
	}
	if !strings.Contains(skill.Instructions, "Override content without frontmatter") {
		t.Errorf("expected override body in instructions, got %q", skill.Instructions)
	}
}

func TestBuildSkillContext_WithTypes(t *testing.T) {
	v := setupTestVault(t)

	ctx, err := BuildSkillContext(v)
	if err != nil {
		t.Fatalf("BuildSkillContext error = %v", err)
	}
	if len(ctx.Types) == 0 {
		t.Fatal("expected at least one type in context")
	}

	// Should include the "book" type from setupTestVault
	found := false
	for _, typ := range ctx.Types {
		if typ.Name == "book" {
			found = true
			break
		}
	}
	if !found {
		names := make([]string, len(ctx.Types))
		for i, typ := range ctx.Types {
			names[i] = typ.Name
		}
		t.Errorf("expected 'book' in context types, got %v", names)
	}
}

func TestBuildSkillContext_EmptyVault(t *testing.T) {
	dir := t.TempDir()
	v := NewVault(dir)
	if err := v.Init(); err != nil {
		t.Fatalf("Init() error = %v", err)
	}
	if err := v.Open(); err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer v.Close()

	ctx, err := BuildSkillContext(v)
	if err != nil {
		t.Fatalf("BuildSkillContext error = %v", err)
	}
	if len(ctx.Types) == 0 {
		t.Fatal("expected built-in types in context")
	}

	// Should include built-in types (tag, page)
	typeNames := make(map[string]bool)
	for _, typ := range ctx.Types {
		typeNames[typ.Name] = true
	}
	if !typeNames["tag"] {
		t.Error("expected 'tag' in context types")
	}
	if !typeNames["page"] {
		t.Error("expected 'page' in context types")
	}
}

func TestBuildSkillContext_TypeWithProperties(t *testing.T) {
	v := setupTestVault(t)

	ctx, err := BuildSkillContext(v)
	if err != nil {
		t.Fatalf("BuildSkillContext error = %v", err)
	}

	for _, typ := range ctx.Types {
		if typ.Name == "book" {
			// book type from setupTestVault has properties
			if len(typ.Properties) == 0 {
				t.Error("expected book type to have properties in context")
			}
			return
		}
	}
	t.Error("book type not found in context")
}
