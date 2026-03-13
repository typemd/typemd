## Context

typemd enforces name uniqueness only for the built-in `tag` type via hardcoded logic in `checkTagNameUnique()` (creation time) and `ValidateTagNameUniqueness()` (validation time). Both query SQLite using `json_extract(properties, '$.name')`. No other type can opt into this behavior.

The `TypeSchema` struct currently has `Name`, `Plural`, `Emoji`, `Properties`, and `NameTemplate` fields. Type schemas are loaded from `.typemd/types/*.yaml` with a built-in fallback for `tag`.

## Goals / Non-Goals

**Goals:**

- Allow any type to enforce name uniqueness via `unique: true` in its schema
- Unify the tag uniqueness mechanism into the general system (no more special-casing)
- Enforce uniqueness at creation time (`NewObject`) and detect violations at validation time (`tmd type validate`)

**Non-Goals:**

- Rename enforcement — out of scope for this change (no rename feature exists yet)
- Migration tooling for fixing existing violations — `tmd type validate` will report them, but auto-fix is not in scope
- Database-level UNIQUE constraints — enforcement remains at the application layer via SQLite queries

## Decisions

### 1. Add `Unique` field to `TypeSchema` struct

Add `Unique bool` with YAML tag `yaml:"unique"`. Defaults to `false` (zero value), so all existing types are unaffected.

**Alternative considered:** Property-level uniqueness (e.g., `unique: true` on individual properties). Rejected because the issue specifically targets name uniqueness, and property-level uniqueness adds significant complexity with limited current need.

### 2. Generalize `checkTagNameUnique` → `checkNameUnique`

Replace `checkTagNameUnique(name string)` with `checkNameUnique(typeName, name string)` that accepts any type name. The SQL query stays the same — just parameterize the type:

```sql
SELECT id FROM objects WHERE type = ? AND json_extract(properties, '$.name') = ? LIMIT 1
```

In `NewObject()`, replace the hardcoded `if typeName == TagTypeName` check with a schema-driven check: load the type schema, and if `schema.Unique` is true, call `checkNameUnique(typeName, name)`.

### 3. Generalize `ValidateTagNameUniqueness` → `ValidateNameUniqueness`

Replace the tag-specific validation with a general function that:
1. Loads all type schemas
2. Filters to those with `Unique: true`
3. For each unique type, queries all objects and checks for duplicate `name` values
4. Reports all violations

### 4. Built-in `tag` schema gets `Unique: true`

The `defaultTypes` map entry for `tag` gains `Unique: true`. This preserves the existing tag behavior through the generalized mechanism. The dedicated `tag.go` functions (`checkTagNameUnique`) can be removed.

### 5. Uniqueness compares the `name` property value

Per user clarification, uniqueness checks compare the `name` property (from frontmatter), not the filename prefix. This aligns with how the existing tag uniqueness works.

## Risks / Trade-offs

- **[Performance]** Uniqueness check queries SQLite on every `NewObject()` call for unique types → Acceptable: single indexed query, negligible overhead
- **[Name template + unique]** Types with both name template and `unique: true` may create conflicts (e.g., daily note template produces the same name each day) → Mitigation: this is intentional behavior — users who set `unique: true` understand only one object per name is allowed
- **[Existing violations]** If a user enables `unique: true` on an existing type with duplicate names, `tmd type validate` will report violations but won't auto-fix → Acceptable: validation surfaces the problem, user decides how to resolve
