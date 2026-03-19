## Why

`tui/view_mode.go` has a `formatPropValue()` function that duplicates value formatting logic already in `core.DisplayProperty.Format()`. The TUI function only handles `string`, `bool`, and a `%v` fallback — missing dates, multi_select, relations, and other type-aware formatting. This creates inconsistent display between view mode tables and other TUI panels (detail, properties) that use `DisplayProperty.Format()`.

## What Changes

- Add `FormatValue() string` method to `DisplayProperty` that returns the formatted value without key prefix
- Refactor `Format()` to delegate to `FormatValue()` (key + ": " + value)
- Unify checkbox/bool formatting to `✓` (true) and empty string (false) across all views
- Replace `tui/view_mode.go` `formatPropValue()` with `DisplayProperty.FormatValue()` calls
- View mode gains proper formatting for date, datetime, multi_select, and relation properties

## Capabilities

### New Capabilities

- `display-property-format`: Covers the DisplayProperty formatting pipeline — FormatValue() for value-only output, Format() for key-value output, and type-aware formatting rules (checkbox, date, datetime, multi_select, relation, backlink, reverse relation)

### Modified Capabilities

## Impact

- `core/display.go` — add `FormatValue()`, refactor `Format()`, change checkbox format from `[x]`/`[ ]` to `✓`/empty
- `core/display_test.go` — update checkbox test expectations, add FormatValue tests
- `tui/view_mode.go` — remove `formatPropValue()`, use `DisplayProperty` with `FormatValue()`
- `tui/detail.go` — checkbox display changes from `[x]`/`[ ]` to `✓`/empty (via Format() change)
- `cmd/show.go` — same checkbox display change (via Format() change)
