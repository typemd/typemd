package core

import (
	"fmt"
	"strings"

	"github.com/cucumber/godog"
	"gopkg.in/yaml.v3"
)

// ── View config step state ───────────────────────────────────────────────────

type viewConfigContext struct {
	dc             *domainContext
	view           *ViewConfig
	yamlOutput     []byte
	yamlInput      string
	deserializedVC *ViewConfig
}

func newViewConfigContext(dc *domainContext) *viewConfigContext {
	return &viewConfigContext{dc: dc}
}

// ── Given steps ─────────────────────────────────────────────────────────────

func (vc *viewConfigContext) aViewConfigWithLayout(name, layout string) {
	vc.view = &ViewConfig{
		Name:   name,
		Layout: ViewLayout(layout),
	}
}

func (vc *viewConfigContext) theViewHasFilterPropertyOperatorValue(property, operator, value string) {
	vc.view.Filter = append(vc.view.Filter, FilterRule{
		Property: property,
		Operator: operator,
		Value:    value,
	})
}

func (vc *viewConfigContext) theViewHasSortPropertyDirection(property, direction string) {
	vc.view.Sort = append(vc.view.Sort, SortRule{
		Property:  property,
		Direction: direction,
	})
}

func (vc *viewConfigContext) theViewHasGroupByProperty(property string) {
	vc.view.GroupBy = append(vc.view.GroupBy, GroupRule{Property: property})
}

func (vc *viewConfigContext) viewYAMLContent(content *godog.DocString) {
	vc.yamlInput = content.Content
}

// ── When steps ──────────────────────────────────────────────────────────────

func (vc *viewConfigContext) iSerializeTheViewConfigToYAML() error {
	data, err := yaml.Marshal(vc.view)
	vc.yamlOutput = data
	vc.dc.lastErr = err
	return nil
}

func (vc *viewConfigContext) iDeserializeTheViewYAML() error {
	var v ViewConfig
	err := yaml.Unmarshal([]byte(vc.yamlInput), &v)
	vc.deserializedVC = &v
	vc.dc.lastErr = err
	return nil
}

// ── Then steps ──────────────────────────────────────────────────────────────

func (vc *viewConfigContext) theViewNameShouldBe(expected string) error {
	if vc.view.Name != expected {
		return fmt.Errorf("expected view name %q, got %q", expected, vc.view.Name)
	}
	return nil
}

func (vc *viewConfigContext) theViewLayoutShouldBe(expected string) error {
	if string(vc.view.Layout) != expected {
		return fmt.Errorf("expected view layout %q, got %q", expected, vc.view.Layout)
	}
	return nil
}

func (vc *viewConfigContext) theViewShouldHaveNFilterRules(n int) error {
	if len(vc.view.Filter) != n {
		return fmt.Errorf("expected %d filter rules, got %d", n, len(vc.view.Filter))
	}
	return nil
}

func (vc *viewConfigContext) theViewShouldHaveNSortRules(n int) error {
	if len(vc.view.Sort) != n {
		return fmt.Errorf("expected %d sort rules, got %d", n, len(vc.view.Sort))
	}
	return nil
}

func (vc *viewConfigContext) theViewShouldHaveNGroupRules(n int) error {
	if len(vc.view.GroupBy) != n {
		return fmt.Errorf("expected %d group rules, got %d", n, len(vc.view.GroupBy))
	}
	return nil
}

func (vc *viewConfigContext) theViewGroupByPropertyNShouldBe(index int, expected string) error {
	idx := index - 1 // 1-based to 0-based
	if idx < 0 || idx >= len(vc.view.GroupBy) {
		return fmt.Errorf("group rule index %d out of range (have %d rules)", index, len(vc.view.GroupBy))
	}
	if vc.view.GroupBy[idx].Property != expected {
		return fmt.Errorf("expected group rule %d property %q, got %q", index, expected, vc.view.GroupBy[idx].Property)
	}
	return nil
}

func (vc *viewConfigContext) theViewYAMLShouldContain(substr string) error {
	if !strings.Contains(string(vc.yamlOutput), substr) {
		return fmt.Errorf("expected YAML to contain %q, got:\n%s", substr, string(vc.yamlOutput))
	}
	return nil
}

func (vc *viewConfigContext) theViewYAMLShouldNotContain(substr string) error {
	if strings.Contains(string(vc.yamlOutput), substr) {
		return fmt.Errorf("expected YAML NOT to contain %q, got:\n%s", substr, string(vc.yamlOutput))
	}
	return nil
}

