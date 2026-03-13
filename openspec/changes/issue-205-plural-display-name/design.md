## Context

Type names in typemd are always singular (e.g., `book`, `person`). The TUI group headers and CLI output display these singular names in collection contexts, which reads awkwardly (e.g., "▼ 📚 book (3)").

The `TypeSchema` struct currently has `Name`, `Emoji`, and `Properties`. The TUI uses a `typeGroup` struct with `Name`, `Emoji`, `Objects`, and `Expanded` to render the left panel.

## Goals / Non-Goals

**Goals:**

- Add optional `plural` field to `TypeSchema` for grammatically correct display in collection contexts
- Provide a `PluralName()` accessor with sensible fallback
- Update TUI group headers to use plural form
- Update built-in `tag` type with plural value

**Non-Goals:**

- Automatic pluralization (no `s` suffix fallback — breaks for non-English type names)
- Localization or i18n support
- Plural form in detail view (detail shows a single object, singular is correct)

## Decisions

### Decision 1: Fallback to Name when Plural is not set

When `plural` is not set, `PluralName()` returns `Name` as-is. No automatic `s` suffix.

**Rationale:** Type names can be in any language (e.g., Chinese `筆記`). Appending `s` would produce nonsense. Users who care about plural correctness can set it explicitly.

**Alternative considered:** Append `s` by default (like the original issue suggested). Rejected because it fails for non-English names and irregular English plurals.

### Decision 2: Store plural on typeGroup struct

The `typeGroup` struct in TUI gains a `Plural` field populated from `TypeSchema.PluralName()` during `buildGroups`. The group header render uses `Plural` instead of `Name` for display.

`Name` is kept as-is for internal identification (sorting, state persistence, expand/collapse tracking).

### Decision 3: YAML field placement

The `plural` field sits at the same level as `name` and `emoji` in the type schema YAML:

```yaml
name: book
emoji: 📚
plural: books
properties:
  - name: title
    type: string
```

## Risks / Trade-offs

- [Risk] Existing type schemas won't show plural forms until users add the field → Acceptable: falls back to singular name, no breakage
- [Risk] Built-in `tag` type hardcodes plural → Only one built-in type, manageable
