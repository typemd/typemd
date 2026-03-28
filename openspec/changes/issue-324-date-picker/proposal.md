## Why

The `date` property currently uses a plain textinput for editing, requiring users to type the full `YYYY-MM-DD` string manually. This is slow, error-prone, and provides no visual feedback (e.g., day of week). A dedicated date picker with segmented input and inline calendar will make date editing intuitive and efficient.

## What Changes

- Add a shared `dateEdit` sub-model with two editing modes: segmented input and inline calendar
- Segmented input: three editable segments (YYYY-MM-DD) with arrow key increment/decrement, direct digit input, and live day-of-week feedback
- Inline calendar: mini month grid with keyboard navigation, month switching, and today marker/jump
- Toggle between modes with `c`
- Route `date` property type through the new date editor in both properties panel (`propEditor`) and table view cell editing (`cellEdit`)
- Help bar shows `[DATE]` during segmented input, `[CAL]` during calendar

## Capabilities

### New Capabilities
- `date-picker`: Dedicated date editing widget with segmented input mode and inline calendar mode, shared between properties panel and table view

### Modified Capabilities
- `property-editing`: Date properties route through the new date picker instead of generic textinput
- `table-view`: Date cell editing routes through the new date picker instead of generic textinput

## Impact

- **TUI**: New file `tui/date_edit.go` for the shared `dateEdit` sub-model; modifications to `tui/prop_editor.go`, `tui/prop_editor_update.go`, `tui/view_mode_cell_edit.go` for routing
- **No core changes**: Validation still uses existing `core.ValidatePropertyValue`
- **No API/dependency changes**: Pure TUI-layer addition
