## 1. TUI: Cell Navigation

- [ ] 1.1 Write BDD scenarios for table cell navigation (horizontal movement, boundary clamping, group header skip)
- [ ] 1.2 Add `colCursor` field to `viewMode` struct and implement left/right/h/l key handling
- [ ] 1.3 Implement Tab/Shift+Tab navigation between editable cells (skip read-only)
- [ ] 1.4 Remap Enter from open-object to edit-cell; add `o` key for opening object detail

## 2. TUI: Crosshair Rendering

- [ ] 2.1 Write BDD scenarios for crosshair highlighting (active cell, cursor row, cursor column header)
- [ ] 2.2 Modify `viewTable()` to render per-cell styles: active cell highlight, dim row highlight, dim column header tint
- [ ] 2.3 Add unit tests for crosshair style selection logic

## 3. TUI: Cell Edit State

- [ ] 3.1 Define `cellEdit` struct (editing bool, row/col index, property name/type, textinput model, picker state)
- [ ] 3.2 Implement `viewMode.activateCellEdit()` — determine property type and initialize appropriate widget
- [ ] 3.3 Implement `viewMode.cancelCellEdit()` — discard edit and return to navigation
- [ ] 3.4 Implement read-only cell detection (relation, created_at, updated_at) — Enter is no-op on these

## 4. TUI: Cell Editing — Text-based Types

- [ ] 4.1 Write BDD scenarios for string/number/date/datetime/url cell editing (activate, confirm, cancel, validation error)
- [ ] 4.2 Implement textinput cell editing for string, number, date, datetime, url property types
- [ ] 4.3 Implement validation via `core.ValidatePropertyValue()` with toast error on failure
- [ ] 4.4 Implement `viewMode.applyCellValue()` save pipeline (update object, call vault.SaveObject)

## 5. TUI: Cell Editing — Picker Types

- [ ] 5.1 Write BDD scenarios for select and multi-select cell editing (activate picker, navigate, confirm, cancel)
- [ ] 5.2 Implement select option picker for select property cells
- [ ] 5.3 Implement multi-select option picker with Space toggle for multi_select property cells

## 6. TUI: Cell Editing — Checkbox and NAME

- [ ] 6.1 Write BDD scenarios for checkbox toggle and NAME column editing
- [ ] 6.2 Implement checkbox direct toggle on Enter/Space in table cells
- [ ] 6.3 Implement NAME column editing via textinput (update object name and save)

## 7. TUI: Edge Cases and Integration

- [ ] 7.1 Handle cancel-on-external-file-change: cancel active cell edit when file watcher triggers object refresh
- [ ] 7.2 Disable cell editing when preview panel or view editor is open
- [ ] 7.3 Update help bar text to reflect new key bindings (h/l, Enter=edit, o=open, Tab)
- [ ] 7.4 Add unit tests for edge cases (empty schema, single NAME column, column cursor reset on view reload)
