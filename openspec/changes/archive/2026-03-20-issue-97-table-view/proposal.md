## Why

The current view mode was built as a "list" but actually renders a table with NAME + property columns. Issue #97 asks for a proper table view. Instead of adding another layout on top, we need to correctly separate the two: rename the existing rendering to `table`, create a true `list` layout (name + optional inline values), and add a `columns` field to ViewConfig so both layouts can specify which properties to display.

## What Changes

- Add `ViewLayoutTable` constant; rename existing column-based rendering to be the table layout
- Redefine `ViewLayoutList` as a true list: `name · val1 · val2` inline format
- Add `Columns` field to `ViewConfig` for explicit column/inline-value selection
- Default columns: list = empty (name only), table = all properties
- Implicit default view layout remains `list` (now shows the new simpler format)
- **BREAKING**: existing `layout: list` YAML files will show the new list behavior (not the old column table)
- Add Layout section to view editor so users can switch between list and table

## Capabilities

### New Capabilities

- `table-view-layout`: Table layout for views — columnar display with NAME + property columns, sort indicators, and configurable column selection
- `view-columns`: Configurable columns/inline-values for both list and table layouts via `columns` field in ViewConfig

### Modified Capabilities

- `tui-layout`: ViewConfig gains `Columns` field and `ViewLayoutTable` constant; list layout redefined as inline name + values

## Impact

- `core/view.go` — ViewConfig struct gains `Columns []string` field, new `ViewLayoutTable` constant
- `tui/view_mode.go` — `View()` dispatches on layout; extract current rendering as `viewTable()`, new `viewList()` for inline format; `viewColumns()` respects `Columns` config
- `tui/view_editor.go` — New Layout section and Columns section
- `.typemd/tui-state.yaml` — No changes (view mode persistence already handles type/view name)
- Existing view YAML files with `layout: list` — **BREAKING**: will now render as simple list, not table
