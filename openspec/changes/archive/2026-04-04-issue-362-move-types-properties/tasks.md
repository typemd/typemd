## 1. Core: Directory Migration

- [x] 1.1 Write BDD scenarios for vault directory migration (old→new move, conflict detection, skip when not needed, pre-check prevents partial migration)
- [x] 1.2 Implement BDD step definitions for directory migration scenarios
- [x] 1.3 Add `migrateDirectoryLayout()` function in `core/migrate.go` (move types dir, move properties file, conflict detection, pre-check all paths)
- [x] 1.4 Add unit tests for migration edge cases (empty old types dir, properties dir creation, cleanup of old directories)

## 2. Core: Path Accessor Changes

- [x] 2.1 Write BDD scenarios for new path accessors (`TypesDir()` returns root-level, `SharedPropertiesPath()` returns root-level)
- [x] 2.2 Implement BDD step definitions for path accessor scenarios
- [x] 2.3 Update `Vault.TypesDir()` to return `filepath.Join(v.Root, "types")`
- [x] 2.4 Update `Vault.SharedPropertiesPath()` to return `filepath.Join(v.Root, "properties", "properties.yaml")`
- [x] 2.5 Update `LocalObjectRepository.sharedPropertiesPath()` to use new `properties/` directory path
- [x] 2.6 Fix all existing tests that set up type schemas or properties at old paths

## 3. Core: Vault Open Integration

- [x] 3.1 Call `migrateDirectoryLayout()` in `Vault.Open()` before SQLite and reconciliation setup
- [x] 3.2 Verify existing BDD scenarios still pass with migration integrated

## 4. Core: Vault Init Changes

- [x] 4.1 Write BDD scenario for `tmd init` creating root-level directories
- [x] 4.2 Update `Vault.Init()` to create `types/` and `properties/` at vault root instead of `.typemd/types/`
- [x] 4.3 Add unit test verifying `.typemd/types/` is NOT created during init

## 5. Fix Existing Tests and Integration

- [x] 5.1 Run full test suite, fix any remaining failures from path changes
- [x] 5.2 Update file watcher paths if any are hardcoded (verify they use `TypesDir()`)

## 6. Documentation

- [x] 6.1 Update CLAUDE.md data model section with new directory layout
- [x] 6.2 Update docs site references to `.typemd/types/` and `.typemd/properties.yaml`
