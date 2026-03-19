## Context

The view system already supports a `Layout` field on `ViewConfig` with only `ViewLayoutList` defined. The current "list" rendering (`viewMode.View()`) displays objects in a columnar table format (NAME + property columns with headers and separators). This is actually a table layout, not a list. Issue #97 asks for a table view, which already exists — we need to relabel it and create a true list layout.

## Goals / Non-Goals

**Goals:**
- Add `ViewLayoutTable` constant for the existing columnar display
- Redefine `ViewLayoutList` as a true list (name + inline property values)
- Add `Columns` field to `ViewConfig` for configurable column/inline-value selection
- Add Layout and Columns sections to the view editor
- Both layouts respect filter, sort, and group_by settings

**Non-Goals:**
- Horizontal scrolling for table (future enhancement)
- Column resizing or drag-to-reorder (future)
- Kanban, gallery, or calendar layouts (future)

## Decisions

### 1. Two layout constants, same underlying view mode

Both `list` and `table` are rendered by the same `viewMode` sub-model. The `View()` method dispatches based on `vm.view.Layout`:
- `ViewLayoutList` → `vm.viewList()` (new)
- `ViewLayoutTable` → `vm.viewTable()` (extracted from current `View()`)

**Alternative considered:** Separate sub-models for each layout. Rejected because both layouts share the same data (objects, groups, cursor, scroll) and interaction model (up/down navigation, Enter to open, Esc to exit). Only the rendering differs.

### 2. Columns field on ViewConfig

```yaml
columns:
  - status
  - rating
```

`Columns []string` in `ViewConfig`. When empty:
- **List layout:** shows only name (no inline values)
- **Table layout:** shows all schema properties (current behavior via `viewColumns()`)

When non-empty, both layouts use exactly the specified columns in order.

### 3. List layout rendering

Each row: `emoji name · val1 · val2` where values come from configured columns.

```
📚 Clean Code · reading · ⭐⭐⭐⭐
📚 DDIA · done · ⭐⭐⭐⭐⭐
📚 The Pragmatic Programmer · to-read
```

- Emoji comes from the type schema (if defined)
- Name is the object's display name (via `GetName()`)
- Column values separated by ` · ` (middle dot)
- Empty values are omitted (no trailing dots)
- Group headers use the same `── Label ──` format as table layout
- Cursor highlighting applies to the entire row

### 4. Table layout is the existing rendering (renamed)

Extract current `View()` body into `viewTable()` with no behavioral changes. The existing column selection logic (`viewColumns()`) becomes the default when `Columns` is empty.

When `Columns` is non-empty, `viewColumns()` returns exactly those columns (in order), skipping its pinned/unpinned auto-detection logic.

### 5. Implicit default view uses list layout

`DefaultView()` already returns `ViewLayoutList`. Since list is now the simple name-only format, new types start with a clean list view. Users can switch to table via the view editor.

### 6. View editor gains Layout and Columns sections

The view editor currently has Filter, Sort, and GroupBy sections. Add:
- **Layout** section: single-choice picker (list / table)
- **Columns** section: multi-select property picker to add/remove columns; reorder via move mode

Layout section goes first (it determines the visual context for other settings). Columns section goes after Layout.

### 7. Sort indicators in table layout

Table column headers show `↑`/`↓` next to sorted columns:

```
NAME          STATUS ↑   RATING
```

This is a display-only enhancement in `viewTable()`. List layout does not show sort indicators.

## Risks / Trade-offs

- **[BREAKING] Existing `layout: list` YAML** → Will now render as simple list instead of table. Users who want the old behavior need to change to `layout: table`. This is acceptable because: (a) the view system is still in early adoption, (b) the fix is a one-word YAML edit, (c) the rename makes the naming correct.
- **[Risk] Columns field naming** → "columns" makes sense for table but is slightly odd for list (where they're inline values). Using the same field for both keeps the API simple. The view editor can label it "Inline Values" for list and "Columns" for table.
