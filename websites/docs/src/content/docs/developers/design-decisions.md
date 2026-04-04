---
title: Design Decisions
description: Key architectural decisions and their rationale.
sidebar:
  order: 3
---

This page documents significant design decisions made during TypeMD's development. Each entry explains what was decided, why, and what alternatives were considered. Understanding these decisions helps contributors make consistent choices when extending the system.

## Relations: dual storage

**Decision**: Relations are stored in both YAML frontmatter and the SQLite `relations` table.

Frontmatter is the source of truth for file portability — objects are self-contained Markdown files that can be version-controlled, synced, and edited manually. The database provides fast reverse lookups and relational queries that would be expensive to compute from files.

When the two stores drift (e.g., after manual file edits), opening the vault rebuilds the database from frontmatter.

**Alternatives considered**:
- Database-only storage — rejected because it violates the local-first, files-are-truth principle
- Frontmatter-only with file scanning for reverse lookups — rejected due to O(n) scan cost on every query

### Single-value vs. multiple-value behavior

Single-value relations (e.g., `author`) overwrite on re-link. Multiple-value relations (e.g., `books`) append and reject duplicates. This matches user expectations: "set the author" vs. "add a book."

### Bidirectional via explicit inverse

Bidirectional relations require both sides to declare their inverse property in the type schema. When linking A→B, the system automatically writes the inverse B→A. This keeps both files consistent without requiring users to manually maintain both sides.

## Wiki-links: separate from relations

**Decision**: Wiki-links (`[[type/name-ulid]]`) are stored in a dedicated `wikilinks` table, separate from schema-defined relations.

Relations are structured (defined in type schemas with typed targets and cardinality). Wiki-links are freeform (written inline in markdown body, any object can link to any other). Keeping them separate avoids conflating two fundamentally different linking mechanisms.

### Full object ID, not display name

Wiki-link targets use full object IDs including the ULID suffix (e.g., `[[person/bob-01kk3gqm8zrrbjjwkx90f727y6]]`). This simplifies target resolution to a direct lookup — no ambiguity, no fuzzy matching needed.

**Trade-off**: The syntax is verbose. This trades usability for implementation simplicity and correctness. Future improvements like auto-complete can reduce the friction without changing the underlying format.

## System property registry

**Decision**: System properties (`name`, `description`, `created_at`, `updated_at`) are defined in a `[]SystemProperty` slice in `core/system_property.go`. Helper functions like `IsSystemProperty()` and `SystemPropertyNames()` derive all behavior from this registry.

Before the registry, `name` was the only system property and its handling was scattered across multiple files as special-case logic. Adding `created_at` and `updated_at` (and eventually more) would have multiplied these special cases.

### Why a slice, not a map

Slices preserve insertion order, which matters for frontmatter output ordering (`name` → `description` → `created_at` → `updated_at` → schema properties). A map would require a separate ordering mechanism.

### Why not callbacks for auto-set logic

Each system property has different auto-set behavior (`name` from slug or template, `created_at` from `time.Now()`, `updated_at` on every save, `description` is user-authored). Encoding these as callbacks in the registry adds abstraction without reducing complexity. Instead, `NewObject` and `SaveObject` handle the setting logic directly, while the registry handles identification and validation.

### RFC 3339 with local timezone

Timestamps use `time.Now().Format(time.RFC3339)`, which includes the local timezone offset. UTC was considered but rejected — in a local-first tool, human-readable local times are more useful. The string is stored as-is in YAML frontmatter to avoid round-trip formatting issues with `time.Time` values.

## Name property: in the properties map

**Decision**: `name` is stored in `Object.Properties["name"]` rather than as a dedicated struct field.

Keeping `name` in the properties map means it flows through existing frontmatter read/write paths naturally — no schema changes to the SQLite `properties` JSON column, no special handling in frontmatter parsing. `GetName()` reads from the map with a fallback to `DisplayName()` (derived from filename) for backward compatibility.

**Alternative rejected**: A dedicated `Object.Name` struct field would require parallel storage, special handling in frontmatter parsing, and database schema changes.

