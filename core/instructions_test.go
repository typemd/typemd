package core

import (
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
	for _, name := range []string{"explore", "importer", "guide"} {
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
