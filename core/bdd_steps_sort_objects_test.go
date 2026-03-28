package core

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cucumber/godog"
)

// ── Sort objects step state ───────────────────────────────────────────────

type sortObjectsContext struct {
	dc      *domainContext
	objects []*Object
	sorted  []*Object
}

func newSortObjectsContext(dc *domainContext) *sortObjectsContext {
	return &sortObjectsContext{dc: dc}
}

// ── Given steps ───────────────────────────────────────────────────────────

func (sc *sortObjectsContext) inMemoryObjects(table *godog.Table) {
	sc.objects = nil
	sc.sorted = nil
	headers := table.Rows[0]

	for _, row := range table.Rows[1:] {
		obj := &Object{
			Properties: make(map[string]any),
		}

		for ci, cell := range row.Cells {
			colName := headers.Cells[ci].Value
			value := strings.TrimSpace(cell.Value)

			switch colName {
			case "type":
				obj.Type = value
			case "name":
				obj.Properties["name"] = value
			default:
				if value == "" {
					// Leave property unset for empty values (nil in map)
					continue
				}
				// Try to parse as number for numeric properties
				if f, err := strconv.ParseFloat(value, 64); err == nil {
					obj.Properties[colName] = f
				} else {
					obj.Properties[colName] = value
				}
			}
		}

		sc.objects = append(sc.objects, obj)
	}
}

// ── When steps ────────────────────────────────────────────────────────────

func (sc *sortObjectsContext) iSortObjectsBy(property, direction string) {
	// Copy the slice so we preserve the original order
	sc.sorted = make([]*Object, len(sc.objects))
	copy(sc.sorted, sc.objects)
	SortObjects(sc.sorted, []SortRule{{Property: property, Direction: direction}})
}

// ── Then steps ────────────────────────────────────────────────────────────

func (sc *sortObjectsContext) theSortedObjectNamesShouldBe(expected string) error {
	expectedNames := strings.Split(expected, ",")
	if len(sc.sorted) != len(expectedNames) {
		return fmt.Errorf("expected %d objects, got %d", len(expectedNames), len(sc.sorted))
	}
	for i, obj := range sc.sorted {
		name := obj.GetName()
		want := strings.TrimSpace(expectedNames[i])
		if name != want {
			return fmt.Errorf("position %d: expected %q, got %q (full order: %s)",
				i, want, name, sortedNames(sc.sorted))
		}
	}
	return nil
}

func (sc *sortObjectsContext) theObjectShouldBeLastInSortedResults(name string) error {
	if len(sc.sorted) == 0 {
		return fmt.Errorf("sorted results are empty")
	}
	last := sc.sorted[len(sc.sorted)-1]
	if last.GetName() != name {
		return fmt.Errorf("expected last object to be %q, got %q (full order: %s)",
			name, last.GetName(), sortedNames(sc.sorted))
	}
	return nil
}

// sortedNames returns a comma-separated list of object names for debugging.
func sortedNames(objects []*Object) string {
	names := make([]string, len(objects))
	for i, obj := range objects {
		names[i] = obj.GetName()
	}
	return strings.Join(names, ",")
}

// ── Init ──────────────────────────────────────────────────────────────────

func initSortObjectsSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	sc := newSortObjectsContext(dc)

	// Given
	ctx.Step(`^in-memory objects:$`, sc.inMemoryObjects)

	// When
	ctx.Step(`^I sort objects by "([^"]*)" "([^"]*)"$`, sc.iSortObjectsBy)

	// Then
	ctx.Step(`^the sorted object names should be "([^"]*)"$`, sc.theSortedObjectNamesShouldBe)
	ctx.Step(`^the object "([^"]*)" should be last in sorted results$`, sc.theObjectShouldBeLastInSortedResults)
}
