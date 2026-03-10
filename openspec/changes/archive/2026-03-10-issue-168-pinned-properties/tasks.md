## 1. Property Struct Update

- [x] 1.1 Add `Pin int \`yaml:"pin,omitempty"\`` field to `Property` struct in `core/type_schema.go`

## 2. Schema Validation

- [x] 2.1 Add pin value validation in `ValidateSchema()` — reject negative or zero pin values, reject duplicate non-zero pin values
- [x] 2.2 Write unit tests for pin validation (positive accepted, negative rejected, duplicate rejected, unpinned no conflict)

## 3. BDD Tests (Core)

- [x] 3.1 Write Gherkin feature file `core/features/pinned_property.feature` covering spec scenarios for pin field parsing and validation
- [x] 3.2 Implement BDD step definitions for pinned property scenarios

## 4. Display Property Update

- [x] 4.1 Add `Pin int` and `Emoji string` fields to `DisplayProperty` struct in `core/display.go`
- [x] 4.2 Populate `Pin` and `Emoji` from schema in `BuildDisplayProperties()`

## 5. TUI Rendering

- [x] 5.1 Update `renderBody()` in `tui/detail.go` to render pinned properties at top with separator
- [x] 5.2 Update `renderProperties()` in `tui/detail.go` to exclude pinned properties
- [x] 5.3 Pass `displayProps` to `renderBody()` so it has access to pinned property data

## 6. Example Schemas

- [x] 6.1 Add pin fields to example vault type schemas (`examples/book-vault/.typemd/types/`)

## 7. Verification

- [x] 7.1 Run `go test ./...` and confirm all tests pass
- [x] 7.2 Run `go build ./...` and confirm clean build
