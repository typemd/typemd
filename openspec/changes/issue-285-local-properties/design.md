## Context

Objects in typemd can contain frontmatter properties that are not defined in the type schema. These "local properties" are already preserved during load/save (round-trip safe) and are filtered out of the SQLite index during sync (#174). However, the display layer treats them identically to schema-defined properties, making it impossible for users to distinguish managed vs. ad-hoc properties.

Currently, `BuildDisplayProperties()` in `QueryService` builds a flat list of `DisplayProperty` items. Local properties appear after schema-defined properties (handled by `OrderedPropKeys`'s "extra" bucket) but lack any metadata to identify them as local.

## Goals / Non-Goals

**Goals:**

- Add `IsLocal` field to `DisplayProperty` so consumers can identify local properties
- Render a visual separator in the TUI property panel before the first local property
- Show a separate "Local Properties" section in CLI `tmd object show`
- Make local properties read-only in TUI (cursor skips them, styled as dim)

**Non-Goals:**

- Local property editing (future issue)
- Local properties in view mode (table/list layouts)
- Local property indexing in SQLite
- Local properties in MCP responses

## Decisions

### 1. `IsLocal` field on `DisplayProperty`

Add `IsLocal bool` to the `DisplayProperty` struct. Set it in `BuildDisplayProperties()` by checking whether a property key exists in the schema property map or system property registry.

**Rationale:** A single boolean field is the simplest way to convey "this property is not schema-defined" to all consumers (TUI, CLI, future web). No need for a separate data structure or a new type.

**Alternative considered:** Returning two separate slices (schema props, local props). Rejected because it forces every consumer to handle two lists, and the current single-list model with metadata fields (`IsReverse`, `IsBacklink`) is well-established.

### 2. TUI separator rendering

Insert a `── Local Properties ──` separator line in the `propEditor.Render()` method before the first item where `IsLocal == true`. The separator is a display-only row, not a navigable item.

**Rationale:** Consistent with the existing "Properties" / "──────────" header pattern. A separator line is minimal yet clear.

### 3. Local properties as read-only

In `isPropertyEditable()`, return `false` when `dp.IsLocal == true`. This makes the cursor skip local properties during navigation, and they render with `dimStyle` (same as reverse relations and backlinks).

**Rationale:** Editing local properties requires type inference (what type is a local property?) and validation (what rules apply?). Deferring to a future issue keeps scope tight.

### 4. CLI section separation

In `cmd/show.go`, split the display properties loop into two groups: non-local (under "Properties") and local (under "Local Properties" with its own header). Only show the local section if there are local properties.

**Rationale:** Simple two-pass approach. No new abstractions needed.

## Risks / Trade-offs

- **[Minor] Property reclassification** — If a user adds a property to the schema that was previously local, it stops being `IsLocal` on next load. This is correct behavior (the property is now schema-managed), but could momentarily confuse users. → Mitigation: this is inherently correct and matches user intent.
- **[Minor] Separator rendering width** — The separator width should match the panel width. → Mitigation: use a fixed-width separator string matching the existing header pattern.
