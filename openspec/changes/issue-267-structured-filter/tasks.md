## 1. Core: Interface & Helper

- [x] 1.1 Write BDD scenarios for structured query filter (empty filter, type filter, multi-filter, with sort)
- [x] 1.2 Add `TypeFilter()` convenience function to `core/view.go`
- [x] 1.3 Change `ObjectIndex.Query()` interface signature from `filter string` to `filter []FilterRule`
- [x] 1.4 Change `QueryService.Query()` signature from `filter string` to `filter []FilterRule`
- [x] 1.5 Change `Vault.QueryObjects()` signature from `filter string` to `filter []FilterRule`

## 2. Core: SQLiteObjectIndex Implementation

- [x] 2.1 Rewrite `SQLiteObjectIndex.Query()` to iterate `[]FilterRule`, using `type` column special case and `FilterRuleToSQL()` for properties
- [x] 2.2 Add unit tests for `type` column mapping vs `json_extract` property mapping

## 3. Core: Migrate Internal Callers

- [x] 3.1 Update `stats.go` — `VaultStats()` and `TypeStats()` to use `TypeFilter()`
- [x] 3.2 Update `type_schema.go` — `CountObjectsByType()` to use `TypeFilter()`
- [x] 3.3 Update `object_service.go` — name uniqueness check to use `[]FilterRule` with type + name
- [x] 3.4 Update `validate.go` — `ValidateAll()` to use `nil` filter
- [x] 3.5 Update `migrate.go` — `MigrateObjects()` to use `TypeFilter()`

## 4. CLI & TUI: Migrate Callers

- [x] 4.1 Update `cmd/object_list.go` to pass `nil` filter
- [x] 4.2 Update `tui/app.go` to pass `nil` filter
- [x] 4.3 Update `tui/view_mode.go` to use `TypeFilter()` for view queries

## 5. Remove Legacy Code

- [x] 5.1 Delete `cmd/query.go` and remove `queryCmd` registration from `root.go`
- [x] 5.2 Remove `parseFilter()` function and `filterCondition` type from `core/query.go`

## 6. Update Tests

- [x] 6.1 Update BDD step definitions in `core/bdd_steps_query_test.go` and `core/bdd_steps_query_sort_test.go`
- [x] 6.2 Update unit tests in `core/query_test.go` and `core/sqlite_object_index_test.go`
- [x] 6.3 Update TUI tests (`tui/create_test.go`, `tui/create_e2e_test.go`, `tui/app.go` related)
- [x] 6.4 Update `core/sync_test.go` filter string usage
- [x] 6.5 Run full test suite and verify all pass