### Migration via sync

Existing objects without a `name` property are backfilled during `Vault.Sync()` using the display name derived from the filename. This piggybacks on the existing sync mechanism — zero extra user steps.

## Name templates: in the properties array

**Decision**: Name templates are defined as `- name: name` entries in the type schema's `properties` array, with only the `template` field permitted.

```yaml
properties:
  - name: name
    template: "Journal {{ date:YYYY-MM-DD }}"
  - name: content
    type: string
```

The `properties` array is the natural place for property configuration. A top-level `name_template` field was considered but rejected — it would establish a precedent of scattering property config across the schema.

### Template evaluation at create time only

`NewObject()` evaluates the template when the name argument is empty, writing the result as a static string to the `name` property. The template string is not stored in the object. This is simple, predictable, and allows users to edit the name afterward.

### User-friendly date format

Templates use `{{ date:YYYY-MM-DD }}` syntax instead of Go's reference time format (`2006-01-02`). The template engine converts tokens internally:

| Token | Go equivalent | Example |
|-------|--------------|---------|
| YYYY  | 2006         | 2026    |
| MM    | 01           | 03      |
| DD    | 02           | 14      |
| HH    | 15           | 09      |
| mm    | 04           | 30      |
| ss    | 05           | 00      |

## Property type system: explicit types over auto-detection

**Decision**: `date` and `datetime` are separate property types rather than a single auto-detecting type.

A `date` property always stores `YYYY-MM-DD`; a `datetime` always stores `YYYY-MM-DDTHH:MM:SS`. Users declare intent in the schema. ISO 8601 strings are naturally sortable in SQLite.

**Alternative rejected**: A single `date` type that accepts both formats — makes validation and display inconsistent.

### Options as object arrays

`select` and `multi_select` types use `options: [{value: x, label: X}]` instead of the old `values: [x]` format. The `label` field enables display names that differ from stored values (e.g., `value: in-progress`, `label: In Progress`).

**Alternative rejected**: Parallel `values` and `labels` arrays — error-prone coupling.

### Deferred typed SQLite storage

Properties are stored as a JSON blob in `objects.properties`. A typed `object_properties` table was considered but deferred — it is only valuable when paired with query syntax that can exploit it (e.g., `rating>4`, date ranges). Building it now would force premature decisions about multi-select storage with no consumers.

## Shared properties: single-level `use` references

**Decision**: Each shared property lives in its own file under `properties/<name>.yaml`, with the property name derived from the filename. Type schemas reference them with `use: <name>`, allowing only `pin`, `emoji`, and `description` as overrides.

```yaml
# properties/due_date.yaml
type: date
emoji: 📅

# types/task/schema.yaml
properties:
  - use: due_date
    pin: 1
```

Resolution happens in `LoadType()` — after parsing, each `use` entry is replaced with a fully resolved `Property` copied from the shared definition with overrides applied. Downstream code never sees `use` entries.

**Alternative rejected**: `ref` keyword — avoided to prevent confusion with JSON Schema `$ref`. No inheritance or multi-level composition — `use` is a single-level lookup with no recursion possible.

## TUI session state: object ID, not cursor index

**Decision**: The TUI persists `selectedObjectID` (e.g., `book/clean-code-01jqr...`) rather than a cursor index.

Object IDs are stable across sessions even when objects are added or deleted. A cursor index of 3 could point to a completely different object after changes.

### Fallback strategy

When the saved object no longer exists, the TUI falls back to the first object in the same type group, then to the overall first object. This keeps the user in the same "neighborhood" of their vault.

### Save on quit only

State is written to `.typemd/tui-state.yaml` only when the user quits — not continuously. A crash loses state, which is acceptable for a convenience feature. The simpler approach avoids filesystem overhead.

### Silent failure on load errors

If the state file is missing, corrupt, or contains invalid data, the TUI silently falls back to default startup behavior. State persistence is a convenience feature, not critical — users should never be blocked by a broken state file.
