## Why

Views can be created in the type editor with default settings, but editing filter rules, sort order, and group_by requires manually editing YAML files. Users need an inline TUI editor to configure views without leaving the terminal.

Additionally, the current `GroupBy` field is a single `string`, which limits grouping to one property. SQL `GROUP BY` supports multiple columns, and multi-level grouping is a natural user expectation. This should be addressed as a **breaking change** now before the view system matures further.

## What Changes

- **BREAKING**: Change `ViewConfig.GroupBy` from `string` to `[]GroupRule` struct slice, enabling multi-level grouping
- Add `GroupRule` struct to core with a `Property` field
- Add a `viewEditor` TUI sub-model for editing ViewConfig inline (filter rules, sort rules, group rules)
- Integrate view editor into view mode as a right-side split panel (reusing existing preview split mechanism)
- Enter view editor via `e` key in view mode; exit with `Esc`
- Update `buildGroups()` in view mode to support multi-level grouping
- Migrate existing single-string `group_by` YAML to `[]GroupRule` format on load

## Capabilities

### New Capabilities
- `view-editor`: TUI sub-model for inline editing of ViewConfig (filter, sort, group rules) within view mode's split panel
- `group-rule`: GroupRule struct and multi-level grouping support replacing single-string GroupBy

### Modified Capabilities
- `structured-query-filter`: No requirement changes (FilterRule struct unchanged)

## Impact

- **core/view.go**: `ViewConfig.GroupBy` type change (`string` → `[]GroupRule`), add `GroupRule` struct, YAML migration on load
- **core/view_test.go**: Update tests for new GroupBy type
- **core/features/view_config.feature**: Update BDD scenarios for GroupRule
- **tui/view_mode.go**: `buildGroups()` rewritten for multi-level grouping, `e` key handler to open editor, split panel integration
- **tui/view_editor.go** (new): `viewEditor` sub-model with filter/sort/group rule editing
- **tui/app.go + app_render.go**: Wire view editor into view mode rendering and update cycle
