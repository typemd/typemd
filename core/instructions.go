package core

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

//go:embed skills/*/SKILL.md
var skillsFS embed.FS

// Skill represents an embedded skill with parsed metadata and instructions.
type Skill struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	Instructions string `json:"instructions"`
}

// SkillContextType represents a type summary in the skill context.
type SkillContextType struct {
	Name        string                `json:"name"`
	Emoji       string                `json:"emoji,omitempty"`
	Description string                `json:"description,omitempty"`
	Properties  []SkillContextProperty `json:"properties,omitempty"`
}

// SkillContextProperty represents a property in a type summary.
type SkillContextProperty struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
}

// SkillContext holds the vault context injected into skill output.
type SkillContext struct {
	Types []SkillContextType `json:"types"`
}

// SkillOutput is the JSON output format for a skill with context.
type SkillOutput struct {
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	Instructions string        `json:"instructions"`
	Context      *SkillContext `json:"context,omitempty"`
}

// SkillListEntry is the JSON output format for listing skills.
type SkillListEntry struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// skillFiles lists embedded SKILL.md file paths.
// The order here determines the display order in listing.
var skillFiles = []string{
	"skills/explore/SKILL.md",
	"skills/importer/SKILL.md",
}

var (
	skillsOnce  sync.Once
	skillsCache []Skill
)

// skillFrontmatter is the YAML frontmatter structure in SKILL.md files.
type skillFrontmatter struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

// parseSkillMD parses a SKILL.md file into frontmatter and body.
// If no frontmatter is found, returns empty frontmatter and the entire content as body.
func parseSkillMD(content []byte) (skillFrontmatter, string) {
	s := string(content)

	if !strings.HasPrefix(s, "---\n") {
		return skillFrontmatter{}, s
	}

	end := strings.Index(s[4:], "\n---\n")
	if end < 0 {
		return skillFrontmatter{}, s
	}

	yamlBlock := s[4 : 4+end]
	body := s[4+end+5:] // skip past closing "---\n"

	var fm skillFrontmatter
	if err := yaml.Unmarshal([]byte(yamlBlock), &fm); err != nil {
		return skillFrontmatter{}, s
	}

	return fm, body
}

// loadEmbeddedSkills loads and parses all embedded skills.
func loadEmbeddedSkills() []Skill {
	skillsOnce.Do(func() {
		skillsCache = make([]Skill, 0, len(skillFiles))
		for _, file := range skillFiles {
			data, err := skillsFS.ReadFile(file)
			if err != nil {
				panic(fmt.Sprintf("embedded skill %s unreadable: %v", file, err))
			}
			fm, body := parseSkillMD(data)
			skillsCache = append(skillsCache, Skill{
				Name:         fm.Name,
				Description:  fm.Description,
				Instructions: strings.TrimSpace(body),
			})
		}
	})
	return skillsCache
}

// ListSkills returns all available skill entries.
func ListSkills() []SkillListEntry {
	skills := loadEmbeddedSkills()
	entries := make([]SkillListEntry, len(skills))
	for i, s := range skills {
		entries[i] = SkillListEntry{Name: s.Name, Description: s.Description}
	}
	return entries
}

// GetSkill returns a skill by name, or an error if not found.
func GetSkill(name string) (*Skill, error) {
	for _, s := range loadEmbeddedSkills() {
		if s.Name == name {
			result := s // copy
			return &result, nil
		}
	}
	return nil, fmt.Errorf("unknown skill %q; available: %s", name, availableSkillNames())
}

// GetSkillRaw returns the raw SKILL.md content (with frontmatter) for a skill.
func GetSkillRaw(name string) ([]byte, error) {
	for _, file := range skillFiles {
		data, err := skillsFS.ReadFile(file)
		if err != nil {
			continue
		}
		fm, _ := parseSkillMD(data)
		if fm.Name == name {
			return data, nil
		}
	}
	return nil, fmt.Errorf("unknown skill %q; available: %s", name, availableSkillNames())
}

// skillOverridePath returns the filesystem path for a vault-level skill override.
func skillOverridePath(vaultRoot, name string) string {
	return filepath.Join(vaultRoot, vaultDir, "instructions", name+".md")
}

// GetSkillWithOverride returns a skill, checking vault override first.
// If vaultRoot is empty, returns the embedded skill.
func GetSkillWithOverride(name, vaultRoot string) (*Skill, error) {
	embedded, err := GetSkill(name)
	if err != nil {
		return nil, err
	}

	if vaultRoot == "" {
		return embedded, nil
	}

	data, err := os.ReadFile(skillOverridePath(vaultRoot, name))
	if err != nil {
		return embedded, nil // no override, use embedded
	}

	fm, body := parseSkillMD(data)
	result := &Skill{
		Name:         embedded.Name,
		Description:  embedded.Description,
		Instructions: strings.TrimSpace(body),
	}
	if fm.Name != "" {
		result.Name = fm.Name
	}
	if fm.Description != "" {
		result.Description = fm.Description
	}
	return result, nil
}

// GetSkillRawWithOverride returns raw SKILL.md content, checking vault override first.
func GetSkillRawWithOverride(name, vaultRoot string) ([]byte, error) {
	if vaultRoot != "" {
		data, err := os.ReadFile(skillOverridePath(vaultRoot, name))
		if err == nil {
			return data, nil
		}
	}

	return GetSkillRaw(name)
}

// BuildSkillContext builds the vault context for skill output.
func BuildSkillContext(v *Vault) (*SkillContext, error) {
	typeNames := v.ListTypes()
	ctx := &SkillContext{
		Types: make([]SkillContextType, 0, len(typeNames)),
	}

	for _, name := range typeNames {
		schema, err := v.LoadType(name)
		if err != nil {
			continue
		}

		st := SkillContextType{
			Name:        schema.Name,
			Emoji:       schema.Emoji,
			Description: schema.Description,
		}

		for _, prop := range schema.Properties {
			st.Properties = append(st.Properties, SkillContextProperty{
				Name:        prop.Name,
				Type:        string(prop.Type),
				Description: prop.Description,
			})
		}

		ctx.Types = append(ctx.Types, st)
	}

	return ctx, nil
}

// availableSkillNames returns a comma-separated list of skill names.
func availableSkillNames() string {
	skills := loadEmbeddedSkills()
	names := make([]string, len(skills))
	for i, s := range skills {
		names[i] = s.Name
	}
	return strings.Join(names, ", ")
}
