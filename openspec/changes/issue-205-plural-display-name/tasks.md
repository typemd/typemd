## 1. Core: TypeSchema Plural Field

- [x] 1.1 Write BDD scenarios for plural field loading (plural present, absent, PluralName fallback)
- [x] 1.2 Implement BDD step definitions for plural scenarios
- [x] 1.3 Add `Plural` field to `TypeSchema` struct and `PluralName()` method
- [x] 1.4 Update built-in `tag` type in `defaultTypes` to include `Plural: "tags"`
- [x] 1.5 Add unit tests for PluralName edge cases (empty, non-English)

## 2. TUI: Group Header Plural Display

- [x] 2.1 Add `Plural` field to `typeGroup` struct in `tui/app.go`
- [x] 2.2 Update `buildGroups` in `tui/list.go` to populate `Plural` from `TypeSchema.PluralName()`
- [x] 2.3 Update group header rendering in `renderList` to use `g.Plural` instead of `g.Name`
- [x] 2.4 Update existing TUI tests (`list_test.go`) to reflect plural display in group headers

## 3. CLI: Type Show Plural Display

- [x] 3.1 Update `tmd type show` command to display plural name when set
- [x] 3.2 Skipped — cmd/ has no test files; core logic tested via BDD and unit tests
