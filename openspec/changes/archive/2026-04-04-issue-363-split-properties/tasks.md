## 1. Core: Per-property file loading

- [x] 1.1 Write BDD scenarios for loading per-property files (single file, multiple files, empty dir, missing dir, non-YAML ignored, name field ignored)
- [x] 1.2 Implement BDD step definitions for per-property loading scenarios
- [x] 1.3 Replace `SharedPropertiesFile` struct and `GetSharedProperties()` to scan `properties/` directory for `*.yaml` files
- [x] 1.4 Replace `Vault.SharedPropertiesPath()` with `SharedPropertiesDir()` and update `Vault.Init()` to ensure directory exists
- [x] 1.5 Update unit tests for new loading logic (caching, malformed YAML, empty files)

## 2. Core: Validation

- [x] 2.1 Write BDD scenarios for per-property validation (reserved name, invalid type, select without options)
- [x] 2.2 Implement step definitions and update `ValidateSharedProperties()` for filename-derived names
- [x] 2.3 Update unit tests for validation edge cases

## 3. Core: Migration

- [x] 3.1 Write BDD scenarios for migration (successful split, empty file, conflict detection, field preservation)
- [x] 3.2 Implement step definitions for migration scenarios
- [x] 3.3 Implement `migrateSharedProperties()` in `migrate.go` — detect legacy file, split into per-property files, remove original
- [x] 3.4 Wire migration into `Vault.Open()` before reconciliation
- [x] 3.5 Update existing migration tests to account for new migration step

## 4. CLI: File watcher update

- [x] 4.1 Update `cmd/validate.go` `addWatchPaths()` to watch `properties/` directory recursively instead of single file

## 5. Docs and references

- [x] 5.1 Update CLAUDE.md data model section for per-property file format
- [x] 5.2 Update example vault (`examples/book-vault`) to use per-property files
- [x] 5.3 Update embedded skills referencing `properties/properties.yaml`
