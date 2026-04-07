## 1. BDD Scenarios

- [ ] 1.1 Write BDD scenarios for source type in `core/features/type_crud.feature` (load defaults, deletion protection, custom override, properties)

## 2. Implementation

- [ ] 2.1 Add `SourceTypeName = "source"` constant in `core/system_property.go`
- [ ] 2.2 Register source type in `defaultTypes` map in `core/type_schema.go` with emoji, plural, and properties

## 3. Unit Tests

- [ ] 3.1 Add unit tests for source type deletion protection and no-event emission in `core/type_crud_test.go`

## 4. Documentation

- [ ] 4.1 Update CLAUDE.md data model section to mention `source` alongside `tag` and `page`
