## Why

Shared properties are stored in a single `properties/properties.yaml` file containing all definitions. This is inconsistent with the `types/` directory structure where each type has its own directory and schema file. Splitting into per-property files (`properties/<name>.yaml`) improves discoverability, reduces merge conflicts, and aligns the vault layout convention.

## What Changes

- Replace single `properties/properties.yaml` with individual `properties/<name>.yaml` files
- Each file contains the property definition without the `name` field (derived from filename)
- Remove the `SharedPropertiesFile` wrapper struct (`properties:` key no longer needed)
- Add auto-migration on vault open: detect legacy `properties/properties.yaml` → split into per-file → remove original
- Update `GetSharedProperties()` to scan directory for `*.yaml` files
- Update validation to handle filename-derived names and detect conflicts (legacy file coexisting with per-property files)
- Update file watcher to monitor `properties/` directory for changes

## Capabilities

### New Capabilities

- `per-property-files`: Loading, validating, and migrating shared properties from per-file format (`properties/<name>.yaml`)

### Modified Capabilities

_(none — existing specs like `property-editing` deal with object-level properties, not shared property file format)_

## Impact

- **Core**: `LocalObjectRepository` (read/cache logic), `SharedPropertiesFile` struct removal, `Vault.SharedPropertiesPath()` path helper
- **Validation**: `ValidateSharedProperties()` adapts to filename-derived names; new error for legacy/per-file coexistence
- **CLI**: `cmd/validate.go` file watcher path update
- **Tests**: BDD feature file and unit tests need fixture updates for per-file format
- **Docs**: Vault layout documentation, shared properties reference
- **Migration**: One-time auto-migration on vault open (non-breaking for existing vaults)
