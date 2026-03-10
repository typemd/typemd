## Context

typemd type schemas define all properties inline. Common properties like `due_date` or `priority` are duplicated across types. The `Property` struct and `LoadType()` currently only handle inline definitions. Validation lives in `ValidateSchema()`.

## Goals / Non-Goals

**Goals:**
- Define shared properties in `.typemd/properties.yaml` using the same Property format as type schemas
- Allow type schemas to reference shared properties via `use: <name>` with `pin`/`emoji` overrides
- Resolve `use` entries at load time so downstream code sees fully resolved TypeSchema
- Validate shared properties file and `use` references

**Non-Goals:**
- Inheritance or composition beyond single-level `use` references
- System properties (#9) — separate concern, different mechanism
- CLI commands for managing shared properties (users edit YAML directly)
- Migration tooling for existing type schemas

## Decisions

### 1. File location: `.typemd/properties.yaml`

Shared property definitions live in `.typemd/properties.yaml` alongside the `types/` directory. The file is optional — vaults without it work as before.

**Format** uses the same `properties` array as type schemas:

```yaml
properties:
  - name: due_date
    type: date
    emoji: 📅
  - name: priority
    type: select
    options:
      - value: high
      - value: medium
      - value: low
```

**Alternative considered**: `.typemd/types/_shared.yaml` — rejected because shared properties aren't a type.

### 2. Reference keyword: `use`

Type schemas reference shared properties with `use: <name>`:

```yaml
properties:
  - name: title
    type: string
  - use: due_date
    pin: 1
    emoji: 🗓️
  - use: priority
```

`use` and `name` are mutually exclusive on a property entry. Only `pin` and `emoji` are allowed as overrides on `use` entries.

**Alternative considered**: `ref` — rejected to avoid confusion with JSON Schema `$ref`.

### 3. Resolution in LoadType()

`LoadType()` gains a resolution step after parsing YAML:

1. Load `.typemd/properties.yaml` (cached on Vault)
2. For each `use` entry in the type schema, look up the shared property by name
3. Create a resolved Property by copying all fields from the shared definition, then applying `pin`/`emoji` overrides
4. Replace the `use` entry with the resolved Property in the schema's Properties slice

After resolution, the `Use` field is cleared — downstream code never sees it.

### 4. Shared properties loading via Vault

Add `LoadSharedProperties()` method on Vault that:
- Reads `.typemd/properties.yaml`
- Returns `[]Property` (empty slice if file doesn't exist)
- Caches result for reuse across multiple `LoadType()` calls

### 5. Validation rules

All validation happens in `ValidateSchema()` (extended) and a new `ValidateSharedProperties()`:

| Rule | Location |
|------|----------|
| Shared properties file: duplicate names | `ValidateSharedProperties()` |
| Shared properties: standard property validation (types, options, etc.) | `ValidateSharedProperties()` |
| `use` references non-existent shared property | `ValidateSchema()` |
| `use` entry has fields other than `pin`, `emoji` | `ValidateSchema()` |
| `name` property conflicts with shared property name | `ValidateSchema()` |
| Resolved properties have duplicate names | `ValidateSchema()` (existing duplicate check, post-resolution) |

## Risks / Trade-offs

- **Circular complexity**: None — `use` is single-level lookup, no recursion possible.
- **Load order**: `properties.yaml` must be loaded before any type schema. Caching on Vault handles this naturally.
- **Breaking change**: None — `use` is additive. Existing schemas without `use` work unchanged. Missing `properties.yaml` is valid.
