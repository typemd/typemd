package core

import (
	"fmt"

	"github.com/cucumber/godog"
)

// ── Property emoji steps ────────────────────────────────────────────────────

func (dc *domainContext) aTypeSchemaWithPropertyHavingEmoji(typeName, propName, emoji string) {
	schema := fmt.Sprintf("name: %s\nproperties:\n  - name: %s\n    type: string\n    emoji: %s\n", typeName, propName, emoji)
	mustWriteTypeSchema(dc.vault, typeName, []byte(schema))
}

func (dc *domainContext) aTypeSchemaWithPropertiesHavingUniqueEmojis(typeName string) {
	schema := fmt.Sprintf(`name: %s
properties:
  - name: title
    type: string
    emoji: 📖
  - name: rating
    type: number
    emoji: ⭐
`, typeName)
	mustWriteTypeSchema(dc.vault, typeName, []byte(schema))
}

func (dc *domainContext) aTypeSchemaWithPropertiesHavingDuplicateEmojis(typeName string) {
	schema := fmt.Sprintf(`name: %s
properties:
  - name: title
    type: string
    emoji: 👤
  - name: author
    type: string
    emoji: 👤
`, typeName)
	mustWriteTypeSchema(dc.vault, typeName, []byte(schema))
}

func (dc *domainContext) aTypeSchemaWithSomePropertiesMissingEmojis(typeName string) {
	schema := fmt.Sprintf(`name: %s
properties:
  - name: title
    type: string
  - name: author
    type: string
  - name: rating
    type: number
    emoji: ⭐
`, typeName)
	mustWriteTypeSchema(dc.vault, typeName, []byte(schema))
}

func initPropertyEmojiSteps(ctx *godog.ScenarioContext, dc *domainContext) {
	ctx.Step(`^a type schema "([^"]*)" with property "([^"]*)" having emoji "([^"]*)"$`, dc.aTypeSchemaWithPropertyHavingEmoji)
	ctx.Step(`^a type schema "([^"]*)" with properties having unique emojis$`, dc.aTypeSchemaWithPropertiesHavingUniqueEmojis)
	ctx.Step(`^a type schema "([^"]*)" with properties having duplicate emojis$`, dc.aTypeSchemaWithPropertiesHavingDuplicateEmojis)
	ctx.Step(`^a type schema "([^"]*)" with some properties missing emojis$`, dc.aTypeSchemaWithSomePropertiesMissingEmojis)
}
