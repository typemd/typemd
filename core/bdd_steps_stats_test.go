package core

import (
	"fmt"
	"strconv"

	"github.com/cucumber/godog"
)

// ── Stats test context ──────────────────────────────────────────────────────

type statsContext struct {
	vaultStats *VaultStats
	typeStats  *TypeStats
}

func newStatsContext() *statsContext {
	return &statsContext{}
}

// ── Setup steps ─────────────────────────────────────────────────────────────

func (sc *statsContext) aTypeWithACheckboxPropertyExists(dc *domainContext, typeName, propName string) {
	yaml := fmt.Sprintf("name: %s\nproperties:\n  - name: %s\n    type: checkbox\n", typeName, propName)
	mustWriteTypeSchema(dc.vault, typeName, []byte(yaml))
}

func (sc *statsContext) aTypeWithADatePropertyExists(dc *domainContext, typeName, propName string) {
	yaml := fmt.Sprintf("name: %s\nproperties:\n  - name: %s\n    type: date\n", typeName, propName)
	mustWriteTypeSchema(dc.vault, typeName, []byte(yaml))
}

func (sc *statsContext) aTypeWithARelationPropertyExists(dc *domainContext, typeName, propName, target string) {
	yaml := fmt.Sprintf("name: %s\nproperties:\n  - name: %s\n    type: relation\n    target: %s\n", typeName, propName, target)
	mustWriteTypeSchema(dc.vault, typeName, []byte(yaml))
}

func (sc *statsContext) aObjectNamedExistsWithTypedProperty(dc *domainContext, typeName, name, prop, value string) {
	dc.aObjectNamedExists(typeName, name)
	schema, err := dc.vault.LoadType(typeName)
	if err != nil {
		panic(fmt.Sprintf("load type schema: %v", err))
	}
	// Find property type and coerce value
	var typedValue any = value
	for _, p := range schema.Properties {
		if p.Name == prop {
			switch p.Type {
			case "number":
				f, err := strconv.ParseFloat(value, 64)
				if err != nil {
					panic(fmt.Sprintf("parse float: %v", err))
				}
				typedValue = f
			case "checkbox":
				b, err := strconv.ParseBool(value)
				if err != nil {
					panic(fmt.Sprintf("parse bool: %v", err))
				}
				typedValue = b
			}
			break
		}
	}
	if err := dc.vault.SetProperty(dc.currentObject.ID, prop, typedValue); err != nil {
		panic(fmt.Sprintf("SetProperty failed: %v", err))
	}
}

// ── VaultStats steps ────────────────────────────────────────────────────────

func (sc *statsContext) iRequestVaultStats(dc *domainContext) {
	stats, err := dc.vault.VaultStats()
	dc.lastErr = err
	sc.vaultStats = stats
}

func (sc *statsContext) theVaultStatsShouldShowNTypes(expected int) error {
	if sc.vaultStats == nil {
		return fmt.Errorf("vault stats is nil")
	}
	if len(sc.vaultStats.Types) != expected {
		return fmt.Errorf("type count = %d, want %d", len(sc.vaultStats.Types), expected)
	}
	return nil
}

func (sc *statsContext) theVaultStatsTotalShouldBe(expected int) error {
	if sc.vaultStats == nil {
		return fmt.Errorf("vault stats is nil")
	}
	if sc.vaultStats.Total != expected {
		return fmt.Errorf("total = %d, want %d", sc.vaultStats.Total, expected)
	}
	return nil
}

func (sc *statsContext) theVaultStatsForTypeShouldShowCount(typeName string, expected int) error {
	if sc.vaultStats == nil {
		return fmt.Errorf("vault stats is nil")
	}
	for _, ts := range sc.vaultStats.Types {
		if ts.Name == typeName {
			if ts.Count != expected {
				return fmt.Errorf("type %q count = %d, want %d", typeName, ts.Count, expected)
			}
			return nil
		}
	}
	return fmt.Errorf("type %q not found in vault stats", typeName)
}

