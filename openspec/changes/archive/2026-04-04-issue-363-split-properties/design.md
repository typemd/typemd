## Context

Shared properties are currently stored in a single `properties/properties.yaml` file with a `properties:` wrapper key and a `SharedPropertiesFile` struct. This was the original design when shared properties were introduced. Now that `types/` uses a per-type directory structure (`types/<name>/schema.yaml`), the single-file approach for properties is inconsistent.

Current file format:
```yaml
# properties/properties.yaml
properties:
  - name: due_date
    type: date
    emoji: "📅"
  - name: priority
    type: select
    options: [low, medium, high]
```

Target per-file format:
```yaml
# properties/due_date.yaml
type: date
emoji: "📅"
```

## Goals / Non-Goals

**Goals:**

- Split `properties/properties.yaml` into `properties/<name>.yaml` (one file per property)
- Derive property name from filename, removing the `name` field from file content
- Auto-migrate legacy `properties/properties.yaml` to per-file format on vault open
- Maintain full backward compatibility for `use` references in type schemas

**Non-Goals:**

- Changing how `use` references work in type schemas (resolution uses property name, which stays the same)
- Adding a save/write API for shared properties (they remain read-only, user-edited files)
- Changing the `properties/` directory location (already moved to vault root by #362)

## Decisions

### 1. Filename is the property name

Property name is derived from the YAML filename without extension (e.g., `due_date.yaml` → name `due_date`). The `name` field inside the file is no longer needed or expected. If present, it is ignored (the filename takes precedence).

**Rationale:** Consistent with how type names work (`types/<name>/schema.yaml`). Eliminates the possibility of name/filename mismatch. Simpler file format.

**Alternative considered:** Keep `name` field in file, validate it matches filename. Rejected: adds complexity without benefit — one source of truth is better.

### 2. YAML files are parsed directly as Property struct

Each per-property file is unmarshalled directly into `Property`, then `Name` is set from the filename. The `SharedPropertiesFile` wrapper struct is removed.

**Rationale:** Simplest parsing approach. The wrapper was only needed to hold the `properties:` array key.

### 3. Migration runs in `Vault.Open()` before reconciliation

Migration detects `properties/properties.yaml`, splits it into per-property files, and removes the original — all before reconciliation runs. This matches the existing migration pattern for `types/` directory layout (#362).

**Rationale:** Migration must happen before `GetSharedProperties()` is called (which happens during reconciliation). The existing `migrateOldLayout()` in `vault.go` is the natural place.

**Alternative considered:** Separate migration command (`tmd migrate`). Rejected: existing pattern is auto-migration on open, and it's proven reliable.

### 4. Conflict detection: legacy file + per-property files = error

If `properties/properties.yaml` exists alongside individual `*.yaml` files (other than itself), vault open returns an error requiring manual resolution.

**Rationale:** Ambiguous state — unclear which is authoritative. Same approach used by the types migration in #362.

### 5. Non-YAML files in `properties/` are silently ignored

When scanning the directory, only `*.yaml` files are loaded. Other files (`.DS_Store`, `.gitkeep`, etc.) are ignored.

**Rationale:** Defensive and pragmatic. Users may have non-property files in the directory.

### 6. Replace `SharedPropertiesPath()` with `SharedPropertiesDir()`

`Vault.SharedPropertiesPath()` and `LocalObjectRepository.sharedPropertiesPath()` currently return the single-file path. Replace with directory-based helpers:
- `SharedPropertiesDir()` — returns `properties/` directory path
- `SharedPropertyPath(name)` — returns `properties/<name>.yaml` (for future use if write support is added)

**Rationale:** The single-file path is no longer meaningful after migration. The directory is the entry point for scanning.

## Risks / Trade-offs

- **[Risk] Property name with invalid filename characters** → Validate property names against filename-safe pattern during loading. Current property names are simple identifiers (`due_date`, `priority`), so this is unlikely but worth guarding.
- **[Risk] Empty properties directory after migration of empty file** → If `properties.yaml` has no properties defined, migration removes it and creates no per-property files. `GetSharedProperties()` handles empty directory gracefully (returns nil).
- **[Risk] File watcher needs directory-level monitoring** → `cmd/validate.go` currently watches the single file. Switch to `watchDirRecursive` on `properties/` directory — same pattern already used for `types/`.
- **[Trade-off] Removing `SharedPropertiesFile` struct is a breaking change for any external code** → Acceptable: this is an internal struct, not part of any public API.
