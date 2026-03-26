## Context

The TUI's view mode (`tui/view_mode.go`) renders objects in a full-width layout with two display modes: list and table. The table layout shows objects in a columnar format (NAME + property columns) with row-only navigation (up/down). The recently merged #315 added inline property editing to the object detail view via `propEditor`, establishing patterns for type-aware editing widgets and the save pipeline.

The existing `propEditor` is tightly coupled to the main `model` struct — its `applyPropertyValue` function accesses `m.selected`, `m.dirty`, `m.saveObject()`, and `m.displayProps`. This cannot be directly reused from `viewMode`, which manages its own object list independently.

## Goals / Non-Goals

**Goals:**
- Enable four-directional cell navigation in table view (row + column)
- Provide inline editing for all non-relation property types and the NAME column
- Crosshair visual feedback (row + column highlights) for easy orientation
- Reuse validation logic (`core.ValidatePropertyValue`) and editing patterns from `propEditor`
- Auto-save on edit confirmation

**Non-Goals:**
- Batch operations (multi-cell select, drag-fill, copy-paste)
- Relation property editing (depends on #88)
- Cell editing in list layout (only table layout)
- Editing while preview panel or view editor is open

## Decisions

### 1. Cell editing state lives inside `viewMode`

Add a `cellEdit` struct to `viewMode` rather than embedding a `propEditor`. The `propEditor` is designed for the vertical properties panel with different navigation semantics. A dedicated `cellEdit` struct is simpler and avoids fighting the `propEditor`'s model-coupled save pipeline.

**Alternative considered:** Extracting a shared editing interface from `propEditor` — rejected because the edit lifecycle is different (cell editing operates on a specific object from the objects list, not `m.selected`), and the abstraction cost exceeds the code savings.

### 2. Column cursor as index into `viewColumns()`

`colCursor = 0` represents the NAME column. `colCursor = 1..N` maps to `viewColumns()[colCursor-1]`. This aligns with the visual layout and makes bounds checking trivial.

### 3. Save pipeline adapted for viewMode

Create a `viewMode.applyCellValue(obj, key, value)` method that:
1. Sets `obj.Properties[key] = value` (or updates Name for NAME column)
2. Calls `vault.SaveObject(obj)`
3. Does NOT rebuild display props or trigger full refresh — the in-memory object is already updated

This is intentionally simpler than `applyPropertyValue` since viewMode doesn't have the `dirty`/`displayProps`/`refreshPropEditor` lifecycle.

### 4. Crosshair rendering via per-cell style selection

In `viewTable()`, each cell checks: (a) is this the cursor row? (b) is this the cursor column? Apply styles:
- Active cell (both): strong foreground + background
- Cursor row only: dim background highlight
- Cursor column only: dim foreground tint on column header
- Neither: default style

This avoids full-row highlighting that would obscure which column is selected.

### 5. Key remapping

| Key | Current behavior | New behavior |
|-----|-----------------|--------------|
| Enter | Open object detail | Edit current cell |
| Space | Open object detail | Edit current cell (or toggle checkbox) |
| `o` | (unused) | Open object detail |
| Left/Right or `h`/`l` | (unused) | Move column cursor |
| Tab | (unused in view mode) | Save + move to next editable cell |
| Esc (during edit) | Exit view mode | Cancel edit |

### 6. Read-only cell detection

A cell is read-only if:
- Property type is `relation` (deferred to #88)
- Property is `created_at` or `updated_at` (immutable system properties)
- Column is a reverse relation or backlink

Tab navigation and Enter-to-edit skip read-only cells. The cursor can land on them (for visual reference) but Enter is a no-op.

## Risks / Trade-offs

- **[Column width constraints]** → Inline textinput must fit within the 12-char column width. Use the existing `colW` constant. For values longer than 12 chars, the textinput scrolls internally (charmbracelet/bubbles textinput supports this natively).
- **[Group header rows]** → Left/right navigation is a no-op on group headers. Only up/down traversal crosses headers.
- **[Concurrent file changes]** → If a file changes during editing (detected by file watcher), cancel the active cell edit and show a toast warning. This is simpler than conflict resolution and matches user expectations.
- **[Preview/editor panel interaction]** → Cell editing is disabled when preview (`p`) or view editor (`e`) is open. Enter/arrows retain their current meaning in those contexts.
