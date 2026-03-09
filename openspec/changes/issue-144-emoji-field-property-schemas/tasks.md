## 1. Property Struct Update

- [x] 1.1 Add `Emoji string \`yaml:"emoji,omitempty"\`` field to `Property` struct in `core/type_schema.go`
- [x] 1.2 Add example emojis to default type properties in `defaultTypes` map

## 2. Schema Validation

- [x] 2.1 Add duplicate property emoji check in `ValidateSchema()` — skip empty emojis, error on duplicates
- [x] 2.2 Write unit tests for emoji uniqueness validation (unique accepted, duplicate rejected, empty skipped)

## 3. BDD Tests

- [x] 3.1 Write Gherkin feature file `core/features/property_emoji.feature` covering spec scenarios
- [x] 3.2 Implement BDD step definitions for property emoji parsing and validation scenarios

## 4. Example Schemas

- [x] 4.1 Add emoji fields to example vault type schemas (`examples/book-vault/.typemd/types/`)

## 5. Verification

- [x] 5.1 Run `go test ./...` and confirm all tests pass
- [x] 5.2 Run `go build ./...` and confirm clean build
