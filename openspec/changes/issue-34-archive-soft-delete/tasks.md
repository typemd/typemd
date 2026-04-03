## 1. Core: Archived system property

- [x] 1.1 Write BDD scenarios for archived system property (IsArchived, omitempty, frontmatter ordering)
- [x] 1.2 Implement BDD step definitions for archived property scenarios
- [x] 1.3 Register `archived` system property in `system_property.go` and add `IsArchived()` to `Object`
- [x] 1.4 Add unit tests for IsArchived edge cases (nil, wrong type, already false)

## 2. Core: SetArchived service method

- [x] 2.1 Write BDD scenarios for SetArchived (archive, unarchive, idempotent, locked bypass)
- [x] 2.2 Implement BDD step definitions for SetArchived scenarios
- [x] 2.3 Add `SetArchived(id, archived)` to `ObjectService` and expose via `Vault` facade
- [x] 2.4 Add unit tests for SetArchived error cases (object not found)

## 3. Core: Default query exclusion

- [x] 3.1 Write BDD scenarios for default archive exclusion in queries (excluded by default, include-archived opt-in, GetObject unaffected)
- [x] 3.2 Implement BDD step definitions for query exclusion scenarios
- [x] 3.3 Add `IncludeArchived` option to `QueryService.Query()` and inject default filter
- [x] 3.4 Add unit tests for query filter injection edge cases

## 4. CLI: Archive and unarchive commands

- [x] 4.1 Write BDD scenarios for `tmd object archive` and `tmd object unarchive` commands (success, idempotent, output messages)
- [x] 4.2 Implement BDD step definitions for CLI archive scenarios
- [x] 4.3 Create `cmd/archive.go` with `archiveCmd` and `unarchiveCmd` following `lock.go` pattern
- [x] 4.4 Add `--include-archived` flag to `tmd object list` and `tmd object query` commands
- [x] 4.5 Write BDD scenarios for `--include-archived` flag on list and query commands
