package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cucumber/godog"
)

// ── Skill instructions steps ────────────────────────────────────────────────

type instructionsContext struct {
	dc           *domainContext
	skills       []SkillListEntry
	skill        *Skill
	rawContent   []byte
	skillContext *SkillContext
}

func newInstructionsContext(dc *domainContext) *instructionsContext {
	return &instructionsContext{dc: dc}
}

func (ic *instructionsContext) iListAvailableSkills() {
	ic.skills = ListSkills()
}

func (ic *instructionsContext) iShouldGetNSkills(expected int) error {
	if len(ic.skills) != expected {
		return fmt.Errorf("expected %d skills, got %d", expected, len(ic.skills))
	}
	return nil
}

func (ic *instructionsContext) theSkillsShouldInclude(names string) error {
	nameList := strings.Split(names, ", ")
	have := make(map[string]bool)
	for _, s := range ic.skills {
		have[s.Name] = true
	}
	for _, name := range nameList {
		name = strings.Trim(name, "\"")
		if !have[name] {
			return fmt.Errorf("skill %q not found, have: %v", name, have)
		}
	}
	return nil
}

func (ic *instructionsContext) eachSkillShouldHaveNameAndDescription() error {
	for _, s := range ic.skills {
		if s.Name == "" {
			return fmt.Errorf("skill has empty name")
		}
		if s.Description == "" {
			return fmt.Errorf("skill %q has empty description", s.Name)
		}
	}
	return nil
}

func (ic *instructionsContext) iGetTheSkill(name string) {
	skill, err := GetSkill(name)
	ic.skill = skill
	ic.dc.lastErr = err
}

func (ic *instructionsContext) theSkillNameShouldBe(expected string) error {
	if ic.skill == nil {
		return fmt.Errorf("no skill loaded")
	}
	if ic.skill.Name != expected {
		return fmt.Errorf("expected skill name %q, got %q", expected, ic.skill.Name)
	}
	return nil
}

func (ic *instructionsContext) theSkillDescriptionShouldNotBeEmpty() error {
	if ic.skill == nil {
		return fmt.Errorf("no skill loaded")
	}
	if ic.skill.Description == "" {
		return fmt.Errorf("skill description is empty")
	}
	return nil
}

func (ic *instructionsContext) theSkillInstructionsShouldNotBeEmpty() error {
	if ic.skill == nil {
		return fmt.Errorf("no skill loaded")
	}
	if ic.skill.Instructions == "" {
		return fmt.Errorf("skill instructions are empty")
	}
	return nil
}

func (ic *instructionsContext) theSkillInstructionsShouldNotContain(substr string) error {
	if ic.skill == nil {
		return fmt.Errorf("no skill loaded")
	}
	if strings.Contains(ic.skill.Instructions, substr) {
		return fmt.Errorf("skill instructions should not contain %q", substr)
	}
	return nil
}

func (ic *instructionsContext) theSkillInstructionsShouldContain(substr string) error {
	if ic.skill == nil {
		return fmt.Errorf("no skill loaded")
	}
	if !strings.Contains(ic.skill.Instructions, substr) {
		return fmt.Errorf("expected skill instructions to contain %q", substr)
	}
	return nil
}

func (ic *instructionsContext) iGetTheRawContentOfSkill(name string) {
	data, err := GetSkillRaw(name)
	ic.rawContent = data
	ic.dc.lastErr = err
}

func (ic *instructionsContext) theRawContentShouldStartWith(prefix string) error {
	s := string(ic.rawContent)
	if !strings.HasPrefix(s, prefix) {
		preview := s
		if len(preview) > 40 {
			preview = preview[:40]
		}
		return fmt.Errorf("expected raw content to start with %q, got %q", prefix, preview)
	}
	return nil
}

func (ic *instructionsContext) theRawContentShouldContain(substr string) error {
	if !strings.Contains(string(ic.rawContent), substr) {
		return fmt.Errorf("expected raw content to contain %q", substr)
	}
	return nil
}

func (ic *instructionsContext) aTypeExistsWithEmojiAndDescription(typeName, emoji, description string) {
	schema := fmt.Sprintf("name: %s\nemoji: %q\ndescription: %q\nproperties: []\n", typeName, emoji, description)
	os.WriteFile(filepath.Join(ic.dc.vault.TypesDir(), typeName+".yaml"), []byte(schema), 0644)
}

func (ic *instructionsContext) theTypeHasAPropertyOfType(typeName, propName, propType string) {
	// Read existing schema file and add property
	schemaPath := filepath.Join(ic.dc.vault.TypesDir(), typeName+".yaml")
	data, _ := os.ReadFile(schemaPath)
	content := string(data)
	content = strings.Replace(content, "properties: []", fmt.Sprintf("properties:\n  - name: %s\n    type: %s", propName, propType), 1)
	os.WriteFile(schemaPath, []byte(content), 0644)
}

func (ic *instructionsContext) iBuildSkillContextFromTheVault() {
	ctx, err := BuildSkillContext(ic.dc.vault)
	ic.skillContext = ctx
	ic.dc.lastErr = err
}

