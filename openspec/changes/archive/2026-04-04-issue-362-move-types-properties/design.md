## Context

typemd stores type schemas at `.typemd/types/<name>/schema.yaml` and shared properties at `.typemd/properties.yaml`. The `objects/` directory is already at vault root. This inconsistency makes schema files hard to discover and easy to exclude from version control.

Current path accessors:
- `Vault.TypesDir()` → `filepath.Join(v.Dir(), "types")` → `.typemd/types/`
- `Vault.SharedPropertiesPath()` → `filepath.Join(v.Dir(), "properties.yaml")` → `.typemd/properties.yaml`
- `LocalObjectRepository.sharedPropertiesPath()` → `filepath.Join(r.root, ".typemd", "properties.yaml")`

All file watchers and consumers use these accessor methods, not hardcoded paths.

## Goals / Non-Goals

**Goals:**
- Move `types/` and `properties/` to vault root
- Auto-migrate existing vaults on `Vault.Open()`
- Create directories at new locations during `tmd init`
- Error clearly when both old and new paths exist

**Non-Goals:**
- Splitting `properties.yaml` into per-property files (issue #363)
- Moving `config.yaml` or other `.typemd/` internal files
- Changing the schema format or structure

## Decisions

### D1: Path accessors point to vault root

Change `TypesDir()` to `filepath.Join(v.Root, "types")` and `SharedPropertiesPath()` to `filepath.Join(v.Root, "properties", "properties.yaml")`.

**Rationale**: All consumers already use these accessors. Changing the return value is sufficient — no caller hardcodes `.typemd/types/`.

**Alternative**: Add a config flag to select old/new paths. Rejected — adds complexity for a one-time migration.

### D2: Migration runs in Vault.Open() before reconciliation

Add a `migrateDirectoryLayout()` call at the start of `Open()`, before SQLite and reconciliation setup. This uses `os.Rename()` for atomic directory moves.

**Rationale**: Migration must complete before any path accessor is used. Running before `Open()`'s existing logic ensures all subsequent operations see new paths.

**Alternative**: Separate `tmd migrate` command only. Rejected — silent breakage on vault open is worse than auto-migration.

### D3: Properties directory wrapper

Shared properties move from `.typemd/properties.yaml` to `properties/properties.yaml` (inside a new `properties/` directory at root).

**Rationale**: Prepares the structure for issue #363 (per-property files). A flat `properties.yaml` at root would require another move later.

### D4: Conflict detection with clear error

If both `.typemd/types/` and `types/` exist (or both property paths), return an error requiring manual resolution instead of guessing.

**Rationale**: Data safety. Silent merging or overwriting could lose user work.

## Risks / Trade-offs

- **[Risk] Existing scripts hardcode `.typemd/types/`** → Mitigated by: this is a documented breaking change; migration handles the common case automatically.
- **[Risk] `os.Rename()` fails across filesystem boundaries** → Mitigated by: vault root and `.typemd/` are always on the same filesystem.
- **[Risk] Empty `.typemd/types/` after move leaves stale directory** → Mitigated by: migration removes old directories after successful move.
- **[Risk] Partial migration (types moved, properties failed)** → Mitigated by: migration checks all preconditions before moving anything; if a conflict is detected for either path, neither is moved.