func (sc *statsContext) theVaultStatsForTypeShouldShowEmoji(typeName, expected string) error {
	if sc.vaultStats == nil {
		return fmt.Errorf("vault stats is nil")
	}
	for _, ts := range sc.vaultStats.Types {
		if ts.Name == typeName {
			if ts.Emoji != expected {
				return fmt.Errorf("type %q emoji = %q, want %q", typeName, ts.Emoji, expected)
			}
			return nil
		}
	}
	return fmt.Errorf("type %q not found in vault stats", typeName)
}

// ── TypeStats steps ─────────────────────────────────────────────────────────

func (sc *statsContext) iRequestTypeStats(dc *domainContext, typeName string) {
	stats, err := dc.vault.TypeStats(typeName)
	dc.lastErr = err
	sc.typeStats = stats
}

func (sc *statsContext) theTypeStatsShouldShowCount(expected int) error {
	if sc.typeStats == nil {
		return fmt.Errorf("type stats is nil")
	}
	if sc.typeStats.Count != expected {
		return fmt.Errorf("count = %d, want %d", sc.typeStats.Count, expected)
	}
	return nil
}

func (sc *statsContext) findProperty(propName string) (*PropertyStats, error) {
	if sc.typeStats == nil {
		return nil, fmt.Errorf("type stats is nil")
	}
	for i, ps := range sc.typeStats.Properties {
		if ps.Name == propName {
			return &sc.typeStats.Properties[i], nil
		}
	}
	return nil, fmt.Errorf("property %q not found in type stats", propName)
}

func (sc *statsContext) theTypeStatsPropertyShouldHaveType(propName, expected string) error {
	ps, err := sc.findProperty(propName)
	if err != nil {
		return err
	}
	if ps.Type != expected {
		return fmt.Errorf("property %q type = %q, want %q", propName, ps.Type, expected)
	}
	return nil
}

func (sc *statsContext) theTypeStatsPropertyShouldShowFilled(propName string, expected int) error {
	ps, err := sc.findProperty(propName)
	if err != nil {
		return err
	}
	if ps.Filled != expected {
		return fmt.Errorf("property %q filled = %d, want %d", propName, ps.Filled, expected)
	}
	return nil
}

// Number stats assertions
func (sc *statsContext) theTypeStatsPropertyNumberAvg(propName string, expected float64) error {
	ps, err := sc.findProperty(propName)
	if err != nil {
		return err
	}
	ns, ok := ps.Stats.(*NumberStats)
	if !ok {
		return fmt.Errorf("property %q stats is not NumberStats (got %T)", propName, ps.Stats)
	}
	if ns.Avg != expected {
		return fmt.Errorf("property %q avg = %v, want %v", propName, ns.Avg, expected)
	}
	return nil
}

func (sc *statsContext) theTypeStatsPropertyNumberMin(propName string, expected float64) error {
	ps, err := sc.findProperty(propName)
	if err != nil {
		return err
	}
	ns, ok := ps.Stats.(*NumberStats)
	if !ok {
		return fmt.Errorf("property %q stats is not NumberStats", propName)
	}
	if ns.Min != expected {
		return fmt.Errorf("property %q min = %v, want %v", propName, ns.Min, expected)
	}
	return nil
}

func (sc *statsContext) theTypeStatsPropertyNumberMax(propName string, expected float64) error {
	ps, err := sc.findProperty(propName)
	if err != nil {
		return err
	}
	ns, ok := ps.Stats.(*NumberStats)
	if !ok {
		return fmt.Errorf("property %q stats is not NumberStats", propName)
	}
	if ns.Max != expected {
		return fmt.Errorf("property %q max = %v, want %v", propName, ns.Max, expected)
	}
	return nil
}

