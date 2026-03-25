## Why

The table view in view mode is currently read-only — users must navigate into the object detail view to edit properties. This creates unnecessary friction for bulk editing workflows. Adding inline cell editing transforms the table view into an Excel-like grid where users can navigate cells and edit values directly, completing the "Edit Everything Inline" vision for v0.7.0.

## What Changes

- **Cell navigation**: Add horizontal cursor (column index) to table view, enabling four-directional cell movement (arrow keys) with crosshair highlighting
- **Inline cell editing**: Press Enter on a cell to edit its value in-place, with type-appropriate widgets (textinput, select picker, multi-select picker, checkbox toggle)
- **NAME column editing**: Allow inline editing of object names directly in the table
- **Key remapping**: Remap Enter from "open object" to "edit cell"; add `o` key for opening object detail view
- **Auto-save**: Edits save immediately on confirm (Enter); Esc cancels
- **Validation**: Reuse `core.ValidatePropertyValue` with toast notifications for errors

## Capabilities

### New Capabilities
- `table-cell-navigation`: Four-directional cell navigation in table view with crosshair highlighting (row + column focus indicators)
- `table-cell-editing`: Inline editing of property values and object names in table view cells, with type-appropriate widgets, validation, and auto-save

### Modified Capabilities
- `table-view-layout`: Table layout rendering changes to support per-cell highlighting and crosshair visual indicators

## Impact

- `tui/view_mode.go` — Add `colCursor`, `cellEdit` to `viewMode` struct; modify `Update()` for cell navigation and editing; modify `viewTable()` for crosshair rendering
- `tui/view_mode.go` — Remap `enter` from open-object to edit-cell; add `o` key binding
- `core/type_schema_validate.go` — Reuse existing `ValidatePropertyValue` (no changes needed)
- `tui/prop_editor.go` — Reuse save pipeline pattern (`applyPropertyValue`) for table cell saves
- `tui/widget/toast.go` — Reuse existing toast for validation errors (no changes needed)
