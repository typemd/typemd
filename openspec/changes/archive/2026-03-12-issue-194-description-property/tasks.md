## 1. BDD Scenarios

- [x] 1.1 Update `core/features/system_property.feature` — update registry scenarios for 4-property order (name, description, created_at, updated_at) and add `IsSystemProperty("description")` scenario
- [x] 1.2 Add BDD scenarios for description behavior — optional field, not auto-populated on creation, editable, not added during sync
- [x] 1.3 Add BDD scenario for frontmatter ordering with description present

## 2. Core Implementation

- [x] 2.1 Add `DescriptionProperty` constant and registry entry in `core/system_property.go` — insert between name and created_at
- [x] 2.2 Update BDD step definitions if needed to support new scenarios

## 3. Unit Tests

- [x] 3.1 Update `core/system_property_test.go` — update registry order assertions, add `IsSystemProperty("description")` test, add validation tests for schema and shared properties rejecting `description`
- [x] 3.2 Add unit tests for description in frontmatter ordering (OrderedPropKeys with description present and absent)
- [x] 3.3 Add unit test verifying NewObject does not include description in properties

## 4. Verify

- [x] 4.1 Run `go test ./...` and confirm all tests pass
- [x] 4.2 Run `go build ./...` and confirm clean build