func (vc *viewConfigContext) theDeserializedViewNameShouldBe(expected string) error {
	if vc.deserializedVC.Name != expected {
		return fmt.Errorf("expected deserialized name %q, got %q", expected, vc.deserializedVC.Name)
	}
	return nil
}

func (vc *viewConfigContext) theDeserializedViewShouldHaveNFilterRules(n int) error {
	if len(vc.deserializedVC.Filter) != n {
		return fmt.Errorf("expected %d filter rules, got %d", n, len(vc.deserializedVC.Filter))
	}
	return nil
}

func (vc *viewConfigContext) theDeserializedViewShouldHaveNSortRules(n int) error {
	if len(vc.deserializedVC.Sort) != n {
		return fmt.Errorf("expected %d sort rules, got %d", n, len(vc.deserializedVC.Sort))
	}
	return nil
}

func (vc *viewConfigContext) theDeserializedViewShouldHaveNGroupRules(n int) error {
	if len(vc.deserializedVC.GroupBy) != n {
		return fmt.Errorf("expected %d group rules, got %d", n, len(vc.deserializedVC.GroupBy))
	}
	return nil
}

func (vc *viewConfigContext) theDeserializedViewGroupByPropertyNShouldBe(index int, expected string) error {
	idx := index - 1
	if idx < 0 || idx >= len(vc.deserializedVC.GroupBy) {
		return fmt.Errorf("deserialized group rule index %d out of range (have %d rules)", index, len(vc.deserializedVC.GroupBy))
	}
	if vc.deserializedVC.GroupBy[idx].Property != expected {
		return fmt.Errorf("expected deserialized group rule %d property %q, got %q", index, expected, vc.deserializedVC.GroupBy[idx].Property)
	}
	return nil
}

// ── Init ────────────────────────────────────────────────────────────────────

func initViewConfigSteps(ctx *godog.ScenarioContext, dc *domainContext) *viewConfigContext {
	vc := newViewConfigContext(dc)

	// Given
	ctx.Step(`^a view config "([^"]*)" with layout "([^"]*)"$`, vc.aViewConfigWithLayout)
	ctx.Step(`^the view has filter property "([^"]*)" operator "([^"]*)" value "([^"]*)"$`, vc.theViewHasFilterPropertyOperatorValue)
	ctx.Step(`^the view has sort property "([^"]*)" direction "([^"]*)"$`, vc.theViewHasSortPropertyDirection)
	ctx.Step(`^the view has group_by property "([^"]*)"$`, vc.theViewHasGroupByProperty)
	ctx.Step(`^view YAML content:$`, vc.viewYAMLContent)

	// When
	ctx.Step(`^I serialize the view config to YAML$`, vc.iSerializeTheViewConfigToYAML)
	ctx.Step(`^I deserialize the view YAML$`, vc.iDeserializeTheViewYAML)

	// Then
	ctx.Step(`^the view name should be "([^"]*)"$`, vc.theViewNameShouldBe)
	ctx.Step(`^the view layout should be "([^"]*)"$`, vc.theViewLayoutShouldBe)
	ctx.Step(`^the view should have (\d+) filter rules?$`, vc.theViewShouldHaveNFilterRules)
	ctx.Step(`^the view should have (\d+) sort rules?$`, vc.theViewShouldHaveNSortRules)
	ctx.Step(`^the view should have (\d+) group rules?$`, vc.theViewShouldHaveNGroupRules)
	ctx.Step(`^the view group_by property (\d+) should be "([^"]*)"$`, vc.theViewGroupByPropertyNShouldBe)
	ctx.Step(`^the view YAML should contain "([^"]*)"$`, vc.theViewYAMLShouldContain)
	ctx.Step(`^the view YAML should not contain "([^"]*)"$`, vc.theViewYAMLShouldNotContain)
	ctx.Step(`^the deserialized view name should be "([^"]*)"$`, vc.theDeserializedViewNameShouldBe)
	ctx.Step(`^the deserialized view should have (\d+) filter rules?$`, vc.theDeserializedViewShouldHaveNFilterRules)
	ctx.Step(`^the deserialized view should have (\d+) sort rules?$`, vc.theDeserializedViewShouldHaveNSortRules)
	ctx.Step(`^the deserialized view should have (\d+) group rules?$`, vc.theDeserializedViewShouldHaveNGroupRules)
	ctx.Step(`^the deserialized view group_by property (\d+) should be "([^"]*)"$`, vc.theDeserializedViewGroupByPropertyNShouldBe)

	return vc
}
