## Context

TypeSchema currently has two fields: `Name` (string) and `Properties` ([]Property). Types are identified solely by text names in CLI output, TUI, and other interfaces. The issue requests adding an optional `emoji` field for visual identification.

The change is small and localized — it adds a single optional field to an existing struct with no cross-cutting concerns.

## Goals / Non-Goals

**Goals:**

- Add `Emoji` field to `TypeSchema` struct with YAML tag `emoji`
- Support emoji in built-in default types
- Display emoji in `tmd type show` and `tmd type list` output
- Maintain full backward compatibility with existing YAML files

**Non-Goals:**

- Emoji validation (e.g., checking if the string is actually an emoji) — any string is accepted
- Uniqueness constraints across types
- Using emoji in object file paths or IDs
- TUI integration (will be handled separately when TUI is built out)

## Decisions

### 1. Emoji as a top-level TypeSchema field

**Decision:** Add `Emoji string` as a top-level field on `TypeSchema`, not as a property.

**Rationale:** Emoji is metadata about the type itself, not a data property of objects. It belongs alongside `Name`, not in the `Properties` slice. This keeps it simple and avoids polluting the property system.

**Alternative considered:** Storing emoji as a reserved property name — rejected because it conflates type metadata with object data.

### 2. No validation on emoji content

**Decision:** Accept any string value for `emoji`, including multi-character strings or non-emoji text.

**Rationale:** Emoji detection is complex (Unicode categories, ZWJ sequences, skin tone modifiers). The field is purely cosmetic — bad values cause no harm. Users self-correct quickly when output looks wrong.

### 3. Built-in defaults include emoji

**Decision:** Set emoji on all three built-in default types: book (📚), person (👤), note (📝).

**Rationale:** Provides a good out-of-box experience. Users creating custom types can choose their own emoji or omit it.

## Risks / Trade-offs

- **[Terminal emoji rendering]** → Some terminals render emoji with inconsistent widths, potentially misaligning tabular output. Mitigation: use emoji as a prefix rather than in fixed-width columns.
- **[Empty emoji display]** → When emoji is not set, display logic must handle the empty case gracefully. Mitigation: simply omit emoji prefix when empty.
