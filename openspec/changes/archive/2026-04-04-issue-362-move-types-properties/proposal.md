## Why

User-editable schema files (`types/` and `properties/`) are currently inside the hidden `.typemd/` directory. This makes them hard to discover, inconsistent with `objects/` (which is already at vault root), and easy to accidentally exclude from version control. Moving them to the vault root establishes a clear principle: **user-editable files at root, internal state in `.typemd/`**.

## What Changes

- **BREAKING**: `TypesDir()` moves from `.typemd/types/` → `types/` at vault root
- **BREAKING**: `SharedPropertiesPath()` moves from `.typemd/properties.yaml` → `properties/properties.yaml` at vault root
- Add auto-migration in `Vault.Open()`: detect old paths → move to new paths
- Update `tmd init` to create `types/` and `properties/` at vault root
- Error if both old and new paths exist simultaneously (require manual resolution)

## Capabilities

### New Capabilities

- `vault-directory-migration`: Auto-migration of `types/` and `properties/` from `.typemd/` to vault root on vault open, with conflict detection and dry-run support

### Modified Capabilities

_(none — existing specs describe behaviors that don't change, only internal paths shift)_

## Impact

- **core/vault.go**: `TypesDir()` and `SharedPropertiesPath()` return new paths
- **core/local_object_repository.go**: `sharedPropertiesPath()` derives from repository root — needs update for new `properties/` directory
- **core/migrate.go**: Add new directory-level migration function
- **core/vault.go `Open()`**: Call migration before reconciliation
- **core/vault.go `Init()`**: Create directories at new locations
- **cmd/**: `tmd init`, `tmd migrate` may need updates
- **File watchers**: Already use `Vault.TypesDir()` — no change needed if accessor is updated
- **Docs, CLAUDE.md**: Path references need updating
