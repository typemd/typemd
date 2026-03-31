package core

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/cucumber/godog"
)

type graphContext struct {
	dotOutput string
}

func initGraphSteps(ctx *godog.ScenarioContext, dc *domainContext) *graphContext {
	gc := &graphContext{}

	ctx.Step(`^I export the graph$`, func() {
		gc.exportGraph(dc, GraphOptions{})
	})
	ctx.Step(`^I export the graph with type filter "([^"]*)"$`, func(typeName string) {
		gc.exportGraph(dc, GraphOptions{Types: []string{typeName}})
	})
	ctx.Step(`^I export the graph without relations$`, func() {
		gc.exportGraph(dc, GraphOptions{NoRelations: true})
	})
	ctx.Step(`^I export the graph without wikilinks$`, func() {
		gc.exportGraph(dc, GraphOptions{NoWikiLinks: true})
	})
	ctx.Step(`^the DOT output should contain "([^"]*)"$`, gc.dotOutputShouldContain)
	ctx.Step(`^the DOT output should have (\d+) nodes?$`, gc.dotOutputShouldHaveNNodes)
	ctx.Step(`^the DOT output should have (\d+) edges?$`, gc.dotOutputShouldHaveNEdges)
	ctx.Step(`^the DOT output should have an edge labeled "([^"]*)"$`, gc.dotOutputShouldHaveEdgeLabeled)
	ctx.Step(`^the DOT output should have an edge labeled "([^"]*)" (\d+) times$`, gc.dotOutputShouldHaveEdgeLabeledNTimes)
	ctx.Step(`^the DOT output should not have an edge labeled "([^"]*)"$`, gc.dotOutputShouldNotHaveEdgeLabeled)

	return gc
}

func (gc *graphContext) exportGraph(dc *domainContext, opts GraphOptions) {
	var buf bytes.Buffer
	dc.lastErr = dc.vault.ExportDOT(&buf, opts)
	gc.dotOutput = buf.String()
}

func (gc *graphContext) dotOutputShouldContain(expected string) error {
	if !strings.Contains(gc.dotOutput, expected) {
		return fmt.Errorf("expected DOT output to contain %q, got:\n%s", expected, gc.dotOutput)
	}
	return nil
}

func (gc *graphContext) dotOutputShouldHaveNNodes(expected int) error {
	count := 0
	for _, line := range strings.Split(gc.dotOutput, "\n") {
		trimmed := strings.TrimSpace(line)
		// A node line has [label=...] but no "->"
		if strings.Contains(trimmed, "[label=") && !strings.Contains(trimmed, "->") {
			count++
		}
	}
	if count != expected {
		return fmt.Errorf("expected %d nodes, got %d\nDOT output:\n%s", expected, count, gc.dotOutput)
	}
	return nil
}

func (gc *graphContext) dotOutputShouldHaveNEdges(expected int) error {
	count := 0
	for _, line := range strings.Split(gc.dotOutput, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(trimmed, "->") {
			count++
		}
	}
	if count != expected {
		return fmt.Errorf("expected %d edges, got %d\nDOT output:\n%s", expected, count, gc.dotOutput)
	}
	return nil
}

func (gc *graphContext) dotOutputShouldHaveEdgeLabeled(label string) error {
	search := fmt.Sprintf(`label="%s"`, label)
	for _, line := range strings.Split(gc.dotOutput, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(trimmed, "->") && strings.Contains(trimmed, search) {
			return nil
		}
	}
	return fmt.Errorf("expected an edge labeled %q in DOT output:\n%s", label, gc.dotOutput)
}

func (gc *graphContext) dotOutputShouldHaveEdgeLabeledNTimes(label string, expected int) error {
	search := fmt.Sprintf(`label="%s"`, label)
	count := 0
	for _, line := range strings.Split(gc.dotOutput, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(trimmed, "->") && strings.Contains(trimmed, search) {
			count++
		}
	}
	if count != expected {
		return fmt.Errorf("expected %d edges labeled %q, got %d\nDOT output:\n%s", expected, label, count, gc.dotOutput)
	}
	return nil
}

func (gc *graphContext) dotOutputShouldNotHaveEdgeLabeled(label string) error {
	search := fmt.Sprintf(`label="%s"`, label)
	for _, line := range strings.Split(gc.dotOutput, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(trimmed, "->") && strings.Contains(trimmed, search) {
			return fmt.Errorf("expected no edge labeled %q, but found one in DOT output:\n%s", label, gc.dotOutput)
		}
	}
	return nil
}