func (sc *statsContext) theTypeStatsPropertyNumberSum(propName string, expected float64) error {
	ps, err := sc.findProperty(propName)
	if err != nil {
		return err
	}
	ns, ok := ps.Stats.(*NumberStats)
	if !ok {
		return fmt.Errorf("property %q stats is not NumberStats", propName)
	}
	if ns.Sum != expected {
		return fmt.Errorf("property %q sum = %v, want %v", propName, ns.Sum, expected)
	}
	return nil
}

// Select stats assertions
func (sc *statsContext) theTypeStatsPropertySelectCount(propName, option string, expected int) error {
	ps, err := sc.findProperty(propName)
	if err != nil {
		return err
	}
	ss, ok := ps.Stats.(*SelectStats)
	if !ok {
		return fmt.Errorf("property %q stats is not SelectStats (got %T)", propName, ps.Stats)
	}
	count, exists := ss.Distribution[option]
	if !exists {
		return fmt.Errorf("property %q option %q not found in distribution", propName, option)
	}
	if count != expected {
		return fmt.Errorf("property %q option %q count = %d, want %d", propName, option, count, expected)
	}
	return nil
}

// Checkbox stats assertions
func (sc *statsContext) theTypeStatsPropertyCheckboxTrueCount(propName string, expected int) error {
	ps, err := sc.findProperty(propName)
	if err != nil {
		return err
	}
	cs, ok := ps.Stats.(*CheckboxStats)
	if !ok {
		return fmt.Errorf("property %q stats is not CheckboxStats (got %T)", propName, ps.Stats)
	}
	if cs.TrueCount != expected {
		return fmt.Errorf("property %q true count = %d, want %d", propName, cs.TrueCount, expected)
	}
	return nil
}

func (sc *statsContext) theTypeStatsPropertyCheckboxFalseCount(propName string, expected int) error {
	ps, err := sc.findProperty(propName)
	if err != nil {
		return err
	}
	cs, ok := ps.Stats.(*CheckboxStats)
	if !ok {
		return fmt.Errorf("property %q stats is not CheckboxStats", propName)
	}
	if cs.FalseCount != expected {
		return fmt.Errorf("property %q false count = %d, want %d", propName, cs.FalseCount, expected)
	}
	return nil
}

// Date stats assertions
func (sc *statsContext) theTypeStatsPropertyDateEarliest(propName, expected string) error {
	ps, err := sc.findProperty(propName)
	if err != nil {
		return err
	}
	ds, ok := ps.Stats.(*DateStats)
	if !ok {
		return fmt.Errorf("property %q stats is not DateStats (got %T)", propName, ps.Stats)
	}
	if ds.Earliest != expected {
		return fmt.Errorf("property %q earliest = %q, want %q", propName, ds.Earliest, expected)
	}
	return nil
}

func (sc *statsContext) theTypeStatsPropertyDateLatest(propName, expected string) error {
	ps, err := sc.findProperty(propName)
	if err != nil {
		return err
	}
	ds, ok := ps.Stats.(*DateStats)
	if !ok {
		return fmt.Errorf("property %q stats is not DateStats", propName)
	}
	if ds.Latest != expected {
		return fmt.Errorf("property %q latest = %q, want %q", propName, ds.Latest, expected)
	}
	return nil
}

// Relation stats assertions
func (sc *statsContext) theTypeStatsPropertyRelationCount(propName string, expected int) error {
	ps, err := sc.findProperty(propName)
	if err != nil {
		return err
	}
	rs, ok := ps.Stats.(*RelationStats)
	if !ok {
		return fmt.Errorf("property %q stats is not RelationStats (got %T)", propName, ps.Stats)
	}
	if rs.Count != expected {
		return fmt.Errorf("property %q relation count = %d, want %d", propName, rs.Count, expected)
	}
	return nil
}

// ── Step registration ───────────────────────────────────────────────────────

func initStatsSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	sc := newStatsContext()

	// Setup steps
	ctx.Step(`^a type "([^"]*)" with a "([^"]*)" checkbox property exists$`, func(typeName, propName string) {
		sc.aTypeWithACheckboxPropertyExists(dc, typeName, propName)
	})
	ctx.Step(`^a type "([^"]*)" with a "([^"]*)" date property exists$`, func(typeName, propName string) {
		sc.aTypeWithADatePropertyExists(dc, typeName, propName)
	})
	ctx.Step(`^a type "([^"]*)" with an "([^"]*)" relation to "([^"]*)" exists$`, func(typeName, propName, target string) {
		sc.aTypeWithARelationPropertyExists(dc, typeName, propName, target)
	})

	// Typed property setup
	ctx.Step(`^a "([^"]*)" object named "([^"]*)" exists with typed property "([^"]*)" set to "([^"]*)"$`, func(typeName, name, prop, value string) {
		sc.aObjectNamedExistsWithTypedProperty(dc, typeName, name, prop, value)
	})

	// VaultStats steps
	ctx.Step(`^I request vault stats$`, func() { sc.iRequestVaultStats(dc) })
	ctx.Step(`^the vault stats should show (\d+) types?$`, sc.theVaultStatsShouldShowNTypes)
	ctx.Step(`^the vault stats total should be (\d+)$`, sc.theVaultStatsTotalShouldBe)
	ctx.Step(`^the vault stats for type "([^"]*)" should show count (\d+)$`, sc.theVaultStatsForTypeShouldShowCount)
	ctx.Step(`^the vault stats for type "([^"]*)" should show emoji "([^"]*)"$`, sc.theVaultStatsForTypeShouldShowEmoji)

	// TypeStats steps
	ctx.Step(`^I request type stats for "([^"]*)"$`, func(typeName string) { sc.iRequestTypeStats(dc, typeName) })
	ctx.Step(`^the type stats should show count (\d+)$`, sc.theTypeStatsShouldShowCount)
	ctx.Step(`^the type stats property "([^"]*)" should have type "([^"]*)"$`, sc.theTypeStatsPropertyShouldHaveType)
	ctx.Step(`^the type stats property "([^"]*)" should show filled (\d+)$`, sc.theTypeStatsPropertyShouldShowFilled)

	// Number
	ctx.Step(`^the type stats property "([^"]*)" number avg should be (\d+(?:\.\d+)?)$`, sc.theTypeStatsPropertyNumberAvg)
	ctx.Step(`^the type stats property "([^"]*)" number min should be (\d+(?:\.\d+)?)$`, sc.theTypeStatsPropertyNumberMin)
	ctx.Step(`^the type stats property "([^"]*)" number max should be (\d+(?:\.\d+)?)$`, sc.theTypeStatsPropertyNumberMax)
	ctx.Step(`^the type stats property "([^"]*)" number sum should be (\d+(?:\.\d+)?)$`, sc.theTypeStatsPropertyNumberSum)

	// Select
	ctx.Step(`^the type stats property "([^"]*)" select "([^"]*)" should have count (\d+)$`, sc.theTypeStatsPropertySelectCount)

	// Checkbox
	ctx.Step(`^the type stats property "([^"]*)" checkbox true count should be (\d+)$`, sc.theTypeStatsPropertyCheckboxTrueCount)
	ctx.Step(`^the type stats property "([^"]*)" checkbox false count should be (\d+)$`, sc.theTypeStatsPropertyCheckboxFalseCount)

	// Date
	ctx.Step(`^the type stats property "([^"]*)" date earliest should be "([^"]*)"$`, sc.theTypeStatsPropertyDateEarliest)
	ctx.Step(`^the type stats property "([^"]*)" date latest should be "([^"]*)"$`, sc.theTypeStatsPropertyDateLatest)

	// Relation
	ctx.Step(`^the type stats property "([^"]*)" relation count should be (\d+)$`, sc.theTypeStatsPropertyRelationCount)
}
