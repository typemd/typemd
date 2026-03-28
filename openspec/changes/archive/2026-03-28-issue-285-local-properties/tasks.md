## 1. Core: IsLocal field and BuildDisplayProperties

- [x] 1.1 Write BDD scenarios for local property identification (extra property is local, schema property is not local, system property is not local, no-schema case)
- [x] 1.2 Implement BDD step definitions for local property identification scenarios
- [x] 1.3 Add `IsLocal bool` field to `DisplayProperty` in `display.go` and set it in `BuildDisplayProperties()` in `query_service.go`
- [x] 1.4 Add unit tests for edge cases (object with only local properties, object with no properties, built-in types with extra frontmatter)

## 2. TUI: Property panel separator and read-only local properties

- [x] 2.1 Update `isPropertyEditable()` to return false for `IsLocal` properties
- [x] 2.2 Add separator rendering in `propEditor.Render()` before the first local property
- [x] 2.3 Add unit tests for separator rendering (with local, without local, only local)

## 3. CLI: Separate local properties section in `tmd object show`

- [x] 3.1 Write BDD scenarios for CLI local property display (separate section, no section when none, format)
- [x] 3.2 Implement BDD step definitions for CLI local property scenarios
- [x] 3.3 Update `cmd/show.go` to split properties into schema and local sections
