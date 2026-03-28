package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/cucumber/godog"
)

func initCompletionSteps(ctx *godog.ScenarioContext, cc *cmdContext) {
	var completions []string

	ctx.Step(`^the book type has a relation "([^"]*)" to "([^"]*)"$`, func(relName, target string) error {
		// Overwrite the book schema to include a relation property.
		schemaYAML := fmt.Sprintf(`name: book
emoji: "\U0001F4DA"
plural: books
properties:
  - name: author
    type: string
  - name: rating
    type: number
  - name: %s
    type: relation
    target: %s
`, relName, target)
		schemaDir := cc.vault.TypesDir() + "/book"
		os.MkdirAll(schemaDir, 0755)
		return os.WriteFile(schemaDir+"/schema.yaml", []byte(schemaYAML), 0644)
	})

	ctx.Step(`^I request object ID completions for "([^"]*)"$`, func(toComplete string) error {
		vaultPath = cc.vaultDir
		completions, _ = completeObjectID(toComplete)
		return nil
	})

	ctx.Step(`^I request type name completions for "([^"]*)"$`, func(toComplete string) error {
		vaultPath = cc.vaultDir
		completions, _ = completeTypeName(toComplete)
		return nil
	})

	ctx.Step(`^I request relation name completions for the created book with prefix "([^"]*)"$`, func(toComplete string) error {
		if len(cc.createdObjectIDs) == 0 {
			return fmt.Errorf("no created object IDs")
		}
		vaultPath = cc.vaultDir
		lastID := cc.createdObjectIDs[len(cc.createdObjectIDs)-1]
		completions, _ = completeRelationName(lastID, toComplete)
		return nil
	})

	ctx.Step(`^I request relation name completions for "([^"]*)" with prefix "([^"]*)"$`, func(fromID, toComplete string) error {
		vaultPath = cc.vaultDir
		completions, _ = completeRelationName(fromID, toComplete)
		return nil
	})

	ctx.Step(`^the completions should include "([^"]*)"$`, func(expected string) error {
		for _, c := range completions {
			if c == expected {
				return nil
			}
		}
		return fmt.Errorf("expected completions to include %q, got: %v", expected, completions)
	})

	ctx.Step(`^the completions should include a book object starting with "([^"]*)"$`, func(prefix string) error {
		for _, c := range completions {
			if strings.HasPrefix(c, prefix) {
				return nil
			}
		}
		return fmt.Errorf("expected a completion starting with %q, got: %v", prefix, completions)
	})

	ctx.Step(`^the completions should be empty$`, func() error {
		if len(completions) > 0 {
			return fmt.Errorf("expected empty completions, got: %v", completions)
		}
		return nil
	})
}