func (ic *instructionsContext) theContextShouldHaveTypes() error {
	if ic.skillContext == nil {
		return fmt.Errorf("no skill context")
	}
	if len(ic.skillContext.Types) == 0 {
		return fmt.Errorf("context has no types")
	}
	return nil
}

func (ic *instructionsContext) theContextTypesShouldInclude(typeName string) error {
	if ic.skillContext == nil {
		return fmt.Errorf("no skill context")
	}
	for _, t := range ic.skillContext.Types {
		if t.Name == typeName {
			return nil
		}
	}
	return fmt.Errorf("type %q not found in context types", typeName)
}

func (ic *instructionsContext) theTypeInContextShouldHaveProperty(typeName, propName string) error {
	if ic.skillContext == nil {
		return fmt.Errorf("no skill context")
	}
	for _, t := range ic.skillContext.Types {
		if t.Name == typeName {
			for _, p := range t.Properties {
				if p.Name == propName {
					return nil
				}
			}
			return fmt.Errorf("property %q not found in type %q", propName, typeName)
		}
	}
	return fmt.Errorf("type %q not found in context types", typeName)
}

func (ic *instructionsContext) aVaultOverrideExistsForSkillWithBody(name, body string) {
	dir := filepath.Join(ic.dc.vault.Dir(), "instructions")
	os.MkdirAll(dir, 0755)
	content := fmt.Sprintf("---\nname: %s\ndescription: Custom override\n---\n\n%s\n", name, body)
	os.WriteFile(filepath.Join(dir, name+".md"), []byte(content), 0644)
}

func (ic *instructionsContext) aVaultOverrideExistsForSkillWithoutFrontmatter(name string) {
	dir := filepath.Join(ic.dc.vault.Dir(), "instructions")
	os.MkdirAll(dir, 0755)
	content := "Override content without frontmatter\n"
	os.WriteFile(filepath.Join(dir, name+".md"), []byte(content), 0644)
}

func (ic *instructionsContext) iGetTheSkillWithVaultOverride(name string) {
	skill, err := GetSkillWithOverride(name, ic.dc.rootDir)
	ic.skill = skill
	ic.dc.lastErr = err
}

func initInstructionsSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	ic := newInstructionsContext(dc)

	ctx.Step(`^I list available skills$`, ic.iListAvailableSkills)
	ctx.Step(`^I should get (\d+) skills$`, ic.iShouldGetNSkills)
	ctx.Step(`^the skills should include (.+)$`, ic.theSkillsShouldInclude)
	ctx.Step(`^each skill should have a name and description$`, ic.eachSkillShouldHaveNameAndDescription)

	ctx.Step(`^I get the skill "([^"]*)"$`, ic.iGetTheSkill)
	ctx.Step(`^the skill name should be "([^"]*)"$`, ic.theSkillNameShouldBe)
	ctx.Step(`^the skill description should not be empty$`, ic.theSkillDescriptionShouldNotBeEmpty)
	ctx.Step(`^the skill instructions should not be empty$`, ic.theSkillInstructionsShouldNotBeEmpty)
	ctx.Step(`^the skill instructions should not contain "([^"]*)"$`, ic.theSkillInstructionsShouldNotContain)
	ctx.Step(`^the skill instructions should contain "([^"]*)"$`, ic.theSkillInstructionsShouldContain)

	ctx.Step(`^I get the raw content of skill "([^"]*)"$`, ic.iGetTheRawContentOfSkill)
	ctx.Step(`^the raw content should start with "([^"]*)"$`, ic.theRawContentShouldStartWith)
	ctx.Step(`^the raw content should contain "([^"]*)"$`, ic.theRawContentShouldContain)

	ctx.Step(`^a type "([^"]*)" exists with emoji "([^"]*)" and description "([^"]*)"$`, ic.aTypeExistsWithEmojiAndDescription)
	ctx.Step(`^the type "([^"]*)" has a property "([^"]*)" of type "([^"]*)"$`, ic.theTypeHasAPropertyOfType)
	ctx.Step(`^I build skill context from the vault$`, ic.iBuildSkillContextFromTheVault)
	ctx.Step(`^the context should have types$`, ic.theContextShouldHaveTypes)
	ctx.Step(`^the context types should include "([^"]*)"$`, ic.theContextTypesShouldInclude)
	ctx.Step(`^the type "([^"]*)" in context should have property "([^"]*)"$`, ic.theTypeInContextShouldHaveProperty)

	ctx.Step(`^a vault override exists for skill "([^"]*)" with body "([^"]*)"$`, ic.aVaultOverrideExistsForSkillWithBody)
	ctx.Step(`^a vault override exists for skill "([^"]*)" without frontmatter$`, ic.aVaultOverrideExistsForSkillWithoutFrontmatter)
	ctx.Step(`^I get the skill "([^"]*)" with vault override$`, ic.iGetTheSkillWithVaultOverride)
}
